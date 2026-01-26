package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Payment Transaction Type System Demo ===\n")

	// Create a new transaction using the constructor function
	transaction := NewTransaction("TXN-001", 150.00, "USD", "MERCH-123")

	// Display transaction information
	fmt.Println(transaction.GetDisplayInfo())
	fmt.Println()

	// Calculate and display fee
	feePercentage := 2.5
	fee := transaction.CalculateFee(feePercentage)
	fmt.Printf("Transaction fee (%.2f%%): %.2f %s\n", feePercentage, fee, transaction.Currency)
	fmt.Println()

	// Demonstrate interface satisfaction
	// We can pass *Transaction to ProcessPayable because Transaction satisfies Payable
	// No explicit "implements" declaration was needed
	fmt.Println("=== Demonstrating Interface Satisfaction ===")
	err := ProcessPayable(transaction)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
	} else {
		fmt.Println("Payment completed successfully")
	}
	fmt.Println()

	// Create another transaction to demonstrate error handling
	fmt.Println("=== Testing Error Handling ===")
	invalidTransaction := NewTransaction("TXN-002", -50.00, "USD", "MERCH-456")
	err = ProcessPayable(invalidTransaction)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
	}
	fmt.Println()

	// Part 6: Demonstrate generic Filter function
	fmt.Println("=== Part 6: Generic Filter Function Demo ===")

	// Create a slice of transactions with varying amounts and currencies
	transactions := []*Transaction{
		NewTransaction("TXN-001", 150.00, "USD", "MERCH-123"),
		NewTransaction("TXN-002", 50.00, "EUR", "MERCH-456"),
		NewTransaction("TXN-003", 250.00, "USD", "MERCH-123"),
		NewTransaction("TXN-004", 75.00, "GBP", "MERCH-789"),
		NewTransaction("TXN-005", 300.00, "USD", "MERCH-456"),
	}

	fmt.Println("All transactions:")
	for _, t := range transactions {
		fmt.Printf("  %s: %.2f %s\n", t.TransactionID, t.Amount, t.Currency)
	}
	fmt.Println()

	// Filter transactions over $100 - type is inferred as []*Transaction
	highValue := Filter(transactions, func(t *Transaction) bool {
		return t.Amount > 100
	})
	fmt.Println("Transactions over $100:")
	for _, t := range highValue {
		fmt.Printf("  %s: %.2f %s\n", t.TransactionID, t.Amount, t.Currency)
	}
	fmt.Println()

	// Filter USD transactions only
	usdOnly := Filter(transactions, func(t *Transaction) bool {
		return t.Currency == "USD"
	})
	fmt.Println("USD transactions only:")
	for _, t := range usdOnly {
		fmt.Printf("  %s: %.2f %s\n", t.TransactionID, t.Amount, t.Currency)
	}
	fmt.Println()

	// Filter by specific merchant
	merchantTransactions := Filter(transactions, func(t *Transaction) bool {
		return t.MerchantID == "MERCH-123"
	})
	fmt.Println("Transactions for MERCH-123:")
	for _, t := range merchantTransactions {
		fmt.Printf("  %s: %.2f %s\n", t.TransactionID, t.Amount, t.Currency)
	}
	fmt.Println()

	fmt.Println("Note: The Filter function works with any type!")
	fmt.Println("Type parameter T is inferred from the slice type - no need to specify Filter[*Transaction]")
}
