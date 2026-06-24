---
title: "Rencana Kerja Developer ERP UNSIA"
source_file: "Rencana_Kerja_Developer_ERP_UNSIA.docx"
format: markdown
---

# Rencana Kerja Developer ERP UNSIA

UNSIA

# RENCANA KERJA DEVELOPER
ERP Pendidikan / SIAKAD Terintegrasi

Dari Finalisasi Boundary sampai Release Bertahap

| Item | Isi |
| --- | --- |
| Dokumen | Developer Workplan dan Delivery Roadmap |
| Produk | ERP Pendidikan / SIAKAD Terintegrasi UNSIA |
| Versi | v1.0 Working Draft |
| Tanggal | 22 Juni 2026 |
| Target Pembaca | Developer Backend, Frontend, QA, DBA, DevOps, Technical Lead, Product Owner, dan System Analyst |
| Basis | PRD, BRD, FSD, OpenAPI Contract, DBML, UAT Scenario, Event Contract, RBAC Matrix, dan State Machine |

Dokumen ini digunakan sebagai pegangan teknis untuk menurunkan requirement ERP menjadi backlog, sprint, testing, dan release bertahap.

# 1. Tujuan Dokumen

Dokumen ini menjelaskan langkah kerja teknis untuk developer dalam membangun ERP Pendidikan / SIAKAD Terintegrasi UNSIA secara bertahap. Fokusnya adalah memastikan setiap modul memiliki batas tanggung jawab yang jelas, role dan permission tervalidasi, API dan event contract konsisten, state machine terdokumentasi, ERD/DBML siap diimplementasikan, backlog tersusun, development berjalan per sprint, pengujian mengikuti quality gate, dan release dilakukan secara bertahap.

Dokumen ini tidak menggantikan PRD, BRD, FSD, OpenAPI, DBML, atau UAT Scenario. Dokumen ini berfungsi sebagai rencana eksekusi teknis agar tim development memiliki urutan kerja yang sama.

# 2. Prinsip Eksekusi Developer

Setiap modul memiliki database dan ownership data masing-masing.

Tidak ada direct cross-database join untuk transaksi online.

Semua integrasi lintas modul memakai API, event, snapshot/read model, dan reconciliation.

Endpoint protected wajib menegakkan token, active role, permission, application, dan data scope di backend.

Critical command wajib idempotent agar retry tidak menghasilkan data ganda.

Perubahan data penting wajib membuat audit log dan outbox event.

Frontend wajib menampilkan loading, empty, error, degraded mode, dan freshness status untuk snapshot/read model.

Setiap fase harus memiliki deliverable, acceptance criteria, evidence, dan sign-off.

# 3. Urutan Kerja Utama

| No | Tahap | Tujuan | Output Utama |
| --- | --- | --- | --- |
| 1 | Finalisasi module boundary | Menetapkan batas domain, source of truth, database owner, dan data lintas modul yang boleh disimpan sebagai snapshot. | Module Boundary Matrix |
| 2 | Finalisasi role dan permission | Menentukan role, permission, menu, endpoint access, action, dan data scope. | RBAC Matrix |
| 3 | Finalisasi API contract | Mengunci endpoint, request, response, error code, auth header, idempotency, dan versioning. | OpenAPI/Swagger Final |
| 4 | Finalisasi event contract | Mengunci event_name, event_version, event_key, payload, producer, consumer, retry, dan DLQ. | Event Catalog dan Event Payload Schema |
| 5 | Finalisasi state machine | Mengunci lifecycle status, allowed transition, guard condition, actor, audit, dan error rule. | State Machine per Domain |
| 6 | Finalisasi ERD/DBML per modul | Menetapkan tabel, field, index, unique constraint, internal FK, external_ref_id, outbox/inbox, dan audit table. | DBML Final per Modul |
| 7 | Buat backlog per modul | Menurunkan requirement menjadi epic, feature, user story, task, test case, dan acceptance criteria. | Product Backlog dan Sprint Backlog |
| 8 | Development per sprint | Membangun modul sesuai prioritas MVP, dependency, dan critical path. | Increment siap diuji |
| 9 | Test per quality gate | Melakukan functional, API, integration, RBAC, state machine, migration, event, degraded, UAT, smoke. | Test Evidence dan Defect Report |
| 10 | Release bertahap | Melakukan release per modul/flow setelah quality gate pass dan sign-off. | Release Candidate, Pilot, Go-Live |

