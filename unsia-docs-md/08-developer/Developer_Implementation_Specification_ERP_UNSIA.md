---
title: "Developer Implementation Specification ERP UNSIA"
source_file: "Developer_Implementation_Specification_ERP_UNSIA.docx"
format: markdown
---

# Developer Implementation Specification ERP UNSIA

UNSIA

DEVELOPER IMPLEMENTATION SPECIFICATION

ERP Pendidikan / SIAKAD Terintegrasi UNSIA

Versi Dokumen: v1.0 | Tanggal: 22 Juni 2026

Disusun untuk kebutuhan Backend Developer, Frontend Developer, DBA, QA/UAT, DevOps, dan Technical Lead

| Item | Isi |
| --- | --- |
| Dokumen | Developer Implementation Specification |
| Produk | ERP Pendidikan / SIAKAD Terintegrasi UNSIA |
| Basis | PRD Global, BRD, FSD Per Modul, DBML, OpenAPI/Swagger, UAT Scenario, Event Contract |
| Tujuan | Menurunkan blueprint system analyst menjadi spesifikasi teknis yang bisa langsung dijadikan backlog dan task development. |
| Pendekatan | Modular-first, API-first, event-driven, RBAC-enforced, audit-ready, idempotent, dan UAT-driven. |

# Daftar Isi

### 1. Prinsip Implementasi Wajib

### 2. Struktur Service dan Repository

### 3. Database Boundary

### 4. Standar Backend Middleware

### 5. Standar Response API

### 6. Shared Library Wajib

### 7. Spesifikasi Modul Core

### 8. Spesifikasi Modul Reference

### 9. Spesifikasi Modul CRM

### 10. Spesifikasi Modul PMB

### 11. Spesifikasi Modul Finance

### 12. Spesifikasi Modul Academic

### 13. Spesifikasi Modul HRIS

### 14. Spesifikasi Modul LMS

### 15. Spesifikasi Modul Assessment

### 16. Spesifikasi Modul Portal

### 17. Event Contract Implementation

### 18. Event Catalog Minimum

### 19. Degraded Mode

### 20. Frontend Implementation Standard

### 21. Critical Flow untuk Developer

### 22. RBAC dan Data Scope

### 23. Testing Wajib Developer

### 24. Backlog Developer per Phase

### 25. Definition of Done Developer

### 26. Catatan Senior System Analyst

# 1. Prinsip Implementasi Wajib

ERP ini harus dibangun sebagai modular distributed ERP. Setiap modul utama memiliki database sendiri, tidak boleh ada foreign key lintas database, dan transaksi online tidak boleh melakukan direct join antar database. Relasi lintas modul wajib menggunakan external_ref_id, API, event, snapshot, read model, dan reconciliation.

| Area | Aturan Teknis |
| --- | --- |
| Database | 1 modul = 1 database fisik. |
| Foreign Key | FK hanya boleh dalam database modul yang sama. |
| Cross-module write | Wajib melalui API command atau event resmi. |
| Cross-module read | Melalui API query, snapshot, read model, atau dashboard projection. |
| Auth | Semua modul memakai Core Auth/JWT/active role. |
| RBAC | Validasi role, permission, dan data scope wajib dilakukan di backend. |
| Critical process | Wajib idempotent. |
| Event | Wajib outbox/inbox. |
| Audit | Semua aksi sensitif wajib audit. |
| Degraded mode | Modul tetap berjalan terbatas saat dependency down. |
| Reconciliation | Snapshot/read model harus bisa dicek ulang ke source of truth. |

# 2. Struktur Service dan Repository

Rekomendasi struktur repository dibuat agar service boundary terlihat jelas dan shared package tidak tersebar di setiap modul.

| unsia-erp/<br>├── services/<br>│ ├── core-service/<br>│ ├── reference-service/<br>│ ├── crm-service/<br>│ ├── pmb-service/<br>│ ├── finance-service/<br>│ ├── academic-service/<br>│ ├── hris-service/<br>│ ├── lms-service/<br>│ ├── assessment-service/<br>│ └── portal-service/<br>├── packages/<br>│ ├── shared-auth/<br>│ ├── shared-rbac/<br>│ ├── shared-idempotency/<br>│ ├── shared-audit/<br>│ ├── shared-event/<br>│ ├── shared-http-client/<br>│ ├── shared-error-envelope/<br>│ └── shared-observability/<br>├── infra/<br>│ ├── docker/<br>│ ├── k8s/<br>│ ├── migrations/<br>│ ├── monitoring/<br>│ └── ci-cd/<br>└── docs/<br> ├── openapi/<br> ├── event-contract/<br> ├── dbml/<br> ├── rbac-matrix/<br> ├── uat/<br> └── runbook/ |
| --- |

