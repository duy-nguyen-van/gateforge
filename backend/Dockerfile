# Build stage
FROM golang:1.26-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Tidy up dependencies and build the application
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:3.23.3

# Install ca-certificates and create non-root user
RUN apk --no-cache add ca-certificates && \
    adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .

# Change ownership to appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/v1/ || exit 1

# Run the application
CMD ["./main"]