# 4. Detail Tahap 1 - Finalisasi Module Boundary

Module boundary menentukan data apa yang dimiliki oleh setiap modul, database apa yang dipakai, data apa yang hanya menjadi snapshot, dan integrasi apa yang diperlukan. Tahap ini harus selesai sebelum developer membuat tabel dan endpoint.

| Modul | Database | Data Milik Modul | Data Lintas Modul Lokal | Source of Truth |
| --- | --- | --- | --- | --- |
| Core | core_db | Person, user, role, permission, session, service client, app registry, audit global | Scope reference, application reference | Credential, role, permission, active role, service token |
| Referensi | reference_db | Study program, academic year, academic period, status code, payment component, document type | User_ref_id untuk audit | Master data lintas modul |
| CRM | crm_db | Campaign, lead, agent, referral, follow-up, commission | Person snapshot, applicant_ref_id | Lead dan marketing pipeline |
| PMB | pmb_db | Applicant, biodata, document, selection, LoA, handover | Person snapshot, invoice/payment snapshot, assessment result reference, student_ref_id | Applicant sebelum menjadi mahasiswa |
| Finance | finance_db | Invoice, invoice item, payment, callback, receipt, clearance, scholarship | Customer snapshot, academic_period_ref_id | Transaksi keuangan dan clearance |
| Academic | academic_db | Student, NIM, curriculum, course, class, KRS, final grade, KHS, transcript, alumni | Applicant_ref_id, lecturer_ref_id, clearance snapshot | Mahasiswa, kelas, KRS, nilai final |
| HRIS | hris_db | Employee, lecturer, homebase, unit, position, active status | Person snapshot, study_program_ref_id | Dosen dan pegawai |
| LMS | lms_db | Online class, enrollment, session, material, assignment, attendance, progress, grade input | Academic class snapshot, student snapshot, lecturer snapshot | Aktivitas pembelajaran online |
| Assessment | assessment_db | Question bank, question version, session, attempt, answer, scoring result | Participant snapshot, context_ref_id | Assessment engine |
| Portal | portal_db | Notification, dashboard read model, preference, shortcut, activity log | User/role snapshot, aggregated payload | Presentation layer dan notifikasi |

## Checklist Tahap 1

Setiap data utama sudah memiliki satu owner modul.

Setiap database hanya berisi tabel domain sendiri dan tabel teknis wajib.

Semua external reference sudah diberi nama konsisten, misalnya person_ref_id, applicant_ref_id, student_ref_id, invoice_ref_id.

Tidak ada kebutuhan transaksi online yang bergantung pada join lintas database.

Snapshot/read model sudah diberi source_event_key, synced_at, source_module, dan freshness_status.

Dependency antar modul sudah dipetakan: API call, event, read model, atau manual reconciliation.

## Definition of Done Tahap 1

Module Boundary Matrix disetujui Product Owner, System Analyst, Technical Lead, DBA, dan Owner Modul.

Tidak ada data ownership ganda.

Tidak ada rencana cross-database FK.

Tidak ada direct cross-database join untuk OLTP.

Dependency modul sudah siap diturunkan menjadi API dan event contract.

# 5. Detail Tahap 2 - Finalisasi Role dan Permission

Role dan permission harus diturunkan sampai level menu, action button, endpoint, dan data scope. Frontend boleh menyembunyikan menu, tetapi backend tetap wajib memvalidasi role, permission, dan scope.

