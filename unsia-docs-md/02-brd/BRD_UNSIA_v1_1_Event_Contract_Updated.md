---
title: "BRD UNSIA"
source_file: "BRD_UNSIA_v1_1_Event_Contract_Updated.docx"
format: markdown
---

# BRD UNSIA

UNSIA

BUSINESS REQUIREMENT DOCUMENT

BRD Global dan Per Modul - Distributed Modular Database FULL

Produk: ERP Pendidikan / SIAKAD Terintegrasi UNSIA

Versi BRD: v1.1 Distributed Database FULL

Tanggal: 19 Juni 2026

Status: Revised Full Draft

Basis Penyusunan: PRD Global UNSIA v6.5 Distributed Database FULL, PRD Global UNSIA v6.4 Detailed FULL, BRD Global UNSIA v1.0 Draft, BRD Per Modul UNSIA v1.0 Draft, dan revisi arsitektur database fisik per modul.

Tujuan Dokumen: Menurunkan keputusan produk v6.5 menjadi kebutuhan bisnis yang dapat divalidasi stakeholder serta menjadi dasar FSD, ERD/DBML per database, API Contract, Event Contract, RBAC Matrix, UAT, migration mapping, dan release plan.

# 1. Kontrol Dokumen

| Versi | Tanggal | Status | Catatan |
| --- | --- | --- | --- |
| v1.0 | 18 Juni 2026 | Draft Awal | BRD Global dan BRD Per Modul disusun berdasarkan PRD Global v6.4 dengan asumsi satu PostgreSQL utama dan schema per domain. |
| v1.1 | 19 Juni 2026 | Revised Full Draft | BRD direvisi mengikuti keputusan physical database per modul, no cross-database FK, no online cross-database join, event-driven integration, snapshot/read model, dan graceful degradation. |
| v1.1.1 | 22 Juni 2026 | Updated Draft | Penambahan Appendix A - Event Contract Standard sebagai standar kebutuhan bisnis dan teknis untuk event identity, trigger bisnis, payload schema, idempotency, retry, DLQ, snapshot/read model, reconciliation, security, error handling, observability, dan UAT event. |

## 1.1 Peran dan Tanggung Jawab

| Peran | Tanggung Jawab |
| --- | --- |
| Product Owner | Mengesahkan kebutuhan bisnis, scope, prioritas, release MVP, dan acceptance criteria global. |
| System Analyst | Menyusun BRD, menjaga traceability ke PRD, mengubah BRD menjadi FSD, API/Event Contract, RBAC, UAT, dan state machine. |
| Owner Modul | Memvalidasi proses bisnis, aturan operasional, laporan, risiko, dan acceptance criteria per modul. |
| Technical Lead | Menilai kelayakan implementasi multi-repo, API, event broker, database terpisah, dan dependency teknis lintas modul. |
| DBA | Memastikan kebutuhan bisnis dapat diturunkan menjadi database fisik per modul, constraint lokal, index, external reference, snapshot, read model, migration, backup, dan reconciliation. |
| Backend Lead | Menyusun service boundary, command/query API, event publisher/consumer, idempotency, dan error handling. |
| Frontend Lead | Menyusun flow UI, degraded-mode behavior, sync timestamp display, dan role-based navigation. |
| QA/UAT Lead | Menurunkan business requirement menjadi skenario UAT normal, negative case, integration case, failure case, dan reconciliation case. |
| DevOps/SRE | Menyediakan deployment, monitoring, alert, log, trace, event broker, retry, DLQ, backup, restore, dan runbook. |

# 2. Ringkasan Eksekutif

UNSIA membutuhkan ERP Pendidikan / SIAKAD Terintegrasi untuk mengelola lifecycle kampus dari lead, applicant, pembayaran, mahasiswa aktif, KRS, LMS, assessment, nilai, KHS, transkrip, sampai alumni. BRD ini memformalkan kebutuhan bisnis lintas modul setelah keputusan arsitektur database direvisi menjadi physical database per modul.

Keputusan bisnis kunci dalam BRD v1.1 adalah sebagai berikut:

### 1. Setiap modul memiliki database fisik sendiri: `core_db`, `reference_db`, `crm_db`, `pmb_db`, `finance_db`, `academic_db`, `hris_db`, `lms_db`, `assessment_db`, dan `portal_db`.

2. Tidak ada foreign key lintas database. Relasi antar modul memakai external reference seperti `person_ref_id`, `applicant_ref_id`, `student_ref_id`, `academic_period_ref_id`, dan `invoice_ref_id`.

3. Tidak ada direct cross-database join untuk transaksi online. Kebutuhan tampilan lintas modul dipenuhi melalui API composition, snapshot, read model, event projection, atau warehouse.

### 4. Jika salah satu database modul mati, proses modul lain yang tidak berhubungan langsung tetap berjalan dalam mode normal atau degraded mode.

### 5. Source of truth tetap tunggal per domain. Snapshot dan read model bukan sumber kebenaran utama.

### 6. Proses kritis wajib memiliki idempotency, audit trail, event log, dan reconciliation.

### 7. Tahun Ajaran, Periode Akademik, Tahun Kurikulum, dan Kurikulum Prodi tetap dipisahkan secara bisnis dan data.

# 3. Latar Belakang dan Masalah Bisnis

Operasional perguruan tinggi memiliki rangkaian proses yang saling bergantung. CRM menghasilkan lead; PMB mengubah lead menjadi applicant; Finance mengelola invoice, payment, dan clearance; Akademik mengubah applicant valid menjadi mahasiswa; HRIS menjadi sumber dosen; LMS menjalankan pembelajaran online berdasarkan kelas dan KRS; Assessment menjadi mesin CBT/quiz/survey; Portal menyajikan dashboard dan notifikasi.

Pada model terpusat, database tunggal memberikan kemudahan join dan constraint, tetapi memiliki risiko operasional: kegagalan database utama atau kesalahan migration dapat berdampak luas ke semua modul. Karena target bisnis UNSIA adalah modul yang lebih tahan gangguan, BRD ini menetapkan pendekatan distributed modular database.

| Masalah Bisnis | Dampak | Arah Solusi Bisnis |
| --- | --- | --- |
| Data calon mahasiswa, mahasiswa, dosen, dan transaksi tersebar tanpa ownership jelas. | Duplikasi data, laporan tidak konsisten, dan audit sulit. | Menetapkan source of truth per domain dan external reference resmi. |
| Satu database utama menjadi titik kegagalan besar. | Gangguan DB dapat menghentikan proses lintas modul. | Memisahkan database fisik per modul dan membuat graceful degradation. |
| Status payment tidak selalu menjadi dasar layanan akademik. | KRS, ujian, KHS, transkrip, atau wisuda dapat berjalan tanpa kontrol finansial. | Finance menjadi source of truth clearance dan publish clearance event. |
| Tahun Ajaran sering rancu dengan Tahun Kurikulum. | Kelas, KRS, invoice, LMS, nilai, dan transkrip rawan salah periode. | Memisahkan kalender operasional dengan versi kurikulum. |
| Role hanya dikontrol di tampilan. | Kebocoran data lintas prodi atau lintas user. | Backend wajib menegakkan permission dan data scope. |
| Integrasi manual antar modul. | Data dobel atau hilang saat retry. | Setiap proses kritis memakai idempotency, outbox/inbox, audit, dan reconciliation. |
| Reporting lintas modul langsung query ke transaksi. | Beban produksi tinggi dan dependency antar DB makin rapat. | Reporting memakai read model, dashboard projection, atau data warehouse. |

# 4. Tujuan Bisnis dan Indikator Keberhasilan

| Tujuan Bisnis | Deskripsi | Indikator Keberhasilan |
| --- | --- | --- |
| Single identity | Semua pengguna menggunakan satu identitas, akun, role, dan session dari Core. | Tidak ada credential dan role table di luar Core. |
| Lifecycle mahasiswa end-to-end | Data dapat dilacak dari lead, applicant, mahasiswa aktif, alumni. | Tidak ada mahasiswa aktif tanpa histori PMB atau mekanisme create sah. |
| Modular database independence | Setiap modul memiliki database fisik sendiri. | Gangguan database satu modul tidak menghentikan modul lain yang tidak terkait langsung. |
| Source of truth tunggal | Setiap data utama memiliki satu owner bisnis dan teknis. | Modul lain tidak melakukan write langsung ke data domain lain. |
| Kontrol pembayaran dan clearance | Finance mengendalikan status kelayakan layanan akademik. | KRS, ujian, KHS, transkrip, dan wisuda dapat dibatasi sesuai clearance. |
| Konsistensi kalender akademik dan kurikulum | Tahun Ajaran, Periode Akademik, Tahun Kurikulum, dan Kurikulum Prodi tidak bercampur. | Kelas, KRS, invoice, LMS, nilai, dan laporan memakai periode; struktur MK memakai kurikulum. |
| Audit bisnis kuat | Aksi sensitif memiliki actor, role, timestamp, reason, old value, new value. | Payment, handover, NIM, clearance, nilai final dapat ditelusuri. |
| Integrasi idempotent | Retry request/event tidak membuat data ganda. | Payment callback, handover, generate NIM, class sync, grade sync tidak dobel. |
| Degraded operation | Modul tetap berjalan terbatas saat dependency down. | UI menampilkan data terakhir, status sync, dan batasan operasi. |
| Delivery bertahap | Implementasi dapat dilakukan per modul tanpa merusak ownership data. | MVP berjalan bertahap dari foundation, PMB, Finance, Akademik, LMS, Assessment, Portal, Reporting. |

# 5. Ruang Lingkup Bisnis

## 5.1 In Scope

