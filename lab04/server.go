package main

import (
	"net/http"
)

// Server holds the HTTP server dependencies.
type Server struct {
	store *MerchantStore
	mux   *http.ServeMux
}

// NewServer creates a new HTTP server with all routes configured.
func NewServer(store *MerchantStore) *Server {
	// TODO: Part 1 - Implement NewServer
	// Create Server with store and new ServeMux
	// Register routes for /merchants and /merchants/
	// Return server
	return nil
}

// ServeHTTP implements http.Handler interface.
// Wraps the mux with logging middleware.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Part 1 - Implement logging middleware
	// Record start time
	// Create responseWriter wrapper to capture status code
	// Call s.mux.ServeHTTP
	// Log the request after completion with method, path, status, duration
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before delegating to the wrapped writer.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// handleMerchants handles both GET and POST to /merchants endpoint.
// GET returns all merchants, POST creates a new merchant.
func (s *Server) handleMerchants(w http.ResponseWriter, r *http.Request) {
	// TODO: Part 3 and Part 5 - Implement method routing
	// Use switch on r.Method
	// Call handleGetMerchants for GET
	// Call handleCreateMerchant for POST
	// Return 405 Method Not Allowed for other methods
}

// handleGetMerchants returns all merchants as JSON.
func (s *Server) handleGetMerchants(w http.ResponseWriter, r *http.Request) {
	// TODO: Part 3 - Implement GET /merchants
	// Get all merchants from store
	// Write JSON response with 200 status
}

// handleCreateMerchant creates a new merchant from JSON request body.
func (s *Server) handleCreateMerchant(w http.ResponseWriter, r *http.Request) {
	// TODO: Part 5 - Implement POST /merchants
	// Decode JSON from request body
	// Validate merchant
	// Create merchant in store
	// Set Location header
	// Return 201 with merchant JSON
}

// handleMerchantByID handles GET /merchants/{id}.
// Extracts the ID from the URL path and retrieves the merchant.
func (s *Server) handleMerchantByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Part 4 - Implement GET /merchants/{id}
	// Check method is GET
	// Extract ID from URL path
	// Get merchant from store
	// Return 404 if not found
	// Return 200 with merchant JSON if found
}

// StartServer starts the HTTP server on the specified port.
func StartServer(port string, store *MerchantStore) error {
	// TODO: Part 1 - Implement StartServer
	// Create server with NewServer
	// Print startup message with available endpoints
	// Call http.ListenAndServe
	return nil
}
