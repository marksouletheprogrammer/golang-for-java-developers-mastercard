package service

import (
	"context"
	"errors"
	"fmt"

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
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	// Validate order meets business rules
	if err := order.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidOrder, err)
	}
	
	// Calculate total from line items
	order.CalculateTotal()
	
	// Set initial status if not set
	if order.Status == "" {
		order.Status = domain.StatusPending
	}
	
	// Persist to repository
	return s.repo.Create(ctx, order)
}

// GetOrder retrieves an order by ID.
func (s *OrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.Get(ctx, id)
}

// ListOrders returns all orders.
func (s *OrderService) ListOrders(ctx context.Context) ([]*domain.Order, error) {
	return s.repo.GetAll(ctx)
}

// UpdateOrderStatus changes order status with validation.
// Business logic: checks if state transition is valid before updating.
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, newStatus domain.OrderStatus) error {
	// Get current order
	order, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	
	// Validate status transition
	if !order.CanTransitionTo(newStatus) {
		return fmt.Errorf("%w: cannot transition from %s to %s", 
			ErrInvalidStatusTransition, order.Status, newStatus)
	}
	
	// Update status in repository
	return s.repo.UpdateStatus(ctx, id, newStatus)
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
func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