| Role | Modul/Menu | Data Scope | Aksi Utama |
| --- | --- | --- | --- |
| Super Admin/Admin BPPTI | Core, registry, role, permission, audit, service token | Global teknis | Manage user, role, permission, impersonation, service token |
| Admin Referensi | Referensi | Global referensi | Manage master data, status code, payment component |
| Admin CRM/Marketing | CRM | CRM domain | Manage campaign, lead, follow-up, conversion |
| Agent/Mitra | CRM limited | Own referral/lead | Create/view own lead, monitor conversion |
| Pendaftar | PMB public/portal | Self | Isi biodata, upload dokumen, lihat invoice, LoA |
| Admin PMB | PMB | PMB domain | Verifikasi applicant, dokumen, seleksi, LoA, handover |
| Admin Finance | Finance | Finance domain | Invoice, payment, manual verification, clearance, report |
| Admin Akademik Biro | Academic | Academic global | Kalender, student, NIM, class, KRS, grade, KHS, transcript |
| Kaprodi/Admin Prodi | Academic | study_program_id | Kurikulum prodi, kelas prodi, monitoring mahasiswa prodi |
| Dosen | LMS/Academic limited | Assigned class | Kelas ajar, materi, tugas, presensi, grade input |
| Dosen PA | Academic advisory | Advisor scope | Approval KRS mahasiswa bimbingan |
| Mahasiswa | Portal/Academic/LMS | Self | KRS, LMS, invoice, KHS, transcript |
| Admin SDM | HRIS | HRIS domain | Employee, lecturer, homebase, active status |
| Admin Assessment | Assessment | Assessment domain | Bank soal, session, scoring, result publish |
| Pimpinan | Portal/Reporting | Read-only aggregate | Dashboard dan KPI |

## Checklist Tahap 2

Daftar role sudah lengkap dan tidak overlap berlebihan.

Setiap role memiliki data scope yang jelas: global, domain, study_program, advisor, assigned_class, self, own_lead, atau read_only.

Permission dibuat granular dengan pola module.resource.action, misalnya pmb.applicant.verify_document.

Action sensitif memiliki permission khusus, bukan sekadar akses menu.

Direct endpoint access tanpa permission harus ditolak dengan error 403.

Setiap perubahan role/permission wajib audit.

## Definition of Done Tahap 2

RBAC Matrix selesai dan disetujui.

Seed role dan permission siap dimasukkan ke Core.

Middleware permission dan data scope sudah memiliki desain final.

Test case RBAC/scope sudah dibuat minimal untuk positive, negative, dan direct URL/API access.

# 6. Detail Tahap 3 - Finalisasi API Contract

API contract menjadi kontrak utama antara backend, frontend, QA, dan integrasi lintas modul. Setelah final, perubahan breaking harus melalui versioning atau change approval.

| Area API | Standar | Catatan Developer |
| --- | --- | --- |
| Standar URL | /api/v1/{module}/{resource} | Gunakan prefix versi dan nama resource plural. |
| Header protected | Authorization, X-Application-Code, X-Active-Role, X-Correlation-Id | Wajib untuk endpoint protected. |
| Idempotency | Idempotency-Key | Wajib untuk command kritis. |
| Response success | success, message, data, meta.trace_id, meta.timestamp | Envelope konsisten lintas modul. |
| Response error | success=false, error.code, error.message, error.details, meta | Error code tidak boleh bebas. |
| Pagination | page, per_page, total, sort, filter | Wajib untuk list page utama. |
| Security | Bearer token + permission + scope | Bukan hanya validasi token. |
| Audit | x-audit-log atau mapping audit action | Aksi sensitif wajib audit. |

## Endpoint Prioritas MVP

