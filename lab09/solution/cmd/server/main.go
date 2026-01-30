package main

import (
	"context"
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

	commonconfig "golang-for-java-developers-training/common/config"
	"lab09/internal/repository"
	"lab09/internal/service"
	grpcTransport "lab09/internal/transport/grpc"
	httpTransport "lab09/internal/transport/http"
	pb "lab09/proto/orders"
)

// Config holds application configuration.
// Loaded from environment variables with sensible defaults.
type Config struct {
	HTTPPort string
	GRPCPort string
}

// loadConfig loads configuration from environment variables.
func loadConfig() Config {
	return Config{
		HTTPPort: commonconfig.GetEnv("HTTP_PORT", "8080"),
		GRPCPort: commonconfig.GetEnv("GRPC_PORT", "9090"),
	}
}

func main() {
	// Load configuration
	config := loadConfig()

	fmt.Println("Starting Order Management Service")
	fmt.Printf("HTTP Port: %s\n", config.HTTPPort)
	fmt.Printf("gRPC Port: %s\n", config.GRPCPort)

	// Create dependencies (dependency injection)
	// Outer layers depend on inner layers, never the reverse.
	// Repository -> Service -> Transport (HTTP + gRPC)

	// 1. Create repository (innermost layer)
	repo := repository.NewMemoryRepository()

	// 2. Create service (inject repository)
	orderService := service.NewOrderService(repo)

	// 3. Create HTTP transport (inject service)
	httpHandler := httpTransport.NewOrderHandler(orderService)

	// 4. Create gRPC transport (inject same service)
	grpcServer := grpcTransport.NewOrderServer(orderService)

	// Setup HTTP server
	httpMux := http.NewServeMux()

	// Register HTTP routes
	// POST /orders - create order
	// GET /orders - list all orders
	// GET /orders/{id} - get specific order
	// PATCH /orders/{id}/status - update order status
	// DELETE /orders/{id} - delete order
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
		// Check if it's a status update
		if len(r.URL.Path) > len("/orders/") && r.URL.Path[len(r.URL.Path)-7:] == "/status" {
			httpHandler.UpdateOrderStatus(w, r)
			return
		}

		// Otherwise it's a single order operation
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
		Addr:    ":" + config.HTTPPort,
		Handler: httpMux,
	}

	// Setup gRPC server
	grpcListener, err := net.Listen("tcp", ":"+config.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcSrv, grpcServer)

	// Register reflection service for grpcurl and other tools
	reflection.Register(grpcSrv)

	// Start both servers in goroutines
	go func() {
		fmt.Printf("HTTP server listening on :%s\n", config.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		fmt.Printf("gRPC server listening on :%s\n", config.GRPCPort)
		if err := grpcSrv.Serve(grpcListener); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down servers...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		fmt.Println("HTTP server stopped gracefully")
	}

	// Shutdown gRPC server
	grpcSrv.GracefulStop()
	fmt.Println("gRPC server stopped gracefully")

	fmt.Println("Service shutdown complete")
}