| Modul | Database | Ruang Lingkup Bisnis | Catatan Ownership |
| --- | --- | --- | --- |
| Core | `core_db` | Identitas, SSO, RBAC, active role, application launcher, service client, impersonation, audit global. | Source of truth person, user, role, permission, session. |
| Referensi | `reference_db` | Master data lintas modul, prodi, tahun ajaran, periode akademik, jalur PMB, komponen pembayaran, status code. | Source of truth master data. |
| CRM | `crm_db` | Lead, campaign, agen, referral, follow-up, pipeline, komisi. | Source of truth peminat sebelum applicant. |
| PMB | `pmb_db` | Applicant, biodata, dokumen, seleksi, daftar ulang, LoA, handover ke Akademik. | Source of truth applicant sebelum mahasiswa. |
| Finance | `finance_db` | Invoice, payment, callback, verifikasi manual, receipt, clearance, cicilan, beasiswa, jurnal dasar. | Source of truth transaksi dan clearance. |
| Akademik | `academic_db` | Student, NIM, kurikulum, mata kuliah, kelas, KRS, nilai final, KHS, transkrip, yudisium, alumni. | Source of truth mahasiswa, kelas, KRS, nilai final. |
| HRIS/SDM | `hris_db` | Pegawai, dosen, homebase, unit kerja, jabatan, status aktif, BKD, payroll source. | Source of truth dosen dan karyawan. |
| LMS | `lms_db` | Kelas online, enrollment dari KRS valid, sesi, materi, tugas, presensi, progress, grade input. | Delivery pembelajaran online; bukan owner kelas akademik. |
| Assessment | `assessment_db` | Bank soal, versi soal, CBT, quiz, survey, attempt, jawaban, scoring, result API. | Mesin assessment reusable. |
| Portal | `portal_db` | Dashboard role-based, notification center, preference, shortcut, activity log, read model dashboard. | Presentation layer dan notification center; bukan source transaksi. |
| Reporting | `reporting_db` / warehouse | Agregasi lintas modul, KPI, laporan pimpinan. | Analitik; bukan source transaksi operasional. |

## 5.2 Out of Scope Awal

Mobile app native penuh belum menjadi scope awal; portal responsive/mobile web menjadi prioritas.

Integrasi PDDIKTI/NeoFeeder full automation belum menjadi MVP pertama, tetapi struktur data harus siap.

Vendor payment gateway spesifik belum ditetapkan; BRD hanya mengatur kebutuhan bisnis validasi pembayaran, callback, idempotency, dan rekonsiliasi.

Migrasi historis lengkap dari sistem lama memerlukan dokumen migration mapping terpisah.

Workflow equivalency kurikulum lama-baru belum menjadi MVP penuh.

Distributed transaction atau 2-phase commit lintas database bukan pattern bisnis utama.

Direct cross-database join untuk dashboard/transaksi tidak menjadi scope resmi.

# 6. Stakeholder dan Persona

| Persona/Role | Aktivitas Utama | Kebutuhan Bisnis |
| --- | --- | --- |
| Pendaftar | Mendaftar, mengisi biodata, upload dokumen, seleksi, membayar, daftar ulang, menerima LoA. | Proses jelas, status transparan, invoice valid, bukti pembayaran tercatat, notifikasi tepat waktu. |
| Mahasiswa | Melihat tagihan, KRS, mengikuti LMS, mengerjakan tugas/quiz, melihat nilai/KHS/transkrip. | Akses tunggal, status keuangan jelas, KRS mudah dipahami, kelas LMS otomatis. |
| Dosen | Mengajar, mengelola sesi LMS, materi, tugas, quiz, presensi, nilai input. | Kelas sesuai plotting Akademik dan data dosen dari HRIS. |
| Dosen PA | Menyetujui KRS mahasiswa bimbingan dan memberi catatan akademik. | Akses scoped hanya ke mahasiswa bimbingan dan status clearance/KRS. |
| Admin PMB | Mengelola gelombang, applicant, dokumen, seleksi, daftar ulang, LoA, handover. | Target periode masuk, status applicant, invoice/payment valid, integrasi Assessment. |
| Admin Keuangan | Mengelola invoice, pembayaran, verifikasi, clearance, kas/bank, jurnal, laporan. | Membaca data pendaftar/mahasiswa dari domain sumber dan mengelola audit transaksi. |
| Admin Akademik Biro | Mengelola kalender akademik, NIM, kurikulum, kelas, KRS, nilai, KHS, transkrip, alumni. | Kontrol global akademik dan integrasi dengan PMB, Finance, LMS, HRIS. |
| Admin Akademik Prodi/Kaprodi | Mengelola kurikulum prodi, kelas, dosen pengampu, monitoring mahasiswa prodi. | Scope `study_program_ref_id` ditegakkan backend. |
| Admin SDM/HRIS | Mengelola dosen, pegawai, homebase, jabatan, status aktif. | Data person dari Core dan employment sebagai source of truth. |
| Admin LMS | Mengelola kelas online, materi, tugas, progress, dan sinkronisasi kelas. | Data kelas dan peserta mengikuti Academic. |
| Admin Assessment | Mengelola bank soal, sesi CBT/quiz, scoring, dan result. | Attempt dan score audit-ready serta reusable lintas konteks. |
| Pimpinan | Melihat KPI lintas modul dan risiko operasional. | Dashboard agregat read-only dengan drilldown terbatas dan timestamp sinkronisasi. |
| DBA/SRE | Menjaga availability, backup, restore, event, dan observability. | Database per modul dapat dipulihkan tanpa merusak modul lain. |

# 7. Keputusan Bisnis Utama

| ID | Keputusan Bisnis | Konsekuensi |
| --- | --- | --- |
| BD-001 | Core menjadi identity authority. | Semua modul membaca user/role dari token dan snapshot, bukan membuat login sendiri. |
| BD-002 | Setiap modul memiliki database fisik sendiri. | Tidak ada cross-database FK dan setiap DB punya backup/migration sendiri. |
| BD-003 | Source of truth tunggal per domain. | Modul lain menyimpan external reference dan snapshot, bukan menyalin kepemilikan. |
| BD-004 | Integrasi memakai API/event/read model. | Reporting dan tampilan lintas modul tidak memakai direct join produksi. |
| BD-005 | Tahun Ajaran dan Tahun Kurikulum dipisahkan. | Periode semester tidak boleh menjadi versi kurikulum. |
| BD-006 | Finance menjadi owner clearance. | Layanan akademik bergantung pada status clearance Finance. |
| BD-007 | Academic menjadi owner final grade. | LMS/Assessment hanya memberi grade input, bukan final grade. |
| BD-008 | LMS tidak membuat kelas akademik. | Kelas LMS berasal dari event/API Academic. |
| BD-009 | Assessment reusable lintas konteks. | CBT PMB, quiz LMS, survey memakai mesin sama dengan context berbeda. |
| BD-010 | Portal bukan source transaksi bisnis. | Portal hanya read model, notification, preference, shortcut. |

# 8. Proses Bisnis End-to-End

| Tahap | Deskripsi Bisnis | Modul Pemilik | Output Bisnis | Integrasi |
| --- | --- | --- | --- | --- |
| Lead/Peminat | Peminat masuk dari campaign/referral/agen/landing page/input manual. | CRM | Lead dengan source dan status follow-up. | Core person snapshot opsional. |
| Applicant/Pendaftar | Lead qualified atau pendaftar publik menjadi applicant. | PMB | Applicant dan akun pendaftar. | CRM event/API, Core user/person. |
| Biodata dan Dokumen | Pendaftar melengkapi biodata dan upload dokumen. | PMB | Biodata snapshot, dokumen pending/verified/rejected. | Referensi snapshot. |
| Invoice dan Payment | PMB/Akademik meminta invoice ke Finance. | Finance | Invoice, payment, receipt, payment status. | Finance API/event ke PMB/Academic/Portal. |
| Seleksi dan Daftar Ulang | Applicant mengikuti seleksi/CBT dan daftar ulang. | PMB + Assessment + Finance | Score, result, re-registration. | Assessment result event, Finance payment status. |
| LoA dan Handover | LoA diterbitkan dan applicant ready for academic. | PMB | LoA, handover request. | PMB event/API ke Academic. |
| Generate NIM | Akademik membuat student dan NIM. | Akademik | Mahasiswa baru, NIM, entry period, curriculum. | PMB handover, Finance clearance. |
| KRS dan Kelas | Mahasiswa mengikuti KRS paket/mandiri. | Akademik | KRS header/item, class enrollment. | Finance clearance snapshot, HRIS lecturer snapshot. |
| Pembelajaran Online | Kelas dan peserta KRS disinkronkan ke LMS. | LMS | LMS class, enrollment, session, material, task. | Academic class/KRS event. |
| Quiz/Assessment | Peserta mengerjakan CBT/quiz/survey. | Assessment | Attempt, answer, score, result. | PMB/LMS as context owner. |
| Nilai Final | Input nilai dari LMS/Assessment diproses Akademik. | Akademik | Final grade dan grade history. | LMS/Assessment event/API. |
| KHS/Transkrip | Nilai final digunakan untuk KHS/transkrip. | Akademik | KHS, transkrip, IPK. | Finance clearance untuk akses layanan. |
| Graduation/Alumni | Mahasiswa lulus dan menjadi alumni. | Akademik | Graduation record, alumni. | Finance clearance wisuda. |
| Dashboard/Notifikasi | Status penting ditampilkan ke user/pimpinan. | Portal | Notification, dashboard read model. | Event dari semua modul. |

# 9. Business Rule Global

