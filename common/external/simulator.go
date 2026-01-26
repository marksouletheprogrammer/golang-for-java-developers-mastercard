package external

import (
	"math/rand"
	"time"
)

// FetchInventoryLevel simulates an external API call to get inventory quantity.
// Adds random latency between 100-300ms to simulate network delay.
func FetchInventoryLevel(sku string) int {
	// Simulate network latency
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
	
	// Return random inventory between 0-100
	return rand.Intn(101)
}

// FetchDynamicPrice simulates an external pricing engine API call.
// Returns adjusted price based on demand, with simulated network delay.
func FetchDynamicPrice(sku string, basePrice float64) float64 {
	// Simulate network latency
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
	
	// Apply random demand multiplier between 0.8 and 1.3
	multiplier := 0.8 + rand.Float64()*0.5
	return basePrice * multiplier
}

// FetchReviewSummary simulates an external review service API call.
// Returns average rating and review count with simulated network delay.
func FetchReviewSummary(sku string) (avgRating float64, reviewCount int) {
	// Simulate network latency
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
	
	// Return random rating between 3.0-5.0 and review count
	avgRating = 3.0 + rand.Float64()*2.0
	reviewCount = rand.Intn(500)
	return
}
