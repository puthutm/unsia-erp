# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.23.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk add --no-cache git

# Copy the Go Modules
COPY services/unsia-core-service/go.mod services/unsia-core-service/go.sum* ./services/unsia-core-service/
COPY packages/ ./packages/

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the entire source code
COPY services/unsia-core-service/ ./

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-core-service ./cmd/core-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /unsia-core-service /app/unsia-core-service

# Expose port 8001
EXPOSE 8001

# Command to run the executable
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-core-service"]
