package repository

import (
	"context"
	"sync"
	"time"

	"lab10/internal/domain"
)

// MemoryRepository is an in-memory implementation of OrderRepository.
// Uses a map with mutex for thread-safe concurrent access.
// Good for testing and demos, not for production (data lost on restart).
type MemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orders: make(map[string]*domain.Order),
	}
}

// Create stores a new order in memory.
func (r *MemoryRepository) Create(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; exists {
		return ErrAlreadyExists
	}

	// Store a copy to prevent external modification
	orderCopy := *order
	orderCopy.CreatedAt = time.Now()
	orderCopy.UpdatedAt = time.Now()
	r.orders[order.ID] = &orderCopy

	return nil
}

// Get retrieves an order by ID.
func (r *MemoryRepository) Get(ctx context.Context, id string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Return a copy to prevent external modification
	orderCopy := *order
	return &orderCopy, nil
}

// GetAll returns all orders.
func (r *MemoryRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orderCopy := *order
		orders = append(orders, &orderCopy)
	}

	return orders, nil
}

// Update replaces an existing order.
func (r *MemoryRepository) Update(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return ErrNotFound
	}

	orderCopy := *order
	orderCopy.UpdatedAt = time.Now()
	r.orders[order.ID] = &orderCopy

	return nil
}

// UpdateStatus changes only the order status.
func (r *MemoryRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[id]
	if !exists {
		return ErrNotFound
	}

	order.Status = status
	order.UpdatedAt = time.Now()

	return nil
}

// Delete removes an order.
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[id]; !exists {
		return ErrNotFound
	}

	delete(r.orders, id)
	return nil
}
