---
title: "FSD Per Modul UNSIA"
source_file: "FSD_Per_Modul_UNSIA_v1_0_Event_Contract_Updated.docx"
format: markdown
---

# FSD Per Modul UNSIA

UNSIA

FUNCTIONAL SPECIFICATION DOCUMENT

FSD Per Modul - ERP Pendidikan / SIAKAD Terintegrasi UNSIA

v1.0 Draft | 18 Juni 2026

| Item | Isi |
| --- | --- |
| Dokumen | Functional Specification Document Per Modul |
| Produk | ERP Pendidikan / SIAKAD Terintegrasi UNSIA |
| Versi | v1.0.1 Updated Draft |
| Basis Penyusunan | PRD Global UNSIA v6.5.1 Event Contract Updated, BRD UNSIA v1.1.1 Event Contract Updated, dan revisi arsitektur database fisik per modul. |
| Tujuan | Mendefinisikan spesifikasi fungsi, menu, form, aksi, validasi, output, integrasi, audit, event contract functional standard, degraded mode behavior, dan UAT starter untuk setiap modul. |
| Status | Updated draft. Ditambahkan standar fungsional Event Contract agar FSD selaras dengan PRD dan BRD revisi distributed modular database. |

# 1. Kontrol Dokumen

| Versi | Tanggal | Status | Catatan |
| --- | --- | --- | --- |
| v1.0 | 18 Juni 2026 | Draft Awal | FSD per modul disusun berdasarkan PRD dan BRD yang telah dibuat sebelumnya. |
| v1.0.1 | 22 Juni 2026 | Updated Draft | Penambahan Event Contract Functional Standard, UI behavior untuk snapshot/read model, degraded mode, event log, retry, DLQ, reconciliation, observability, dan UAT event contract. |

| Peran | Tanggung Jawab |
| --- | --- |
| Product Owner | Menyetujui scope fungsi dan prioritas. |
| System Analyst | Menyusun spesifikasi fungsi, menu, validasi, dan alur modul. |
| Owner Modul | Memvalidasi kebutuhan operasional dan aturan bisnis. |
| Technical Lead | Menurunkan FSD menjadi desain teknis, API, dan task development. |
| DBA | Menurunkan field dan data ownership menjadi ERD/DBML dan data dictionary. |
| QA/UAT Lead | Menurunkan FSD menjadi test case dan UAT scenario. |

# 2. Prinsip FSD

FSD menjelaskan fungsi sistem, bukan kode program dan bukan DDL database final.

Setiap fungsi harus mengacu pada role, data scope, validasi, audit, dan output yang jelas.

Modul tidak boleh melakukan write langsung ke domain lain. Semua integrasi memakai API, event, service token, atau read model resmi.

Proses kritis wajib idempotent, terutama payment callback, PMB handover, generate NIM, class sync, enrollment sync, dan grade sync.

Status bisnis kritis harus memakai status code standar dan memiliki histori jika berubah.

Setiap halaman list utama wajib mendukung pagination, filter, sorting, dan search dengan index yang sesuai.

# 3. Global UI dan Functional Standard

| Area | Standar Fungsi |
| --- | --- |
| List Page | Pagination, sorting, search, filter, column visibility opsional, export jika role diizinkan, empty state, loading state. |
| Form Page | Save draft jika proses panjang, mandatory marker, validation message spesifik, confirmation untuk aksi sensitif. |
| Detail Page | Ringkasan data utama, status, histori perubahan, related records, action button sesuai permission. |
| Status Change | Wajib memiliki transition rule, actor, timestamp, reason/note jika sensitif, old status, new status. |
| File Upload | Allowed type, max size, virus scan jika tersedia, signed URL, versioning/reupload, rejection reason. |
| Audit | Aksi create, update, delete/deactivate, approve/reject, issue LoA, handover, payment verify, grade correction, impersonation wajib audit. |
| Error Handling | User message jelas, error code, trace_id, timestamp, detail teknis hanya untuk admin/developer. |
| Security | Backend harus validasi token, permission, active role, and data scope untuk setiap endpoint protected. |
| Event Sync Status | Halaman yang memakai snapshot/read model wajib menampilkan synced_at/refreshed_at, source module, dan status data bila tidak real-time. |
| Degraded Mode UI | Jika dependency down, UI wajib menampilkan status layanan, membatasi aksi yang berisiko, dan memberi label pending_review untuk proses yang perlu verifikasi ulang. |
| Integration Log UI | Halaman teknis wajib menampilkan event_key, idempotency_key, trace_id, status retry, last_error, processed_at, dan action replay sesuai permission. |
| Reconciliation UI | Mismatch antara source of truth dan snapshot/read model wajib tampil dalam daftar rekonsiliasi dengan status open, corrected, ignored, atau pending_review. |

