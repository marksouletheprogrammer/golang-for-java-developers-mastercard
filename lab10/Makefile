# Makefile for Lab 10 - Deployment-Ready Build
# 
# Windows Command Prompt Note:
# The 'make' command is not available in Windows Command Prompt by default.
# Options: (1) Use WSL, Git Bash, or PowerShell, (2) Install Make for Windows,
# or (3) run the commands inside each target directly (e.g., the commands under 'build:').

.PHONY: proto proto-clean build run test clean docker docker-build docker-run version

# TODO: Part 5 - Add build targets for cross-compilation
# Version information (can be overridden)
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags with version injection
LDFLAGS := -ldflags="-w -s -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/orders.proto

# Clean generated protobuf files
proto-clean:
	rm -f proto/orders/*.pb.go

# TODO: Part 5 - Build the server binary with version info
# Build the server binary
build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/server cmd/server/main.go

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/server-linux-amd64 cmd/server/main.go

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/server-darwin-amd64 cmd/server/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/server-darwin-arm64 cmd/server/main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/server-windows-amd64.exe cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Run with environment variables
run-prod:
	ENVIRONMENT=production go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detector
test-race:
	go test -v -race ./...

# TODO: Part 6 - Build Docker image with version info
# Docker commands
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t order-service:$(VERSION) \
		-t order-service:latest \
		..

docker-run:
	docker run -p 8080:8080 -p 9090:9090 --env-file .env.example order-service:latest

# TODO: Part 7 - Use docker-compose for local development
docker-compose-up:
	docker-compose up --build

docker-compose-down:
	docker-compose down

# Display version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Development helpers
dev:
	@echo "Starting development environment..."
	ENVIRONMENT=development LOG_LEVEL=debug go run cmd/server/main.go

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run ./...

# Tidy dependencies
tidy:
	go mod tidy

# Install required tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
