package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	// Seed random number generator for consistent results
	rand.Seed(time.Now().UnixNano())

	// Run performance comparison demo
	runPerformanceDemo()

	fmt.Println("\n=== Starting API Server ===")
	fmt.Println()

	// Create merchant store with sample data
	store := NewMerchantStore()
	initSampleMerchants(store)

	// Start HTTP server on port 8080
	if err := StartServer("8080", store); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// runPerformanceDemo runs the concurrent enrichment performance comparison.
func runPerformanceDemo() {
	// Create real external API client for production use
	client := &RealExternalAPIClient{}

	// Create sample products
	products := []Product{
		{SKU: "PROD-001", Name: "Laptop", BasePrice: 999.99},
		{SKU: "PROD-002", Name: "Mouse", BasePrice: 29.99},
		{SKU: "PROD-003", Name: "Keyboard", BasePrice: 79.99},
		{SKU: "PROD-004", Name: "Monitor", BasePrice: 299.99},
		{SKU: "PROD-005", Name: "Webcam", BasePrice: 89.99},
		{SKU: "PROD-006", Name: "Headset", BasePrice: 149.99},
		{SKU: "PROD-007", Name: "USB Cable", BasePrice: 12.99},
		{SKU: "PROD-008", Name: "Desk Lamp", BasePrice: 39.99},
		{SKU: "PROD-009", Name: "Phone Stand", BasePrice: 19.99},
		{SKU: "PROD-010", Name: "Notebook", BasePrice: 9.99},
	}

	fmt.Println("=== Product Enrichment Performance Comparison ===")
	fmt.Println()

	// Test 1: Sequential processing
	fmt.Println("--- Sequential Processing ---")
	start := time.Now()
	enrichedSeq := EnrichProductDataSequential(client, products)
	seqDuration := time.Since(start)
	fmt.Printf("Processed %d products in %v\n", len(enrichedSeq), seqDuration)
	fmt.Printf("Average time per product: %v\n", seqDuration/time.Duration(len(products)))
	fmt.Println()

	// Test 2: Concurrent processing (goroutine per product)
	fmt.Println("--- Concurrent Processing (Goroutine per Product) ---")
	start = time.Now()
	enrichedConcurrent := EnrichProductDataConcurrent(client, products)
	concurrentDuration := time.Since(start)
	fmt.Printf("Processed %d products in %v\n", len(enrichedConcurrent), concurrentDuration)
	fmt.Printf("Average time per product: %v\n", concurrentDuration/time.Duration(len(products)))
	fmt.Printf("Speedup: %.2fx faster than sequential\n", float64(seqDuration)/float64(concurrentDuration))
	fmt.Println()

	// Test 3: Fan-out pattern (goroutines per API call)
	fmt.Println("--- Fan-Out Pattern (Parallel API Calls per Product) ---")
	start = time.Now()
	enrichedFanOut := EnrichProductDataFanOut(client, products)
	fanOutDuration := time.Since(start)
	fmt.Printf("Processed %d products in %v\n", len(enrichedFanOut), fanOutDuration)
	fmt.Printf("Average time per product: %v\n", fanOutDuration/time.Duration(len(products)))
	fmt.Printf("Speedup: %.2fx faster than sequential\n", float64(seqDuration)/float64(fanOutDuration))
	fmt.Printf("Speedup: %.2fx faster than concurrent\n", float64(concurrentDuration)/float64(fanOutDuration))
	fmt.Println()

	// Show a sample enriched product
	fmt.Println("--- Sample Enriched Product ---")
	sample := enrichedFanOut[0]
	fmt.Printf("SKU: %s\n", sample.SKU)
	fmt.Printf("Name: %s\n", sample.Name)
	fmt.Printf("Base Price: $%.2f\n", sample.BasePrice)
	fmt.Printf("Dynamic Price: $%.2f\n", sample.DynamicPrice)
	fmt.Printf("Inventory: %d units\n", sample.InventoryQty)
	fmt.Printf("Rating: %.1f (%d reviews)\n", sample.AvgRating, sample.ReviewCount)
}
