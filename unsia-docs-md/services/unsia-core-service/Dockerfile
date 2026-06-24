# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.mod.sum ./
COPY ../packages/*/go.mod ../packages/
COPY ../packages/*/go.sum ../packages/ || true

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=1 go build -o /unsia-core-service ./cmd/core-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Create keys directory
RUN mkdir -p /app/keys

# Copy binary from builder
COPY --from=builder /unsia-core-service /app/unsia-core-service

# Expose port
EXPOSE 8001

# Run the service
CMD ["/app/unsia-core-service"]
