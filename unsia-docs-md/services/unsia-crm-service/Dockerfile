# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.23.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk add --no-cache git

# Copy the Go Modules
COPY services/unsia-crm-service/go.mod services/unsia-crm-service/go.sum* ./services/unsia-crm-service/
COPY packages/ ./packages/

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the entire source code
COPY services/unsia-crm-service/ ./

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-crm-service ./cmd/crm-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /unsia-crm-service /app/unsia-crm-service

# Expose port 8009
EXPOSE 8009

# Command to run the executable
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-crm-service"]
