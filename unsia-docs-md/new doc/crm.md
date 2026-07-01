# PRD Final - CRM Module

## 1. Pendahuluan & Latar Belakang
Modul CRM mengoptimalkan pencatatan leads/calon pendaftar dari berbagai kanal pemasaran (sosial media, agen kemitraan, referral, event) dan memantau status tindak lanjut (*follow-up pipeline*) hingga dikonversi menjadi pendaftar PMB.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `crm_db`
* **Source of Truth (Kepemilikan Data)**:
  * `crm.campaigns` (Kampanye promosi).
  * `crm.agents` (Agen kemitraan rujukan).
  * `crm.referrals` (Kode rujukan agen).
  * `crm.leads` (Prospek calon pendaftar).
  * `crm.lead_activities` (Log interaksi/follow-up).
  * `crm.lead_status_histories` (Histori transisi status leads).
  * `crm.commission_rules` & `crm.commission_records` (Parameter komisi rujukan agen).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `person_snapshots`: Berisi salinan data nama dan kontak lead yang bersumber dari Core.
* Menyimpan `applicant_ref_id` setelah prospek dikonversi menjadi applicant resmi di modul PMB.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-CRM-001** | P0 | Kemandirian Pencatatan Leads. | CRM menyimpan leads secara mandiri dan tidak bergantung langsung pada PMB untuk aktivitas follow-up. |
| **PRD-CRM-002** | P0 | Konversi Prospek Idempotent. | Konversi prospek dilakukan via API/event PMB secara idempotent, bukan manipulasi database PMB langsung. |
| **PRD-CRM-003** | P1 | Rekaman Referral Komisi Agen. | Pencatatan data pendaftar yang dirujuk terhubung otomatis dengan aturan komisi untuk dikirim ke modul Finance. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `crm.lead_qualified`: Dikirim saat prospek lolos kualifikasi dan siap dikonversi.
    * *Payload Minimum*: `lead_id`, `person_ref_id`, `source_code`, `campaign_id`, `agent_id`, `occurred_at`.
* **Event yang Dikonsumsi (Consumer)**:
  * `pmb.applicant_created`: Menghubungkan lead dengan applicant resmi dan merubah status di CRM menjadi "Converted".

## 6. Degraded Mode & Resilience Guardrails
* **PMB Outage Resilience**: Jika PMB mengalami gangguan, proses input leads dan follow-up di CRM tetap dapat berjalan normal. Handover leads yang lolos kualifikasi ditunda di tabel outbox CRM dan akan dikirim ulang secara otomatis setelah PMB online.
