# shared-idempotency

**Type:** Go module (shared library)

## Fungsi

Library untuk mencegah duplikasi command kritis (payment callback, handover, generate NIM, class sync, dll).

## Cara Kerja

```
Request masuk dengan Idempotency-Key header
    ↓
Cek idempotency_keys table
    ↓
Jika sudah ada → return cached response (HTTP 200)
    ↓
Jika belum ada → lock key → proses command → simpan response → release lock
```

## Penggunaan

```go
import "github.com/unsia-erp/shared-idempotency"

result, cached, err := sharedidempotency.CheckAndLock(ctx, idempotencyKey, ttl)
if cached {
    return result // return response dari cache
}
// proses command...
sharedidempotency.SaveResponse(ctx, idempotencyKey, response, ttl)
```

## Command yang WAJIB Memakai Idempotency-Key

- Payment callback (`finance.payment-callbacks`)
- PMB handover ke Academic
- Generate NIM
- Class sync (LMS dari Academic)
- Enrollment sync
- Grade sync
- Replay DLQ
- OAuth token issuance dari authorization code
