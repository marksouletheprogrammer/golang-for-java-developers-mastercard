package domain

// Product represents a basic product in the e-commerce system.
type Product struct {
	SKU       string
	Name      string
	BasePrice float64
}

// EnrichedProduct contains product data enriched from multiple external sources.
type EnrichedProduct struct {
	SKU          string
	Name         string
	BasePrice    float64
	InventoryQty int
	DynamicPrice float64
	AvgRating    float64
	ReviewCount  int
}
