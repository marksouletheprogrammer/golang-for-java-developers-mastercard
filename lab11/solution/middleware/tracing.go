package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware creates a span for each HTTP request.
// Spans record the execution flow and timing of operations.
func TracingMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("http-server")
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start a new span for this request
		ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.host", r.Host),
			),
		)
		defer span.End()
		
		// Add request ID to span if available
		if requestID := GetRequestID(ctx); requestID != "" {
			span.SetAttributes(attribute.String("request.id", requestID))
		}
		
		// Wrap response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Process request with span context
		next.ServeHTTP(wrapped, r.WithContext(ctx))
		
		// Record response status in span
		span.SetAttributes(attribute.Int("http.status_code", wrapped.statusCode))
		
		// Mark span as error if status >= 400
		if wrapped.statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
		}
	})
}
