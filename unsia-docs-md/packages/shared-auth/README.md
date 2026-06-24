# shared-auth

**Type:** Go module (shared library)

## Fungsi

Library auth yang digunakan semua service Go UNSIA agar validasi JWT, JWKS caching, dan service token konsisten.

| Fungsi | Deskripsi |
|--------|-----------|
| `ValidateJWT(token string)` | Validasi JWT RS256 menggunakan public key dari JWKS cache |
| `FetchJWKS(url string)` | Fetch dan cache JSON Web Key Set dari Core Service |
| `ValidateServiceToken(token string)` | Validasi machine-to-machine service token |
| `ExtractClaims(token string)` | Extract claims: `sub`, `active_role`, `permissions`, `scope` |
| `RefreshJWKSIfStale()` | Refresh JWKS jika TTL cache sudah lewat (default 5 menit) |

## Penggunaan

```go
import "github.com/unsia-erp/shared-auth"

claims, err := sharedauth.ValidateJWT(bearerToken)
if err != nil {
    // return 401
}
```

## Degraded Mode

Jika Core Service tidak tersedia, JWKS cache tetap digunakan selama TTL belum expired. Token yang sudah expired tetap ditolak.
