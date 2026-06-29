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

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-pmb-service ./cmd/pmb-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /unsia-pmb-service /app/unsia-pmb-service
COPY services/unsia-pmb-service/.env .

# Expose port 8003
EXPOSE 8003

# Command to run
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-pmb-service"]
