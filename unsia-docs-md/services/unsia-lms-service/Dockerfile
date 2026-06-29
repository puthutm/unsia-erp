# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-lms-service/go.mod services/unsia-lms-service/go.sum* ./services/unsia-lms-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-lms-service

# Copy the service source code
COPY services/unsia-lms-service/ ./

# Download dependencies
RUN go mod tidy

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-lms-service ./cmd/lms-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /unsia-lms-service /app/unsia-lms-service
COPY services/unsia-lms-service/.env .

# Expose port 8006
EXPOSE 8006

# Command to run
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-lms-service"]