| Modul | Endpoint | Tujuan |
| --- | --- | --- |
| Core | POST /auth/login, POST /auth/refresh, GET /auth/me, POST /auth/switch-role, GET /applications | Foundation login dan app launcher |
| Referensi | GET/POST /ref/study-programs, /ref/academic-years, /ref/academic-periods, /ref/status-codes | Master data awal |
| CRM | GET/POST /crm/leads, POST /crm/leads/{id}/convert-to-applicant | Lead to applicant |
| PMB | POST /pmb/applicants, POST /submit, POST /documents, POST /request-invoice, POST /issue-loa, POST /handover-to-academic | Applicant lifecycle |
| Finance | POST /finance/invoices, GET /finance/invoices/{id}, POST /payment-callbacks/{provider}, POST /payment-verifications, GET /clearances | Invoice dan clearance |
| Academic | POST /students/generate-from-applicant, POST /classes, POST /krs, POST /krs/{id}/submit, POST /krs/{id}/approve, POST /grades/finalize | Student dan KRS |
| LMS | POST /classes/sync-from-academic, POST /enrollments/sync-from-krs, POST /grade-syncs | Class/enrollment/grade input |
| Assessment | POST /sessions, POST /attempts, POST /results/publish | Assessment engine |
| Portal | POST /notifications, GET /dashboard | Notification dan dashboard |

## Definition of Done Tahap 3

OpenAPI/Swagger dapat dibuka tanpa error.

Semua endpoint P0 memiliki request schema, response schema, error response, security, role/permission, dan idempotency flag.

QA dapat membuat API test dari contract tanpa bertanya ulang ke developer.

Frontend dapat membuat integration stub/mock berdasarkan OpenAPI.

# 7. Detail Tahap 4 - Finalisasi Event Contract

Event contract mengatur komunikasi asynchronous lintas modul. Event wajib memiliki identitas unik, payload stabil, producer jelas, consumer jelas, dan strategi retry/DLQ.

| Event Name | Producer | Consumer | Fungsi |
| --- | --- | --- | --- |
| core.person_updated | Core | CRM, PMB, HRIS, Academic, Portal | Update person snapshot |
| reference.study_program_updated | Referensi | PMB, Academic, HRIS, LMS, Portal | Update prodi snapshot |
| reference.academic_period_updated | Referensi | PMB, Finance, Academic, LMS, Assessment, Portal | Update periode akademik |
| crm.lead_qualified | CRM | PMB | Lead siap dikonversi |
| pmb.applicant_created | PMB | Finance, Assessment, Portal | Applicant baru tersedia |
| finance.invoice_created | Finance | PMB, Portal | Invoice dibuat |
| finance.payment_paid | Finance | PMB, Academic, Portal | Pembayaran berhasil |
| finance.clearance_changed | Finance | PMB, Academic, LMS, Portal | Status clearance berubah |
| pmb.ready_for_academic | PMB | Academic | Applicant siap menjadi mahasiswa |
| academic.student_created | Academic | PMB, Finance, LMS, Portal | Mahasiswa/NIM dibuat |
| academic.class_opened | Academic | LMS, Portal | Kelas akademik dibuka |
| academic.krs_approved | Academic | LMS, Finance, Portal | KRS disetujui |
| lms.grade_input_submitted | LMS | Academic | Nilai input dari LMS |
| assessment.result_calculated | Assessment | PMB, LMS, Academic, Portal | Hasil assessment tersedia |

## Event Envelope Wajib

event_name: nama event stabil, misalnya finance.payment_paid.

event_version: versi schema, misalnya v1.

event_key: key unik deterministik untuk mencegah duplicate processing.

publisher_service: service penerbit event.

aggregate_type dan aggregate_id: objek bisnis utama.

correlation_id dan causation_id: trace lintas request/event.

occurred_at: waktu kejadian bisnis.

payload: data minimum yang disepakati consumer.

## Definition of Done Tahap 4

Event catalog disetujui producer dan consumer.

Setiap event memiliki payload schema dan sample payload.

Outbox dan inbox table siap di semua database modul.

Retry policy, max retry, DLQ, dan replay permission sudah ditentukan.

Duplicate event test case sudah dibuat untuk semua event kritis.

# 8. Detail Tahap 5 - Finalisasi State Machine

State machine memastikan status bisnis tidak berubah sembarangan. Setiap status transition harus memiliki allowed actor, guard condition, audit log, dan error response untuk invalid transition.

