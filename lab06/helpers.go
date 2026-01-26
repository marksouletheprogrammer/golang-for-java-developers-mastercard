package main

import "golang-for-java-developers-training/common/domain"

// initSampleMerchants adds sample merchants to the store for demo purposes.
func initSampleMerchants(store *domain.MerchantStore) {
	store.Create(&domain.Merchant{
		ID:       "MERCH-001",
		Name:     "Tech Solutions Inc",
		Category: "Technology",
		Country:  "USA",
		Status:   "active",
	})
	store.Create(&domain.Merchant{
		ID:       "MERCH-002",
		Name:     "Global Retail Co",
		Category: "Retail",
		Country:  "UK",
		Status:   "active",
	})
	store.Create(&domain.Merchant{
		ID:       "MERCH-003",
		Name:     "Food Services Ltd",
		Category: "Food & Beverage",
		Country:  "Canada",
		Status:   "inactive",
	})
}
