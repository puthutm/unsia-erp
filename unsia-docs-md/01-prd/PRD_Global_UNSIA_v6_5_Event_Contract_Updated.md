---
title: "PRD Global UNSIA"
source_file: "PRD_Global_UNSIA_v6_5_Event_Contract_Updated.docx"
format: markdown
---

# PRD Global UNSIA

UNSIA

PRODUCT REQUIREMENT DOCUMENT

PRD Global - Distributed Modular Database Revision

ERP Pendidikan / SIAKAD Terintegrasi UNSIA

v6.5 Revised Draft | 19 Juni 2026

| Item | Isi |
| --- | --- |
| Dokumen | Product Requirement Document Global - Revisi Arsitektur Database Modular Terdistribusi |
| Produk | ERP Pendidikan / SIAKAD Terintegrasi UNSIA |
| Versi | v6.5 Revised Draft |
| Basis Revisi | PRD Global UNSIA v6.4 Detailed FULL, BRD Global v1.0, BRD Per Modul v1.0, dan keputusan desain database fisik per modul. |
| Status | Draft revisi untuk review Product Owner, System Analyst, Technical Lead, DBA, Security, DevOps, dan Owner Modul. |
| Keputusan Utama | Mengganti asumsi satu PostgreSQL utama dengan schema per domain menjadi database fisik berdiri sendiri per modul, tanpa cross-database foreign key dan tanpa direct cross-database join pada transaksi online. |

# 1. Kontrol Dokumen

| Versi | Tanggal | Status | Catatan Perubahan |
| --- | --- | --- | --- |
| v6.5 | 19 Juni 2026 | Revised Draft | Revisi arsitektur data menjadi physical database per modul, event-driven integration, snapshot/read model, outbox/inbox, degraded mode, dan failure isolation. |
| v6.4 | 17 Juni 2026 | Detailed PRD | Baseline detailed PRD dengan multi-repo dan satu PostgreSQL utama/schema per domain. |
| v6.5.1 | 22 Juni 2026 | Updated Draft | Penambahan Appendix B - Event Contract Standard untuk memperjelas event identity, payload schema, idempotency, retry, DLQ, snapshot/read model, reconciliation, security, error contract, observability, dan UAT event. |

| Peran | Tanggung Jawab pada Revisi v6.5 |
| --- | --- |
| Product Owner | Mengesahkan scope bisnis, prioritas release, dan batasan toleransi eventual consistency. |
| System Analyst | Menjaga traceability PRD ke BRD, API contract, event contract, state machine, UAT, dan release plan. |
| Technical Lead | Memastikan arsitektur multi-service, API composition, outbox/inbox worker, retry, dan deployment dapat diimplementasikan. |
| DBA | Menetapkan database boundary, internal FK, index, idempotency, backup/restore, RPO/RTO, replication, dan recovery per database. |
| Security/DevOps | Menentukan service authentication, secret management, network policy, observability, backup isolation, dan disaster recovery. |
| Owner Modul | Memvalidasi data ownership, snapshot minimum, event yang dibutuhkan, degraded mode, dan proses rekonsiliasi. |
| QA/UAT Lead | Menguji failure scenario, duplicate event, delayed event, retry, partial outage, data reconciliation, dan permission scope. |

# 2. Ringkasan Eksekutif Revisi

Revisi v6.5 menetapkan bahwa ERP Pendidikan / SIAKAD Terintegrasi UNSIA tidak lagi diasumsikan berjalan pada satu database PostgreSQL utama dengan schema per domain. Produk direvisi menjadi ekosistem modular dengan database fisik berdiri sendiri per modul: Core, Referensi, CRM, PMB, Finance, Akademik, HRIS, LMS, Assessment, dan Portal.

Tujuan utama perubahan ini adalah failure isolation. Jika satu database modul mengalami gangguan, modul lain yang tidak bergantung langsung pada data real-time modul tersebut harus tetap dapat membaca dan memproses data lokalnya. Gangguan Finance tidak boleh menghentikan input biodata PMB, gangguan LMS tidak boleh menghentikan KRS Akademik, dan gangguan Portal tidak boleh menghentikan transaksi operasional modul sumber.

Konsekuensinya, integritas lintas modul tidak lagi dijaga dengan foreign key lintas database. Integritas lintas modul dijaga dengan API contract, event contract, idempotency key, outbox/inbox pattern, snapshot lokal, read model, reconciliation job, dan data ownership yang tegas.

PRD ini tetap mempertahankan prinsip produk v6.4: single identity melalui Core, source of truth tunggal per domain, backend-enforced RBAC, academic calendar first, curriculum preserved per student, finance clearance terintegrasi, audit, dan idempotency. Revisi ini hanya mengubah model persistence dan pola integrasi agar lebih tahan terhadap partial outage.

| Keputusan DBA: Agar klaim “satu DB mati tidak mempengaruhi DB yang tidak berhubungan” valid secara operasional, database per modul harus ditempatkan minimal sebagai database instance/cluster terpisah untuk modul kritikal. Jika semua database masih berada di satu server/cluster fisik yang sama, maka kegagalan server/cluster tetap dapat menjatuhkan seluruh modul. |
| --- |

# 3. Latar Belakang Revisi Arsitektur

Baseline v6.4 sudah menetapkan domain ownership dan modul tidak boleh mengubah langsung data domain lain. Namun asumsi satu database utama masih memiliki satu availability boundary. Desain tersebut kuat untuk konsistensi dan reporting lintas modul, tetapi kurang optimal bila target operasional adalah isolasi kegagalan antar database.

| Masalah pada Desain Satu DB Utama | Risiko | Arah Revisi v6.5 |
| --- | --- | --- |
| Satu database utama menjadi shared persistence boundary. | Gangguan database/cluster dapat mempengaruhi semua modul sekaligus. | Database fisik dipisah per modul dengan backup, restore, dan recovery masing-masing. |
| Cross-schema join mudah dilakukan. | Aplikasi berisiko bergantung pada tabel domain lain dan melanggar ownership. | Online transaction dilarang melakukan direct cross-database join; gunakan API/read model. |
| Cross-schema FK menggoda dipakai untuk semua relasi. | Modul menjadi tightly coupled dan sulit dipisah. | FK hanya internal database. Relasi lintas modul memakai external_ref_id dan contract validation. |
| Reporting lintas modul bisa langsung dari DB produksi. | Query analitik berat dapat mengganggu transaksi operasional. | Reporting lintas modul melalui warehouse, data mart, atau portal read model. |
| Retry integrasi bisa membuat data ganda bila tidak seragam. | Duplicate applicant, duplicate payment, duplicate grade input. | Semua proses kritis wajib punya idempotency key dan event_key deterministik. |

# 4. Visi, Tujuan, dan Prinsip Produk v6.5

Visi produk tetap membangun backbone operasional kampus dari peminat sampai alumni. Revisi v6.5 menambahkan prinsip resilience dan data independence agar sistem tetap dapat beroperasi secara parsial ketika salah satu modul mengalami gangguan.

| Tujuan Produk | Penjelasan Revisi | Indikator Keberhasilan |
| --- | --- | --- |
| Single identity dan SSO | Core tetap menjadi authority akun, role, permission, token, app registry, dan service client. | Tidak ada credential table di luar Core; modul dapat memvalidasi token yang belum expired menggunakan cache/public key. |
| Data ownership jelas | Setiap database modul hanya menjadi source of truth untuk domainnya sendiri. | Tidak ada write langsung ke database modul lain; semua write lintas modul melalui API/event. |
| Failure isolation | Gangguan database satu modul tidak menghentikan data dan proses modul lain yang tidak bergantung langsung. | UAT partial outage membuktikan modul lain tetap berjalan dengan snapshot terakhir atau degraded mode. |
| Eventual consistency terkendali | Data lintas modul disinkronkan melalui event dan read model dengan staleness yang terlihat. | Setiap read model memiliki source_event_key, synced_at, dan status rekonsiliasi. |
| Finance clearance terintegrasi | Finance tetap source of truth clearance, sementara Academic/PMB menyimpan clearance snapshot untuk operasional terbatas. | KRS/ujian/KHS/transkrip/wisuda mematuhi policy clearance dengan fallback pending_review saat Finance tidak tersedia. |
| Audit dan idempotency | Setiap database memiliki audit lokal, idempotency lokal, outbox, dan inbox. | Duplicate event, retry API, dan callback berulang tidak membuat data ganda. |

