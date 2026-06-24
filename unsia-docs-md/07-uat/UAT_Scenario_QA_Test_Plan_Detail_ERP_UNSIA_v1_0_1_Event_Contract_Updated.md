---
title: "UAT Scenario dan QA Test Plan Detail UNSIA"
source_file: "UAT_Scenario_QA_Test_Plan_Detail_ERP_UNSIA_v1_0_1_Event_Contract_Updated.docx"
format: markdown
---

# UAT Scenario dan QA Test Plan Detail UNSIA

UNSIA

UAT SCENARIO DAN QA TEST PLAN DETAIL

ERP Pendidikan / SIAKAD Terintegrasi UNSIA

v1.0.1 Event Contract Updated - 22 Juni 2026

| Item | Isi |
| --- | --- |
| Dokumen | UAT Scenario dan QA Test Plan Detail |
| Produk | ERP Pendidikan / SIAKAD Terintegrasi UNSIA |
| Versi | v1.0.1 Event Contract Updated |
| Basis | PRD Global v6.5.1, BRD v1.1.1, FSD v1.0.1, DBML v1.0.1, API Contract/OpenAPI, Event Contract, RBAC Matrix, Data Dictionary, dan State Machine. |
| Tujuan | Menetapkan strategi pengujian, skenario UAT, QA test plan, API test, event contract test, outbox/inbox test, retry/DLQ test, reconciliation test, RBAC/scope test, state machine test, regression smoke, test data, defect severity, dan sign-off control. |
| Catatan | Update v1.0.1 menambahkan skenario Event Contract, termasuk outbox, inbox, duplicate event, retry, DLQ, replay, snapshot/read model freshness, dan reconciliation evidence. |

# 1. Tujuan Dokumen

Dokumen ini menetapkan rencana pengujian terperinci untuk ERP Pendidikan / SIAKAD Terintegrasi UNSIA. Fokusnya bukan hanya memeriksa apakah fitur dapat diklik, tetapi memastikan setiap proses bisnis berjalan sesuai requirement, role, data scope, API contract, state machine, audit, idempotency, dan integrasi lintas modul.

Dokumen ini menjadi dasar bagi QA/UAT Lead, Product Owner, Owner Modul, Technical Lead, DBA, Backend Lead, Frontend Lead, dan stakeholder operasional untuk melakukan validasi sebelum modul masuk release candidate atau go-live terbatas.

# 2. Basis dan Ruang Lingkup Pengujian

| Area | Ruang Lingkup |
| --- | --- |
| Modul | Core, Referensi, CRM, PMB, Finance, Akademik, HRIS/SDM, LMS, Assessment, dan Portal. |
| Test level | Requirement review, functional test, API test, integration test, RBAC/scope test, state machine test, migration validation, UAT, regression, smoke, dan NFR baseline. |
| Process coverage | Lead to applicant, applicant to LoA, payment and clearance, handover and generate NIM, KRS, LMS sync, assessment attempt, grade sync, KHS/transcript, notification, dan dashboard. |
| Security coverage | Token, active role, permission, backend data scope, direct URL/API access, service token, idempotency, dan audit trail. |
| Data coverage | Master referensi, periode akademik, kurikulum, applicant, invoice, payment, clearance, student, KRS, class, enrollment, attempt, score, notification, dan audit log. |

# 3. Prinsip Pengujian

| Prinsip | Makna Pengujian |
| --- | --- |
| Backend enforced scope | Setiap endpoint protected wajib memvalidasi token, active role, permission, dan data scope. Menu yang disembunyikan di UI tidak cukup. |
| Source of truth | Modul hanya boleh mengubah data milik domainnya sendiri. Data lintas modul dibaca lewat API, event, read model, atau service resmi. |
| Idempotent critical flow | Payment callback, handover, generate NIM, class sync, enrollment sync, grade sync, dan result publish tidak boleh menghasilkan data ganda saat retry. |
| Audit before convenience | Aksi sensitif harus mencatat actor, active role, timestamp, request/correlation id, reason bila diperlukan, old value, dan new value. |
| State machine compliance | Perubahan status hanya boleh terjadi melalui allowed transition, guard condition, actor yang sah, dan error code yang konsisten. |
| Traceability | Setiap test case harus dapat ditelusuri ke requirement, modul, endpoint/API, table/domain, role, dan acceptance criteria. |

# 4. Peran dan Tanggung Jawab QA/UAT

