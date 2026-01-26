package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	commonhttp "golang-for-java-developers-training/common/http"
)

// Server holds the HTTP server dependencies.
type Server struct {
	store *MerchantStore
	mux   *http.ServeMux
}

// NewServer creates a new HTTP server with all routes configured.
func NewServer(store *MerchantStore) *Server {
	s := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}

	// Register routes
	s.mux.HandleFunc("/merchants", s.handleMerchants)
	s.mux.HandleFunc("/merchants/", s.handleMerchantByID)

	// Product enrichment endpoint (Lab 5)
	s.mux.HandleFunc("/products/", s.handleProductEnriched)

	return s
}

// ServeHTTP implements http.Handler interface.
// Wraps the mux with logging middleware.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Logging middleware - logs every request
	start := time.Now()

	// Create a response writer wrapper to capture status code
	wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

	// Call the actual handler
	s.mux.ServeHTTP(wrapped, r)

	// Log after request completes
	duration := time.Since(start)
	log.Printf("%s %s - %d (%v)", r.Method, r.URL.Path, wrapped.statusCode, duration)
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
	switch r.Method {
	case http.MethodGet:
		s.handleGetMerchants(w, r)
	case http.MethodPost:
		s.handleCreateMerchant(w, r)
	default:
		commonhttp.WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetMerchants returns all merchants as JSON.
func (s *Server) handleGetMerchants(w http.ResponseWriter, r *http.Request) {
	merchants := s.store.GetAll()
	commonhttp.WriteJSON(w, merchants, http.StatusOK)
}

// handleCreateMerchant creates a new merchant from JSON request body.
func (s *Server) handleCreateMerchant(w http.ResponseWriter, r *http.Request) {
	// Decode JSON from request body
	var merchant Merchant
	if err := json.NewDecoder(r.Body).Decode(&merchant); err != nil {
		commonhttp.WriteError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate merchant
	if err := ValidateMerchant(&merchant); err != nil {
		commonhttp.WriteError(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create merchant
	if err := s.store.Create(&merchant); err != nil {
		commonhttp.WriteError(w, err.Error(), http.StatusConflict)
		return
	}

	// Set Location header to point to the new resource
	w.Header().Set("Location", "/merchants/"+merchant.ID)
	commonhttp.WriteJSON(w, merchant, http.StatusCreated)
}

// handleMerchantByID handles GET /merchants/{id}.
// Extracts the ID from the URL path and retrieves the merchant.
func (s *Server) handleMerchantByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		commonhttp.WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path: /merchants/{id}
	id := strings.TrimPrefix(r.URL.Path, "/merchants/")
	if id == "" || id == "/" {
		commonhttp.WriteError(w, "Merchant ID is required", http.StatusBadRequest)
		return
	}

	// Handle trailing slash
	id = strings.TrimSuffix(id, "/")

	// Get merchant
	merchant, err := s.store.GetByID(id)
	if err != nil {
		commonhttp.WriteError(w, "Merchant not found", http.StatusNotFound)
		return
	}

	commonhttp.WriteJSON(w, merchant, http.StatusOK)
}

// handleProductEnriched handles GET /products/{sku}/enriched.
// Enriches a single product with inventory, pricing, and review data using concurrent API calls.
func (s *Server) handleProductEnriched(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		commonhttp.WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract product path: /products/{sku}/enriched
	path := strings.TrimPrefix(r.URL.Path, "/products/")

	// Parse SKU and endpoint
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "enriched" {
		commonhttp.WriteError(w, "Invalid endpoint. Use /products/{sku}/enriched", http.StatusBadRequest)
		return
	}

	sku := parts[0]
	if sku == "" {
		commonhttp.WriteError(w, "SKU is required", http.StatusBadRequest)
		return
	}

	// For demo purposes, create a product with dummy base data
	// In a real system, this would come from a product database
	product := Product{
		SKU:       sku,
		Name:      "Product " + sku,
		BasePrice: 99.99,
	}

	// Create real external API client
	client := &RealExternalAPIClient{}

	// Use fan-out pattern to enrich with concurrent API calls
	enriched := EnrichSingleProductFanOut(client, product)

	commonhttp.WriteJSON(w, enriched, http.StatusOK)
}

// StartServer starts the HTTP server on the specified port.
func StartServer(port string, store *MerchantStore) error {
	server := NewServer(store)

	addr := ":" + port
	fmt.Printf("Starting merchant API server on http://localhost%s\n", addr)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET    /merchants              - List all merchants")
	fmt.Println("  GET    /merchants/{id}         - Get merchant by ID")
	fmt.Println("  POST   /merchants              - Create new merchant")
	fmt.Println("  GET    /products/{sku}/enriched - Get enriched product data")

	return http.ListenAndServe(addr, server)
}
