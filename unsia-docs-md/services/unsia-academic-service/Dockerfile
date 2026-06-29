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
COPY services/unsia-academic-service/go.mod services/unsia-academic-service/go.sum* ./services/unsia-academic-service/
COPY packages/ ./packages/

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the entire source code
COPY services/unsia-academic-service/ ./

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-academic-service ./cmd/academic-service

RUN rm -f /root/.netrc

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /unsia-academic-service /app/unsia-academic-service
COPY services/unsia-academic-service/.env .

# Expose port 8004
EXPOSE 8004

# Command to run the executable
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-academic-service"]
