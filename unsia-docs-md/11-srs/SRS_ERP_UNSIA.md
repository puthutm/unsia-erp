---
title: "SRS ERP UNSIA"
source_file: "SRS_ERP_UNSIA.docx"
format: markdown
---

# SRS ERP UNSIA

# UNSIA

# SOFTWARE REQUIREMENTS SPECIFICATION (SRS)

# ERP Pendidikan / SIAKAD Terintegrasi UNSIA

## Versi 1.0 Draft | 22 Juni 2026

Dokumen ini menyatukan kebutuhan sistem perangkat lunak untuk ERP Pendidikan / SIAKAD Terintegrasi UNSIA sebagai acuan implementasi developer, QA, DBA, DevOps, dan owner modul.

| Item | Isi |
| --- | --- |
| Dokumen | Software Requirements Specification (SRS) |
| Produk | ERP Pendidikan / SIAKAD Terintegrasi UNSIA |
| Versi | v1.0 Draft |
| Basis | PRD Global v6.5.1, BRD v1.1.1, FSD v1.0.1, OpenAPI v1.0.1, Event Contract, DBML, UAT/QA Test Plan, dan dokumen struktur repo. |
| Arsitektur | Modular distributed ERP dengan database fisik per modul, API/event-driven integration, outbox/inbox, idempotency, audit, RBAC, data scope, degraded mode, dan reconciliation. |
| Status | Draft untuk review Product Owner, Technical Lead, DBA, Security/DevOps, QA/UAT Lead, Backend Lead, Frontend Lead, dan Owner Modul. |

# Kontrol Dokumen

| Versi | Tanggal | Status | Catatan |
| --- | --- | --- | --- |
| v1.0 | 22 Juni 2026 | Draft | Penyusunan SRS awal berdasarkan baseline PRD, BRD, FSD, OpenAPI, Event Contract, DBML, UAT, dan rancangan multi-repo. |

## Daftar Approval

| Peran | Nama | Status | Tanggal | Catatan |
| --- | --- | --- | --- | --- |
| Product Owner |  | Belum disetujui |  |  |
| System Analyst |  | Drafted |  |  |
| Technical Lead |  | Belum direview |  |  |
| Backend Lead |  | Belum direview |  |  |
| Frontend Lead |  | Belum direview |  |  |
| DBA |  | Belum direview |  |  |
| Security/DevOps |  | Belum direview |  |  |
| QA/UAT Lead |  | Belum direview |  |  |
| Owner Core |  | Belum direview |  |  |
| Owner PMB |  | Belum direview |  |  |
| Owner Finance |  | Belum direview |  |  |
| Owner Akademik |  | Belum direview |  |  |

# Daftar Isi

Catatan: daftar isi dapat diperbarui otomatis di Microsoft Word melalui References > Update Table setelah dokumen dibuka.

### 1. Pendahuluan

### 2. Deskripsi Umum Sistem

### 3. Arsitektur dan Batasan Teknis

### 4. User Class dan Role

### 5. Kebutuhan Fungsional Global

### 6. Kebutuhan Fungsional Per Modul

### 7. External Interface Requirements

### 8. Data Requirements

### 9. Non-Functional Requirements

### 10. Security, Audit, dan Compliance

### 11. State Machine Requirements

### 12. Testing dan Acceptance Criteria

### 13. Deployment, Operasional, dan Release

### 14. Requirement Traceability Matrix

### 15. Lampiran

# 1. Pendahuluan

## 1.1 Tujuan Dokumen

SRS ini menetapkan kebutuhan perangkat lunak untuk ERP Pendidikan / SIAKAD Terintegrasi UNSIA. Dokumen ini menjadi acuan tunggal bagi tim developer, QA, DBA, DevOps/SRE, Product Owner, dan Owner Modul dalam merancang, membangun, menguji, merilis, dan mengoperasikan sistem.

SRS ini menerjemahkan kebutuhan produk, bisnis, fungsi, API, event, database, role/permission, dan UAT menjadi requirement yang dapat diimplementasikan dan diuji.

## 1.2 Ruang Lingkup Sistem

Sistem mencakup lifecycle kampus dari lead, applicant, pembayaran, mahasiswa aktif, KRS, LMS, assessment, nilai, KHS, transkrip, alumni, dashboard, notification, sampai reporting. Sistem terdiri dari modul Core, Referensi, CRM, PMB, Finance, Akademik, HRIS/SDM, LMS, Assessment, Portal, Integration Worker, dan Reporting/warehouse.

## 1.3 Definisi, Akronim, dan Istilah

| Istilah | Definisi |
| --- | --- |
| SRS | Software Requirements Specification, dokumen kebutuhan perangkat lunak. |
| ERP | Enterprise Resource Planning untuk mengintegrasikan proses operasional kampus. |
| SIAKAD | Sistem Informasi Akademik. |
| Source of Truth | Modul/domain yang menjadi pemilik data utama. |
| External Reference | ID lintas modul seperti person_ref_id, applicant_ref_id, student_ref_id, invoice_ref_id tanpa FK lintas database. |
| Snapshot/Read Model | Salinan ringkas data lintas modul untuk tampilan/operasi terbatas, bukan sumber kebenaran final. |
| Outbox/Inbox | Pola event-driven untuk publish dan consume event secara idempotent. |
| Idempotency | Kemampuan request/event diproses berulang tanpa menghasilkan data ganda. |
| DLQ | Dead Letter Queue untuk event yang gagal diproses setelah retry. |
| Degraded Mode | Mode operasi terbatas saat dependency modul lain down. |
| RBAC | Role-Based Access Control. |
| Data Scope | Batas data yang boleh diakses role tertentu. |

