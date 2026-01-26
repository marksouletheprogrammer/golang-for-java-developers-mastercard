# Lab 8 Solution

## How to Run

Install Protocol Buffer compiler (if not already installed):
```bash
# macOS
brew install protobuf

# Or download from: https://github.com/protocolbuffers/protobuf/releases
```

Install Go protobuf plugins:
```bash
cd lab08/solution
make install-tools
```

Add Go bin to PATH (required for protoc to find the plugins):
```bash
# Unix/macOS/WSL
export PATH="$PATH:$(go env GOPATH)/bin"

# Windows Command Prompt
set PATH=%PATH%;%GOPATH%\bin
# Or if GOPATH not set: for /f %i in ('go env GOPATH') do set PATH=%PATH%;%i\bin
```

Generate protobuf code and update dependencies:
```bash
make proto
go mod tidy
```

Run the service:
```bash
go run cmd/server/main.go
```

The service starts two servers:
- HTTP on port 8080
- gRPC on port 9090

## Verifying the Solution

**Note for Windows Command Prompt users:** The curl commands below use `\` for line continuation and single quotes. For Windows Command Prompt, use `^` instead of `\` and replace single quotes with escaped double quotes (`\"` inside the JSON). Or use PowerShell/Git Bash which support Unix syntax.

Test HTTP API:
```bash
# Create an order
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order-1",
    "customer_id": "customer-123",
    "items": [
      {
        "product_id": "prod-1",
        "product_name": "Widget",
        "quantity": 2,
        "unit_price": 29.99
      }
    ]
  }'

# Get order
curl http://localhost:8080/orders/order-1

# List all orders
curl http://localhost:8080/orders

# Update order status
curl -X PATCH http://localhost:8080/orders/order-1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "confirmed"}'

# Delete order
curl -X DELETE http://localhost:8080/orders/order-1
```

Test gRPC API (requires grpcurl):
```bash
# Install grpcurl (macOS)
brew install grpcurl

# Install grpcurl (WSL/Linux)
# Option 1: Using Go
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Option 2: Download binary
# wget https://github.com/fullstorydev/grpcurl/releases/latest/download/grpcurl_*_linux_x86_64.tar.gz
# tar -xvf grpcurl_*_linux_x86_64.tar.gz
# sudo mv grpcurl /usr/local/bin/

# Create order
grpcurl -plaintext -d '{
  "id": "order-2",
  "customer_id": "customer-456",
  "items": [
    {
      "product_id": "prod-2",
      "product_name": "Gadget",
      "quantity": 1,
      "unit_price": 49.99
    }
  ]
}' localhost:9090 orders.OrderService/CreateOrder

# Get order
grpcurl -plaintext -d '{"id": "order-2"}' localhost:9090 orders.OrderService/GetOrder

# List orders
grpcurl -plaintext localhost:9090 orders.OrderService/ListOrders

# Update status
grpcurl -plaintext -d '{
  "id": "order-2",
  "status": "CONFIRMED"
}' localhost:9090 orders.OrderService/UpdateOrderStatus
```

Verify both transports share data:
```bash
# Create order via HTTP
curl -X POST http://localhost:8080/orders -H "Content-Type: application/json" \
  -d '{"id": "shared-order", "customer_id": "cust-1", "items": [{"product_id": "p1", "product_name": "Product", "quantity": 1, "unit_price": 10}]}'

# Retrieve via gRPC
grpcurl -plaintext -d '{"id": "shared-order"}' localhost:9090 orders.OrderService/GetOrder
```

## Key Results

**Clean Architecture Benefits:**
- Business logic in service layer is transport-agnostic
- Same OrderService used by both HTTP and gRPC
- Easy to add new transports (CLI, WebSocket, etc.)
- Each layer has clear responsibility

**Dependency Flow:**
```
Transport Layer (HTTP, gRPC)
         ↓
   Service Layer
         ↓
  Repository Layer
         ↓
    Domain Layer
```

Dependencies flow inward. Outer layers import inner layers, never the reverse.

## Clean Architecture vs Traditional Layering

Traditional layered architecture often tightly couples layers. Clean architecture enforces dependency inversion - high-level policy (business logic) doesn't depend on low-level details (database, transport).

**Key Difference:**
- Traditional: Controller → Service → Repository → Database
- Clean: Domain ← Service ← Repository Interface, with concrete implementations injected

The repository interface lives in the service layer, but implementations live outside. This is the Dependency Inversion Principle - depend on abstractions, not concretions.

**Benefits:**
- Easy to swap implementations (in-memory → PostgreSQL)
- Test with mocks without touching real infrastructure
- Transport changes don't affect business logic
- Business rules are pure and testable