Setiap service hanya boleh mengakses database miliknya sendiri.

Jika modul membutuhkan data modul lain, gunakan API query, event projection, atau snapshot lokal.

Shared package harus dibuat sejak awal agar standar auth, RBAC, audit, idempotency, error, dan event seragam.

# 3. Database Boundary

## 3.1 Database per Modul

| Service | Database | Ownership |
| --- | --- | --- |
| Core Service | core_db | Person, user, role, permission, session, service token. |
| Reference Service | reference_db | Prodi, tahun ajaran, periode akademik, status code, payment component. |
| CRM Service | crm_db | Campaign, lead, agent, referral. |
| PMB Service | pmb_db | Applicant, biodata, dokumen, LoA, handover. |
| Finance Service | finance_db | Invoice, payment, clearance, receipt. |
| Academic Service | academic_db | Student, NIM, KRS, kelas, nilai final, KHS, transkrip. |
| HRIS Service | hris_db | Employee, lecturer, homebase, status dosen. |
| LMS Service | lms_db | Online class, enrollment, session, materi, tugas, presensi, grade input. |
| Assessment Service | assessment_db | Bank soal, session, attempt, answer, score. |
| Portal Service | portal_db | Notification, dashboard, shortcut, preference. |

## 3.2 Tabel Teknis Wajib di Setiap Modul

| audit_logs<br>idempotency_keys<br>outbox_events<br>inbox_events<br>reconciliation_mismatch_logs |
| --- |

Tambahan khusus pada core_db:

| event_contracts<br>event_consumers<br>event_replay_logs<br>integration_event_logs |
| --- |

# 4. Standar Backend Middleware

Setiap endpoint protected wajib melewati middleware berikut secara konsisten.

| Request<br> ↓<br>CorrelationIdMiddleware<br> ↓<br>AuthMiddleware<br> ↓<br>ApplicationCodeMiddleware<br> ↓<br>ActiveRoleMiddleware<br> ↓<br>PermissionMiddleware<br> ↓<br>DataScopeMiddleware<br> ↓<br>IdempotencyMiddleware untuk command kritis<br> ↓<br>Controller<br> ↓<br>Service Layer<br> ↓<br>Repository<br> ↓<br>Audit + Outbox Event<br> ↓<br>Response Envelope |
| --- |

## 4.1 Header Wajib

| Authorization: Bearer <jwt_token><br>X-Application-Code: PMB<br>X-Active-Role: admin_pmb<br>X-Correlation-Id: 3fb4b7f1-7d28-4d13-a812-9cc5e1c0c011<br>Idempotency-Key: pmb:handover:applicant-uuid |
| --- |

Authorization wajib untuk seluruh endpoint protected.

X-Application-Code digunakan untuk memastikan request datang dari aplikasi/modul yang diizinkan.

X-Active-Role menentukan konteks role yang sedang aktif.

X-Correlation-Id wajib diteruskan ke log, audit, event, dan service call lintas modul.

Idempotency-Key wajib untuk command kritis seperti payment callback, handover, generate NIM, class sync, enrollment sync, grade sync, dan replay DLQ.

# 5. Standar Response API

## 5.1 Success Envelope

| {<br> "success": true,<br> "message": "Request processed successfully",<br> "data": {},<br> "meta": {<br> "trace_id": "3fb4b7f1-7d28-4d13-a812-9cc5e1c0c011",<br> "timestamp": "2026-06-22T10:00:00+07:00"<br> }<br>} |
| --- |

## 5.2 Error Envelope

| {<br> "success": false,<br> "error": {<br> "code": "FORBIDDEN_SCOPE",<br> "message": "Anda tidak memiliki akses ke data ini.",<br> "details": {}<br> },<br> "meta": {<br> "trace_id": "3fb4b7f1-7d28-4d13-a812-9cc5e1c0c011",<br> "timestamp": "2026-06-22T10:00:00+07:00"<br> }<br>} |
| --- |

