# Multi-stage build for Go Zeo++ API
FROM golang:1.21-alpine AS builder

# Install required packages
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o zeo-api cmd/server/main.go

# Final stage
FROM alpine:latest

# Install Zeo++ dependencies
RUN apk add --no-cache \
    g++ \
    make \
    wget \
    git \
    bash

# Install Zeo++
WORKDIR /opt
RUN wget http://www.zeoplusplus.org/zeo++-0.3.tar.gz && \
    tar -xzf zeo++-0.3.tar.gz && \
    cd zeo++-0.3 && \
    make && \
    make install

# Copy the binary from builder stage
COPY --from=builder /app/zeo-api /usr/local/bin/zeo-api

# Create non-root user
RUN addgroup -g 1001 -S zeo && \
    adduser -S zeo -u 1001

# Create workspace directory
RUN mkdir -p /app/workspace /app/config && \
    chown -R zeo:zeo /app

# Copy config file
COPY --from=builder /app/config/config.yaml /app/config/config.yaml

# Switch to non-root user
USER zeo

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["/usr/local/bin/zeo-api"]