## 1.4 Referensi Dokumen

| Dokumen Referensi | Kegunaan dalam SRS |
| --- | --- |
| PRD Global UNSIA v6.5.1 | Menjadi dasar scope produk, arsitektur database modular, event-driven integration, dan NFR. |
| BRD UNSIA v1.1.1 | Menjadi dasar kebutuhan bisnis, source of truth, role stakeholder, dan aturan global. |
| FSD Per Modul v1.0.1 | Menjadi dasar fungsi per modul, UI standard, audit, event behavior, dan UAT starter. |
| OpenAPI/Swagger v1.0.1 | Menjadi dasar API contract, endpoint, header, security, dan response envelope. |
| UAT Scenario dan QA Test Plan v1.0.1 | Menjadi dasar test level, quality gate, defect severity, dan go/no-go rule. |
| DBML Global v1.0.1 | Menjadi dasar struktur database, outbox/inbox, idempotency, dan reconciliation table. |

## 1.5 Prioritas Requirement

| Kode Prioritas | Makna |
| --- | --- |
| P0 | Wajib untuk MVP/go-live terbatas. Tanpa ini proses utama tidak boleh dirilis. |
| P1 | Penting untuk stabilitas operasional dan release tahap berikutnya. |
| P2 | Enhancement atau optimasi setelah proses utama stabil. |

# 2. Deskripsi Umum Sistem

## 2.1 Perspektif Produk

ERP UNSIA adalah sistem kampus modular yang mengelola proses end-to-end dari peminat sampai alumni. Sistem tidak dibangun sebagai satu aplikasi monolitik dengan satu database utama, melainkan sebagai ekosistem service/modul dengan database fisik per modul. Integrasi lintas modul dilakukan melalui API, event, snapshot/read model, dan reconciliation.

## 2.2 Modul Sistem

| Modul | Fungsi Utama | Database |
| --- | --- | --- |
| Core | Identity, SSO, RBAC, permission, active role, service token, audit, idempotency, integration control. | core_db |
| Referensi | Master data lintas modul, prodi, tahun ajaran, periode akademik, status code, komponen pembayaran. | reference_db |
| CRM | Campaign, lead, agent, referral, follow-up, conversion, commission. | crm_db |
| PMB | Applicant, biodata, dokumen, seleksi, LoA, invoice request, handover akademik. | pmb_db |
| Finance | Invoice, invoice item, payment, callback, manual verification, receipt, clearance, beasiswa. | finance_db |
| Akademik | Student, NIM, kurikulum, mata kuliah, kelas, KRS, final grade, KHS, transkrip, alumni. | academic_db |
| HRIS/SDM | Employee, lecturer, homebase, unit kerja, jabatan, status aktif, BKD. | hris_db |
| LMS | Online class, enrollment, sesi, materi, tugas, presensi, progress, grade input. | lms_db |
| Assessment | Question bank, assessment session, attempt, answer, scoring, result publish. | assessment_db |
| Portal | Dashboard, notification, read marker, shortcut, preference, activity log, dashboard read model. | portal_db |

## 2.3 User Class

| User Class | Kebutuhan Utama |
| --- | --- |
| Pendaftar | Mendaftar, mengisi biodata, upload dokumen, melihat invoice, pembayaran, status, dan LoA. |
| Mahasiswa | Mengelola KRS, melihat tagihan, mengikuti LMS, melihat nilai, KHS, dan transkrip. |
| Dosen | Mengelola kelas LMS, materi, tugas, presensi, dan grade input. |
| Dosen PA | Menyetujui/menolak KRS mahasiswa bimbingan dan memberi catatan akademik. |
| Admin PMB | Mengelola pendaftar, dokumen, seleksi, LoA, dan handover. |
| Admin Finance | Mengelola invoice, payment, verifikasi, clearance, dan laporan keuangan. |
| Admin Akademik Biro | Mengelola mahasiswa, NIM, kurikulum, kelas, KRS, nilai, KHS, transkrip, alumni. |
| Kaprodi/Admin Prodi | Mengelola dan memonitor data akademik sesuai prodi. |
| Admin SDM | Mengelola pegawai, dosen, homebase, jabatan, dan status aktif. |
| Admin LMS | Mengelola kelas online, enrollment, dan sinkronisasi LMS. |
| Admin Assessment | Mengelola bank soal, sesi assessment, scoring, dan result publish. |
| Pimpinan | Melihat dashboard KPI dan laporan agregat. |
| Technical Admin/DevOps/SRE | Mengelola observability, service token, retry, DLQ, reconciliation, dan runbook. |

## 2.4 Asumsi dan Dependensi

Setiap modul memiliki database fisik atau instance/cluster terpisah minimal untuk modul kritikal.

Semua modul menggunakan Core sebagai identity dan access authority.

API dan event contract tersedia sebelum implementasi endpoint P0 dimulai.

Deployment awal dapat menggunakan Docker, Nginx, PostgreSQL, Redis, dan RabbitMQ/BullMQ.

Sistem pembayaran eksternal akan dihubungkan melalui provider callback yang wajib divalidasi signature-nya.

Reporting final lintas modul menggunakan read model/warehouse yang telah direkonsiliasi.

# 3. Arsitektur dan Batasan Teknis

## 3.1 Arsitektur Target

