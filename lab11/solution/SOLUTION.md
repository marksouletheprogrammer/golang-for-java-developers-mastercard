# Lab 11 Solution

## How to Run

Start all services with docker-compose:
```bash
cd lab11/solution
docker-compose up --build
```

This starts:
- **Order Service** on http://localhost:8080
- **Jaeger UI** on http://localhost:16686
- **Prometheus** on http://localhost:9090
- **Grafana** on http://localhost:3000

## Testing Observability

**Note for Windows Command Prompt users:** The curl commands below use `\` for line continuation and single quotes. For Windows Command Prompt, use `^` instead of `\` and replace single quotes with escaped double quotes (`\"` inside the JSON). The bash loop syntax won't work - use PowerShell or Git Bash, or manually run the curl command multiple times with different values.

Create an order (triggers all observability features):
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"id": "ORD-001", "customer_id": "CUST-123", "amount": 99.99}'
```

Create multiple orders to generate data:
```bash
# Unix/macOS/WSL/Git Bash
for i in {1..10}; do
  curl -X POST http://localhost:8080/orders \
    -H "Content-Type: application/json" \
    -d "{\"id\": \"ORD-$i\", \"customer_id\": \"CUST-$i\", \"amount\": $((RANDOM % 500 + 10))}"
  sleep 0.5
done

# PowerShell alternative
for ($i=1; $i -le 10; $i++) {
  curl.exe -X POST http://localhost:8080/orders `
    -H "Content-Type: application/json" `
    -d "{`"id`": `"ORD-$i`", `"customer_id`": `"CUST-$i`", `"amount`": $((Get-Random -Minimum 10 -Maximum 500))}"
  Start-Sleep -Milliseconds 500
}
```

Test validation failure (see how errors appear in logs, traces, and metrics):
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"id": "", "customer_id": "CUST-123", "amount": -10}'
```

## Viewing Observability Data

### Logs (Console)

Check service logs:
```bash
docker-compose logs -f order-service
```

Every log entry includes:
- Request ID for correlation
- Timestamp
- Log level (INFO, WARN, ERROR)
- Context (order ID, customer ID, etc.)

### Metrics (Prometheus)

Open Prometheus: http://localhost:9090

Try these queries:
```
# Request rate
rate(http_requests_total[1m])

# Request duration 95th percentile
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Failed orders by reason
orders_failed_total

# Order value distribution
histogram_quantile(0.50, rate(order_value_dollars_bucket[5m]))
```

Or view raw metrics: http://localhost:8080/metrics

### Traces (Jaeger)

Open Jaeger UI: http://localhost:16686

1. Select **order-service** from the Service dropdown
2. Click **Find Traces**
3. Click on a trace to see the full request flow

Each trace shows:
- Complete request timeline
- Span hierarchy (HTTP request → ProcessOrder)
- Span attributes (request ID, order ID, customer ID, amount)
- Error details if any
- Duration for each operation

### Dashboards (Grafana)

Open Grafana: http://localhost:3000
- Login: `admin` / `admin`

Grafana is pre-configured with:
- **Prometheus** datasource (metrics)
- **Jaeger** datasource (traces)
- **Order Service Observability** dashboard (auto-loaded on startup)

The dashboard includes:
- **HTTP Request Rate**: Requests per second by endpoint and status
- **Request Duration**: p50, p95, p99 latency percentiles
- **Total Orders**: Counter of successfully created orders
- **Failed Orders**: Total failures and breakdown by reason (pie chart)
- **Order Value**: Median order value over time
- **HTTP Status Codes**: Stacked view of 2xx, 4xx, 5xx responses

The dashboard auto-refreshes every 5 seconds to show real-time metrics.

## Key Concepts

### The Three Pillars of Observability

**Logs**: Discrete events with timestamps and context. Good for debugging specific issues.
**Metrics**: Aggregated measurements over time. Good for alerting and capacity planning.
**Traces**: Request flow through distributed systems. Good for understanding performance bottlenecks.

Together, these provide complete system visibility.

### Structured Logging with log/slog

Every log entry has machine-parseable structure:
```go
logger.Info("Order created",
    slog.String("request_id", requestID),
    slog.String("order_id", order.ID),
    slog.Float64("amount", order.Amount),
)
```

Benefits: Easy to search/filter, correlate by request ID, no string parsing needed.

### Request ID Correlation

Generate unique ID for each request:
```go
requestID := uuid.New().String()
ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
```

The request ID appears in:
- All log entries for the request
- Response headers (`X-Request-ID`)
- Trace span attributes

This enables correlation across all three pillars.

### Prometheus Metrics

Exposed at `/metrics` endpoint in Prometheus format. Three types:

**Counter**: Monotonically increasing value (requests, orders created)
```go
httpRequests.WithLabelValues(method, path, statusCode).Inc()
```

**Histogram**: Distribution of values with buckets (latency, order value)
```go
requestDuration.WithLabelValues(method, path).Observe(durationSec)
```

**Gauge**: Current value that can go up or down (active connections, queue depth)

Prometheus scrapes `/metrics` every 15 seconds to collect time-series data.

### Distributed Tracing with OpenTelemetry

Traces show request flow. Each operation is a span:
```go
ctx, span := tracer.Start(ctx, "ProcessOrder")
defer span.End()

span.SetAttributes(
    attribute.String("order.id", orderID),
    attribute.Float64("order.amount", amount),
)
```

Spans are exported to Jaeger via OTLP (OpenTelemetry Protocol). Jaeger stores and visualizes the complete trace tree.

### Error Recording in All Three Pillars

When an error occurs:
```go
logger.Error("Order failed", slog.String("error", err.Error()))
metrics.RecordOrderFailed(ctx, "validation_failed")
span.RecordError(err)
span.SetStatus(codes.Error, "Validation failed")
```

- **Log**: Provides error message and context for debugging
- **Metric**: Increments error counter for alerting
- **Trace**: Marks span as error, visible in Jaeger UI

### Middleware Chain

Observability is added via HTTP middleware:
1. **RequestID**: Generate unique ID, add to context and response headers
2. **Logging**: Log request start/completion with request ID
3. **Tracing**: Create root span for request, propagate trace context
4. **Metrics**: Record request counts, durations, status codes

Middleware executes in order, each layer wrapping the next.

### OTLP Exporters

OpenTelemetry Protocol (OTLP) is the standard for exporting telemetry:
- **Traces** → Jaeger via HTTP (port 4318)
- **Metrics** → Prometheus via `/metrics` scraping (port 8080)
- **Logs** → Console (could export to Loki, Elasticsearch, etc.)

This architecture separates instrumentation from backends. You can switch from Jaeger to Zipkin without changing application code.

### Docker Compose Architecture

All services run in a shared `observability` network:
- **order-service**: Application instrumented with OpenTelemetry
- **jaeger**: All-in-one Jaeger for trace storage and UI
- **prometheus**: Scrapes `/metrics` endpoint every 15s
- **grafana**: Visualizes data from Prometheus and Jaeger

No external infrastructure needed - runs locally with `docker-compose up`.

## Cleanup

Stop all services:
```bash
docker-compose down
```

Remove volumes:
```bash
docker-compose down -v
```
