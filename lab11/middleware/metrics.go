package middleware

import (
	"net/http"

	"lab11/observability"
)

// MetricsMiddleware records metrics for each HTTP request.
// TODO: Part 4 - Implement metrics middleware
func MetricsMiddleware(metrics *observability.Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Record start time

		// TODO: Wrap response writer to capture status code

		// TODO: Process request

		// TODO: Calculate duration in seconds
		// TODO: Record HTTP request metrics with method, path, status code
		// TODO: Record request duration with method, path

		next.ServeHTTP(w, r)
	})
}