# 4. Role dan Menu Coverage Global

| Role | Modul/Menu Utama | Scope Data |
| --- | --- | --- |
| Super Admin/Admin BPPTI | Core, application registry, role, permission, audit, service token | Global sesuai kewenangan |
| Admin Referensi | Master data umum, status code, payment component, document type | Global referensi |
| Admin CRM/Marketing | Campaign, lead, follow-up, conversion, commission | CRM domain |
| Agent/Mitra | Lead/referral milik sendiri | Agent scope |
| Pendaftar | PMB public, biodata, document upload, invoice status, LoA | Self scope |
| Admin PMB | Gelombang, applicant, dokumen, seleksi, LoA, handover | PMB domain |
| Admin Finance | Invoice, payment, verification, clearance, scholarship, report | Finance domain |
| Admin Akademik Biro | Calendar, student, NIM, class, KRS, grade, KHS, transcript, alumni | Academic global |
| Admin Akademik Prodi/Kaprodi | Curriculum, course, class, monitoring prodi | study_program_id |
| Dosen | LMS class, session, material, task, attendance, grade input | Assigned class |
| Dosen PA | KRS approval mahasiswa bimbingan | Advisor scope |
| Mahasiswa | KRS, LMS, invoice, grade, KHS, transcript | Self scope |
| Admin SDM | Employee, lecturer, homebase, status, BKD, certification | HRIS domain |
| Admin Assessment | Question bank, session, scoring, result | Assessment domain |
| Pimpinan | Executive dashboard and KPI | Read-only aggregate |

# 5. FSD Modul Core

| Item | Isi |
| --- | --- |
| Tujuan Modul | Identity authority, SSO, RBAC, audit, service token, idempotency, and integration control. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 5.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Login | Semua user | email/username, password, captcha opsional | submit login, forgot password | credential valid, user aktif, failed attempt policy | Core session aktif, audit login |
| Active Role Selector | User multi-role | user_id, role_id, scope | pilih role, switch role | role harus assigned dan aktif | token/session memakai active role |
| App Launcher | Semua user | application_id, menu code, shortcut | open app, search app | aplikasi tampil sesuai permission | navigasi ke modul |
| User Management | Super Admin, Admin BPPTI | person, username, email, status | create, update, activate, deactivate | email unik, person valid, reason jika deactivate | user tersedia untuk role assignment |
| Role and Permission | Super Admin | role, permission, menu, endpoint, action | create role, assign permission | permission tidak duplikat, perubahan audit | RBAC aktif di backend |
| Data Scope Management | Super Admin | role_id, scope_type, scope_value | assign scope | scope sesuai tipe role | admin prodi/dosen/agent/self terbatas |
| Service Token | Technical Admin | service name, token, expiry, allowed scope | generate, rotate, revoke | hanya service aktif, expiry valid | integrasi antar modul tervalidasi |
| Impersonation | Role tertentu | target user, reason, duration | start, stop | reason wajib, role allowed | audit impersonation |
| Audit Log Viewer | Auditor, Super Admin | actor, action, module, request_id, old/new value | filter, export | read-only, masked sensitive data | audit trail dapat ditelusuri |
| Idempotency and Integration Log | Technical Admin | event_key, idempotency_key, status, payload hash | view, retry marker | tidak edit transaksi | retry proses kritis aman |

## 5.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 6. FSD Modul Referensi

| Item | Isi |
| --- | --- |
| Tujuan Modul | Master data and status code authority for all modules. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 6.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Master Data Catalog | Admin Referensi | code, name, category, status | create, update, deactivate | code unik, status valid | dropdown lintas modul konsisten |
| Program Studi | Admin Referensi, Admin Akademik | study_program_code, name, faculty, level, status | create, update, deactivate | code unik, tidak hard delete jika dipakai | prodi aktif tersedia |
| Tahun Ajaran | Admin Akademik Biro | academic_year_code, start_date, end_date, status | create, activate, close | periode tanggal valid, hanya satu aktif sesuai policy | kalender operasional |
| Periode Akademik | Admin Akademik Biro | academic_year_id, period_type, start/end date, status | create, update, activate | wajib berada di bawah Tahun Ajaran | periode untuk PMB, invoice, kelas, KRS |
| Status Code | Admin Referensi | domain, code, label, sort_order, active | create, update, deactivate | domain dan code unik | status transaksi tidak string bebas |
| Komponen Pembayaran | Admin Finance | component_code, name, type, amount policy | create, update, deactivate | dipakai invoice tidak boleh hard delete | invoice item konsisten |
| Metode Pembayaran | Admin Finance | method_code, provider, active | create, update, deactivate | provider valid | payment method tersedia |
| Jenis Dokumen PMB | Admin PMB | document_code, name, required rule | create, update, deactivate | terkait jalur/prodi bila perlu | upload dokumen mengikuti rule |