| Prinsip Produk v6.5 | Makna Praktis |
| --- | --- |
| Physical database per module | Core, Referensi, CRM, PMB, Finance, Akademik, HRIS, LMS, Assessment, dan Portal memiliki database fisik sendiri. |
| No cross-database foreign key | Relasi lintas modul memakai external_ref_id, bukan FK database. FK hanya boleh di dalam database modul yang sama. |
| No direct cross-database join for OLTP | Transaksi online tidak boleh melakukan join langsung ke database modul lain. Data lintas modul diperoleh melalui API composition, snapshot, read model, atau event projection. |
| Snapshot is not source of truth | Snapshot dipakai agar modul tetap operasional saat dependency down, tetapi kebenaran final tetap di modul pemilik. |
| Outbox/inbox mandatory | Setiap perubahan penting diterbitkan ke outbox dan dikonsumsi via inbox secara idempotent. |
| Degraded mode by design | Setiap modul wajib memiliki perilaku jelas ketika dependency utama tidak tersedia. |
| Reconciliation before report finalization | Laporan final lintas modul harus melewati rekonsiliasi, bukan hanya membaca snapshot yang mungkin stale. |

# 5. Scope Produk Global Revisi

| Modul | Database Fisik | Source of Truth | Data Lintas Modul yang Boleh Disimpan Lokal |
| --- | --- | --- | --- |
| Core | core_db | persons, users, roles, permissions, sessions, service clients, app registry | Scope_ref_id dari Referensi/HRIS untuk role assignment; tidak menjadi pemilik master prodi/unit. |
| Referensi | reference_db | regions, religions, study_programs, academic_years, academic_periods, payment components, status codes | User_ref_id untuk audit; tidak menyimpan credential. |
| CRM | crm_db | campaigns, leads, agents, referrals, follow-ups, commissions | person snapshot dari Core; applicant_ref_id dari PMB setelah convert. |
| PMB | pmb_db | applicants, applicant biodata, documents, selection status, re-registration, LoA, handover | person snapshot, reference snapshot, invoice/payment status snapshot, assessment result reference, student_ref_id setelah handover. |
| Finance | finance_db | invoices, invoice items, payments, receipts, callbacks, clearances, journal entries | customer snapshot dari PMB/Akademik/Core; academic_period_ref_id dari Referensi. |
| Akademik | academic_db | students, NIM, curriculums, courses, course offerings, KRS, final grades, KHS, transcripts, alumni | person snapshot, applicant_ref_id, lecturer_ref_id, clearance snapshot, reference snapshot. |
| HRIS | hris_db | employees, lecturers, positions, work units, homebase, BKD, payroll source | person snapshot dari Core; study_program_ref_id dari Referensi. |
| LMS | lms_db | online classes, LMS enrollment, sessions, materials, assignments, attendance, progress, LMS grade input | academic class snapshot, student snapshot, lecturer snapshot, assessment_session_ref_id. |
| Assessment | assessment_db | question banks, question versions, assessment sessions, attempts, answers, scoring results | participant snapshot dari PMB/Akademik/Core; context_ref_id dari consumer. |
| Portal | portal_db | notifications, read markers, preferences, shortcuts, dashboard read models | user/role snapshot, notification events, aggregated dashboard payload. |

## 5.1 Out of Scope Revisi

Distributed transaction two-phase commit lintas database tidak menjadi requirement. Konsistensi lintas modul memakai saga/eventual consistency.

Direct database link, FDW, atau cross-database join untuk transaksi online tidak menjadi pola resmi aplikasi.

Data warehouse final dan BI enterprise penuh dapat menjadi dokumen turunan, bukan bagian detail DDL PRD ini.

Integrasi PDDIKTI/NeoFeeder full automation tetap bukan MVP pertama, tetapi desain data akademik harus siap untuk mapping dan rekonsiliasi.

Mobile app native penuh tetap di luar MVP awal; portal responsive/mobile web menjadi prioritas.

# 6. Definisi Kunci Tambahan v6.5

| Istilah | Definisi Produk | Contoh |
| --- | --- | --- |
| Physical module database | Database fisik yang dimiliki satu modul dan menjadi batas ownership serta recovery. | academic_db hanya dimiliki modul Akademik. |
| External reference / *_ref_id | UUID milik domain lain yang disimpan tanpa FK database lintas modul. | pmb_db.applicants.person_ref_id menunjuk core_db.persons.id. |
| Snapshot | Salinan ringkas data domain lain untuk kebutuhan tampilan/proses lokal saat dependency down. | lms_db.student_snapshots berisi nim dan nama mahasiswa. |
| Read model | Projection lokal hasil event untuk query cepat atau dashboard. | pmb_db.applicant_invoice_statuses dari event Finance. |
| Outbox event | Event yang ditulis dalam transaksi lokal database pemilik sebelum dipublish ke broker. | finance.payment_paid. |
| Inbox event | Catatan event masuk yang sudah diterima/diproses consumer untuk mencegah duplicate processing. | academic_db.inbox_events menyimpan finance.clearance_changed. |
| Eventual consistency | Kondisi data antar modul menjadi konsisten setelah event diproses, bukan harus serentak dalam satu transaksi. | Payment paid di Finance muncul di PMB beberapa detik kemudian. |
| Degraded mode | Mode operasi terbatas ketika dependency tidak tersedia. | Academic menahan finalisasi KRS sebagai pending_review jika clearance real-time tidak tersedia. |
| Reconciliation job | Proses berkala untuk membandingkan snapshot/read model dengan source of truth. | PMB mencocokkan applicant invoice status dengan Finance setiap malam. |

# 7. Stakeholder dan Kebutuhan Tambahan

| Persona/Role | Kebutuhan Tambahan Akibat Database Terdistribusi |
| --- | --- |
| Pendaftar | Tetap dapat mengisi biodata dan dokumen meskipun Finance sementara tidak tersedia; status pembayaran menampilkan status terakhir dan label waktu sinkronisasi. |
| Mahasiswa | Tetap dapat melihat KRS/nilai terakhir meskipun modul Finance atau LMS sedang gangguan, dengan indikator data mungkin belum real-time. |
| Dosen | Tetap dapat mengakses kelas LMS yang sudah tersinkron meskipun Academic sementara tidak tersedia. |
| Admin Modul | Memiliki informasi jelas apakah data yang tampil real-time dari source API atau snapshot lokal. |
| Pimpinan | Dashboard agregat menampilkan timestamp refresh dan status kesehatan sumber data. |
| DBA/DevOps | Membutuhkan dashboard health per database, replication status, backup status, lag event, dan retry queue. |

# 8. Proses Bisnis End-to-End dan Mode Integrasi

