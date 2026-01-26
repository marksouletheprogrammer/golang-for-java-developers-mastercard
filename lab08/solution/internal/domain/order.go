package domain

import (
	"errors"
	"time"
)

// OrderStatus represents the current state of an order in its lifecycle.
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

// LineItem represents a single item in an order with quantity and price.
type LineItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

// Subtotal calculates the total price for this line item (quantity * unit price).
func (li LineItem) Subtotal() float64 {
	return float64(li.Quantity) * li.UnitPrice
}

// Order represents a customer order with line items and status tracking.
// This is the core domain entity - no database or transport concerns here.
type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	Items       []LineItem  `json:"items"`
	Status      OrderStatus `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// CalculateTotal computes the total amount from all line items.
// Called when order is created or items are modified.
func (o *Order) CalculateTotal() {
	total := 0.0
	for _, item := range o.Items {
		total += item.Subtotal()
	}
	o.TotalAmount = total
}

// Validate checks if the order meets business rules.
// Returns error if validation fails.
func (o *Order) Validate() error {
	if o.CustomerID == "" {
		return errors.New("customer ID is required")
	}
	if len(o.Items) == 0 {
		return errors.New("order must have at least one item")
	}
	for i, item := range o.Items {
		if item.ProductID == "" {
			return errors.New("product ID is required for all items")
		}
		if item.Quantity <= 0 {
			return errors.New("quantity must be positive")
		}
		if item.UnitPrice < 0 {
			return errors.New("unit price cannot be negative")
		}
		// Validate product name
		if item.ProductName == "" {
			return errors.New("product name is required")
		}
		_ = i // Suppress unused variable warning if needed
	}
	return nil
}

// CanTransitionTo checks if order can transition to a new status.
// Enforces valid state transitions (e.g., can't ship a cancelled order).
func (o *Order) CanTransitionTo(newStatus OrderStatus) bool {
	switch o.Status {
	case StatusPending:
		return newStatus == StatusConfirmed || newStatus == StatusCancelled
	case StatusConfirmed:
		return newStatus == StatusShipped || newStatus == StatusCancelled
	case StatusShipped:
		return newStatus == StatusDelivered
	case StatusDelivered, StatusCancelled:
		return false // Terminal states
	default:
		return false
	}
}
