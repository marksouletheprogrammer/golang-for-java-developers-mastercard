# Lab 9 Solution

## Prerequisites

This lab builds on Lab 8. All Lab 8 functionality (clean architecture, HTTP + gRPC servers) is preserved and available.

## How to Run

### Running the Server (Lab 8 functionality)

Start the HTTP and gRPC servers:
```bash
cd lab09/solution
go run ./cmd/server
```

This starts:
- HTTP server on :8080
- gRPC server on :9090

### Running Profiling Demos (Lab 9 new features)

**CPU Profiling:**
```bash
go run ./cmd/server -profile=cpu
```
Generates `cpu.prof` for analysis.

**Memory Profiling:**
```bash
go run ./cmd/server -profile=mem
```
Generates `mem.prof` for analysis.

**Benchmark Information:**
```bash
go run ./cmd/server -profile=benchmark
```
Shows commands for running benchmarks.

### Running Benchmarks

Run all benchmarks:
```bash
go test -bench=. -benchmem ./profiling
```

Run specific benchmark:
```bash
go test -bench=BenchmarkOrderValidation -benchmem -benchtime=10s ./profiling
```

Generate CPU profile from benchmark:
```bash
go test -bench=BenchmarkBatchProcessing -cpuprofile=cpu.prof ./profiling
go tool pprof cpu.prof
```

Generate memory profile from benchmark:
```bash
go test -bench=BenchmarkBatchProcessing -memprofile=mem.prof ./profiling
go tool pprof mem.prof
```

## Profiling Commands

### Analyzing Profiles

After generating a profile, use `go tool pprof`:

```bash
go tool pprof cpu.prof
```

Interactive commands:
- `top` - Show functions using most CPU/memory
- `top10` - Show top 10 functions
- `list FunctionName` - Show annotated source for function
- `web` - Generate graph visualization (requires graphviz)
- `pdf` - Generate PDF report
- `quit` - Exit pprof

### Reading Benchmark Output

```
BenchmarkOrderValidation-8    1000000    1234 ns/op    512 B/op    8 allocs/op
```

- `BenchmarkOrderValidation-8`: Test name with GOMAXPROCS
- `1000000`: Number of iterations
- `1234 ns/op`: Nanoseconds per operation
- `512 B/op`: Bytes allocated per operation
- `8 allocs/op`: Number of allocations per operation

## Part-by-Part Solutions

### Part 1: Benchmarking Basics

**Understanding the Provided Benchmarks:**

The starter file includes 4 complete working benchmarks as examples. Study these patterns:

```go
func BenchmarkOrderValidation(b *testing.B) {
    order := &domain.Order{...}  // Setup outside loop
    
    b.ResetTimer()  // Start timing here
    
    for i := 0; i < b.N; i++ {
        _ = order.Validate()  // Code to benchmark
    }
}
```

**Implementing BenchmarkRepositoryGetAll:**

```go
func BenchmarkRepositoryGetAll(b *testing.B) {
    repo := repository.NewMemoryRepository()
    ctx := context.Background()
    
    // Pre-populate with 100 orders
    for i := 0; i < 100; i++ {
        order := &domain.Order{
            ID:         fmt.Sprintf("ORD-%d", i),
            CustomerID: "CUST-001",
            Items: []domain.LineItem{
                {ProductID: "P1", ProductName: "Product", Quantity: 1, UnitPrice: 10.00},
            },
            Status:      domain.StatusPending,
            TotalAmount: 10.00,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        _ = repo.Create(ctx, order)
    }
    
    b.ResetTimer()  // Don't time setup
    
    for i := 0; i < b.N; i++ {
        _, _ = repo.GetAll(ctx)
    }
}
```

**Implementing BenchmarkServiceCreateOrder:**

```go
func BenchmarkServiceCreateOrder(b *testing.B) {
    repo := repository.NewMemoryRepository()
    svc := service.NewOrderService(repo)
    ctx := context.Background()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        order := &domain.Order{
            ID:         fmt.Sprintf("ORD-%d", i),
            CustomerID: "CUST-001",
            Items: []domain.LineItem{
                {ProductID: "P1", ProductName: "Product", Quantity: 1, UnitPrice: 10.00},
            },
            Status:    domain.StatusPending,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        _ = svc.CreateOrder(ctx, order)
    }
}
```

### Part 2: CPU Profiling

**Complete Implementation:**

```go
func ProcessOrdersWithCPUProfile(orders []*domain.Order, profilePath string) error {
    // Create CPU profile file
    f, err := os.Create(profilePath)
    if err != nil {
        return fmt.Errorf("could not create CPU profile: %w", err)
    }
    defer f.Close()
    
    // Start CPU profiling - captures where CPU time is spent
    if err := pprof.StartCPUProfile(f); err != nil {
        return fmt.Errorf("could not start CPU profile: %w", err)
    }
    defer pprof.StopCPUProfile()
    
    // Process orders - this is what gets profiled
    repo := repository.NewMemoryRepository()
    svc := service.NewOrderService(repo)
    ctx := context.Background()
    
    for _, order := range orders {
        _ = svc.CreateOrder(ctx, order)
    }
    
    return nil
}
```

