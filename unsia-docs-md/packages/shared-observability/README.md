# shared-observability

**Type:** Go module (shared library)

## Fungsi

Structured logging, trace ID propagation, correlation ID middleware, dan metrics exporter untuk semua service Go UNSIA.

## Komponen

| Komponen | Library | Deskripsi |
|----------|---------|-----------|
| Structured Logging | `zerolog` atau `zap` | JSON log dengan level, trace_id, service_name |
| Tracing | OpenTelemetry | Distributed tracing lintas service |
| Metrics | Prometheus | HTTP metrics, event lag, DB connection pool |
| Correlation ID | Middleware Gin | Propagate `X-Correlation-Id` ke semua request |

## Middleware Gin

```go
import "github.com/unsia-erp/shared-observability"

r := gin.New()
r.Use(sharedobservability.CorrelationIDMiddleware())
r.Use(sharedobservability.RequestLoggerMiddleware())
r.Use(sharedobservability.MetricsMiddleware())
```

## Log Format

```json
{
  "level": "info",
  "service": "pmb-service",
  "trace_id": "3fb4b7f1-...",
  "correlation_id": "abc123",
  "method": "POST",
  "path": "/api/v1/pmb/applicants",
  "status": 201,
  "duration_ms": 45,
  "message": "request completed"
}
```