| Domain | Transition Utama | Guard Condition | Audit/History |
| --- | --- | --- | --- |
| Applicant | DRAFT -> SUBMITTED -> VERIFIED -> ACCEPTED -> LOA_ISSUED -> HANDED_OVER | Biodata lengkap, dokumen valid, policy payment valid, LoA terbit | Applicant status history |
| Invoice | DRAFT -> ISSUED -> PARTIALLY_PAID -> PAID -> CANCELLED/EXPIRED | Nominal valid, invoice aktif, payment valid | Invoice status history |
| Payment | RECEIVED -> VERIFIED -> POSTED -> FAILED/REVERSED | Signature valid, provider_event_id unik, amount sesuai | Payment audit |
| Clearance | BLOCKED -> CONDITIONAL -> CLEARED -> REVOKED | Policy terpenuhi atau ada approval khusus | Clearance history |
| KRS | DRAFT -> SUBMITTED -> APPROVED -> FINALIZED -> CANCELLED | Mahasiswa aktif, periode KRS buka, clearance valid, dosen PA approve | KRS status history |
| Grade | INPUTTED -> SUBMITTED -> VERIFIED -> FINALIZED -> CORRECTED | Role dosen/admin valid, periode nilai buka, koreksi pakai reason | Grade correction audit |
| Outbox Event | PENDING -> PUBLISHED -> FAILED -> DLQ | Publish success atau retry limit exceeded | Outbox log |
| Inbox Event | RECEIVED -> PROCESSED -> IGNORED_DUPLICATE -> FAILED -> DLQ | event_key unik, schema valid, consumer berhasil | Inbox log |

## Definition of Done Tahap 5

Setiap status punya daftar allowed transition.

Invalid transition menghasilkan error code konsisten.

Transition sensitif meminta reason/note.

Status history mencatat actor, active role, old status, new status, timestamp, dan correlation_id.

QA memiliki positive dan negative test case untuk setiap transition P0.

# 9. Detail Tahap 6 - Finalisasi ERD/DBML per Modul

ERD/DBML harus mengikuti module boundary. FK hanya boleh internal database. Relasi ke modul lain menggunakan external_ref_id tanpa FK database.

| Database/Modul | Tabel Prioritas | Catatan |
| --- | --- | --- |
| Semua modul | audit_logs, idempotency_keys, outbox_events, inbox_events, reconciliation_mismatch_logs | Tabel teknis wajib |
| Core | persons, users, roles, permissions, user_roles, active_role_sessions, service_clients, applications, event_contracts | Identity dan contract catalog |
| Reference | study_programs, academic_years, academic_periods, status_codes, payment_components, document_types | Master data |
| CRM | campaigns, agents, leads, follow_ups, referrals, commissions | Marketing pipeline |
| PMB | applicants, applicant_biodata, applicant_documents, applicant_invoice_statuses, loa_documents, handover_logs | Applicant lifecycle |
| Finance | invoices, invoice_items, payments, payment_callbacks, receipts, clearances, manual_verifications | Payment dan clearance |
| Academic | students, nim_sequences, curriculums, courses, classes, krs_headers, krs_items, grades, khs, transcripts | Akademik inti |
| HRIS | employees, lecturers, lecturer_homebases, positions, work_units, employment_statuses | SDM dan dosen |
| LMS | lms_classes, lms_enrollments, sessions, materials, assignments, submissions, attendance, lms_grade_inputs | Pembelajaran online |
| Assessment | question_banks, question_versions, assessment_sessions, attempts, answers, scores, result_publications | Assessment engine |
| Portal | notifications, read_markers, preferences, shortcuts, dashboard_read_models, activity_logs | Presentation dan dashboard |

## Index dan Constraint Wajib

idempotency_keys.idempotency_key harus unique.

outbox_events.event_key harus unique.

inbox_events.event_key harus unique per consumer.

students.applicant_ref_id harus unique.

students.nim harus unique.

payment_callbacks.provider_event_id harus unique.

payments.provider_payment_ref harus unique.

KRS item harus unique per student, period, dan course offering.

LMS enrollment harus unique per lms_class_ref_id dan student_ref_id.

Snapshot/read model harus memiliki source_event_key dan synced_at.