| Peran | Tanggung Jawab |
| --- | --- |
| Product Owner | Menyetujui prioritas pengujian, menerima atau menolak release berdasarkan quality gate, dan memutuskan go/no-go. |
| System Analyst | Menjaga traceability PRD, BRD, FSD, BPMN, State Machine, API Contract, RBAC, dan UAT scenario. |
| Owner Modul | Memvalidasi proses bisnis, acceptance criteria, data operasional, laporan, dan hasil UAT modul. |
| QA/UAT Lead | Menyusun test plan, mengatur test cycle, memantau defect, menutup evidence, dan menyiapkan sign-off UAT. |
| Technical Lead | Memastikan defect teknis, API, service, performance, dan dependency lintas modul ditangani sebelum release. |
| DBA | Memvalidasi schema, constraint, index, seed data, migration dry-run, rollback rehearsal, dan post-migration validation. |
| Backend Lead | Memvalidasi service logic, endpoint, integration event, idempotency, audit, dan error handling. |
| Frontend Lead | Memvalidasi UI behavior, form validation, loading/empty/error state, responsive layout, dan role-based view. |

# 5. Test Level dan Quality Gate

| Level | Tujuan | Exit Criteria |
| --- | --- | --- |
| Requirement Review | Memastikan requirement tidak ambigu sebelum development atau test execution. | Tidak ada ambiguitas P0, open issue kritis dicatat, owner modul memahami scope. |
| Functional Test | Memastikan fungsi, form, list, action, validasi, output, dan audit berjalan sesuai FSD. | Seluruh P0 pass, P1 mayor pass atau memiliki approved workaround. |
| API Contract Test | Memastikan endpoint, header, payload, response envelope, error envelope, dan security scheme konsisten. | Endpoint P0 pass, error code konsisten, trace_id/correlation_id tersedia. |
| Integration Test | Memastikan komunikasi lintas modul berjalan dengan service token, event key, retry, dan audit. | Tidak ada duplicate record, event retry aman, integration log tersedia. |
| RBAC/Scope Test | Memastikan role, permission, action, endpoint, dan data scope ditegakkan di backend. | Tidak ada unauthorized read/write, direct URL/API access ditolak. |
| State Machine Test | Memastikan status lifecycle bergerak sesuai transition, guard condition, actor, audit, dan error. | Allowed transition pass, invalid transition ditolak, audit/status history lengkap. |
| Migration Validation | Memastikan schema, seed, FK, unique index, duplicate check, dan rollback rehearsal aman. | Tidak ada blocking migration defect, validation query pass. |
| UAT | Memastikan flow operasional diterima owner bisnis dan user perwakilan. | Owner modul sign-off, semua defect Sev-1/Sev-2 closed. |
| Regression/Smoke | Memastikan perubahan baru tidak merusak critical path. | Smoke P0 pass sebelum release candidate. |
| Event Contract Test | Memastikan event_name, event_version, event_key, payload schema, producer, consumer, outbox, inbox, retry, DLQ, dan reconciliation berjalan sesuai kontrak. | Event P0 pass, duplicate event aman, retry tidak membuat data ganda, DLQ dapat direplay dengan audit. |
| Degraded Mode Test | Memastikan modul tetap berjalan terbatas ketika dependency down. | UI menampilkan synced_at/refreshed_at, source status, dan action yang berisiko masuk pending_review atau retry queue. |

# 6. Entry Criteria dan Exit Criteria

| Kategori | Entry Criteria | Exit Criteria |
| --- | --- | --- |
| Requirement | PRD, BRD, FSD, API, DB, State Machine, RBAC, dan Data Dictionary tersedia minimal final draft. | Tidak ada requirement P0 yang belum punya acceptance criteria. |
| Environment | QA/Staging environment tersedia, database migrated, seed data minimal tersedia, API base URL stabil. | Environment dapat direbuild dan smoke test pass. |
| Access | Akun test semua role tersedia dengan scope yang jelas. | Seluruh role P0 dapat login dan diuji sesuai permission. |
| Data | Seed data meliputi applicant, invoice, payment, student, KRS, class, lecturer, attempt, notification. | Test data mendukung positive, negative, role-scope, idempotency, dan regression case. |
| Defect | Defect tracker aktif, severity disepakati, evidence format tersedia. | Tidak ada Sev-1/Sev-2 open; Sev-3/Sev-4 memiliki keputusan defer atau fix. |
| Sign-off | Owner modul dan QA/UAT Lead tersedia untuk review hasil. | Sign-off per modul lengkap atau approved with notes yang tidak mengubah scope kritis. |

# 7. Test Environment dan Test Data

