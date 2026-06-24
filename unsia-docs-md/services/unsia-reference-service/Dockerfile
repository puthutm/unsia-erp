# ==============================
# Stage 1: Builder
# ==============================
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.mod.sum ./
COPY ../packages/*/go.mod ../packages/ || true
COPY ../packages/*/go.sum ../packages/ || true

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /unsia-reference-service ./cmd/reference-service

# ==============================
# Stage 2: Runner
# ==============================
FROM alpine:3.19 AS runner

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /unsia-reference-service /app/unsia-reference-service

EXPOSE 8002

CMD ["/app/unsia-reference-service"]