Sistem menggunakan modular distributed architecture dengan service per modul, database per modul, API/event-driven integration, outbox/inbox, idempotency, audit trail, degraded mode, dan reconciliation. Untuk pilihan stack full Next.js, setiap modul dapat dibuat sebagai Next.js/Node service terpisah dengan PostgreSQL masing-masing, sedangkan portal web dan integration worker dipisah.

## 3.2 Struktur Repo Target

unsia-core-service

unsia-reference-service

unsia-crm-service

unsia-pmb-service

unsia-finance-service

unsia-academic-service

unsia-hris-service

unsia-lms-service

unsia-assessment-service

unsia-portal-web

unsia-integration-worker

unsia-shared-contracts

unsia-infra

unsia-docs

## 3.3 Constraint Arsitektur

Tidak ada credential, password, atau session authority selain Core.

Tidak ada write langsung ke database modul lain.

Tidak ada cross-database foreign key.

Tidak ada direct cross-database join untuk transaksi online.

Relasi lintas modul wajib menggunakan external reference.

Setiap snapshot/read model wajib menyimpan source_module, source_event_key, synced_at/refreshed_at, dan reconciliation status bila relevan.

Payment callback, PMB handover, generate NIM, class sync, enrollment sync, grade sync, dan notification delivery wajib idempotent.

Final grade hanya dimiliki Academic.

Finance tetap source of truth untuk invoice, payment, dan clearance.

Portal tidak boleh menjadi source transaksi bisnis.

## 3.4 Technology Stack Rekomendasi

| Layer | Teknologi Rekomendasi |
| --- | --- |
| Frontend/Portal | Next.js, React, TypeScript, Tailwind CSS, Shadcn UI, TanStack Query |
| Service/API | Next.js Route Handler atau Node.js TypeScript service per modul |
| ORM/DB Access | Prisma atau Drizzle |
| Database | PostgreSQL per modul |
| Cache/Lock | Redis |
| Queue/Event | RabbitMQ atau BullMQ + Redis |
| Worker | Node.js TypeScript worker untuk outbox/inbox, retry, DLQ, reconciliation |
| Storage | MinIO/S3-compatible object storage |
| API Docs | OpenAPI/Swagger |
| Observability | Pino/Loki, Prometheus/Grafana, Sentry |
| Deployment | Docker, Nginx/API gateway, CI/CD GitHub Actions/GitLab CI |

# 4. User Role, Permission, dan Data Scope

## 4.1 Role Final

| Role Code | Nama Role | Data Scope Utama |
| --- | --- | --- |
| super_admin | Super Admin | Global |
| admin_bppti | Admin BPPTI | Global teknis sesuai assignment |
| technical_admin | Technical Admin | Technical scope |
| auditor | Auditor | Read-only audit |
| admin_referensi | Admin Referensi | Global referensi |
| admin_crm | Admin CRM/Marketing | CRM domain |
| agen_mitra | Agent/Mitra | Own lead/referral |
| pendaftar | Pendaftar | Self scope |
| admin_pmb | Admin PMB | PMB domain |
| admin_finance | Admin Finance | Finance domain |
| admin_akademik_biro | Admin Akademik Biro | Academic global |
| kaprodi | Kaprodi | study_program_id |
| admin_akademik_prodi | Admin Akademik Prodi | study_program_id |
| dosen | Dosen | Assigned class |
| dosen_pa | Dosen PA | Advisor scope |
| mahasiswa | Mahasiswa | Self scope |
| admin_sdm | Admin SDM/HRIS | HRIS domain |
| admin_lms | Admin LMS | LMS domain |
| admin_assessment | Admin Assessment | Assessment domain |
| pimpinan | Pimpinan | Read-only aggregate |
| service_account | Service Account | Service scope sesuai client |

## 4.2 Permission Naming

Permission menggunakan format module.resource.action, misalnya pmb.applicant.read, finance.payment.verify, academic.krs.approve, lms.grade_input.submit, integration.dlq.replay.

## 4.3 Authorization Flow

Request

- > Validate Authorization Bearer Token

- > Validate X-Application-Code

- > Validate X-Active-Role

- > Load Role Assignment

- > Check Permission

- > Resolve Data Scope

- > Validate Resource Scope

- > Process Request

- > Write Audit Log if sensitive action

- > Return Success/Error Envelope

# 5. Kebutuhan Fungsional Global

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-G-001 | P0 | Sistem harus menyediakan ERP/SIAKAD terintegrasi untuk lifecycle lead sampai alumni. | Data dapat ditelusuri dari CRM, PMB, Finance, Academic, LMS, Assessment, Portal, dan Reporting. |
| SRS-G-002 | P0 | Sistem harus memakai Core sebagai pusat identity dan access. | Semua modul menggunakan token/claim Core dan tidak memiliki credential table sendiri. |
| SRS-G-003 | P0 | Setiap modul harus memiliki database fisik sendiri. | Setiap modul dapat backup, restore, migration, dan recovery terpisah. |
| SRS-G-004 | P0 | Sistem tidak boleh memakai cross-database FK. | ERD fisik hanya memiliki FK internal database; relasi lintas modul memakai external reference. |
| SRS-G-005 | P0 | Sistem tidak boleh memakai direct cross-database join pada transaksi online. | UI transaksi membaca API/snapshot/read model, bukan join lintas database. |
| SRS-G-006 | P0 | Setiap modul harus memiliki outbox/inbox event. | Event dapat dipublish, diterima, diproses, diretry, dan ditelusuri. |
| SRS-G-007 | P0 | Setiap proses kritis harus idempotent. | Retry tidak membuat applicant, student, payment, class, enrollment, grade, atau notification duplikat. |
| SRS-G-008 | P0 | Source of truth per domain harus jelas. | Modul lain tidak mengubah data owner secara langsung. |
| SRS-G-009 | P0 | Sistem harus mendukung degraded mode saat dependency down. | Modul tidak terkait tetap berjalan dan UI menampilkan data terakhir atau pending_review. |
| SRS-G-010 | P1 | Reporting lintas modul harus memakai read model/warehouse. | Dashboard pimpinan tidak query langsung ke semua database transaksi. |
| SRS-G-011 | P1 | Reconciliation lintas modul harus tersedia. | Selisih source vs snapshot/read model dapat dideteksi, dipantau, dan diperbaiki. |
| SRS-G-012 | P1 | Setiap dashboard/read model harus menampilkan waktu sinkronisasi terakhir. | User mengetahui freshness data melalui synced_at/refreshed_at. |

