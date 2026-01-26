package main

// User represents a user in the system.
type User struct {
	Name    string
	Email   string
	Age     int
	Country string
}

// ValidateUser validates a user and returns an error if validation fails.
// This is the Go way - return errors instead of throwing exceptions.
// The caller must explicitly check for errors using if err != nil.
func ValidateUser(u *User) error {
	// TODO: Part 1 - Implement validation logic
	// Check if name is empty
	// Check if email contains @
	// Check if age is at least 18
	return nil
}

// NewUser is a constructor that validates the user before creating it.
// Returns both *User and error, following Go's convention of returning errors.
// The caller must check if err != nil before using the returned User.
func NewUser(name, email string, age int, country string) (*User, error) {
	// TODO: Part 2 - Implement user constructor
	// Create User struct
	// Call ValidateUser
	// Return nil, err if validation fails
	// Return user, nil if validation succeeds
	return nil, nil
}

// ValidateUserComplete checks all validation rules and collects all errors.
// Instead of returning on first error, it accumulates all failures.
// This is useful for form validation where you want to show all errors at once.
func ValidateUserComplete(u *User) error {
	// TODO: Part 3 - Implement complete validation
	// Create a slice to collect error messages
	// Check all validation rules and append errors
	// Return combined error with all messages
	return nil
}

// ValidationError is a custom error type that provides structured information
// about validation failures. It implements the error interface by having an Error() method.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface for ValidationError.
// Any type with an Error() string method automatically satisfies the error interface.
func (ve *ValidationError) Error() string {
	// TODO: Part 4 - Implement Error() method
	// Return formatted string with Field and Message
	return ""
}

// ValidateUserCustomError validates user and returns a custom ValidationError.
// This allows callers to type-assert the error to get structured information.
func ValidateUserCustomError(u *User) error {
	// TODO: Part 4 - Implement validation with custom error type
	// Check each field and return &ValidationError{Field: "...", Message: "..."}
	return nil
}

// MustCreateUser creates a user and panics if validation fails.
// The "Must" prefix is a Go convention indicating the function will panic on error.
// Use this only in initialization code where failure should stop the program.
func MustCreateUser(name, email string, age int, country string) *User {
	// TODO: Part 5 - Implement MustCreateUser
	// Call NewUser
	// If error occurs, panic with descriptive message
	// Otherwise return the user
	return nil
}
