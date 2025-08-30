# =============================================================================
# COVID-19 KMS Backend Go Application Container
# =============================================================================

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Final stage
FROM alpine:latest

# Accept Railway environment variables as build arguments
ARG DATABASE_HOST
ARG DATABASE_PORT
ARG DATABASE_NAME
ARG DATABASE_USER
ARG DATABASE_SSL_MODE

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main ./main

# Copy environment configuration files
COPY --from=builder /app/internal/config/env.example ./env.example

# Set environment variables from build arguments
ENV HOST=0.0.0.0
ENV PORT=8000
ENV ENV=production
ENV DATABASE_HOST=$DATABASE_HOST
ENV DATABASE_PORT=$DATABASE_PORT
ENV DATABASE_NAME=$DATABASE_NAME
ENV DATABASE_USER=$DATABASE_USER
ENV DATABASE_SSL_MODE=$DATABASE_SSL_MODE

# Debug: Print environment variables during build
RUN echo "🔍 Debug: Environment Variables" && \
    echo "HOST: $HOST" && \
    echo "PORT: $PORT" && \
    echo "ENV: $ENV" && \
    echo "DATABASE_HOST: $DATABASE_HOST" && \
    echo "DATABASE_PORT: $DATABASE_PORT" && \
    echo "DATABASE_NAME: $DATABASE_NAME" && \
    echo "DATABASE_USER: $DATABASE_USER" && \
    echo "DATABASE_SSL_MODE: $DATABASE_SSL_MODE"

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/health || exit 1

# Labels
LABEL maintainer="COVID-19 KMS Team"
LABEL description="COVID-19 Knowledge Management System Backend API"
LABEL version="1.0"

# Run the application (Railway will use this)
CMD ["/app/main"]

# Alternative: Use ENTRYPOINT for more explicit control
# ENTRYPOINT ["./main"]
