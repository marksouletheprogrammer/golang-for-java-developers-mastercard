package main

import (
	"fmt"
	"time"
)

// Transaction represents a payment transaction in the system.
// Fields are exported (capitalized) to make them accessible from other packages.
// In Go, encapsulation is achieved through package boundaries, not access modifiers.
type Transaction struct {
	TransactionID string
	Amount        float64
	Currency      string
	MerchantID    string
	Timestamp     time.Time
}

// NewTransaction is a constructor function that creates and returns a new Transaction.
// Go doesn't have constructors like Java, but the convention is to use NewX functions.
// Returning a pointer (*Transaction) allows methods to modify the instance efficiently.
func NewTransaction(id string, amount float64, currency string, merchantID string) *Transaction {
	return &Transaction{
		TransactionID: id,
		Amount:        amount,
		Currency:      currency,
		MerchantID:    merchantID,
		Timestamp:     time.Now(),
	}
}

// CalculateFee calculates the transaction fee based on a percentage.
// Uses a pointer receiver (*Transaction) because we might want to modify the transaction
// in the future, and it's more efficient for larger structs (avoids copying).
// Returns the fee amount without modifying the original transaction amount.
func (t *Transaction) CalculateFee(feePercentage float64) float64 {
	return t.Amount * (feePercentage / 100.0)
}

// GetDisplayInfo returns a formatted string with transaction details.
// Uses a pointer receiver for consistency, though a value receiver would work here too.
// In Go, we don't need getter methods - exported fields are directly accessible.
// This method provides a formatted view rather than raw field access.
func (t *Transaction) GetDisplayInfo() string {
	return fmt.Sprintf(
		"Transaction[ID: %s, Amount: %.2f %s, Merchant: %s, Time: %s]",
		t.TransactionID,
		t.Amount,
		t.Currency,
		t.MerchantID,
		t.Timestamp.Format("2006-01-02 15:04:05"),
	)
}

// ProcessPayment processes the payment transaction.
// This method satisfies the Payable interface implicitly.
// Note: There's no explicit "implements Payable" declaration in Go.
// Any type that has a ProcessPayment() error method automatically satisfies Payable.
func (t *Transaction) ProcessPayment() error {
	// Simulate payment processing logic
	if t.Amount <= 0 {
		return fmt.Errorf("invalid amount: %.2f", t.Amount)
	}

	fmt.Printf("Processing payment: %s for %.2f %s\n", t.TransactionID, t.Amount, t.Currency)
	return nil
}

// Payable is an interface that defines the contract for payment processing.
// In Go, interfaces are satisfied implicitly - no "implements" keyword needed.
// Any type with a ProcessPayment() error method automatically implements this interface.
type Payable interface {
	ProcessPayment() error
}

// ProcessPayable accepts any type that implements the Payable interface.
// This demonstrates Go's implicit interface satisfaction.
// We can pass a *Transaction to this function because Transaction has ProcessPayment method.
func ProcessPayable(p Payable) error {
	fmt.Println("Starting payment processing...")
	return p.ProcessPayment()
}

// Filter returns a new slice containing only elements that satisfy the predicate function.
// This is a generic function that works with any type T.
// The [T any] syntax defines a type parameter where 'any' is the constraint (no restriction).
// Similar to Java's <T> but Go uses square brackets and explicit constraints.
// Go's generics do NOT use type erasure - type information is preserved at runtime.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	// Pre-allocate with capacity to potentially avoid reallocations
	result := make([]T, 0, len(slice))

	// Iterate through the input slice
	for _, item := range slice {
		// Apply the predicate function to each element
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}