# 6. Kebutuhan Fungsional Per Modul

## 6.1 Modul Core

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-CORE-001 | P0 | Sistem harus menyediakan login dan session user. | User aktif dapat login dan menerima access token/refresh token. |
| SRS-CORE-002 | P0 | Sistem harus menyediakan active role selector. | User multi-role dapat memilih active role dan menu/permission berubah sesuai role. |
| SRS-CORE-003 | P0 | Sistem harus mengelola role, permission, data scope, dan application launcher. | Role/permission/scope dapat diassign dan ditegakkan backend. |
| SRS-CORE-004 | P0 | Sistem harus menyediakan service token untuk service-to-service call. | Service token dapat dibuat, dirotasi, dicabut, dan divalidasi. |
| SRS-CORE-005 | P1 | Sistem harus menyediakan impersonation dengan reason. | Impersonation terbatas durasi, wajib reason, dan semua aksi audit. |

## 6.2 Modul Referensi

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-REF-001 | P0 | Sistem harus mengelola prodi sebagai master data. | Prodi aktif dapat dipakai PMB, Academic, HRIS, dan Portal. |
| SRS-REF-002 | P0 | Sistem harus mengelola Tahun Ajaran dan Periode Akademik. | Periode berada di bawah tahun ajaran dan dapat dipakai PMB, Finance, Academic, LMS. |
| SRS-REF-003 | P0 | Sistem harus mengelola status code lintas modul. | Status transaksi tidak menggunakan string bebas. |
| SRS-REF-004 | P1 | Sistem harus publish event perubahan master data. | Consumer dapat memperbarui reference snapshot. |

## 6.3 Modul CRM

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-CRM-001 | P0 | Sistem harus mengelola campaign, lead, agent, referral, dan follow-up. | Admin CRM dapat mengelola lead; agent hanya lead miliknya. |
| SRS-CRM-002 | P0 | Sistem harus melakukan convert qualified lead ke applicant PMB. | Convert menggunakan PMB API dan idempotent. |
| SRS-CRM-003 | P1 | Sistem harus menyediakan dashboard funnel dan conversion. | Conversion rate dapat dilihat sesuai role/scope. |

## 6.4 Modul PMB

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-PMB-001 | P0 | Sistem harus membuat dan mengelola applicant. | Applicant memiliki biodata, status, prodi tujuan, dan periode masuk. |
| SRS-PMB-002 | P0 | Sistem harus menyediakan upload dan verifikasi dokumen. | Dokumen dapat verified/rejected dengan reason. |
| SRS-PMB-003 | P0 | Sistem harus request invoice ke Finance. | PMB tidak membuat invoice lokal; Finance menjadi source of truth. |
| SRS-PMB-004 | P0 | Sistem harus menerbitkan LoA. | LoA hanya terbit setelah dokumen dan payment policy valid. |
| SRS-PMB-005 | P0 | Sistem harus handover applicant ke Academic. | Handover idempotent dan tidak membuat student/NIM ganda. |

## 6.5 Modul Finance

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-FIN-001 | P0 | Sistem harus membuat invoice dan invoice item. | Invoice memiliki bill_to_ref_id dan status resmi Finance. |
| SRS-FIN-002 | P0 | Sistem harus memproses payment callback. | Signature valid, provider_event_id unique, duplicate callback aman. |
| SRS-FIN-003 | P0 | Sistem harus menyediakan manual payment verification. | Admin Finance dapat verifikasi bukti pembayaran dengan audit. |
| SRS-FIN-004 | P0 | Sistem harus mengelola clearance. | Clearance resmi berasal dari Finance dan dapat dibaca Academic/PMB/LMS/Portal. |

## 6.6 Modul Akademik

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-ACD-001 | P0 | Sistem harus generate student dan NIM dari applicant valid. | applicant_ref_id unique dan NIM unique. |
| SRS-ACD-002 | P0 | Sistem harus mengelola kurikulum, mata kuliah, dan kelas. | Kelas menjadi source untuk LMS class sync. |
| SRS-ACD-003 | P0 | Sistem harus mengelola KRS dan approval Dosen PA. | KRS mematuhi periode aktif dan clearance policy. |
| SRS-ACD-004 | P0 | Sistem harus mengelola source grade dan final grade. | LMS/Assessment hanya input; final grade milik Academic. |
| SRS-ACD-005 | P1 | Sistem harus menghasilkan KHS, transkrip, yudisium, dan alumni. | Mahasiswa dapat melihat/download dokumen sesuai self scope. |

