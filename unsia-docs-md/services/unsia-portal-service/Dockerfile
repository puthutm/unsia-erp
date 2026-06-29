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
COPY services/unsia-portal-service/go.mod services/unsia-portal-service/go.sum* ./services/unsia-portal-service/
COPY packages/ ./packages/

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the entire source code
COPY services/unsia-portal-service/ ./

# Build the Go app
RUN CGO_ENABLED=1 go build -o /unsia-portal-service ./cmd/portal-service

RUN rm -f /root/.netrc

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /unsia-portal-service /app/unsia-portal-service
COPY services/unsia-portal-service/.env .
COPY services/unsia-portal-service/*.pem .

# Expose port 8010
EXPOSE 8010

# Command to run the executable
CMD ["/bin/sh", "-c", "apk update && apk add bash"]
CMD ["/app/unsia-portal-service"]
