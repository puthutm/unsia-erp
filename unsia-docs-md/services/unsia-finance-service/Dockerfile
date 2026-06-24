# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.mod.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /unsia-finance-service ./cmd/finance-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-finance-service /app/unsia-finance-service

EXPOSE 8005

CMD ["/app/unsia-finance-service"]
