package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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
// TODO: Part 1, 6 - Implement order handler with logging, tracing, and error recording
func (s *Server) handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Part 6 - Create a child span for order processing
	// TODO: Use otel.Tracer("order-service") and Start() a span named "ProcessOrder"
	// TODO: Defer span.End()

	// Decode request
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		// TODO: Part 1 - Log error with structured logging
		// TODO: Include request_id, error message

		// TODO: Part 6 - Record error in span
		// TODO: Use span.RecordError() and span.SetStatus()

		// TODO: Part 4 - Record failed order metric

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// TODO: Part 6 - Add order details to span attributes
	// TODO: order.id, order.customer_id, order.amount

	// Validate order
	if order.ID == "" || order.CustomerID == "" || order.Amount <= 0 {
		// TODO: Part 1 - Log validation failure with order details

		// TODO: Part 6 - Set span status to error

		// TODO: Part 4 - Record failed order metric with reason "validation_failed"

		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	// Set default status
	order.Status = "pending"

	// TODO: Part 1 - Log business event "Order created"
	// TODO: Include request_id, order_id, customer_id, amount, status

	// TODO: Part 4 - Record success metrics
	// TODO: RecordOrderCreated() and RecordOrderValue()

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
// TODO: Part 5 - Implement metrics endpoint
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: Use promhttp.Handler() to expose Prometheus metrics
	promhttp.Handler().ServeHTTP(w, r)
}