| Tahap | Modul Pemilik | Database Pemilik | Integrasi Keluar | Fallback Jika Dependency Down |
| --- | --- | --- | --- | --- |
| Lead/Peminat | CRM | crm_db | crm.lead_created, crm.lead_qualified | Lead tetap tercatat lokal; handover ke PMB masuk retry queue. |
| Applicant/Pendaftar | PMB | pmb_db | pmb.applicant_created, pmb.document_verified | PMB tetap menerima biodata/dokumen; invoice creation retry jika Finance down. |
| Invoice dan Payment | Finance | finance_db | finance.invoice_created, finance.payment_paid, finance.clearance_changed | Finance outage hanya mempengaruhi billing/payment/clearance real-time; modul lain memakai snapshot terakhir. |
| Seleksi/CBT | Assessment + PMB | assessment_db + pmb_db | assessment.result_calculated, pmb.selection_decided | Jika Assessment down, jadwal CBT/quiz tertahan; PMB non-CBT tetap jalan. |
| LoA dan Handover | PMB | pmb_db | pmb.ready_for_academic, pmb.handover_requested | Jika Academic down, handover masuk pending retry, applicant tetap tersimpan. |
| Generate NIM | Akademik | academic_db | academic.student_created | Jika PMB down, Akademik tetap memproses request yang sudah diterima dan idempotent. |
| KRS dan Kelas | Akademik | academic_db | academic.krs_approved, academic.class_opened | KRS baru butuh clearance policy; LMS sync retry jika LMS down. |
| Pembelajaran Online | LMS | lms_db | lms.progress_updated, lms.grade_input_submitted | Kelas yang sudah tersinkron tetap berjalan meskipun Academic down. |
| Nilai Final/KHS/Transkrip | Akademik | academic_db | academic.final_grade_published, academic.khs_issued | Nilai final tidak dipindah ke LMS; input LMS/Assessment masuk sebagai grade input. |
| Dashboard/Notifikasi | Portal | portal_db | portal.notification_created | Jika Portal down, modul sumber tetap jalan; notifikasi diproses ulang setelah pulih. |

# 9. Arsitektur Produk dan Asumsi Implementasi v6.5

| Aspek | Keputusan Produk v6.5 | Implikasi |
| --- | --- | --- |
| Aplikasi | Multi-repo/multi-service per modul. | Deployment dapat bertahap dan failure modul lebih terisolasi. |
| Database | Physical database per modul. | Migration, backup, restore, RPO/RTO, indexing, dan tuning dilakukan per database. |
| Availability boundary | Minimal per database modul; modul kritikal direkomendasikan berada pada instance/cluster terpisah. | DB mati pada satu modul tidak otomatis menjatuhkan DB modul lain jika infrastruktur fisiknya terpisah. |
| Auth | Core tetap identity authority dengan OIDC/JWT/service client. | Token validasi dapat dilakukan lokal memakai cached JWKS/public key. Login baru bergantung pada Core. |
| Integrasi sync | REST/gRPC API + service client + circuit breaker + timeout. | Dipakai untuk command/query real-time yang memang membutuhkan data terbaru. |
| Integrasi async | Event broker + transactional outbox/inbox. | Dipakai untuk sinkronisasi status, snapshot, read model, dan notification. |
| Data reference | External_ref_id dan snapshot. | Tidak ada cross-database FK; validasi via API/event. |
| Reporting | Portal read model untuk dashboard operasional; warehouse/data mart untuk laporan lintas modul final. | Query analitik tidak membebani database transaksi. |
| Audit | Audit lokal per database + optional audit aggregation. | Aksi sensitif tetap terlacak walau database modul lain down. |
| Idempotency | Idempotency key lokal per modul dan event_key deterministik lintas modul. | Retry aman, duplicate event tidak menghasilkan data ganda. |

Frontend/API Gateway
 -> Core Service / core_db
 -> PMB Service / pmb_db
 -> Finance Service / finance_db
 -> Academic Service / academic_db
 -> LMS Service / lms_db

Event Broker
 <- outbox_events dari setiap DB
 -> inbox_events consumer
 -> snapshot/read model lokal

Reporting Warehouse/Data Mart
 <- CDC/Event dari semua DB
 -> dashboard/laporan lintas modul

# 10. Source of Truth dan Ownership Data v6.5

| Domain Data | Source of Truth | Aturan Produk v6.5 |
| --- | --- | --- |
| Identitas orang | core_db.persons | Modul lain menyimpan person_ref_id dan snapshot minimum; tidak membuat identitas paralel. |
| Akun, role, permission | core_db.users, roles, permissions | Tidak ada password/session/role authority di modul lain. Permission cache boleh ada untuk degraded read. |
| Master data umum | reference_db.* | Modul lain menyimpan ref_id/code/name snapshot; perubahan master dipublish sebagai event. |
| Lead/peminat | crm_db.* | PMB menerima lead_ref_id saat convert; CRM tetap pemilik histori lead dan campaign. |
| Applicant PMB | pmb_db.* | Finance/Akademik menyimpan applicant_ref_id/customer snapshot; tidak mengubah applicant langsung. |
| Invoice/payment/clearance | finance_db.* | PMB/Akademik/LMS/Portal hanya membaca via API atau snapshot/event. |
| Mahasiswa/KRS/nilai final | academic_db.* | LMS/Portal menyimpan snapshot; final grade tetap di Akademik. |
| Dosen/karyawan | hris_db.* | Academic/LMS memakai lecturer_ref_id dan snapshot; tidak membuat dosen mandiri. |
| Kelas online/progress | lms_db.* | LMS kelas berasal dari class snapshot Academic, tetapi progress/tugas/presensi LMS dimiliki LMS. |
| Assessment | assessment_db.* | Consumer menerima result event/API; assessment menyimpan attempt, answer, score. |
| Notifikasi/preferensi | portal_db.* | Portal tidak menjadi pemilik status bisnis sumber. |

# 11. Pola Integrasi Antar Modul

| Pola | Kapan Dipakai | Aturan Guardrail |
| --- | --- | --- |
| API Command | Saat modul meminta modul pemilik melakukan perubahan data. | Contoh PMB meminta Finance membuat invoice. Request wajib idempotent. |
| API Query | Saat butuh data real-time dari source of truth. | Gunakan timeout pendek, circuit breaker, dan fallback snapshot bila tersedia. |
| Event Publication | Saat terjadi perubahan status penting. | Event ditulis ke outbox dalam transaksi lokal yang sama dengan perubahan domain. |
| Event Consumption | Saat modul lain perlu update snapshot/read model. | Consumer mencatat event_key di inbox agar tidak memproses duplikat. |
| Snapshot | Untuk tampilan/proses lokal saat dependency down. | Harus memiliki source_event_key/synced_at agar staleness terlihat. |
| Read Model | Untuk dashboard atau query lintas domain yang sering dibaca. | Tidak boleh menjadi sumber kebenaran final; harus bisa direbuild dari event/source. |
| Reconciliation | Untuk memastikan snapshot tidak menyimpang dari source. | Job berkala wajib menghasilkan mismatch report dan retry correction. |
| Warehouse/Data Mart | Untuk laporan lintas modul dan analitik historis. | Tidak boleh digunakan untuk transaksi operasional real-time. |

# 12. Event Catalog Minimum