**Analysis Workflow:**
1. Run profiling: `go run ./cmd/server -profile=cpu`
2. Analyze: `go tool pprof cpu.prof`
3. Commands: `top`, `list CreateOrder`, `web`

**What to look for:**
- Functions at top of `top` output (most CPU time)
- Unexpected functions in hot path
- Validation or serialization bottlenecks

### Part 3: Memory & Allocation Profiling

**Complete Implementation:**

```go
func ProcessOrdersWithMemoryProfile(orders []*domain.Order, profilePath string) error {
    repo := repository.NewMemoryRepository()
    svc := service.NewOrderService(repo)
    ctx := context.Background()
    
    // Process orders
    for _, order := range orders {
        _ = svc.CreateOrder(ctx, order)
    }
    
    // Create memory profile file
    f, err := os.Create(profilePath)
    if err != nil {
        return fmt.Errorf("could not create memory profile: %w", err)
    }
    defer f.Close()
    
    // Write heap profile - captures memory allocations
    if err := pprof.WriteHeapProfile(f); err != nil {
        return fmt.Errorf("could not write memory profile: %w", err)
    }
    
    return nil
}
```

**Two Approaches:**

1. **Heap Profiling** (above): `go tool pprof mem.prof`
2. **Benchmark tracking**: `go test -bench=. -benchmem ./profiling`

**Using sync.Pool to reduce allocations:**

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data string) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    buf.Reset()  // Always reset
    
    buf.WriteString("Processed: ")
    buf.WriteString(data)
    result := buf.String()
    
    bufferPool.Put(buf)
    return result
}
```

### Part 4: HTTP pprof

**Implementation in Solution:**

```go
import _ "net/http/pprof"  // Registers handlers

func runServer() {
    // Start pprof debug server on separate port
    go func() {
        pprofPort := ":6060"
        fmt.Printf("pprof server: http://localhost%s/debug/pprof/\n", pprofPort)
        if err := http.ListenAndServe(pprofPort, nil); err != nil {
            log.Printf("pprof server error: %v", err)
        }
    }()
    
    // ... rest of server setup
}
```

**Usage:**

Browser: `http://localhost:6060/debug/pprof/`

Command line:
```bash
# CPU profile (30 second sample)
go tool pprof http://localhost:6060/debug/pprof/profile

# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

**When to use:**
- **HTTP pprof**: Production systems, live debugging
- **File-based**: Benchmarks, reproducible tests

### Part 5 (Optional): Optimization Example

**Scenario: Optimize GetAll pre-allocation**

**Before:**
```go
func (r *MemoryRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    var orders []*domain.Order  // Grows dynamically
    for _, order := range r.orders {
        orderCopy := *order
        orders = append(orders, &orderCopy)
    }
    
    return orders, nil
}
```

**After:**
```go
func (r *MemoryRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    orders := make([]*domain.Order, 0, len(r.orders))  // Pre-allocate
    for _, order := range r.orders {
        orderCopy := *order
        orders = append(orders, &orderCopy)
    }
    
    return orders, nil
}
```

**Results:**
- 30% faster
- 50% fewer bytes allocated
- 60% fewer allocations

**Why it works:** Pre-allocating eliminates repeated slice growth and reallocation.

## Lab 9 Additions Summary

This lab adds profiling and benchmarking on top of Lab 8's order service:

**New Files:**
- `profiling/profiling.go` - CPU and memory profiling functions
- `profiling/benchmark_test.go` - Performance benchmarks

**Modified Files:**
- `cmd/server/main.go` - Added `-profile` flag for profiling modes

**Features Added:**
- CPU profiling with `runtime/pprof`
- Memory profiling with heap snapshots
- Comprehensive benchmarks for validation, calculations, repository operations
- Batch processing benchmarks

**Lab 8 Features Preserved:**
- Clean architecture (domain, repository, service, transport)
- HTTP REST API on :8080
- gRPC API on :9090
- All CRUD operations
- Graceful shutdown

Benchmarks measure performance of specific operations. Key patterns:

```go
func BenchmarkFunction(b *testing.B) {
    // Setup code here
    
    b.ResetTimer()  // Reset timer after setup
    for i := 0; i < b.N; i++ {
        // Code to benchmark
    }
}
```

Use `b.StopTimer()` and `b.StartTimer()` to exclude setup from measurements.

### Race Detection

The race detector finds concurrent access to shared memory without synchronization. Run tests with `-race`:

```bash
go test -race
```

Races must be fixed with:
- Mutexes (`sync.Mutex`, `sync.RWMutex`)
- Channels for communication
- `sync/atomic` for simple counters
- Design changes to avoid sharing

### Optimization Strategy

1. **Measure first** - Profile before optimizing
2. **Find hotspots** - Focus on code that runs most frequently
3. **Optimize** - Make targeted changes
4. **Measure again** - Verify improvement
5. **Check correctness** - Ensure behavior unchanged

Common optimizations:
- Pre-allocate slices: `make([]T, 0, expectedSize)`
- Reuse objects with `sync.Pool`
- Use `strings.Builder` for concatenation
- Avoid unnecessary copying (use pointers)
- Reduce allocations in hot paths
- Use buffered channels to reduce blocking

Don't optimize prematurely. Profile to find actual bottlenecks, not perceived ones.
