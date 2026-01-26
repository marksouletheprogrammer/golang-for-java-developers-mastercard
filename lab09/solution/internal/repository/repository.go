package repository

import (
	"context"
	"errors"

	"lab09/internal/domain"
)

var (
	// ErrNotFound indicates the requested order doesn't exist.
	ErrNotFound = errors.New("order not found")

	// ErrAlreadyExists indicates an order with this ID already exists.
	ErrAlreadyExists = errors.New("order already exists")
)

// OrderRepository defines the interface for order data access.
// Using an interface allows us to swap implementations (in-memory, database, etc.)
// without changing business logic. This is the Repository pattern.
type OrderRepository interface {
	// Create stores a new order. Returns ErrAlreadyExists if ID is duplicate.
	Create(ctx context.Context, order *domain.Order) error

	// Get retrieves an order by ID. Returns ErrNotFound if it doesn't exist.
	Get(ctx context.Context, id string) (*domain.Order, error)

	// GetAll returns all orders. Empty slice if none exist.
	GetAll(ctx context.Context) ([]*domain.Order, error)

	// Update replaces an existing order. Returns ErrNotFound if it doesn't exist.
	Update(ctx context.Context, order *domain.Order) error

	// UpdateStatus changes only the status field. Returns ErrNotFound if it doesn't exist.
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error

	// Delete removes an order. Returns ErrNotFound if it doesn't exist.
	Delete(ctx context.Context, id string) error
}