# 10. Detail Tahap 7 - Buat Backlog per Modul

Backlog harus dibuat dari requirement final, bukan dari asumsi developer. Setiap backlog item wajib memiliki acceptance criteria dan test case awal.

| Level Backlog | Isi | Contoh |
| --- | --- | --- |
| Epic | Kelompok besar pekerjaan | PMB Applicant Lifecycle |
| Feature | Kemampuan sistem dalam epic | Upload dan verifikasi dokumen applicant |
| User Story | Kebutuhan dari sisi pengguna | Sebagai Admin PMB, saya ingin memverifikasi dokumen agar applicant bisa lanjut seleksi. |
| Technical Task | Pekerjaan teknis developer | Create table applicant_documents, endpoint verify, audit log, event update |
| Acceptance Criteria | Syarat fitur diterima | Dokumen bisa approved/rejected, rejection wajib reason, audit tercatat |
| Test Case | Skenario QA/UAT | Positive approve, reject tanpa reason ditolak, direct access non-PMB ditolak |
| Definition of Done | Syarat selesai developer | Unit test, API test, RBAC test, audit, evidence pass |

## Template User Story

Judul: [Modul] - [Fungsi] - [Aksi].

User story: Sebagai [role], saya ingin [aksi], agar [tujuan bisnis].

Business rule: aturan validasi utama.

API: endpoint yang digunakan.

Database: tabel yang terdampak.

Event: event yang dipublish/consume.

Audit: aksi yang wajib dicatat.

RBAC/scope: role dan scope yang boleh melakukan aksi.

Acceptance criteria: positive, negative, duplicate/idempotency, scope, audit.

Evidence: screenshot, request/response, query validasi, audit log, event log.

# 11. Detail Tahap 8 - Development per Sprint

Development dilakukan berdasarkan dependency. Foundation harus selesai lebih dulu sebelum flow bisnis lintas modul dibangun.

| Sprint | Fokus | Deliverable | Owner Utama |
| --- | --- | --- | --- |
| Sprint 0 | Architecture Foundation | Repo/service structure, database per modul, migration pipeline, shared auth, shared envelope, logging, local environment | Tech Lead, DevOps, DBA |
| Sprint 1 | Core Auth | Login, refresh, me, switch role, application launcher, user seed, role seed | Backend Core, Frontend Shell |
| Sprint 2 | Reference Foundation | Prodi, tahun ajaran, periode akademik, status code, payment component, document type | Backend Reference, Frontend Admin |
| Sprint 3 | CRM + PMB Applicant | Lead, convert lead, applicant, biodata, document upload | Backend CRM/PMB |
| Sprint 4 | Finance Invoice + Payment | Invoice request, invoice detail, payment callback, manual verification, payment_paid event | Backend Finance |
| Sprint 5 | PMB LoA + Handover | Document verification, selection, issue LoA, ready_for_academic, handover retry | PMB + Academic |
| Sprint 6 | Academic Student + NIM | Generate student, NIM sequence, student profile, academic period binding | Academic |
| Sprint 7 | Academic KRS | Class, KRS draft, submit, approval, clearance check, krs_approved event | Academic + Finance |
| Sprint 8 | HRIS + LMS Sync | Lecturer active read model, class sync, enrollment sync, basic LMS class | HRIS + LMS |
| Sprint 9 | Assessment + Grade Input | Question bank, session, attempt, scoring, result publish, LMS grade input to Academic | Assessment + LMS + Academic |
| Sprint 10 | Portal + Dashboard | Notification center, role dashboard, refreshed_at, activity log, read model | Portal + Frontend |
| Sprint 11 | Hardening | API contract test, event test, RBAC test, state machine test, migration dry-run, performance baseline | All Teams + QA |
| Sprint 12 | UAT + Pilot Release | UAT scenario execution, defect closure, sign-off, pilot release, rollback rehearsal | PO, QA, SA, DevOps |

## Sprint Ceremony dan Evidence

Sprint planning harus mengacu pada backlog yang sudah memiliki acceptance criteria.

