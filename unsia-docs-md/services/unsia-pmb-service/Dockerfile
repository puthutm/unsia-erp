# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.23.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk add --no-cache git

# Copy the Go Modules
COPY services/unsia-pmb-service/go.mod services/unsia-pmb-service/go.sum* ./services/unsia-pmb-service/
COPY packages/ ./packages/

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the entire source code
COPY services/unsia-pmb-service/ ./

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-pmb-service ./cmd/pmb-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /unsia-pmb-service /app/unsia-pmb-service

# Expose port 8003
EXPOSE 8003

# Command to run the executable
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-pmb-service"]
