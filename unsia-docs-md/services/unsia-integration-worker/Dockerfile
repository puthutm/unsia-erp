# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-integration-worker/go.mod services/unsia-integration-worker/go.sum* ./services/unsia-integration-worker/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-integration-worker

# Copy the service source code
COPY services/unsia-integration-worker/ ./

# Download dependencies
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=1 go build -o /unsia-integration-worker ./cmd/integration-worker

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-integration-worker /app/unsia-integration-worker

CMD ["/app/unsia-integration-worker"]
