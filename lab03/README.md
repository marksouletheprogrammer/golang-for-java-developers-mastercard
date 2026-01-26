# Lab 3: Error Handling the Go Way

Error handling is one of the most distinctive features of Go compared to Java. In this lab, you'll learn to handle errors explicitly using return values rather than exceptions, and you'll understand when (rarely) to use panic and recover.

## Background: The Java Version

In Java, you might have a method like:
```java
public void validateUser(User user) throws InvalidUserException {
    if (user.getAge() < 18) {
        throw new InvalidUserException("User must be 18 or older");
    }
    if (user.getEmail() == null || !user.getEmail().contains("@")) {
        throw new InvalidUserException("Valid email is required");
    }
    // more validations...
}
```

In Go, we don't throw exceptions. We return errors.

**Starter files provided:** `user.go`, `main.go`, `database.go`, and `go.mod`.

### Part 1: Basic Error Handling
1. In this lab, use the provided starter files.
2. In `user.go`, define a `User` struct with fields for name, email, age, and country.
3. Create a function called `ValidateUser` that accepts a pointer to a User and returns an error.
4. Implement validation rules: age must be 18+, email must contain "@", name must not be empty.
5. Return descriptive error messages using `errors.New()` or `fmt.Errorf()`.

### Part 2: Use the Validation Function
1. Create a user constructor that calls `ValidateUser` and returns both `*User` and `error`.
2. In the provided `main.go`, create several User instances with both valid and invalid data.
3. Check for errors using the `if err != nil` pattern and print appropriate messages.

### Part 3: Multiple Validation Errors
1. Create a `ValidateUserComplete` function that checks ALL validation rules and collects all errors.
2. Return a single error containing information about all failures.

### Part 4: Custom Error Types (Optional)
1. Create a custom error type called `ValidationError` that implements the `error` interface.
2. Store information about which field failed and why.
3. Update your validation function to return this custom error type.
4. Demonstrate type assertion to check if an error is a `ValidationError`

**Example custom error type:**
```go
// Custom error type
type ValidationError struct {
    //...
}

// Implement the error interface
func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Return custom error
func validateUser(user *User) error {
    if user.Age < 18 {
        return &ValidationError{Field: "Age", Message: "must be 18 or older"}
    }
    return nil
}

// Type assertion to check error type
err := validateUser(user)
if err != nil {
    if valErr, ok := err.(*ValidationError); ok {
        fmt.Printf("Validation failed - Field: %s, Message: %s\n", valErr.Field, valErr.Message)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

### Part 5: Understanding Panic and Recover (Optional)
1. Create a `MustCreateUser` function that panics on validation failure.
2. Experiment with `recover()` to catch the panic in a deferred function.

### Part 6: Defer for Cleanup (Optional)
1. The provided `database.go` file contains helper functions for database operations. Review and use these, or create additional functions that simulate opening and closing a database connection.
2. Use `defer` to ensure cleanup happens even when errors occur.