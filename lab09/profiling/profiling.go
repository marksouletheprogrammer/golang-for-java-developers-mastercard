package profiling

import (
	"fmt"
	"time"

	"lab09/internal/domain"
)

// ProcessOrdersWithCPUProfile processes a batch of orders while CPU profiling.
// This demonstrates how to add CPU profiling to identify performance bottlenecks.
// TODO: Part 2 - Implement CPU profiling
func ProcessOrdersWithCPUProfile(orders []*domain.Order, profilePath string) error {
	// TODO: Create CPU profile file using os.Create(profilePath)
	// TODO: Start CPU profiling with pprof.StartCPUProfile(f)
	// TODO: Defer pprof.StopCPUProfile()
	// TODO: Create repository using repository.NewMemoryRepository()
	// TODO: Create service using service.NewOrderService(repo)
	// TODO: Create context using context.Background()
	// TODO: Loop through orders and call svc.CreateOrder(ctx, order) for each
	// TODO: Return any errors with proper error wrapping
	return nil
}

// ProcessOrdersWithMemoryProfile processes orders and creates a memory profile.
// This helps identify allocation hotspots and memory usage patterns.
// TODO: Part 3 - Implement memory profiling
func ProcessOrdersWithMemoryProfile(orders []*domain.Order, profilePath string) error {
	// TODO: Create repository using repository.NewMemoryRepository()
	// TODO: Create service using service.NewOrderService(repo)
	// TODO: Create context using context.Background()
	// TODO: Loop through orders and call svc.CreateOrder(ctx, order) for each
	// TODO: Create memory profile file using os.Create(profilePath)
	// TODO: Write heap profile with pprof.WriteHeapProfile(f)
	// TODO: Return any errors with proper error wrapping
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
