# Lab 5: Intro to Concurrency

In this lab, you'll add concurrent external API calls to enrich product data and see dramatic performance improvements from Go's concurrency model.

### Continuing from Lab 4
This lab continues directly from Lab 4. You can continue to iterate on that lab or start over with the provided starter files in this directory. Note that all labs contain a solution directory (if you are stuck).

## The Scenario

Your e-commerce platform needs to enrich product data from multiple sources: inventory levels from the warehouse API, pricing from the pricing engine, and customer reviews from the review service. Each API call takes 100-300ms. Processing sequentially is slow - goroutines can parallelize this work.

**Starter files provided:** `enrichment.go`, `external.go`, `product.go`, `server.go`, `main.go`, and `merchant.go` with partial implementations.

**Note:** The performance demo harness in `main.go` is fully provided so you can run and observe the benchmarks. Your implementation work focuses on the enrichment functions in `enrichment.go`.


### Part 1: Review the Mock External API
1. In this lab, you'll work with the provided starter files.
2. Review the `Product` struct in `product.go` with fields: SKU, name, basePrice.
3. In `external.go`, review the functions that simulate external API calls:
   - `FetchInventoryLevel` - returns stock quantity after a delay.
   - `FetchDynamicPrice` - returns current price based on demand.
   - `FetchReviewSummary` - returns average rating and review count.

### Part 2: Sequential Implementation
1. In `enrichment.go`, complete `EnrichProductDataSequential` that calls all three external functions in sequence for each SKU
2. Test with 10 SKUs and measure total execution time

### Part 3: Concurrent Implementation with WaitGroups
1. Implement `EnrichProductDataConcurrent` that launches a goroutine for each SKU
2. Use a WaitGroup to wait for all goroutines to complete
3. Compare execution time with the sequential version

**Example WaitGroup pattern:**
```go
var wg sync.WaitGroup

for _, sku := range skus {
    wg.Add(1) // Increment counter before launching goroutine
    go func(s string) {
        defer wg.Done() // Decrement counter when goroutine completes
        
        // Process the SKU
        result := EnrichSingleProduct(s)
        // Store result somehow...
    }(sku) // Pass sku as parameter to avoid closure issues
}

wg.Wait() // Block until all goroutines call Done()
```
Important: Pass loop variables as parameters to avoid all goroutines seeing the last value.

### Part 4: Handle Results with Channels
1. Create a struct to hold enriched product data
2. Update your concurrent function to use a buffered channel for collecting results
3. Ensure all data is correctly gathered

**Example buffered channel for collecting results:**
```go
// Create buffered channel - buffer size = number of expected results
results := make(chan EnrichedProduct, len(skus))

var wg sync.WaitGroup
for _, sku := range skus {
    wg.Add(1)
    go func(s string) {
        defer wg.Done()
        
        enriched := EnrichSingleProduct(s)
        results <- enriched // Send to channel (won't block because of buffer)
    }(sku)
}

wg.Wait()
close(results) // Close channel after all sends complete

// Collect results
var enrichedProducts []EnrichedProduct
for product := range results { // Range over channel until closed
    enrichedProducts = append(enrichedProducts, product)
}
```
The `<-` operator sends/receives on channels. Buffered channels prevent blocking when sender is faster than receiver.

### Part 5: Fan-Out Pattern (Optional)
1. Implement `EnrichSingleProductFanOut` that launches three goroutines (one per API call) for a single product
2. Use channels to collect results from each goroutine
3. Combine the results into a single enriched product struct

### Part 6: Performance Analysis (Optional)
1. Run the provided performance demo in `main.go` that compares sequential, concurrent (goroutine per product), and fan-out approaches
2. The demo tests with varying numbers of products and displays execution times
3. Observe the speedup factors and understand which concurrency pattern performs best for this I/O-bound workload

### Part 7: Add to Your API (Optional)
1. Add a GET /products/{sku}/enriched endpoint to your merchant API (or create a new product API)
2. Use your concurrent enrichment function to fetch and return enriched data
3. Test and observe the response time