| Area | Kebutuhan Minimum |
| --- | --- |
| Environment | QA environment untuk functional/API test, staging environment untuk SIT/UAT, dan rehearsal environment untuk migration dry-run. |
| User Role | Super Admin, Admin BPPTI, Admin Referensi, Admin CRM, Agen, Pendaftar, Admin PMB, Admin Finance, Admin Akademik Biro, Kaprodi/Admin Prodi, Dosen, Dosen PA, Mahasiswa, Admin SDM, Admin LMS, Admin Assessment, Pimpinan. |
| Seed Data | Tahun Ajaran, Periode Akademik, Prodi, Kurikulum aktif, gelombang PMB, applicant, invoice, payment, student, class, KRS, dosen, LMS class, assessment session, notification. |
| Evidence | Screenshot UI, API request/response, database validation query, audit log, integration log, idempotency log, dan defect ID bila gagal. |
| Data Reset | Setiap siklus SIT/UAT harus memiliki prosedur reset data atau baseline snapshot agar hasil test bisa direproduksi. |

# 8. Defect Severity dan Aturan Triage

| Severity | Kriteria | Contoh | Tindak Lanjut |
| --- | --- | --- | --- |
| Sev-1 Critical | Production/blocking risk | Data loss, security breach, duplicate student/NIM/payment, payment wrong posting, unauthorized cross-scope access | Stop release; fix immediately; PO+Tech Lead+QA approval required |
| Sev-2 High | Major business flow blocked | PMB cannot submit, Finance cannot verify, KRS cannot finalize, handover fails, grade cannot finalize | Must fix before UAT sign-off/release |
| Sev-3 Medium | Important function impaired with workaround | Filter/export wrong, notification delayed, validation message incomplete | Fix before release if affects P0/P1 flow |
| Sev-4 Low | Minor UI/content/cosmetic defect | Text typo, alignment issue, non-critical layout inconsistency | Can be deferred with PO approval |

# 9. UAT Scenario Per Modul

Bagian ini memuat skenario UAT utama per modul. Detail status eksekusi, tester, evidence, dan defect ID dikelola pada workbook Excel agar mudah difilter dan dipantau.

