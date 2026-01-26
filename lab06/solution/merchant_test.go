package main

import (
	"testing"
)

// TestValidateMerchant uses table-driven tests to validate merchant validation logic.
// Each test case has a name, input merchant, and expected error message.
func TestValidateMerchant(t *testing.T) {
	tests := []struct {
		name        string
		merchant    *Merchant
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid merchant",
			merchant: &Merchant{
				ID:       "MERCH-001",
				Name:     "Test Merchant",
				Category: "Retail",
				Country:  "USA",
				Status:   "active",
			},
			expectError: false,
		},
		{
			name: "missing ID",
			merchant: &Merchant{
				Name:     "Test Merchant",
				Category: "Retail",
				Country:  "USA",
				Status:   "active",
			},
			expectError: true,
			errorMsg:    "id is required",
		},
		{
			name: "missing name",
			merchant: &Merchant{
				ID:       "MERCH-001",
				Category: "Retail",
				Country:  "USA",
				Status:   "active",
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "missing category",
			merchant: &Merchant{
				ID:      "MERCH-001",
				Name:    "Test Merchant",
				Country: "USA",
				Status:  "active",
			},
			expectError: true,
			errorMsg:    "category is required",
		},
		{
			name: "missing country",
			merchant: &Merchant{
				ID:       "MERCH-001",
				Name:     "Test Merchant",
				Category: "Retail",
				Status:   "active",
			},
			expectError: true,
			errorMsg:    "country is required",
		},
		{
			name: "missing status",
			merchant: &Merchant{
				ID:       "MERCH-001",
				Name:     "Test Merchant",
				Category: "Retail",
				Country:  "USA",
			},
			expectError: true,
			errorMsg:    "status is required",
		},
		{
			name: "invalid status",
			merchant: &Merchant{
				ID:       "MERCH-001",
				Name:     "Test Merchant",
				Category: "Retail",
				Country:  "USA",
				Status:   "pending",
			},
			expectError: true,
			errorMsg:    "status must be 'active' or 'inactive'",
		},
	}

	// Run all test cases
	for _, tt := range tests {
		// t.Run creates a subtest with the given name
		// This provides clear failure reporting and allows running specific tests
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMerchant(tt.merchant)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestInitSampleMerchants tests that sample merchants are correctly initialized.
func TestInitSampleMerchants(t *testing.T) {
	store := NewMerchantStore()

	// NewMerchantStore already initializes with 3 sample merchants from common package
	merchants := store.GetAll()
	if len(merchants) != 3 {
		t.Errorf("expected 3 sample merchants from NewMerchantStore, got %d", len(merchants))
	}

	// Verify specific merchants exist
	expectedIDs := []string{"MERCH-001", "MERCH-002", "MERCH-003"}
	for _, id := range expectedIDs {
		merchant, err := store.GetByID(id)
		if err != nil {
			t.Errorf("expected sample merchant %s to exist", id)
		}
		if merchant.ID != id {
			t.Errorf("expected merchant ID %s, got %s", id, merchant.ID)
		}
	}

	// initSampleMerchants should handle duplicates gracefully (will fail to add)
	initSampleMerchants(store)

	// Should still have 3 merchants (duplicates rejected)
	merchants = store.GetAll()
	if len(merchants) != 3 {
		t.Errorf("expected 3 merchants after initSampleMerchants, got %d", len(merchants))
	}
}

// TestMerchantStore tests the merchant store operations.
func TestMerchantStore(t *testing.T) {
	t.Run("create and retrieve merchant", func(t *testing.T) {
		store := NewMerchantStore()

		merchant := &Merchant{
			ID:       "MERCH-TEST-001",
			Name:     "Test Merchant",
			Category: "Retail",
			Country:  "USA",
			Status:   "active",
		}

		// Create merchant
		err := store.Create(merchant)
		if err != nil {
			t.Fatalf("failed to create merchant: %v", err)
		}

		// Retrieve by ID
		retrieved, err := store.GetByID("MERCH-TEST-001")
		if err != nil {
			t.Fatalf("failed to retrieve merchant: %v", err)
		}

		if retrieved.ID != merchant.ID || retrieved.Name != merchant.Name {
			t.Errorf("retrieved merchant doesn't match: got %+v, want %+v", retrieved, merchant)
		}
	})

	t.Run("create duplicate merchant fails", func(t *testing.T) {
		store := NewMerchantStore()

		// Use pre-populated merchant ID to test duplicate detection
		merchant := &Merchant{
			ID:       "MERCH-001",
			Name:     "Duplicate Merchant",
			Category: "Retail",
			Country:  "USA",
			Status:   "active",
		}

		// Try to create duplicate - should fail (MERCH-001 already exists)
		err := store.Create(merchant)
		if err == nil {
			t.Error("expected error when creating duplicate merchant, got none")
		}
	})

	t.Run("retrieve non-existent merchant fails", func(t *testing.T) {
		store := NewMerchantStore()

		_, err := store.GetByID("NON-EXISTENT")
		if err == nil {
			t.Error("expected error when retrieving non-existent merchant, got none")
		}
	})

	t.Run("get all merchants", func(t *testing.T) {
		store := NewMerchantStore()

		// Store already has 3 pre-populated merchants
		initialCount := len(store.GetAll())

		merchants := []*Merchant{
			{ID: "MERCH-TEST-101", Name: "Merchant 1", Category: "Retail", Country: "USA", Status: "active"},
			{ID: "MERCH-TEST-102", Name: "Merchant 2", Category: "Tech", Country: "UK", Status: "inactive"},
		}

		// Create all merchants
		for _, m := range merchants {
			if err := store.Create(m); err != nil {
				t.Fatalf("failed to create merchant %s: %v", m.ID, err)
			}
		}

		// Get all
		all := store.GetAll()
		expectedCount := initialCount + len(merchants)
		if len(all) != expectedCount {
			t.Errorf("expected %d merchants, got %d", expectedCount, len(all))
		}
	})
}
