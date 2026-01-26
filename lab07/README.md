# Lab 7: Advanced Concurrency Patterns

You'll implement advanced concurrency patterns for processing high volumes of payment transactions: worker pools, fan-out/fan-in, select statements, and context-based cancellation.

## The Scenario

Your payment processing system receives batches of thousands of transactions that need to be validated, scored for fraud risk, and recorded. Sequential processing would create unacceptable delays. A worker pool with controlled concurrency is the solution.

**Starter files provided:** `main.go` with substantial implementation including `Transaction` struct and test data generation, plus `main_test.go`.

## Assignment

### Part 1: Review the Test Data
1. In this lab, review the provided `main.go` file.
2. Examine the `Transaction` struct with fields: ID, amount, currency, merchantID, cardNumber (last 4 digits).
3. Review the function that generates test data with transaction records.
4. Note the invalid transactions included (negative amounts, invalid currencies).

### Part 2: Review Sequential Processing Baseline
1. Review the provided `ProcessTransactionsSequential` function (already implemented as a baseline for comparison)
2. Note the fraud scoring simulation with 50-100ms delay
3. Observe how it tracks success and error counts
4. Run with 100 transactions and measure execution time

### Part 3: Basic Worker Pool
1. Implement `ProcessTransactionsWithWorkerPool` that accepts a configurable number of workers
2. Use channels for jobs and results
3. Launch worker goroutines that process transactions from the jobs channel
4. Use a WaitGroup to coordinate shutdown
5. Test with different worker counts (1, 5, 10, 50, 100) to find the optimal number

**Example worker pool structure:**
```go
func ProcessTransactionsWithWorkerPool(transactions []Transaction, numWorkers int) []Result {
    jobs := make(chan Transaction, len(transactions))
    results := make(chan Result, len(transactions))
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(jobs, results, &wg)
    }
    
    // Send jobs
    for _, txn := range transactions {
        jobs <- txn
    }
    close(jobs) // Signal no more jobs
    
    // Wait for workers to finish, then close results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    var allResults []Result
    for result := range results {
        allResults = append(allResults, result)
    }
    
    return allResults
}

func worker(jobs <-chan Transaction, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    for txn := range jobs { // Loop until jobs channel is closed
        result := processTransaction(txn) // Your processing logic
        results <- result
    }
}
```
The worker pool limits concurrency to a fixed number of goroutines, preventing resource exhaustion.

### Part 4: Fan-Out/Fan-In Pattern
1. Implement a `FanOut` function that distributes work across multiple worker channels
2. Implement a `FanIn` function that merges results from multiple channels into one
3. Refactor your worker pool to use these patterns

### Part 5: Select Statement for Multiplexing
1. Update your worker pool to support graceful shutdown, timeouts, and progress updates
2. Use select statements to handle multiple channels
3. Test normal completion, timeout, and manual cancellation scenarios

### Part 6: Context-Based Cancellation (Optional)
1. Update your worker pool to accept a `context.Context`
2. Workers should check `context.Done()` and stop processing when cancelled
3. Test with context timeout and manual cancellation

### Part 7: Bounded Concurrency with Semaphore Pattern (Optional)
1. Implement a semaphore using a buffered channel to limit concurrent fraud scoring calls
2. Test with different semaphore sizes
3. Compare to unlimited concurrency

### Part 8: Performance Comparison (Optional)
1. Benchmark all implementations with various configurations
2. Measure execution time, throughput, and resource usage
3. Create a report showing optimal configurations

### Part 9: Java Comparison (Optional)
Consider how Go's worker pools, channels, and context compare to Java's ExecutorService, BlockingQueue, and Future patterns.