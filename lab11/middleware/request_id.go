package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

const RequestIDKey ContextKey = "request_id"

// RequestIDMiddleware generates a unique request ID for each request.
// The request ID is added to the context and logged with every log entry.
// This enables correlation of all logs for a single request.
// TODO: Part 2 - Implement request ID middleware
func RequestIDMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Generate unique request ID using uuid.New().String()
		
		// TODO: Add request ID to context using context.WithValue()
		// TODO: Update request with new context using r.WithContext()
		
		// TODO: Add request ID to response headers (X-Request-ID)
		
		// TODO: Log request start with request ID, method, path, remote_addr
		
		// TODO: Pass to next handler
		next.ServeHTTP(w, r)
	})
}

// GetRequestID extracts the request ID from context.
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
