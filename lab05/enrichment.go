package main

// EnrichProductDataSequential processes products one at a time, calling all three
// external APIs sequentially for each product. This is slow because we wait for
// each API call to complete before starting the next one.
// TODO: Part 2 - Implement sequential processing baseline
func EnrichProductDataSequential(products []Product) []EnrichedProduct {
	// TODO: Create a slice to hold enriched products
	// TODO: Loop through each product
	// TODO: Call FetchInventoryLevel, FetchDynamicPrice, and FetchReviewSummary sequentially
	// TODO: Build EnrichedProduct and append to slice
	// TODO: Return the slice
	return nil
}

// EnrichProductDataConcurrent processes multiple products concurrently using goroutines.
// Each product is enriched in its own goroutine, allowing parallel API calls.
// Uses a WaitGroup to wait for all goroutines to complete.
// TODO: Part 3 - Implement concurrent processing with WaitGroups
func EnrichProductDataConcurrent(products []Product) []EnrichedProduct {
	// TODO: Create a buffered channel sized to hold all results
	// TODO: Create a WaitGroup to track goroutines
	// TODO: Launch a goroutine for each product
	// TODO: In each goroutine: fetch data, create EnrichedProduct, send to channel
	// TODO: Use another goroutine to wait for all workers and close the channel
	// TODO: Collect all results from the channel
	// TODO: Return the slice
	return nil
}

// EnrichSingleProductFanOut enriches a single product using fan-out pattern.
// Launches three goroutines (one per API call) to fetch data in parallel.
// This is faster than sequential because all three API calls happen concurrently.
// TODO: Part 5 - Implement fan-out pattern for single product
func EnrichSingleProductFanOut(product Product) EnrichedProduct {
	// TODO: Create channels for inventory, price, and reviews
	// TODO: Launch three goroutines, one for each API call
	// TODO: Each goroutine sends its result to the appropriate channel
	// TODO: Receive results from all three channels
	// TODO: Build and return EnrichedProduct
	return EnrichedProduct{}
}

// EnrichProductDataFanOut processes multiple products, using fan-out pattern
// for each product. This combines both approaches: goroutine per product,
// and goroutines for each API call within a product.
// TODO: Part 5 - Implement fan-out pattern for multiple products
func EnrichProductDataFanOut(products []Product) []EnrichedProduct {
	// TODO: Create results channel and WaitGroup
	// TODO: Launch a goroutine for each product
	// TODO: Each goroutine calls EnrichSingleProductFanOut and sends result to channel
	// TODO: Wait for all goroutines and close channel
	// TODO: Collect and return results
	return nil
}
