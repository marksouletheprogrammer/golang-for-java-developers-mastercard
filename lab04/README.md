# Lab 4: Build HTTP Services

You'll build a RESTful API for managing merchant information. This service will allow clients to create, retrieve, and list merchants using JSON over HTTP.

**Starter files provided:** `server.go`, `merchant.go`, `main.go`, and `go.mod`.

### Part 1: Set Up the HTTP Server
1. In this lab, use the provided starter files.
2. In the provided `server.go`, implement a `StartServer` function that listens on port 8080.
3. Update the provided `main.go` to call it.

### Part 2: In-Memory Merchant Storage
1. In `merchant.go`, create a `Merchant` struct with fields: ID, name, category, country, status.
2. Create a thread-safe in-memory data structure to store merchants (consider using a map with a mutex).
3. Initialize it with 2-3 sample merchants.

### Part 3: Create a GET /merchants Endpoint
1. Create a handler that retrieves all merchants and returns them as JSON.
2. Set appropriate headers and status codes.
3. Test with curl or a browser.

### Part 4: Create a GET /merchants/{id} Endpoint
1. Create a handler that retrieves a single merchant by ID.
2. Return 200 with the merchant if found, 404 if not found.
3. Try it with both existing and non-existing IDs.

### Part 5: Create a POST /merchants Endpoint (Optional)
1. Create a handler that accepts JSON, validates it, and adds a new merchant.
2. Return appropriate status codes: 201 for success, 400 for validation errors, 409 for duplicate IDs.
3. Set the Location header to point to the new resource.
 
### Part 6: Error Handling in HTTP Handlers (Optional)
1. Create a helper function for writing consistent JSON error responses.
2. Update all handlers to use it.

### Part 7: Request Logging Middleware (Optional)
Create middleware that logs HTTP method, path, status code, and duration for each request.

