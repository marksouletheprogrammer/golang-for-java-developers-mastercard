package observability

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
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
	otelHTTPRequests    metric.Int64Counter
	otelRequestDuration metric.Float64Histogram
}

// NewMetrics creates and registers all application metrics.
// Creates both Prometheus metrics (for /metrics endpoint) and OpenTelemetry metrics (for tracing).
// TODO: Part 4 - Implement metrics creation
func NewMetrics() (*Metrics, error) {
	// TODO: Create Prometheus counters using promauto.NewCounterVec()
	// TODO: Create HTTPRequests counter with labels: method, path, status_code
	// TODO: Create OrdersCreated counter (no labels)
	// TODO: Create OrdersFailed counter with label: reason

	// TODO: Create Prometheus histograms using promauto.NewHistogramVec()
	// TODO: Create RequestDuration histogram with labels: method, path
	// TODO: Create OrderValue histogram with custom buckets for dollar amounts

	// TODO: Create OpenTelemetry meter using otel.Meter("order-service")
	// TODO: Create OpenTelemetry counters and histograms for trace correlation

	return &Metrics{
		// TODO: Assign all metrics
	}, nil
}

// RecordHTTPRequest records an HTTP request with labels.
// Records to both Prometheus (for /metrics) and OpenTelemetry (for traces).
// TODO: Part 4 - Implement HTTP request recording
func (m *Metrics) RecordHTTPRequest(ctx context.Context, method, path string, statusCode int) {
	// TODO: Convert status code to status class (e.g., "2xx", "4xx")
	// TODO: Increment Prometheus counter with labels
	// TODO: Add to OpenTelemetry counter with attributes
}

// RecordRequestDuration records request duration in seconds.
// TODO: Part 4 - Implement request duration recording
func (m *Metrics) RecordRequestDuration(ctx context.Context, durationSec float64, method, path string) {
	// TODO: Observe duration in Prometheus histogram
	// TODO: Record duration in OpenTelemetry histogram
}

// RecordOrderCreated records a successfully created order.
// TODO: Part 4 - Implement order created recording
func (m *Metrics) RecordOrderCreated(ctx context.Context) {
	// TODO: Increment OrdersCreated counter
}

// RecordOrderFailed records a failed order creation.
// TODO: Part 4 - Implement order failed recording
func (m *Metrics) RecordOrderFailed(ctx context.Context, reason string) {
	// TODO: Increment OrdersFailed counter with reason label
}

// RecordOrderValue records the value of an order.
// TODO: Part 4 - Implement order value recording
func (m *Metrics) RecordOrderValue(ctx context.Context, value float64) {
	// TODO: Observe value in OrderValue histogram
}