## 6.7 Modul HRIS/SDM

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-HRIS-001 | P0 | Sistem harus mengelola employee dan lecturer. | Dosen aktif dapat dibaca Academic dan LMS. |
| SRS-HRIS-002 | P0 | Sistem harus mengelola homebase, unit kerja, jabatan, dan status aktif. | Dosen nonaktif tidak boleh diplot ke kelas baru. |
| SRS-HRIS-003 | P1 | Sistem harus publish event perubahan status dosen. | Academic dan LMS memperbarui lecturer snapshot. |

## 6.8 Modul LMS

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-LMS-001 | P0 | Sistem harus sync class dari Academic. | Duplicate class sync tidak membuat kelas ganda. |
| SRS-LMS-002 | P0 | Sistem harus sync enrollment dari KRS approved. | Enrollment LMS berasal dari KRS valid dan idempotent. |
| SRS-LMS-003 | P0 | Sistem harus mengelola sesi, materi, tugas, presensi, progress. | Dosen assigned dapat mengelola kelas; mahasiswa enrolled dapat mengikuti. |
| SRS-LMS-004 | P0 | Sistem harus mengirim grade input ke Academic. | Grade sync idempotent dan tidak menimpa final grade. |

## 6.9 Modul Assessment

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-ASM-001 | P0 | Sistem harus mengelola question bank dan versioning soal. | Soal yang sudah dipakai attempt tidak diedit langsung. |
| SRS-ASM-002 | P0 | Sistem harus mengelola assessment session, participant, attempt, answer, dan scoring. | Attempt submitted bersifat immutable. |
| SRS-ASM-003 | P0 | Sistem harus publish result ke context owner. | Result publish idempotent dan hanya menjadi input bagi PMB/LMS/Academic. |

## 6.10 Modul Portal

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-POR-001 | P0 | Sistem harus menyediakan dashboard role-based. | Dashboard mengikuti active role dan menampilkan refreshed_at. |
| SRS-POR-002 | P0 | Sistem harus menyediakan notification center. | Notification idempotent dan user hanya melihat notifikasi miliknya. |
| SRS-POR-003 | P1 | Sistem harus menyediakan executive dashboard agregat. | Pimpinan hanya read-only aggregate dan tidak dapat mengubah transaksi. |

## 6.11 Modul Integration Worker

| ID | Prioritas | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| SRS-INT-001 | P0 | Sistem harus publish outbox event setelah transaksi lokal commit. | Outbox status berubah menjadi published dan traceable. |
| SRS-INT-002 | P0 | Sistem harus consume event dengan inbox idempotent. | event_key diproses satu kali dan duplicate ditandai ignored. |
| SRS-INT-003 | P0 | Sistem harus menyediakan retry, DLQ, dan replay. | Replay wajib reason dan audit, tidak membuat data duplikat. |
| SRS-INT-004 | P1 | Sistem harus menjalankan reconciliation source vs snapshot. | Mismatch dapat ditampilkan, dikoreksi, diabaikan dengan reason, atau pending_review. |

# 7. External Interface Requirements

## 7.1 User Interface Requirements

| Area UI | Requirement |
| --- | --- |
| List Page | Harus mendukung pagination, sorting, search, filter, loading state, empty state, dan export jika role diizinkan. |
| Form Page | Harus menampilkan mandatory marker, validation message spesifik, confirmation untuk aksi sensitif, dan save draft jika proses panjang. |
| Detail Page | Harus menampilkan ringkasan data, status, histori perubahan, related records, dan action button sesuai permission. |
| Snapshot Page | Harus menampilkan source_module, synced_at/refreshed_at, dan label data tidak real-time jika stale. |
| Error State | Harus menampilkan pesan user-friendly dan trace_id untuk support teknis. |
| Degraded Mode UI | Harus menampilkan status dependency dan membatasi aksi berisiko menjadi read-only atau pending_review. |

## 7.2 API Interface Requirements

Authorization: Bearer <access_token>

X-Application-Code: <application_code>

X-Active-Role: <role_code>

X-Correlation-Id: <uuid>

Idempotency-Key: <unique_business_key> # for critical commands

Semua response sukses dan error wajib memakai envelope standar. Endpoint protected harus menolak request tanpa token, active role, permission, atau scope yang valid.

### 7.3 Success Envelope

{

"success": true,

"message": "Request processed successfully",

"data": {},

"meta": {

"trace_id": "uuid",

"correlation_id": "uuid",

"timestamp": "2026-06-22T10:00:00+07:00",

"api_version": "v1"

}

}

### 7.4 Error Envelope

{

"success": false,

"error": {

"code": "FORBIDDEN_SCOPE",

"message": "Anda tidak memiliki akses ke data ini.",

"details": {}

},

"meta": {

"trace_id": "uuid",

"correlation_id": "uuid",

"timestamp": "2026-06-22T10:00:00+07:00",

"api_version": "v1"

}

}

## 7.5 Integration/Event Interface Requirements

| Requirement | Acceptance Criteria |
| --- | --- |
| Setiap event penting harus memiliki event_name, event_version, event_key, publisher, consumer, aggregate_id, correlation_id, causation_id, occurred_at, dan payload. | Event contract test dapat memvalidasi schema dan event identity. |
| Outbox event hanya dipublish setelah transaksi lokal commit. | Rollback transaksi tidak menghasilkan event published. |
| Consumer wajib mencatat event_key di inbox. | Duplicate event diproses satu kali dan ditandai ignored. |
| Event gagal sementara masuk retry queue. | Retry_count, next_retry_at, dan last_error tercatat. |
| Event gagal permanen masuk DLQ. | DLQ dapat direplay dengan reason dan audit. |
| Reconciliation mismatch harus tersedia. | Mismatch source vs snapshot muncul di monitor dan memiliki status open/corrected/ignored/pending_review. |

