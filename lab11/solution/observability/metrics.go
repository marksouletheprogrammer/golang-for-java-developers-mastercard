package observability

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Metrics holds all application metrics.
// Using Prometheus metrics for direct export to /metrics endpoint.
type Metrics struct {
	// Prometheus metrics for /metrics endpoint
	HTTPRequests    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	OrdersCreated   prometheus.Counter
	OrdersFailed    *prometheus.CounterVec
	OrderValue      prometheus.Histogram
	
	// OpenTelemetry metrics for tracing correlation
	otelHTTPRequests   metric.Int64Counter
	otelRequestDuration metric.Float64Histogram
}

// NewMetrics creates and registers all application metrics.
// Creates both Prometheus metrics (for /metrics endpoint) and OpenTelemetry metrics (for tracing).
func NewMetrics() (*Metrics, error) {
	// Create Prometheus metrics
	httpRequests := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)
	
	requestDuration := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	
	ordersCreated := promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total number of orders created",
		},
	)
	
	ordersFailed := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_failed_total",
			Help: "Total number of failed orders",
		},
		[]string{"reason"},
	)
	
	orderValue := promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "order_value_dollars",
			Help:    "Order value in dollars",
			Buckets: []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
		},
	)
	
	// Create OpenTelemetry metrics for tracing correlation
	meter := otel.Meter("order-service")
	
	otelHTTPRequests, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, err
	}
	
	otelRequestDuration, err := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
	)
	if err != nil {
		return nil, err
	}
	
	return &Metrics{
		HTTPRequests:        httpRequests,
		RequestDuration:     requestDuration,
		OrdersCreated:       ordersCreated,
		OrdersFailed:        ordersFailed,
		OrderValue:          orderValue,
		otelHTTPRequests:    otelHTTPRequests,
		otelRequestDuration: otelRequestDuration,
	}, nil
}

// RecordHTTPRequest records an HTTP request with labels.
// Records to both Prometheus (for /metrics) and OpenTelemetry (for traces).
func (m *Metrics) RecordHTTPRequest(ctx context.Context, method, path string, statusCode int) {
	// Prometheus metric - convert status code to string like "2xx", "4xx", etc.
	statusClass := fmt.Sprintf("%dxx", statusCode/100)
	m.HTTPRequests.WithLabelValues(method, path, statusClass).Inc()
	
	// OpenTelemetry metric for trace correlation
	m.otelHTTPRequests.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.Int("status_code", statusCode),
		),
	)
}

// RecordRequestDuration records request duration in seconds.
func (m *Metrics) RecordRequestDuration(ctx context.Context, durationSec float64, method, path string) {
	// Prometheus metric
	m.RequestDuration.WithLabelValues(method, path).Observe(durationSec)
	
	// OpenTelemetry metric for trace correlation
	m.otelRequestDuration.Record(ctx, durationSec,
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.String("path", path),
		),
	)
}

// RecordOrderCreated records a successfully created order.
func (m *Metrics) RecordOrderCreated(ctx context.Context) {
	m.OrdersCreated.Inc()
}

// RecordOrderFailed records a failed order creation.
func (m *Metrics) RecordOrderFailed(ctx context.Context, reason string) {
	m.OrdersFailed.WithLabelValues(reason).Inc()
}

// RecordOrderValue records the value of an order.
func (m *Metrics) RecordOrderValue(ctx context.Context, value float64) {
	m.OrderValue.Observe(value)
}
