# unsia-portal-service (Backend)

**Stack:** Go 1.22+ · Gin · GORM + sqlc · golang-migrate

**Database:** `portal_db` (PostgreSQL)

## Tanggung Jawab

Modul Portal backend mengelola **notification center, dashboard read model, preferensi user, dan shortcut**. Portal adalah presentation layer — **bukan source transaksi bisnis**.

| Domain | Deskripsi |
|--------|-----------|
| Notification | Pesan notifikasi per user dari semua modul |
| Read Marker | Status baca/belum baca per notifikasi per user |
| Dashboard Read Model | Agregasi data lintas modul untuk dashboard |
| Preferensi User | Pengaturan tampilan dan notifikasi |
| Shortcut | Pintasan menu per role |
| Activity Log | Jejak aktivitas user di portal |

## Endpoint Utama

```
GET    /api/v1/portal/notifications
PATCH  /api/v1/portal/notifications/{id}/read
GET    /api/v1/portal/dashboard
GET    /api/v1/portal/dashboard/executive
GET    /api/v1/portal/shortcuts
PUT    /api/v1/portal/preferences
POST   /api/v1/portal/notifications    (internal — dipanggil oleh modul lain)
```

## Struktur Direktori (akan diisi saat development)

```
unsia-portal-service/
├── cmd/portal-service/main.go
├── internal/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   ├── handler/
│   └── middleware/
├── migrations/
├── tests/
├── Dockerfile
├── .env.example
└── go.mod
```

## Dependencies

- **Upstream:** Core (auth, user snapshot)
- **Event consume:** Semua modul (notification events, dashboard payload events)
- **API call:** —

## Aturan Penting

- Dashboard **wajib menampilkan `refreshed_at`** — data bisa dari read model, bukan real-time
- Portal **tidak boleh menjadi source status bisnis** — hanya read model dan notifikasi
- Jika Portal down, modul sumber **tetap berjalan normal** — notifikasi diproses ulang setelah pulih
- Read marker bekerja **per user, per notifikasi**

## Owner

Phase 3, Sprint 8.