# 6. Shared Library Wajib

| Package | Fungsi |
| --- | --- |
| shared-auth | JWT validation, JWKS cache, service token validation. |
| shared-rbac | Permission check dan data scope resolver. |
| shared-idempotency | Simpan request hash, response cache, lock, expiry. |
| shared-audit | Audit actor, role, old value, new value, reason. |
| shared-event | Outbox writer, inbox consumer, retry, DLQ, event envelope. |
| shared-http-client | Service-to-service call, timeout, retry, circuit breaker. |
| shared-error-envelope | Format error konsisten. |
| shared-observability | Trace ID, correlation ID, structured log, metrics. |

# 7. Spesifikasi Modul Core

Core adalah pusat identitas, login, SSO, user, role, permission, active role, application launcher, service token, impersonation, audit global, dan event contract catalog.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | /api/v1/auth/login | Login user |
| POST | /api/v1/auth/refresh | Refresh token |
| GET | /api/v1/auth/me | Profil user, active role, permission, dan scope |
| POST | /api/v1/auth/switch-role | Ganti active role |
| GET | /api/v1/applications | Application launcher |
| POST | /api/v1/impersonations/start | Impersonation |

## Rule Developer

Login hanya di Core.

Modul lain validate JWT via cached JWKS/public key.

Active role harus terbaca dari token/session.

Permission harus granular: module.resource.action.

Data scope harus ikut active role.

Impersonation wajib reason dan audit.

Service-to-service call wajib service token.

## Acceptance Criteria

User multi-role bisa login dan memilih role aktif.

Menu berubah mengikuti active role.

Admin Prodi tidak bisa akses data prodi lain.

Service token invalid ditolak.

Core down: token valid tetap bisa diverifikasi selama JWKS cache masih valid.

# 8. Spesifikasi Modul Reference

Reference Service menjadi source of truth untuk master data lintas modul: program studi, tahun ajaran, periode akademik, jalur PMB, gelombang PMB, document type, payment component, payment method, status code, region, dan religion.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| GET | /api/v1/ref/study-programs | List prodi |
| GET | /api/v1/ref/academic-years | List tahun ajaran |
| POST | /api/v1/ref/academic-years | Buat tahun ajaran |
| GET | /api/v1/ref/academic-periods | List periode akademik |
| POST | /api/v1/ref/academic-periods | Buat periode akademik |

## Rule Developer

Master data yang sudah dipakai transaksi tidak boleh hard delete.

Gunakan inactive atau archived.

Status bisnis kritis harus menggunakan managed status code.

Perubahan master data harus publish event.

Consumer menyimpan snapshot lokal.

## Acceptance Criteria

Admin dapat membuat tahun ajaran dan periode akademik.

Hanya satu periode aktif sesuai policy.

Master data yang sudah dipakai tidak bisa dihapus permanen.

Event update master data terkirim ke consumer.

# 9. Spesifikasi Modul CRM

CRM mengelola lead, campaign, agent, referral, follow-up, conversion, dan commission.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| GET | /api/v1/crm/leads | List lead sesuai scope |
| POST | /api/v1/crm/leads | Buat lead |
| POST | /api/v1/crm/leads/{lead_id}/convert-to-applicant | Convert lead ke applicant PMB |

## Rule Developer

Lead hanya bisa dikonversi jika status qualified.

Convert lead wajib idempotent.

Agent hanya bisa melihat lead/referral miliknya sendiri.

CRM menyimpan applicant_ref_id setelah berhasil convert.

## Acceptance Criteria

Duplicate convert tidak membuat applicant ganda.

applicant_ref_id tersimpan di CRM.

Outbox CRM dan PMB tercatat.

Akses lintas agent ditolak.

# 10. Spesifikasi Modul PMB

PMB menjadi source of truth applicant, biodata, dokumen, seleksi, LoA, daftar ulang, dan handover ke Academic.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | /api/v1/pmb/applicants | Membuat applicant |
| POST | /api/v1/pmb/applicants/{id}/submit | Submit pendaftaran |
| POST | /api/v1/pmb/applicants/{id}/documents | Upload dokumen |
| POST | /api/v1/pmb/applicants/{id}/documents/{document_id}/verify | Verifikasi dokumen |
| POST | /api/v1/pmb/applicants/{id}/request-invoice | Request invoice ke Finance |
| POST | /api/v1/pmb/applicants/{id}/issue-loa | Terbitkan LoA |
| POST | /api/v1/pmb/applicants/{id}/handover-to-academic | Handover ke Academic |

