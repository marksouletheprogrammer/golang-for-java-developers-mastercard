package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
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

var cpuProfileFile *os.File

func main() {
	// Parse profiling flags
	profileMode := flag.String("profile", "", "Enable profiling mode: cpu, mem, or benchmark")
	flag.Parse()

	// Handle profiling modes
	switch *profileMode {
	case "cpu":
		startCPUProfile()
		defer stopCPUProfile()
	case "mem":
		defer writeMemProfile()
	case "benchmark":
		printBenchmarkInfo()
		return
	case "":
		// Normal server mode, continue
	default:
		fmt.Printf("Unknown profile mode: %s\n", *profileMode)
		fmt.Println("Valid modes: cpu, mem, benchmark")
		os.Exit(1)
	}

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

// startCPUProfile starts CPU profiling and writes to cpu.prof.
func startCPUProfile() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatalf("Could not create CPU profile: %v", err)
	}
	cpuProfileFile = f

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		log.Fatalf("Could not start CPU profile: %v", err)
	}

	fmt.Println("CPU profiling enabled. Profile will be written to cpu.prof")
	fmt.Println("Server will run for 30 seconds, then shut down...")

	// Auto-shutdown after 30 seconds for CPU profiling
	go func() {
		time.Sleep(30 * time.Second)
		fmt.Println("\nProfiling complete. Shutting down...")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
}

// stopCPUProfile stops CPU profiling.
func stopCPUProfile() {
	pprof.StopCPUProfile()
	if cpuProfileFile != nil {
		cpuProfileFile.Close()
	}
	fmt.Println("CPU profile written to cpu.prof")
	fmt.Println("Analyze with: go tool pprof cpu.prof")
}

// writeMemProfile writes a memory profile to mem.prof.
func writeMemProfile() {
	runtime.GC() // Get up-to-date statistics

	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatalf("Could not create memory profile: %v", err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatalf("Could not write memory profile: %v", err)
	}

	fmt.Println("Memory profile written to mem.prof")
	fmt.Println("Analyze with: go tool pprof mem.prof")
}

// printBenchmarkInfo prints benchmark instructions and exits.
func printBenchmarkInfo() {
	fmt.Println("=== Benchmark Commands ===")
	fmt.Println()
	fmt.Println("Run all benchmarks:")
	fmt.Println("  go test -bench=. -benchmem ./profiling")
	fmt.Println()
	fmt.Println("Run specific benchmark:")
	fmt.Println("  go test -bench=BenchmarkOrderValidation -benchmem ./profiling")
	fmt.Println()
	fmt.Println("Generate CPU profile from benchmark:")
	fmt.Println("  go test -bench=BenchmarkBatchProcessing -cpuprofile=cpu.prof ./profiling")
	fmt.Println("  go tool pprof cpu.prof")
	fmt.Println()
	fmt.Println("Generate memory profile from benchmark:")
	fmt.Println("  go test -bench=BenchmarkBatchProcessing -memprofile=mem.prof ./profiling")
	fmt.Println("  go tool pprof mem.prof")
	fmt.Println()
	fmt.Println("Interactive pprof commands:")
	fmt.Println("  top      - Show top functions")
	fmt.Println("  list FN  - Show annotated source")
	fmt.Println("  web      - Generate graph (requires graphviz)")
	fmt.Println("  quit     - Exit pprof")
}
