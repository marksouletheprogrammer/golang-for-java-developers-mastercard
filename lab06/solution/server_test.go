package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	commonhttp "golang-for-java-developers-training/common/http"
)

// createTestMerchant is a helper function that creates a test merchant.
// Marked with t.Helper() so failures report the calling test line number.
func createTestMerchant(t *testing.T, id, name, category, country, status string) *Merchant {
	t.Helper()
	return &Merchant{
		ID:       id,
		Name:     name,
		Category: category,
		Country:  country,
		Status:   status,
	}
}

// setupTestServer is a helper that creates a server with a fresh store.
func setupTestServer(t *testing.T) *Server {
	t.Helper()
	return NewServer(NewMerchantStore())
}

// TestHandleGetMerchants tests the GET /merchants endpoint.
func TestHandleGetMerchants(t *testing.T) {
	t.Run("returns pre-populated sample merchants", func(t *testing.T) {
		server := setupTestServer(t)

		req := httptest.NewRequest(http.MethodGet, "/merchants", nil)
		rec := httptest.NewRecorder()

		server.handleMerchants(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var merchants []*Merchant
		if err := json.NewDecoder(rec.Body).Decode(&merchants); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		// NewMerchantStore pre-populates with 3 sample merchants
		if len(merchants) != 3 {
			t.Errorf("expected 3 pre-populated merchants, got %d", len(merchants))
		}
	})

	t.Run("returns all merchants including new ones", func(t *testing.T) {
		server := setupTestServer(t)

		// Store already has 3 pre-populated merchants
		initialCount := len(server.store.GetAll())

		// Add test merchants with unique IDs
		testMerchants := []*Merchant{
			createTestMerchant(t, "MERCH-TEST-001", "Merchant 1", "Retail", "USA", "active"),
			createTestMerchant(t, "MERCH-TEST-002", "Merchant 2", "Tech", "UK", "inactive"),
		}

		for _, m := range testMerchants {
			server.store.Create(m)
		}

		req := httptest.NewRequest(http.MethodGet, "/merchants", nil)
		rec := httptest.NewRecorder()

		server.handleMerchants(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var merchants []*Merchant
		if err := json.NewDecoder(rec.Body).Decode(&merchants); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		expectedCount := initialCount + len(testMerchants)
		if len(merchants) != expectedCount {
			t.Errorf("expected %d merchants, got %d", expectedCount, len(merchants))
		}
	})
}

// TestHandleGetMerchantByID tests the GET /merchants/{id} endpoint.
func TestHandleGetMerchantByID(t *testing.T) {
	tests := []struct {
		name           string
		merchantID     string
		setupStore     func(*MerchantStore)
		expectedStatus int
	}{
		{
			name:       "returns merchant when exists",
			merchantID: "MERCH-001",
			setupStore: func(store *MerchantStore) {
				store.Create(createTestMerchant(t, "MERCH-001", "Test", "Retail", "USA", "active"))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "returns 404 when merchant not found",
			merchantID:     "NON-EXISTENT",
			setupStore:     func(store *MerchantStore) {},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(t)
			tt.setupStore(server.store)

			req := httptest.NewRequest(http.MethodGet, "/merchants/"+tt.merchantID, nil)
			rec := httptest.NewRecorder()

			server.handleMerchantByID(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

// TestHandleCreateMerchant tests the POST /merchants endpoint.
func TestHandleCreateMerchant(t *testing.T) {
	tests := []struct {
		name           string
		merchant       *Merchant
		setupStore     func(*MerchantStore)
		expectedStatus int
	}{
		{
			name:           "creates valid merchant",
			merchant:       createTestMerchant(t, "MERCH-TEST-999", "Test", "Retail", "USA", "active"),
			setupStore:     func(store *MerchantStore) {},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "fails with missing ID",
			merchant:       createTestMerchant(t, "", "Test", "Retail", "USA", "active"),
			setupStore:     func(store *MerchantStore) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "fails with invalid status",
			merchant:       createTestMerchant(t, "MERCH-TEST-888", "Test", "Retail", "USA", "pending"),
			setupStore:     func(store *MerchantStore) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "fails with duplicate ID",
			merchant: createTestMerchant(t, "MERCH-001", "Test", "Retail", "USA", "active"),
			setupStore: func(store *MerchantStore) {
				// MERCH-001 already exists in pre-populated data, no need to create
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(t)
			tt.setupStore(server.store)

			body, _ := json.Marshal(tt.merchant)
			req := httptest.NewRequest(http.MethodPost, "/merchants", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			server.handleCreateMerchant(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Verify Location header on success
			if tt.expectedStatus == http.StatusCreated {
				location := rec.Header().Get("Location")
				expected := "/merchants/" + tt.merchant.ID
				if location != expected {
					t.Errorf("expected Location header %q, got %q", expected, location)
				}

				// Verify merchant was actually created in store
				_, err := server.store.GetByID(tt.merchant.ID)
				if err != nil {
					t.Errorf("merchant was not created in store: %v", err)
				}
			}
		})
	}
}

// TestInvalidJSON tests that invalid JSON is rejected.
func TestInvalidJSON(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/merchants", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.handleCreateMerchant(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestMethodNotAllowed tests that wrong HTTP methods are rejected.
func TestMethodNotAllowed(t *testing.T) {
	server := setupTestServer(t)

	t.Run("PUT not allowed on /merchants", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/merchants", nil)
		rec := httptest.NewRecorder()

		server.handleMerchants(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("POST not allowed on /merchants/{id}", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/merchants/MERCH-001", nil)
		rec := httptest.NewRecorder()

		server.handleMerchantByID(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

// TestGetMerchantByIDEdgeCases tests edge cases for merchant ID extraction.
func TestGetMerchantByIDEdgeCases(t *testing.T) {
	server := setupTestServer(t)
	server.store.Create(createTestMerchant(t, "MERCH-001", "Test", "Retail", "USA", "active"))

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "empty ID",
			path:           "/merchants/",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "valid ID",
			path:           "/merchants/MERCH-001",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ID with trailing slash",
			path:           "/merchants/MERCH-001/",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			server.handleMerchantByID(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

// TestResponseFormat tests that responses have correct format and headers.
func TestResponseFormat(t *testing.T) {
	server := setupTestServer(t)

	t.Run("JSON responses have correct content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/merchants", nil)
		rec := httptest.NewRecorder()

		server.handleGetMerchants(rec, req)

		contentType := rec.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}
	})

	t.Run("error responses are valid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/merchants/NON-EXISTENT", nil)
		rec := httptest.NewRecorder()

		server.handleMerchantByID(rec, req)

		var errResp commonhttp.ErrorResponse
		if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
			t.Fatalf("failed to decode error response: %v", err)
		}

		if errResp.Error == "" {
			t.Error("expected error message in response")
		}
	})
}

// TestHandleProductEnriched tests the GET /products/{sku}/enriched endpoint.
func TestHandleProductEnriched(t *testing.T) {
	server := setupTestServer(t)

	t.Run("returns enriched product data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/LAPTOP-001/enriched", nil)
		rec := httptest.NewRecorder()

		server.handleProductEnriched(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var enriched EnrichedProduct
		if err := json.NewDecoder(rec.Body).Decode(&enriched); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if enriched.SKU != "LAPTOP-001" {
			t.Errorf("expected SKU LAPTOP-001, got %s", enriched.SKU)
		}

		// Verify enriched data is present
		if enriched.InventoryQty == 0 {
			t.Error("expected non-zero inventory quantity")
		}
		if enriched.DynamicPrice == 0 {
			t.Error("expected non-zero dynamic price")
		}
	})

	t.Run("returns 400 for invalid path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products/LAPTOP-001/invalid", nil)
		rec := httptest.NewRecorder()

		server.handleProductEnriched(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 400 for missing SKU", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products//enriched", nil)
		rec := httptest.NewRecorder()

		server.handleProductEnriched(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns 405 for non-GET method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products/LAPTOP-001/enriched", nil)
		rec := httptest.NewRecorder()

		server.handleProductEnriched(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}
