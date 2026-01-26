package profiling

import (
	"context"
	"fmt"
	"testing"
	"time"

	"lab09/internal/domain"
	"lab09/internal/repository"
)

// BenchmarkOrderValidation benchmarks order validation logic.
func BenchmarkOrderValidation(b *testing.B) {
	order := &domain.Order{
		ID:         "ORD-001",
		CustomerID: "CUST-001",
		Items: []domain.LineItem{
			{
				ProductID:   "PROD-001",
				ProductName: "Test Product 1",
				Quantity:    2,
				UnitPrice:   19.99,
			},
			{
				ProductID:   "PROD-002",
				ProductName: "Test Product 2",
				Quantity:    1,
				UnitPrice:   29.99,
			},
		},
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = order.Validate()
	}
}

// BenchmarkCalculateTotal benchmarks total calculation.
func BenchmarkCalculateTotal(b *testing.B) {
	order := &domain.Order{
		ID:         "ORD-001",
		CustomerID: "CUST-001",
		Items: []domain.LineItem{
			{ProductID: "P1", ProductName: "Product 1", Quantity: 2, UnitPrice: 10.00},
			{ProductID: "P2", ProductName: "Product 2", Quantity: 3, UnitPrice: 15.00},
			{ProductID: "P3", ProductName: "Product 3", Quantity: 1, UnitPrice: 25.00},
			{ProductID: "P4", ProductName: "Product 4", Quantity: 5, UnitPrice: 5.00},
		},
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		order.CalculateTotal()
	}
}

// BenchmarkRepositoryCreate benchmarks creating orders in the repository.
func BenchmarkRepositoryCreate(b *testing.B) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
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
}

// BenchmarkRepositoryGet benchmarks retrieving orders by ID.
func BenchmarkRepositoryGet(b *testing.B) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// Pre-populate with 1000 orders
	for i := 0; i < 1000; i++ {
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		orderID := fmt.Sprintf("ORD-%d", i%1000)
		_, _ = repo.Get(ctx, orderID)
	}
}

// BenchmarkRepositoryGetAll benchmarks retrieving all orders.
// TODO: Part 1 - Implement this benchmark
// Hint: Pre-populate with 100 orders, then call GetAll in the loop
func BenchmarkRepositoryGetAll(b *testing.B) {
	// TODO: Create repository and context
	// TODO: Pre-populate with 100 orders
	// TODO: Call b.ResetTimer()
	// TODO: Loop b.N times calling repo.GetAll(ctx)
}

// BenchmarkServiceCreateOrder benchmarks the full service layer create operation.
// TODO: Part 1 - Implement this benchmark
// Hint: Use service.CreateOrder which includes validation and repository operations
func BenchmarkServiceCreateOrder(b *testing.B) {
	// TODO: Create repository and service
	// TODO: Create context
	// TODO: Call b.ResetTimer()
	// TODO: Loop b.N times creating orders with unique IDs through service
}