## 6.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 7. FSD Modul CRM

| Item | Isi |
| --- | --- |
| Tujuan Modul | Lead acquisition and conversion to PMB applicant. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 7.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Campaign Management | Admin CRM | campaign_code, name, period, budget, source | create, update, close | periode valid, code unik | campaign aktif untuk lead |
| Lead Capture | Marketing, Admin CRM, Agent | name, phone, email, source, interest program | create lead | phone/email format valid, duplicate check | lead tercatat |
| Lead Pipeline | Marketing | status, owner, next action, priority | move status, assign owner | status transition valid | pipeline terkendali |
| Follow-up Activity | Marketing, Agent | channel, note, date, next_action | add follow-up | lead aktif, note wajib | histori follow-up |
| Agent/Referral | Admin CRM | agent data, referral code, commission rule | create, update, deactivate | referral code unik | agent scope aktif |
| Convert to Applicant | Admin CRM, Marketing | lead_id, target wave, study program | convert | lead qualified, idempotent, no duplicate applicant | applicant PMB terbentuk |
| Commission Record | Admin CRM | agent_id, applicant_id, event, amount | calculate, approve | rule terpenuhi | komisi siap review |
| CRM Dashboard | Admin CRM, Pimpinan | lead count, source, conversion, follow-up | filter, export | read-only untuk pimpinan | funnel akuisisi |

## 7.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 8. FSD Modul PMB

| Item | Isi |
| --- | --- |
| Tujuan Modul | Applicant lifecycle from registration to LoA and academic handover. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 8.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Public Registration | Pendaftar | gelombang, prodi, jalur, identitas awal | register, submit | gelombang aktif, target_entry_period_id valid | applicant draft/submitted |
| Applicant Dashboard | Pendaftar | status biodata, dokumen, invoice, seleksi, LoA | view status, continue process | self scope | status PMB transparan |
| Biodata Form | Pendaftar | personal data, address, education, family, financial profile | save draft, submit | mandatory field, format NIK/email/phone | biodata lengkap |
| Document Upload | Pendaftar | document_type, file, note | upload, replace | file type, size, required document | dokumen pending verification |
| Document Verification | Admin PMB | document_id, status, reason | verify, reject | reason wajib jika reject | dokumen verified/rejected |
| Selection/CBT Request | Admin PMB | applicant, session, context | assign CBT, read result | applicant eligible | score seleksi |
| Invoice Request | Admin PMB | applicant, fee component, due date | request invoice | komponen valid, applicant valid | invoice dari Finance |
| Payment Status Viewer | Pendaftar, Admin PMB | invoice, payment, status | view | read-only dari Finance | status pembayaran resmi |
| Re-registration | Pendaftar, Admin PMB | accepted applicant, payment status, documents | submit daftar ulang, validate | accepted dan payment policy valid | daftar ulang valid |
| LoA Issuance | Admin PMB | applicant, template, number | issue, download | seleksi, daftar ulang, payment valid | LoA issued |
| Handover to Academic | Admin PMB | applicant, target period, curriculum candidate | handover | ready for academic, idempotent | handover log dan proses NIM |

## 8.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 9. FSD Modul Finance

| Item | Isi |
| --- | --- |
| Tujuan Modul | Invoice, payment, verification, clearance, and financial reporting. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 9.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Invoice Management | Admin Finance | payer, component, amount, due date, period | create, update draft, issue | payer valid, component valid | invoice issued |
| Invoice Request Queue | Admin Finance | request source, applicant/student, component | approve, reject | request dari PMB/Academic valid | invoice resmi |
| Payment Gateway Callback Log | System, Admin Finance | provider event id, invoice, amount, signature, status | receive, validate | signature valid, idempotent | payment tercatat |
| Manual Payment Verification | Admin Finance | proof file, bank, amount, date, reason | approve, reject | bukti valid, reason wajib jika reject | payment verified/rejected |
| Receipt | Admin Finance, Payer | payment, receipt number, date | generate, download | payment valid | receipt resmi |
| Clearance Policy | Admin Finance | service, rule, threshold, status | create, update | policy aktif dan tidak konflik | policy layanan akademik |
| Student Clearance | Admin Finance, Academic | student, period, status, reason | evaluate, override conditional | override butuh reason | clear/blocked/conditional |
| Scholarship/Discount | Admin Finance | student, component, amount/percent, period | create, approve | approval sesuai policy | potongan tagihan |
| Installment/Dispensation | Admin Finance | student, schedule, amount, due date | create, approve | jadwal valid, reason wajib | status conditional |
| Finance Dashboard | Admin Finance, Pimpinan | issued, paid, overdue, aging, clearance | filter, export | read-only untuk pimpinan | monitoring keuangan |

