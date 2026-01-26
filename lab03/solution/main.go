package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("=== Go Error Handling Demo ===\n")
	
	// Part 1: Basic error handling with if err != nil pattern
	fmt.Println("--- Part 1: Basic Error Handling ---")
	validUser, err := NewUser("Alice", "alice@example.com", 25, "USA")
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", validUser)
	}
	
	// Try creating invalid users
	invalidUser1, err := NewUser("", "bob@example.com", 30, "UK")
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", invalidUser1)
	}
	
	invalidUser2, err := NewUser("Charlie", "invalid-email", 20, "Canada")
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", invalidUser2)
	}
	
	invalidUser3, err := NewUser("David", "david@example.com", 16, "Australia")
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", invalidUser3)
	}
	
	// Part 2: Multiple validation errors
	fmt.Println("\n--- Part 2: Multiple Validation Errors ---")
	badUser := &User{Name: "", Email: "bad", Age: 15, Country: "USA"}
	err = ValidateUserComplete(badUser)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}
	
	// Part 3: Custom error types with type assertion
	fmt.Println("\n--- Part 3: Custom Error Types ---")
	testUser := &User{Name: "Eve", Email: "invalid", Age: 20, Country: "France"}
	err = ValidateUserCustomError(testUser)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
		
		// Type assertion to check if it's a ValidationError
		// This allows us to access structured error information
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf("Field with error: %s\n", validationErr.Field)
			fmt.Printf("Error message: %s\n", validationErr.Message)
		}
	}
	
	// Part 4: Panic and recover
	fmt.Println("\n--- Part 4: Panic and Recover ---")
	
	// Use MustCreateUser which panics on error
	// Wrap it in a function so we can recover from the panic
	func() {
		// defer with recover catches panics
		// recover() returns nil if no panic occurred
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
		}()
		
		// This will panic because age is invalid
		_ = MustCreateUser("Frank", "frank@example.com", 10, "Germany")
		fmt.Println("This line won't execute because of panic")
	}()
	
	fmt.Println("Program continues after recovery")
	
	// Part 5: Defer for cleanup
	fmt.Println("\n--- Part 5: Defer for Cleanup ---")
	
	// Successful query
	err = ProcessData("SELECT * FROM users")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println()
	
	// Failed query - defer still ensures cleanup happens
	err = ProcessData("SELECT * FROM invalid_table")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	// Part 6: Multiple defers execute in reverse order
	DemoMultipleDefers()
}