Daily check harus mencatat blocker lintas modul, terutama API/event dependency.

Sprint review wajib menunjukkan demo fitur, API response, audit log, dan event log bila relevan.

Sprint retro mencatat masalah requirement, defect, integrasi, environment, dan action improvement.

Setiap sprint wajib menghasilkan build yang bisa diuji QA minimal pada environment QA.

# 12. Detail Tahap 9 - Test per Quality Gate

| Quality Gate | Fokus Uji | Exit Criteria |
| --- | --- | --- |
| Requirement Review | Requirement tidak ambigu dan acceptance criteria ada | No ambiguity P0, owner modul setuju |
| Functional Test | Form, list, validasi, status, output, audit | P0 pass, P1 mayor pass atau ada approved workaround |
| API Contract Test | Endpoint, header, payload, envelope, error code, security | Endpoint P0 pass dan trace_id tersedia |
| Integration Test | Komunikasi API/event lintas modul | Tidak ada duplicate record dan integration log tersedia |
| RBAC/Scope Test | Role, permission, endpoint, data scope backend | Unauthorized read/write ditolak |
| State Machine Test | Allowed/invalid transition dan guard condition | Invalid transition ditolak, history lengkap |
| Migration Validation | Schema, seed, FK lokal, unique index, rollback | Migration dry-run pass dan rollback rehearsal aman |
| Event Contract Test | Outbox, inbox, duplicate, retry, DLQ, replay | Duplicate event aman dan DLQ dapat direplay dengan audit |
| Degraded Mode Test | Dependency down dan fallback snapshot/read model | UI menampilkan sync status dan action berisiko dibatasi |
| Regression/Smoke | Critical path setelah perubahan baru | Smoke P0 pass sebelum release candidate |
| UAT | Validasi proses bisnis oleh owner modul | Owner modul sign-off dan Sev-1/Sev-2 closed |

## Severity Defect

| Severity | Kriteria | Contoh | Keputusan |
| --- | --- | --- | --- |
| Sev-1 Critical | Production/blocking risk | Data loss, duplicate NIM/payment, security breach, unauthorized cross-scope access | Stop release, fix immediately |
| Sev-2 High | Major business flow blocked | PMB submit gagal, payment verify gagal, KRS finalize gagal, handover gagal | Must fix before UAT sign-off |
| Sev-3 Medium | Important function impaired with workaround | Filter/export salah, validasi pesan kurang jelas | Fix sesuai prioritas sprint |
| Sev-4 Low | Minor UI/text issue | Typo, spacing, cosmetic issue | Can defer with approval |

# 13. Detail Tahap 10 - Release Bertahap

| Release | Nama | Scope | Kriteria Rilis |
| --- | --- | --- | --- |
| Release 0 | Technical Foundation | Core Auth, RBAC seed, Reference master, shared package, migration, event broker | Internal technical demo |
| Release 1 | PMB + Finance MVP | Lead, applicant, invoice, payment, LoA, notification | Pilot PMB terbatas |
| Release 2 | Academic MVP | Generate NIM, student profile, class, KRS, approval, clearance check | Pilot akademik 1 periode |
| Release 3 | LMS + Assessment MVP | Class sync, enrollment, basic LMS, assessment, grade input | Pilot kelas terpilih |
| Release 4 | Portal + Dashboard | Notification center, dashboard role-based, executive KPI, activity log | Pilot pimpinan dan user utama |
| Release 5 | Hardening + Go-Live Bertahap | Regression, UAT sign-off, backup/restore, rollback, monitoring, runbook | Go-live bertahap per modul/flow |

## Checklist Go/No-Go

Seluruh P0 critical path pass.

Tidak ada defect Sev-1 dan Sev-2 terbuka.

RBAC/scope backend pass.

API contract dan event contract test pass.

Migration dry-run dan rollback rehearsal pass.

Monitoring, alert, logging, dan dashboard teknis aktif.

Backup dan restore sudah diuji.

Owner modul, QA/UAT Lead, Product Owner, Technical Lead, DBA, dan DevOps sign-off.

