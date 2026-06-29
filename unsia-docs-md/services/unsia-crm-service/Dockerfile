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

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-crm-service ./cmd/crm-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /unsia-crm-service /app/unsia-crm-service
COPY services/unsia-crm-service/.env .

# Expose port 8009
EXPOSE 8009

# Command to run
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-crm-service"]