# 8. Data Requirements

## 8.1 Source of Truth

| Domain Data | Source of Truth | Consumer |
| --- | --- | --- |
| Person/user/role/permission | Core | Semua modul |
| Master data prodi/periode/status | Referensi | PMB, Finance, Academic, HRIS, LMS, Assessment, Portal |
| Lead/peminat | CRM | PMB, Portal, Reporting |
| Applicant | PMB | CRM, Finance, Academic, Assessment, Portal |
| Invoice/payment/clearance | Finance | PMB, Academic, LMS, Portal, Reporting |
| Student/KRS/final grade | Academic | Finance, LMS, Assessment, Portal, Reporting |
| Dosen/pegawai | HRIS | Academic, LMS, Portal |
| Kelas online/progress | LMS | Academic, Portal, Reporting |
| Attempt/scoring | Assessment | PMB, LMS, Academic, Portal |
| Notification/preference | Portal | User dan modul sumber |

## 8.2 External Reference Standard

| Field | Modul Penyimpan | Mengarah ke Source |
| --- | --- | --- |
| person_ref_id | PMB, HRIS, Academic, LMS, Portal | Core |
| user_ref_id | Portal, audit lokal | Core |
| study_program_ref_id | PMB, Academic, HRIS | Referensi |
| academic_period_ref_id | PMB, Finance, Academic, LMS | Referensi |
| lead_ref_id | PMB | CRM |
| applicant_ref_id | Finance, Academic | PMB |
| invoice_ref_id | PMB, Portal | Finance |
| student_ref_id | Finance, LMS, Portal | Academic |
| lecturer_ref_id | Academic, LMS | HRIS |
| academic_class_ref_id | LMS | Academic |
| assessment_session_ref_id | PMB, LMS, Academic | Assessment |

## 8.3 Tabel Teknis Wajib per Database Modul

audit_logs

idempotency_keys

outbox_events

inbox_events

reconciliation_mismatch_logs

Core dapat memiliki tambahan event_contracts, event_consumers, event_replay_logs, integration_event_logs, service_clients, roles, permissions, user_role_assignments, dan active_role_sessions.

# 9. Non-Functional Requirements

| ID | Kategori | Requirement | Acceptance Criteria |
| --- | --- | --- | --- |
| NFR-001 | Availability | Setiap modul memiliki health check aplikasi, database, event publisher, event consumer, dan dependency eksternal. | Status health dapat dilihat DevOps dan admin teknis. |
| NFR-002 | Security | Endpoint protected wajib validasi token, active role, permission, dan data scope. | Direct URL/API access tanpa hak ditolak. |
| NFR-003 | Reliability | Retry integrasi tidak boleh membuat record duplikat. | Idempotency dan unique business key terbukti pada test. |
| NFR-004 | Observability | Setiap request lintas modul membawa correlation_id. | Trace dapat mengikuti flow PMB invoice hingga Finance payment dan Portal notification. |
| NFR-005 | Event Observability | Metric event_lag_seconds, outbox_pending_count, inbox_pending_count, retry_count_by_event, dlq_count_by_consumer, reconciliation_mismatch_count tersedia. | Dashboard observability tersedia. |
| NFR-006 | Performance | Query transaksi hanya memakai database lokal modul. | Tidak ada direct OLTP query ke DB modul lain. |
| NFR-007 | Backup/Restore | Backup per database modul dan restore diuji berkala. | Restore test menghasilkan bukti dan waktu pemulihan. |
| NFR-008 | Accessibility | Main flow mendukung keyboard navigation, label jelas, validation message spesifik, dan responsive layout. | QA aksesibilitas main flow pass. |
| NFR-009 | Auditability | Aksi sensitif, retry, replay, DLQ, dan reconciliation memiliki audit evidence. | Audit log dapat ditelusuri berdasarkan actor/system, timestamp, trace_id/correlation_id. |

# 10. Security, Audit, dan Compliance

## 10.1 Security Requirements

Semua endpoint protected wajib menggunakan Bearer Token dari Core.

Setiap request protected wajib membawa X-Application-Code, X-Active-Role, dan X-Correlation-Id.

Service-to-service call wajib menggunakan service token resmi.

Payment callback wajib validasi provider signature.

Sensitive data pada log, event log, dan export wajib dimasking sesuai permission.

Service credential hanya dapat mengakses database modul sendiri.

Secret management dan rotasi credential wajib tersedia sebelum production.

## 10.2 Audit Requirements

| Aksi | Audit Wajib |
| --- | --- |
| Login/logout dan switch active role | actor, role, session, IP, user agent, timestamp |
| Assign role/permission/scope | old value, new value, actor, reason bila sensitif |
| Verify/reject dokumen PMB | status lama/baru, reason rejection, actor |
| Issue LoA dan handover | applicant_id, student_ref_id jika ada, idempotency_key, correlation_id |
| Payment callback/manual verification | provider_event_id, invoice_ref_id, amount, signature status, actor/system |
| Generate NIM | applicant_ref_id, student_id, nim, actor/system |
| Approve/reject KRS | krs_id, actor dosen PA, status lama/baru, reason |
| Finalize/correct grade | old grade, new grade, reason, actor |
| Replay DLQ dan reconciliation resolve | reason, payload hash, event_key, actor/system |

# 11. State Machine Requirements

