package main

import (
	"log"
)

func main() {
	// Create merchant store with sample data
	store := NewMerchantStore()
	
	// Start HTTP server on port 8080
	if err := StartServer("8080", store); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
