# Lab 9: Performance & Profiling

Learn to benchmark and profile your order service using Go's built-in tools to identify performance characteristics and bottlenecks.

### Continuing from Lab 8
This lab continues directly from Lab 8. You can continue to iterate on that lab or start over with the provided starter files in this directory. Note that all labs contain a solution directory (if you are stuck).

**Starter files provided:** Complete working project from Lab 8 solution, plus `profiling/` directory with working benchmark examples and profiling function skeletons.

---

## Part 1: Benchmarking Basics

The `profiling/benchmark_test.go` file includes 4 working benchmarks as examples:
- `BenchmarkOrderValidation` - validates an order
- `BenchmarkCalculateTotal` - calculates order total
- `BenchmarkRepositoryCreate` - creates orders in repository
- `BenchmarkRepositoryGet` - retrieves orders by ID

Run the benchmarks:
```bash
cd profiling
go test -bench=. -benchmem
```

**Understanding benchmark output:**
```
BenchmarkOrderValidation-8    1000000    1234 ns/op    512 B/op    8 allocs/op
```
- `1000000`: iterations run
- `1234 ns/op`: nanoseconds per operation
- `512 B/op`: bytes allocated per operation
- `8 allocs/op`: number of allocations per operation

### Add Your Own Benchmarks

Implement these 2 benchmarks in `benchmark_test.go`:
1. `BenchmarkRepositoryGetAll` - benchmark retrieving all orders (pre-populate with 100 orders)
2. `BenchmarkServiceCreateOrder` - benchmark the full service layer create operation

**Benchmark template:**
```go
func BenchmarkYourFunction(b *testing.B) {
    // Setup code (not timed)
    testData := setupTestData()
    
    b.ResetTimer() // Start timing here
    
    for i := 0; i < b.N; i++ {
        yourFunction(testData)
    }
}
```

---

## Part 2: CPU Profiling

In `profiling/profiling.go`, implement the `ProcessOrdersWithCPUProfile` function:
1. Create a CPU profile file using `os.Create()`
2. Start CPU profiling with `pprof.StartCPUProfile()`
3. Process all orders using the service layer
4. Stop profiling with `defer pprof.StopCPUProfile()`

**Example structure:**
```go
func ProcessOrdersWithCPUProfile(orders []*domain.Order, profilePath string) error {
    f, err := os.Create(profilePath)
    if err != nil {
        return err
    }
    defer f.Close()
    
    if err := pprof.StartCPUProfile(f); err != nil {
        return err
    }
    defer pprof.StopCPUProfile()
    
    // Process orders here
    // ...
    
    return nil
}
```

### Analyze the Profile

Run your profiling function with 10,000 orders, then analyze:
```bash
go tool pprof cpu.prof
```

**Interactive pprof commands:**
- `top` - show top functions by CPU time
- `top10` - show top 10 functions
- `list <function>` - show annotated source for a function
- `web` - generate a graph (requires graphviz)
- `png` - save graph as PNG

**Look for:**
- Functions consuming the most CPU time
- Unexpected hotspots
- Opportunities for optimization

---

## Part 3: Memory & Allocation Profiling

Implement `ProcessOrdersWithMemoryProfile` in `profiling/profiling.go`:
1. Process all orders
2. Create memory profile file
3. Write heap profile with `pprof.WriteHeapProfile()`

**Example structure:**
```go
func ProcessOrdersWithMemoryProfile(orders []*domain.Order, profilePath string) error {
    // Process orders
    // ...
    
    f, err := os.Create(profilePath)
    if err != nil {
        return err
    }
    defer f.Close()
    
    if err := pprof.WriteHeapProfile(f); err != nil {
        return err
    }
    
    return nil
}
```

Analyze with:
```bash
go tool pprof mem.prof
```

### Approach 2: Benchmark Memory Tracking

Run benchmarks with memory tracking:
```bash
go test -bench=. -benchmem ./profiling
```

This shows allocations per operation without generating profile files.

### Reducing Allocations with sync.Pool

When profiling reveals excessive allocations, consider using `sync.Pool` for frequently allocated objects:

```go
import "sync"

var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data string) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset() // Clear before reuse
    
    buf.WriteString("Processed: ")
    buf.WriteString(data)
    result := buf.String()
    
    bufferPool.Put(buf) // Return to pool
    return result
}
```

**When to use sync.Pool:**
- High allocation rate in hot paths
- Objects are expensive to create
- Object reuse is safe (no lingering state)

---

## Part 4: HTTP pprof

HTTP pprof enables live profiling of running services without stopping them. This matches production monitoring scenarios.

### Enable HTTP pprof

Add the pprof HTTP handler to your server. Two approaches:

**Approach 1: Use existing HTTP server**
```go
import _ "net/http/pprof"  // Registers handlers automatically
```

This registers pprof endpoints on your existing HTTP mux.

**Approach 2: Separate debug server** (recommended for production)
```go
import (
    "net/http"
    _ "net/http/pprof"
)

func init() {
    go func() {
        log.Println("Starting pprof server on :6060")
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

### Access Profiles

Start your server and access profiles:

**Web interface:**
```
http://localhost:6060/debug/pprof/
```

**Profile types available:**
- `/debug/pprof/heap` - memory allocations
- `/debug/pprof/goroutine` - goroutine stack traces
- `/debug/pprof/profile` - CPU profile (30s default)
- `/debug/pprof/allocs` - allocation sampling
- `/debug/pprof/block` - blocking profile
- `/debug/pprof/mutex` - mutex contention

**Analyze with go tool pprof:**
```bash
# CPU profile (samples for 30 seconds)
go tool pprof http://localhost:6060/debug/pprof/profile

# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

**When to use HTTP pprof vs file-based:**
- **HTTP pprof**: Production systems, live debugging, continuous monitoring
- **File-based**: Benchmarks, controlled tests, CI/CD pipelines

---

## Part 5 (Optional)

Now that you have profiling tools, use them to find and fix a real bottleneck.

### Scenario: Slow GetAll Endpoint

The `GetAll` repository method becomes slow with large datasets (1000+ orders). Use profiling to diagnose and optimize.

### Steps

1. **Create test with large dataset**
   - Populate repository with 1000+ orders
   - Benchmark `GetAll` operation
   - Record baseline performance

2. **Profile the operation**
   - Use CPU profiling to find hot paths
   - Use memory profiling to find allocation hotspots
   - Identify the bottleneck (likely: slice copying or JSON serialization)

3. **Apply optimization**
   - Pre-allocate slices with known capacity
   - Reuse encoders/decoders
   - Consider using pointers to reduce copying

4. **Measure improvement**
   - Re-run benchmarks
   - Compare before/after metrics
   - Document the improvement (e.g., "50% faster, 30% fewer allocations")

**Example optimization - pre-allocation:**
```go
// Before
func (r *Repository) GetAll() ([]Order, error) {
    var orders []Order
    for _, order := range r.data {
        orders = append(orders, order) // Grows slice dynamically
    }
    return orders, nil
}

// After
func (r *Repository) GetAll() ([]Order, error) {
    orders := make([]Order, 0, len(r.data)) // Pre-allocate capacity
    for _, order := range r.data {
        orders = append(orders, order)
    }
    return orders, nil
}
```

This part is optional because the focus of this lab is learning profiling tools, not optimization techniques.