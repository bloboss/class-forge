# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN make build

# Runtime stage for API server
FROM alpine:latest AS api-server

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl

# Create app user
RUN addgroup -S fgc && adduser -S fgc -G fgc

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/fgc-server /usr/local/bin/fgc-server

# Copy configuration example
COPY --from=builder /app/config.yaml.example /app/config.yaml.example

# Change ownership
RUN chown -R fgc:fgc /app

# Switch to app user
USER fgc

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run server
CMD ["fgc-server"]

# Runtime stage for CLI
FROM alpine:latest AS cli

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create app user
RUN addgroup -S fgc && adduser -S fgc -G fgc

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/fgc /usr/local/bin/fgc

# Change ownership
RUN chown -R fgc:fgc /app

# Switch to app user
USER fgc

# Run CLI
ENTRYPOINT ["fgc"]