Release note dan runbook sudah tersedia.

# 14. RACI Delivery

| Aktivitas | Responsible | Accountable | Consulted | Informed |
| --- | --- | --- | --- | --- |
| Module Boundary | System Analyst | Product Owner, Owner Modul | Technical Lead, DBA | Backend Lead, QA Lead |
| Role & Permission | System Analyst, Core Backend | Product Owner | Owner Modul, Security | Frontend, QA |
| API Contract | Backend Lead | Technical Lead | Frontend Lead, QA Lead, SA | Product Owner |
| Event Contract | Backend Lead | Technical Lead | DevOps, DBA, QA | Owner Modul |
| State Machine | System Analyst | Product Owner | Owner Modul, QA | Backend Lead |
| ERD/DBML | DBA | Technical Lead | Backend Lead, SA | QA, DevOps |
| Backlog | Product Owner, SA | Product Owner | Tech Lead, QA, Owner Modul | Developer |
| Development | Developer | Technical Lead | QA, DevOps, SA | PO, Owner Modul |
| Quality Gate | QA/UAT Lead | Product Owner | SA, Tech Lead, DBA, DevOps | Owner Modul |
| Release | DevOps, Tech Lead | Product Owner | QA, DBA, Backend Lead, Frontend Lead | All Stakeholders |

# 15. Final Delivery Checklist

| Deliverable | Owner | Status | Catatan |
| --- | --- | --- | --- |
| Module Boundary Matrix | SA + PO + Owner Modul + DBA | Belum/Selesai | Wajib sebelum ERD/API final |
| RBAC Matrix | SA + Core Backend + Security | Belum/Selesai | Wajib sebelum frontend menu dan backend middleware |
| OpenAPI Contract | Backend Lead + QA | Belum/Selesai | Wajib sebelum API test dan frontend integration |
| Event Contract Catalog | Backend Lead + DevOps | Belum/Selesai | Wajib sebelum asynchronous integration |
| State Machine | SA + QA + Owner Modul | Belum/Selesai | Wajib sebelum status workflow coding |
| ERD/DBML per Modul | DBA + Backend Lead | Belum/Selesai | Wajib sebelum migration |
| Product Backlog | PO + SA | Belum/Selesai | Wajib sebelum sprint planning |
| Sprint Backlog | Tech Lead + Developer | Belum/Selesai | Wajib setiap sprint |
| Unit/Integration/API Test | Developer + QA | Belum/Selesai | Wajib sebelum UAT |
| UAT Evidence | QA/UAT Lead + Owner Modul | Belum/Selesai | Wajib sebelum sign-off |
| Release Note | Tech Lead + DevOps | Belum/Selesai | Wajib sebelum deployment |
| Runbook & Rollback Plan | DevOps + DBA | Belum/Selesai | Wajib sebelum go-live |

# 16. Catatan Implementasi Penting

Jangan membuat payment status sebagai source of truth di PMB atau Academic. Finance tetap owner payment dan clearance.

Jangan membuat final grade di LMS. LMS hanya mengirim grade input; nilai final tetap di Academic.

Jangan membuat user/password di modul selain Core.

Jangan membuat event tanpa event_key deterministik.

Jangan membuat retry tanpa idempotency.

Jangan hanya mengandalkan hide menu di frontend. Backend wajib validasi permission dan scope.

Jangan membuat dashboard lintas modul tanpa refreshed_at dan source status.

Jangan membuat handover PMB ke Academic tanpa unique applicant_ref_id di Academic.

Jangan melakukan direct DB connection ke database modul lain dari service manapun.

# 17. Ringkasan Keputusan untuk Developer

Urutan kerja yang disarankan adalah menyelesaikan foundation architecture terlebih dahulu, kemudian menutup module boundary, RBAC, API contract, event contract, state machine, dan ERD/DBML. Setelah itu backlog disusun per modul dan development berjalan per sprint dengan quality gate ketat. Release dilakukan bertahap berdasarkan flow bisnis yang sudah terbukti pass, bukan berdasarkan modul yang sekadar selesai coding.
