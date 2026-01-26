# Makefile for Lab 8 Order Service
# 
# Windows Command Prompt Note:
# The 'make' command is not available in Windows Command Prompt by default.
# Options: (1) Use WSL, Git Bash, or PowerShell, (2) Install Make for Windows,
# or (3) run the commands inside each target directly (e.g., the commands under 'proto:').

.PHONY: proto run test clean install-tools

# Generate protobuf code
# Note: Requires protoc-gen-go and protoc-gen-go-grpc in PATH
proto:
	mkdir -p proto/orders
	protoc --go_out=. --go_opt=module=lab08 \
		--go-grpc_out=. --go-grpc_opt=module=lab08 \
		proto/orders.proto

# Run the service (assumes proto files already generated)
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Clean generated files
clean:
	rm -f proto/orders/*.pb.go
	rm -f order-service

# Install required tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