## Rule Developer

Applicant tidak boleh dibuat ganda dari lead yang sama.

Submit hanya bisa jika biodata minimum lengkap.

Dokumen rejected wajib memiliki rejection_reason.

LoA hanya bisa terbit jika dokumen dan payment policy valid.

Handover wajib idempotent.

Jika Academic down, handover masuk retry queue atau pending.

PMB hanya menyimpan payment status sebagai read model, bukan source of truth.

## Acceptance Criteria

Applicant lifecycle berjalan dari draft sampai handed_over.

Duplicate request invoice mengembalikan invoice existing.

Duplicate handover tidak membuat student/NIM ganda.

PMB menampilkan status pembayaran dari snapshot Finance.

# 11. Spesifikasi Modul Finance

Finance adalah source of truth untuk invoice, invoice item, payment, callback, manual verification, receipt, clearance, scholarship, dan laporan keuangan.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | /api/v1/finance/invoices | Membuat invoice |
| GET | /api/v1/finance/invoices/{invoice_id} | Detail invoice |
| POST | /api/v1/finance/payment-callbacks/{provider} | Callback payment gateway |
| POST | /api/v1/finance/payment-verifications | Verifikasi pembayaran manual |
| GET | /api/v1/finance/clearances | Cek clearance applicant/mahasiswa |

## Rule Developer

Callback wajib validasi signature provider.

provider_event_id harus unique.

Duplicate callback tidak boleh membuat payment ganda.

Payment paid harus publish event finance.payment_paid.

Clearance adalah source of truth Finance.

## Acceptance Criteria

Invoice bisa dibuat dari PMB/Academic request.

Payment callback duplicate aman.

Clearance berubah sesuai payment policy.

PMB/Academic/Portal menerima event payment dan clearance.

# 12. Spesifikasi Modul Academic

Academic adalah source of truth mahasiswa, NIM, kurikulum, mata kuliah, kelas, KRS, nilai final, KHS, transkrip, yudisium, dan alumni.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | /api/v1/academic/students/generate-from-applicant | Generate student dan NIM |
| GET | /api/v1/academic/students | List mahasiswa |
| POST | /api/v1/academic/classes | Buka kelas |
| POST | /api/v1/academic/krs | Buat draft KRS |
| POST | /api/v1/academic/krs/{krs_id}/submit | Submit KRS |
| POST | /api/v1/academic/krs/{krs_id}/approve | Approval KRS |
| POST | /api/v1/academic/grades/source-imports | Import source grade |
| POST | /api/v1/academic/grades/{grade_id}/finalize | Finalisasi nilai |

## Rule Developer

applicant_ref_id pada students harus unique.

NIM harus unique dan sequence wajib dikunci saat generate.

Mahasiswa hanya bisa submit KRS miliknya sendiri.

Dosen PA hanya approve mahasiswa bimbingannya.

KRS final wajib cek clearance Finance.

Final grade hanya milik Academic.

## Acceptance Criteria

Duplicate handover tidak membuat student/NIM ganda.

KRS blocked jika clearance tidak memenuhi policy.

Duplicate krs_approved tidak membuat LMS enrollment ganda.

LMS/Assessment hanya menjadi sumber grade input.

# 13. Spesifikasi Modul HRIS

HRIS adalah source of truth dosen, pegawai, homebase, unit kerja, jabatan, status aktif, BKD, dan data SDM.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| GET | /api/v1/hris/lecturers | List dosen aktif |
| GET | /api/v1/hris/lecturers/{lecturer_id} | Detail dosen aktif |

## Rule Developer

Academic dan LMS tidak boleh membuat data dosen mandiri.

Academic/LMS hanya menyimpan lecturer_ref_id dan snapshot.

Dosen nonaktif tidak boleh diplot ke kelas baru.

Perubahan dosen harus publish event hris.lecturer_updated.

## Acceptance Criteria

Dosen aktif dapat dibaca Academic/LMS.

Dosen nonaktif ditolak saat plotting kelas.

