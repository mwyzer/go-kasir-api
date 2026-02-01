# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install git (needed for go mod download)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first (layer caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o kasir-api .

# Stage 2: Run
FROM alpine:latest

# Install CA certificates for HTTPS (needed for database)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/kasir-api .

# Expose port (Zeabur will override with PORT env)
EXPOSE 8080

# Run
CMD ["./kasir-api"]