Setiap perubahan status harus mengikuti allowed transition, actor yang sah, guard condition, postcondition, audit, dan error code konsisten. Invalid transition harus ditolak dengan STATE_TRANSITION_INVALID atau BUSINESS_RULE_VIOLATION.

## 11.1 Applicant State Machine

DRAFT

- > SUBMITTED

- > DOCUMENT_VERIFIED

- > ACCEPTED

- > LOA_ISSUED

- > HANDED_OVER

## 11.2 Invoice/Payment State Machine

INVOICE_DRAFT

- > ISSUED

- > PARTIALLY_PAID

- > PAID

- > CANCELLED/EXPIRED

PAYMENT_RECEIVED

- > VERIFIED

- > REJECTED/REVERSED

## 11.3 KRS State Machine

DRAFT

- > SUBMITTED

- > APPROVED

- > REJECTED

- > FINALIZED

- > CANCELLED

## 11.4 Assessment Attempt State Machine

NOT_STARTED

- > IN_PROGRESS

- > SUBMITTED

- > SCORED

- > RESULT_PUBLISHED

## 11.5 Event Processing State Machine

PENDING

- > PUBLISHED

- > PROCESSED

- > RETRYING

- > FAILED

- > DLQ

- > REPLAYED

- > IGNORED_DUPLICATE

# 12. Testing dan Acceptance Criteria

## 12.1 Test Level dan Quality Gate

| Test Level | Tujuan | Exit Criteria |
| --- | --- | --- |
| Requirement Review | Memastikan requirement tidak ambigu sebelum development/test. | Tidak ada ambiguitas P0 dan semua P0 punya acceptance criteria. |
| Functional Test | Memastikan fungsi, form, list, action, validasi, output, audit berjalan sesuai FSD. | Seluruh P0 pass dan P1 mayor pass atau workaround disetujui. |
| API Contract Test | Memastikan endpoint, header, payload, response envelope, error envelope, security scheme konsisten. | Endpoint P0 pass, error code konsisten, trace_id/correlation_id tersedia. |
| Integration Test | Memastikan komunikasi lintas modul berjalan dengan service token, event key, retry, dan audit. | Tidak ada duplicate record, event retry aman, integration log tersedia. |
| RBAC/Scope Test | Memastikan role, permission, endpoint, dan data scope ditegakkan backend. | Tidak ada unauthorized read/write dan direct URL/API ditolak. |
| State Machine Test | Memastikan status bergerak sesuai transition, guard, actor, audit, error. | Allowed transition pass, invalid transition ditolak. |
| Migration Validation | Memastikan schema, seed, FK internal, unique index, rollback aman. | Tidak ada blocking migration defect. |
| Regression/Smoke | Memastikan perubahan baru tidak merusak critical path. | Smoke P0 pass sebelum release candidate. |
| Event Contract Test | Memastikan event identity, schema, outbox/inbox, retry, DLQ, reconciliation sesuai kontrak. | Duplicate event aman, retry tidak membuat data ganda, DLQ replay audit-ready. |

## 12.2 Critical UAT Scenario

| ID | Scenario | Expected Result |
| --- | --- | --- |
| UAT-001 | User login dan memilih active role. | Menu, permission, dan scope mengikuti active role. |
| UAT-002 | Admin Prodi A membuka data Prodi B. | Ditolak oleh backend. |
| UAT-003 | Lead dikonversi dua kali menjadi applicant. | Hanya satu applicant terbentuk. |
| UAT-004 | Payment callback provider dikirim dua kali. | Payment tidak dobel dan response idempotent. |
| UAT-005 | Applicant handover dua kali ke Academic. | Student/NIM tidak dobel. |
| UAT-006 | Finance down saat PMB input biodata. | PMB tetap berjalan untuk biodata/dokumen dan status bayar memakai snapshot/degraded label. |
| UAT-007 | Mahasiswa dengan clearance blocked submit KRS final. | KRS ditolak/pending sesuai policy. |
| UAT-008 | Academic publish krs_approved dua kali. | LMS enrollment tidak dobel. |
| UAT-009 | LMS mengirim grade input duplicate. | Academic tidak membuat grade input ganda. |
| UAT-010 | Event consumer down lalu pulih. | Event retry dan diproses setelah consumer pulih. |
| UAT-011 | DLQ replay dilakukan. | Replay audit-ready dan tidak membuat duplicate record. |
| UAT-012 | Snapshot stale. | UI menampilkan synced_at/refreshed_at dan label freshness. |

## 12.3 Go/No-Go Criteria

Seluruh P0 UAT dan SIT pass.

Tidak ada defect Sev-1 dan Sev-2 terbuka.

RBAC/scope test pass.

Migration dry-run pass dan rollback rehearsal aman.

Smoke test pass sebelum release candidate.

Sign-off Product Owner, QA/UAT Lead, Technical Lead, DBA, dan Owner Modul terkait lengkap.

Release ditahan jika ada duplikasi payment, duplikasi student/NIM, kebocoran data lintas scope, bypass state machine, migration failure tanpa rollback aman, atau proses akademik/keuangan P0 tidak berjalan.

# 13. Deployment, Operasional, dan Release

## 13.1 Deployment Environment

| Environment | Tujuan |
| --- | --- |
| Local Development | Development per repo/modul menggunakan Docker Compose lokal. |
| QA | Functional test, API contract test, RBAC/scope test, state machine test. |
| Staging | SIT, UAT, migration rehearsal, smoke, integration event test. |
| Production | Pilot/go-live bertahap dengan monitoring dan rollback plan. |

