# shared-errorenvelope

**Type:** Go module (shared library)

## Fungsi

Format error response konsisten untuk seluruh service Go UNSIA.

## Format Response

**Success:**
```json
{
  "success": true,
  "message": "Request processed successfully",
  "data": {},
  "meta": {
    "trace_id": "3fb4b7f1-7d28-4d13-a812-9cc5e1c0c011",
    "timestamp": "2026-06-22T10:00:00+07:00"
  }
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN_SCOPE",
    "message": "Anda tidak memiliki akses ke data ini.",
    "details": {}
  },
  "meta": {
    "trace_id": "3fb4b7f1-7d28-4d13-a812-9cc5e1c0c011",
    "timestamp": "2026-06-22T10:00:00+07:00"
  }
}
```

## Penggunaan

```go
import "github.com/unsia-erp/shared-errorenvelope"

// Success response
c.JSON(200, sharederr.Success(data))

// Error response
c.JSON(403, sharederr.Error("FORBIDDEN_SCOPE", "Anda tidak memiliki akses ke data ini."))

// Validation error
c.JSON(400, sharederr.ValidationError(fieldErrors))
```