| Event | Publisher | Consumer Utama | Payload Minimum |
| --- | --- | --- | --- |
| core.person_updated | Core | CRM, PMB, HRIS, Academic, Portal | person_id, full_name, email, phone, status_code, occurred_at |
| reference.study_program_updated | Referensi | PMB, Academic, HRIS, LMS, Portal | study_program_id, code, name, status_code, occurred_at |
| reference.academic_period_updated | Referensi | PMB, Finance, Academic, LMS, Assessment, Portal | academic_period_id, academic_year_id, code, name, term_code, status_code |
| crm.lead_qualified | CRM | PMB | lead_id, person_ref_id, source_code, campaign_id, agent_id, occurred_at |
| pmb.applicant_created | PMB | Finance, Assessment, Portal | applicant_id, person_ref_id, applicant_no, target_period_ref_id, study_program_ref_id |
| finance.invoice_created | Finance | PMB, Portal | invoice_id, invoice_no, bill_to_type, bill_to_ref_id, amount_total, status_code |
| finance.payment_paid | Finance | PMB, Academic, Portal | payment_id, invoice_id, bill_to_type, bill_to_ref_id, paid_amount, paid_at |
| finance.clearance_changed | Finance | PMB, Academic, LMS, Portal | subject_type, subject_ref_id, academic_period_ref_id, service_code, status_code |
| pmb.ready_for_academic | PMB | Academic | applicant_id, person_ref_id, target_period_ref_id, study_program_ref_id, curriculum_candidate_ref |
| academic.student_created | Academic | PMB, Finance, LMS, Portal | student_id, person_ref_id, nim, entry_period_ref_id, study_program_ref_id, curriculum_id |
| academic.class_opened | Academic | LMS, Portal | course_offering_id, academic_period_ref_id, course_id, class_code, lecturer_refs |
| academic.krs_approved | Academic | LMS, Finance, Portal | student_id, krs_id, academic_period_ref_id, krs_item_ids |
| lms.grade_input_submitted | LMS | Academic | student_ref_id, course_offering_ref_id, source_ref_id, score, submitted_at |
| assessment.result_calculated | Assessment | PMB, LMS, Academic, Portal | assessment_session_id, participant_type, participant_ref_id, total_score, passed |

# 13. Kebutuhan Produk Global Revisi

| ID | Prioritas | Kebutuhan Produk | Acceptance Criteria |
| --- | --- | --- | --- |
| PRD-DB-001 | P0 | Setiap modul utama harus memiliki database fisik sendiri. | Core, Referensi, CRM, PMB, Finance, Akademik, HRIS, LMS, Assessment, dan Portal memiliki database terpisah dengan credential dan migration pipeline masing-masing. |
| PRD-DB-002 | P0 | Tidak boleh ada cross-database foreign key. | DDL setiap database hanya mendefinisikan FK internal; relasi lintas modul memakai external_ref_id. |
| PRD-DB-003 | P0 | Transaksi online tidak boleh melakukan direct join lintas database. | Code review dan query audit tidak menemukan join dari service satu modul ke DB modul lain. |
| PRD-DB-004 | P0 | Setiap database modul wajib memiliki audit_logs, idempotency_keys, outbox_events, dan inbox_events. | Semua modul dapat mencatat audit lokal, mencegah retry ganda, publish event, dan consume event secara idempotent. |
| PRD-INT-001 | P0 | Semua command lintas modul harus melalui API resmi modul pemilik. | Modul peminta tidak memiliki credential write ke database modul pemilik. |
| PRD-INT-002 | P0 | Setiap event lintas modul harus memiliki event_key deterministik dan payload contract. | Consumer dapat menolak/mengabaikan event duplikat berdasarkan event_key. |
| PRD-INT-003 | P0 | Setiap snapshot/read model harus menyimpan sumber dan waktu sinkronisasi. | Tabel snapshot/read model memiliki source_event_key/source_module dan synced_at/refreshed_at. |
| PRD-INT-004 | P0 | Setiap modul harus memiliki fallback behavior saat dependency down. | UAT partial outage membuktikan service tidak crash dan memberi status degraded/pending/retry. |
| PRD-INT-005 | P0 | Setiap proses kritis wajib idempotent. | Payment callback, PMB handover, generate NIM, class sync, KRS sync, grade sync, dan notification tidak membuat duplikasi saat retry. |
| PRD-RES-001 | P0 | Sistem harus mendukung partial availability. | Saat satu database non-Core mati, modul lain yang tidak memerlukan real-time data dari DB tersebut tetap dapat membaca/menulis data lokalnya. |
| PRD-RES-002 | P0 | Core dependency harus dimitigasi dengan token validation cache. | Service tetap dapat memvalidasi token aktif yang belum expired menggunakan cached public key/JWKS. |
| PRD-RES-003 | P1 | Dashboard harus menampilkan data freshness. | Setiap widget lintas modul menampilkan refreshed_at atau status data source. |
| PRD-REP-001 | P1 | Laporan lintas modul final harus melewati rekonsiliasi. | Ada mismatch report untuk snapshot PMB-Finance, Academic-Finance, Academic-LMS, LMS-Assessment. |
| PRD-SEC-001 | P0 | Database credential harus scoped per service. | Service hanya memiliki akses ke database miliknya sendiri, kecuali read-only reporting connector yang disetujui. |
| PRD-OPS-001 | P0 | Setiap database memiliki backup/restore dan monitoring sendiri. | Ada health check, backup status, restore test, event lag, dead letter queue, dan alert per modul. |

# 14. Data Independence dan Failure Isolation

| Database Mati | Terdampak Langsung | Tetap Bisa Berjalan | Fallback/Guardrail |
| --- | --- | --- | --- |
| core_db | Login baru, role switching, user provisioning, app launcher real-time. | Token yang belum expired dapat diverifikasi lokal jika JWKS/public key cache tersedia. | Batasi write sensitif; aktifkan cached permission dengan TTL pendek. |
| reference_db | Perubahan master data dan lookup real-time. | Modul tetap memakai reference snapshot terakhir. | Tampilkan warning data referensi mungkin tidak terbaru; retry sync. |
| crm_db | Lead capture, follow-up, campaign, agent commission. | PMB, Finance, Academic, HRIS, LMS, Assessment tetap berjalan. | Handover lead tertunda; tidak berdampak pada applicant yang sudah ada di PMB. |
| pmb_db | Applicant baru, biodata, dokumen, LoA, handover. | CRM lead, Finance pembayaran applicant existing, Academic mahasiswa existing, LMS, HRIS. | Finance/Academic memakai customer/applicant snapshot yang sudah ada. |
| finance_db | Invoice baru, payment verification, clearance real-time, receipt, jurnal. | PMB input biodata/dokumen, Academic data existing, LMS existing class, Assessment. | Gunakan clearance/payment snapshot terakhir; proses yang butuh clear real-time masuk pending_review. |
| academic_db | Generate NIM, KRS, kelas akademik, nilai final, KHS, transkrip. | CRM, PMB, Finance applicant payment, LMS kelas yang sudah sync, Assessment CBT PMB, HRIS. | LMS memakai academic class/student snapshot; class sync baru masuk retry. |
| hris_db | Perubahan dosen/pegawai/homebase/jabatan. | Academic dan LMS tetap memakai lecturer snapshot yang sudah tersinkron. | Plotting dosen baru ditahan bila butuh validasi real-time. |
| lms_db | Pembelajaran online, materi, tugas, presensi LMS, progress. | Academic KRS/kelas/nilai final, Finance, PMB, HRIS, Assessment non-LMS. | Grade sync dari LMS tertunda; final grade tidak boleh tergantung langsung pada DB LMS. |
| assessment_db | CBT/quiz/survey attempt dan scoring. | PMB non-CBT, LMS non-quiz, Academic, Finance, Portal. | Jadwal quiz/CBT masuk postponed; result export retry setelah pulih. |
| portal_db | Dashboard, notification center, user preference, shortcut. | Semua modul operasional sumber tetap berjalan. | Notification event tertahan di outbox/inbox dan diproses ulang setelah Portal pulih. |

# 15. Revisi Requirement Per Modul

## 15.1 Core

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-CORE-DB-001 | P0 | Core harus menyediakan OIDC/JWT dan service client yang dapat divalidasi oleh modul lain tanpa query database Core pada setiap request. |
| PRD-CORE-DB-002 | P0 | Core harus menerbitkan event person/user/role/application update agar modul lain dapat memperbarui snapshot/cache. |
| PRD-CORE-DB-003 | P1 | Core harus menyediakan endpoint introspection untuk kasus verifikasi real-time, tetapi service wajib memiliki fallback cached public key. |