| ID | Modul | Scenario | Priority | Role/User | Expected Result |
| --- | --- | --- | --- | --- | --- |
| UAT-001 | Core | Login multi-role dan memilih active role | P0 | Super Admin/Admin BPPTI | Menu, token, permission, dan scope mengikuti active role yang dipilih. |
| UAT-002 | Core | Admin Prodi akses data prodi lain via URL langsung | P0 | Super Admin/Admin BPPTI | Backend menolak akses dengan SCOPE_DENIED/403 dan mencatat access denied log. |
| UAT-003 | Core | Impersonation dengan reason | P0 | Super Admin/Admin BPPTI | Actor asli, user target, reason, waktu mulai, dan waktu selesai tercatat. |
| UAT-004 | Referensi | Membuat Tahun Ajaran dan Periode Akademik | P0 | Admin Referensi/Admin Akademik | Periode tersedia untuk PMB/Finance/Academic dan tidak rancu dengan Tahun Kurikulum. |
| UAT-005 | Referensi | Status code standar digunakan transaksi | P0 | Admin Referensi/Admin Akademik | Transaksi memakai status code standar, bukan string bebas. |
| UAT-006 | CRM | Convert lead qualified menjadi applicant PMB | P1 | Admin CRM/Agen | Applicant PMB terbentuk satu kali, lead status berubah, conversion audit tercatat. |
| UAT-007 | CRM | Agen hanya melihat lead miliknya | P1 | Admin CRM/Agen | Agen hanya melihat lead/referral miliknya; akses langsung ke lead lain ditolak. |
| UAT-008 | PMB | Pendaftar mengisi biodata dan upload dokumen | P0 | Pendaftar/Admin PMB | Status applicant/dokumen berubah sesuai rule; dokumen masuk antrian verifikasi. |
| UAT-009 | PMB | Admin PMB verifikasi dokumen dan issue LoA | P0 | Pendaftar/Admin PMB | LoA hanya terbit jika semua guard condition terpenuhi dan audit tersimpan. |
| UAT-010 | PMB | Handover PMB ke Akademik | P0 | Pendaftar/Admin PMB | Student dan NIM terbentuk satu kali; handover log dan NIM audit tercatat. |
| UAT-011 | Finance | Payment callback gateway berhasil | P0 | Admin Finance/Payment Gateway | Payment posted, invoice paid/partial sesuai amount, journal/audit tercatat. |
| UAT-012 | Finance | Callback duplikat tidak menggandakan payment | P0 | Admin Finance/Payment Gateway | Callback kedua masuk duplicate_ignored; tidak ada payment/journal dobel. |
| UAT-013 | Finance | Clearance KRS blocked karena tunggakan | P0 | Admin Finance/Payment Gateway | Finance mengembalikan blocked/conditional; KRS tidak bisa finalized sesuai policy. |
| UAT-014 | Akademik | Generate NIM dari applicant PMB | P0 | Mahasiswa/Dosen PA/Admin Akademik | NIM unik, student memiliki entry_period_id dan curriculum_id permanen. |
| UAT-015 | Akademik | KRS finalized setelah validasi | P0 | Mahasiswa/Dosen PA/Admin Akademik | KRS finalized dan item KRS siap disinkronkan ke LMS. |
| UAT-016 | Akademik | Grade finalization dengan history | P0 | Mahasiswa/Dosen PA/Admin Akademik | Nilai final tercatat, koreksi membutuhkan reason dan grade history tersimpan. |
| UAT-017 | HRIS/SDM | Dosen aktif dapat dipakai Akademik/LMS | P1 | Admin SDM | Academic/LMS dapat membaca dosen aktif; tidak dapat mengubah data HRIS. |
| UAT-018 | HRIS/SDM | Dosen nonaktif tidak bisa diplot kelas baru | P1 | Admin SDM | Sistem menolak plotting baru atau memberi validasi status tidak aktif. |
| UAT-019 | LMS | Sync kelas dari Academic ke LMS | P0 | Dosen/Mahasiswa/Admin LMS | LMS class terbentuk berdasarkan academic_class_id unique; LMS tidak membuat kelas mandiri. |
| UAT-020 | LMS | Enrollment dari KRS valid | P0 | Dosen/Mahasiswa/Admin LMS | Mahasiswa hanya masuk kelas dari KRS valid. |
| UAT-021 | LMS | Grade source dikirim ke Academic | P0 | Dosen/Mahasiswa/Admin LMS | Academic menerima source grade; final grade tetap milik Academic. |
| UAT-022 | Assessment | CBT PMB session dan participant | P0 | Admin Assessment/Peserta | Attempt berjalan dari assigned sampai scored/published; result API idempotent. |
| UAT-023 | Assessment | Soal yang sudah digunakan dibuat versi baru | P0 | Admin Assessment/Peserta | Soal historis tidak berubah; versioning aman untuk audit. |
| UAT-024 | Portal | Dashboard mahasiswa role-based | P1 | Mahasiswa/Pimpinan/Semua Role | Portal menampilkan ringkasan dari modul sumber sesuai role; tidak mengubah data bisnis langsung. |
| UAT-025 | Portal | Notification center menerima event lintas modul | P1 | Mahasiswa/Pimpinan/Semua Role | Notifikasi tampil, read marker tersimpan, link menuju modul sumber benar. |
| UAT-026 | Portal | Executive dashboard read-only | P1 | Mahasiswa/Pimpinan/Semua Role | Pimpinan melihat data agregat dan tidak dapat mengubah transaksi. |
| E2E-001 | CRM, PMB, Finance, Academic, LMS, Assessment, Portal | Lead to Alumni baseline | P0 | Multi-role | Seluruh lifecycle memiliki traceability, correlation_id, audit, dan tidak ada data orphan. |
| E2E-002 | PMB, Finance, Assessment, Portal | PMB registration sampai LoA | P0 | Multi-role | LoA hanya terbit setelah dokumen dan payment policy valid. |
| E2E-003 | Finance, PMB, Academic, Portal | Payment callback and reconciliation | P0 | Multi-role | Status bayar resmi dari Finance dan dipakai modul lain. |
| E2E-004 | Academic, Finance, LMS, Portal | KRS sampai LMS enrollment | P0 | Multi-role | Kelas LMS sama dengan kelas Academic dan enrollment hanya dari KRS valid. |
| E2E-005 | Assessment, LMS, Academic | Assessment result to final grade | P0 | Multi-role | Final grade tetap di Academic dan history lengkap. |
| UAT-027 | Event Contract | Event contract catalog dapat dibaca technical admin | P0 | Technical Admin/DevOps | Daftar event menampilkan event_name, version, publisher, consumer, status, dan payload schema. |
| UAT-028 | Event Contract | Outbox event dibuat setelah transaksi lokal commit | P0 | System/QA | Event masuk outbox dengan event_key unik, correlation_id, causation_id, dan status PENDING/PUBLISHED. |
| UAT-029 | Event Contract | Inbox consumer memproses event satu kali | P0 | System/QA | Consumer mencatat event_key pada inbox dan tidak memproses ulang event duplikat. |
| UAT-030 | Event Contract | Duplicate finance.payment_paid tidak menggandakan payment status snapshot | P0 | Admin Finance/Admin PMB | Snapshot PMB/Academic tidak dobel dan duplicate event diberi status IGNORED_DUPLICATE. |
| UAT-031 | Event Contract | Consumer down membuat event masuk retry queue | P0 | DevOps/QA | Event memiliki retry_count, next_retry_at, last_error, dan diproses ulang saat consumer pulih. |
| UAT-032 | Event Contract | Event gagal permanen masuk DLQ dan dapat direplay | P0 | DevOps/SRE | Replay membutuhkan reason, tercatat audit, dan tidak membuat data duplikat. |
| UAT-033 | Event Contract | Snapshot stale menampilkan synced_at dan label data tidak real-time | P1 | Admin Modul/User | UI menampilkan freshness data dan membatasi aksi yang bergantung pada source real-time. |
| UAT-034 | Event Contract | Reconciliation mismatch source vs snapshot terdeteksi | P0 | Admin Modul/DevOps | Mismatch report terbentuk dengan status OPEN, CORRECTED, IGNORED, atau PENDING_REVIEW. |
| UAT-035 | Event Contract | Event version unsupported ditolak dengan error standar | P1 | System/QA | Consumer mengembalikan EVENT_VERSION_UNSUPPORTED dan event masuk DLQ/compatibility queue. |

