# Build stage
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/app

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy web assets
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]
