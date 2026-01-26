package main

import (
	"testing"
)

// MockExternalAPIClient is a test double for ExternalAPIClient.
// It allows you to control the return values for testing.
type MockExternalAPIClient struct {
	InventoryQty int
	DynamicPrice float64
	AvgRating    float64
	ReviewCount  int
}

func (m *MockExternalAPIClient) FetchInventoryLevel(sku string) int {
	return m.InventoryQty
}

func (m *MockExternalAPIClient) FetchDynamicPrice(sku string, basePrice float64) float64 {
	return m.DynamicPrice
}

func (m *MockExternalAPIClient) FetchReviewSummary(sku string) (avgRating float64, reviewCount int) {
	return m.AvgRating, m.ReviewCount
}

// TestEnrichSingleProductFanOutWithMock tests the fan-out enrichment using a mock.
// This test runs instantly and returns deterministic results.
func TestEnrichSingleProductFanOutWithMock(t *testing.T) {
	mock := &MockExternalAPIClient{
		InventoryQty: 50,
		DynamicPrice: 89.99,
		AvgRating:    4.5,
		ReviewCount:  100,
	}

	product := Product{
		SKU:       "TEST-001",
		Name:      "Test Product",
		BasePrice: 99.99,
	}

	enriched := EnrichSingleProductFanOut(mock, product)

	// Verify basic fields are preserved
	if enriched.SKU != product.SKU {
		t.Errorf("expected SKU %s, got %s", product.SKU, enriched.SKU)
	}
	if enriched.Name != product.Name {
		t.Errorf("expected Name %s, got %s", product.Name, enriched.Name)
	}
	if enriched.BasePrice != product.BasePrice {
		t.Errorf("expected BasePrice %.2f, got %.2f", product.BasePrice, enriched.BasePrice)
	}

	// Verify enriched fields match mock values exactly
	if enriched.InventoryQty != 50 {
		t.Errorf("expected InventoryQty 50, got %d", enriched.InventoryQty)
	}
	if enriched.DynamicPrice != 89.99 {
		t.Errorf("expected DynamicPrice 89.99, got %.2f", enriched.DynamicPrice)
	}
	if enriched.AvgRating != 4.5 {
		t.Errorf("expected AvgRating 4.5, got %.1f", enriched.AvgRating)
	}
	if enriched.ReviewCount != 100 {
		t.Errorf("expected ReviewCount 100, got %d", enriched.ReviewCount)
	}
}

// TestEnrichProductDataSequentialWithMock tests sequential enrichment using a mock.
func TestEnrichProductDataSequentialWithMock(t *testing.T) {
	mock := &MockExternalAPIClient{
		InventoryQty: 25,
		DynamicPrice: 15.50,
		AvgRating:    3.8,
		ReviewCount:  42,
	}

	products := []Product{
		{SKU: "PROD-001", Name: "Product 1", BasePrice: 10.00},
		{SKU: "PROD-002", Name: "Product 2", BasePrice: 20.00},
	}

	enriched := EnrichProductDataSequential(mock, products)

	if len(enriched) != len(products) {
		t.Errorf("expected %d enriched products, got %d", len(products), len(enriched))
	}

	// All products should have the same mock values
	for i, e := range enriched {
		if e.SKU != products[i].SKU {
			t.Errorf("product %d: expected SKU %s, got %s", i, products[i].SKU, e.SKU)
		}
		if e.InventoryQty != 25 {
			t.Errorf("product %d: expected inventory 25, got %d", i, e.InventoryQty)
		}
		if e.DynamicPrice != 15.50 {
			t.Errorf("product %d: expected price 15.50, got %.2f", i, e.DynamicPrice)
		}
		if e.AvgRating != 3.8 {
			t.Errorf("product %d: expected rating 3.8, got %.1f", i, e.AvgRating)
		}
		if e.ReviewCount != 42 {
			t.Errorf("product %d: expected review count 42, got %d", i, e.ReviewCount)
		}
	}
}