## 9.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 10. FSD Modul Akademik

| Item | Isi |
| --- | --- |
| Tujuan Modul | Academic calendar, curriculum, student, NIM, KRS, class, grade, transcript, alumni. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 10.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Academic Calendar | Admin Akademik Biro | Tahun Ajaran, Periode Akademik, dates, status | create, activate, close | periode di bawah Tahun Ajaran | periode operasional |
| Curriculum Management | Admin Akademik Prodi/Biro | curriculum_year, study program, status | create, activate, archive | Tahun Kurikulum bukan Tahun Ajaran | curriculum_id aktif |
| Course Catalog | Admin Akademik Prodi | course code, name, credits, semester, prerequisite | create, update | code unik per curriculum/prodi | mata kuliah kurikulum |
| Class Offering | Admin Akademik | course, period, class, capacity, lecturer, schedule | create, update, close | academic_period_id wajib, dosen aktif | kelas kuliah |
| Generate NIM | Admin Akademik Biro | handover applicant, prodi, period, curriculum_id | generate | PMB/Finance valid, sequence lock | student dan NIM unik |
| Student Profile | Admin Akademik, Mahasiswa | student data, NIM, entry period, curriculum | view, update limited | source PMB/Core preserved | data mahasiswa |
| KRS Paket | Admin Akademik, Mahasiswa | student, package, period, classes | generate, approve | semester 1-2, package valid, clearance valid | KRS paket |
| KRS Mandiri | Mahasiswa | course, class, credits, period | select, drop, submit | semester 3+, SKS, prerequisite, quota, clash, clearance | KRS submitted |
| PA Approval | Dosen PA | student, KRS, note | approve, reject | advisor scope | KRS approved/rejected |
| Grade Input Review | Admin Akademik, Dosen | grade components, input source | review, finalize | rule nilai valid | nilai final |
| Grade Correction | Admin Akademik authorized | student, course, old grade, new grade, reason | submit correction | reason wajib, approval sesuai policy | grade history |
| KHS and Transcript | Mahasiswa, Admin Akademik | period, grades, IPS/IPK | view, publish, download | clearance policy valid | KHS/transkrip |
| Yudisium and Alumni | Admin Akademik | student, graduation status, date | process graduation | graduation final | alumni record |

## 10.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 11. FSD Modul HRIS

| Item | Isi |
| --- | --- |
| Tujuan Modul | Employee and lecturer source of truth for academic and LMS operations. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 11.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Employee Master | Admin SDM | person, employee number, unit, status | create, update, deactivate | person Core valid, employee number unik | pegawai aktif |
| Lecturer Profile | Admin SDM | NIDN/NIDK, homebase, academic rank, status | create, update | person valid, status aktif/nonaktif | dosen source of truth |
| Work Unit and Position | Admin SDM | unit, position, effective date | assign, update | effective date valid | riwayat jabatan |
| Employment Status | Admin SDM | active, inactive, leave, retired | change status | reason wajib | status dosen/pegawai |
| Lecturer Read Model | Admin Akademik, LMS | lecturer, status, homebase | view/search | read-only | plotting dosen valid |
| BKD Record | Admin SDM | lecturer, period, activity, workload | create, update | periode valid | data BKD |
| Certification and Performance | Admin SDM | certificate, performance score, period | create, update | dokumen/score valid | rekam jejak SDM |
| Payroll Source | Admin SDM, Finance | employee, salary base, status | prepare, export/read | scope sesuai release | data dasar payroll |

## 11.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 12. FSD Modul LMS

| Item | Isi |
| --- | --- |
| Tujuan Modul | Online learning delivery based on Academic class and valid KRS enrollment. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 12.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Class Sync Receiver | System/Admin LMS | academic_class_id, period, course, lecturer | receive sync | unique academic_class_id | kelas LMS |
| Enrollment Sync | System/Admin LMS | student, class, KRS status | receive, update | KRS approved | peserta LMS |
| Class Dashboard | Dosen, Mahasiswa | class info, sessions, materials, tasks | view | assigned class/self scope | kelas online |
| Session Management | Dosen | session date, topic, description | create, update, publish | dosen assigned | sesi belajar |
| Material Management | Dosen | title, file/link, description | upload, publish | file type/size valid | materi tersedia |
| Assignment | Dosen, Mahasiswa | title, due date, instruction, submission | create, submit, grade input | due date valid, student enrolled | tugas dan submission |
| Discussion | Dosen, Mahasiswa | topic, post, reply | create, reply, moderate | enrollment valid | diskusi kelas |
| Vicon Link | Dosen | provider, meeting link, schedule | create, update | dosen assigned | link vicon |
| Learning Attendance | Dosen | session, student, status | mark attendance | student enrolled | presensi LMS |
| Progress Tracking | Dosen, Mahasiswa | material read, task submitted, quiz status | view | self/assigned scope | progress mahasiswa |
| Grade Sync | System/Admin LMS | activity score, student, class, component | send to Academic | idempotent, final grade not overwritten | grade input akademik |