| ID | Business Rule |
| --- | --- |
| BR-G-001 | Tidak ada password, session, dan role table di luar Core. |
| BR-G-002 | Tidak ada cross-database foreign key. |
| BR-G-003 | Tidak ada direct cross-database join untuk transaksi online. |
| BR-G-004 | Data bisnis utama hanya memiliki satu source of truth. |
| BR-G-005 | Modul tidak boleh write langsung ke database modul lain. |
| BR-G-006 | Semua external reference wajib dapat divalidasi melalui API, event, snapshot, atau reconciliation. |
| BR-G-007 | Setiap event consumer wajib idempotent. |
| BR-G-008 | Snapshot dan read model bukan source of truth. |
| BR-G-009 | Read model lintas modul wajib memiliki `synced_at` atau `refreshed_at`. |
| BR-G-010 | Tahun Ajaran adalah kalender operasional, bukan Tahun Kurikulum. |
| BR-G-011 | Periode Akademik wajib berada di bawah Tahun Ajaran. |
| BR-G-012 | Tahun Kurikulum adalah atribut versi di dalam Kurikulum Prodi. |
| BR-G-013 | Mahasiswa wajib menyimpan `curriculum_id` saat NIM dibuat. |
| BR-G-014 | Gelombang PMB wajib memiliki `target_entry_period_ref_id`. |
| BR-G-015 | Payment status resmi hanya berasal dari Finance. |
| BR-G-016 | Clearance layanan akademik hanya berasal dari Finance. |
| BR-G-017 | Handover PMB ke Academic wajib idempotent. |
| BR-G-018 | Generate NIM wajib idempotent dan unique. |
| BR-G-019 | Kelas LMS wajib berasal dari kelas Academic. |
| BR-G-020 | Enrollment LMS wajib berasal dari KRS valid. |
| BR-G-021 | LMS/Assessment tidak boleh langsung menimpa final grade Academic. |
| BR-G-022 | Question yang sudah dipakai attempt tidak boleh diedit langsung; harus versi baru. |
| BR-G-023 | Portal tidak boleh menjadi owner data transaksi bisnis. |
| BR-G-024 | Proses kritis wajib mencatat audit trail. |

# 10. Kebutuhan Bisnis Global

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan Bisnis |
| --- | --- | --- | --- |
| BRD-G-001 | P0 | Sistem harus menyediakan ERP/SIAKAD terintegrasi untuk lifecycle lead sampai alumni. | Data dapat ditelusuri dari CRM, PMB, Finance, Academic, LMS, Assessment, Portal. |
| BRD-G-002 | P0 | Sistem harus memakai Core sebagai pusat identity dan access. | Semua modul menggunakan token/claim Core. |
| BRD-G-003 | P0 | Setiap modul harus memiliki database fisik sendiri. | Setiap modul dapat backup, restore, migration, dan recovery terpisah. |
| BRD-G-004 | P0 | Sistem tidak boleh memakai cross-database FK. | ERD fisik hanya memiliki FK internal database; relasi lintas modul memakai external reference. |
| BRD-G-005 | P0 | Sistem tidak boleh memakai direct cross-database join pada transaksi online. | UI transaksi membaca API/snapshot/read model. |
| BRD-G-006 | P0 | Sistem harus memiliki outbox/inbox event per modul. | Event dapat dipublish, diterima, diproses, diretry, dan ditelusuri. |
| BRD-G-007 | P0 | Setiap proses kritis harus idempotent. | Retry tidak membuat applicant/student/payment/class/grade duplikat. |
| BRD-G-008 | P0 | Source of truth per domain harus jelas. | Modul lain tidak mengubah data owner secara langsung. |
| BRD-G-009 | P0 | Sistem harus mendukung degraded mode saat dependency down. | Modul tidak terkait tetap berjalan dan UI menampilkan data terakhir. |
| BRD-G-010 | P1 | Reporting lintas modul harus memakai read model/warehouse. | Dashboard pimpinan tidak query langsung ke semua database transaksi. |
| BRD-G-011 | P1 | Reconciliation lintas modul harus tersedia. | Selisih source vs snapshot/read model dapat dideteksi dan diperbaiki. |
| BRD-G-012 | P1 | Setiap dashboard/read model harus menampilkan waktu sinkronisasi terakhir. | User mengetahui freshness data. |

# 11. Source of Truth dan Ownership Data

| Domain Data | Source of Truth | Konsumen | Aturan Bisnis |
| --- | --- | --- | --- |
| Person/user/role/permission | Core | Semua modul | Modul lain simpan `person_ref_id`, `user_ref_id`, snapshot jika perlu. |
| Master data prodi/periode/status | Referensi | PMB, Finance, Academic, HRIS, LMS, Assessment, Portal | Modul lain simpan `*_ref_id` dan reference snapshot. |
| Lead/peminat | CRM | PMB, Portal, Reporting | PMB menerima lead qualified, bukan mengubah lead langsung. |
| Applicant | PMB | CRM, Finance, Academic, Assessment, Portal | Academic baru membuat student setelah PMB ready for handover. |
| Invoice/payment/clearance | Finance | PMB, Academic, Portal, Reporting | Status bayar dan clearance hanya dari Finance. |
| Student/KRS/final grade | Academic | Finance, LMS, Assessment, Portal, Reporting | LMS/Assessment hanya input, final grade milik Academic. |
| Dosen/pegawai | HRIS | Academic, LMS, Finance, Portal | Academic/LMS tidak membuat biodata dosen sendiri. |
| Kelas online/progress | LMS | Academic, Portal, Reporting | LMS mengikuti kelas dan peserta dari Academic. |
| Attempt/scoring | Assessment | PMB, LMS, Academic, Portal | Assessment result dikirim ke context owner. |
| Notification/preference | Portal | User dan modul sumber | Portal bukan source transaksi bisnis. |

# 12. Pola Relasi dan Join Data Bisnis

## 12.1 Prinsip Relasi Lintas Modul

Relasi lintas modul tidak diwujudkan sebagai FK database. Relasi diwujudkan melalui:

### 1. External reference ID.

### 2. Snapshot lokal.

### 3. Read model lokal.

### 4. API query/command.

### 5. Event contract.

### 6. Reconciliation job.

### 7. Warehouse untuk reporting.

## 12.2 Contoh Relasi Bisnis

| Relasi Bisnis | Implementasi Data | Cara Validasi |
| --- | --- | --- |
| Applicant memiliki person | `pmb_db.applicants.person_ref_id` | Core API atau `person_snapshots`. |
| Applicant berasal dari lead | `pmb_db.applicants.lead_ref_id` | CRM event/API. |
| Applicant punya status payment | `pmb_db.applicant_invoice_statuses` | Finance event `invoice_created/payment_paid`. |
| Student berasal dari applicant | `academic_db.students.applicant_ref_id` | PMB handover event dan idempotency key. |
| Student punya clearance KRS | `academic_db.student_clearance_snapshots` | Finance event `clearance_updated`. |
| LMS class berasal dari course offering | `lms_db.academic_class_snapshots.course_offering_ref_id` | Academic event `course_offering_created`. |
| LMS enrollment berasal dari KRS item | `lms_db.lms_enrollments.krs_item_ref_id` | Academic event `krs_approved`. |
| Quiz LMS memakai Assessment | `lms_db.quiz_activities.assessment_session_ref_id` | Assessment API/event. |
| Portal dashboard membaca semua modul | `portal_db.dashboard_read_models` | Event/read model/warehouse. |

## 12.3 Business Rule Join

| ID | Rule |
| --- | --- |
| JOIN-001 | Join lintas modul untuk transaksi online dilarang. |
| JOIN-002 | Join lokal di dalam satu database modul diperbolehkan. |
| JOIN-003 | Untuk data lintas modul yang sering dibaca, modul wajib memakai read model. |
| JOIN-004 | Untuk data lintas modul yang harus real-time, modul menggunakan API owner. |
| JOIN-005 | Untuk laporan agregat, sistem menggunakan warehouse/data mart. |
| JOIN-006 | UI wajib membedakan data real-time dan data snapshot dengan timestamp. |

# 13. BRD Per Modul

# 13.1 Modul Core

## Tujuan Bisnis

Core menjadi pusat identitas, akses, role, permission, session, application registry, service client, impersonation, audit global, idempotency, dan security foundation untuk seluruh modul.

## Stakeholder Utama

Super Admin, Admin BPPTI, Admin Modul, Pendaftar, Mahasiswa, Dosen, Pegawai, Pimpinan, Service Account antar modul.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Person dan user identity; SSO; active role; role, permission, data scope; application registry; service client; impersonation; audit global; token validation; idempotency; event outbox/inbox. | Data bisnis PMB, Finance, Akademik, HRIS, LMS, Assessment, Portal; pengaturan KRS/nilai/invoice/payment. |

## Proses Bisnis Utama

| Proses | Deskripsi | Output |
| --- | --- | --- |
| User onboarding | Admin membuat/menghubungkan person dengan user dan role. | User aktif dengan role dan scope. |
| Login dan role selection | User login melalui Core dan memilih active role. | Token/session membawa identity, active role, application, scope. |
| Role switching | User berpindah role sesuai hak. | Active role berubah dan audit tercatat. |
| Service authentication | Modul memakai service client/token. | Permintaan antar modul dapat divalidasi. |
| Impersonation | Admin tertentu masuk sebagai user lain untuk support. | Audit impersonation lengkap. |
| Token degraded mode | Modul memverifikasi JWT dari cache saat Core tidak tersedia. | Operasi terbatas tetap berjalan untuk token valid. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-CORE-001 | P0 | Core harus menjadi satu-satunya pusat login dan session. | Tidak ada modul operasional yang menyimpan credential/session sendiri. |
| BRD-CORE-002 | P0 | Core harus menyediakan active role session untuk user multi-role. | Menu, permission, data scope mengikuti active role. |
| BRD-CORE-003 | P0 | Core harus mengelola role, permission, dan data scope lintas modul. | Endpoint protected menolak akses tanpa permission/scope valid. |
| BRD-CORE-004 | P0 | Core harus menyediakan application registry dan app launcher. | User hanya melihat aplikasi sesuai active role. |
| BRD-CORE-005 | P0 | Core harus menyediakan service client/token. | Modul dapat memvalidasi service-to-service call. |
| BRD-CORE-006 | P0 | Core harus menyediakan audit global untuk aksi lintas modul. | Actor, role, timestamp, request id, reason tercatat. |
| BRD-CORE-007 | P0 | Core harus mendukung offline token validation terbatas. | JWT valid dapat diverifikasi lokal selama TTL/cache valid. |
| BRD-CORE-008 | P1 | Core harus mendukung impersonation terbatas. | Impersonation wajib reason dan mencatat actor asli. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| persons, users, roles, permissions, role_assignments, sessions, applications, service_clients, audit_logs, idempotency_keys, outbox_events, inbox_events. | Master data scope dari Referensi atau modul pemilik melalui external reference/snapshot. |

