# shared-rbac

**Type:** Go module (shared library)

## Fungsi

Library RBAC untuk permission check dan data scope resolver yang digunakan semua service Go UNSIA.

| Fungsi | Deskripsi |
|--------|-----------|
| `CheckPermission(claims, permission string)` | Cek apakah active role memiliki permission tertentu |
| `ResolveDataScope(claims)` | Resolve data scope berdasarkan active role (global, study_program, self, assigned_class, own_lead) |
| `EnforceScope(resource, scope)` | Validasi bahwa resource yang diakses sesuai data scope user |

## Pattern Permission

```
module.resource.action
Contoh: pmb.applicant.verify_document
         academic.krs.approve
         finance.clearance.override
```

## Penggunaan

```go
import "github.com/unsia-erp/shared-rbac"

if err := sharedrbac.CheckPermission(claims, "pmb.applicant.verify_document"); err != nil {
    // return 403
}
scope := sharedrbac.ResolveDataScope(claims)
```