# 10. Integration Test Matrix

| Test ID | Integration ID | Provider | Consumer | Trigger | Event Key | Failure Handling |
| --- | --- | --- | --- | --- | --- | --- |
| IT-001 | INT-CORE-001 | Core | Semua modul | Token validation dan active role | N/A | Reject 401/403 when invalid |
| IT-002 | INT-CRM-PMB-001 | CRM | PMB | Convert lead qualified menjadi applicant | crm:lead:{id}:convert | No duplicate applicant |
| IT-003 | INT-PMB-FIN-001 | PMB | Finance | Request invoice PMB | pmb:invoice:{applicant_id}:{context} | Return existing invoice if duplicate |
| IT-004 | INT-FIN-PMB-001 | Finance | PMB/Portal | Publish payment status | finance:payment:{provider_event_id} | Retry with event key |
| IT-005 | INT-FIN-ACD-001 | Finance | Academic | Clearance check untuk KRS/ujian/KHS | N/A | Return blocked/conditional/clear |
| IT-006 | INT-PMB-ASM-001 | PMB | Assessment | Membuat CBT session dan participant PMB | pmb:cbt:{wave_id}:{applicant_id} | No duplicate participant |
| IT-007 | INT-ASM-PMB-001 | Assessment | PMB | Kirim hasil CBT PMB | assessment:result:{attempt_id} | Retry publish result |
| IT-008 | INT-PMB-ACD-001 | PMB | Academic | Handover applicant menjadi student | pmb:handover:{applicant_id} | No duplicate student/NIM |
| IT-009 | INT-HRIS-ACD-001 | HRIS | Academic/LMS | Lecturer read model | N/A | Consumer cannot edit lecturer |
| IT-010 | INT-ACD-LMS-001 | Academic | LMS | Class sync | academic:class:{id}:sync | Upsert by academic_class_id |
| IT-011 | INT-LMS-ACD-001 | LMS | Academic | Grade source import | lms:grade:{source_ref_id} | No duplicate source grade |
| IT-012 | INT-PORTAL-ALL-001 | All modules | Portal | Notification event publish | notification:{source}:{event_id} | No duplicate notification |
| IT-013 | INT-EVT-001 | All modules | Event Broker | Publish outbox event after commit | {module}:{event}:{aggregate_id}:v1 | Rollback transaction tidak boleh publish event. |
| IT-014 | INT-EVT-002 | Event Broker | Consumer module | Consume inbox event | event_key | Duplicate event menjadi IGNORED_DUPLICATE. |
| IT-015 | INT-EVT-003 | Finance | PMB/Academic | finance.payment_paid duplicate publish | finance.payment_paid:{payment_id}:v1 | Snapshot tidak ganda. |
| IT-016 | INT-EVT-004 | Academic | LMS | academic.krs_approved delayed event | academic.krs_approved:{krs_id}:v1 | Enrollment sync retry tanpa duplikat. |
| IT-017 | INT-EVT-005 | Assessment | PMB/LMS/Academic | assessment.result_calculated retry | assessment.result_calculated:{attempt_id}:v1 | Result consumed satu kali. |
| IT-018 | INT-EVT-006 | Consumer | Retry Queue | Consumer temporary failure | event_key | retry_count dan next_retry_at terisi. |
| IT-019 | INT-EVT-007 | Retry Queue | DLQ | Retry limit exceeded | event_key | Event masuk DLQ dengan last_error. |
| IT-020 | INT-EVT-008 | Reconciliation Job | Consumer module | Source vs snapshot mismatch | source_event_key | Mismatch report dibuat. |

