# Lab 6: Automated Testing & Table-Driven Tests

In this lab, you'll write comprehensive tests for your merchant service using Go's testing conventions, particularly table-driven tests.

**Starter files provided:** All files from Lab 5 including `merchant_test.go` with starter tests, `helpers.go`, `enrichment.go`, `server.go`, and other implementation files.

**Note:** The product enrichment functions from Lab 5 are intentionally provided in `enrichment.go`. This lab focuses on writing comprehensive tests for existing code, not implementing new functionality.

### Continuing from Lab 5
This lab continues directly from Lab 5. You can continue to iterate on that lab or start over with the provided starter files in this directory. Note that all labs contain a solution directory (if you are stuck).

### Part 1: Test Merchant Validation
1. Create a `ValidateMerchant` function that checks merchant data (e.g., ID not empty, name not empty, valid status like "active" or "inactive", valid country code)
2. In the provided `merchant_test.go`, review the existing starter tests and expand the table-driven test for validation
3. Define test cases covering valid merchants and various validation failures (empty ID, empty name, invalid status, etc.)
4. Use subtests (with `t.Run()`) for clear failure reporting

**Example table-driven test structure:**
```go
func TestValidateMerchant(t *testing.T) {
    tests := []struct {
        name     string
        merchant Merchant
        wantErr  bool
    }{
        {
            name:     "valid merchant",
            merchant: Merchant{ID: "M001", Name: "Test Co", Status: "active", Country: "USA"},
            wantErr:  false,
        },
        {
            name:     "empty ID",
            merchant: Merchant{ID: "", Name: "Test Co", Status: "active", Country: "USA"},
            wantErr:  true,
        },
        //...
    }
    // Essentially just executed the function for each item in the slice.
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateMerchant(&tt.merchant)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateMerchant() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```
This pattern makes it easy to add new test cases by just adding to the slice.

### Part 2: Test Merchant Methods
1. Write table-driven tests for any business logic methods on your Merchant struct
2. Test edge cases and boundary conditions (e.g., empty input, nil values, zero values, very long strings, special characters)

### Part 3: Test HTTP Handlers
1. Create `server_test.go` and test your HTTP endpoints using `httptest.NewRecorder()`
2. Test GET /merchants, GET /merchants/{id}, POST /merchants, and GET /products/{sku}/enriched endpoints
3. Verify status codes, response bodies, and that operations actually modify storage
4. Test both success and error cases

### Part 4: Test Helpers and Setup (Optional)
1. Create helper functions like `createTestMerchant()` and `setupTestServer()` to reduce code duplication in tests
2. Call `t.Helper()` at the start of these helper functions - this makes test failure messages show the line where the helper was called, not the line inside the helper, making failures easier to debug

### Part 5: Mock External Dependencies
The enrichment functions have been refactored to use dependency injection with the `ExternalAPIClient` interface. This allows you to replace slow, non-deterministic external API calls with fast, predictable mocks in your tests.

**What's provided:**
- `ExternalAPIClient` interface with three methods: `FetchInventoryLevel`, `FetchDynamicPrice`, `FetchReviewSummary`
- `RealExternalAPIClient` struct that wraps the actual external API calls
- All enrichment functions refactored to accept `ExternalAPIClient` as the first parameter

**Your task:**
1. In `enrichment_test.go`, create a `MockExternalAPIClient` struct with fields to store return values:
   ```go
   type MockExternalAPIClient struct {
       InventoryQty int
       DynamicPrice float64
       AvgRating    float64
       ReviewCount  int
   }
   ```

2. Implement the three interface methods to return the configured values:
   ```go
   func (m *MockExternalAPIClient) FetchInventoryLevel(sku string) int {
       return m.InventoryQty
   }
   // ... implement FetchDynamicPrice and FetchReviewSummary
   ```

3. Write tests using the mock with deterministic data:
   ```go
   func TestEnrichSingleProductFanOutWithMock(t *testing.T) {
       mock := &MockExternalAPIClient{
           InventoryQty: 50,
           DynamicPrice: 99.99,
           AvgRating:    4.5,
           ReviewCount:  100,
       }
       
       product := Product{SKU: "TEST-001", Name: "Test", BasePrice: 100.00}
       enriched := EnrichSingleProductFanOut(mock, product)
       
       // Assert exact values - no randomness!
       if enriched.InventoryQty != 50 {
           t.Errorf("expected inventory 50, got %d", enriched.InventoryQty)
       }
       // ... more assertions
   }
   ```

4. Write similar tests for other enrichment functions using the mock

**Comparison to Java:**
- **Java (Mockito)**: Uses reflection and dynamic proxies to create mocks at runtime
  ```java
  ExternalAPIClient mock = Mockito.mock(ExternalAPIClient.class);
  when(mock.fetchInventoryLevel("TEST-001")).thenReturn(50);
  ```
- **Go**: Manual mock implementation using structs
  - More explicit and type-safe
  - No magic, easier to debug
  - No external library needed for simple cases
  - For complex scenarios, libraries like `testify/mock` or `gomock` are available

### Part 6: Measure Coverage (Optional)
1. Run tests with coverage reporting enabled
2. Generate an HTML coverage report
3. Add tests to improve coverage (aim for >80%)