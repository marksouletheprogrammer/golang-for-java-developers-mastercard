package middleware

import (
	"net/http"
	"time"

	"lab11/observability"
)

// MetricsMiddleware records metrics for each HTTP request.
func MetricsMiddleware(metrics *observability.Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap response writer
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Process request
		next.ServeHTTP(wrapped, r)
		
		// Record metrics
		duration := time.Since(start).Seconds()
		metrics.RecordHTTPRequest(r.Context(), r.Method, r.URL.Path, wrapped.statusCode)
		metrics.RecordRequestDuration(r.Context(), duration, r.Method, r.URL.Path)
	})
}