// TestEnrichProductDataConcurrentWithMock tests concurrent enrichment using a mock.
func TestEnrichProductDataConcurrentWithMock(t *testing.T) {
	mock := &MockExternalAPIClient{
		InventoryQty: 75,
		DynamicPrice: 125.00,
		AvgRating:    4.9,
		ReviewCount:  250,
	}

	products := []Product{
		{SKU: "PROD-001", Name: "Product 1", BasePrice: 100.00},
		{SKU: "PROD-002", Name: "Product 2", BasePrice: 150.00},
		{SKU: "PROD-003", Name: "Product 3", BasePrice: 200.00},
	}

	enriched := EnrichProductDataConcurrent(mock, products)

	if len(enriched) != len(products) {
		t.Errorf("expected %d enriched products, got %d", len(products), len(enriched))
	}

	// Verify all products were enriched with deterministic mock values
	for _, e := range enriched {
		if e.InventoryQty != 75 {
			t.Errorf("product %s: expected inventory 75, got %d", e.SKU, e.InventoryQty)
		}
		if e.DynamicPrice != 125.00 {
			t.Errorf("product %s: expected price 125.00, got %.2f", e.SKU, e.DynamicPrice)
		}
		if e.AvgRating != 4.9 {
			t.Errorf("product %s: expected rating 4.9, got %.1f", e.SKU, e.AvgRating)
		}
		if e.ReviewCount != 250 {
			t.Errorf("product %s: expected review count 250, got %d", e.SKU, e.ReviewCount)
		}
	}
}

// TestEnrichProductDataFanOutWithMock tests fan-out pattern using a mock.
func TestEnrichProductDataFanOutWithMock(t *testing.T) {
	mock := &MockExternalAPIClient{
		InventoryQty: 10,
		DynamicPrice: 99.99,
		AvgRating:    5.0,
		ReviewCount:  500,
	}

	products := []Product{
		{SKU: "PROD-001", Name: "Product 1", BasePrice: 99.99},
		{SKU: "PROD-002", Name: "Product 2", BasePrice: 79.99},
	}

	enriched := EnrichProductDataFanOut(mock, products)

	if len(enriched) != len(products) {
		t.Errorf("expected %d enriched products, got %d", len(products), len(enriched))
	}

	// Verify all products were enriched with all fields
	for _, e := range enriched {
		if e.InventoryQty != 10 {
			t.Errorf("product %s: expected inventory 10, got %d", e.SKU, e.InventoryQty)
		}
		if e.DynamicPrice != 99.99 {
			t.Errorf("product %s: expected price 99.99, got %.2f", e.SKU, e.DynamicPrice)
		}
		if e.AvgRating != 5.0 {
			t.Errorf("product %s: expected rating 5.0, got %.1f", e.SKU, e.AvgRating)
		}
		if e.ReviewCount != 500 {
			t.Errorf("product %s: expected review count 500, got %d", e.SKU, e.ReviewCount)
		}
	}
}

// TestMockEdgeCases tests edge cases with mock values.
func TestMockEdgeCases(t *testing.T) {
	t.Run("zero inventory", func(t *testing.T) {
		mock := &MockExternalAPIClient{
			InventoryQty: 0,
			DynamicPrice: 50.00,
			AvgRating:    3.0,
			ReviewCount:  10,
		}

		product := Product{SKU: "TEST-001", Name: "Test", BasePrice: 50.00}
		enriched := EnrichSingleProductFanOut(mock, product)

		if enriched.InventoryQty != 0 {
			t.Errorf("expected zero inventory, got %d", enriched.InventoryQty)
		}
	})

	t.Run("very high price", func(t *testing.T) {
		mock := &MockExternalAPIClient{
			InventoryQty: 1,
			DynamicPrice: 9999.99,
			AvgRating:    5.0,
			ReviewCount:  1000,
		}

		product := Product{SKU: "TEST-002", Name: "Luxury", BasePrice: 5000.00}
		enriched := EnrichSingleProductFanOut(mock, product)

		if enriched.DynamicPrice != 9999.99 {
			t.Errorf("expected price 9999.99, got %.2f", enriched.DynamicPrice)
		}
	})
}
