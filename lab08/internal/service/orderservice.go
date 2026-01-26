package service

import (
	"context"
	"errors"

	"lab08/internal/domain"
	"lab08/internal/repository"
)

var (
	// ErrInvalidOrder indicates order validation failed.
	ErrInvalidOrder = errors.New("invalid order")
	
	// ErrInvalidStatusTransition indicates status change not allowed.
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// OrderService contains business logic for order operations.
// Depends on repository interface (not concrete implementation) for flexibility.
// This is dependency injection - repository is injected via constructor.
type OrderService struct {
	repo repository.OrderRepository
}

// NewOrderService creates a new order service with the given repository.
func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

// CreateOrder validates and creates a new order.
// Business logic: validates order, calculates total, sets initial status.
// TODO: Part 2 - Implement CreateOrder
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	// TODO: Validate the order using order.Validate()
	// TODO: If validation fails, wrap error with ErrInvalidOrder using fmt.Errorf("%w: %v", ...)
	// TODO: Call order.CalculateTotal() to set the total amount
	// TODO: If order.Status is empty, set it to domain.StatusPending
	// TODO: Call s.repo.Create() to persist the order
	// TODO: Return any errors from the repository
	return nil
}

// GetOrder retrieves an order by ID.
// TODO: Part 2 - Implement GetOrder
func (s *OrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return nil, nil
}

// ListOrders returns all orders.
// TODO: Part 2 - Implement ListOrders
func (s *OrderService) ListOrders(ctx context.Context) ([]*domain.Order, error) {
	return nil, nil
}

// UpdateOrderStatus changes order status with validation.
// Business logic: checks if state transition is valid before updating.
// TODO: Part 2 - Implement UpdateOrderStatus
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, newStatus domain.OrderStatus) error {
	// TODO: Get the current order using s.repo.Get()
	// TODO: Check if the status transition is valid using order.CanTransitionTo()
	// TODO: If not valid, return ErrInvalidStatusTransition with a descriptive message
	// TODO: Update the status using s.repo.UpdateStatus()
	// TODO: Return any errors
	return nil
}

// CalculateOrderTotal recalculates and returns the total for an order.
// Useful if prices change or items are modified.
func (s *OrderService) CalculateOrderTotal(ctx context.Context, id string) (float64, error) {
	order, err := s.repo.Get(ctx, id)
	if err != nil {
		return 0, err
	}
	
	order.CalculateTotal()
	return order.TotalAmount, nil
}

// DeleteOrder removes an order.
// TODO: Part 2 - Implement DeleteOrder
func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return nil
}