## 12.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 13. FSD Modul Assessment

| Item | Isi |
| --- | --- |
| Tujuan Modul | Reusable assessment engine for CBT, quiz, survey, scoring, and result API. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 13.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Question Bank | Admin Assessment, Dosen | question text, type, option, answer, difficulty | create, update draft | mandatory field, type valid | bank soal |
| Question Versioning | Admin Assessment | question_id, version, status | create new version | soal terpakai tidak diedit langsung | histori versi soal |
| Question Set | Admin Assessment | set name, context, questions, random rule | create, publish | question active | paket assessment |
| Assessment Session | Admin Assessment, PMB/LMS | context, schedule, duration, participant source | create, open, close | context valid, schedule valid | session CBT/quiz/survey |
| Participant Management | Admin Assessment | participant, source module, eligibility | assign, remove | participant valid | peserta assessment |
| Attempt Engine | Peserta | session, answers, timestamp | start, answer, submit | session open, participant eligible | attempt tercatat |
| Scoring | System/Admin Assessment | answer key, score rule, result | auto score, manual review | rule valid | score final assessment |
| Survey | Admin Assessment, Peserta | questionnaire, response | submit survey | context survey | survey result |
| Result API | System | context, participant, score, status | send result | consumer mapping valid | result ke PMB/LMS/Akademik |
| Attempt Report | Admin Assessment | session, participant, score distribution | filter, export | role authorized | laporan assessment |

## 13.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 14. FSD Modul Portal

| Item | Isi |
| --- | --- |
| Tujuan Modul | Role-based access layer, notification center, shortcut, and dashboard aggregation. |
| Tipe Spesifikasi | Functional design level. Detail API, database field final, and UI high fidelity akan diturunkan pada API Contract, ERD/DBML, Data Dictionary, dan UI Design. |
| Prinsip Modul | Setiap aksi mengikuti active role, permission, data scope, audit, validasi bisnis, dan source of truth domain. |

## 14.1 Functional Specification Table

| Screen/Fungsi | Role | Field Utama | Action | Validasi/Business Rule | Output/Integrasi |
| --- | --- | --- | --- | --- | --- |
| Role-based Dashboard | Semua user | active role, widget, shortcut | view | widget sesuai role dan scope | dashboard personal |
| Notification Center | Semua user | notification, source module, status | view, mark read | user scope valid | notifikasi terbaca |
| Shortcut Management | User/Admin | shortcut, module, menu | add, remove, reorder | permission valid | akses cepat |
| User Preference | User | language, display, notification setting | update | self scope | preferensi tersimpan |
| Applicant Portal | Pendaftar | PMB status, document status, invoice, LoA | view, continue action | self scope | status pendaftaran |
| Student Portal | Mahasiswa | invoice, clearance, KRS, class, grade, KHS | view, download | self scope dan clearance policy | layanan mahasiswa |
| Lecturer Portal | Dosen | assigned class, LMS shortcut, PA approval, grade tasks | view, navigate | assigned scope | workspace dosen |
| Executive Dashboard | Pimpinan | PMB funnel, finance, academic, LMS, HRIS KPI | view, drilldown limited | read-only | monitoring eksekutif |
| Activity Log | User/Admin | activity type, timestamp, source | view own/admin view | scope valid | jejak aktivitas portal |

## 14.2 Audit, Error, and UAT Notes

| Area | Catatan FSD |
| --- | --- |
| Audit | Aksi create/update/approve/reject/sync/status change wajib mencatat actor, active role, timestamp, request_id, old value, new value jika relevan. |
| Error | Setiap error menampilkan pesan spesifik untuk user dan trace_id untuk support teknis. |
| RBAC | Semua action button hanya tampil jika permission tersedia, tetapi backend tetap wajib memvalidasi permission dan scope. |
| UAT | UAT minimal menguji positive case, negative case, duplicate/idempotency case, access scope case, dan audit trail. |

# 15. Integration Functional Matrix

