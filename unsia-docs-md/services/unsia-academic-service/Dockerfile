# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod/go.sum of the service and the shared packages
COPY services/unsia-academic-service/go.mod services/unsia-academic-service/go.sum* ./services/unsia-academic-service/
COPY packages/ ./packages/

# Change working directory to the service directory
WORKDIR /app/services/unsia-academic-service

# Copy the service source code
COPY services/unsia-academic-service/ ./

# Download dependencies
RUN go mod tidy

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-academic-service ./cmd/academic-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /unsia-academic-service /app/unsia-academic-service
COPY services/unsia-academic-service/.env .

# Expose port 8004
EXPOSE 8004

# Command to run
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-academic-service"]
