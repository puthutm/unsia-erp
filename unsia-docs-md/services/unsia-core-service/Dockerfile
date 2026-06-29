# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-core-service/go.mod services/unsia-core-service/go.sum* ./services/unsia-core-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-core-service

# Copy the service source code
COPY services/unsia-core-service/ ./

# Download dependencies
RUN go mod tidy

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-core-service ./cmd/core-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /unsia-core-service /app/unsia-core-service
COPY services/unsia-core-service/.env .

# Expose port 8001
EXPOSE 8001

# Command to run
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-core-service"]
