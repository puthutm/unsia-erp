# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-reference-service/go.mod services/unsia-reference-service/go.sum* ./services/unsia-reference-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-reference-service

# Copy the service source code
COPY services/unsia-reference-service/ ./

# Download dependencies
RUN go mod tidy

# Build binary
RUN CGO_ENABLED=1 go build -o /unsia-reference-service ./cmd/reference-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-reference-service /app/unsia-reference-service

EXPOSE 8002

CMD ["/app/unsia-reference-service"]
