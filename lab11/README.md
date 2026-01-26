# Lab 11: Observability

You'll instrument your order service with comprehensive observability using OpenTelemetry. You'll emit structured logs, collect metrics, and create distributed traces.

## The Three Pillars

**Logs**: Discrete events ("Order created", "Payment failed").

**Metrics**: Aggregated measurements (orders/sec, error rate, latency).

**Traces**: Request flow through distributed systems.

Together, these give you complete visibility into your service's behavior.

## Important Note
This lab includes a Dockerfile and docker-compose file. These already have everything that you need for this lab. In this lab, you will add instrumentation to the Golang application, then use Docker to see it in action.

**Starter files provided:** `middleware/` directory with `logging.go`, `metrics.go`, `request_id.go`, `tracing.go`, plus `observability/` directory with `metrics.go` and `telemetry.go`. Docker configuration files (`Dockerfile`, `docker-compose.yml`, `prometheus.yml`) are also provided.

### Part 1: Structured Logging Setup
1. Replace `fmt.Println` and basic `log` calls with `log/slog`.
2. Configure log level and format (JSON for production, human-readable for development).
3. Add context to the logs: request ID, user ID, order ID, operation name.

### Part 2: Request Logging Middleware
1. In the provided `middleware/` directory, complete the HTTP middleware and gRPC interceptor implementations that generate request IDs and log request start/completion.
2. Pass request ID through context.
3. Verify request ID appears in all logs for a single request.

### Part 3: OpenTelemetry Setup
1. Add OpenTelemetry SDK and exporters.
2. Initialize tracer and meter providers in `main.go`.
3. Create helper functions for getting tracers/meters and starting spans.

### Part 4: Metrics Collection
1. Implement at least 1 metric of one of these types:
   - Counter: HTTP requests, orders created, orders failed, status changes.
   - Histograms: Request duration, order processing time, order value.
   - Gauges: Pending orders, active connections.
2. Increment metrics in appropriate places with labels/attributes.
3. Add labels for HTTP method, status code, endpoint name, order status, etc.

### Part 5: Expose Metrics Endpoint
1. Add a `/metrics` HTTP endpoint.
2. Use the included docker-compose file to run `docker-compose up --build`, which will start this service and all dependencies.
3. Prometheus should start scraping from the `/metrics` endpoint. Navigate to the prometheus home page in your browser to query metrics.

### Part 6: Distributed Tracing
1. Create at least 1 span for HTTP handlers, gRPC handlers, service methods, repository operations, and/or external calls.
2. Set span names, attributes, kinds, and record errors.
3. Verify spans are created and linked at runtime.

### Part 7: Context Propagation (Optional)
1. Propagate trace context in HTTP headers and gRPC metadata.
2. Pass context through all service layers.
3. Verify a single request creates a complete trace tree.
4. (Optional) Create spans for handlers or methods with names, attributes, etc. Verify these spans appear in the traces after restart.

### Part 8: Export to Observability Backend  (Optional)
1. Start up local Jaeger, Prometheus, and Grafana via docker-compose.
2. Test the API. Trigger validation errors and see how they appear in the logs, trace, and metrics.
