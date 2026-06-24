# unsia-crm

**Stack:** Next.js 14+ (App Router) · TypeScript · TanStack Query · Tailwind CSS

**Backend:** `unsia-crm-service` → `crm_db`

## Tanggung Jawab

Frontend modul CRM — manajemen lead, campaign, agen, follow-up, dan konversi ke PMB.

## Route Structure

```
app/
├── (agent)/                       → Agent / Mitra (own_lead scope)
│   ├── dashboard/
│   ├── leads/
│   │   ├── page.tsx               → Hanya lead milik sendiri
│   │   ├── buat/
│   │   └── [id]/
│   └── komisi/
│
└── (admin)/                       → Admin CRM / Marketing
    ├── dashboard/
    │   └── funnel/                → Pipeline akuisisi
    ├── campaigns/
    │   ├── page.tsx
    │   └── [id]/
    ├── leads/
    │   ├── page.tsx               → List semua lead + pipeline
    │   └── [id]/
    │       ├── follow-ups/
    │       └── convert/           → Convert ke applicant PMB
    ├── agents/
    └── komisi/
        ├── page.tsx
        └── approval/
```

## Role yang Dilayani

| Role | Scope | Akses Utama |
|------|-------|-------------|
| Agent / Mitra | `own_lead` | Buat lead, lihat lead sendiri, cek komisi |
| Admin CRM / Marketing | CRM domain | Semua lead, campaign, pipeline, komisi |
| Pimpinan | Read-only | Dashboard funnel akuisisi |

## Integrasi API

- `unsia-crm-service` — semua data CRM
- `unsia-pmb-service` — convert lead → applicant (via API)
- `unsia-core-service` — auth/token

## Aturan UI

- Agent **hanya bisa melihat lead miliknya sendiri** — backend enforced, bukan hanya hide UI
- Tombol "Convert ke Applicant" hanya aktif jika status lead = `QUALIFIED`
- Convert bersifat **idempotent** — tombol disabled setelah berhasil + tampilkan `applicant_ref_id`
- Dashboard funnel menampilkan `refreshed_at` untuk data agregat
