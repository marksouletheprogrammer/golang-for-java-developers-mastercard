# Lab 6 Solution

## How to Run Tests

Run all tests:
```bash
cd lab06/solution
go test -v ./...
```

Run tests with race detector:
```bash
go test -race ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate HTML coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Running the Server

The server includes all Lab 4 + Lab 5 functionality:
```bash
go run .
```

This will:
1. Run the performance demo (sequential vs concurrent vs fan-out)
2. Start the API server with merchant and product endpoints

## Key Testing Concepts

### Table-Driven Tests

Table-driven tests define test cases as a slice of structs, then loop through them. This makes it easy to add new test cases without duplicating test logic:

```go
tests := []struct {
    name     string
    input    string
    expected int
}{
    {"case 1", "input1", 1},
    {"case 2", "input2", 2},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := function(tt.input)
        if result != tt.expected {
            t.Errorf("got %v, want %v", result, tt.expected)
        }
    })
}
```

Each test case has a descriptive name for clear failure reporting. Adding new cases is just adding a struct to the slice.

### Subtests with t.Run

`t.Run(name, func)` creates a subtest. Benefits:
- Clear failure reporting showing which specific case failed
- Can run individual subtests with `-run TestName/SubtestName`
- Each subtest runs independently
- Parallel execution possible with `t.Parallel()`

### httptest Package

`httptest.NewRecorder()` captures HTTP responses without starting a real server. `httptest.NewRequest()` creates test requests. This makes HTTP handler testing fast and doesn't require network ports.

### Test Helpers

Mark helper functions with `t.Helper()` so test failures report the calling line number, not the helper's line. Helpers reduce boilerplate and make tests more readable.

### Testing HTTP Handlers

Test handlers verify:
- Status codes match expected values
- Response bodies contain correct data
- Headers are set properly (e.g., Location on create)
- Operations have side effects (e.g., merchant actually created in store)

Use `httptest.NewRecorder()` to capture responses and `json.Decoder` to parse JSON responses for verification.

### Race Detection

Run tests with `-race` flag to detect data races in concurrent code. The race detector finds unsynchronized access to shared memory. Fix races with mutexes, channels, or atomic operations.

### Coverage

`-cover` shows percentage of code executed by tests. `-coverprofile` generates detailed coverage data. `go tool cover -html` creates an HTML report showing which lines are covered. Aim for >80% coverage, but don't chase 100% if it leads to brittle tests.

### Mocking with Interfaces (Part 5)

Go uses interface-based dependency injection for mocking, not reflection-based frameworks like Java's Mockito.

**The Pattern:**
1. Define an interface for the dependency
2. Production code accepts the interface as a parameter
3. Tests pass a mock implementation

**Java approach (Mockito):**
```java
@Mock
private ExternalAPIClient mockClient;

@Test
public void testEnrichment() {
    when(mockClient.fetchInventoryLevel("TEST-001")).thenReturn(50);
    when(mockClient.fetchDynamicPrice("TEST-001", 99.99)).thenReturn(89.99);
    
    EnrichedProduct result = enricher.enrich(mockClient, product);
    
    verify(mockClient).fetchInventoryLevel("TEST-001");
}
```

**Go approach (manual mocks):**
```go
type MockExternalAPIClient struct {
    InventoryQty int
    DynamicPrice float64
}

func (m *MockExternalAPIClient) FetchInventoryLevel(sku string) int {
    return m.InventoryQty
}

func TestEnrichment(t *testing.T) {
    mock := &MockExternalAPIClient{
        InventoryQty: 50,
        DynamicPrice: 89.99,
    }
    
    enriched := EnrichSingleProductFanOut(mock, product)
    
    if enriched.InventoryQty != 50 {
        t.Errorf("expected 50, got %d", enriched.InventoryQty)
    }
}
```

**Key differences:**
- **Go is explicit**: Mock is a regular struct, no magic
- **Java uses reflection**: Mockito creates proxies at runtime
- **Go verification**: Assert return values directly
- **Java verification**: `verify()` checks method calls
- **Go libraries**: `testify/mock` or `gomock` available for complex cases
- **Type safety**: Go's approach is compile-time safe, Java's is runtime

**Benefits of Go's approach:**
- No external library needed for simple mocks
- Easy to debug (no reflection magic)
- Compile-time safety
- Clear and explicit

**When to use libraries:**
- Complex mocking scenarios (multiple calls with different returns)
- Need to verify call counts or arguments
- Mocking many methods

The key insight: Go's interfaces are small and focused, making manual mocks simple. Java interfaces tend to be larger, making frameworks like Mockito more valuable.