# 11. RBAC dan Data Scope Test

| ID | Scope Category | Role | Data Scope | Allowed Action | Expected Control |
| --- | --- | --- | --- | --- | --- |
| RBAC-001 | RBAC-001 | Global Admin | Core global management | Global | CRUD role/permission/application/audit |
| RBAC-002 | RBAC-002 | Module Admin | PMB applicant management | PMB domain | CRUD applicant, document verification, LoA, handover |
| RBAC-003 | RBAC-003 | Study Program Scope | Academic data by study_program_id | study_program_id | Read/update class/KRS/monitoring for assigned prodi |
| RBAC-004 | RBAC-004 | Assigned Class Scope | LMS/grade for assigned class | assigned_class | Manage session/material/task/grade input for assigned class |
| RBAC-005 | RBAC-005 | Advisor Scope | KRS approval mahasiswa bimbingan | advisor_scope | Approve/reject KRS assigned students |
| RBAC-006 | RBAC-006 | Self Scope | Own KRS, invoice, LMS, grade, transcript | self | Read/update own eligible data |
| RBAC-007 | RBAC-007 | Agent Scope | Own referral/lead | agent_scope | Read own lead/referral |
| RBAC-008 | RBAC-008 | Executive Read-Only | Dashboard aggregate | read-only aggregate | Read dashboard/report |
| RBAC-009 | RBAC-009 | Portal Shortcut Scope | Role-based shortcut/menu | active_role | Menu follows active role |
| RBAC-010 | RBAC-010 | Backend Direct URL | Direct endpoint access | scope-specific | Attempt direct API call |

# 12. State Machine Test

| ID | State Machine | State Path | Event | Guard Condition | Expected Audit/Control |
| --- | --- | --- | --- | --- | --- |
| STM-001 | SM-001 | Applicant PMB Lifecycle | draft -> submitted -> accepted -> handed_over | submit/approve/handover | Required biodata, document, payment policy |
| STM-002 | SM-002 | Applicant Document Verification | uploaded -> verified/rejected | verify/reject | Reviewer role and rejection reason if rejected |
| STM-003 | SM-003 | Invoice Lifecycle | draft -> issued -> paid/overdue/cancelled/voided | issue/pay/overdue/cancel/void | Invoice item valid, payment policy, reason for void |
| STM-004 | SM-004 | Payment Lifecycle | callback_received -> verified -> posted / duplicate_ignored | receive_callback/post_payment/duplicate_event | Valid signature, amount, idempotency key |
| STM-005 | SM-005 | Finance Clearance Lifecycle | blocked/conditional/clear | evaluate_clearance | Policy and invoice status valid |
| STM-006 | SM-006 | PMB Handover to Academic | not_started -> requested -> validating -> completed/failed/duplicate_ignored | request_handover/start_validation | Applicant ready, finance clear, curriculum/period valid |
| STM-007 | SM-007 | Student Academic Status | created/active/leave/dropout/graduated/alumni | activate/leave/graduate | Academic rule and required approval |
| STM-008 | SM-008 | KRS Lifecycle | draft -> submitted -> approved -> finalized/cancelled | submit/approve/finalize | Clearance, SKS, prerequisite, quota, schedule, PA approval |
| STM-009 | SM-009 | LMS Sync Lifecycle | pending -> processing -> synced/failed/ignored | sync/retry/duplicate_event | academic_class_id or krs_item_id valid |
| STM-010 | SM-010 | Assessment Attempt Lifecycle | assigned -> opened -> in_progress -> submitted -> scored/published/expired | open/submit/score/publish/expire | Session active, participant valid, time valid |
| STM-011 | SM-011 | Grade Finalization Lifecycle | source_imported -> reviewed -> finalized/corrected | review/finalize/correct | Source grade valid, reason for correction |

# 13. Regression Smoke Test

