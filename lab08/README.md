# Lab 8: Clean Architecture & gRPC

You'll build an order management service with clean architecture principles, supporting both REST and gRPC transports that share the same business logic.

## The Architecture

Clean architecture organizes code into concentric layers:
- **Domain Layer**: Core business entities (Order, LineItem)
- **Service Layer**: Business logic and use cases
- **Repository Layer**: Data access abstraction
- **Transport Layer**: HTTP, gRPC - how users interact with your service

Dependencies flow inward. Outer layers depend on inner layers, never the reverse.

**Starter files provided:** Project directory structure including `cmd/server/`, `internal/domain/`, `internal/service/`, `internal/repository/`, `internal/transport/`, `proto/`, plus `Makefile` and `go.mod`.

## Lab Prerequisite Setup

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

When proto files are ready, run:
```bash
make proto
go mod tidy
```

### Part 1: Review the Project Structure
In this lab, the project structure has been created for you with the following layout:
- `cmd/server/` for main.go
- `internal/domain/` for entities (Order, LineItem structs)
- `internal/service/` for business logic
- `internal/repository/` for data access
- `internal/transport/http/` and `internal/transport/grpc/` for transports
- `proto/` for Protocol Buffer definitions

### Part 2: Create the Service Layer
1. Note that the repository layer is provided for you in `internal/repository/`.
1. In `internal/service/orderservice.go`, create an `OrderService` struct.
2. Inject the repository as a dependency (use the interface, not concrete type).
3. Implement business operations: CreateOrder (with validation), GetOrder, ListOrders, UpdateOrderStatus, CalculateOrderTotal.

### Part 3: Refactor HTTP Transport
1. See `OrderHandler` struct in `internal/transport/http/handlers.go`. Note the `OrderService` dependency.
2. Implement the handler methods to delegate to the service.

### Part 4: Define Protocol Buffers
1. Install the Protocol Buffers compiler and Go plugins.
2. In the provided `proto/orders.proto`, we define Order, LineItem messages and OrderService RPCs.
3. Use the provided `Makefile` to generate Go code from the proto file.

```bash
make proto
go mod tidy
```

### Part 5: Implement gRPC Server
1. In `internal/transport/grpc/server.go`, implement the generated gRPC server interface
2. Inject the same OrderService used by HTTP.
3. Convert between protobuf messages and domain entities.
4. Map service errors to gRPC status codes.

### Part 6: Run Both Transports
1. Update `cmd/server/main.go` to create the repository, service, and both transports.
2. Start HTTP server on port 8080 and gRPC server on port 9090.
3. Handle graceful shutdown for both.

### Part 7: Test Both Transports
1. Test HTTP API with curl.
2. Test gRPC with grpcurl or a custom client.
3. Verify both transports see the same data.

### Part 8: Dependency Injection and Configuration (Optional)
1. Create a configuration struct for ports, endpoints, and feature flags.
2. Load configuration from environment variables.
3. Use constructor functions that accept dependencies.