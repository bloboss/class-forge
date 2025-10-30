# Forgejo Classroom Makefile

# Variables
APP_NAME := forgejo-classroom
CLI_NAME := fgc
SERVER_NAME := fgc-server
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Docker related variables
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := $(VERSION)

# Test related variables
TEST_TIMEOUT := 10m
TEST_PACKAGES := ./...
INTEGRATION_PACKAGES := ./test/...

.PHONY: help build clean test test-unit test-integration test-contract lint vet fmt deps update-deps run-cli run-server docker-build docker-test-up docker-test-down migrate-up migrate-down

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build all binaries
build: build-cli build-server

## build-cli: Build CLI binary
build-cli:
	@echo "Building CLI binary..."
	@mkdir -p $(GOBIN)
	go build $(LDFLAGS) -o $(GOBIN)/$(CLI_NAME) ./cmd/fgc

## build-server: Build server binary
build-server:
	@echo "Building server binary..."
	@mkdir -p $(GOBIN)
	go build $(LDFLAGS) -o $(GOBIN)/$(SERVER_NAME) ./cmd/fgc-server

## clean: Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(GOBIN)
	@go clean ./...

## test: Run all tests
test: test-unit test-integration

## test-unit: Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -timeout $(TEST_TIMEOUT) -race -coverprofile=coverage.out $(TEST_PACKAGES)

## test-integration: Run integration tests (requires Docker)
test-integration: docker-test-up
	@echo "Running integration tests..."
	@sleep 5  # Wait for services to be ready
	go test -timeout $(TEST_TIMEOUT) -tags=integration $(INTEGRATION_PACKAGES) || (make docker-test-down && exit 1)
	@make docker-test-down

## test-contract: Run contract tests against OpenAPI spec
test-contract:
	@echo "Running contract tests..."
	@echo "Contract tests not yet implemented"

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet $(TEST_PACKAGES)

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt $(TEST_PACKAGES)

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

## update-deps: Update dependencies
update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## run-cli: Run CLI application (with sample command)
run-cli: build-cli
	@echo "Running CLI application..."
	$(GOBIN)/$(CLI_NAME) --help

## run-server: Run API server
run-server: build-server
	@echo "Running API server..."
	$(GOBIN)/$(SERVER_NAME)

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest

## docker-test-up: Start test environment
docker-test-up:
	@echo "Starting test environment..."
	docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for services to be ready..."
	@timeout 60 bash -c 'until docker-compose -f docker-compose.test.yml exec -T postgres-test pg_isready -U fgc_test; do sleep 1; done'
	@timeout 30 bash -c 'until docker-compose -f docker-compose.test.yml exec -T redis-test redis-cli ping; do sleep 1; done'

## docker-test-down: Stop test environment
docker-test-down:
	@echo "Stopping test environment..."
	docker-compose -f docker-compose.test.yml down -v

## docker-dev-up: Start development environment
docker-dev-up:
	@echo "Starting development environment..."
	docker-compose up -d

## docker-dev-down: Stop development environment
docker-dev-down:
	@echo "Stopping development environment..."
	docker-compose down

## migrate-up: Run database migrations (up)
migrate-up:
	@echo "Running database migrations (up)..."
	@echo "Migration tool not yet implemented - migrations are in ./migrations/"

## migrate-down: Run database migrations (down)
migrate-down:
	@echo "Running database migrations (down)..."
	@echo "Migration tool not yet implemented - migrations are in ./migrations/"

## ci: Run CI checks (lint, vet, test)
ci: deps fmt vet lint test

## install: Install binaries to GOPATH/bin
install: build
	@echo "Installing binaries..."
	cp $(GOBIN)/$(CLI_NAME) $(shell go env GOPATH)/bin/
	cp $(GOBIN)/$(SERVER_NAME) $(shell go env GOPATH)/bin/

## coverage: Generate test coverage report
coverage: test-unit
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## benchmark: Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem $(TEST_PACKAGES)

## dev-setup: Set up development environment
dev-setup: deps docker-test-up
	@echo "Development environment ready!"
	@echo "Run 'make run-server' to start the API server"
	@echo "Run 'make run-cli' to test the CLI"

## release: Build release binaries for multiple platforms
release:
	@echo "Building release binaries..."
	@mkdir -p releases
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o releases/$(CLI_NAME)-linux-amd64 ./cmd/fgc
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o releases/$(CLI_NAME)-linux-arm64 ./cmd/fgc
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o releases/$(CLI_NAME)-darwin-amd64 ./cmd/fgc
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o releases/$(CLI_NAME)-darwin-arm64 ./cmd/fgc
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o releases/$(CLI_NAME)-windows-amd64.exe ./cmd/fgc
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o releases/$(SERVER_NAME)-linux-amd64 ./cmd/fgc-server
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o releases/$(SERVER_NAME)-darwin-amd64 ./cmd/fgc-server
	@echo "Release binaries built in ./releases/"

# Default target
.DEFAULT_GOAL := help