Snapshot lecturer di consumer ter-update via event.

# 14. Spesifikasi Modul LMS

LMS mengelola online class, enrollment, sesi, materi, video, vicon, assignment, submission, discussion, attendance, progress, dan grade input.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | /api/v1/lms/classes/sync-from-academic | Sync kelas dari Academic |
| POST | /api/v1/lms/enrollments/sync-from-krs | Sync enrollment dari KRS |
| POST | /api/v1/lms/grade-syncs | Kirim nilai aktivitas ke Academic |

## Rule Developer

LMS tidak membuat kelas akademik sendiri.

LMS class berasal dari academic.class_opened.

LMS enrollment berasal dari academic.krs_approved.

LMS boleh menyimpan progress, presensi, tugas, dan grade input.

Final grade tetap milik Academic.

## Acceptance Criteria

Class sync duplicate aman.

Enrollment sync duplicate aman.

Grade input terkirim ke Academic tanpa menjadi final grade otomatis.

# 15. Spesifikasi Modul Assessment

Assessment adalah mesin CBT, quiz, survey, question bank, attempt, answer, scoring, dan result publish.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | /api/v1/assessment/sessions | Membuat assessment session |
| POST | /api/v1/assessment/attempts | Membuat attempt peserta |
| POST | /api/v1/assessment/results/publish | Publish hasil ke consumer |

## Rule Developer

Question yang sudah dipakai attempt tidak boleh diedit langsung.

Jika butuh perubahan, buat question_version baru.

Attempt harus immutable setelah submitted.

Result publish wajib idempotent.

Result hanya menjadi input untuk PMB/LMS/Academic, bukan otomatis final decision.

## Acceptance Criteria

Question versioning berjalan.

Attempt submitted tidak dapat diubah.

Duplicate result publish aman.

Consumer menerima result sebagai input sesuai konteks.

# 16. Spesifikasi Modul Portal

Portal adalah presentation layer, bukan source transaksi. Portal mengelola dashboard, notification center, read marker, user preference, shortcut, dan executive dashboard.

## Endpoint Utama

| Method | Endpoint | Fungsi |
| --- | --- | --- |
| POST | /api/v1/portal/notifications | Membuat notifikasi portal |
| GET | /api/v1/portal/dashboard | Dashboard sesuai role |

## Rule Developer

Portal tidak boleh menjadi source status bisnis.

Portal hanya membaca read model atau snapshot.

Dashboard wajib menampilkan refreshed_at.

Jika Portal down, transaksi PMB/Finance/Academic tetap sukses.

Notification event diproses ulang setelah Portal pulih.

## Acceptance Criteria

Dashboard sesuai active role.

Data dashboard punya refreshed_at.

Portal down tidak mengganggu transaksi modul sumber.

Read marker bekerja per user.

# 17. Event Contract Implementation

## 17.1 Event Envelope Standar

| {<br> "event_name": "finance.payment_paid",<br> "event_version": "v1",<br> "event_key": "finance.payment_paid:payment-uuid:v1",<br> "publisher_service": "finance-service",<br> "aggregate_type": "payment",<br> "aggregate_id": "payment-uuid",<br> "correlation_id": "3fb4b7f1-7d28-4d13-a812-9cc5e1c0c011",<br> "causation_id": "callback-provider-event-id",<br> "occurred_at": "2026-06-22T10:00:00+07:00",<br> "payload": {}<br>} |
| --- |

## 17.2 Outbox Flow

| Business Transaction<br> ↓<br>Update domain table<br> ↓<br>Insert audit_logs<br> ↓<br>Insert outbox_events status = PENDING<br> ↓<br>Commit transaction<br> ↓<br>Outbox worker publish to event broker<br> ↓<br>Update outbox_events status = PUBLISHED |
| --- |

## 17.3 Inbox Flow

| Receive event<br> ↓<br>Check inbox_events by event_key<br> ↓<br>If exists → mark IGNORED_DUPLICATE<br> ↓<br>If not exists → insert inbox event<br> ↓<br>Validate event version + payload schema<br> ↓<br>Update snapshot/read model<br> ↓<br>Mark inbox_events PROCESSED |
| --- |

# 18. Event Catalog Minimum

