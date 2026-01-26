package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Transaction represents a payment transaction to be processed.
type Transaction struct {
	ID         int
	Amount     float64
	Currency   string
	MerchantID string
	CardNumber string // Last 4 digits only
}

// ProcessingResult contains the outcome of processing a transaction.
type ProcessingResult struct {
	Transaction Transaction
	FraudScore  float64 // 0.0 to 1.0, higher = more suspicious
	Valid       bool
	Error       error
}

// GenerateTestData creates test transactions with some intentionally invalid records.
// Every 20th transaction has a negative amount, every 25th has an invalid currency.
func GenerateTestData(count int) []Transaction {
	transactions := make([]Transaction, count)
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CAD"}
	invalidCurrencies := []string{"XYZ", "ABC"}

	for i := 0; i < count; i++ {
		amount := rand.Float64() * 10000
		currency := currencies[rand.Intn(len(currencies))]

		// Inject invalid data for testing error handling
		if i%20 == 0 {
			amount = -100.0
		}
		if i%25 == 0 {
			currency = invalidCurrencies[rand.Intn(len(invalidCurrencies))]
		}

		transactions[i] = Transaction{
			ID:         i + 1,
			Amount:     amount,
			Currency:   currency,
			MerchantID: fmt.Sprintf("MERCHANT_%d", rand.Intn(100)),
			CardNumber: fmt.Sprintf("****%04d", rand.Intn(10000)),
		}
	}

	return transactions
}

// ValidateTransaction checks if a transaction meets basic business rules.
func ValidateTransaction(tx Transaction) error {
	if tx.Amount < 0 {
		return fmt.Errorf("invalid amount: %f", tx.Amount)
	}
	validCurrencies := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true, "CAD": true,
	}
	if !validCurrencies[tx.Currency] {
		return fmt.Errorf("invalid currency: %s", tx.Currency)
	}
	return nil
}

// CalculateFraudScore simulates fraud detection with 50-100ms delay (I/O-bound operation).
// Higher amounts increase fraud score to simulate risk-based scoring.
func CalculateFraudScore(tx Transaction) float64 {
	// Simulate external API call delay
	delay := 50 + rand.Intn(51)
	time.Sleep(time.Duration(delay) * time.Millisecond)

	score := rand.Float64()
	if tx.Amount > 5000 {
		score += 0.3
	}
	if score > 1.0 {
		score = 1.0
	}
	return score
}

// NOTE: This baseline implementation is provided for performance comparison.
// Your task is to implement the concurrent patterns below that improve upon this sequential approach.
//
// ProcessTransactionsSequential processes transactions one at a time (baseline for comparison).
// This is simple but slow - each transaction takes 50-100ms, so 1000 tx = ~75 seconds.
func ProcessTransactionsSequential(transactions []Transaction) (success, errors int, results []ProcessingResult) {
	results = make([]ProcessingResult, 0, len(transactions))

	for _, tx := range transactions {
		err := ValidateTransaction(tx)
		if err != nil {
			errors++
			results = append(results, ProcessingResult{
				Transaction: tx,
				Valid:       false,
				Error:       err,
			})
			continue
		}

		fraudScore := CalculateFraudScore(tx)
		success++
		results = append(results, ProcessingResult{
			Transaction: tx,
			FraudScore:  fraudScore,
			Valid:       true,
			Error:       nil,
		})
	}

	return success, errors, results
}

// ProcessTransactionsWithWorkerPool implements the worker pool pattern.
// Pattern: Fixed number of goroutines process jobs from a shared channel.
// - Jobs channel: Distributes work to workers
// - Results channel: Collects processed results
// - WaitGroup: Coordinates worker lifecycle and shutdown
// Optimal workers: 10 (found through testing) - more doesn't help for I/O-bound workload.
// TODO: Part 3 - Implement worker pool pattern
func ProcessTransactionsWithWorkerPool(transactions []Transaction, numWorkers int) (success, errors int, results []ProcessingResult) {
	// TODO: Create jobs channel (buffered with transaction count)
	// TODO: Create results channel (buffered with transaction count)
	// TODO: Create WaitGroup
	// TODO: Launch numWorkers goroutines
	// TODO: Each worker: read from jobs, validate, calculate fraud score, send to results
	// TODO: Send all transactions to jobs channel, then close it
	// TODO: Launch goroutine to wait for workers and close results channel
	// TODO: Collect all results and count successes/errors
	return 0, 0, nil
}

// FanOut distributes work across multiple channels using round-robin.
// Pattern: Split work from one source to multiple channels for parallel processing.
// Use case: Heterogeneous workloads, load balancing across workers.
// TODO: Part 4 - Implement fan-out pattern
func FanOut(transactions []Transaction, numWorkers int) []chan Transaction {
	// TODO: Create slice of channels (one per worker)
	// TODO: Launch goroutine to distribute transactions round-robin
	// TODO: Close all channels when done
	return nil
}

// FanIn merges results from multiple channels into one output channel.
// Pattern: Combine results from parallel workers into a single stream.
// Uses goroutines per input channel to avoid blocking.
// TODO: Part 4 - Implement fan-in pattern
func FanIn(channels ...<-chan ProcessingResult) <-chan ProcessingResult {
	// TODO: Create output channel
	// TODO: Create WaitGroup
	// TODO: For each input channel, launch goroutine to forward results to output
	// TODO: Launch goroutine to wait for all and close output channel
	return nil
}

// ProcessTransactionsWithFanOutFanIn uses fan-out/fan-in for distributed processing.
// Pattern: Distribute work (fan-out) -> Process in parallel -> Merge results (fan-in).
// More complex than basic worker pool, best for heterogeneous workloads.
// TODO: Part 4 - Implement fan-out/fan-in processing
func ProcessTransactionsWithFanOutFanIn(transactions []Transaction, numWorkers int) (success, errors int, results []ProcessingResult) {
	// TODO: Use FanOut to distribute transactions
	// TODO: Create result channels for each worker
	// TODO: Launch workers to process from each job channel
	// TODO: Use FanIn to merge results
	// TODO: Collect and count results
	return 0, 0, nil
}

// ProcessTransactionsWithSelect demonstrates select statement for multiplexing.
// Pattern: Use select to handle multiple channel operations simultaneously.
// Features:
// - Progress updates: Non-blocking progress reporting
// - Timeout handling: Graceful termination on timeout
// - Manual cancellation: Done channel for explicit shutdown
// Select chooses the first ready channel operation (non-deterministic if multiple ready).
// TODO: Part 5 - Implement select-based processing with progress and timeout
func ProcessTransactionsWithSelect(transactions []Transaction, numWorkers int, timeout time.Duration) (success, errors int, results []ProcessingResult, timedOut bool) {
	// TODO: Create jobs, results, done, and progress channels
	// TODO: Launch workers using select to handle jobs or done signal
	// TODO: Send progress updates (non-blocking)
	// TODO: In main loop, use select to handle: results, progress, or timeout
	// TODO: Return early if timeout occurs
	return 0, 0, nil, false
}

// ProcessTransactionsWithContext implements context-based cancellation.
// Pattern: Use context.Context for cancellation, timeouts, and deadline propagation.
// Benefits:
// - Standard Go pattern for cancellation
// - Automatic propagation through call hierarchy
// - Built-in deadline/timeout support
// - No goroutine leaks on cancellation
// Workers check ctx.Done() to detect cancellation and stop processing.
// TODO: Part 6 - Implement context-based cancellation
func ProcessTransactionsWithContext(ctx context.Context, transactions []Transaction, numWorkers int) (success, errors int, results []ProcessingResult, cancelled bool) {
	// TODO: Create jobs and results channels
	// TODO: Create WaitGroup
	// TODO: Launch workers that check ctx.Done() before and after expensive operations
	// TODO: Send jobs with select on ctx.Done()
	// TODO: Collect results with select on ctx.Done()
	// TODO: Detect if cancellation occurred
	return 0, 0, nil, false
}

// ProcessTransactionsWithSemaphore limits concurrent operations using a semaphore.
// Pattern: Buffered channel as semaphore to bound concurrency.
// - Buffered channel capacity = max concurrent operations
// - Send to acquire, receive to release
// - Blocks when semaphore is full (all permits taken)
// Use case: Rate-limiting external API calls (e.g., fraud scoring service).
// Trade-off: Lower semaphore = less load on external service, lower throughput.
// TODO: Part 7 - Implement semaphore-based rate limiting
func ProcessTransactionsWithSemaphore(transactions []Transaction, numWorkers int, semaphoreSize int) (success, errors int, results []ProcessingResult) {
	// TODO: Create jobs, results, and semaphore channels
	// TODO: Semaphore channel buffer size = semaphoreSize
	// TODO: Workers acquire semaphore before CalculateFraudScore, release after
	// TODO: This limits concurrent fraud score calculations
	return 0, 0, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("=== Lab 7: Advanced Concurrency Patterns ===\n")

	fmt.Println("Part 1: Generating Test Data...")
	transactions := GenerateTestData(1000)
	fmt.Printf("Generated %d transactions\n\n", len(transactions))

	testTransactions := transactions[:100]

	fmt.Println("Part 2: Sequential Processing Baseline")
	start := time.Now()
	success, errors, _ := ProcessTransactionsSequential(testTransactions)
	duration := time.Since(start)
	fmt.Printf("Sequential: %d success, %d errors, Duration: %v\n", success, errors, duration)
	fmt.Printf("Throughput: %.2f tx/sec\n\n", float64(len(testTransactions))/duration.Seconds())

	fmt.Println("Part 3: Basic Worker Pool")
	workerCounts := []int{1, 5, 10, 50, 100}
	for _, workers := range workerCounts {
		start = time.Now()
		success, errors, _ = ProcessTransactionsWithWorkerPool(testTransactions, workers)
		duration = time.Since(start)
		fmt.Printf("Workers: %3d | %d success, %d errors | Duration: %v | Throughput: %.2f tx/sec\n",
			workers, success, errors, duration, float64(len(testTransactions))/duration.Seconds())
	}
	fmt.Println()

	fmt.Println("Part 4: Fan-Out/Fan-In Pattern")
	start = time.Now()
	success, errors, _ = ProcessTransactionsWithFanOutFanIn(testTransactions, 10)
	duration = time.Since(start)
	fmt.Printf("Fan-Out/Fan-In: %d success, %d errors, Duration: %v\n", success, errors, duration)
	fmt.Printf("Throughput: %.2f tx/sec\n\n", float64(len(testTransactions))/duration.Seconds())

	fmt.Println("Part 5: Select Statement for Multiplexing")
	start = time.Now()
	success, errors, _, timedOut := ProcessTransactionsWithSelect(testTransactions, 10, 30*time.Second)
	duration = time.Since(start)
	fmt.Printf("Select Pattern: %d success, %d errors, TimedOut: %v, Duration: %v\n", success, errors, timedOut, duration)
	fmt.Printf("Throughput: %.2f tx/sec\n\n", float64(len(testTransactions))/duration.Seconds())

	fmt.Println("Part 6: Context-Based Cancellation")

	fmt.Println("  - Normal completion:")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()
	start = time.Now()
	success, errors, _, cancelled := ProcessTransactionsWithContext(ctx1, testTransactions, 10)
	duration = time.Since(start)
	fmt.Printf("    %d success, %d errors, Cancelled: %v, Duration: %v\n", success, errors, cancelled, duration)

	fmt.Println("  - With timeout:")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	start = time.Now()
	success, errors, _, cancelled = ProcessTransactionsWithContext(ctx2, testTransactions, 10)
	duration = time.Since(start)
	fmt.Printf("    %d success, %d errors, Cancelled: %v, Duration: %v\n\n", success, errors, cancelled, duration)

	fmt.Println("Part 7: Bounded Concurrency with Semaphore Pattern")
	semaphoreSizes := []int{5, 10, 20, 50}
	for _, semSize := range semaphoreSizes {
		start = time.Now()
		success, errors, _ = ProcessTransactionsWithSemaphore(testTransactions, 10, semSize)
		duration = time.Since(start)
		fmt.Printf("Semaphore: %2d | %d success, %d errors | Duration: %v | Throughput: %.2f tx/sec\n",
			semSize, success, errors, duration, float64(len(testTransactions))/duration.Seconds())
	}
	fmt.Println()

	fmt.Println("Part 8: Performance Comparison (1000 transactions - concurrent patterns only)")
	fullTransactions := transactions

	fmt.Println("\n1. Worker Pool (Optimal: 10 workers):")
	start = time.Now()
	success, errors, _ = ProcessTransactionsWithWorkerPool(fullTransactions, 10)
	duration = time.Since(start)
	fmt.Printf("   Duration: %v | Throughput: %.2f tx/sec | Success: %d | Errors: %d\n",
		duration, float64(len(fullTransactions))/duration.Seconds(), success, errors)

	fmt.Println("\n2. Fan-Out/Fan-In (10 workers):")
	start = time.Now()
	success, errors, _ = ProcessTransactionsWithFanOutFanIn(fullTransactions, 10)
	duration = time.Since(start)
	fmt.Printf("   Duration: %v | Throughput: %.2f tx/sec | Success: %d | Errors: %d\n",
		duration, float64(len(fullTransactions))/duration.Seconds(), success, errors)

	fmt.Println("\n3. Context-Based (10 workers):")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel3()
	start = time.Now()
	success, errors, _, _ = ProcessTransactionsWithContext(ctx3, fullTransactions, 10)
	duration = time.Since(start)
	fmt.Printf("   Duration: %v | Throughput: %.2f tx/sec | Success: %d | Errors: %d\n",
		duration, float64(len(fullTransactions))/duration.Seconds(), success, errors)

	fmt.Println("\n4. Semaphore (10 workers, semaphore: 20):")
	start = time.Now()
	success, errors, _ = ProcessTransactionsWithSemaphore(fullTransactions, 10, 20)
	duration = time.Since(start)
	fmt.Printf("   Duration: %v | Throughput: %.2f tx/sec | Success: %d | Errors: %d\n",
		duration, float64(len(fullTransactions))/duration.Seconds(), success, errors)

	fmt.Println("\n  Note: Sequential processing of 1000 tx would take ~75 seconds")
	fmt.Println("  (already demonstrated in Part 2 with 100 tx)")

	fmt.Println("\n=== Performance Report ===")
	fmt.Println("Optimal Configuration: Worker Pool with 10 workers")
	fmt.Println("Key Findings:")
	fmt.Println("- Worker pools provide ~10x speedup over sequential processing")
	fmt.Println("- Optimal worker count depends on workload (I/O-bound vs CPU-bound)")
	fmt.Println("- Context provides clean cancellation without resource leaks")
	fmt.Println("- Semaphores useful for rate-limiting external service calls")
	fmt.Println("- Fan-out/fan-in adds complexity but good for heterogeneous workloads")

	fmt.Println("\nLab 7 Complete!")
}
