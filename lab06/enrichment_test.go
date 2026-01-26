package main

import (
	"testing"
)

// MockExternalAPIClient is a test double for ExternalAPIClient.
// It allows you to control the return values for testing.
// TODO: Part 5 - Add fields to store return values
type MockExternalAPIClient struct {
	// TODO: Add fields for InventoryQty, DynamicPrice, AvgRating, ReviewCount
}

// TODO: Part 5 - Implement FetchInventoryLevel method
func (m *MockExternalAPIClient) FetchInventoryLevel(sku string) int {
	// TODO: Return the configured mock value
	return 0
}

// TODO: Part 5 - Implement FetchDynamicPrice method
func (m *MockExternalAPIClient) FetchDynamicPrice(sku string, basePrice float64) float64 {
	// TODO: Return the configured mock value
	return 0.0
}

// TODO: Part 5 - Implement FetchReviewSummary method
func (m *MockExternalAPIClient) FetchReviewSummary(sku string) (avgRating float64, reviewCount int) {
	// TODO: Return the configured mock values
	return 0.0, 0
}

// TODO: Part 5 - Write test for EnrichSingleProductFanOut using mock
func TestEnrichSingleProductFanOutWithMock(t *testing.T) {
	// TODO: Create a mock with specific test values
	// TODO: Create a test product
	// TODO: Call EnrichSingleProductFanOut with mock and product
	// TODO: Assert exact expected values (no randomness!)
}

// TODO: Part 5 - Write test for EnrichProductDataSequential using mock
func TestEnrichProductDataSequentialWithMock(t *testing.T) {
	// TODO: Create a mock with specific test values
	// TODO: Create test products
	// TODO: Call EnrichProductDataSequential with mock and products
	// TODO: Assert exact expected values for all products
}

// TODO: Part 5 - Write test for EnrichProductDataConcurrent using mock
func TestEnrichProductDataConcurrentWithMock(t *testing.T) {
	// TODO: Create a mock with specific test values
	// TODO: Create test products
	// TODO: Call EnrichProductDataConcurrent with mock and products
	// TODO: Assert exact expected values for all products
}

// TODO: Part 5 - Write test for EnrichProductDataFanOut using mock
func TestEnrichProductDataFanOutWithMock(t *testing.T) {
	// TODO: Create a mock with specific test values
	// TODO: Create test products
	// TODO: Call EnrichProductDataFanOut with mock and products
	// TODO: Assert exact expected values for all products
}
