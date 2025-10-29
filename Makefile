# Makefile for iam-services
# Quick commands for development

.PHONY: help lint lint-fix lint-fast test build clean check-all

# Default target
help:
	@echo "Available commands:"
	@echo "  make lint       - Run golangci-lint"
	@echo "  make lint-fix   - Run golangci-lint with auto-fix"
	@echo "  make lint-fast  - Run golangci-lint in fast mode"
	@echo "  make test       - Run all tests"
	@echo "  make build      - Build the project"
	@echo "  make check-all  - Run all checks (lint + build + test)"
	@echo "  make clean      - Clean build artifacts"

# Linting
lint:
	golangci-lint run --config .golangci.yml

lint-fix:
	golangci-lint run --config .golangci.yml --fix

lint-fast:
	golangci-lint run --config .golangci.yml --fast

# Lint specific packages
lint-model:
	golangci-lint run --config .golangci.yml internal/domain/model/...

lint-handler:
	golangci-lint run --config .golangci.yml internal/handler/...

lint-dao:
	golangci-lint run --config .golangci.yml internal/dao/...

lint-service:
	golangci-lint run --config .golangci.yml internal/service/...

# Testing
test:
	go test -v -race -coverprofile=coverage.out ./...

test-short:
	go test -v -short ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Building
build:
	go build -v -o bin/iam-service ./cmd/server

build-linux:
	GOOS=linux GOARCH=amd64 go build -v -o bin/iam-service-linux ./cmd/server

# Development
run:
	go run cmd/server/main.go

# Complete check before push
check-all: lint build test
	@echo "✅ All checks passed! Ready to push."

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out
	go clean -cache -testcache

# Docker
docker-build:
	docker build -t iam-service:latest .

docker-run:
	docker-compose up -d

# Proto generation
proto:
	./scripts/generate-proto.ps1
