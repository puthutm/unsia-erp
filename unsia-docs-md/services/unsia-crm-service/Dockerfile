# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-crm-service/go.mod services/unsia-crm-service/go.sum* ./services/unsia-crm-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-crm-service

# Copy the service source code
COPY services/unsia-crm-service/ ./

# Download dependencies
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=1 go build -o /unsia-crm-service ./cmd/crm-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-crm-service /app/unsia-crm-service

EXPOSE 8009

CMD ["/app/unsia-crm-service"]
