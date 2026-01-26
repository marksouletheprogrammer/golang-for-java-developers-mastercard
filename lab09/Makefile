# Makefile for Lab 9 Order Service with Profiling
# 
# Windows Command Prompt Note:
# The 'make' command is not available in Windows Command Prompt by default.
# Options: (1) Use WSL, Git Bash, or PowerShell, (2) Install Make for Windows,
# or (3) run the commands inside each target directly (e.g., the commands under 'proto:').

.PHONY: proto run test clean install-tools bench bench-mem profile-cpu profile-mem

# Generate protobuf code
# Note: Requires protoc-gen-go and protoc-gen-go-grpc in PATH
proto:
	mkdir -p proto/orders
	protoc --go_out=. --go_opt=module=lab09 \
		--go-grpc_out=. --go-grpc_opt=module=lab09 \
		proto/orders.proto

# Run the service (assumes proto files already generated)
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Run benchmarks
bench:
	go test -bench=. -benchmem ./profiling

# Run benchmarks with memory allocation tracking
bench-mem:
	go test -bench=. -benchmem -memprofile=mem.prof ./profiling

# Run CPU profiling
profile-cpu:
	go test -bench=BenchmarkBatchProcessing -cpuprofile=cpu.prof ./profiling
	go tool pprof cpu.prof

# Run memory profiling
profile-mem:
	go test -bench=BenchmarkBatchProcessing -memprofile=mem.prof ./profiling
	go tool pprof mem.prof

# Clean generated files
clean:
	rm -f proto/orders/*.pb.go
	rm -f order-service
	rm -f *.prof
	rm -f profiling/*.test

# Install required tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