| Event | Publisher | Consumer |
| --- | --- | --- |
| core.person_updated | Core | CRM, PMB, HRIS, Academic, Portal |
| reference.study_program_updated | Referensi | PMB, Academic, HRIS, LMS, Portal |
| reference.academic_period_updated | Referensi | PMB, Finance, Academic, LMS, Assessment, Portal |
| crm.lead_qualified | CRM | PMB |
| pmb.applicant_created | PMB | Finance, Assessment, Portal |
| finance.invoice_created | Finance | PMB, Portal |
| finance.payment_paid | Finance | PMB, Academic, Portal |
| finance.clearance_changed | Finance | PMB, Academic, LMS, Portal |
| pmb.ready_for_academic | PMB | Academic |
| academic.student_created | Academic | PMB, Finance, LMS, Portal |
| academic.class_opened | Academic | LMS, Portal |
| academic.krs_approved | Academic | LMS, Finance, Portal |
| lms.grade_input_submitted | LMS | Academic |
| assessment.result_calculated | Assessment | PMB, LMS, Academic, Portal |

# 19. Degraded Mode

| Dependency Down | Perilaku Sistem |
| --- | --- |
| Finance down | PMB tetap menerima biodata/dokumen; payment status memakai snapshot terakhir. |
| Academic down | LMS tetap membuka kelas/enrollment yang sudah tersinkron. |
| LMS down | Academic tetap membuka kelas, KRS, dan final grade. |
| Portal down | Modul sumber tetap transaksi; notifikasi pending. |
| Core down | Token valid tetap diverifikasi via cached JWKS; login baru ditahan. |

UI wajib menampilkan synced_at/refreshed_at jika memakai snapshot atau read model.

Aksi yang berisiko harus dibatasi atau diberi status pending_review.

Retry queue dan last_error harus bisa dilihat oleh role teknis yang berwenang.

# 20. Frontend Implementation Standard

| Page | Requirement |
| --- | --- |
| List page | Pagination, search, sorting, filter, empty state, loading state. |
| Form page | Mandatory marker, validation message, confirmation untuk aksi sensitif. |
| Detail page | Status, histori, related records, action sesuai permission. |
| Snapshot page | Tampilkan synced_at, source_module, dan freshness label. |
| Error state | User message jelas, trace ID tersedia. |
| Degraded state | Label data tidak real-time, pending_review, atau retrying. |
| Integration log | Tampilkan event_key, idempotency_key, retry_count, last_error. |

# 21. Critical Flow untuk Developer

## 21.1 Lead to Applicant

| CRM lead qualified<br> ↓<br>POST /api/v1/crm/leads/{lead_id}/convert-to-applicant<br> ↓<br>CRM validate lead status<br> ↓<br>CRM call PMB create applicant<br> ↓<br>PMB create applicant<br> ↓<br>PMB publish pmb.applicant_created<br> ↓<br>Finance/Assessment/Portal consume event |
| --- |

Acceptance:

Duplicate convert tidak membuat applicant ganda.

applicant_ref_id tersimpan di CRM.

Outbox CRM dan PMB tercatat.

## 21.2 Applicant to Invoice

| Admin PMB request invoice<br> ↓<br>POST /api/v1/pmb/applicants/{id}/request-invoice<br> ↓<br>PMB validate applicant<br> ↓<br>PMB call Finance POST /api/v1/finance/invoices<br> ↓<br>Finance create invoice<br> ↓<br>Finance publish finance.invoice_created<br> ↓<br>PMB update applicant_invoice_statuses |
| --- |

Acceptance:

Duplicate request invoice return existing invoice.

PMB tidak insert langsung ke finance_db.

Invoice source of truth tetap Finance.

## 21.3 Payment Callback

| Provider callback<br> ↓<br>POST /api/v1/finance/payment-callbacks/{provider}<br> ↓<br>Validate signature<br> ↓<br>Check provider_event_id<br> ↓<br>Create payment<br> ↓<br>Update invoice status<br> ↓<br>Publish finance.payment_paid<br> ↓<br>PMB/Academic/Portal update snapshot |
| --- |

Acceptance:

Duplicate callback tidak membuat payment ganda.

PMB/Academic snapshot tidak dobel.

Payment event memiliki event_key deterministik.

## 21.4 Handover Applicant to Academic

