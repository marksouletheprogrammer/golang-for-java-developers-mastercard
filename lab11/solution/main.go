package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lab11/observability"
	"lab11/server"
)

func main() {
	// Setup structured logging
	logger := setupLogger()
	logger.Info("Starting order service with observability")
	
	// Initialize OpenTelemetry
	shutdownTelemetry, err := observability.InitTelemetry(logger)
	if err != nil {
		logger.Error("Failed to initialize telemetry", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer shutdownTelemetry()
	
	// Create HTTP server with observability
	srv := server.NewServer(logger)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}
	
	// Start HTTP server
	go func() {
		logger.Info("Starting HTTP server", slog.String("port", "8080"))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
	
	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	
	logger.Info("Shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", slog.String("error", err.Error()))
	}
	
	logger.Info("Server stopped")
}

// setupLogger configures structured logging for the application.
// Uses log/slog with text format for development.
func setupLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
