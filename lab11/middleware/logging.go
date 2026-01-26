package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware logs request completion with duration and status code.
// Uses structured logging with context including request ID.
// TODO: Part 2 - Implement logging middleware
func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Record start time
		
		// TODO: Wrap response writer to capture status code
		
		// TODO: Process request
		
		// TODO: Calculate duration
		// TODO: Get request ID from context
		
		// TODO: Log request completion with:
		// - request_id
		// - method
		// - path
		// - status
		// - duration
		
		next.ServeHTTP(w, r)
	})
}
