# Lab 5 Solution

## How to Run

Start the application:
```bash
cd lab05/solution
go run .
```

The program will:
1. **Run performance demo** - Compare sequential vs concurrent vs fan-out enrichment approaches
2. **Start API server** - Merchant API (Lab 4) + new product enrichment endpoint (Lab 5)

The server runs on http://localhost:8080 with all Lab 4 endpoints plus the new product endpoint.

## Testing Product Enrichment

Test the new product enrichment endpoint:
```bash
# Get enriched product data for a SKU
curl http://localhost:8080/products/LAPTOP-001/enriched

# Try different SKUs
curl http://localhost:8080/products/MOUSE-500/enriched
curl http://localhost:8080/products/MONITOR-X/enriched
```

The endpoint uses the **fan-out pattern** to fetch inventory, pricing, and reviews concurrently, returning enriched product data with ~200ms response time (vs ~600ms if sequential).

## Performance Comparison Results

**Typical Performance (10 products, ~200ms avg API latency):**

| Approach | Total Time | Speedup vs Sequential |
|----------|-----------|----------------------|
| Sequential | ~6.0s | 1.0x baseline |
| Concurrent (goroutine per product) | ~600ms | ~10x faster |
| Fan-Out (parallel API calls) | ~200ms | ~30x faster |

**Why Fan-Out is Fastest:**
- **Sequential**: Each product waits for 3 API calls in sequence (200ms × 3 = 600ms per product)
- **Concurrent**: Products processed in parallel, but each product's APIs are sequential (600ms per product, but all products run together)
- **Fan-Out**: Products AND their API calls run in parallel (only limited by slowest API call ~200ms)

## Concurrency Patterns Explained

**WaitGroup Pattern (Concurrent):**
- Launch a goroutine for each product
- Each goroutine processes one product (3 API calls sequentially)
- WaitGroup ensures we wait for all goroutines to complete
- Good for I/O-bound per-item processing

**Fan-Out Pattern:**
- Launch goroutines for both products AND each product's API calls
- Use channels to collect results from parallel operations
- Maximizes parallelism when each item has multiple independent operations
- Best for maximizing throughput with multiple I/O operations per item

**Buffered Channels:**
- Buffer size = number of expected results prevents blocking
- Allows goroutines to send results without waiting for receiver
- Important for performance in high-concurrency scenarios

## Comparison to Java

**Java Approach:**
```java
// Java equivalent using ExecutorService
ExecutorService executor = Executors.newFixedThreadPool(10);
List<Future<EnrichedProduct>> futures = new ArrayList<>();

for (Product product : products) {
    futures.add(executor.submit(() -> enrichProduct(product)));
}

// Wait for all
for (Future<EnrichedProduct> future : futures) {
    enrichedProducts.add(future.get());
}
executor.shutdown();
```

**Key Differences:**
- Go: Lightweight goroutines (thousands possible), Java: Heavy threads (limited)
- Go: Channels for communication, Java: Shared memory with locks or BlockingQueue
- Go: Built-in `select` for channel operations, Java: Manual polling or callbacks
- Go: WaitGroup for synchronization, Java: CountDownLatch or join()
- Go: Simple syntax for concurrency, Java: More verbose with executors

### Avoiding Closure Pitfalls

When launching goroutines in a loop, pass loop variables as parameters:
```go
for _, item := range items {
    go func(i Item) {  // Pass as parameter
        process(i)
    }(item)
}
```

Without passing as parameter, all goroutines would see the final loop value due to closure capturing the loop variable reference.
