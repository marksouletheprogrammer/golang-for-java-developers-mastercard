package main

import (
	"sync"

	"golang-for-java-developers-training/common/domain"
	"golang-for-java-developers-training/common/external"
)

// Type aliases for convenience
type Product = domain.Product
type EnrichedProduct = domain.EnrichedProduct

// ExternalAPIClient defines the contract for external API calls.
// This interface allows dependency injection for testing.
type ExternalAPIClient interface {
	FetchInventoryLevel(sku string) int
	FetchDynamicPrice(sku string, basePrice float64) float64
	FetchReviewSummary(sku string) (avgRating float64, reviewCount int)
}

// RealExternalAPIClient wraps the actual external API calls.
type RealExternalAPIClient struct{}

func (r *RealExternalAPIClient) FetchInventoryLevel(sku string) int {
	return external.FetchInventoryLevel(sku)
}

func (r *RealExternalAPIClient) FetchDynamicPrice(sku string, basePrice float64) float64 {
	return external.FetchDynamicPrice(sku, basePrice)
}

func (r *RealExternalAPIClient) FetchReviewSummary(sku string) (avgRating float64, reviewCount int) {
	return external.FetchReviewSummary(sku)
}

// EnrichProductDataSequential processes products one at a time, calling all three
// external APIs sequentially for each product. This is slow because we wait for
// each API call to complete before starting the next one.
func EnrichProductDataSequential(client ExternalAPIClient, products []Product) []EnrichedProduct {
	enriched := make([]EnrichedProduct, 0, len(products))

	for _, product := range products {
		// Call each external API sequentially - waits for each to complete
		inventory := client.FetchInventoryLevel(product.SKU)
		price := client.FetchDynamicPrice(product.SKU, product.BasePrice)
		avgRating, reviewCount := client.FetchReviewSummary(product.SKU)

		enriched = append(enriched, EnrichedProduct{
			SKU:          product.SKU,
			Name:         product.Name,
			BasePrice:    product.BasePrice,
			InventoryQty: inventory,
			DynamicPrice: price,
			AvgRating:    avgRating,
			ReviewCount:  reviewCount,
		})
	}

	return enriched
}

// EnrichProductDataConcurrent processes multiple products concurrently using goroutines.
// Each product is enriched in its own goroutine, allowing parallel API calls.
// Uses a WaitGroup to wait for all goroutines to complete.
func EnrichProductDataConcurrent(client ExternalAPIClient, products []Product) []EnrichedProduct {
	// Buffered channel sized to hold all results
	results := make(chan EnrichedProduct, len(products))

	// WaitGroup tracks how many goroutines are running
	var wg sync.WaitGroup

	// Launch a goroutine for each product
	for _, product := range products {
		wg.Add(1)

		// Launch goroutine - must pass product as parameter to avoid closure issues
		go func(p Product) {
			defer wg.Done()

			// These three calls still happen sequentially within this goroutine
			inventory := client.FetchInventoryLevel(p.SKU)
			price := client.FetchDynamicPrice(p.SKU, p.BasePrice)
			avgRating, reviewCount := client.FetchReviewSummary(p.SKU)

			// Send result to channel
			results <- EnrichedProduct{
				SKU:          p.SKU,
				Name:         p.Name,
				BasePrice:    p.BasePrice,
				InventoryQty: inventory,
				DynamicPrice: price,
				AvgRating:    avgRating,
				ReviewCount:  reviewCount,
			}
		}(product)
	}

	// Wait for all goroutines to complete, then close the channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results from the channel
	enriched := make([]EnrichedProduct, 0, len(products))
	for result := range results {
		enriched = append(enriched, result)
	}

	return enriched
}

// EnrichSingleProductFanOut enriches a single product using fan-out pattern.
// Launches three goroutines (one per API call) to fetch data in parallel.
// This is faster than sequential because all three API calls happen concurrently.
func EnrichSingleProductFanOut(client ExternalAPIClient, product Product) EnrichedProduct {
	// Create channels for each piece of data
	inventoryCh := make(chan int, 1)
	priceCh := make(chan float64, 1)
	reviewCh := make(chan struct {
		rating float64
		count  int
	}, 1)

	// Launch three goroutines to fetch data in parallel
	go func() {
		inventoryCh <- client.FetchInventoryLevel(product.SKU)
	}()

	go func() {
		priceCh <- client.FetchDynamicPrice(product.SKU, product.BasePrice)
	}()

	go func() {
		avgRating, reviewCount := client.FetchReviewSummary(product.SKU)
		reviewCh <- struct {
			rating float64
			count  int
		}{avgRating, reviewCount}
	}()

	// Collect results from all three channels
	// These operations block until data is available, but that's fine
	// because all three API calls are running in parallel
	inventory := <-inventoryCh
	price := <-priceCh
	reviews := <-reviewCh

	return EnrichedProduct{
		SKU:          product.SKU,
		Name:         product.Name,
		BasePrice:    product.BasePrice,
		InventoryQty: inventory,
		DynamicPrice: price,
		AvgRating:    reviews.rating,
		ReviewCount:  reviews.count,
	}
}

// EnrichProductDataFanOut processes multiple products, using fan-out pattern
// for each product. This combines both approaches: goroutine per product,
// and goroutines for each API call within a product.
func EnrichProductDataFanOut(client ExternalAPIClient, products []Product) []EnrichedProduct {
	results := make(chan EnrichedProduct, len(products))
	var wg sync.WaitGroup

	for _, product := range products {
		wg.Add(1)

		go func(p Product) {
			defer wg.Done()
			results <- EnrichSingleProductFanOut(client, p)
		}(product)
	}

	// Wait and close channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	enriched := make([]EnrichedProduct, 0, len(products))
	for result := range results {
		enriched = append(enriched, result)
	}

	return enriched
}