| Integrasi | Trigger | Data Utama | Validasi | Output |
| --- | --- | --- | --- | --- |
| Core ke semua modul | User login/access protected page | token, active role, permission, scope | token valid, session aktif | akses modul sesuai role |
| CRM ke PMB | Convert lead qualified | lead_id, applicant basic data, target wave | lead qualified, no duplicate | applicant PMB |
| PMB ke Finance | Request invoice | applicant_id, component, amount policy | applicant valid, component valid | invoice issued |
| Finance ke PMB/Academic/Portal | Payment status update | invoice_id, payment_status, clearance | callback/manual verification valid | status resmi terbaca |
| PMB ke Assessment | CBT assignment | applicant_id, session_id, context | eligible applicant | CBT participant |
| Assessment ke PMB/LMS/Academic | Result ready | participant_id, score, status, context | consumer mapping valid | result consumed |
| PMB ke Academic | Handover ready for academic | applicant, target period, prodi | accepted, LoA, payment policy valid | student/NIM process |
| Academic ke LMS | Class/enrollment sync | class, course, lecturer, students | class active, KRS approved | LMS class/enrollment |
| HRIS ke Academic/LMS | Lecturer read | lecturer_id, status, homebase | lecturer active | valid lecturer selection |
| LMS ke Academic | Grade input sync | activity score, component, student | idempotent, not final grade | grade input available |
| Modul ke Portal | Notification event | user, message, source module, link | user/scope valid | notification/read marker |
| Semua modul ke Event Broker | Status bisnis penting berubah | event_name, event_version, event_key, aggregate_id, payload | event_key unik, payload valid, outbox committed | event dipublish dan dapat dikonsumsi modul lain |
| Event Broker ke Consumer | Event diterima consumer | event_key, payload, source_event_key | consumer idempotent, event_version supported | inbox tercatat dan snapshot/read model diperbarui |
| Consumer ke Retry Queue | Proses event gagal sementara | event_key, retry_count, last_error | retry_count belum melewati batas | event diproses ulang sesuai retry policy |
| Retry Queue ke DLQ | Retry melewati batas | event_key, last_error, failed_at, raw_payload | gagal permanen atau versi tidak didukung | event masuk DLQ dan siap investigasi/replay |
| Reconciliation Job | Snapshot/read model berbeda dari source | source_ref_id, snapshot_value, source_value | mismatch terdeteksi | mismatch report dan correction job/pending_review |

# 16. Report Specification Summary

| Area | Report | Filter Minimal | Output |
| --- | --- | --- | --- |
| Core | User, role, audit, impersonation, integration log | date, role, module, actor | table, export authorized |
| Referensi | Master data status, academic period, study program | status, category, date | table, export |
| CRM | Lead funnel, campaign conversion, agent performance | date, campaign, source, agent | dashboard, table, export |
| PMB | Applicant funnel, document verification, payment PMB, LoA, handover | wave, prodi, status, period | dashboard, table, export |
| Finance | Invoice, payment, overdue, aging receivable, clearance | period, component, status, prodi | dashboard, table, export |
| Akademik | Student, KRS, class, grade, KHS, transcript, alumni | period, prodi, semester, status | dashboard, table, document output |
| HRIS | Lecturer active, homebase, workload, BKD, certification | unit, prodi, status, period | dashboard, table, export |
| LMS | Class activity, progress, attendance, assignment, grade input | period, class, lecturer, status | dashboard, table, export |
| Assessment | Question bank, attempt statistics, score distribution, survey result | context, session, date, participant | dashboard, table, export |
| Portal | Notification delivery, unread, user activity, executive dashboard | role, date, module | dashboard, table |

# 17. Open Items untuk Dokumen Lanjutan

| Dokumen Lanjutan | Yang Harus Diperinci |
| --- | --- |
| API Contract | Endpoint, method, request, response, error code, permission, service token, idempotency, retry policy. |
| ERD/DBML | Tabel final, relasi, FK, unique index, enum/status, audit/history table, migration strategy. |
| Data Dictionary | Nama field, tipe data, nullable, default, validation rule, source field, sample value. |
| State Machine | Status, transition, allowed actor, precondition, postcondition, reason, audit, notification. |
| RBAC Matrix | Role, menu, action, endpoint, data scope, read/write/delete/export permission. |
| UAT Scenario | Positive case, negative case, role-scope case, idempotency case, audit case, regression case. |
| UI Wireframe | Layout halaman, form grouping, user journey, responsive behavior, error state. |
| Event Contract Catalog | event_name, event_version, publisher, consumer, trigger bisnis, payload schema, idempotency rule, retry policy, DLQ, snapshot impact, reconciliation rule, security rule, observability, dan UAT scenario. |

# 18. Approval