## 15.2 Referensi

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-REF-DB-001 | P0 | Referensi harus menerbitkan event perubahan master data untuk study program, academic year, academic period, status code, payment component, payment method, dan document type. |
| PRD-REF-DB-002 | P0 | Modul consumer harus menyimpan reference snapshot minimum yang dibutuhkan untuk operasi lokal. |
| PRD-REF-DB-003 | P1 | Perubahan master data sensitif wajib memiliki effective_from/effective_to dan tidak memutus histori transaksi lama. |

## 15.3 CRM

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-CRM-DB-001 | P0 | CRM menyimpan lead secara mandiri dan tidak bergantung pada PMB untuk operasi lead capture/follow-up. |
| PRD-CRM-DB-002 | P0 | Konversi lead ke applicant PMB dilakukan melalui API/event idempotent, bukan insert langsung ke pmb_db. |
| PRD-CRM-DB-003 | P1 | CRM menyimpan applicant_ref_id setelah PMB berhasil membuat applicant. |

## 15.4 PMB

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-PMB-DB-001 | P0 | PMB menyimpan applicant dan dokumen secara mandiri dengan person_ref_id dan reference snapshot. |
| PRD-PMB-DB-002 | P0 | PMB tidak menyimpan payment sebagai source of truth; PMB menyimpan applicant_invoice_statuses sebagai read model dari Finance. |
| PRD-PMB-DB-003 | P0 | Handover PMB ke Academic wajib idempotent dan aman ketika Academic down. |

## 15.5 Finance

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-FIN-DB-001 | P0 | Finance menjadi satu-satunya source of truth invoice, payment, receipt, clearance, dan jurnal. |
| PRD-FIN-DB-002 | P0 | Finance menyimpan customer snapshot untuk APPLICANT/STUDENT/PERSON agar invoice/payment tetap dapat diproses walau PMB/Academic sementara down. |
| PRD-FIN-DB-003 | P0 | Payment callback harus idempotent berdasarkan provider_code + provider_event_id dan internal idempotency key. |

## 15.6 Akademik

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-ACA-DB-001 | P0 | Akademik menyimpan student, NIM, curriculum, course offering, KRS, final grade, KHS, transcript secara mandiri. |
| PRD-ACA-DB-002 | P0 | Akademik tidak query langsung finance_db saat transaksi KRS; validasi menggunakan Finance API atau clearance snapshot sesuai policy. |
| PRD-ACA-DB-003 | P0 | Akademik menerbitkan class_opened, krs_approved, final_grade_published untuk LMS/Portal/reporting. |

## 15.7 HRIS

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-HRIS-DB-001 | P0 | HRIS menjadi source of truth employee dan lecturer. Academic/LMS memakai lecturer_ref_id dan lecturer snapshot. |
| PRD-HRIS-DB-002 | P1 | Perubahan status dosen aktif/nonaktif dipublish sebagai event agar Academic/LMS dapat menolak plotting baru. |

## 15.8 LMS

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-LMS-DB-001 | P0 | LMS tidak membuat kelas akademik; LMS membuat online class dari event/snapshot course_offering Academic. |
| PRD-LMS-DB-002 | P0 | Enrollment LMS berasal dari krs_approved event dan disimpan sebagai enrollment lokal. |
| PRD-LMS-DB-003 | P0 | Grade input dari LMS ke Academic wajib idempotent dan tidak boleh menimpa final grade. |

## 15.9 Assessment

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-ASM-DB-001 | P0 | Assessment menyimpan question bank, version, session, participant snapshot, attempt, answer, dan scoring result secara mandiri. |
| PRD-ASM-DB-002 | P0 | Result dikirim ke consumer melalui API/event dengan result_export dan retry policy. |

## 15.10 Portal

| ID | Prioritas | Requirement Revisi |
| --- | --- | --- |
| PRD-POR-DB-001 | P0 | Portal menyimpan notification, read marker, user preference, shortcut, dan dashboard read model, bukan data bisnis sumber. |
| PRD-POR-DB-002 | P1 | Dashboard widget harus menyimpan refreshed_at/source_modules agar data freshness terlihat. |

# 16. Non-Functional Requirement Tambahan

| Kategori | Requirement | Acceptance Criteria |
| --- | --- | --- |
| Availability | Setiap modul memiliki health check aplikasi, database, event publisher, event consumer, dan dependency eksternal. | Status health modul dapat dilihat oleh DevOps dan Portal admin. |
| RTO/RPO | RTO/RPO ditetapkan per database modul sesuai kritikalitas. | Core/Finance/Academic memiliki target lebih ketat daripada Portal/CRM bila disepakati bisnis. |
| Observability | Setiap request lintas modul membawa request_id/correlation_id. | Trace dapat mengikuti flow PMB invoice hingga Finance payment dan Portal notification. |
| Event Reliability | Outbox publisher dan inbox consumer memiliki retry dan dead letter queue. | Event gagal tidak hilang dan dapat direprocess. |
| Security | Service credential hanya dapat mengakses database modul sendiri. | Secret scan dan database role audit tidak menemukan credential lintas DB tidak sah. |
| Data Freshness | Read model/snapshot memiliki timestamp sinkronisasi. | UI lintas modul menampilkan atau menyimpan freshness metadata. |
| Backup/Restore | Backup dilakukan per database dan restore diuji berkala. | Restore test menghasilkan bukti dan waktu pemulihan. |
| Performance | Query transaksi hanya memakai database lokal modul. | Tidak ada direct query OLTP ke DB modul lain. |

# 17. Business Rule Global Revisi

Tidak ada modul yang boleh menyimpan credential, password, atau session authority selain Core.

Tidak ada modul yang boleh melakukan write langsung ke database modul lain.

Tidak ada cross-database foreign key dalam desain fisik.

Tidak ada direct cross-database join untuk transaksi online.

Setiap external_ref_id harus dapat divalidasi melalui API atau event source modul pemilik.

Setiap snapshot/read model harus menyimpan sumber, waktu sinkronisasi, dan status pemrosesan.

Snapshot/read model tidak boleh menjadi dasar laporan final tanpa rekonsiliasi.

Payment callback, PMB handover, generate NIM, class sync, KRS sync, grade sync, dan notification delivery wajib idempotent.

Final grade tetap dimiliki Academic walaupun input berasal dari LMS atau Assessment.

Finance tetap menjadi source of truth clearance walaupun Academic/LMS menyimpan clearance snapshot.

Portal tidak boleh mengubah status bisnis sumber kecuali melalui API modul pemilik.

Modul consumer harus aman terhadap event duplikat, event terlambat, dan event out of order.

# 18. Data Quality dan Reconciliation

| Area Rekonsiliasi | Sumber A | Sumber B | Mismatch yang Harus Dideteksi |
| --- | --- | --- | --- |
| Applicant Payment | pmb_db.applicant_invoice_statuses | finance_db.invoices/payments | Status invoice/payment berbeda, invoice hilang, amount mismatch. |
| PMB Handover | pmb_db.handover_logs | academic_db.students | Applicant ready tetapi student belum dibuat, duplicate student_ref_id. |
| Student Clearance | academic_db.student_clearance_snapshots | finance_db.clearances | Clearance KRS/ujian/KHS/transkrip tidak sinkron. |
| Academic-LMS Class | academic_db.course_offerings/krs_items | lms_db.lms_classes/lms_enrollments | Class belum sync, enrollment kurang/lebih, student inactive masih enrolled. |
| LMS/Assessment Grade Input | lms_db.grade_inputs/assessment_db.scoring_results | academic_db.grade_inputs | Score belum terkirim, score duplikat, source_ref_id ganda. |
| Portal Dashboard | portal_db.dashboard_read_models | Source module summary API/warehouse | Widget stale, source unavailable, angka KPI berbeda dari source. |

