package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"lab11/middleware"
	"lab11/observability"
)

// Server handles HTTP requests with full observability.
type Server struct {
	logger  *slog.Logger
	metrics *observability.Metrics
	mux     *http.ServeMux
}

// NewServer creates a server with observability middleware.
func NewServer(logger *slog.Logger) *Server {
	metrics, err := observability.NewMetrics()
	if err != nil {
		logger.Error("Failed to initialize metrics", slog.String("error", err.Error()))
		panic(err)
	}

	s := &Server{
		logger:  logger,
		metrics: metrics,
		mux:     http.NewServeMux(),
	}

	// Register routes
	s.mux.HandleFunc("/orders", s.handleOrders)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/metrics", s.handleMetrics)

	return s
}

// ServeHTTP implements http.Handler with middleware stack.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Build middleware chain from innermost to outermost
	var handler http.Handler = s.mux
	handler = middleware.MetricsMiddleware(s.metrics, handler)
	handler = middleware.TracingMiddleware(handler)
	handler = middleware.LoggingMiddleware(s.logger, handler)
	handler = middleware.RequestIDMiddleware(s.logger, handler)

	handler.ServeHTTP(w, r)
}

// Order represents a simple order for demonstration.
type Order struct {
	ID         string  `json:"id"`
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

// handleOrders handles order creation with full observability.
func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a child span for order processing
	ctx := r.Context()
	tracer := otel.Tracer("order-service")
	ctx, span := tracer.Start(ctx, "ProcessOrder")
	defer span.End()

	// Decode request
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		// Log error with structured logging
		requestID := middleware.GetRequestID(ctx)
		s.logger.Error("Failed to decode order JSON",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)

		// Record error in span
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid JSON")

		// Record failed order metric
		s.metrics.RecordOrderFailed(ctx, "invalid_json")

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Add order details to span attributes
	span.SetAttributes(
		attribute.String("order.id", order.ID),
		attribute.String("order.customer_id", order.CustomerID),
		attribute.Float64("order.amount", order.Amount),
	)

	// Validate order
	if order.ID == "" || order.CustomerID == "" || order.Amount <= 0 {
		// Log validation failure with order details
		requestID := middleware.GetRequestID(ctx)
		s.logger.Warn("Order validation failed",
			slog.String("request_id", requestID),
			slog.String("order_id", order.ID),
			slog.String("customer_id", order.CustomerID),
			slog.Float64("amount", order.Amount),
		)

		// Set span status to error
		span.SetStatus(codes.Error, "Validation failed")

		// Record failed order metric with reason
		s.metrics.RecordOrderFailed(ctx, "validation_failed")

		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	// Set default status
	order.Status = "pending"

	// Log business event
	requestID := middleware.GetRequestID(ctx)
	s.logger.Info("Order created",
		slog.String("request_id", requestID),
		slog.String("order_id", order.ID),
		slog.String("customer_id", order.CustomerID),
		slog.Float64("amount", order.Amount),
		slog.String("status", order.Status),
	)

	// Record success metrics
	s.metrics.RecordOrderCreated(ctx)
	s.metrics.RecordOrderValue(ctx, order.Amount)

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// HealthResponse represents health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// handleHealth returns service health.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "healthy"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMetrics exposes Prometheus metrics for scraping.
// This endpoint is scraped by Prometheus to collect metrics.
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
