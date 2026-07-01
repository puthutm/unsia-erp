# PRD Final - Portal Module

## 1. Pendahuluan & Latar Belakang
Modul Portal bertindak sebagai agregator antar muka pengguna, pusat notifikasi (*notification center*), manajemen preferensi pengguna, pintasan modul (*shortcuts*), serta read model teragregasi untuk penyajian data dashboard bagi pimpinan institusi.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `portal_db`
* **Source of Truth (Kepemilikan Data)**:
  * `portal.notifications` (Pesan notifikasi pengguna).
  * `portal.notification_reads` (Log penanda status baca notifikasi).
  * `portal.user_preferences` (Pengaturan bahasa/tema visual).
  * `portal.menu_shortcuts` (Pintasan akses menu cepat).
  * `portal.portal_activity_logs` (Rekaman jejak aktivitas user di portal).
  * `portal.dashboard_read_models` (Agregasi draf data visual grafik pimpinan).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `user_snapshots` & `role_snapshots` dari Core.
* `dashboard_read_models`: Salinan data proyeksi KPI lintas modul (PMB, Akademik, Finance) hasil event projection untuk menyajikan grafik performa dashboard eksekutif.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-POR-001** | P0 | Larangan Mutasi Bisnis Langsung. | Portal dilarang mengubah status bisnis sumber; seluruh update didelegasikan ke API modul pemilik. |
| **PRD-POR-002** | P1 | Widget Data Freshness. | Setiap grafik data di dashboard pimpinan wajib menampilkan label `refreshed_at` secara jelas. |
| **PRD-POR-003** | P0 | Notifikasi Idempotent. | Penyaluran notifikasi transaksi wajib diamankan menggunakan `event_key` agar tidak tampil ganda di inbox pengguna. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Dikonsumsi (Consumer)**:
  * Mengonsumsi seluruh event bisnis penting (seperti `finance.payment_paid`, `academic.krs_approved`) untuk digenerasikan sebagai pesan notifikasi pengguna secara idempotent.

## 6. Degraded Mode & Resilience Guardrails
* **Source Outage Resilience**: Jika salah satu database modul transaksi down, modul Portal tetap dapat diakses. Widget dashboard yang merujuk pada modul yang down akan menampilkan status offline beserta label waktu sinkronisasi data snapshot terakhir.
