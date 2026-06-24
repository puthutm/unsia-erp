# shared-httpclient

**Type:** Go module (shared library)

## Fungsi

HTTP client untuk service-to-service call dengan timeout, retry, dan circuit breaker bawaan.

## Konfigurasi Default

| Parameter | Default |
|-----------|---------|
| Timeout per request | 5 detik |
| Max retries | 3 kali |
| Retry delay | Exponential backoff (1s, 2s, 4s) |
| Circuit breaker threshold | 5 consecutive failures → open |
| Circuit breaker reset | 30 detik |

## Penggunaan

```go
import "github.com/unsia-erp/shared-httpclient"

client := sharedhttpclient.New(sharedhttpclient.Config{
    BaseURL: "http://finance-service",
    ServiceToken: os.Getenv("FINANCE_SERVICE_TOKEN"),
    Timeout: 5 * time.Second,
})

resp, err := client.Post("/api/v1/finance/invoices", payload)
```

## Headers yang Otomatis Ditambahkan

- `Authorization: Bearer <service_token>`
- `X-Correlation-Id` (diteruskan dari context request)
- `X-Source-Service: {service_name}`
