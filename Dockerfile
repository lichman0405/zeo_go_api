# Multi-stage build for Go Zeo++ API
FROM golang:1.24-bullseye AS builder

# Install required packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    git \
    wget \
    make \
    tar \
    && rm -rf /var/lib/apt/lists/*

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
FROM debian:bullseye-slim

# Install Zeo++ dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    wget \
    make \
    tar \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install Zeo++
WORKDIR /opt

# Define environment variables for Zeo++ versioning
ENV ZEO_VERSION=0.3
ENV ZEO_FILENAME=zeo++-${ZEO_VERSION}.tar.gz
ENV ZEO_URL=http://www.zeoplusplus.org/${ZEO_FILENAME}
ENV ZEO_FOLDER=zeo++-${ZEO_VERSION}

# Download and build Zeo++
RUN wget ${ZEO_URL} && \
    tar -xzf ${ZEO_FILENAME} && \
    cd ${ZEO_FOLDER}/voro++/src && make && \
    cd ../.. && make && \
    mv network /usr/local/bin/network && \
    chmod +x /usr/local/bin/network && \
    rm -rf ${ZEO_FILENAME} ${ZEO_FOLDER}

# Copy the binary from builder stage
COPY --from=builder /app/zeo-api /usr/local/bin/zeo-api

# Create workspace directory
RUN mkdir -p /app/workspace /app/config

# Copy config file
COPY --from=builder /app/config/config.yaml /app/config/config.yaml

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["/usr/local/bin/zeo-api"]