| PMB applicant accepted + LoA issued + payment policy valid<br> ↓<br>POST /api/v1/pmb/applicants/{id}/handover-to-academic<br> ↓<br>PMB call Academic generate student<br> ↓<br>Academic validate applicant_ref_id unique<br> ↓<br>Academic generate NIM<br> ↓<br>Academic publish academic.student_created<br> ↓<br>PMB update student_ref_id |
| --- |

Acceptance:

Duplicate handover tidak membuat student/NIM ganda.

Jika Academic down, handover masuk retry.

Academic.students.applicant_ref_id harus unique.

## 21.5 KRS to LMS Enrollment

| Mahasiswa submit KRS<br> ↓<br>Dosen PA approve<br> ↓<br>Academic check Finance clearance<br> ↓<br>Academic publish academic.krs_approved<br> ↓<br>LMS consume event<br> ↓<br>LMS upsert enrollment |
| --- |

Acceptance:

Clearance blocked membuat KRS final ditolak.

Duplicate krs_approved tidak membuat enrollment ganda.

LMS class/enrollment berasal dari Academic.

# 22. RBAC dan Data Scope

| Role | Data Scope |
| --- | --- |
| Super Admin | Global |
| Admin PMB | PMB domain |
| Admin Finance | Finance domain |
| Admin Akademik Biro | Academic global |
| Kaprodi/Admin Prodi | study_program_id |
| Dosen | Assigned class |
| Dosen PA | Advisor scope |
| Mahasiswa | Self scope |
| Agent/Mitra | Own referral/lead |
| Pimpinan | Read-only aggregate |

## 22.1 Pseudocode Scope Check

| function authorize(request, permission):<br> user = validateJwt(request.Authorization)<br> activeRole = request.header["X-Active-Role"]<br> appCode = request.header["X-Application-Code"]<br><br> assert user.hasActiveRole(activeRole)<br> assert user.hasPermission(activeRole, permission)<br> assert appCode in user.allowedApplications<br><br> scope = resolveScope(user, activeRole)<br><br> if scope.type == "global":<br> return allow<br><br> if scope.type == "study_program":<br> assert request.resource.study_program_id == scope.study_program_id<br><br> if scope.type == "self":<br> assert request.resource.person_ref_id == user.person_id<br><br> if scope.type == "assigned_class":<br> assert request.resource.class_id in user.assigned_classes<br><br> return allow |
| --- |

# 23. Testing Wajib Developer

## 23.1 Unit Test

Validation rule.

State transition.

Duplicate check.

Idempotency service.

Event payload builder.

Permission resolver.

## 23.2 Integration Test

CRM ke PMB convert lead.

PMB ke Finance request invoice.

Finance ke PMB payment status.

PMB ke Academic handover.

Academic ke LMS class sync.

LMS ke Academic grade input.

Assessment ke PMB/LMS/Academic result.

Modul ke Portal notification.

## 23.3 Event Contract Test

Outbox dibuat setelah transaksi commit.

Inbox mencatat event_key.

Duplicate event menjadi IGNORED_DUPLICATE.

Retry mengisi retry_count dan next_retry_at.

Retry limit exceeded masuk DLQ.

Replay DLQ butuh reason dan audit.

Reconciliation mismatch menghasilkan report.

# 24. Backlog Developer per Phase

## Phase 0 - Architecture Foundation

| Task | Output/Owner |
| --- | --- |
| Setup repo/service structure | Tech Lead |
| Setup database per modul | DBA |
| Setup migration pipeline per database | DBA/DevOps |
| Setup API gateway | Backend/DevOps |
| Setup JWT/JWKS validation | Backend |
| Setup shared response envelope | Backend |
| Setup shared idempotency package | Backend |
| Setup shared audit package | Backend |
| Setup outbox/inbox worker | Backend/DevOps |
| Setup event broker | DevOps |
| Setup logging, metrics, tracing | DevOps |

## Phase 1 - Core + Reference

| Task | Output/Owner |
| --- | --- |
| Login | POST /auth/login |
| Refresh token | POST /auth/refresh |
| Auth me | GET /auth/me |
| Switch role | POST /auth/switch-role |
| Application launcher | GET /applications |
| Role permission seed | RBAC matrix |
| Master prodi | /ref/study-programs |
| Tahun ajaran | /ref/academic-years |
| Periode akademik | /ref/academic-periods |

## Phase 2 - CRM + PMB + Finance