| Peran | Nama | Status | Tanggal | Catatan |
| --- | --- | --- | --- | --- |
| Product Owner |  | Belum disetujui |  |  |
| System Analyst |  | Drafted |  |  |
| Technical Lead |  | Belum direview |  |  |
| DBA |  | Belum direview |  |  |
| QA/UAT Lead |  | Belum direview |  |  |
| Owner Core |  | Belum direview |  |  |
| Owner PMB |  | Belum direview |  |  |
| Owner Finance |  | Belum direview |  |  |
| Owner Akademik |  | Belum direview |  |  |

# Appendix A - Event Contract Functional Standard

Appendix ini menambahkan standar fungsional untuk Event Contract pada FSD Per Modul UNSIA. Fokusnya adalah perilaku layar, aksi pengguna, validasi, output, integrasi, audit, degraded mode, retry, DLQ, reconciliation, dan UAT. Standar ini tidak menggantikan API Contract atau Event Contract Catalog, tetapi menjadi acuan fungsional bagi tim analis, developer, QA, DBA, dan DevOps.

## A.1 Prinsip Fungsional Event Contract

| ID | Prinsip | Acceptance Criteria |
| --- | --- | --- |
| FSD-EC-001 | Setiap event penting harus terlihat pada integration log. | Admin teknis dapat menelusuri event_key, source, consumer, status, retry_count, dan last_error. |
| FSD-EC-002 | Halaman berbasis snapshot wajib menampilkan freshness data. | User melihat synced_at/refreshed_at dan label data mungkin belum real-time. |
| FSD-EC-003 | Aksi yang bergantung pada modul down harus masuk mode terbatas. | Sistem menahan aksi sebagai pending_review atau menampilkan read-only sesuai rule. |
| FSD-EC-004 | Duplicate event tidak boleh membuat data ganda. | Consumer memproses satu event_key satu kali. |
| FSD-EC-005 | DLQ dan reconciliation harus dapat dipantau. | Role teknis melihat daftar event gagal, mismatch, dan status perbaikan. |

## A.2 Event Functional Coverage per Modul

| Modul | Event Utama | Consumer | Dampak Fungsional |
| --- | --- | --- | --- |
| Core | core.person_updated, core.user_role_changed | Semua modul | Update person/user/role snapshot dan permission cache. |
| Referensi | reference.master_data_updated, reference.period_changed | PMB, Finance, Academic, HRIS, LMS, Assessment, Portal | Update reference snapshot untuk dropdown, status code, periode, dan prodi. |
| CRM | crm.lead_created, crm.lead_qualified | PMB, Portal, Reporting | Lead dapat dikonversi ke applicant secara idempotent. |
| PMB | pmb.applicant_created, pmb.document_verified, pmb.ready_for_academic | Finance, Assessment, Academic, Portal, Reporting | Invoice request, assessment assignment, handover, dan dashboard PMB. |
| Finance | finance.invoice_created, finance.payment_paid, finance.clearance_changed | PMB, Academic, LMS, Portal, Reporting | Payment status, clearance snapshot, dan notifikasi layanan akademik. |
| Akademik | academic.student_created, academic.krs_approved, academic.final_grade_published | PMB, LMS, Portal, Reporting | NIM, enrollment LMS, KHS, transkrip, dan dashboard akademik. |
| HRIS | hris.lecturer_status_changed | Academic, LMS, Portal | Validasi plotting dosen dan kelas LMS. |
| LMS | lms.progress_updated, lms.grade_input_submitted | Academic, Portal, Reporting | Progress belajar dan grade input ke Academic. |
| Assessment | assessment.result_calculated | PMB, LMS, Academic, Portal | Hasil CBT/quiz/survey dikirim ke context owner. |
| Portal | portal.notification_read | Portal | Read marker dan preferensi notifikasi. |

## A.3 Standar Field pada Layar Integration Log

| Field | Fungsi Tampilan |
| --- | --- |
| event_name | Nama event yang diproses. |
| event_version | Versi payload schema. |
| event_key | Kunci unik event untuk idempotency. |
| publisher_service | Service pengirim event. |
| consumer_service | Service penerima event. |
| aggregate_type dan aggregate_id | Objek bisnis utama yang berubah. |
| status | PENDING, PROCESSED, RETRYING, FAILED, DLQ, atau IGNORED_DUPLICATE. |
| retry_count | Jumlah retry yang sudah dilakukan. |
| last_error | Error terakhir jika proses gagal. |
| occurred_at, published_at, processed_at | Waktu kejadian, publish, dan consume event. |
| trace_id/correlation_id | ID pelacakan lintas service. |

## A.4 Screen Specification Tambahan

