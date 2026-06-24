---
title: "Event Contract ERP UNSIA"
source_file: "OpenAPI_Swagger_Final_ERP_UNSIA_v1_0_1_Event_Contract_Updated.json"
format: markdown
---

# Event Contract ERP UNSIA

Dokumen ini diturunkan dari OpenAPI/Event Contract baseline dan dipakai sebagai kontrak integrasi event, outbox/inbox, DLQ replay, dan reconciliation.

## Prinsip Event Contract

- Setiap event wajib memiliki `event_name`, `event_version`, dan `event_key` deterministik.
- Producer menulis event ke outbox setelah transaksi domain berhasil commit.
- Consumer wajib mencatat event pada inbox dan memproses secara idempotent.
- Duplicate event tidak boleh membuat data ganda.
- Event gagal diproses harus masuk retry dan DLQ sesuai policy.
- Replay DLQ wajib reason dan audit.
- Snapshot/read model wajib memiliki `source_event_key`, `source_module`, dan `synced_at`.

## Event/Integration Endpoints

| Method | Path | Operation ID | Summary | Idempotent | Permission/Role |
|---|---|---|---|---|---|
| `GET` | `/api/v1/integration/event-contracts` | `listEventContracts` | List event contract catalog | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/event-contracts/{event_name}` | `getEventContractByName` | Detail event contract berdasarkan event_name | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/outbox-events` | `listOutboxEvents` | List outbox events | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/inbox-events` | `listInboxEvents` | List inbox events | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/dlq-events` | `listDlqEvents` | List DLQ events | `False` | technical_admin \| devops_sre \| auditor |
| `POST` | `/api/v1/integration/dlq-events/{event_key}/replay` | `replayDlqEvent` | Replay event dari DLQ | `True` | devops_sre |
| `GET` | `/api/v1/integration/reconciliation-mismatches` | `listReconciliationMismatches` | List reconciliation mismatch logs | `False` | technical_admin \| devops_sre \| auditor |
| `POST` | `/api/v1/integration/reconciliation-mismatches/{mismatch_id}/resolve` | `resolveReconciliationMismatch` | Resolve reconciliation mismatch | `True` | technical_admin \| devops_sre \| owner_modul |

## Event Envelope Standard

```json
{
  "event_name": "finance.payment_paid",
  "event_version": "v1",
  "event_key": "finance.payment_paid:payment_uuid:v1",
  "publisher_service": "finance-service",
  "aggregate_type": "payment",
  "aggregate_id": "uuid",
  "correlation_id": "uuid",
  "causation_id": "uuid-or-provider-event-id",
  "occurred_at": "2026-06-22T10:00:00+07:00",
  "payload": {}
}
```

## Event Catalog Minimum

| Event | Publisher | Consumer |
|---|---|---|
| `core.person_updated` | Core | CRM, PMB, HRIS, Academic, Portal |
| `reference.study_program_updated` | Referensi | PMB, Academic, HRIS, LMS, Portal |
| `reference.academic_period_updated` | Referensi | PMB, Finance, Academic, LMS, Assessment, Portal |
| `crm.lead_qualified` | CRM | PMB |
| `pmb.applicant_created` | PMB | Finance, Assessment, Portal |
| `finance.invoice_created` | Finance | PMB, Portal |
| `finance.payment_paid` | Finance | PMB, Academic, Portal |
| `finance.clearance_changed` | Finance | PMB, Academic, LMS, Portal |
| `pmb.ready_for_academic` | PMB | Academic |
| `academic.student_created` | Academic | PMB, Finance, LMS, Portal |
| `academic.class_opened` | Academic | LMS, Portal |
| `academic.krs_approved` | Academic | LMS, Finance, Portal |
| `lms.grade_input_submitted` | LMS | Academic |
| `assessment.result_calculated` | Assessment | PMB, LMS, Academic, Portal |