| Task | Output/Owner |
| --- | --- |
| Lead capture | Lead tercatat |
| Convert lead | Applicant PMB |
| Applicant registration | Applicant + biodata |
| Upload document | Applicant document |
| Verify document | Verified/rejected |
| Request invoice | Invoice Finance |
| Payment callback | Payment paid |
| Applicant invoice read model | PMB payment status |
| Issue LoA | LoA document |
| Handover | Request ke Academic |

## Phase 3 - Academic Core

| Task | Output/Owner |
| --- | --- |
| Generate student | Student + NIM |
| Student list scoped | Admin/Kaprodi |
| Curriculum | Kurikulum |
| Course offering | Kelas dibuka |
| KRS draft | KRS header/item |
| KRS approval | Approved KRS |
| Clearance check | clear/blocked/conditional |
| Publish KRS approved | LMS enrollment event |

## Phase 4 - HRIS + LMS

| Task | Output/Owner |
| --- | --- |
| Lecturer read model | Dosen aktif |
| Class sync | LMS class |
| Enrollment sync | LMS enrollment |
| LMS session | Session pembelajaran |
| Material/task/attendance | Aktivitas LMS |
| Grade sync | Source grade ke Academic |

## Phase 5 - Assessment + Grade

| Task | Output/Owner |
| --- | --- |
| Question bank | Bank soal |
| Question version | Versioning soal |
| Assessment session | CBT/quiz/survey session |
| Attempt | Jawaban peserta |
| Scoring | Score |
| Result publish | Event result |
| Grade source import | Academic grade input |

## Phase 6 - Portal + Reporting

| Task | Output/Owner |
| --- | --- |
| Notification center | Notification/read marker |
| Dashboard role-based | Dashboard per role |
| Executive dashboard | KPI agregat |
| Activity log | Jejak aktivitas |
| Reconciliation monitor | Mismatch report |
| Integration log viewer | Event tracing |

# 25. Definition of Done Developer

Endpoint sesuai OpenAPI.

Request/response memakai envelope standar.

Auth, active role, permission, dan data scope sudah diuji.

Command kritis memakai Idempotency-Key.

Audit log tercatat untuk aksi sensitif.

Status transition sesuai state machine.

Event outbox dibuat jika ada perubahan penting.

Consumer inbox idempotent.

Error code konsisten.

Unit test pass.

Integration test pass.

RBAC/scope test pass.

Evidence UAT tersedia.

Tidak ada Sev-1/Sev-2 terbuka.

# 26. Catatan Senior System Analyst

## 26.1 Urutan Prioritas Developer

Buat shared architecture terlebih dahulu.

Buat Core Auth dan RBAC.

Buat database migration per modul.

Buat idempotency dan audit.

Buat outbox/inbox.

Masuk ke flow bisnis PMB, Finance, dan Academic.

Lanjutkan LMS, Assessment, Portal, dan Reporting.

## 26.2 Kesalahan Besar yang Harus Dihindari

Membuat join langsung PMB ke Finance.

Membuat payment status sebagai source of truth di PMB.

Membuat final grade di LMS.

Membuat user/password di modul selain Core.

Membuat event tanpa event_key.

Membuat retry tanpa idempotency.

Menyembunyikan menu di frontend tetapi lupa validasi backend.

Membuat dashboard tanpa refreshed_at.

Membuat snapshot tanpa source_event_key.

Membuat handover tanpa unique applicant_ref_id.

# 27. Lampiran: Dokumen Baseline

| Dokumen | Fungsi dalam Implementasi |
| --- | --- |
| PRD Global UNSIA | Menjadi sumber arah produk, arsitektur modular, dan prinsip failure isolation. |
| BRD Global dan Per Modul | Menjadi sumber proses bisnis, scope modul, ownership data, dan acceptance bisnis. |
| FSD Per Modul | Menjadi sumber spesifikasi fungsi, menu, form, validasi, output, dan integrasi. |
| DBML Global | Menjadi sumber struktur database, tabel teknis, dan boundary antar modul. |
| OpenAPI/Swagger | Menjadi kontrak endpoint, request, response, header, security, dan error. |
| UAT Scenario dan QA Test Plan | Menjadi sumber quality gate, test scenario, severity, dan sign-off. |

- Akhir Dokumen ---
