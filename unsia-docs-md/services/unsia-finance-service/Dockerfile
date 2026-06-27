# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-finance-service/go.mod services/unsia-finance-service/go.sum* ./services/unsia-finance-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-finance-service

# Copy the service source code
COPY services/unsia-finance-service/ ./

# Download dependencies
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=1 go build -o /unsia-finance-service ./cmd/finance-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-finance-service /app/unsia-finance-service

EXPOSE 8005

CMD ["/app/unsia-finance-service"]
