# Lab 7 Solution

## How to Run

```bash
cd lab07/solution
go run main.go
```

Benchmarks:
```bash
go test -bench=. -benchmem
```

## Key Results

**Optimal Configuration**: 10 workers provides ~8-10x speedup over sequential processing

| Pattern | Throughput (100 tx) |
|---------|---------------------|
| Sequential | ~13 tx/sec |
| Worker Pool (10) | ~111 tx/sec |
| All concurrent patterns | Similar performance |

**Why 10 workers?** I/O-bound workload (50-100ms delays) - more workers don't help.

## Go vs Java Concurrency

### Worker Pools

Java:
```java
ExecutorService executor = Executors.newFixedThreadPool(numWorkers);
executor.submit(() -> process(tx));
executor.shutdown();
executor.awaitTermination(timeout, TimeUnit.SECONDS);
```

Go:
```go
for w := 0; w < numWorkers; w++ {
    go func() {
        for tx := range jobs {
            process(tx)
        }
    }()
}
```

No pool object, no explicit shutdown. Just goroutines and channels.

### Passing Work

Java uses `BlockingQueue`, Go uses channels. `for tx := range jobs` replaces repeatedly calling `queue.poll()`.

### Waiting for Completion

Java's `CountDownLatch` equals Go's `sync.WaitGroup`. Both block until workers finish.

### Cancellation

Java has multiple approaches (`Future.cancel()`, volatile flags, interrupts). Go standardizes on `context.Context`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
select {
case <-ctx.Done():
    return
}
```

### Select Statement

Go's `select` waits on multiple channels simultaneously. Java has no equivalent - you'd use `BlockingQueue.poll(timeout)` with multiple try-catch blocks.

### Scale

Goroutines are 2KB, threads are 1MB. Go handles 100,000+ goroutines easily. Java maxes out around 10,000 threads.
