# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-pmb-service/go.mod services/unsia-pmb-service/go.sum* ./services/unsia-pmb-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-pmb-service

# Copy the service source code
COPY services/unsia-pmb-service/ ./

# Download dependencies
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=1 go build -o /unsia-pmb-service ./cmd/pmb-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-pmb-service /app/unsia-pmb-service

EXPOSE 8003

CMD ["/app/unsia-pmb-service"]
