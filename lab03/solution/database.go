package main

import (
	"fmt"
	"time"
)

// DatabaseConnection simulates a database connection for demonstrating defer.
type DatabaseConnection struct {
	connected bool
}

// Connect simulates opening a database connection.
func (db *DatabaseConnection) Connect() error {
	fmt.Println("Opening database connection...")
	db.connected = true
	return nil
}

// Close simulates closing a database connection.
// This should always be called when done, even if errors occur.
func (db *DatabaseConnection) Close() error {
	if !db.connected {
		return nil
	}
	fmt.Println("Closing database connection...")
	db.connected = false
	return nil
}

// Query simulates running a query that might fail.
func (db *DatabaseConnection) Query(sql string) error {
	if !db.connected {
		return fmt.Errorf("database not connected")
	}
	
	fmt.Printf("Executing query: %s\n", sql)
	
	// Simulate a query that might fail
	if sql == "SELECT * FROM invalid_table" {
		return fmt.Errorf("table does not exist")
	}
	
	time.Sleep(100 * time.Millisecond) // Simulate work
	fmt.Println("Query executed successfully")
	return nil
}

// ProcessData demonstrates using defer for cleanup.
// The deferred Close() will execute even if Query returns an error.
// defer statements execute in LIFO order when the function returns.
func ProcessData(sql string) error {
	db := &DatabaseConnection{}
	
	// Connect to database
	if err := db.Connect(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	
	// Schedule cleanup with defer - this will run when function exits
	// even if an error occurs below. This ensures resources are always released.
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("Warning: failed to close database: %v\n", err)
		}
	}()
	
	// Run query - if this fails, defer will still execute
	if err := db.Query(sql); err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	
	return nil
}

// DemoMultipleDefers shows how multiple defer statements execute in reverse order.
func DemoMultipleDefers() {
	fmt.Println("\n=== Multiple Defer Demo ===")
	defer fmt.Println("Third - this executes third (first defer registered)")
	defer fmt.Println("Second - this executes second")
	defer fmt.Println("First - this executes first (last defer registered)")
	fmt.Println("Function body")
}
