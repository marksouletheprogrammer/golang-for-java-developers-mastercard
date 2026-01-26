package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
)

// TracingMiddleware creates a span for each HTTP request.
// Spans record the execution flow and timing of operations.
// TODO: Part 6 - Implement tracing middleware
func TracingMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("http-server")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Start a new span for this request
		// TODO: Use tracer.Start() with span name like "GET /orders"
		// TODO: Set span kind to trace.SpanKindServer
		// TODO: Add attributes: http.method, http.url, http.host
		// TODO: Defer span.End()

		// TODO: Add request ID to span if available

		// TODO: Wrap response writer to capture status code

		// TODO: Process request with span context

		// TODO: Part 6 - Record response status in span
		// TODO: Mark span as error if status >= 400

		next.ServeHTTP(w, r)
	})
}
