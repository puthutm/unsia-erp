# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.23.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk add --no-cache git

ARG NETRC_FILE=.netrc

# Salin file netrc (sementaa untuk build)
COPY ${NETRC_FILE} /root/.netrc

# Amankan file
RUN chmod 600 /root/.netrc

# Copy the Go Modules and Sum files
COPY services/unsia-finance-service/go.mod services/unsia-finance-service/go.sum* ./services/unsia-finance-service/
COPY packages/ ./packages/

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the entire source code
COPY services/unsia-finance-service/ ./

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-finance-service ./cmd/finance-service

RUN rm -f /root/.netrc

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /unsia-finance-service /app/unsia-finance-service

# Expose port 8005
EXPOSE 8005

# Command to run the executable
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-finance-service"]
