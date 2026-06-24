# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.mod.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /unsia-portal-service ./cmd/portal-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-portal-service /app/unsia-portal-service

EXPOSE 8010

CMD ["/app/unsia-portal-service"]