## 13.2 Release Bertahap

| Release | Scope |
| --- | --- |
| Release 0 - Foundation | Infra, shared contracts, Core Auth/RBAC, Reference master, audit, idempotency, outbox/inbox baseline. |
| Release 1 - PMB & Finance MVP | Applicant, document, request invoice, payment callback/manual verification, clearance, LoA. |
| Release 2 - Academic MVP | Handover, generate student/NIM, curriculum, class, KRS, KHS basic. |
| Release 3 - LMS & Assessment | Class sync, enrollment, LMS activity, grade input, assessment session/attempt/scoring/result. |
| Release 4 - Portal & Dashboard | Portal multi-role, notification center, executive dashboard, activity log. |
| Release 5 - Hardening & Reporting | Reconciliation, warehouse/read model, observability, backup/restore rehearsal, performance tuning. |

## 13.3 Operational Requirements

Setiap service memiliki health endpoint.

Setiap database memiliki backup dan restore runbook.

Setiap event consumer memiliki metric lag, retry count, DLQ count.

Setiap incident integration memiliki correlation_id untuk tracing.

Secret dan service token harus dirotasi sesuai policy.

Rollback plan tersedia per service dan per database migration.

# 14. Requirement Traceability Matrix

| Requirement Area | Source Basis | Test Basis |
| --- | --- | --- |
| Module boundary dan database per modul | PRD, BRD, DBML | Migration validation, integration test, degraded mode test |
| Role, permission, dan data scope | FSD, RBAC Matrix, OpenAPI | RBAC/scope test, direct API access test |
| API contract | OpenAPI, FSD | API contract test, integration test |
| Event contract dan outbox/inbox | PRD, BRD, FSD Appendix Event, UAT Appendix | Event contract test, duplicate event test, retry/DLQ test |
| State machine | FSD, UAT | State machine test dan negative transition test |
| Audit dan idempotency | BRD, FSD, OpenAPI, UAT | Audit evidence test, duplicate request test |
| Degraded mode dan reconciliation | PRD, BRD, FSD, UAT | Partial outage test dan mismatch reconciliation test |
| Release readiness | UAT/QA Test Plan | Go/no-go checklist, smoke test, sign-off |

# 15. Lampiran

## 15.1 Endpoint P0

| Modul | Endpoint/Fungsi P0 |
| --- | --- |
| Core | Login, refresh token, auth me, switch role, application launcher, role/permission management. |
| Referensi | Study programs, academic years, academic periods, status code, payment components. |
| CRM | Lead create/read/update, follow-up, convert lead to applicant. |
| PMB | Applicant create/read/update, submit, document upload/verify/reject, request invoice, issue LoA, handover. |
| Finance | Invoice create/read, payment callback, manual verification, clearance read/update. |
| Academic | Generate student/NIM, student read, class create/read, KRS submit/approve/reject, grade import/finalize. |
| HRIS | Lecturer active read, lecturer status change. |
| LMS | Class sync, enrollment sync, session/material/task/attendance, grade sync. |
| Assessment | Session create, attempt start/submit, scoring, result publish. |
| Portal | Dashboard role-based, notification read/mark-read, preference/shortcut. |
| Integration | Event catalog, outbox monitoring, inbox monitoring, DLQ replay, reconciliation mismatch. |

## 15.2 Error Code Minimum

| Error Code | HTTP | Makna |
| --- | --- | --- |
| AUTH_REQUIRED | 401 | Token tidak ditemukan. |
| TOKEN_EXPIRED | 401 | Token expired. |
| TOKEN_INVALID | 401 | Token tidak valid. |
| ROLE_NOT_ASSIGNED | 403 | User tidak memiliki active role. |
| PERMISSION_DENIED | 403 | Permission tidak tersedia. |
| FORBIDDEN_SCOPE | 403 | Data di luar scope user. |
| VALIDATION_ERROR | 422 | Validasi field gagal. |
| RESOURCE_NOT_FOUND | 404 | Data tidak ditemukan. |
| STATE_TRANSITION_INVALID | 409 | Perubahan status tidak diizinkan. |
| DUPLICATE_REQUEST | 409 | Request duplikat. |
| IDEMPOTENCY_KEY_REQUIRED | 422 | Idempotency-Key wajib tetapi tidak dikirim. |
| IDEMPOTENCY_PAYLOAD_MISMATCH | 409 | Key sama tetapi payload berbeda. |
| BUSINESS_RULE_VIOLATION | 409 | Melanggar aturan bisnis. |
| INTEGRATION_FAILED | 502 | Service dependency gagal. |
| DEPENDENCY_UNAVAILABLE | 503 | Modul dependency down. |
| PROVIDER_SIGNATURE_INVALID | 401 | Signature payment callback invalid. |
| EVENT_SCHEMA_INVALID | 422 | Payload event tidak sesuai schema. |
| RECONCILIATION_REQUIRED | 409 | Snapshot berbeda dari source. |

## 15.3 Definition of Done Requirement

Requirement memiliki ID, prioritas, deskripsi, owner, dan acceptance criteria.

Requirement dapat ditelusuri ke modul, API/event, tabel/domain, role/scope, dan test case.

Endpoint terkait memiliki API contract dan error code.

Proses kritis memiliki idempotency rule dan unique business key.

Aksi sensitif memiliki audit rule.

Event terkait memiliki event contract, outbox/inbox, retry, DLQ, dan reconciliation rule jika perlu.

QA memiliki positive, negative, RBAC/scope, idempotency, audit, dan regression scenario.
