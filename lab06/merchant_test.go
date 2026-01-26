package main

import (
	"testing"
)



// TestValidateMerchant uses table-driven tests to validate merchant validation logic.
// Each test case has a name, input merchant, and expected error message.
// This is a sample test to demonstrate the table-driven testing pattern in Go.
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
		// TODO: Part 1 - Add more test cases for other validation scenarios:
		// - missing name
		// - missing category
		// - missing country
		// - missing status
	}

	// Run each test case as a subtest
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMerchant(tt.merchant)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

// TODO: Part 2 - Add tests for MerchantStore methods:
// - TestMerchantStore_Create
// - TestMerchantStore_GetByID
// - TestMerchantStore_GetAll

// TODO: Part 3 - Create server_test.go for HTTP handler tests:
// - Test GET /merchants
// - Test GET /merchants/{id}
// - Test POST /merchants
// - Test GET /products/{sku}/enriched

// TODO: Part 4 - Add test helper functions:
// - createTestMerchant() - creates a valid test merchant
// - setupTestServer() - creates a test server with a test store
// Mark helpers with t.Helper() so test failures point to the actual test line

// TODO: Part 5 - Run tests with coverage:
// go test -cover
// go test -coverprofile=coverage.out
// go tool cover -html=coverage.out
