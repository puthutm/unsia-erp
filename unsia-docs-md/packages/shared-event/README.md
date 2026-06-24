# shared-event

**Type:** Go module (shared library)

## Fungsi

Library untuk outbox writer, inbox consumer, event envelope builder, retry handler, dan DLQ — digunakan semua service Go UNSIA.

## Event Envelope Standard

```go
type EventEnvelope struct {
    EventName       string    // "finance.payment_paid"
    EventVersion    string    // "v1"
    EventKey        string    // deterministik: "finance.payment_paid:uuid:v1"
    PublisherService string
    AggregateType   string
    AggregateID     string
    CorrelationID   string
    CausationID     string
    OccurredAt      time.Time
    Payload         any
}
```

## Fungsi Tersedia

| Fungsi | Deskripsi |
|--------|-----------|
| `WriteOutbox(tx, event)` | Tulis event ke `outbox_events` dalam transaksi yang sama dengan domain update |
| `ConsumeInbox(event)` | Cek duplikat berdasarkan `event_key`, proses jika baru |
| `BuildEventKey(name, id, version)` | Generate deterministik event key |
| `MarkPublished(outboxID)` | Update status outbox ke PUBLISHED |
| `MarkIgnored(inboxID)` | Mark inbox sebagai IGNORED_DUPLICATE |
| `SendToDLQ(event, err)` | Kirim ke Dead Letter Queue setelah retry habis |
