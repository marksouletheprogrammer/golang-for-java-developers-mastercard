# Lab 2 Solution

## How to Run

```bash
cd lab02/solution
go run .
```

**Why `go run .` instead of `go run main.go`?**

This lab has multiple `.go` files (`main.go` and `transaction.go`) in the same package. The `go run .` command compiles all Go files in the current directory together. If you tried `go run main.go` alone, it would fail because `main.go` references types and functions defined in `transaction.go`.

For single-file programs, `go run main.go` works fine, but `go run .` works for both single and multi-file packages, making it more consistent.

## Go vs Java: Key Differences

### Implicit Interface Satisfaction

Java requires explicit declaration:
```java
public class Transaction implements Payable {
    // Must explicitly declare we implement Payable
}
```

Go uses implicit satisfaction:
```go
type Transaction struct { }
func (t *Transaction) ProcessPayment() error { }
// Automatically satisfies Payable interface - no declaration needed
```

Any type with the right method signature satisfies the interface. This enables duck typing while maintaining type safety. You can define interfaces after the types exist, which is impossible in Java.

### Encapsulation Without Access Modifiers

Java uses `private`, `protected`, `public`:
```java
private String transactionId;
public String getTransactionId() { return transactionId; }
```

Go uses capitalization and package boundaries:
```go
type Transaction struct {
    TransactionID string  // Exported (public)
    privateField  string  // Unexported (private to package)
}
```

Exported (capitalized) names are public across packages. Unexported (lowercase) names are private to the package. No getter/setter boilerplate needed - access exported fields directly.

### No Getters/Setters in Idiomatic Go

Java requires getters and setters for encapsulation:
```java
public BigDecimal getAmount() { return amount; }
public void setAmount(BigDecimal amount) { this.amount = amount; }
```

Go accesses exported fields directly:
```go
transaction.Amount = 200.00  // Direct access
fmt.Println(transaction.Amount)
```

Only create methods when you need computed values or validation, not for simple field access.

### Why No Classes?

Go chose composition over inheritance. Instead of class hierarchies, Go uses:
- **Structs** for data
- **Interfaces** for behavior contracts  
- **Embedding** for composition

This eliminates complexity from inheritance (diamond problem, fragile base class, etc.) and makes code easier to reason about. Behavior is defined by what a type can do (interfaces), not what it inherits from.

### Generics in Go (Part 6)

Go introduced generics in version 1.18. The syntax differs from Java but the concepts are familiar.

**Java generics:**
```java
public <T> List<T> filter(List<T> list, Predicate<T> predicate) {
    List<T> result = new ArrayList<>();
    for (T item : list) {
        if (predicate.test(item)) {
            result.add(item);
        }
    }
    return result;
}
```

**Go generics:**
```go
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0, len(slice))
    for _, item := range slice {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}
```

**Key differences:**
- **Syntax**: Go uses `[T any]` instead of `<T>`
- **Constraints**: `any` means no constraint (like Java's unbounded type parameter). Go also supports `comparable` and custom constraints
- **Type Erasure**: Java uses type erasure at runtime - generic type info is lost. Go preserves type information, so you can use type switches and reflection with generics
- **Type Inference**: Both support type inference, but Go's is more aggressive. You rarely need to specify type parameters explicitly

**Usage example:**
```go
// Type is inferred from the slice - no need to write Filter[*Transaction]
highValue := Filter(transactions, func(t *Transaction) bool {
    return t.Amount > 100
})
```

**When to use generics vs interfaces:**
- Use **generics** when you need to preserve the specific type (like `Filter` returning the same type).
- Use **interfaces** when you need polymorphic behavior (like `Payable` working with different payment types).
- Generics are compile-time polymorphism, interfaces are runtime polymorphism.