| Smoke ID | Module | Critical Path | Priority | Frequency |
| --- | --- | --- | --- | --- |
| SMOKE-001 | Core | Login, role selection, token validation | P0 | Every build |
| SMOKE-002 | Referensi | Academic year, period, prodi master active | P0 | Every release candidate |
| SMOKE-003 | PMB | Applicant registration, document upload, LoA guard | P0 | Every release candidate |
| SMOKE-004 | Finance | Invoice create, payment posted, duplicate callback ignored | P0 | Every build for finance changes |
| SMOKE-005 | Academic | Generate NIM, KRS finalization, grade finalization | P0 | Every release candidate |
| SMOKE-006 | LMS | Class sync and enrollment sync | P0 | Every LMS/Academic change |
| SMOKE-007 | Assessment | Attempt open, submit, score, publish result | P0 | Every Assessment change |
| SMOKE-008 | Portal | Dashboard, notification, shortcut by role | P1 | Every release candidate |
| SMOKE-009 | RBAC | Unauthorized endpoint direct access rejected | P0 | Every build |
| SMOKE-010 | Audit | Critical action writes audit and correlation id | P0 | Every release candidate |

# 14. Negative Case Minimum

| ID | Area | Expected Negative Control |
| --- | --- | --- |
| NEG-001 | Auth | Protected endpoint tanpa token ditolak dengan AUTH_REQUIRED/401. |
| NEG-002 | RBAC | User tanpa permission mencoba action/endpoint langsung dan ditolak dengan FORBIDDEN/403. |
| NEG-003 | Scope | Admin Prodi membuka data prodi lain melalui URL/API dan ditolak dengan SCOPE_DENIED/403. |
| NEG-004 | Payment | Callback dengan invalid signature ditolak dan tidak membuat payment. |
| NEG-005 | Idempotency | Request payment/handover/sync dengan key sama tidak membuat record ganda. |
| NEG-006 | KRS | Mahasiswa dengan clearance blocked tidak dapat finalisasi KRS. |
| NEG-007 | LMS | LMS tidak boleh membuat kelas mandiri tanpa academic_class_id. |
| NEG-008 | Assessment | Attempt expired tidak dapat submit jawaban. |
| NEG-009 | Grade | Koreksi nilai final tanpa reason ditolak. |
| NEG-010 | Portal | Portal tidak dapat mengubah data bisnis sumber secara langsung. |
| NEG-011 | Event Contract | Publish event tanpa event_key ditolak dengan EVENT_SCHEMA_INVALID. |
| NEG-012 | Event Contract | Consumer menerima event_version tidak didukung dan memasukkan event ke DLQ. |
| NEG-013 | Event Contract | Replay DLQ tanpa reason/permission ditolak dan tercatat audit. |
| NEG-014 | Event Contract | Payload hash berbeda untuk event_key sama ditolak sebagai conflict payload. |
| NEG-015 | Event Contract | Snapshot stale melebihi SLA tidak boleh dipakai untuk finalisasi report. |

# 15. NFR Baseline Test

| ID | Area | Minimum Test Focus |
| --- | --- | --- |
| NFR-001 | Security | Token, active role, service token, endpoint scope, direct API access, and sensitive action audit. |
| NFR-002 | Observability | Every API error has code, message, details when allowed, trace_id/correlation_id, timestamp. |
| NFR-003 | Performance | Critical list and dashboard pages load within agreed threshold on seeded data. |
| NFR-004 | Accessibility | Main flows support keyboard navigation, clear labels, specific validation messages, and responsive mobile layout. |
| NFR-005 | Reliability | Retry on integration does not create duplicate record and is traceable in integration log. |
| NFR-006 | Backup/Rollback | Migration rehearsal, backup verification, rollback rehearsal, and post-migration validation pass. |
| NFR-007 | Event Observability | event_lag_seconds, outbox_pending_count, inbox_pending_count, retry_count_by_event, dlq_count_by_consumer, dan reconciliation_mismatch_count tersedia. |
| NFR-008 | Event Resilience | Event retry tidak menyebabkan duplicate record pada payment, handover, student, enrollment, grade input, dan notification. |
| NFR-009 | Event Auditability | Publish, consume, retry, replay, DLQ, dan reconciliation wajib punya actor/system, timestamp, trace_id/correlation_id, dan audit evidence. |

# 16. UAT Sign-Off Control

