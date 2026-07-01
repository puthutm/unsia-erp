# PRD Final - Core Module

## 1. Pendahuluan & Latar Belakang
Modul Core bertanggung jawab sebagai otoritas tunggal untuk manajemen identitas pengguna (*single identity*), otentikasi (OIDC/JWT), manajemen peran (*role*), perizinan (*permission*), pendaftaran aplikasi (*app registry*), dan manajemen klien layanan (*service client*). Desain v6.5 mewajibkan pembagian database fisik terpisah (`core_db`) dengan isolasi penuh tanpa foreign key eksternal ke database modul lain.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `core_db`
* **Source of Truth (Kepemilikan Data)**:
  * `core.persons`: Menyimpan profil identitas individu fisik (nama lengkap, gender, NIK, kontak).
  * `core.users`: Menyimpan credential akun (username, email, hash password, status akun aktif/inactive/suspended).
  * `core.roles`: Master data peran sistem beserta tipe cakupan data (*scope type*).
  * `core.permissions`: Daftar perizinan otorisasi dengan format `module.resource.action`.
  * `core.user_roles` & `core.role_permissions`: Penugasan peran ke pengguna dan izin ke peran.
  * `core.sessions` & `core.active_role_sessions`: Manajemen sesi otentikasi dan peran aktif.
  * `core.idempotency_keys` & `core.integration_event_logs`: Tabel pengaman idempotency transaksi.
  * `core.audit_logs`: Log jejak audit global sensitif.

## 3. Data Lintas Modul (Snapshot/Read Model)
Modul Core bertindak sebagai landasan sistem dan tidak menyimpan data snapshot dari modul lain secara lokal. Sebaliknya, modul Core membagikan snapshot data personal dan perizinan pengguna ke seluruh modul hilir.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-CORE-001** | P0 | Core menyediakan token OIDC/JWT yang membawa klaim user info, active role, dan perizinan. | Token JWT dapat divalidasi oleh modul lain secara mandiri menggunakan Cached JWKS/Public Key tanpa query ke `core_db`. |
| **PRD-CORE-002** | P0 | Penerbitan event perubahan data personal pengguna. | Perubahan pada nama, email, telepon di `core.persons` menerbitkan event `core.person_updated` untuk disinkronkan ke snapshot modul lain. |
| **PRD-CORE-003** | P1 | Endpoint Token Introspection real-time. | Core menyediakan endpoint introspection untuk kasus verifikasi sesi sensitif (seperti plot nilai akhir) dengan fallback lokal di sisi client. |
| **PRD-CORE-004** | P0 | Log Audit Keamanan Terpusat. | Setiap manipulasi data sensitif (impersonasi akun, pergantian password) wajib tercatat dengan detail old/new value, actor, request_id, dan IP address. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `core.person_updated`: Dikirim ketika profil diri individu diubah.
    * *Payload Minimum*: `person_id`, `full_name`, `email`, `phone`, `status_code`, `occurred_at`.
  * `core.user_role_changed`: Dikirim ketika pemetaan peran pengguna mengalami penyesuaian.
    * *Payload Minimum*: `user_id`, `role_id`, `study_program_id`, `action_type` (add/remove), `occurred_at`.

## 6. Degraded Mode & Resilience Guardrails
* **Core Downtime Fallback**: Saat database Core down, pengguna yang sudah memiliki sesi login aktif tetap dapat mengakses sistem. Modul-modul hilir ERP memvalidasi token JWT secara lokal menggunakan Cached Public Key/JWKS yang memiliki TTL (Time-to-Live) pendek. Write sensitif dibatasi, dan login baru ditahan hingga Core pulih.
