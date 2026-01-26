package observability

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTelemetry initializes OpenTelemetry with OTLP exporters for Jaeger.
// Returns a shutdown function that must be called before program exit.
// Sends traces to Jaeger via OTLP HTTP endpoint.
func InitTelemetry(logger *slog.Logger) (func(), error) {
	ctx := context.Background()
	
	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("order-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}
	
	// Get Jaeger endpoint from environment or use default
	jaegerEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "http://jaeger:4318"
	}
	
	// Setup OTLP trace exporter for Jaeger
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(jaegerEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	
	// Setup trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	
	// Setup propagators for context propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	
	logger.Info("OpenTelemetry initialized",
		slog.String("jaeger_endpoint", jaegerEndpoint),
		slog.String("service", "order-service"),
	)
	
	// Return shutdown function
	return func() {
		logger.Info("Shutting down OpenTelemetry")
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down tracer provider", slog.String("error", err.Error()))
		}
	}, nil
}