| Area/Modul | Sign-off Role | Status | Catatan |
| --- | --- | --- | --- |
| Core | Owner Modul | Not Signed |  |
| Referensi | Owner Modul | Not Signed |  |
| CRM | Owner Modul | Not Signed |  |
| PMB | Owner Modul | Not Signed |  |
| Finance | Owner Modul | Not Signed |  |
| Akademik | Owner Modul | Not Signed |  |
| HRIS/SDM | Owner Modul | Not Signed |  |
| LMS | Owner Modul | Not Signed |  |
| Assessment | Owner Modul | Not Signed |  |
| Portal | Owner Modul | Not Signed |  |
| Cross-Module UAT | Product Owner | Not Signed |  |
| QA Completion | QA/UAT Lead | Not Signed |  |
| Technical Readiness | Technical Lead | Not Signed |  |
| Data/Migration Readiness | DBA | Not Signed |  |

# 17. Go/No-Go Recommendation Rule

Release dapat masuk go-live atau pilot hanya jika seluruh P0 UAT dan SIT pass, tidak ada defect Sev-1 dan Sev-2 terbuka, RBAC/scope test pass, migration dry-run pass, smoke test pass, dan sign-off minimal Product Owner, QA/UAT Lead, Technical Lead, DBA, serta Owner Modul terkait sudah lengkap.

Release harus ditahan jika ditemukan duplikasi payment, duplikasi student/NIM, kebocoran data lintas scope, bypass state machine, migration failure tanpa rollback aman, atau proses akademik/keuangan P0 tidak dapat dijalankan.

# 18. Lampiran

Lampiran utama untuk eksekusi test adalah workbook Excel UAT_QA_Test_Plan_Detail_ERP_UNSIA.xlsx. Workbook tersebut memuat Dashboard, Test Plan, UAT Scenarios, QA Test Cases, Integration Tests, RBAC Scope Tests, State Machine Tests, Regression Smoke, Data Setup, Defect Severity, UAT Sign-off, dan Source Basis.

# Appendix A - Event Contract QA/UAT Addendum

Appendix ini menambahkan standar pengujian Event Contract untuk memastikan event-driven integration berjalan idempotent, observable, resilient, dan dapat direkonsiliasi. Skenario ini berlaku untuk Core, Referensi, CRM, PMB, Finance, Academic, HRIS, LMS, Assessment, Portal, dan Reporting.

## A.1 Event Contract Quality Gate

| ID | Quality Gate |
| --- | --- |
| EC-QG-001 | Setiap event P0 memiliki event_name, event_version, event_key, publisher, consumer, payload schema, dan error contract. |
| EC-QG-002 | Outbox event hanya publish setelah transaksi lokal commit. |
| EC-QG-003 | Consumer mencatat event_key pada inbox sebelum atau saat memproses event. |
| EC-QG-004 | Duplicate event tidak membuat applicant, student, payment, enrollment, grade input, atau notification ganda. |
| EC-QG-005 | Retry, DLQ, replay, dan reconciliation memiliki audit evidence. |
| EC-QG-006 | Dashboard/read model menampilkan synced_at/refreshed_at dan source status. |

## A.2 Event Contract Test Data

| Area | Minimum Test Data |
| --- | --- |
| Payment event | invoice valid, payment valid, provider_event_id unik, duplicate callback payload. |
| Handover event | applicant accepted, LoA issued, payment policy valid, duplicate handover request. |
| KRS event | student active, clearance clear/blocked, KRS finalized, duplicate krs_approved event. |
| LMS event | academic_class_id, krs_item_id, enrollment projection, consumer down simulation. |
| Assessment event | attempt_id, participant_ref_id, score final, duplicate result publish. |
| Reconciliation data | source value berbeda dengan snapshot value, stale synced_at, missing source, missing snapshot. |

## A.3 Evidence Minimum untuk Event Test

| Evidence | Minimum Bukti |
| --- | --- |
| Outbox evidence | screenshot/query outbox_events berisi event_key, payload_hash, status, published_at. |
| Inbox evidence | screenshot/query inbox_events berisi event_key, consumer_module, status, processed_at. |
| Duplicate evidence | bukti jumlah record business tidak bertambah saat event sama dikirim ulang. |
| Retry evidence | retry_count, next_retry_at, last_error, dan status RETRYING. |
| DLQ evidence | dead_letter_at, last_error, raw_payload/payload, dan status DLQ. |
| Replay evidence | reason replay, actor/role, audit log, dan hasil replay. |
| Reconciliation evidence | mismatch report, source_value, snapshot_value, status correction. |

## A.4 Template Test Case Event Contract

Test ID :

Event Name :

Publisher :

Consumer :

Pre-condition :

Test Step :

Expected Result :

Validation Query :

Evidence Required :

Severity if Failed:
