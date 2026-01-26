package profiling

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"lab09/internal/domain"
	"lab09/internal/repository"
	"lab09/internal/service"
)

// ProcessOrdersWithCPUProfile processes a batch of orders while CPU profiling.
// This demonstrates how to add CPU profiling to identify performance bottlenecks.
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

// ProcessOrdersWithMemoryProfile processes orders and creates a memory profile.
// This helps identify allocation hotspots and memory usage patterns.
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
	// This shows which functions allocate the most memory
	if err := pprof.WriteHeapProfile(f); err != nil {
		return fmt.Errorf("could not write memory profile: %w", err)
	}
	
	return nil
}

// GenerateTestOrders creates a batch of test orders for profiling.
func GenerateTestOrders(count int) []*domain.Order {
	orders := make([]*domain.Order, count)
	
	for i := 0; i < count; i++ {
		orders[i] = &domain.Order{
			ID:         fmt.Sprintf("ORD-%06d", i),
			CustomerID: fmt.Sprintf("CUST-%04d", i%1000),
			Items: []domain.LineItem{
				{
					ProductID:   fmt.Sprintf("PROD-%03d", i%100),
					ProductName: fmt.Sprintf("Product %d", i%100),
					Quantity:    (i % 5) + 1,
					UnitPrice:   float64((i%50)+10) * 1.99,
				},
				{
					ProductID:   fmt.Sprintf("PROD-%03d", (i+1)%100),
					ProductName: fmt.Sprintf("Product %d", (i+1)%100),
					Quantity:    (i % 3) + 1,
					UnitPrice:   float64((i%30)+5) * 2.49,
				},
			},
			Status:    domain.StatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	
	return orders
}
