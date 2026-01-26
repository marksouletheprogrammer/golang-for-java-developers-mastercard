package main

import (
	"errors"
	"fmt"
	"strings"
)

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
	// Check if name is empty
	if u.Name == "" {
		return errors.New("name must not be empty")
	}
	
	// Check if email contains @
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		return errors.New("valid email is required")
	}
	
	// Check if age is at least 18
	if u.Age < 18 {
		return fmt.Errorf("user must be 18 or older, got %d", u.Age)
	}
	
	return nil
}

// NewUser is a constructor that validates the user before creating it.
// Returns both *User and error, following Go's convention of returning errors.
// The caller must check if err != nil before using the returned User.
func NewUser(name, email string, age int, country string) (*User, error) {
	u := &User{
		Name:    name,
		Email:   email,
		Age:     age,
		Country: country,
	}
	
	// Validate the user
	if err := ValidateUser(u); err != nil {
		return nil, err
	}
	
	return u, nil
}

// ValidateUserComplete checks all validation rules and collects all errors.
// Instead of returning on first error, it accumulates all failures.
// This is useful for form validation where you want to show all errors at once.
func ValidateUserComplete(u *User) error {
	var errorMessages []string
	
	if u.Name == "" {
		errorMessages = append(errorMessages, "name must not be empty")
	}
	
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		errorMessages = append(errorMessages, "valid email is required")
	}
	
	if u.Age < 18 {
		errorMessages = append(errorMessages, fmt.Sprintf("user must be 18 or older, got %d", u.Age))
	}
	
	if len(errorMessages) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}
	
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
	return fmt.Sprintf("validation failed on field '%s': %s", ve.Field, ve.Message)
}

// ValidateUserCustomError validates user and returns a custom ValidationError.
// This allows callers to type-assert the error to get structured information.
func ValidateUserCustomError(u *User) error {
	if u.Name == "" {
		return &ValidationError{Field: "Name", Message: "must not be empty"}
	}
	
	if u.Email == "" || !strings.Contains(u.Email, "@") {
		return &ValidationError{Field: "Email", Message: "valid email is required"}
	}
	
	if u.Age < 18 {
		return &ValidationError{Field: "Age", Message: fmt.Sprintf("must be 18 or older, got %d", u.Age)}
	}
	
	return nil
}

// MustCreateUser creates a user and panics if validation fails.
// The "Must" prefix is a Go convention indicating the function will panic on error.
// Use this only in initialization code where failure should stop the program.
func MustCreateUser(name, email string, age int, country string) *User {
	u, err := NewUser(name, email, age, country)
	if err != nil {
		panic(fmt.Sprintf("failed to create user: %v", err))
	}
	return u
}