| Screen/Fungsi | Role | Action | Output |
| --- | --- | --- | --- |
| Outbox Event Viewer | Technical Admin, DevOps/SRE | Filter event keluar, lihat payload, status publish, retry marker. | Event keluar dapat ditelusuri. |
| Inbox Event Viewer | Technical Admin, DevOps/SRE | Filter event masuk, status proses, duplicate flag, processed_at. | Consumer processing transparan. |
| Retry Queue | DevOps/SRE | View retry_count, next_retry_at, last_error, manual retry jika diizinkan. | Event gagal sementara dapat diproses ulang. |
| Dead Letter Queue | DevOps/SRE | View DLQ, inspect payload, mark resolved, replay terbatas. | Event gagal permanen tidak hilang. |
| Reconciliation Monitor | Admin Modul, DevOps/SRE | View mismatch source vs snapshot, correction status, pending_review. | Selisih data dapat ditindaklanjuti. |
| Event Detail Page | Technical Admin, Auditor | Lihat event envelope, payload masked, consumer result, audit trail. | Audit event lengkap tersedia. |

## A.5 Functional Behavior saat Dependency Down

| Dependency Down | Perilaku Fungsional | UI/Output |
| --- | --- | --- |
| Finance down | PMB tetap menerima biodata/dokumen. Payment status memakai snapshot terakhir. | Label data tidak real-time dan tampilkan synced_at. Invoice baru masuk retry/pending. |
| Academic down | LMS tetap membuka kelas dan enrollment yang sudah tersinkron. | Class sync baru masuk retry queue. |
| LMS down | Academic tetap membuka kelas, KRS, dan final grade. | Sinkronisasi LMS masuk outbox dan retry. |
| Portal down | Modul sumber tetap melakukan transaksi. | Notifikasi diproses ulang setelah Portal pulih. |
| Core down | Token valid yang belum expired tetap dapat divalidasi memakai cache/public key. | Login baru dan role switching ditahan sampai Core pulih. |

## A.6 Validation dan Error Message Standard

| Error Code | Kondisi | User Message |
| --- | --- | --- |
| EVENT_SCHEMA_INVALID | Payload event tidak sesuai schema. | Event tidak dapat diproses karena format data tidak valid. Hubungi admin teknis. |
| EVENT_DUPLICATE | Event dengan event_key sama sudah diproses. | Data sudah diproses sebelumnya. Tidak ada perubahan baru. |
| EVENT_VERSION_UNSUPPORTED | Consumer belum mendukung versi event. | Versi event belum didukung oleh modul penerima. |
| SOURCE_REF_NOT_FOUND | External reference belum ditemukan. | Data referensi belum tersedia. Sistem akan mencoba ulang. |
| SNAPSHOT_STALE | Snapshot melewati batas freshness. | Data yang ditampilkan mungkin belum real-time. |
| RECONCILIATION_REQUIRED | Snapshot/read model berbeda dari source. | Data perlu rekonsiliasi sebelum finalisasi. |

## A.7 Audit dan Permission Event

| Action | Role | Aturan |
| --- | --- | --- |
| View Event Log | Technical Admin, Auditor | Read-only, payload sensitif dimasked. |
| Manual Retry | DevOps/SRE | Hanya untuk event gagal sementara. Wajib audit. |
| Replay DLQ | DevOps/SRE authorized | Butuh alasan dan approval sesuai policy. |
| Mark Resolved | DevOps/SRE, Owner Modul | Wajib mencatat reason. |
| Ignore Mismatch | Owner Modul authorized | Wajib alasan bisnis dan audit trail. |
| Export Event Log | Auditor, Super Admin | Masked sensitive data dan sesuai permission. |

## A.8 UAT Event Contract Functional Scenario

| Skenario UAT | Expected Result |
| --- | --- |
| Event valid dipublish. | Outbox event berubah menjadi published dan consumer menerima event. |
| Consumer berhasil memproses event. | Inbox mencatat event_key, status PROCESSED, dan snapshot/read model berubah. |
| Event yang sama dikirim dua kali. | Consumer memproses satu kali dan menandai duplicate sebagai ignored. |
| Consumer down. | Event masuk retry queue dan diproses setelah consumer pulih. |
| Event gagal permanen. | Event masuk DLQ dengan last_error dan raw_payload. |
| Replay DLQ dilakukan. | Replay tercatat audit dan tidak membuat data duplikat. |
| Snapshot stale. | UI menampilkan synced_at dan label data mungkin belum real-time. |
| Mismatch source vs snapshot. | Reconciliation monitor menampilkan mismatch report. |

## A.9 Template Tambahan untuk Setiap Fungsi Berbasis Event

Screen/Fungsi :
Role :
Trigger Bisnis :
Event Published :
Event Consumed :
Field Utama :
Action :
Validasi :
Output :
Snapshot Impact :
Degraded Behavior :
Idempotency Rule :
Audit :
Error Code :
UAT Scenario :
