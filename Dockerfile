# Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy gokits dependency first (needed for local replace directive)
COPY gokits/ ./gokits/

# Copy iam-services files
WORKDIR /app/iam-services

# Copy go mod files
COPY iam-services/go.mod iam-services/go.sum ./
RUN go mod download

# Copy source code
COPY iam-services/ ./

# Build the application
RUN go build -o bin/iam-service cmd/server/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/iam-services/bin/iam-service .

# Copy migrations (optional, for reference)
COPY --from=builder /app/iam-services/migrations ./migrations

# Copy configs (required for Casbin models)
COPY --from=builder /app/iam-services/configs ./configs

# Expose gRPC port and HTTP port
EXPOSE 50051
EXPOSE 8080

# Run the application
CMD ["./iam-service"]

