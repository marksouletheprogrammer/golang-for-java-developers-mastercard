package main

import (
	"fmt"
)

func main() {
	fmt.Println("=== Go Error Handling Demo ===\n")

	// TODO: Part 1 - Basic error handling with if err != nil pattern
	// Create valid and invalid users using NewUser
	// Check for errors and print appropriate messages

	// TODO: Part 2 - Multiple validation errors
	// Create a user with multiple validation issues
	// Call ValidateUserComplete and display all errors

	// TODO: Part 3 - Custom error types with type assertion
	// Create a user that will fail validation
	// Call ValidateUserCustomError
	// Use errors.As() to check if it's a ValidationError
	// Access the Field and Message if it is

	// TODO: Part 4 - Panic and recover
	// Use MustCreateUser in a function with defer/recover
	// Show that program continues after recovery

	// TODO: Part 5 - Defer for cleanup
	// Call ProcessData with valid and invalid queries
	// Observe that cleanup happens even on errors

	// TODO: Part 6 - Multiple defers
	// Call DemoMultipleDefers to see LIFO execution order
}
