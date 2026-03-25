# ============================================================
# Multi-stage Dockerfile for Mini Bank API
# Tech Stack: Go Echo + Oracle DB + Celery
# ============================================================

# Stage 1: Go builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/api \
    ./cmd/api/main.go

# ============================================================

# Stage 2: Go API image
FROM alpine:3.19 AS api

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/api .

# Copy config
COPY config.yaml .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run
CMD ["./api"]

# ============================================================

# Stage 3: Python Celery worker image
FROM python:3.11-slim AS worker

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libffi-dev \
    libssl-dev \
    redis-tools \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python deps
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy worker code
COPY workers/ ./workers/
COPY cmd/worker/ ./cmd/worker/

# Set Python path
ENV PYTHONPATH=/app

# Run Celery worker
CMD ["python", "-m", "celery", "-A", "workers.celery_app", \
     "worker", "--loglevel=info", "--concurrency=4"]
