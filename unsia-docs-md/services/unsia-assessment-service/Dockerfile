# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-assessment-service/go.mod services/unsia-assessment-service/go.sum* ./services/unsia-assessment-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-assessment-service

# Copy the service source code
COPY services/unsia-assessment-service/ ./

# Download dependencies
RUN go mod tidy

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-assessment-service ./cmd/assessment-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the compiled binary
COPY --from=builder /unsia-assessment-service /app/unsia-assessment-service

USER appuser

# Expose port 8007
EXPOSE 8007

# Start the service (use env vars or mounted config instead of .env)
CMD ["/app/unsia-assessment-service"]
