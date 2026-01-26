# Lab 10: Deployment-Ready Build

You'll add configuration management, feature flags, Docker support, graceful shutdown, and cross-compilation to make your order service ready to deploy anywhere.

### Continuing from Lab 9
This lab continues directly from Lab 9. You can continue to iterate on that lab or start over with the provided starter files in this directory. Note that all labs contain a solution directory (if you are stuck).

**Starter files provided:** `Makefile`, `Dockerfile`, `docker-compose.yml`, `.env.example`, and the complete project structure from Lab 9.

### Part 1: Configuration Management
1. In the `config` package there is a `Config` struct for ports, URLs, timeouts, database settings, log level, and environment name.
2. Implement `LoadConfig()` that reads from environment variables with sensible defaults and validation.
3. Support multiple configuration sources (environment variables, optional config file, defaults).
4. Run `main.go` to load and use configuration. It will print out what configuration it has.

### Part 2: Feature Flags
1. Add feature flags to your config for enabling/disabling gRPC server, profiling endpoints, experimental features, and discount calculations.
2. Implement feature flag checks in your code.
3. Test that features can be toggled without recompiling.

### Part 3: Secrets Management
1. Add support for loading secrets from environment variables and files
2. Never log secrets.
3. Provide a `ValidateSecrets()` function.
4. Create example configuration files with placeholders.

### Part 4: Graceful Shutdown (Optional)
1. Create a shutdown coordinator that listens for OS signals and gracefully closes servers, drains queues, and closes connections.
2. Set a maximum shutdown timeout (e.g., 30 seconds).
3. Test that in-flight requests complete during shutdown.

### Part 5: Build a Production Binary (Optional)
1. Review and complete the provided `Makefile` with targets for building on different platforms.
2. Use build flags to strip debug info and inject version information.
3. Add a version endpoint to your API.

### Part 6: Complete the Multi-Stage Dockerfile
1. Review the provided `Dockerfile` and complete the build stage: Use `golang:1.22`, copy source, download dependencies, run tests, build binary.
2. Complete the runtime stage: Use minimal base image (`alpine` or `distroless`), copy binary, run as non-root user.
3. Optimize for caching by copying `go.mod` and `go.sum` first.
4. Verify image size (should be <20MB for distroless).

### Part 7: Docker Compose for Local Development
1. Review and complete the provided `docker-compose.yml` with your service and any dependencies.
2. Configure environment variables as needed.
3. Start the service with `docker-compose up --build` and test with local client apps.

### Part 8: Health Checks and Readiness (Optional)
1. Create `/health` endpoint that returns service info.
2. Create `/ready` endpoint that checks dependencies and returns 503 if not ready.
3. Add to both HTTP and gRPC.
4. Configure Docker health check.

### Part 9: Logging for Production (Optional)
1. Use structured logging (`log/slog`) with context.
2. Configure JSON format for production, human-readable for development.
3. Add request logging middleware.
4. Never log sensitive data (like full card numbers or order details containing PII).