## Integrasi

| Integrasi | Kebutuhan | Guardrail |
| --- | --- | --- |
| Core ke semua modul | Token validation, active role, permission, app launcher. | Modul tidak membuat login sendiri. |
| Semua modul ke Core | User/person lookup, audit summary, service client validation. | Gunakan API/cache; tidak direct DB join. |

## Reporting dan Output

Daftar user aktif, role assignment, permission matrix, audit log sensitif, impersonation log, failed login, access denied, service client usage.

## UAT Starter

### 1. User multi-role login dan memilih active role berbeda.

### 2. Admin Prodi mencoba akses data prodi lain dan ditolak.

### 3. Service token invalid ditolak.

### 4. Impersonation tanpa reason ditolak.

### 5. Token valid tetap bisa diverifikasi saat Core API sementara tidak tersedia.

## Dependency dan Open Issue

Daftar role/permission final, struktur data scope, kebijakan MFA, masa berlaku session/token, TTL permission cache, kebijakan lockout.

# 13.2 Modul Referensi

## Tujuan Bisnis

Referensi menyediakan master data standar lintas modul agar dropdown, status, periode, prodi, jalur PMB, komponen pembayaran, dan kode bisnis konsisten.

## Stakeholder Utama

Admin Referensi, Admin Akademik Biro, Admin PMB, Admin Finance, Admin HRIS, Admin LMS, QA, DBA.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Wilayah, agama, kewarganegaraan, prodi, jenjang, tahun ajaran, periode akademik, jalur PMB, jenis dokumen, komponen pembayaran, metode pembayaran, status code. | Transaksi PMB, invoice, payment, KRS, nilai, attempt assessment. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-REF-001 | P0 | Referensi harus menyediakan master data lintas modul. | Dropdown/status lintas modul menggunakan kode standar. |
| BRD-REF-002 | P0 | Status bisnis kritis harus berupa managed code. | Transaksi tidak dapat menyimpan status tidak terdaftar. |
| BRD-REF-003 | P0 | Tahun Ajaran dan Periode Akademik harus tersedia lintas modul. | Gelombang, invoice, kelas, KRS, LMS, nilai, laporan punya konteks periode. |
| BRD-REF-004 | P1 | Master data yang sudah dipakai transaksi tidak boleh hard delete. | Tersedia inactive/archived. |
| BRD-REF-005 | P1 | Referensi harus publish event saat master berubah. | Modul consumer memperbarui snapshot lokal. |
| BRD-REF-006 | P1 | Referensi harus mendukung valid_from/valid_to untuk data tertentu. | Master periodik dapat dikontrol masa berlaku. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| regions, religions, academic_levels, study_programs, academic_years, academic_periods, pmb_paths, document_types, payment_components, payment_methods, status_codes. | Role/user dari Core untuk otorisasi pengelolaan master data. |

## Integrasi

Referensi publish master-data event ke PMB, Finance, Academic, HRIS, LMS, Assessment, Portal. Consumer menyimpan `reference_snapshots` atau cache. Tidak ada direct DB join.

## Reporting dan Output

Daftar master aktif/nonaktif, audit perubahan master, status code catalog, periode akademik per tahun ajaran, prodi aktif.

## UAT Starter

### 1. Membuat Periode Akademik tanpa Tahun Ajaran ditolak.

### 2. Menghapus master yang sudah dipakai transaksi ditolak.

### 3. PMB hanya menampilkan jalur PMB aktif.

### 4. Status bebas yang tidak terdaftar ditolak.

### 5. Perubahan nama prodi memicu event dan snapshot consumer diperbarui.

# 13.3 Modul CRM

## Tujuan Bisnis

CRM mengelola peminat sebelum menjadi applicant: campaign, lead, agen, referral, follow-up, pipeline, dan komisi.

## Stakeholder Utama

Owner CRM, Admin Marketing, Admin PMB, Agen/Mitra, Pimpinan, Finance untuk komisi.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Lead, campaign, agen, referral, follow-up, pipeline, commission eligibility, handover ke PMB. | Applicant biodata lengkap, dokumen PMB, invoice/payment resmi, generate NIM. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-CRM-001 | P0 | CRM harus mengelola lead dari berbagai sumber. | Lead memiliki source, campaign, status, dan follow-up. |
| BRD-CRM-002 | P0 | Lead qualified dapat dikonversi menjadi applicant PMB. | Conversion idempotent dan tidak membuat applicant duplikat. |
| BRD-CRM-003 | P1 | Agen hanya melihat lead/referral miliknya. | Agent scope ditegakkan backend. |
| BRD-CRM-004 | P1 | CRM harus menyimpan `applicant_ref_id` setelah handover berhasil. | CRM dapat menelusuri lead menjadi applicant tanpa FK lintas DB. |
| BRD-CRM-005 | P1 | CRM tetap dapat mencatat lead/follow-up saat PMB down. | Handover masuk retry/pending queue. |
| BRD-CRM-006 | P1 | CRM mendukung komisi/referral berdasarkan rule. | Komisi hanya eligible setelah kondisi bisnis terpenuhi. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| campaigns, agents, leads, follow_ups, lead_stage_histories, referrals, commission_rules, commissions, handover_to_pmb_logs. | Person/user dari Core via snapshot/API; applicant status dari PMB via event/API; payment status komisi dari Finance bila diperlukan. |

## Integrasi

CRM ke PMB melalui event/API `lead_qualified`. PMB mengembalikan `applicant_ref_id`. CRM tidak mengubah applicant langsung.

## Reporting dan Output

Lead funnel, conversion rate, campaign performance, agent performance, follow-up aging, komisi referral.

## UAT Starter

### 1. Lead baru dibuat dari campaign aktif.

### 2. Lead qualified dikirim dua kali ke PMB dan hanya satu applicant terbentuk.

### 3. Agen mencoba melihat lead agen lain dan ditolak.

### 4. PMB down saat handover; CRM tetap menyimpan status pending.

# 13.4 Modul PMB

## Tujuan Bisnis

PMB mengelola applicant, biodata, dokumen, seleksi, daftar ulang, LoA, dan handover ke Akademik.

## Stakeholder Utama

Pendaftar, Admin PMB, Owner PMB, Admin Finance, Admin Akademik, Admin Assessment, Pimpinan.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Wave, applicant, biodata, dokumen, seleksi, daftar ulang, LoA, handover, invoice/payment status read model. | Payment resmi, clearance resmi, NIM, KRS, nilai final. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-PMB-001 | P0 | PMB harus mengelola gelombang dengan target periode masuk. | Wave tidak dapat dibuka tanpa `target_entry_period_ref_id`. |
| BRD-PMB-002 | P0 | Applicant dapat mengisi biodata dan upload dokumen. | Status kelengkapan dan verifikasi dokumen tercatat. |
| BRD-PMB-003 | P0 | PMB harus membaca status pembayaran dari Finance. | Admin PMB tidak dapat mengubah status paid langsung. |
| BRD-PMB-004 | P0 | LoA hanya diterbitkan setelah syarat terpenuhi. | Seleksi, daftar ulang, payment wajib valid. |
| BRD-PMB-005 | P0 | Handover ke Akademik harus idempotent. | Request berulang tidak membuat mahasiswa ganda. |
| BRD-PMB-006 | P0 | PMB menyimpan invoice/payment read model. | PMB tetap menampilkan status terakhir saat Finance down. |
| BRD-PMB-007 | P1 | PMB menerima result dari Assessment. | Hasil CBT dapat menjadi dasar seleksi. |
| BRD-PMB-008 | P1 | PMB menyimpan person/reference snapshot. | Data dasar tetap tampil saat Core/Referensi down. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| waves, applicants, applicant_biodata, applicant_program_choices, applicant_documents, applicant_invoice_statuses, selections, re_registrations, loa_documents, handover_logs. | Person/user Core, master Referensi, lead CRM, invoice/payment Finance, assessment result, student Academic. |

## Integrasi

PMB menerima lead dari CRM, meminta invoice ke Finance, menerima payment event, meminta/menautkan session Assessment, dan mengirim handover ke Academic.

## Reporting dan Output

Applicant per wave, kelengkapan dokumen, status seleksi, status pembayaran read model, LoA issued, handover pending/success/failed.

## UAT Starter

### 1. Wave tanpa target period ditolak.

### 2. Applicant upload dokumen dan admin verifikasi.

### 3. LoA tanpa payment valid ditolak.

### 4. Finance down; PMB menampilkan payment status terakhir dengan timestamp.

### 5. Handover dua kali tidak membuat student duplikat.

# 13.5 Modul Finance

## Tujuan Bisnis

Finance menjadi source of truth invoice, payment, callback, manual verification, receipt, clearance, cicilan, beasiswa, jurnal dasar, dan laporan keuangan.

## Stakeholder Utama

Admin Finance, Pendaftar, Mahasiswa, Admin PMB, Admin Akademik, Pimpinan, Auditor.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Invoice, invoice item, payment, callback, manual verification, receipt, clearance, installment, scholarship, journal entry, customer snapshot. | Biodata applicant final, NIM generation, KRS, final grade. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-FIN-001 | P0 | Finance harus mengelola invoice pendaftar/mahasiswa. | Invoice memiliki item, due date, status, dan audit. |
| BRD-FIN-002 | P0 | Finance menjadi sumber status pembayaran dan clearance. | Modul lain membaca payment/clearance dari Finance/API/event. |
| BRD-FIN-003 | P0 | Callback payment gateway idempotent. | Callback berulang tidak membuat payment/jurnal ganda. |
| BRD-FIN-004 | P0 | Manual verification harus audit-ready. | Actor, timestamp, bukti, status, reason tercatat. |
| BRD-FIN-005 | P0 | Clearance mengendalikan layanan akademik. | KRS/ujian/KHS/transkrip/wisuda dapat clear/blocked/conditional. |
| BRD-FIN-006 | P1 | Finance mendukung cicilan dan beasiswa. | Tagihan dapat disesuaikan sesuai policy dan audit. |
| BRD-FIN-007 | P1 | Finance publish event invoice/payment/clearance. | PMB/Academic/Portal memperbarui read model. |
| BRD-FIN-008 | P1 | Finance tetap dapat memproses existing invoice saat PMB/Academic down. | Customer snapshot cukup untuk transaksi existing. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| customer_snapshots, invoices, invoice_items, payments, payment_callbacks, manual_verifications, receipts, clearances, installment_plans, scholarships, journal_entries, journal_lines. | Applicant PMB, student Academic, person Core, payment component/method Referensi. |

