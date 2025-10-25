.PHONY: help proto build run test clean docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/iam.proto
	@echo "Done!"

build: ## Build the application
	@echo "Building..."
	go build -o bin/iam-service cmd/server/main.go
	@echo "Build complete: bin/iam-service"

run: ## Run the application
	@echo "Running IAM Service..."
	go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	@echo "Clean complete!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies downloaded!"

db-create: ## Create database
	@echo "Creating database..."
	psql -U postgres -c "CREATE DATABASE iam_db;"
	@echo "Database created!"

db-migrate: ## Run database migrations
	@echo "Running migrations..."
	psql -U postgres -d iam_db -f migrations/001_init_schema.sql
	psql -U postgres -d iam_db -f migrations/002_seed_data.sql
	@echo "Migrations complete!"

db-drop: ## Drop database
	@echo "Dropping database..."
	psql -U postgres -c "DROP DATABASE IF EXISTS iam_db;"
	@echo "Database dropped!"

db-reset: db-drop db-create db-migrate ## Reset database

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t iam-service:latest .
	@echo "Docker image built!"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 50051:50051 --env-file .env iam-service:latest

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted!"

install-tools: ## Install development tools
	@echo "Installing tools..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@echo "Tools installed!"

