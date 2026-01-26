package main

import (
	"time"
)

// Transaction represents a payment transaction in the system.
type Transaction struct {
	// TODO Part 1
}

// NewTransaction is a constructor function that creates and returns a new Transaction.
// TODO: Part 2 - Implement constructor
func NewTransaction(id string, amount float64, currency string, merchantID string) *Transaction {
	return nil
}

// CalculateFee calculates the transaction fee based on a percentage.
// TODO: Part 3 - Implement fee calculation method
func (t *Transaction) CalculateFee(feePercentage float64) float64 {
	return 0.0
}

// GetDisplayInfo returns a formatted string with transaction details.
// TODO: Part 3 - Implement display info method
func (t *Transaction) GetDisplayInfo() string {
	return ""
}

// ProcessPayment processes the payment transaction.
// TODO: Part 4 - Implement this method to satisfy the Payable interface
func (t *Transaction) ProcessPayment() error {
	// TODO: Validate the transaction (check for positive amount)
	// TODO: Print a message indicating payment is being processed
	// TODO: Return nil on success, or an error if validation fails
	return nil
}

// Payable is an interface that defines the contract for payment processing.
// TODO: Part 4 - Define the Payable interface
type Payable interface {
}

// Filter returns a new slice containing only elements that satisfy the predicate function.
// This is a generic function that works with any type T.
// TODO: Part 6 - Implement generic filter function
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// TODO: Create a new slice to hold filtered results
	// TODO: Iterate through the input slice
	// TODO: For each element, call the predicate function
	// TODO: If predicate returns true, add the element to the result slice
	// TODO: Return the filtered slice
	return nil
}