## Integrasi

Finance menerima request invoice dari PMB/Academic, publish invoice/payment/clearance event, dan memberi API query clearance real-time bila dibutuhkan.

## Reporting dan Output

Invoice aging, payment completion, manual verification queue, clearance blocked list, receipt report, journal basic, scholarship/discount report.

## UAT Starter

### 1. Invoice dibuat untuk applicant PMB.

### 2. Callback valid membuat payment dan receipt.

### 3. Callback sama dikirim ulang tidak membuat payment ganda.

### 4. Academic membaca clearance snapshot dari event Finance.

### 5. Manual verification tanpa bukti/reason ditolak sesuai policy.

# 13.6 Modul Akademik

## Tujuan Bisnis

Akademik menjadi source of truth student, NIM, kurikulum, mata kuliah, kelas, KRS, nilai final, KHS, transkrip, yudisium, dan alumni.

## Stakeholder Utama

Admin Akademik Biro, Admin Akademik Prodi, Kaprodi, Mahasiswa, Dosen, Dosen PA, Admin Finance, Admin LMS, Pimpinan.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Student, NIM, curriculum, course, curriculum course, period setting, course offering, schedule, lecturer assignment, KRS, grade final, KHS, transcript, graduation, alumni, clearance snapshot. | Payment resmi, LMS material/progress, assessment attempt. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-ACA-001 | P0 | Academic harus membuat student dan NIM dari handover PMB valid. | Student tidak dibuat tanpa applicant ready dan syarat Finance. |
| BRD-ACA-002 | P0 | NIM generation harus unik dan idempotent. | Handover berulang tidak membuat NIM/mahasiswa ganda. |
| BRD-ACA-003 | P0 | Student wajib menyimpan entry period dan curriculum. | Histori kurikulum mahasiswa stabil. |
| BRD-ACA-004 | P0 | Academic mengelola kurikulum dan mata kuliah. | Kurikulum dapat versioned dan tidak merusak data lama. |
| BRD-ACA-005 | P0 | Course offering wajib berada pada periode akademik. | Kelas tanpa `academic_period_ref_id` ditolak. |
| BRD-ACA-006 | P0 | KRS harus memvalidasi clearance Finance. | KRS final ditolak jika clearance blocked. |
| BRD-ACA-007 | P0 | LMS enrollment berasal dari KRS approved. | Academic publish KRS/class event ke LMS. |
| BRD-ACA-008 | P0 | Final grade hanya milik Academic. | LMS/Assessment input tidak menimpa final grade otomatis. |
| BRD-ACA-009 | P1 | KHS/transkrip/wisuda dapat dibatasi clearance. | Layanan ditolak/conditional sesuai Finance policy. |
| BRD-ACA-010 | P1 | Academic menyimpan reference/person/clearance snapshot. | Operasi tertentu tetap berjalan saat dependency down. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| students, nim_sequences, curriculums, courses, curriculum_courses, course_offerings, course_schedules, course_lecturers, krs_headers, krs_items, grades, grade_inputs, khs, transcripts, graduations, alumni. | PMB applicant, Finance clearance, HRIS lecturer, Reference period/prodi, Core person/user, LMS/Assessment grade input. |

## Integrasi

Academic menerima handover PMB, menerima clearance Finance, membaca dosen HRIS, publish class/KRS event ke LMS, menerima grade input LMS/Assessment.

## Reporting dan Output

Mahasiswa aktif, NIM generated, kelas per periode, KRS status, dosen pengampu, nilai final, KHS, transkrip, alumni, clearance blocking report.

## UAT Starter

### 1. Generate NIM tanpa handover valid ditolak.

### 2. Handover sama dua kali tidak membuat student ganda.

### 3. KRS final blocked saat clearance blocked.

### 4. LMS grade input masuk sebagai input, bukan final grade.

### 5. Perubahan final grade wajib reason dan audit.

# 13.7 Modul HRIS/SDM

## Tujuan Bisnis

HRIS menjadi source of truth pegawai, dosen, homebase, jabatan, unit kerja, status aktif, BKD, performance, sertifikasi, dan payroll source.

## Stakeholder Utama

Admin SDM, Dosen, Pegawai, Admin Akademik, Admin LMS, Finance Payroll, Pimpinan.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Employee, lecturer, work unit, position, employment status, employment record, homebase history, BKD, performance, certification, payroll source. | Kelas akademik, KRS, LMS material, payroll processing final. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-HRIS-001 | P0 | HRIS menjadi sumber data dosen dan pegawai. | Academic/LMS tidak membuat biodata dosen sendiri. |
| BRD-HRIS-002 | P0 | HRIS mengelola status aktif dosen. | Dosen nonaktif tidak dapat diplot ke kelas baru. |
| BRD-HRIS-003 | P1 | HRIS mengelola homebase dan histori homebase. | Perubahan homebase tidak menghapus histori lama. |
| BRD-HRIS-004 | P1 | HRIS publish event dosen/pegawai update. | Academic/LMS memperbarui lecturer snapshot. |
| BRD-HRIS-005 | P1 | HRIS menyediakan data payroll source. | Finance dapat membaca payroll source sesuai scope. |
| BRD-HRIS-006 | P1 | HRIS tetap berjalan saat Academic/LMS down. | Update dosen tersimpan dan event dikirim saat dependency pulih. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| employees, lecturers, work_units, positions, employment_statuses, employment_records, homebase_histories, bkd_records, performance_records, certifications, payroll_sources. | Person Core, study program Reference. |

## Integrasi

HRIS publish lecturer/employee event ke Academic, LMS, Finance, Portal. Academic dan LMS memakai `lecturer_ref_id` + snapshot.

## Reporting dan Output

Daftar pegawai, daftar dosen aktif, homebase dosen, BKD per periode, payroll source, sertifikasi, status aktif/nonaktif.

## UAT Starter

### 1. Membuat dosen dari person Core.

### 2. Dosen nonaktif ditolak saat plotting kelas.

### 3. Perubahan homebase menyimpan histori.

### 4. LMS menerima update dosen dari event HRIS.

# 13.8 Modul LMS

## Tujuan Bisnis

LMS mengelola pembelajaran online berdasarkan kelas akademik dan enrollment KRS valid.

## Stakeholder Utama

Dosen, Mahasiswa, Admin LMS, Admin Akademik, Admin Assessment, Pimpinan.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| LMS class, enrollment dari KRS, session, material, assignment, submission, attendance, progress, quiz activity, grade input. | Kelas akademik source, KRS approval, final grade, transkrip. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-LMS-001 | P0 | LMS class harus berasal dari course offering Academic. | LMS tidak dapat membuat kelas akademik sendiri. |
| BRD-LMS-002 | P0 | LMS enrollment harus berasal dari KRS approved. | Mahasiswa tidak dapat self-enroll. |
| BRD-LMS-003 | P1 | LMS mengelola sesi, materi, tugas, presensi, progress. | Dosen hanya mengelola kelas yang diampu. |
| BRD-LMS-004 | P1 | LMS dapat membuat quiz activity melalui Assessment. | Quiz memiliki `assessment_session_ref_id`. |
| BRD-LMS-005 | P0 | LMS grade input tidak menimpa final grade. | Grade input dikirim ke Academic sebagai source input. |
| BRD-LMS-006 | P1 | LMS menyimpan class/student/lecturer snapshot. | Kelas tetap tampil dengan data terakhir saat Academic/HRIS down. |
| BRD-LMS-007 | P1 | LMS down tidak menghentikan KRS Academic. | Academic tetap dapat mengelola KRS walau LMS sync pending. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| lms_classes, lms_enrollments, lms_lecturers, sessions, materials, assignments, submissions, attendance, progress, quiz_activities, grade_inputs. | Academic class/KRS/student, HRIS lecturer, Assessment session/result, Core user. |

## Integrasi

LMS menerima class/KRS event dari Academic, lecturer event dari HRIS, membuat quiz ke Assessment, dan mengirim grade input ke Academic.

## Reporting dan Output

Kelas online aktif, enrollment LMS, materi tersedia, tugas/pengumpulan, presensi, progress mahasiswa, grade input status.

## UAT Starter

### 1. Academic membuat kelas dan LMS menerima class sync.

### 2. Mahasiswa tanpa KRS approved tidak muncul di kelas LMS.

### 3. Dosen hanya melihat kelas yang diampu.

### 4. Grade input LMS dikirim ke Academic dan tidak menjadi final otomatis.

# 13.9 Modul Assessment

## Tujuan Bisnis

Assessment menyediakan mesin assessment reusable untuk CBT PMB, quiz LMS, survey, dan scoring lain dengan bank soal berversi.

## Stakeholder Utama

Admin Assessment, Admin PMB, Admin LMS, Dosen, Pendaftar, Mahasiswa, Pimpinan.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Question bank, question version, assessment session, session question, participant snapshot, attempt, answer, scoring result, result export. | Status seleksi final PMB, final grade Academic, LMS material. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-ASM-001 | P0 | Assessment harus mendukung bank soal dan versioning. | Soal yang sudah dipakai tidak diedit langsung. |
| BRD-ASM-002 | P0 | Assessment reusable untuk CBT PMB, quiz LMS, survey. | Context menentukan consumer result. |
| BRD-ASM-003 | P0 | Attempt dan jawaban harus audit-ready. | Peserta, waktu, jawaban, score, status tercatat. |
| BRD-ASM-004 | P0 | Scoring result harus idempotent. | Result event berulang tidak membuat score duplikat. |
| BRD-ASM-005 | P1 | Assessment menyimpan participant snapshot. | Attempt tetap bisa ditelusuri saat PMB/LMS down. |
| BRD-ASM-006 | P1 | Assessment mengirim result ke context owner. | PMB/LMS menerima result sesuai session context. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| question_banks, questions, question_versions, assessment_sessions, session_questions, participant_snapshots, attempts, attempt_answers, scoring_results, result_exports. | Participant dari PMB/LMS/Core, context owner PMB/LMS, result consumer Academic bila relevan. |

