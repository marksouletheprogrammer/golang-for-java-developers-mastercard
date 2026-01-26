# Lab 3 Solution

## How to Run

```bash
cd lab03/solution
go run .
```

## Go vs Java Error Handling

### Explicit Error Returns vs Exceptions

Java uses exceptions:
```java
public User createUser(String name, String email) throws ValidationException {
    if (name.isEmpty()) {
        throw new ValidationException("Name required");
    }
    return new User(name, email);
}

try {
    User user = createUser("", "test@example.com");
} catch (ValidationException e) {
    System.err.println(e.getMessage());
}
```

Go returns errors as values:
```go
func createUser(name, email string) (*User, error) {
    if name == "" {
        return nil, errors.New("name required")
    }
    return &User{name, email}, nil
}

user, err := createUser("", "test@example.com")
if err != nil {
    fmt.Println(err)
}
```

In Go, errors are explicit return values, not exceptions. You cannot ignore errors - the compiler forces you to handle or explicitly discard them. This makes error handling visible in the code flow.

### No Try-Catch Blocks

Go doesn't have try-catch. Every function that can fail returns an error as its last return value. You check errors with `if err != nil` immediately after the function call. This makes error handling linear and easy to follow.

Exception-based error handling in Java can skip multiple stack frames, making control flow hard to trace. Go's approach keeps error handling local and explicit.

### Panic and Recover (Rarely Used)

Go has `panic` and `recover`, but they're not for normal error handling. Use them only for unrecoverable situations like programmer errors or initialization failures.

Panic stops normal execution and begins unwinding the stack, running deferred functions. Recover can catch a panic in a deferred function. Most Go code never uses panic/recover - errors as values are preferred.

Functions prefixed with `Must` (like `MustCreateUser`) are a convention indicating they panic on failure. Use these only in init code where failure should stop the program.

### Defer for Cleanup

Java uses try-finally or try-with-resources:
```java
Connection conn = db.connect();
try {
    conn.query("SELECT * FROM users");
} finally {
    conn.close();
}
```

Go uses defer:
```go
conn := db.Connect()
defer conn.Close()
conn.Query("SELECT * FROM users")
```

Defer schedules a function to run when the surrounding function returns, regardless of how it returns (normal return, panic, or error). Deferred functions execute in LIFO order. This ensures cleanup always happens without cluttering error paths with duplicate cleanup code.
