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

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-reference-service ./cmd/reference-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /unsia-reference-service /app/unsia-reference-service
COPY services/unsia-reference-service/.env .

# Expose port 8002
EXPOSE 8002

# Command to run
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-reference-service"]
