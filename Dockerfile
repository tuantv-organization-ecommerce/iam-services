# Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o bin/iam-service cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/iam-service .

# Copy migrations (optional, for reference)
COPY --from=builder /app/migrations ./migrations

# Expose gRPC port
EXPOSE 50051

# Run the application
CMD ["./iam-service"]

