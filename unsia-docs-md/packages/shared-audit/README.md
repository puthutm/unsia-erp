# shared-audit

**Type:** Go module (shared library)

## Fungsi

Library pencatatan audit log untuk semua aksi sensitif di service Go UNSIA.

## Schema Audit Log

```go
type AuditEntry struct {
    ID            string    // UUID
    Actor         string    // user_id
    ActiveRole    string    // role yang digunakan saat aksi
    Action        string    // contoh: "pmb.applicant.handover"
    Module        string    // contoh: "pmb"
    ResourceType  string    // contoh: "applicant"
    ResourceID    string    // UUID resource
    OldValue      any       // JSON sebelum perubahan (nullable)
    NewValue      any       // JSON setelah perubahan (nullable)
    Reason        string    // alasan jika aksi sensitif (nullable)
    CorrelationID string    // X-Correlation-Id dari request
    IPAddress     string
    UserAgent     string
    OccurredAt    time.Time
}
```

## Aksi yang Wajib Diaudit

- Create, update, deactivate/delete data penting
- Approve / reject dokumen, invoice, KRS
- Handover PMB → Academic
- Generate NIM
- Payment verify / override clearance
- Grade correction
- Impersonation
- OAuth client approve / suspend / revoke
- DLQ replay

## Penggunaan

```go
import "github.com/unsia-erp/shared-audit"

sharedaudit.Log(ctx, sharedaudit.AuditEntry{
    Action:    "pmb.applicant.handover",
    ResourceType: "applicant",
    ResourceID: applicantID,
    Reason:    reason,
})
```