# 19. Pola Query Data Lintas Modul

Dalam desain v6.5, join lintas modul dilakukan secara logis, bukan SQL join langsung antar database produksi. Berikut pola resminya.

| Kebutuhan | Pola Resmi | Contoh |
| --- | --- | --- |
| Detail applicant + payment status | PMB API membaca pmb_db.applicants dan pmb_db.applicant_invoice_statuses; jika butuh real-time, call Finance API. | GET /pmb/applicants/{id} + optional GET /finance/invoices?bill_to=APPLICANT:{id} |
| KRS + clearance | Academic API membaca krs lokal dan clearance snapshot; jika policy real-time, call Finance API sebelum finalisasi. | POST /academic/krs/{id}/finalize memvalidasi finance.clearance snapshot/API. |
| LMS class + student list | LMS membaca lms_db.lms_classes, academic_class_snapshots, student_snapshots. | Tidak join ke academic_db saat kelas dibuka. |
| Dashboard pimpinan | Portal membaca portal_db.dashboard_read_models atau summary API masing-masing modul. | Widget menampilkan refreshed_at dan source status. |
| Laporan final semester | Warehouse/data mart membaca event/CDC semua modul dan menjalankan rekonsiliasi. | Laporan KRS vs pembayaran vs nilai final. |

- Contoh anti-pattern yang dilarang:
SELECT * FROM pmb_db.applicants a
JOIN finance_db.invoices i ON i.bill_to_ref_id = a.id;

-- Pola yang benar:
-- 1) PMB membaca applicant lokal.
-- 2) PMB membaca applicant_invoice_statuses lokal, atau call Finance API.
-- 3) Finance tetap source of truth invoice/payment.

# 20. Release Plan Revisi

| Fase | Fokus | Output Wajib |
| --- | --- | --- |
| Phase 0 - Architecture Foundation | Database boundary, service identity, event broker, outbox/inbox library, idempotency standard, observability. | Template DB modul, migration standard, event contract base, retry/DLQ, correlation_id. |
| Phase 1 - Core + Referensi | Identity, role, permission, service client, master data event. | core_db, reference_db, JWKS cache, reference snapshots. |
| Phase 2 - CRM + PMB + Finance Basic | Lead-to-applicant, invoice, payment, applicant invoice status read model. | crm_db, pmb_db, finance_db, event PMB-Finance. |
| Phase 3 - Academic Core | Student/NIM, curriculum, course offering, KRS, clearance snapshot. | academic_db, handover PMB-Academic, Finance clearance sync. |
| Phase 4 - HRIS + LMS | Lecturer source, class sync, KRS enrollment, LMS delivery. | hris_db, lms_db, academic class/student/lecturer snapshots. |
| Phase 5 - Assessment + Grade Input | CBT/quiz/survey, scoring, result export, grade input to Academic. | assessment_db, LMS/PMB/Academic result integration. |
| Phase 6 - Portal + Reporting | Notification, dashboard read model, data freshness, warehouse/reporting plan. | portal_db, dashboard_read_models, reconciliation reports. |

# 21. Acceptance Criteria dan UAT Partial Outage

Matikan finance_db. PMB tetap dapat membuat/mengubah applicant dan upload dokumen. Status invoice menampilkan snapshot terakhir atau “payment status unavailable”.

Matikan academic_db. LMS tetap dapat membuka kelas dan enrollment yang sudah tersinkron; class sync baru masuk retry queue.

Matikan lms_db. Academic tetap dapat membuka kelas, KRS, dan final grade tanpa error database LMS.

Matikan portal_db. PMB, Finance, Academic, LMS, HRIS, dan Assessment tetap dapat melakukan transaksi sumber; notification event diproses ulang setelah Portal pulih.

Kirim payment callback yang sama dua kali. Finance hanya membuat satu payment/receipt/journal effect.

Kirim event finance.payment_paid dua kali ke PMB. PMB applicant_invoice_statuses tidak duplikat dan inbox mencatat event sudah diproses.

Kirim handover applicant yang sama dua kali. Academic hanya membuat satu student dan satu NIM.

Ubah status dosen menjadi nonaktif di HRIS. Academic/LMS menerima event dan menolak plotting baru setelah snapshot tersinkron.

Dashboard Portal menampilkan refreshed_at dan status source ketika salah satu modul down.

Reconciliation job mendeteksi mismatch antara pmb applicant invoice status dan Finance invoice/payment source.

# 22. Risiko dan Mitigasi

| Risiko | Dampak | Mitigasi |
| --- | --- | --- |
| Kompleksitas arsitektur naik. | Tim development dan QA membutuhkan disiplin event/API yang lebih tinggi. | Sediakan platform library untuk outbox/inbox, idempotency, contract testing, dan tracing. |
| Event terlambat atau gagal diproses. | Snapshot/read model stale. | Retry, DLQ, event lag monitoring, reconciliation job. |
| Data lintas modul berbeda sementara. | User melihat status yang belum terbaru. | Tampilkan freshness metadata dan status pending/retry. |
| Core menjadi dependency kritikal. | Login baru dan role switching terganggu. | JWT public key cache, permission snapshot, token TTL yang terukur, Core HA. |
| Reference snapshot stale. | Dropdown/status lokal tidak terbaru. | Effective dating, event update, background sync, validasi saat command sensitif. |
| Reporting lintas modul berat. | Database transaksi terganggu bila query langsung. | Gunakan warehouse/data mart dan read replica khusus reporting. |
| Infrastruktur fisik tetap single point of failure. | Semua database bisa mati jika satu cluster/server sama gagal. | Pisahkan instance/cluster untuk modul kritikal sesuai target availability. |

# 23. Dokumen Turunan Setelah Revisi PRD

| Dokumen Turunan | Fungsi | Owner Awal |
| --- | --- | --- |
| Distributed ERD/DBML per Modul | Menentukan tabel internal, external_ref_id, snapshot, read model, outbox/inbox per database. | DBA + System Analyst |
| API Contract per Modul | Menentukan command/query resmi antar service dan error contract. | System Analyst + Technical Lead |
| Event Contract Catalog | Menentukan event_key, schema payload, publisher, consumer, versioning, dan retry policy. | System Analyst + Technical Lead |
| Data Dictionary per Database | Menjelaskan field, tipe data, constraint internal, index, PII classification. | DBA |
| State Machine | Menentukan status applicant, invoice, payment, clearance, KRS, grade, handover. | System Analyst |
| RBAC Matrix + Service Scope | Menentukan role user dan permission service-to-service. | System Analyst + Security |
| Resilience Test Plan | Menguji DB down, API timeout, duplicate event, delayed event, DLQ, restore. | QA/UAT + DevOps |
| Migration Strategy | Menentukan migrasi dari schema-per-domain ke physical DB per modul jika sudah ada implementasi sebelumnya. | DBA + Technical Lead |

# Appendix A - Database Boundary Minimum