## Integrasi

Assessment menerima request session dari PMB/LMS, menerima participant snapshot, publish result event, dan menyediakan result API.

## Reporting dan Output

Daftar bank soal, versi soal, session aktif, attempt completion, score distribution, export status.

## UAT Starter

### 1. Membuat soal versi pertama.

### 2. Soal sudah dipakai attempt tidak dapat diedit langsung.

### 3. CBT PMB menghasilkan result ke PMB.

### 4. Quiz LMS menghasilkan result ke LMS.

### 5. Event result duplikat tidak membuat result ganda.

# 13.10 Modul Portal

## Tujuan Bisnis

Portal menjadi presentation layer, dashboard role-based, notification center, preference, shortcut, dan activity log.

## Stakeholder Utama

Pendaftar, Mahasiswa, Dosen, Admin Modul, Pimpinan, Product Owner.

## Scope Bisnis

| In Scope | Out of Scope |
| --- | --- |
| Dashboard role-based, notifications, user preferences, shortcuts, activity logs, dashboard read models, user/role snapshots. | Source of truth transaksi PMB, Finance, Academic, HRIS, LMS, Assessment. |

## Kebutuhan Bisnis Modul

| ID | Prioritas | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| BRD-POR-001 | P1 | Portal menampilkan dashboard berbasis active role. | User hanya melihat widget sesuai role dan scope. |
| BRD-POR-002 | P1 | Portal menjadi notification center. | Notifikasi dikirim ke user sesuai event dan role. |
| BRD-POR-003 | P1 | Portal menyimpan preference dan shortcut user/role. | Shortcut mengikuti active role. |
| BRD-POR-004 | P1 | Portal tidak menjadi source transaksi bisnis. | Portal tidak mengubah payment/KRS/nilai langsung. |
| BRD-POR-005 | P1 | Dashboard memakai read model/summary API. | Widget memiliki `refreshed_at`. |
| BRD-POR-006 | P1 | Portal down tidak menghentikan modul sumber. | PMB/Finance/Academic/LMS tetap berjalan. |

## Data Ownership

| Data Dimiliki | Data Dibaca/Dikonsumsi |
| --- | --- |
| notification_events, notifications, user_preferences, shortcuts, activity_logs, dashboard_read_models, dashboard_widgets, user_snapshots, role_snapshots. | Event/status dari semua modul, user/role Core. |

## Integrasi

Portal menerima event dari semua modul untuk notification/read model, membaca Core role/user snapshot, dan menyediakan dashboard berbasis role.

## Reporting dan Output

Unread notification, user activity, dashboard KPI, role shortcut list, read model freshness.

## UAT Starter

### 1. Mahasiswa menerima notifikasi payment paid.

### 2. Dosen menerima notifikasi kelas LMS.

### 3. Pimpinan melihat dashboard read-only.

### 4. Portal down tidak memblokir transaksi Finance.

# 14. Kebutuhan Integrasi Lintas Modul

| ID | Prioritas | Integrasi | Kebutuhan Bisnis | Guardrail |
| --- | --- | --- | --- | --- |
| INT-001 | P0 | Core ke semua modul | Token, role, permission, scope. | Modul tidak membuat login sendiri. |
| INT-002 | P0 | Referensi ke semua modul | Master data dan status code. | Consumer memakai snapshot/cache/event. |
| INT-003 | P0 | CRM ke PMB | Lead qualified menjadi applicant. | Conversion idempotent. |
| INT-004 | P0 | PMB ke Finance | Request invoice dan baca payment status. | Payment status hanya dari Finance. |
| INT-005 | P0 | Finance ke PMB | Invoice/payment event. | PMB read model bukan source. |
| INT-006 | P0 | PMB ke Assessment | CBT selection session. | Result dikembalikan sesuai context. |
| INT-007 | P0 | PMB ke Academic | Handover applicant ready. | Student/NIM dibuat oleh Academic. |
| INT-008 | P0 | Finance ke Academic | Clearance event/API. | KRS/layanan akademik mengikuti clearance. |
| INT-009 | P0 | HRIS ke Academic/LMS | Lecturer active status dan homebase. | Dosen nonaktif tidak diplot. |
| INT-010 | P0 | Academic ke LMS | Class and KRS enrollment sync. | LMS tidak self-enroll. |
| INT-011 | P0 | LMS ke Academic | Grade input. | Final grade tetap milik Academic. |
| INT-012 | P0 | Assessment ke PMB/LMS | Result event/API. | Consumer memproses idempotent. |
| INT-013 | P1 | Semua modul ke Portal | Notification/dashboard event. | Portal tidak menjadi source transaksi. |
| INT-014 | P1 | Semua modul ke Reporting | Data warehouse/read model. | Tidak query langsung DB produksi lintas modul. |

# 15. Reporting dan Dashboard

| Kategori Laporan | Modul Sumber | Output Bisnis | Catatan |
| --- | --- | --- | --- |
| PMB Funnel | CRM + PMB | Lead, applicant, submitted, accepted, LoA, handover. | Dibangun via warehouse/read model. |
| Finance Collection | Finance | Invoice issued, paid, overdue, outstanding. | Source Finance. |
| Payment vs Applicant | PMB + Finance | Applicant paid/unpaid per wave/prodi. | PMB read model atau warehouse. |
| Student Intake | PMB + Academic | Applicant handover vs student created. | Reconciliation critical. |
| KRS Monitoring | Academic + Finance | KRS submitted/approved/finalized vs clearance. | Academic read model clearance. |
| LMS Engagement | LMS + Academic | Class active, attendance, progress, submission. | LMS source, Academic class context. |
| Assessment Result | Assessment + PMB/LMS | CBT/quiz completion and score. | Context-specific. |
| Lecturer Load | HRIS + Academic + LMS | Dosen aktif, kelas diampu, BKD. | HRIS source. |
| Executive KPI | Semua modul | Dashboard pimpinan. | Portal/warehouse, timestamp wajib. |

# 16. Non-Functional Business Requirement

| ID | Kategori | Kebutuhan Bisnis | Kriteria Penerimaan |
| --- | --- | --- | --- |
| NFBR-001 | Availability | Modul kritikal harus memiliki rencana availability. | Core, Finance, Academic memiliki HA/recovery plan. |
| NFBR-002 | Reliability | Event tidak hilang saat failure sementara. | Outbox pending dapat dipublish ulang. |
| NFBR-003 | Resilience | Modul tetap berjalan dalam degraded mode. | UI menampilkan data terakhir dan batasan operasi. |
| NFBR-004 | Performance | Query transaksi utama memakai database lokal/read model. | Tidak ada online direct cross-database join. |
| NFBR-005 | Security | RBAC dan data scope ditegakkan backend. | Akses lintas scope ditolak. |
| NFBR-006 | Auditability | Aksi sensitif audit-ready. | Actor, role, timestamp, reason, old/new value tercatat. |
| NFBR-007 | Data Integrity | External reference dapat direkonsiliasi. | Mismatch dapat dideteksi dan diperbaiki. |
| NFBR-008 | Observability | Request/event lintas modul dapat ditelusuri. | Correlation id dan event key tersedia. |
| NFBR-009 | Backup/Restore | Database per modul dapat dipulihkan terpisah. | Restore drill tersedia sebelum production. |
| NFBR-010 | Compliance | Data pribadi dan transaksi sensitif dilindungi. | Akses, audit, dan retention sesuai policy. |

# 17. Failure Scenario Business Requirement

| Scenario | Dampak yang Diterima | Modul yang Tetap Harus Jalan | Kebutuhan UI/Operasional |
| --- | --- | --- | --- |
| `finance_db` down | Invoice baru, payment verification, clearance real-time tertunda. | CRM, PMB biodata/dokumen, Academic baca clearance snapshot, LMS, Assessment non-Finance. | Tampilkan payment/clearance terakhir + timestamp. |
| `academic_db` down | Generate NIM, KRS, final grade, KHS, transkrip tertunda. | CRM, PMB pendaftaran, Finance payment existing, HRIS, Assessment CBT PMB. | Handover masuk pending/retry. |
| `lms_db` down | Pembelajaran online dan progress LMS tertunda. | Academic KRS/kelas, PMB, Finance, HRIS, Assessment non-LMS. | Academic tidak menunggu LMS untuk KRS. |
| `core_db` down | Login baru dan role switch terganggu. | Modul dapat memverifikasi token valid dengan cache terbatas. | Tampilkan keterbatasan login/role switch. |
| `portal_db` down | Dashboard dan notification center tidak tersedia. | Semua transaksi sumber tetap berjalan. | Event notifikasi pending. |
| Event broker down | Sinkronisasi read model tertunda. | Transaksi lokal tetap commit dan event tertahan di outbox. | Monitoring outbox backlog. |

# 18. Asumsi, Dependency, Risiko, dan Mitigasi

## 18.1 Asumsi

| ID | Asumsi |
| --- | --- |
| ASM-001 | Setiap modul dapat memiliki database fisik sendiri. |
| ASM-002 | Semua modul memakai UUID global sebagai primary/external reference. |
| ASM-003 | Event broker tersedia untuk integrasi asinkron. |
| ASM-004 | API Gateway/service discovery tersedia untuk integrasi sinkron. |
| ASM-005 | Core menyediakan token yang dapat diverifikasi lokal. |
| ASM-006 | Setiap modul memiliki owner bisnis dan technical owner. |
| ASM-007 | Data warehouse/read model tersedia untuk reporting lintas modul. |

## 18.2 Dependency

