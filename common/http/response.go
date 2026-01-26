package http

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// WriteError writes a JSON error response with the given status code.
// Provides consistent error response format across all endpoints.
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	WriteJSON(w, ErrorResponse{Error: message}, statusCode)
}