| Database | Internal FK Boleh | External Ref Yang Dipakai | Read Model/Snapshot Minimum |
| --- | --- | --- | --- |
| core_db | users -> persons; role_assignments -> users/roles | scope_ref_id ke Ref/HRIS sesuai kebutuhan | Tidak wajib, kecuali app health/status modul. |
| reference_db | academic_periods -> academic_years; study_programs -> academic_levels | actor_user_ref_id ke Core untuk audit | Tidak wajib. |
| crm_db | leads -> campaigns/agents; follow_ups -> leads | person_ref_id, user_ref_id, applicant_ref_id | person_snapshots. |
| pmb_db | applicant_biodata/documents/selection -> applicants | person_ref_id, lead_ref_id, invoice_ref_id, assessment_session_ref_id, student_ref_id | person_snapshots, reference_snapshots, applicant_invoice_statuses. |
| finance_db | payments -> invoices; invoice_items -> invoices; journal_lines -> journal_entries | bill_to_ref_id, person_ref_id, academic_period_ref_id | customer_snapshots. |
| academic_db | krs_items -> krs_headers; grades -> krs_items; courses -> curriculums internal mapping | person_ref_id, applicant_ref_id, lecturer_ref_id, academic_period_ref_id | person_snapshots, reference_snapshots, student_clearance_snapshots. |
| hris_db | lecturers -> employees; employment_records -> employees | person_ref_id, study_program_ref_id, academic_period_ref_id | person_snapshots. |
| lms_db | enrollments -> lms_classes; sessions/materials/assignments -> lms_classes | course_offering_ref_id, student_ref_id, lecturer_ref_id, assessment_session_ref_id | academic_class_snapshots, student_snapshots, lecturer_snapshots. |
| assessment_db | questions -> banks; versions -> questions; answers -> attempts | context_ref_id, participant_ref_id, user_ref_id | participant_snapshots. |
| portal_db | notifications -> notification_events | user_ref_id, role_ref_id, application_ref_id, source_entity_id | user_snapshots, role_snapshots, dashboard_read_models. |

# Appendix B - Event Contract Standard

Appendix ini melengkapi Section 12 Event Catalog Minimum. Tujuannya adalah mengubah daftar event menjadi kontrak teknis yang dapat langsung dipakai oleh Backend, QA/UAT, DBA, Security, DevOps, dan Owner Modul. Standar ini berlaku untuk seluruh event pada Core, Referensi, CRM, PMB, Finance, Akademik, HRIS, LMS, Assessment, Portal, dan Reporting.

## B.1 Prinsip Umum Event Contract

| ID | Prinsip | Ketentuan |
| --- | --- | --- |
| EC-G-001 | Setiap event wajib memiliki identitas unik. | Event harus memiliki event_name, event_version, event_key, aggregate_type, aggregate_id, occurred_at, dan publisher_service. |
| EC-G-002 | Event ditulis melalui transactional outbox. | Event hanya boleh dipublish setelah transaksi lokal database pemilik berhasil commit. |
| EC-G-003 | Setiap consumer wajib idempotent. | Consumer wajib mencatat event_key ke inbox_events sebelum atau saat pemrosesan event. |
| EC-G-004 | Snapshot bukan source of truth. | Snapshot harus memuat source_event_key dan synced_at agar staleness data dapat dibaca. |
| EC-G-005 | Event contract wajib versioned. | Perubahan payload yang tidak backward compatible harus menaikkan event_version. |
| EC-G-006 | Retry dan DLQ wajib tersedia. | Event gagal proses tidak boleh hilang. Event harus masuk retry queue dan DLQ jika gagal permanen. |
| EC-G-007 | Reconciliation wajib untuk data lintas modul kritis. | Snapshot/read model yang berbeda dari source of truth harus menghasilkan mismatch report. |

## B.2 Struktur Wajib Event Identity

| Field | Tipe | Deskripsi |
| --- | --- | --- |
| event_name | string | Nama event dengan format domain.action, contoh finance.payment_paid. |
| event_version | string | Versi schema event, contoh v1. |
| event_key | string | Kunci unik global untuk idempotency dan duplicate handling. |
| event_type | string | DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, atau SNAPSHOT_EVENT. |
| publisher_service | string | Nama service pengirim event. |
| publisher_database | string | Database source of truth pemilik event. |
| aggregate_type | string | Objek bisnis utama, contoh payment, invoice, applicant, student, krs. |
| aggregate_id | uuid/string | ID objek bisnis utama pada database pemilik. |
| correlation_id | string | ID untuk melacak satu proses end-to-end lintas service. |
| causation_id | string | ID event atau command yang memicu event ini. |
| occurred_at | datetime | Waktu kejadian bisnis terjadi. |
| published_at | datetime | Waktu event berhasil dikirim ke broker. |

Contoh envelope event:

{
 "event_name": "finance.payment_paid",
 "event_version": "v1",
 "event_key": "finance.payment_paid:payment_id:8f2c:v1",
 "event_type": "INTEGRATION_EVENT",
 "publisher_service": "finance-service",
 "publisher_database": "finance_db",
 "aggregate_type": "payment",
 "aggregate_id": "8f2c",
 "correlation_id": "corr-20260619-001",
 "causation_id": "payment_callback:gateway:trx-9001",
 "occurred_at": "2026-06-19T10:15:00Z",
 "published_at": "2026-06-19T10:15:05Z"
}

## B.3 Business Trigger dan State Change

| Komponen | Penjelasan |
| --- | --- |
| Field | Keterangan |
| business_trigger | Kondisi bisnis yang menyebabkan event diterbitkan. |
| pre_condition | Status atau kondisi data sebelum event terjadi. |
| post_condition | Status atau kondisi data setelah event terjadi. |
| source_table | Tabel utama yang menjadi sumber perubahan. |
| state_transition | Perubahan status, contoh UNPAID menjadi PAID. |
| publish_timing | Waktu event boleh dipublish, umumnya setelah commit transaksi lokal. |

## B.4 Publisher, Consumer, dan Tujuan Konsumsi

| Event | Publisher | Consumer | Tujuan Konsumsi |
| --- | --- | --- | --- |
| core.person_updated | Core | CRM, PMB, Academic, HRIS, LMS, Assessment, Portal | Memperbarui person snapshot lokal. |
| pmb.applicant_created | PMB | Finance, Assessment, Portal, Reporting | Membuat konteks applicant, notifikasi, dan projection dashboard. |
| finance.invoice_created | Finance | PMB, Academic, Portal, Reporting | Memperbarui status tagihan pada consumer. |
| finance.payment_paid | Finance | PMB, Academic, Portal, Reporting | Memperbarui payment status, clearance snapshot, notifikasi, dan laporan. |
| finance.clearance_changed | Finance | PMB, Academic, LMS, Portal | Mengatur kelayakan layanan akademik. |
| academic.student_created | Academic | PMB, LMS, Portal, Reporting | Menghubungkan applicant dengan student/NIM. |
| academic.krs_approved | Academic | LMS, Portal, Reporting | Membuat atau memperbarui enrollment LMS. |
| assessment.result_calculated | Assessment | PMB, LMS, Academic, Portal | Mengirim hasil assessment ke context owner. |

## B.5 Payload Schema dan Validation Rule

Setiap event wajib memiliki payload schema. Payload hanya memuat data yang dibutuhkan consumer. Data sensitif dan PII harus dibatasi sesuai kebutuhan bisnis dan prinsip least privilege.

| Field | Required | Validation Rule |
| --- | --- | --- |
| payment_id | Ya | UUID dan harus ada pada finance_db.payments. |
| invoice_id | Ya | UUID dan harus ada pada finance_db.invoices. |
| invoice_no | Ya | String, tidak kosong. |
| bill_to_type | Ya | APPLICANT, STUDENT, atau PERSON. |
| bill_to_ref_id | Ya | UUID external reference subject pembayaran. |
| paid_amount | Ya | Decimal lebih dari 0. |
| payment_method_code | Ya | Kode metode pembayaran dari master/reference. |
| paid_at | Ya | Datetime valid. |
| status_code | Ya | Harus PAID untuk event finance.payment_paid. |

Contoh payload finance.payment_paid:

{
 "payment_id": "8f2c",
 "invoice_id": "inv-1001",
 "invoice_no": "INV/PMB/2026/0001",
 "bill_to_type": "APPLICANT",
 "bill_to_ref_id": "app-2201",
 "paid_amount": 2500000,
 "payment_method_code": "VA_BNI",
 "paid_at": "2026-06-19T10:15:00Z",
 "status_code": "PAID"
}