| Dependency | Keterangan | Risiko Jika Tidak Tersedia |
| --- | --- | --- |
| Core SSO/RBAC | Fondasi akses seluruh modul. | Modul membuat auth sendiri. |
| Reference master | Fondasi periode/prodi/status. | Transaksi tidak konsisten. |
| Event broker | Fondasi sync read model. | Integrasi menjadi terlalu sinkron dan rapuh. |
| API contract | Fondasi command/query lintas modul. | Integrasi tidak stabil. |
| Idempotency standard | Fondasi retry aman. | Data ganda. |
| Reconciliation job | Fondasi konsistensi snapshot. | Mismatch tidak terdeteksi. |
| Observability | Fondasi tracing dan audit teknis. | Incident sulit ditelusuri. |

## 18.3 Risiko dan Mitigasi

| ID | Risiko | Dampak | Mitigasi |
| --- | --- | --- | --- |
| RSK-001 | Database per modul masih berada dalam satu server fisik. | Failure server tetap menjatuhkan semua DB. | Pisahkan cluster/instance untuk modul kritikal. |
| RSK-002 | Event consumer tidak idempotent. | Data dobel. | Unique event_key dan inbox processing rule. |
| RSK-003 | Snapshot stale dipakai untuk keputusan kritis. | Keputusan bisnis salah. | Timestamp, stale policy, real-time API fallback. |
| RSK-004 | Reporting query langsung ke DB transaksi. | Beban tinggi dan dependency rapat. | Warehouse/read model. |
| RSK-005 | Handover applicant tidak terkunci. | Student/NIM ganda. | Idempotency, unique applicant_ref_id di Academic, NIM sequence lock. |
| RSK-006 | Payment callback dobel. | Payment/jurnal ganda. | Provider event id unique, idempotency, callback log. |
| RSK-007 | RBAC hanya di UI. | Kebocoran data. | Backend permission/scope test. |
| RSK-008 | Tahun Ajaran dan Tahun Kurikulum rancu. | Salah periode dan transkrip. | Validasi form, label UI, data model terpisah. |
| RSK-009 | Migration data lama kotor. | Constraint gagal dan data dobel. | Pre-flight cleansing dan dry-run. |
| RSK-010 | Core down terlalu lama. | Login dan akses terganggu. | HA Core, token cache, disaster recovery. |

# 19. Release Bisnis dan MVP

| Release | Fokus Bisnis | Modul | Acceptance Business Milestone |
| --- | --- | --- | --- |
| R0 Foundation | Identity, master data, event foundation, DB per modul, audit, idempotency. | Core, Referensi, Infrastructure | Login, RBAC, master period/prodi, outbox/inbox baseline siap. |
| R1 PMB Basic | Lead dan applicant registration. | CRM, PMB, Core, Referensi | Lead/applicant/biodata/dokumen berjalan. |
| R2 Finance Basic | Invoice/payment/receipt/clearance basic. | Finance, PMB | PMB dapat request invoice dan membaca payment status. |
| R3 Academic Onboarding | Handover dan generate NIM. | PMB, Academic, Finance | Applicant valid menjadi student/NIM tanpa duplikasi. |
| R4 KRS dan LMS Sync | KRS, class offering, LMS class/enrollment. | Academic, HRIS, LMS | KRS approved muncul di LMS. |
| R5 Assessment | CBT PMB dan quiz LMS. | Assessment, PMB, LMS | Result assessment terkirim ke context owner. |
| R6 Grade/KHS/Transcript | Final grade, KHS, transkrip. | Academic, LMS, Assessment, Finance | Grade final dan layanan akademik terkontrol clearance. |
| R7 Portal/Reporting | Dashboard, notification, KPI. | Portal, Reporting, semua modul | Dashboard role-based dengan sync timestamp. |
| R8 Hardening | Failure handling, reconciliation, migration, performance. | Semua modul | UAT failure scenario dan reconciliation lulus. |

# 20. Acceptance Criteria Global dan UAT Starter

## 20.1 Acceptance Criteria Global

### 1. Semua user login melalui Core.

### 2. User multi-role dapat memilih active role dan scope berubah sesuai role.

### 3. Tidak ada cross-database FK pada ERD fisik.

### 4. Tidak ada direct cross-database join pada transaksi online.

### 5. Lead qualified tidak membuat applicant duplikat saat retry.

### 6. Applicant handover tidak membuat student/NIM duplikat saat retry.

### 7. Payment callback duplikat tidak membuat payment/jurnal ganda.

### 8. PMB tidak dapat mengubah status paid langsung.

### 9. Academic tidak dapat finalize KRS jika clearance blocked.

### 10. LMS tidak dapat membuat kelas akademik sendiri.

### 11. LMS enrollment hanya berasal dari KRS approved.

### 12. LMS/Assessment grade input tidak menimpa final grade.

### 13. Portal tidak dapat mengubah transaksi sumber.

### 14. Finance down tidak menghentikan input biodata PMB.

### 15. LMS down tidak menghentikan KRS Academic.

### 16. Portal down tidak menghentikan transaksi sumber.

### 17. Read model/dashboard menampilkan timestamp sinkronisasi.

### 18. Reconciliation job dapat mendeteksi mismatch source vs snapshot.

## 20.2 UAT Starter Matrix

| ID | Scenario | Modul | Expected Result |
| --- | --- | --- | --- |
| UAT-G-001 | User login dan pilih active role. | Core | Menu/scope sesuai active role. |
| UAT-G-002 | Admin Prodi akses prodi lain via API. | Core/Academic | Ditolak. |
| UAT-G-003 | Lead qualified dikirim ulang. | CRM/PMB | Hanya satu applicant. |
| UAT-G-004 | Applicant bayar dan callback provider duplikat. | Finance/PMB | Payment satu, status PMB update via event. |
| UAT-G-005 | Applicant ready handover dikirim dua kali. | PMB/Academic | Hanya satu student dan NIM. |
| UAT-G-006 | Finance down saat applicant detail dibuka. | PMB | Status pembayaran terakhir tampil dengan timestamp. |
| UAT-G-007 | Clearance blocked saat mahasiswa final KRS. | Finance/Academic | KRS final ditolak. |
| UAT-G-008 | Academic class sync duplikat ke LMS. | Academic/LMS | LMS class tidak dobel. |
| UAT-G-009 | Dosen nonaktif diplot ke kelas. | HRIS/Academic | Ditolak. |
| UAT-G-010 | LMS grade input dikirim ke Academic. | LMS/Academic | Masuk grade input, bukan final grade otomatis. |
| UAT-G-011 | Soal digunakan attempt lalu diedit. | Assessment | Edit langsung ditolak, harus versi baru. |
| UAT-G-012 | Portal down saat Finance payment. | Finance/Portal | Payment tetap sukses, notifikasi pending. |
| UAT-G-013 | Event broker down saat transaksi lokal. | Semua modul | Transaksi commit, event pending di outbox. |

# 21. Dokumen Turunan

| Dokumen | Tujuan | Owner |
| --- | --- | --- |
| FSD Global dan Per Modul | Detail fungsi, form, flow, state machine, validation. | System Analyst |
| ERD/DBML per Database | Struktur tabel lokal, FK internal, external reference, index. | DBA |
| API Contract | Endpoint, request, response, error code, permission, idempotency. | Technical Lead/Backend Lead |
| Event Contract | Event name, payload, version, owner, consumer, retry, DLQ. | Technical Lead |
| RBAC Matrix | Role, permission, menu, endpoint, scope. | System Analyst/Core Owner |
| UAT Scenario | Test case bisnis normal, negatif, failure, retry. | QA/UAT Lead |
| Migration Mapping | Mapping data lama, cleansing, dry-run, rollback. | DBA/Owner Data |
| Reconciliation Spec | Aturan source vs snapshot/read model. | DBA/Technical Lead |
| Operations Runbook | Backup, restore, failover, monitoring, incident. | DevOps/SRE |

# 22. Approval

| Area | Owner | Status Review | Tanggal | Catatan |
| --- | --- | --- | --- | --- |
| Product Owner | - | Belum direview | - | - |
| Core | Admin BPPTI/Core Owner | Belum direview | - | - |
| Referensi | Admin Referensi/Admin Akademik | Belum direview | - | - |
| CRM | Owner CRM/Marketing PMB | Belum direview | - | - |
| PMB | Owner PMB/Admin PMB | Belum direview | - | - |
| Finance | Owner Finance/Admin Keuangan | Belum direview | - | - |
| Akademik | Biro Akademik/Prodi | Belum direview | - | - |
| HRIS/SDM | Owner HRIS/Admin SDM | Belum direview | - | - |
| LMS | Owner LMS/Admin LMS | Belum direview | - | - |
| Assessment | Owner Assessment/Admin Assessment | Belum direview | - | - |
| Portal | Owner Portal | Belum direview | - | - |
| DBA | DBA Lead | Belum direview | - | - |
| Technical Lead | Technical Lead | Belum direview | - | - |
| QA/UAT | QA/UAT Lead | Belum direview | - | - |
| DevOps/SRE | SRE Lead | Belum direview | - | - |

# 23. Penutup

BRD v1.1 ini menetapkan kebutuhan bisnis UNSIA setelah keputusan arsitektur database berubah menjadi physical database per modul. Fokusnya bukan hanya memisahkan database, tetapi memastikan proses bisnis tetap konsisten melalui source of truth tunggal, external reference, API/event contract, snapshot, read model, audit, idempotency, reconciliation, dan graceful degradation.

Dokumen ini harus divalidasi oleh Product Owner, Owner Modul, Technical Lead, DBA, QA/UAT, dan SRE sebelum diturunkan ke FSD, ERD/DBML, API Contract, Event Contract, RBAC Matrix, UAT, migration mapping, dan release plan.

# Appendix A - Event Contract Standard

