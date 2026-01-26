# Lab 2: Type System & Idiomatic Go

In this lab, you'll translate a Java class into idiomatic Go code. This exercise will help you understand the fundamental differences between Java's object-oriented approach and Go's composition-based design.

## The Java Class

See the sample Java class `transaction.java` in this directory. The class has:
- Private fields: transactionId (String), amount (BigDecimal), currency (String), merchantId (String), timestamp (Date)
- Public getters and setters for all fields
- A method `calculateFee(double feePercentage)` that returns the transaction fee.
- A method `getDisplayInfo()` that returns a formatted string.
- The class implements an interface `Payable` with a method signature for processing payments.

**Starter files provided:** `transaction.go`, `main.go`, `go.mod`, and `transaction.java` (for reference).


### Part 1: Create the Transaction Struct
1. In this lab, use the provided starter files.
2. In the provided `transaction.go` file, define a struct named `Transaction` with appropriate fields.
3. Consider which fields should be exported (public) vs. unexported (private).

**Example of exported vs unexported fields:**
```go
type Transaction struct {
    TransactionID string    // Exported (public) - capitalized first letter
    //...
    merchantID    string    // unexported (private) - lowercase first letter
    // ...
}
```
Go uses capitalization to control visibility: capitalized names are exported (public), lowercase names are unexported (private).

### Part 2: Implement a Constructor Pattern
1. Create a constructor function for Transaction following Go naming conventions.
2. The function should accept all necessary parameters and return a pointer to a new Transaction instance.

**Example constructor pattern:**
```go
// Constructor function - conventionally named NewTypeName
func NewTransaction(id string, /* other params */) *Transaction {
    return &Transaction{
        TransactionID: id,
        //...
    }
}
```
Note: Go doesn't have constructors like Java. Instead, use functions that return pointers to new instances.

### Part 3: Add Methods to Transaction
1. Implement a method that calculates the transaction fee (without modifying the original amount).
2. Implement a method that returns formatted display information.
3. Decide whether to use pointer receivers or value receivers.

**Example of pointer vs value receivers:**
```go
// Pointer receiver - use when method needs to modify the struct.
func (t *Transaction) UpdateAmount(newAmount float64) {
    t.Amount = newAmount // Can modify the original
}

// Value receiver - use for methods that don't modify.
func (t Transaction) CalculateFee(feePercentage float64) float64 {
    return t.Amount * feePercentage // Cannot modify original Transaction
}

func (t Transaction) GetDisplayInfo() string {
    return fmt.Sprintf("Transaction %s: %.2f %s", t.TransactionID, t.Amount, t.Currency)
}
```
Rule of thumb: Use pointer receivers if the method modifies the receiver or the struct is large. Use value receivers for read-only operations on small structs.

### Part 4: Define a Payable Interface
1. Create an interface called `Payable` that declares a method for processing payments.
2. Ensure your Transaction struct satisfies this interface implicitly.

### Part 5: Demonstrate Interface Satisfaction
1. In the provided `main.go`, create a function that accepts a `Payable` parameter.
2. Create a Transaction instance and pass it to this function.
3. Observe that no explicit "implements" declaration is needed.

### Part 6: Introduction to Generics (Optional)
1. Implement a generic `Filter` function that works with any slice type:
   - Function signature: `Filter[T any](slice []T, predicate func(T) bool) []T`
   - The function should return a new slice containing only elements that satisfy the predicate
   - The `[T any]` syntax defines a type parameter (similar to Java's `<T>`)
2. Create predicate functions to filter transactions:
   - Filter transactions with amount greater than a threshold
   - Filter transactions by specific currency
   - Filter transactions by merchant ID pattern
3. Test your generic Filter function with transaction slices and observe type inference in action.

**Comparison to Java:**
- Go: `Filter[T any](slice []T, predicate func(T) bool) []T`
- Java: `<T> List<T> filter(List<T> list, Predicate<T> predicate)`
- Go uses square brackets `[]` for type parameters instead of angle brackets `<>`
- Go requires explicit constraint (`any` means no constraint, similar to Java's unbounded type parameter)
- Go's generics do NOT use type erasure (unlike Java), so type information is preserved at runtime

**Example usage:**
```go
transactions := []*Transaction{
    NewTransaction("TXN-001", 150.00, "USD", "MERCH-123"),
    NewTransaction("TXN-002", 50.00, "EUR", "MERCH-456"),
    NewTransaction("TXN-003", 250.00, "USD", "MERCH-123"),
}

// Filter transactions over $100 - type inferred as []*Transaction
highValue := Filter(transactions, func(t *Transaction) bool {
    return t.Amount > 100
})

// Filter USD transactions
usdOnly := Filter(transactions, func(t *Transaction) bool {
    return t.Currency == "USD"
})
```

### Part 7: Compare with Java (Optional)
Document the key differences:
1. How does Go's implicit interface satisfaction differ from Java's explicit `implements`?
2. How does Go handle encapsulation without access modifiers?
3. What replaces Java's getters and setters in idiomatic Go?
4. Why doesn't Go have classes?