## B.6 Idempotency dan Duplicate Handling

| Area | Aturan |
| --- | --- |
| Event key | event_key harus unik dan deterministik. Format disarankan: {event_name}:{aggregate_id}:{event_version}. |
| Consumer inbox | Consumer wajib menyimpan event_key pada inbox_events. |
| Duplicate event | Jika event_key sudah pernah diproses, consumer tidak memproses ulang payload. |
| Retry command | Command/API yang memicu event wajib membawa idempotency_key. |
| Conflict payload | Jika event_key sama tetapi payload berbeda, consumer menolak event dan mencatat mismatch untuk investigasi. |

## B.7 Ordering, Dependency, dan Causality

| Event | Wajib Setelah | Catatan |
| --- | --- | --- |
| finance.payment_paid | finance.invoice_created | Payment tidak valid tanpa invoice. |
| finance.clearance_changed | finance.payment_paid atau finance.clearance_reviewed | Clearance berubah berdasarkan status finance resmi. |
| academic.student_created | pmb.ready_for_academic atau pmb.handover_requested | Mahasiswa dibuat setelah applicant siap diserahkan ke Akademik. |
| lms.enrollment_synced | academic.krs_approved | Enrollment LMS hanya berasal dari KRS valid. |
| academic.final_grade_published | academic.grade_input_received | LMS/Assessment hanya memberi input, final grade tetap di Akademik. |

## B.8 Retry Policy dan Dead Letter Queue

| Komponen | Aturan |
| --- | --- |
| Retry schedule | 1 menit, 5 menit, 15 menit, lalu exponential backoff. |
| Maksimal retry | 10 kali atau sesuai SLA modul. |
| Temporary failure | Event tetap berada pada retry queue. |
| Permanent failure | Event masuk DLQ setelah retry maksimum. |
| DLQ payload | Wajib menyimpan event_key, consumer, last_error, retry_count, failed_at, dan raw_payload. |
| Recovery | Event dari DLQ dapat direplay secara manual oleh role DevOps/SRE yang berwenang. |

Field teknis yang disarankan pada outbox/inbox:

retry_count
next_retry_at
last_error
locked_at
locked_by
processed_at
dead_letter_at
schema_version
correlation_id
causation_id

## B.9 Snapshot dan Read Model Impact

| Event | Consumer | Tabel Lokal yang Diupdate |
| --- | --- | --- |
| core.person_updated | PMB, Academic, LMS | person_snapshots |
| reference.study_program_updated | PMB, Academic, HRIS, LMS | reference_snapshots atau study_program_snapshots |
| finance.invoice_created | PMB, Academic, Portal | applicant_invoice_statuses, student_finance_snapshots, dashboard_read_models |
| finance.clearance_changed | Academic, LMS, Portal | student_clearance_snapshots, lms_clearance_snapshots, dashboard_read_models |
| academic.krs_approved | LMS, Portal | lms_enrollments, dashboard_read_models |
| assessment.result_calculated | PMB, LMS, Academic | assessment_result_snapshots atau grade_inputs |

Setiap snapshot/read model minimal memiliki source_event_key, source_event_name, source_updated_at, synced_at, dan sync_status.

## B.10 Reconciliation Rule

| Relasi Kritis | Reconciliation Rule |
| --- | --- |
| PMB invoice snapshot vs Finance invoice/payment | Cocokkan invoice_id, payment status, paid_amount, dan paid_at. |
| Academic clearance snapshot vs Finance clearance | Cocokkan subject_ref_id, service_code, academic_period_ref_id, dan status_code. |
| LMS enrollment vs Academic KRS | Cocokkan krs_item_ref_id, student_ref_id, course_offering_ref_id, dan enrollment status. |
| Academic grade input vs LMS/Assessment source | Cocokkan source_module, source_ref_id, score, weight, dan submitted_at. |
| Portal dashboard read model vs source events | Cocokkan refreshed_at, source status, dan payload aggregate. |

Jika snapshot berbeda dari source of truth, sistem membuat mismatch report. Jika dapat diperbaiki otomatis, correction job dijalankan. Jika perlu keputusan admin, status data menjadi pending_review.

## B.11 Security dan Authorization Event

| Aspek | Aturan |
| --- | --- |
| Publisher authorization | Hanya service owner domain yang boleh publish event domainnya. |
| Consumer authorization | Hanya consumer terdaftar yang boleh subscribe event tertentu. |
| PII minimization | Payload tidak boleh membawa data pribadi berlebihan. Gunakan ref_id dan snapshot minimum. |
| Service authentication | Publish dan consume event memakai service credential yang dikelola Core/Security. |
| Audit | Publish, consume, retry, DLQ, dan replay event wajib tercatat audit. |
| Replay control | Replay DLQ hanya boleh dilakukan role DevOps/SRE atau admin teknis yang diberi izin. |

## B.12 Error Contract

| Error Code | Keterangan | Tindakan |
| --- | --- | --- |
| EVENT_DUPLICATE | event_key sudah pernah diproses. | Abaikan payload dan catat sebagai duplicate. |
| EVENT_SCHEMA_INVALID | Payload tidak sesuai schema. | Tolak event dan catat error. |
| EVENT_VERSION_UNSUPPORTED | Consumer belum mendukung event_version. | Masukkan ke DLQ atau compatibility queue. |
| SOURCE_REF_NOT_FOUND | External reference tidak ditemukan. | Retry jika kemungkinan event pendahulu belum masuk. |
| SNAPSHOT_UPDATE_FAILED | Snapshot lokal gagal diperbarui. | Retry dan catat last_error. |
| RECONCILIATION_REQUIRED | Data source dan snapshot berbeda. | Buat mismatch report. |
| CONSUMER_TEMPORARY_FAILURE | Consumer gagal sementara. | Retry sesuai retry policy. |
| CONSUMER_PERMANENT_FAILURE | Consumer gagal permanen. | Masuk DLQ. |

## B.13 Observability dan Monitoring

| Metric | Tujuan |
| --- | --- |
| outbox_pending_count | Mengukur event yang belum dipublish. |
| inbox_pending_count | Mengukur event masuk yang belum diproses. |
| event_lag_seconds | Mengukur keterlambatan event dari occurred_at ke processed_at. |
| retry_count_by_event | Melihat event yang sering gagal. |
| dlq_count_by_consumer | Menemukan consumer bermasalah. |
| reconciliation_mismatch_count | Mengukur selisih source of truth dan snapshot/read model. |

## B.14 UAT Event Contract

| Skenario UAT | Expected Result |
| --- | --- |
| Event valid diterbitkan. | Event masuk outbox dan dipublish ke broker setelah transaksi lokal commit. |
| Event diterima consumer. | Consumer menyimpan event_key pada inbox_events dan memperbarui snapshot/read model. |
| Event yang sama dikirim dua kali. | Consumer hanya memproses satu kali. |
| Consumer down. | Event masuk retry queue dan diproses ulang saat consumer pulih. |
| Payload tidak valid. | Event ditolak dengan EVENT_SCHEMA_INVALID. |
| Event version tidak didukung. | Event masuk DLQ atau compatibility queue. |
| Snapshot berbeda dengan source. | Reconciliation job membuat mismatch report. |
| DLQ direplay. | Replay tercatat audit dan tidak membuat data duplikat. |

## B.15 Template Final Event Contract

Event Name :
Event Version :
Event Type :
Publisher Service :
Publisher DB :
Source Table :
Aggregate Type :
Aggregate ID :
Business Trigger :
Pre-condition :
Post-condition :
Consumer :
Payload Schema :
Validation Rule :
Idempotency Rule :
Ordering Rule :
Retry Policy :
DLQ Policy :
Snapshot Impact :
Reconciliation :
Security Rule :
Error Contract :
Observability :
UAT Scenario :
