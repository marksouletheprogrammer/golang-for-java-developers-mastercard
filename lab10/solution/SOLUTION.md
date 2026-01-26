# Lab 10 Solution

## Prerequisites

This lab builds on Lab 8's order service (clean architecture, HTTP + gRPC). Lab 9's profiling features are not included as Lab 10 focuses on deployment readiness.

## How to Run

### Local Development (no Docker)

Run with default development configuration:
```bash
cd lab10/solution
go run ./cmd/server
```

Run with custom environment variables:
```bash
# Unix/macOS/WSL
ENVIRONMENT=development LOG_LEVEL=debug HTTP_PORT=8080 go run ./cmd/server

# Windows Command Prompt
set ENVIRONMENT=development
set LOG_LEVEL=debug
set HTTP_PORT=8080
go run ./cmd/server
```

Run in production mode (requires API_KEY and JWT_SECRET):
```bash
# Unix/macOS/WSL
ENVIRONMENT=production API_KEY=my-test-key JWT_SECRET=my-test-secret go run ./cmd/server

# Windows Command Prompt
set ENVIRONMENT=production
set API_KEY=my-test-key
set JWT_SECRET=my-test-secret
go run ./cmd/server
```

**Note:** The `API_KEY` and `JWT_SECRET` are configuration placeholders for this lab. They can be any non-empty strings - the service doesn't actually use them. They demonstrate proper secrets management (validation, no logging) and production environment requirements. In a real application, these would be obtained from your infrastructure or secret management system.

### Using Makefile

Build and run:
```bash
make build
./bin/server
```

Build for all platforms:
```bash
make build-all
```

This creates binaries in `bin/`:
- `server-linux-amd64`
- `server-darwin-amd64`
- `server-darwin-arm64`
- `server-windows-amd64.exe`

### Using Docker

Build Docker image:
```bash
make docker-build
```

Run container:
```bash
make docker-run
```

### Using Docker Compose

Start entire stack:
```bash
make docker-compose-up
```

Stop:
```bash
make docker-compose-down
```

## Testing Endpoints

**Note for Windows Command Prompt users:** The curl commands below use `\` for line continuation and single quotes. For Windows Command Prompt, use `^` instead of `\` and replace single quotes with escaped double quotes (`\"` inside the JSON). Or use PowerShell/Git Bash which support Unix syntax.

**Health check:**
```bash
curl http://localhost:8080/health
```

**Readiness check:**
```bash
curl http://localhost:8080/ready
```

**Version info:**
```bash
curl http://localhost:8080/version
```

**Create order:**
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "id": "ORD-001",
    "customer_id": "CUST-001",
    "items": [
      {
        "product_id": "PROD-001",
        "product_name": "Widget",
        "quantity": 2,
        "unit_price": 19.99
      }
    ],
    "status": "pending"
  }'
```

## Key Concepts

### 12-Factor App Principles

**Config in Environment**: All configuration through environment variables, no hardcoded values. Secrets loaded separately from regular config.

**Stateless Processes**: Service doesn't rely on in-memory state. Can be killed and restarted without data loss.

**Port Binding**: Service binds to ports specified by environment, not hardcoded.

**Graceful Shutdown**: Listens for SIGTERM/SIGINT, completes in-flight requests before stopping.

**Dev/Prod Parity**: Same code runs in all environments, only configuration differs.

### Configuration Management

Load from environment with defaults:
```go
port := getEnv("HTTP_PORT", "8080")
```

Validate on startup:
```go
if err := cfg.Validate(); err != nil {
    log.Fatal(err)
}
```

Separate validation for secrets prevents accidental logging. Never log secrets - check them separately and fail fast.

### Feature Flags

Toggle features without recompiling:
```go
if cfg.EnableProfiling {
    // Register profiling endpoints
}
```

Useful for:
- Gradual rollouts
- A/B testing
- Disabling expensive features
- Environment-specific features

### Graceful Shutdown

Listen for OS signals:
```go
shutdown := make(chan os.Signal, 1)
signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
<-shutdown
```

Shutdown with timeout:
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
httpServer.Shutdown(ctx)
```

Ensures in-flight requests complete before process exits. Critical for zero-downtime deployments.

### Build Optimization

Inject version at build time:
```bash
go build -ldflags "-X main.Version=v1.0.0 -s -w"
```

- `-X main.Version=...` sets variable at compile time
- `-s` strips symbol table
- `-w` strips DWARF debugging info
- Result: smaller binary, no debug symbols

### Multi-Stage Docker

**Build stage**: Full Go toolchain, compile binary
**Runtime stage**: Minimal Alpine image, just binary and CA certs

Benefits:
- Small image size (~15-20MB vs 1GB+)
- Reduced attack surface
- Faster deploys
- Layer caching optimizes rebuilds

Copy go.mod/go.sum first so dependency layer is cached unless dependencies change.

### Health vs Readiness

**Health** (`/health`): Is the process alive? Always returns 200 if running.

**Readiness** (`/ready`): Can the service handle requests? Checks dependencies (database, cache, etc.). Returns 503 if not ready.

Kubernetes uses both:
- Liveness probe uses `/health` - restarts if failing
- Readiness probe uses `/ready` - removes from load balancer if failing

### Structured Logging

Use `log/slog` for structured logging:
```go
logger.Info("Request completed",
    slog.String("method", r.Method),
    slog.Int("status", statusCode),
    slog.Duration("duration", duration),
)
```

**Text format** (development): Human-readable
**JSON format** (production): Machine-parseable for log aggregation

Never log sensitive data (passwords, API keys, full card numbers, PII). Log request IDs for tracing instead.

## Lab 10 Additions Summary

This lab adds deployment-ready features on top of Lab 8's order service:

**New Files:**
- `config/config.go` - Configuration management with environment variables, validation, feature flags
- `.env.example` - Example configuration file
- `Dockerfile` - Multi-stage build with non-root user and health check
- `docker-compose.yml` - Local development orchestration
- `Makefile` - Enhanced with version injection, cross-compilation, Docker commands

**Modified Files:**
- `cmd/server/main.go` - Uses config package, adds health/ready/version endpoints, feature-flagged gRPC, improved graceful shutdown

**Features Added:**
- Environment-based configuration with validation
- Feature flags (EnableGRPC, EnableMetrics, EnableHealthz, EnableDebugMode)
- Secrets management with production validation
- Health check endpoints (`/health`, `/ready`, `/version`)
- Graceful shutdown with configurable timeout
- Version information injected at build time
- Multi-stage Docker build (Alpine-based, ~15-20MB)
- Cross-platform compilation support
- Docker Compose for local development

**Lab 8 Features Preserved:**
- Clean architecture (domain, repository, service, transport layers)
- HTTP REST API with full CRUD operations
- gRPC API (now feature-flagged)
- In-memory repository
- All order management functionality

**12-Factor App Compliance:**
- ✅ Config in environment variables
- ✅ Stateless processes
- ✅ Port binding via config
- ✅ Graceful shutdown
- ✅ Dev/prod parity through configuration
