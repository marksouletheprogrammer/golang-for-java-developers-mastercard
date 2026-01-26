# Lab 4 Solution

## How to Run

Start the server:
```bash
cd lab04/solution
go run .
```

The server starts on http://localhost:8080

## Testing the API

**Note for Windows Command Prompt users:** The curl commands below use `\` for line continuation and single quotes. For Windows Command Prompt, use `^` instead of `\` and replace single quotes with escaped double quotes (`\"` inside the JSON). Or use PowerShell/Git Bash which support Unix syntax.

List all merchants:
```bash
curl http://localhost:8080/merchants
```

Get a specific merchant:
```bash
curl http://localhost:8080/merchants/MERCH-001
```

Create a new merchant:
```bash
curl -X POST http://localhost:8080/merchants \
  -H "Content-Type: application/json" \
  -d '{
    "id": "MERCH-004",
    "name": "New Merchant",
    "category": "Services",
    "country": "USA",
    "status": "active"
  }'
```

Test error cases:
```bash
# Not found (404)
curl http://localhost:8080/merchants/INVALID-ID

# Duplicate ID (409)
curl -X POST http://localhost:8080/merchants \
  -H "Content-Type: application/json" \
  -d '{"id": "MERCH-001", "name": "Duplicate", "category": "Test", "country": "USA", "status": "active"}'

# Validation error (400)
curl -X POST http://localhost:8080/merchants \
  -H "Content-Type: application/json" \
  -d '{"id": "", "name": "Invalid"}'
```

## Key Concepts

### Thread-Safe Storage

Go's `sync.RWMutex` allows multiple concurrent readers or one exclusive writer. Using `RLock/RUnlock` for reads and `Lock/Unlock` for writes prevents race conditions without blocking all concurrent operations.

### HTTP Routing

Go's `http.ServeMux` matches patterns:
- `/merchants` matches exactly
- `/merchants/` matches as prefix (handles `/merchants/MERCH-001`)

Extract path parameters manually using `strings.TrimPrefix`. For complex routing, consider using a router library like `gorilla/mux` or `chi`.

### Middleware Pattern

The `ServeHTTP` method wraps the mux to add logging. This middleware pattern logs every request's method, path, status code, and duration. You can chain multiple middleware functions for authentication, rate limiting, etc.

### JSON Encoding

`json.Encoder` writes directly to `http.ResponseWriter` for efficiency. Use struct tags like `json:"id"` to control JSON field names. The encoder automatically handles marshaling Go types to JSON.

### Error Responses

Consistent error format across all endpoints makes API easier to consume. Always set `Content-Type: application/json` and appropriate HTTP status codes (400 for validation, 404 for not found, 409 for conflict, etc.).