Appendix ini melengkapi BRD UNSIA v1.1 sebagai standar kebutuhan event lintas modul. Tujuannya adalah memastikan proses bisnis yang berjalan melalui API, outbox/inbox, snapshot, read model, dan reconciliation memiliki kontrak event yang jelas, konsisten, teruji, dan dapat diaudit. Standar ini menjadi dasar turunan untuk FSD, API Contract, Event Contract Catalog, ERD/DBML, UAT, observability, dan release plan.

## A.1 Tujuan Bisnis Event Contract

| ID | Tujuan Bisnis | Kriteria Penerimaan |
| --- | --- | --- |
| EC-B-001 | Menjamin proses lintas modul tetap konsisten tanpa cross-database FK. | Relasi lintas modul diproses melalui event, external reference, snapshot, dan reconciliation. |
| EC-B-002 | Mencegah data ganda saat retry terjadi. | Semua event consumer wajib idempotent dan mencatat event_key. |
| EC-B-003 | Menjaga layanan tetap berjalan saat dependency down. | Consumer memakai snapshot/read model dengan informasi waktu sinkronisasi. |
| EC-B-004 | Memastikan audit bisnis tersedia. | Publish, consume, retry, DLQ, replay, dan reconciliation wajib tercatat. |
| EC-B-005 | Memastikan reporting lintas modul tidak membebani transaksi. | Dashboard memakai read model atau warehouse, bukan join langsung ke database produksi. |

## A.2 Struktur Wajib Event Identity

| Field | Deskripsi |
| --- | --- |
| event_name | Nama event dengan format domain.action, contoh finance.payment_paid. |
| event_version | Versi schema event, contoh v1. |
| event_key | Kunci unik global untuk duplicate handling dan idempotency. |
| event_type | DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, atau SNAPSHOT_EVENT. |
| publisher_service | Service yang menerbitkan event. |
| publisher_database | Database source of truth pemilik event. |
| aggregate_type | Objek bisnis utama, contoh payment, invoice, applicant, student, krs. |
| aggregate_id | ID objek bisnis utama pada database pemilik. |
| correlation_id | ID untuk melacak satu proses bisnis end-to-end. |
| causation_id | ID command atau event yang memicu event saat ini. |
| occurred_at | Waktu kejadian bisnis terjadi. |
| published_at | Waktu event berhasil dikirim ke broker. |

Contoh event envelope:

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

## A.3 Trigger Bisnis dan Perubahan Status

| Komponen | Penjelasan |
| --- | --- |
| business_trigger | Kondisi bisnis yang menyebabkan event diterbitkan. |
| pre_condition | Status atau kondisi data sebelum event terjadi. |
| post_condition | Status atau kondisi data setelah event terjadi. |
| source_table | Tabel utama yang menjadi sumber perubahan. |
| state_transition | Perubahan status, contoh UNPAID menjadi PAID. |
| publish_timing | Waktu event boleh dipublish, umumnya setelah commit transaksi lokal. |

## A.4 Publisher, Consumer, dan Tujuan Konsumsi

| Event | Publisher | Consumer | Tujuan Konsumsi |
| --- | --- | --- | --- |
| core.person_updated | Core | CRM, PMB, Academic, HRIS, LMS, Assessment, Portal | Update person snapshot lokal. |
| pmb.applicant_created | PMB | Finance, Assessment, Portal, Reporting | Membuat konteks applicant dan dashboard. |
| finance.invoice_created | Finance | PMB, Academic, Portal, Reporting | Update status tagihan. |
| finance.payment_paid | Finance | PMB, Academic, Portal, Reporting | Update status bayar, clearance snapshot, notifikasi, dan laporan. |
| finance.clearance_changed | Finance | PMB, Academic, LMS, Portal | Mengatur kelayakan layanan akademik. |
| academic.student_created | Academic | PMB, LMS, Portal, Reporting | Menghubungkan applicant dengan student/NIM. |
| academic.krs_approved | Academic | LMS, Portal, Reporting | Membuat atau memperbarui enrollment LMS. |
| assessment.result_calculated | Assessment | PMB, LMS, Academic, Portal | Mengirim hasil assessment ke context owner. |

## A.5 Payload Schema dan Validation Rule

Setiap event wajib memiliki payload schema. Payload harus cukup untuk kebutuhan consumer, tetapi tidak membawa data pribadi berlebihan. Untuk data sensitif, gunakan external_ref_id dan snapshot minimum.

| Field | Required | Validation Rule |
| --- | --- | --- |
| payment_id | Ya | UUID dan harus ada pada finance_db.payments. |
| invoice_id | Ya | UUID dan harus ada pada finance_db.invoices. |
| invoice_no | Ya | String dan tidak kosong. |
| bill_to_type | Ya | APPLICANT, STUDENT, atau PERSON. |
| bill_to_ref_id | Ya | UUID external reference subject pembayaran. |
| paid_amount | Ya | Decimal lebih dari 0. |
| payment_method_code | Ya | Kode metode pembayaran. |
| paid_at | Ya | Datetime valid. |
| status_code | Ya | Harus PAID untuk finance.payment_paid. |

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

## A.6 Idempotency dan Duplicate Handling

| Area | Aturan |
| --- | --- |
| Event key | event_key harus unik dan deterministik. Format disarankan: {event_name}:{aggregate_id}:{event_version}. |
| Consumer inbox | Consumer wajib menyimpan event_key pada inbox_events. |
| Duplicate event | Jika event_key sudah pernah diproses, consumer tidak memproses ulang payload. |
| Retry command | Command/API yang memicu event wajib membawa idempotency_key. |
| Conflict payload | Jika event_key sama tetapi payload berbeda, consumer menolak event dan mencatat mismatch. |

## A.7 Ordering, Dependency, dan Causality

| Event | Wajib Setelah | Catatan |
| --- | --- | --- |
| finance.payment_paid | finance.invoice_created | Payment tidak valid tanpa invoice. |
| finance.clearance_changed | finance.payment_paid atau finance.clearance_reviewed | Clearance berubah berdasarkan status finance resmi. |
| academic.student_created | pmb.ready_for_academic atau pmb.handover_requested | Mahasiswa dibuat setelah applicant siap diserahkan ke Academic. |
| lms.enrollment_synced | academic.krs_approved | Enrollment LMS berasal dari KRS valid. |
| academic.final_grade_published | academic.grade_input_received | LMS/Assessment hanya memberi input, final grade tetap di Academic. |

## A.8 Retry Policy dan Dead Letter Queue

| Komponen | Aturan |
| --- | --- |
| Retry schedule | 1 menit, 5 menit, 15 menit, lalu exponential backoff. |
| Maksimal retry | 10 kali atau sesuai SLA modul. |
| Temporary failure | Event tetap berada pada retry queue. |
| Permanent failure | Event masuk DLQ setelah retry maksimum. |
| DLQ payload | Wajib menyimpan event_key, consumer, last_error, retry_count, failed_at, dan raw_payload. |
| Recovery | Event dari DLQ dapat direplay manual oleh role DevOps/SRE yang berwenang. |

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

## A.9 Snapshot dan Read Model Impact

| Event | Consumer | Tabel Lokal yang Diupdate |
| --- | --- | --- |
| core.person_updated | PMB, Academic, LMS | person_snapshots |
| reference.study_program_updated | PMB, Academic, HRIS, LMS | reference_snapshots atau study_program_snapshots |
| finance.invoice_created | PMB, Academic, Portal | applicant_invoice_statuses, student_finance_snapshots, dashboard_read_models |
| finance.clearance_changed | Academic, LMS, Portal | student_clearance_snapshots, lms_clearance_snapshots, dashboard_read_models |
| academic.krs_approved | LMS, Portal | lms_enrollments, dashboard_read_models |
| assessment.result_calculated | PMB, LMS, Academic | assessment_result_snapshots atau grade_inputs |

Setiap snapshot/read model minimal memiliki source_event_key, source_event_name, source_updated_at, synced_at, dan sync_status.

## A.10 Reconciliation Rule

| Relasi Kritis | Reconciliation Rule |
| --- | --- |
| PMB invoice snapshot vs Finance invoice/payment | Cocokkan invoice_id, payment status, paid_amount, dan paid_at. |
| Academic clearance snapshot vs Finance clearance | Cocokkan subject_ref_id, service_code, academic_period_ref_id, dan status_code. |
| LMS enrollment vs Academic KRS | Cocokkan krs_item_ref_id, student_ref_id, course_offering_ref_id, dan enrollment status. |
| Academic grade input vs LMS/Assessment source | Cocokkan source_module, source_ref_id, score, weight, dan submitted_at. |
| Portal dashboard read model vs source events | Cocokkan refreshed_at, source status, dan payload aggregate. |

Jika snapshot berbeda dari source of truth, sistem membuat mismatch report. Jika dapat diperbaiki otomatis, correction job dijalankan. Jika perlu keputusan admin, status data menjadi pending_review.

## A.11 Security dan Authorization Event

| Aspek | Aturan |
| --- | --- |
| Publisher authorization | Hanya service owner domain yang boleh publish event domainnya. |
| Consumer authorization | Hanya consumer terdaftar yang boleh subscribe event tertentu. |
| PII minimization | Payload tidak boleh membawa data pribadi berlebihan. Gunakan ref_id dan snapshot minimum. |
| Service authentication | Publish dan consume event memakai service credential yang dikelola Core/Security. |
| Audit | Publish, consume, retry, DLQ, dan replay event wajib tercatat audit. |
| Replay control | Replay DLQ hanya boleh dilakukan role DevOps/SRE atau admin teknis yang diberi izin. |

## A.12 Error Contract

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

## A.13 Observability dan Monitoring

| Metric | Tujuan |
| --- | --- |
| outbox_pending_count | Mengukur event yang belum dipublish. |
| inbox_pending_count | Mengukur event masuk yang belum diproses. |
| event_lag_seconds | Mengukur keterlambatan event dari occurred_at ke processed_at. |
| retry_count_by_event | Melihat event yang sering gagal. |
| dlq_count_by_consumer | Menemukan consumer bermasalah. |
| reconciliation_mismatch_count | Mengukur selisih source of truth dan snapshot/read model. |

## A.14 UAT Event Contract

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

## A.15 Template Final Event Contract

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
