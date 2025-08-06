# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Binary names
BINARY_NAME=mcp-feishu
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_DARWIN=$(BINARY_NAME)_darwin

# Version and build info
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date '+%Y-%m-%d_%H:%M:%S')
COMMIT=$(shell git rev-parse HEAD)

# LDFLAGS for embedding version info
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

.PHONY: all build clean test coverage deps fmt vet lint help

# Default target
all: test build

# Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_DARWIN)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	$(GOTEST) -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	$(GOFMT) -s -w .

# Vet code
vet:
	$(GOVET) ./...

# Lint code (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) -v

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DARWIN) -v

# Build for all platforms
build-all: build-linux build-windows build-darwin

# Docker targets
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest

docker-run:
	docker run --rm -it \
		-e FEISHU_WEBHOOK_URL=$(FEISHU_WEBHOOK_URL) \
		-e FEISHU_SECURITY_TYPE=$(FEISHU_SECURITY_TYPE) \
		$(BINARY_NAME):latest

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

# Development targets
dev: deps fmt vet lint test build

# Install the binary
install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Create release archives
release: build-all
	mkdir -p release
	tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_UNIX) README.md LICENSE examples/
	zip -r release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_WINDOWS) README.md LICENSE examples/
	tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_DARWIN) README.md LICENSE examples/
	@echo "Release packages created in release/ directory"

# Run with example config
run-example:
	cp examples/env.example .env
	@echo "Please edit .env file with your Feishu webhook URL, then run:"
	@echo "make run"

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) && ./$(BINARY_NAME)

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage report"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  build-all    - Cross compile for all platforms"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  dev          - Run all development checks"
	@echo "  install      - Install binary to /usr/local/bin"
	@echo "  release      - Create release packages"
	@echo "  run-example  - Copy example config"
	@echo "  run          - Build and run the application"
	@echo "  help         - Show this help message"
