package observability

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// InitTelemetry initializes OpenTelemetry with OTLP exporters for Jaeger.
// Returns a shutdown function that must be called before program exit.
// Sends traces to Jaeger via OTLP HTTP endpoint.
// TODO: Part 3 - Implement OpenTelemetry initialization
func InitTelemetry(logger *slog.Logger) (func(), error) {
	ctx := context.Background()

	// TODO: Create resource with service information using resource.New()
	// TODO: Add semconv.ServiceName("order-service") and semconv.ServiceVersion("1.0.0")

	// TODO: Get Jaeger endpoint from environment variable OTEL_EXPORTER_OTLP_ENDPOINT
	// TODO: Default to "http://jaeger:4318" if not set
	jaegerEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "http://jaeger:4318"
	}

	// TODO: Setup OTLP trace exporter using otlptracehttp.New()
	// TODO: Configure with endpoint and insecure connection

	// TODO: Setup trace provider with trace.NewTracerProvider()
	// TODO: Use trace.WithBatcher() for the exporter
	// TODO: Use trace.WithResource() for the resource

	// TODO: Set global tracer provider with otel.SetTracerProvider()

	// TODO: Part 7 - Setup propagators for context propagation
	// TODO: Use propagation.NewCompositeTextMapPropagator with TraceContext and Baggage
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info("OpenTelemetry initialized",
		slog.String("jaeger_endpoint", jaegerEndpoint),
		slog.String("service", "order-service"),
	)

	// TODO: Return shutdown function that calls tp.Shutdown(ctx)
	return func() {
		logger.Info("Shutting down OpenTelemetry")
		// TODO: Implement shutdown logic
	}, nil
}
