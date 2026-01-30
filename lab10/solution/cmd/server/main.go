package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"lab10/config"
	"lab10/internal/repository"
	"lab10/internal/service"
	grpcTransport "lab10/internal/transport/grpc"
	httpTransport "lab10/internal/transport/http"
	pb "lab10/proto"
)

// version information injected at build time via -ldflags
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	// TODO: Part 1 - Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Display startup information
	fmt.Println("=== Order Management Service ===")
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Build Time: %s\n", buildTime)
	fmt.Printf("Environment: %s\n", cfg.Environment)
	fmt.Printf("HTTP Port: %s\n", cfg.HTTPPort)
	if cfg.Features.EnableGRPC {
		fmt.Printf("gRPC Port: %s\n", cfg.GRPCPort)
	}
	fmt.Println()

	// Start server
	if err := runServer(cfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// runServer starts the HTTP and gRPC servers with graceful shutdown.
// TODO: Part 4 - Implement graceful shutdown
func runServer(cfg *config.Config) error {
	// Create dependencies (dependency injection)
	repo := repository.NewMemoryRepository()
	orderService := service.NewOrderService(repo)
	httpHandler := httpTransport.NewOrderHandler(orderService)
	grpcServer := grpcTransport.NewOrderServer(orderService)

	// Setup HTTP server
	httpMux := http.NewServeMux()

	// TODO: Part 8 - Add health check endpoints
	if cfg.Features.EnableHealthz {
		httpMux.HandleFunc("/health", handleHealth)
		httpMux.HandleFunc("/ready", handleReady)
		httpMux.HandleFunc("/version", handleVersion)
	}

	// Register HTTP routes
	httpMux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			httpHandler.CreateOrder(w, r)
		case http.MethodGet:
			httpHandler.ListOrders(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	httpMux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > len("/orders/") && r.URL.Path[len(r.URL.Path)-7:] == "/status" {
			httpHandler.UpdateOrderStatus(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			httpHandler.GetOrder(w, r)
		case http.MethodDelete:
			httpHandler.DeleteOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      httpMux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// TODO: Part 2 - Setup gRPC server only if feature flag is enabled
	var grpcListener net.Listener
	var grpcSrv *grpc.Server

	if cfg.Features.EnableGRPC {
		var err error
		grpcListener, err = net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			return fmt.Errorf("failed to listen on gRPC port: %w", err)
		}

		grpcSrv = grpc.NewServer()
		pb.RegisterOrderServiceServer(grpcSrv, grpcServer)
		reflection.Register(grpcSrv)
	}

	// TODO: Part 4 - Start servers in goroutines
	go func() {
		fmt.Printf("HTTP server listening on :%s\n", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	if cfg.Features.EnableGRPC {
		go func() {
			fmt.Printf("gRPC server listening on :%s\n", cfg.GRPCPort)
			if err := grpcSrv.Serve(grpcListener); err != nil {
				log.Fatalf("gRPC server error: %v", err)
			}
		}()
	}

	// TODO: Part 4 - Implement graceful shutdown
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down servers...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server gracefully
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		fmt.Println("HTTP server stopped gracefully")
	}

	// Shutdown gRPC server gracefully
	if cfg.Features.EnableGRPC && grpcSrv != nil {
		grpcSrv.GracefulStop()
		fmt.Println("gRPC server stopped gracefully")
	}

	fmt.Println("Service shutdown complete")
	return nil
}

// handleHealth returns basic health status.
// TODO: Part 8 - Implement health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	// TODO: Return JSON with status: "ok", version, uptime
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// handleReady checks if the service is ready to accept traffic.
// TODO: Part 8 - Implement readiness check
func handleReady(w http.ResponseWriter, r *http.Request) {
	// TODO: Check dependencies (e.g., database connection)
	// TODO: Return 200 if ready, 503 if not
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
	})
}

// handleVersion returns version information.
// TODO: Part 5 - Implement version endpoint
func handleVersion(w http.ResponseWriter, r *http.Request) {
	// TODO: Return JSON with version, commit, buildTime
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version":   version,
		"commit":    commit,
		"buildTime": buildTime,
	})
}
