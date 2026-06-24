# unsia-integration-worker

**Stack:** Go 1.22+ · RabbitMQ · amqp091-go

**Database:** Tidak memiliki DB sendiri — membaca/menulis ke `outbox_events` dan `inbox_events` di masing-masing database modul.

## Tanggung Jawab

Service terpisah yang menjalankan seluruh proses **event asynchronous** lintas modul: outbox publisher, inbox consumer, retry, DLQ, reconciliation, dan snapshot refresh.

| Worker | Deskripsi |
|--------|-----------|
| Outbox Publisher | Membaca `outbox_events` PENDING dari semua DB modul dan publish ke RabbitMQ |
| Inbox Consumer | Menerima event dari RabbitMQ, cek duplikat, proses ke snapshot/read model |
| DLQ Replay | Replay event dari Dead Letter Queue dengan reason dan audit |
| Reconciliation | Cek selisih antara snapshot dan source of truth secara periodik |
| Snapshot Refresh | Refresh read model yang stale |
| Notification Dispatcher | Forward notification event ke Portal |

## Struktur Direktori (akan diisi saat development)

```
unsia-integration-worker/
├── cmd/integration-worker/main.go
├── internal/
│   ├── workers/
│   │   ├── outbox_publisher.go
│   │   ├── inbox_consumer.go
│   │   ├── dlq_replay.go
│   │   ├── reconciliation.go
│   │   └── snapshot_refresh.go
│   ├── consumers/
│   │   ├── finance_payment_paid.go
│   │   ├── pmb_ready_for_academic.go
│   │   ├── academic_student_created.go
│   │   ├── academic_krs_approved.go
│   │   ├── lms_grade_input_submitted.go
│   │   └── assessment_result_calculated.go
│   ├── publishers/
│   │   └── event_bus.go
│   ├── queues/
│   │   └── rabbitmq.go
│   └── clients/
│       ├── core.client.go
│       ├── pmb.client.go
│       ├── finance.client.go
│       ├── academic.client.go
│       ├── lms.client.go
│       └── portal.client.go
├── tests/
├── Dockerfile
├── .env.example
└── go.mod
```

## Retry Policy

| Attempt | Delay |
|---------|-------|
| 1 | 5 detik |
| 2 | 30 detik |
| 3 | 2 menit |
| 4 | 10 menit |
| 5 | 30 menit |
| > 5 | Masuk DLQ |

## Aturan Penting

- Event duplikat (berdasarkan `event_key`) ditandai `IGNORED_DUPLICATE` dan tidak diproses ulang
- DLQ replay wajib menyertakan `reason` dan dicatat di audit log
- Health check endpoint wajib ada untuk monitoring lag queue dan koneksi RabbitMQ

## Owner

Phase 0, Sprint 0 — harus tersedia sebelum modul lain mulai development.
