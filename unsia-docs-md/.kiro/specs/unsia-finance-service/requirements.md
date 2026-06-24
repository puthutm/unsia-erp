# Requirements Document

## Introduction

`unsia-finance-service` adalah backend service ke-5 dalam rantai dependency ERP UNSIA yang dibangun dengan Go dan Gin framework. Service ini merupakan **otoritas tunggal (source of truth)** untuk seluruh operasi keuangan institusi: invoice dan tagihan mahasiswa/pendaftar, penerimaan pembayaran via payment gateway maupun manual, verifikasi pembayaran, status bebas tanggungan (clearance) mahasiswa, beasiswa, cicilan, kas internal, penggajian, pencairan komisi, pencatatan pajak dan BPJS, serta akuntansi (chart of accounts, jurnal, anggaran).

Service ini beroperasi pada port `:8005`, menggunakan database `finance_db` (PostgreSQL), dan berinteraksi dengan service lain **hanya** melalui event outbox/inbox dan HTTP API resmi — tidak ada akses lintas database secara langsung.

### Posisi dalam Dependency Chain

```
core-service (auth/SSO) → reference-service (master data)
→ crm-service (leads) → pmb-service (applicants)
→ **finance-service** (invoice, payment, clearance)
→ academic-service (mahasiswa, KRS, nilai)
```

### Teknologi

- **Bahasa & Framework**: Go + Gin
- **Database**: PostgreSQL (`finance_db`) dengan golang-migrate
- **Auth**: JWT RS256 validasi via JWKS dari core-service
- **Event**: Transactional Outbox/Inbox Pattern (RabbitMQ)
- **Shared packages**: `shared-auth`, `shared-rbac`, `shared-errorenvelope`, `shared-audit`, `shared-idempotency`, `shared-event`, `shared-observability`


## Glossary

- **Finance_Service**: Backend service Go/Gin yang mengelola domain keuangan ERP UNSIA, berjalan di port `:8005`.
- **Invoice**: Tagihan resmi yang diterbitkan kepada pendaftar (applicant) atau mahasiswa, berisi satu atau lebih item pembayaran.
- **Invoice_Item**: Baris detail tagihan yang mengacu pada `payment_component_id` dari reference-service.
- **Payment**: Catatan penerimaan uang yang dikaitkan dengan invoice; dapat berasal dari payment gateway atau pembayaran manual.
- **Payment_Gateway**: Penyedia layanan pembayaran pihak ketiga (Midtrans / Xendit) yang mengirimkan callback ke Finance_Service.
- **Payment_Callback**: Notifikasi asinkron dari Payment_Gateway yang diterima oleh Finance_Service untuk memperbarui status Payment.
- **Payment_Verification**: Proses verifikasi pembayaran manual oleh staf finance, menghasilkan status `VERIFIED` atau `REJECTED`.
- **Scholarship**: Beasiswa yang diberikan kepada mahasiswa, dapat berupa potongan nominal atau persentase tagihan.
- **Installment_Request**: Permintaan cicilan pembayaran yang diajukan mahasiswa dan disetujui admin finance.
- **Clearance_Policy**: Aturan yang mendefinisikan kondisi kapan seorang mahasiswa dinyatakan bebas tanggungan keuangan untuk layanan tertentu.
- **Student_Clearance**: Status bebas tanggungan keuangan seorang mahasiswa untuk suatu periode akademik dan lingkup layanan.
- **Clearance_Dispensation**: Pengecualian sementara atas status clearance yang diblokir, disetujui oleh admin finance.
- **Cash_Account**: Rekening kas/bank internal institusi yang digunakan sebagai akun penerimaan atau pengeluaran.
- **Cash_Transaction**: Mutasi debet/kredit pada Cash_Account yang dipicu oleh pembayaran, pencairan, atau transaksi lain.
- **Payroll_Run**: Proses penggajian yang mencakup satu periode tertentu dan menghasilkan sejumlah Payroll_Item.
- **Payroll_Item**: Baris detail gaji per karyawan dalam satu Payroll_Run, berisi komponen bruto, potongan, dan neto.
- **Disbursement**: Pencairan komisi kepada agen CRM atau penerima lain berdasarkan `commission_record_id`.
- **Tax_Record**: Catatan kewajiban pajak yang dikaitkan dengan transaksi tertentu.
- **BPJS_Record**: Catatan iuran BPJS ketenagakerjaan/kesehatan per karyawan per periode.
- **COA_Account**: Akun dalam Chart of Accounts (bagan akun) keuangan institusi.
- **Journal**: Header jurnal akuntansi yang mengacu pada transaksi sumber (source_type, source_id).
- **Journal_Entry**: Baris debet/kredit dalam satu Journal yang mengacu pada COA_Account.
- **Budget**: Rencana anggaran tahunan institusi dengan total alokasi dan status.
- **Budget_Line**: Baris detail anggaran per akun COA dalam satu Budget.
- **Outbox_Event**: Event yang ditulis secara atomik bersama transaksi domain ke tabel `finance.outbox_events`, siap dipublikasikan ke message broker.
- **Inbox_Event**: Event yang diterima dari modul lain dan dicatat ke tabel `finance.inbox_events` untuk diproses secara idempoten.
- **Idempotency_Key**: Header `Idempotency-Key` pada request kritis yang memastikan operasi tidak diproses lebih dari satu kali.
- **JWKS**: JSON Web Key Set yang diterbitkan core-service untuk validasi JWT RS256.
- **X-Application-Code**: Header wajib yang mengidentifikasi aplikasi pengirim request.
- **X-Active-Role**: Header wajib yang menyatakan peran aktif pengguna saat request dikirim.
- **X-Correlation-Id**: Header wajib untuk tracing request lintas service.
- **BLOCKED**: Status clearance mahasiswa yang berarti akses layanan ditolak.
- **CONDITIONAL**: Status clearance mahasiswa yang berarti akses diberikan dengan syarat tertentu.
- **CLEARED**: Status clearance mahasiswa yang berarti bebas tanggungan penuh.
- **REVOKED**: Status clearance yang sebelumnya CLEARED namun kemudian dicabut.


## Requirements

---

### Persyaratan 1: Autentikasi dan Otorisasi Request

**User Story:** Sebagai sistem keamanan ERP UNSIA, saya ingin setiap request ke Finance_Service divalidasi token dan perannya, sehingga hanya pengguna berwenang yang dapat mengakses atau memodifikasi data keuangan.

#### Kriteria Penerimaan

1. WHEN sebuah request masuk ke endpoint protected, THE Finance_Service SHALL memvalidasi token JWT RS256 dengan memverifikasi signature menggunakan public key yang diambil dari JWKS endpoint core-service.
2. WHEN sebuah request masuk ke endpoint protected tanpa header `Authorization`, THE Finance_Service SHALL mengembalikan respons HTTP 401 dengan error code `AUTH_TOKEN_MISSING`.
3. WHEN sebuah request memiliki token JWT yang sudah kedaluwarsa, THE Finance_Service SHALL mengembalikan respons HTTP 401 dengan error code `AUTH_TOKEN_EXPIRED`.
4. WHEN sebuah request memiliki token JWT dengan signature tidak valid atau format tidak sesuai, THE Finance_Service SHALL mengembalikan respons HTTP 401 dengan error code `AUTH_TOKEN_INVALID`.
5. WHEN sebuah request masuk ke endpoint protected, THE Finance_Service SHALL memeriksa keberadaan header `X-Application-Code`, `X-Active-Role`, dan `X-Correlation-Id`.
6. WHEN header `X-Application-Code`, `X-Active-Role`, atau `X-Correlation-Id` tidak ada pada request ke endpoint protected, THE Finance_Service SHALL mengembalikan respons HTTP 400 dengan pesan error yang menyebutkan nama header yang hilang secara spesifik.
7. WHEN pengguna dengan role `admin_finance` mengirim request GET, POST, PUT, atau DELETE ke endpoint manajemen invoice, payment, clearance, payroll, dan anggaran, THE Finance_Service SHALL mengizinkan akses sesuai dengan permission yang terdaftar untuk role tersebut.
8. WHEN pengguna dengan role `admin_pmb` mengirim request GET ke endpoint invoice, THE Finance_Service SHALL mengizinkan akses baca hanya untuk invoice dengan `target_type = 'applicant'`.
9. WHEN pengguna dengan role `mahasiswa` mengirim request ke endpoint invoice atau clearance, THE Finance_Service SHALL membatasi akses hanya pada data milik mahasiswa tersebut berdasarkan `student_id` yang sesuai dengan `sub` claim JWT.
10. WHEN pengguna dengan role `pimpinan` mengirim request GET ke endpoint Finance_Service, THE Finance_Service SHALL mengizinkan akses baca pada seluruh data keuangan.
11. IF pengguna tidak memiliki permission yang diperlukan untuk suatu aksi, THEN THE Finance_Service SHALL mengembalikan respons HTTP 403 dengan error code `PERMISSION_DENIED`.
12. WHEN JWKS endpoint core-service tidak dapat dijangkau saat validasi token, THE Finance_Service SHALL menggunakan JWKS yang terakhir berhasil di-cache selama TTL cache belum habis (minimum 5 menit), dan mengembalikan HTTP 503 dengan error code `AUTH_SERVICE_UNAVAILABLE` jika cache sudah expired.

---

### Persyaratan 2: Manajemen Invoice (P0)

**User Story:** Sebagai admin finance atau sistem PMB, saya ingin dapat membuat dan mengelola invoice untuk pendaftar maupun mahasiswa, sehingga tagihan keuangan tercatat secara akurat dan dapat digunakan sebagai dasar penerimaan pembayaran.

#### Kriteria Penerimaan

1. WHEN `POST /api/v1/finance/invoices` dipanggil dengan payload valid oleh service PMB atau Academic menggunakan service token, THE Finance_Service SHALL membuat invoice baru dengan status `DRAFT` dan menghasilkan `invoice_number` yang unik, lalu mengembalikan HTTP 201 dengan data invoice yang dibuat.
2. IF request `POST /api/v1/finance/invoices` tidak memiliki field wajib (`target_type`, `target_id`, `due_date`, minimal satu `items` dengan `payment_component_id` dan `final_amount` positif), THEN THE Finance_Service SHALL mengembalikan HTTP 422 dengan daftar field yang tidak valid.
3. WHEN Finance_Service memvalidasi `payment_component_id` dalam `items`, THE Finance_Service SHALL memanggil reference-service untuk memverifikasi setiap komponen aktif, dan mengembalikan HTTP 422 dengan error code `PAYMENT_COMPONENT_NOT_FOUND` jika salah satu komponen tidak valid atau tidak aktif.
4. WHEN invoice berhasil dibuat, THE Finance_Service SHALL menjumlahkan seluruh `final_amount` dari `invoice_items` dan menyimpannya sebagai `total_amount` pada tabel `invoices` secara sinkron sebelum respons dikembalikan.
5. WHEN invoice berhasil dibuat, THE Finance_Service SHALL menulis Outbox_Event `finance.invoice_created` secara atomik dalam satu transaksi database yang sama.
6. THE Finance_Service SHALL memastikan bahwa `invoice_number` bersifat unik di seluruh tabel `invoices` dan tidak dapat diduplikasi.
7. WHEN `GET /api/v1/finance/invoices/:id` dipanggil dengan `invoice_id` yang valid, THE Finance_Service SHALL mengembalikan detail invoice beserta seluruh invoice_items dalam satu respons.
8. IF `invoice_id` pada `GET /api/v1/finance/invoices/:id` tidak ditemukan di database, THEN THE Finance_Service SHALL mengembalikan HTTP 404 dengan error code `INVOICE_NOT_FOUND`.
9. WHEN `GET /api/v1/finance/invoices` dipanggil, THE Finance_Service SHALL mendukung filter berdasarkan `status`, `target_type`, `applicant_id`, `student_id`, `academic_period_id`, dan rentang tanggal `due_date_from`/`due_date_to`, serta mendukung pagination dengan `page` (default 1) dan `limit` (default 20, maksimal 100).
10. WHEN header `Idempotency-Key` dikirimkan pada `POST /api/v1/finance/invoices`, THE Finance_Service SHALL mengembalikan respons yang identik untuk request kedua dengan `Idempotency-Key` yang sama tanpa membuat invoice baru.
11. IF `due_date` pada request pembuatan invoice sudah lewat dari tanggal saat ini, THEN THE Finance_Service SHALL mengembalikan HTTP 422 dengan error code `INVALID_DUE_DATE`.
12. WHEN invoice berstatus `ISSUED` dan `paid_amount` sama dengan `total_amount` setelah pembaruan pembayaran, THE Finance_Service SHALL mengubah status invoice menjadi `PAID` secara otomatis dalam transaksi yang sama.
13. WHEN `paid_amount` diperbarui menjadi lebih besar dari nol tetapi kurang dari `total_amount`, THE Finance_Service SHALL mengubah status invoice menjadi `PARTIALLY_PAID`.
14. IF operasi pembaruan `paid_amount` akan menyebabkan `paid_amount` melebihi `total_amount`, THEN THE Finance_Service SHALL mengembalikan HTTP 422 dengan error code `OVERPAYMENT_NOT_ALLOWED` dan tidak mengubah `paid_amount`.
15. WHILE invoice berstatus `DRAFT`, THE Finance_Service SHALL mengizinkan perubahan status ke `ISSUED` hanya oleh pengguna dengan role `admin_finance`.
16. IF invoice berstatus `CANCELLED` atau `EXPIRED` dan menerima request pembuatan Payment baru, THEN THE Finance_Service SHALL mengembalikan HTTP 409 dengan error code `INVOICE_NOT_PAYABLE`.


---

### Persyaratan 3: State Machine Invoice

**User Story:** Sebagai admin finance, saya ingin status invoice mengikuti aturan transisi yang ketat, sehingga tidak ada perubahan status yang tidak valid yang dapat merusak integritas data tagihan.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL mengimplementasikan state machine invoice dengan status valid: `DRAFT`, `ISSUED`, `PARTIALLY_PAID`, `PAID`, `CANCELLED`, dan `EXPIRED`.
2. THE Finance_Service SHALL mengizinkan transisi status invoice hanya pada jalur yang valid: `DRAFT → ISSUED`, `ISSUED → PARTIALLY_PAID`, `ISSUED → PAID`, `ISSUED → CANCELLED`, `ISSUED → EXPIRED`, `PARTIALLY_PAID → PAID`, `PARTIALLY_PAID → CANCELLED`, `PARTIALLY_PAID → EXPIRED`.
3. IF sebuah operasi mencoba mengubah status invoice ke transisi yang tidak valid, THEN THE Finance_Service SHALL mengembalikan HTTP 409 dengan error code `INVALID_STATUS_TRANSITION` dan menyebutkan status asal dan status tujuan secara eksplisit.
4. WHEN invoice berstatus `PAID` atau `CANCELLED`, THE Finance_Service SHALL menolak semua permintaan transisi status lebih lanjut dengan HTTP 409.
5. THE Finance_Service SHALL memastikan bahwa setiap perubahan status invoice dicatat dengan `actor_id`, `old_status`, `new_status`, dan `changed_at` melalui modul `shared-audit`.
6. FOR ALL operasi transisi status invoice yang valid, THE Finance_Service SHALL memastikan bahwa status invoice sebelum dan sesudah operasi tetap konsisten dengan aturan state machine (properti invariansi state machine).

---

### Persyaratan 4: Penerimaan Payment Gateway Callback (P0)

**User Story:** Sebagai sistem payment gateway (Midtrans/Xendit), saya ingin dapat mengirimkan notifikasi pembayaran ke Finance_Service, sehingga status pembayaran mahasiswa dapat diperbarui secara real-time dan akurat.

#### Kriteria Penerimaan

1. WHEN `POST /api/v1/finance/payment-callbacks/:provider` diterima, THE Finance_Service SHALL menyimpan raw payload ke tabel `payment_gateway_callbacks` dengan status `received` sebelum melakukan validasi bisnis apapun.
2. IF payload callback tidak memiliki field wajib (`provider_event_id`, `order_id`, `amount`, `status`), THEN THE Finance_Service SHALL mengembalikan HTTP 400 dengan error code `CALLBACK_PAYLOAD_INVALID` dan tidak menyimpan record.
3. THE Finance_Service SHALL memvalidasi signature callback dari payment gateway menggunakan secret key yang dikonfigurasi per provider, dan menandai `signature_valid = false` serta status `ignored` jika signature tidak cocok.
4. IF secret key untuk `:provider` tidak dikonfigurasi di environment, THEN THE Finance_Service SHALL mencatat log ERROR dan mengembalikan HTTP 500 dengan error code `PROVIDER_NOT_CONFIGURED`.
5. WHEN callback diterima dengan kombinasi `(provider, provider_event_id)` yang sudah ada di database, THE Finance_Service SHALL mengembalikan HTTP 200 dengan status `ignored` tanpa memproses ulang.
6. WHEN callback diterima dengan signature valid dan `provider_event_id` belum pernah diproses, THE Finance_Service SHALL membuat atau memperbarui record Payment dan memperbarui `paid_amount` pada Invoice secara atomik dalam satu transaksi database.
7. IF transaksi atomik pada kriteria 6 gagal, THEN THE Finance_Service SHALL melakukan rollback seluruh perubahan Payment dan Invoice, namun tetap mempertahankan record `payment_gateway_callbacks` dengan status `error`.
8. WHEN payment berhasil diproses dari callback, THE Finance_Service SHALL menulis Outbox_Event `finance.payment_paid` secara atomik dalam transaksi yang sama.
9. IF `amount` pada callback berbeda dari `total_amount` invoice (toleransi nol), THEN THE Finance_Service SHALL menandai payment dengan status `RECEIVED` dan menambahkan record ke antrian verifikasi manual, bukan langsung `VERIFIED`.
10. WHEN header `Idempotency-Key` dikirim pada payment callback request dengan nilai antara 1–255 karakter, THE Finance_Service SHALL menggunakan nilai tersebut sebagai `idempotency_key` pada tabel `payment_gateway_callbacks` dan `payments`.
11. WHEN callback diterima tanpa error internal, THE Finance_Service SHALL mengembalikan HTTP 200 agar payment gateway tidak melakukan retry yang tidak perlu.
12. WHEN terjadi error internal (HTTP 500) saat memproses callback, THE Finance_Service SHALL mengembalikan HTTP 500 sehingga payment gateway dapat melakukan retry.

---

### Persyaratan 5: Verifikasi Pembayaran Manual (P0)

**User Story:** Sebagai admin finance, saya ingin dapat memverifikasi bukti pembayaran yang dikirimkan secara manual oleh mahasiswa, sehingga pembayaran non-gateway dapat diakui dan invoice dapat diperbarui statusnya.

#### Kriteria Penerimaan

1. WHEN `POST /api/v1/finance/payment-verifications` dipanggil dengan `payment_id` valid dan `verification_status = 'approved'`, THE Finance_Service SHALL mengubah status payment menjadi `VERIFIED` dan memperbarui `paid_amount` pada invoice terkait.
2. WHEN verifikasi pembayaran manual berhasil, THE Finance_Service SHALL menulis Outbox_Event `finance.payment_paid` secara atomik dalam satu transaksi.
3. IF verifikasi pembayaran ditolak (`verification_status = 'rejected'`), THEN THE Finance_Service SHALL mewajibkan field `reason` tidak kosong dan mengubah status payment menjadi `FAILED`.
4. THE Finance_Service SHALL memastikan bahwa satu Payment hanya dapat memiliki satu PaymentVerification aktif pada satu waktu.
5. WHEN verifikasi pembayaran manual dilakukan, THE Finance_Service SHALL mencatat `verified_by`, `verification_status`, `note`, dan `verified_at` ke tabel `payment_verifications`.
6. THE Finance_Service SHALL memastikan operasi verifikasi pembayaran bersifat idempoten menggunakan `Idempotency-Key`.
7. IF payment yang akan diverifikasi sudah berstatus `VERIFIED` atau `POSTED`, THEN THE Finance_Service SHALL mengembalikan HTTP 409 dengan error code `PAYMENT_ALREADY_VERIFIED`.


---

### Persyaratan 6: State Machine Payment

**User Story:** Sebagai admin finance, saya ingin status pembayaran mengikuti aturan transisi yang ketat dan terdefinisi, sehingga alur pembayaran dari penerimaan hingga pencatatan buku besar berjalan dengan benar.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL mengimplementasikan state machine payment dengan status valid: `RECEIVED`, `VERIFIED`, `POSTED`, `FAILED`, dan `REVERSED`.
2. THE Finance_Service SHALL mengizinkan transisi status payment hanya pada jalur valid: `RECEIVED → VERIFIED`, `RECEIVED → FAILED`, `VERIFIED → POSTED`, `VERIFIED → REVERSED`, `POSTED → REVERSED`.
3. IF sebuah operasi mencoba mengubah status payment ke transisi yang tidak valid, THEN THE Finance_Service SHALL mengembalikan HTTP 409 dengan error code `INVALID_PAYMENT_STATUS_TRANSITION`.
4. WHEN payment berstatus `POSTED`, THE Finance_Service SHALL menolak semua perubahan status kecuali `REVERSED`, dan mewajibkan `reason` untuk transaksi REVERSED.
5. THE Finance_Service SHALL memastikan bahwa setiap perubahan status payment dicatat melalui `shared-audit` dengan informasi `actor_id`, `old_status`, `new_status`, dan `changed_at`.

---

### Persyaratan 7: Manajemen Clearance Mahasiswa (P0)

**User Story:** Sebagai admin finance atau sistem akademik, saya ingin dapat memeriksa dan mengelola status bebas tanggungan keuangan mahasiswa, sehingga akses layanan akademik (KRS, KHS, wisuda) hanya diberikan kepada mahasiswa yang memenuhi kewajiban keuangannya.

#### Kriteria Penerimaan

1. WHEN `GET /api/v1/finance/clearances` dipanggil dengan parameter `student_id` dan `service_scope`, THE Finance_Service SHALL mengembalikan status clearance terkini untuk mahasiswa tersebut pada periode akademik yang aktif, atau HTTP 404 jika belum ada record clearance.
2. WHEN `POST /api/v1/finance/clearances` dipanggil oleh admin_finance dengan field wajib (`student_id`, `academic_period_id`, `service_scope`, `status`) yang valid, THE Finance_Service SHALL membuat record baru dengan HTTP 201 jika belum ada, atau memperbarui record yang sudah ada dengan HTTP 200.
3. WHEN `POST /api/v1/finance/clearances` dipanggil, THE Finance_Service SHALL mengevaluasi `clearance_policies` yang aktif dan relevan dengan `service_scope` yang diminta untuk menentukan status awal `BLOCKED`, `CONDITIONAL`, atau `CLEARED`. IF tidak ada policy yang cocok, THE Finance_Service SHALL menetapkan status default `BLOCKED`.
4. WHEN status clearance berubah dari nilai sebelumnya, THE Finance_Service SHALL menulis Outbox_Event `finance.clearance_changed` secara atomik dalam satu transaksi database yang sama.
5. THE Finance_Service SHALL mengizinkan transisi status clearance hanya pada jalur valid: `BLOCKED → CONDITIONAL`, `BLOCKED → CLEARED`, `CONDITIONAL → CLEARED`, `CONDITIONAL → BLOCKED`, `CLEARED → REVOKED`.
6. IF status clearance berubah menjadi `REVOKED`, THEN THE Finance_Service SHALL mewajibkan field `reason` tidak kosong dan mengembalikan HTTP 422 dengan error code `REASON_REQUIRED` jika `reason` kosong.
7. THE Finance_Service SHALL memastikan setiap perubahan status clearance dicatat melalui `shared-audit` dengan `updated_by`, `old_status`, `new_status`, dan `updated_at`.
8. WHEN `GET /api/v1/finance/clearances` dipanggil tanpa `student_id`, THE Finance_Service SHALL mendukung filter berdasarkan `academic_period_id`, `service_scope`, dan `status`, dengan pagination `page` (default 1) dan `limit` (default 20, maksimal 100).
9. WHEN mahasiswa memiliki Clearance_Dispensation yang aktif (`expires_at > NOW()`) dan belum kedaluwarsa, THE Finance_Service SHALL mengembalikan status `CONDITIONAL` untuk mahasiswa tersebut meskipun terdapat tunggakan.
10. THE Finance_Service SHALL memastikan bahwa Outbox_Event `finance.clearance_changed` hanya ditulis ketika status clearance benar-benar berubah dari nilai sebelumnya.
11. IF `student_id` pada request tidak ditemukan di cache atau referensi mahasiswa yang valid, THEN THE Finance_Service SHALL mengembalikan HTTP 404 dengan error code `STUDENT_NOT_FOUND`. IF tidak ada periode akademik aktif yang ditemukan, THE Finance_Service SHALL mengembalikan HTTP 422 dengan error code `NO_ACTIVE_ACADEMIC_PERIOD`.

---

### Persyaratan 8: State Machine Clearance

**User Story:** Sebagai admin finance, saya ingin aturan transisi status clearance mahasiswa diterapkan secara konsisten, sehingga tidak ada perubahan clearance yang tidak sah yang dapat mengakibatkan mahasiswa mendapat akses yang tidak semestinya.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL mengimplementasikan state machine clearance dengan status valid: `BLOCKED`, `CONDITIONAL`, `CLEARED`, dan `REVOKED`.
2. IF sebuah operasi mencoba mengubah status clearance ke transisi yang tidak valid, THEN THE Finance_Service SHALL mengembalikan HTTP 409 dengan error code `INVALID_CLEARANCE_STATUS_TRANSITION`.
3. FOR ALL operasi perubahan status clearance yang valid, THE Finance_Service SHALL memastikan bahwa event `finance.clearance_changed` yang dipublikasikan berisi `old_status`, `new_status`, `student_id`, `service_scope`, dan `academic_period_id` yang akurat.


---

### Persyaratan 9: Beasiswa (P1)

**User Story:** Sebagai admin finance, saya ingin dapat mengelola data beasiswa mahasiswa, sehingga potongan tagihan yang bersumber dari beasiswa dapat diterapkan secara akurat pada invoice yang relevan.

#### Kriteria Penerimaan

1. WHEN `POST /api/v1/finance/scholarships` dipanggil dengan `student_id`, `scholarship_type`, `amount`, dan data valid, THE Finance_Service SHALL membuat record beasiswa baru dengan status `PENDING_APPROVAL`.
2. WHEN beasiswa disetujui oleh admin_finance, THE Finance_Service SHALL mengubah status scholarship menjadi `APPROVED` dan mencatat `approved_by` serta `approved_at`.
3. WHEN `GET /api/v1/finance/scholarships` dipanggil, THE Finance_Service SHALL mendukung filter berdasarkan `student_id`, `status`, dan `scholarship_type`, serta mendukung pagination.
4. IF `amount` beasiswa melebihi `total_amount` invoice yang terkait, THEN THE Finance_Service SHALL membatasi penerapan potongan maksimal sebesar `total_amount` invoice tersebut dan tidak mengizinkan nilai negatif pada tagihan.
5. THE Finance_Service SHALL memastikan setiap operasi create dan approve beasiswa dicatat melalui `shared-audit`.

---

### Persyaratan 10: Cicilan Pembayaran (P1)

**User Story:** Sebagai mahasiswa atau admin finance, saya ingin dapat mengajukan permohonan cicilan untuk tagihan yang belum lunas, sehingga mahasiswa dengan keterbatasan finansial dapat tetap memenuhi kewajibannya secara bertahap.

#### Kriteria Penerimaan

1. WHEN `POST /api/v1/finance/invoices/:id/request-installment` dipanggil dengan `student_id` dan `reason` valid, THE Finance_Service SHALL membuat Installment_Request baru dengan status `PENDING`.
2. THE Finance_Service SHALL memvalidasi bahwa invoice yang dirujuk berstatus `ISSUED` atau `PARTIALLY_PAID` sebelum membuat Installment_Request, dan mengembalikan HTTP 422 jika invoice sudah `PAID` atau `CANCELLED`.
3. WHEN Installment_Request disetujui oleh admin_finance, THE Finance_Service SHALL mengubah status clearance mahasiswa terkait menjadi `CONDITIONAL` dan mencatat `approved_by` serta `approved_at`.
4. THE Finance_Service SHALL memastikan bahwa persetujuan Installment_Request bersifat idempoten menggunakan `Idempotency-Key`.

---

### Persyaratan 11: Integrasi Event Inbound dari PMB dan Academic

**User Story:** Sebagai Finance_Service, saya ingin menerima dan memproses event dari PMB dan Academic secara idempoten, sehingga data invoice dan referensi mahasiswa selalu sinkron dengan modul upstream.

#### Kriteria Penerimaan

1. WHEN Finance_Service menerima event `pmb.applicant_created` dari PMB, THE Finance_Service SHALL mencatat event tersebut di tabel `finance.inbox_events` dengan status `received`, kemudian memeriksa `event_key` untuk memastikan event belum pernah diproses. IF `event_key` sudah ada dengan status `processed`, THE Finance_Service SHALL menandai event baru sebagai `duplicate` tanpa efek samping apapun.
2. WHEN Finance_Service memproses event `pmb.applicant_created` dengan flag `auto_create_invoice = true` dalam payload, THE Finance_Service SHALL membuat invoice baru menggunakan daftar `payment_component_id` yang tercantum dalam payload event. IF salah satu `payment_component_id` tidak valid, THE Finance_Service SHALL menandai event sebagai `failed` dan mencatat error tanpa membuat invoice parsial.
3. WHEN Finance_Service menerima event `academic.student_created`, THE Finance_Service SHALL memperbarui `student_id` pada invoice yang memiliki `applicant_id` yang sesuai. IF tidak ada invoice dengan `applicant_id` tersebut, THE Finance_Service SHALL menandai event sebagai `processed` tanpa error karena invoice mungkin belum dibuat.
4. THE Finance_Service SHALL memastikan bahwa event inbound yang sama berdasarkan `event_key` ditandai `duplicate` dan dilewati tanpa memproses ulang atau menghasilkan efek samping.
5. IF pemrosesan event inbound gagal setelah jumlah retry yang dikonfigurasi (minimum 1, maksimum 10), THEN THE Finance_Service SHALL memindahkan event ke status `dead_letter` pada tabel `finance.inbox_events` untuk ditangani secara manual, dengan mencatat detail kegagalan terakhir.

### Persyaratan 12: Publikasi Event Outbound

**User Story:** Sebagai modul downstream (PMB, Academic, Portal), saya ingin menerima event dari Finance_Service ketika terjadi perubahan status yang relevan, sehingga data di modul saya selalu terkini tanpa harus polling.

#### Kriteria Penerimaan

1. WHEN invoice berhasil dibuat, THE Finance_Service SHALL menulis Outbox_Event dengan `event_name = 'finance.invoice_created'` ke tabel `finance.outbox_events` secara atomik dalam transaksi yang sama dengan pembuatan invoice.
2. WHEN payment dikonfirmasi dari verifikasi manual (status berubah menjadi `VERIFIED`), THE Finance_Service SHALL menulis Outbox_Event dengan `event_name = 'finance.payment_paid'` secara atomik dalam transaksi verifikasi yang sama.
3. WHEN payment berhasil diproses dari gateway callback valid, THE Finance_Service SHALL menulis Outbox_Event dengan `event_name = 'finance.payment_paid'` secara atomik dalam transaksi callback yang sama.
4. WHEN status clearance mahasiswa berubah dari nilai sebelumnya, THE Finance_Service SHALL menulis Outbox_Event dengan `event_name = 'finance.clearance_changed'` secara atomik.
5. THE Finance_Service SHALL memastikan setiap Outbox_Event berisi field standar: `event_name`, `event_version`, `event_key`, `publisher_service`, `aggregate_type`, `aggregate_id`, `correlation_id`, `causation_id`, `occurred_at`, dan `payload`.
6. THE Finance_Service SHALL memastikan `event_key` bersifat deterministik mengikuti format `{event_name}:{aggregate_id}:{occurred_at_unix}` sehingga event yang sama tidak pernah menghasilkan `event_key` duplikat di tabel `outbox_events`.


---

### Persyaratan 13: Kas dan Transaksi Kas (P2)

**User Story:** Sebagai admin finance, saya ingin mengelola rekening kas/bank internal institusi dan mencatat setiap mutasinya, sehingga posisi kas institusi selalu terpantau secara akurat.

#### Kriteria Penerimaan

1. WHEN `GET /api/v1/finance/cash-accounts` dipanggil oleh admin_finance, THE Finance_Service SHALL mengembalikan daftar semua Cash_Account dengan `account_code`, `account_name`, `bank_name`, dan `is_active`.
2. WHEN sebuah Cash_Transaction dibuat untuk Cash_Account, THE Finance_Service SHALL memvalidasi bahwa `cash_account_id` merujuk pada akun yang aktif (`is_active = true`), dan mengembalikan HTTP 422 jika akun tidak aktif.
3. THE Finance_Service SHALL memastikan bahwa setiap Cash_Transaction memiliki `transaction_type` (debet/kredit), `amount` positif, dan `transaction_at` yang valid.
4. THE Finance_Service SHALL memastikan bahwa total saldo Cash_Account dapat dihitung dengan benar dari jumlah debet dikurangi jumlah kredit seluruh transaksi yang terkait.

---

### Persyaratan 14: Penggajian / Payroll (P2)

**User Story:** Sebagai admin finance, saya ingin dapat membuat dan mengelola siklus penggajian karyawan, sehingga proses pembayaran gaji berjalan secara terstruktur dan dapat diaudit.

#### Kriteria Penerimaan

1. WHEN `POST /api/v1/finance/payroll-runs` dipanggil dengan `payroll_period`, `run_date`, dan data Payroll_Item valid, THE Finance_Service SHALL membuat Payroll_Run baru dengan status `DRAFT`.
2. THE Finance_Service SHALL memvalidasi bahwa `net_amount` pada setiap Payroll_Item sama dengan `gross_amount` dikurangi `deduction_amount`, dan mengembalikan HTTP 422 jika tidak konsisten.
3. THE Finance_Service SHALL memastikan bahwa `total_amount` pada Payroll_Run selalu sama dengan jumlah `net_amount` dari seluruh Payroll_Item yang terkait.
4. WHEN Payroll_Run disetujui oleh admin_finance, THE Finance_Service SHALL mengubah status Payroll_Run menjadi `APPROVED` dan mengubah status seluruh Payroll_Item menjadi `APPROVED`.
5. THE Finance_Service SHALL memastikan bahwa `total_amount` Payroll_Run bersifat read-only setelah status `APPROVED` dan tidak dapat diubah.

---

### Persyaratan 15: Chart of Accounts dan Jurnal Akuntansi (P2)

**User Story:** Sebagai admin finance, saya ingin mencatat setiap transaksi keuangan ke dalam jurnal akuntansi dengan double-entry bookkeeping, sehingga laporan keuangan institusi dapat disusun dengan akurat.

#### Kriteria Penerimaan

1. WHEN `GET /api/v1/finance/coa-accounts` dipanggil oleh admin_finance, THE Finance_Service SHALL mengembalikan daftar COA_Account dengan `account_code`, `account_name`, `normal_balance`, dan `is_active`, mendukung filter `is_active` dan pagination.
2. WHEN `POST /api/v1/finance/journals` dipanggil dengan data Journal dan Journal_Entries valid, THE Finance_Service SHALL membuat Journal baru beserta seluruh Journal_Entry-nya dalam satu transaksi atomik.
3. THE Finance_Service SHALL memvalidasi prinsip double-entry bookkeeping: jumlah total debet seluruh Journal_Entry dalam satu Journal harus sama dengan jumlah total kredit, dan mengembalikan HTTP 422 jika tidak seimbang.
4. THE Finance_Service SHALL memvalidasi bahwa setiap `coa_account_id` dalam Journal_Entry merujuk pada COA_Account yang aktif.
5. THE Finance_Service SHALL memastikan bahwa `journal_number` bersifat unik di seluruh tabel `journals`.
6. FOR ALL Journal yang valid, THE Finance_Service SHALL memastikan properti: `SUM(debit of entries) = SUM(credit of entries)` (properti keseimbangan jurnal).

---

### Persyaratan 16: Anggaran dan Realisasi (P2)

**User Story:** Sebagai admin finance atau pimpinan, saya ingin dapat mengelola anggaran tahunan dan memantau realisasinya per akun COA, sehingga penggunaan anggaran dapat dikendalikan dan dilaporkan.

#### Kriteria Penerimaan

1. WHEN `GET /api/v1/finance/budgets` dipanggil, THE Finance_Service SHALL mengembalikan daftar Budget dengan `budget_code`, `name`, `fiscal_year`, `total_amount`, dan `status`, mendukung filter `fiscal_year` dan `status`.
2. WHEN Budget dibuat, THE Finance_Service SHALL memvalidasi bahwa `budget_code` bersifat unik dan `total_amount` lebih besar dari nol.
3. THE Finance_Service SHALL memastikan bahwa `total_amount` pada Budget selalu sama dengan jumlah `amount` dari seluruh Budget_Line yang terkait.
4. THE Finance_Service SHALL memastikan bahwa `realized_amount` pada setiap Budget_Line tidak melebihi `amount` Budget_Line tersebut pada kondisi normal, kecuali terdapat override eksplisit dari admin_finance.


---

### Persyaratan 17: Pencairan Komisi (Disbursement)

**User Story:** Sebagai admin finance, saya ingin dapat mengelola pencairan komisi kepada agen CRM berdasarkan data komisi yang dikirim dari CRM service, sehingga pembayaran komisi dapat dilakukan secara terstruktur dan terekam dengan baik.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL membuat record Disbursement berdasarkan `commission_record_id` dari CRM dengan status `PENDING`.
2. WHEN Disbursement disetujui dan diproses, THE Finance_Service SHALL mengubah status menjadi `DISBURSED` dan mencatat `disbursed_at`.
3. THE Finance_Service SHALL memastikan bahwa satu `commission_record_id` hanya menghasilkan satu Disbursement yang aktif (tidak ada duplikasi pencairan untuk komisi yang sama).
4. THE Finance_Service SHALL memastikan `amount` pada Disbursement selalu berupa nilai positif.

---

### Persyaratan 18: Pencatatan Pajak dan BPJS

**User Story:** Sebagai admin finance, saya ingin dapat mencatat kewajiban pajak dan iuran BPJS yang timbul dari transaksi keuangan dan penggajian, sehingga kepatuhan kewajiban fiskal institusi dapat dipantau dan dilaporkan.

#### Kriteria Penerimaan

1. WHEN Tax_Record dibuat, THE Finance_Service SHALL memvalidasi bahwa `tax_type`, `source_type`, `source_id`, `amount`, dan `tax_period` semuanya diisi.
2. WHEN BPJS_Record dibuat, THE Finance_Service SHALL memvalidasi bahwa `employee_id`, `amount`, `period`, dan `status` semuanya diisi.
3. THE Finance_Service SHALL memastikan `amount` pada Tax_Record dan BPJS_Record selalu berupa nilai positif.

---

### Persyaratan 19: Format Response API

**User Story:** Sebagai developer yang mengkonsumsi Finance_Service API, saya ingin semua respons menggunakan format envelope standar, sehingga integrasi dengan service lain dan tampilan frontend menjadi konsisten dan mudah diimplementasikan.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL membungkus semua respons sukses dalam format `SuccessEnvelope` yang berisi field: `success: true`, `code`, `message`, `data`, dan `meta` (dengan `request_id`, `correlation_id`, `timestamp`).
2. THE Finance_Service SHALL membungkus semua respons error dalam format `ErrorEnvelope` yang berisi field: `success: false`, `code`, `message`, `errors` (array `ErrorDetail` dengan `field` dan `message`), dan `meta`.
3. THE Finance_Service SHALL menggunakan kode HTTP yang tepat: 200 untuk GET sukses, 201 untuk POST yang membuat resource baru, 400 untuk request tidak valid, 401 untuk tidak terautentikasi, 403 untuk tidak terotorisasi, 404 untuk resource tidak ditemukan, 409 untuk konflik bisnis, 422 untuk validasi gagal, dan 500 untuk error internal.
4. THE Finance_Service SHALL menyertakan `X-Correlation-Id` dari request dalam field `correlation_id` pada response `meta`.
5. THE Finance_Service SHALL mengembalikan `trace_id` pada setiap respons error untuk keperluan debugging, tanpa mengekspos detail teknis internal kepada pengguna non-admin.

---

### Persyaratan 20: Idempotency pada Operasi Kritis

**User Story:** Sebagai sistem yang memanggil Finance_Service, saya ingin operasi kritis seperti pembuatan invoice, pemrosesan payment, dan perubahan clearance bersifat idempoten, sehingga request yang dikirim ulang karena timeout atau error jaringan tidak menghasilkan data duplikat.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL menerima header `Idempotency-Key` pada semua endpoint POST/PUT yang bersifat command (membuat atau mengubah data).
2. WHEN request dengan `Idempotency-Key` yang sudah ada sebelumnya diterima dan request sebelumnya telah berhasil diproses, THE Finance_Service SHALL mengembalikan respons yang identik dengan respons asli tanpa memproses ulang.
3. WHEN request dengan `Idempotency-Key` yang sama diterima sementara request pertama masih diproses (status `processing`), THE Finance_Service SHALL mengembalikan HTTP 409 dengan error code `IDEMPOTENCY_KEY_IN_PROGRESS`.
4. THE Finance_Service SHALL menyimpan hasil respons per `Idempotency-Key` ke tabel `finance.idempotency_keys` dengan `response_json` dan `status`.
5. THE Finance_Service SHALL menggunakan `shared-idempotency` package untuk implementasi mekanisme ini secara konsisten di seluruh handler.
6. FOR ALL operasi yang menggunakan `Idempotency-Key`, THE Finance_Service SHALL memastikan bahwa mengirim ulang request yang sama menghasilkan output yang identik dengan request pertama (properti idempotency).


---

### Persyaratan 21: Observabilitas dan Audit

**User Story:** Sebagai admin teknis atau auditor, saya ingin setiap operasi penting pada Finance_Service tercatat dalam audit log dan dapat di-trace, sehingga investigasi insiden dan audit kepatuhan dapat dilakukan dengan mudah.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL mencatat audit log untuk setiap operasi create, update status, approve, reject, dan cancel pada entitas: Invoice, Payment, Payment_Verification, Student_Clearance, Scholarship, Installment_Request, Payroll_Run, dan Disbursement.
2. THE Finance_Service SHALL menggunakan `shared-audit` package dan mencatat: `user_id` (actor), `active_role_id`, `module`, `action`, `entity_name`, `entity_id`, `old_value`, `new_value`, `request_id`, `ip_address`, dan `created_at`.
3. THE Finance_Service SHALL menggunakan `shared-observability` package untuk structured logging dengan level INFO, WARN, dan ERROR, serta menyertakan `trace_id` dan `correlation_id` di setiap log entry.
4. THE Finance_Service SHALL mengekspos metrics Prometheus standar: request count per endpoint, response time histogram, error rate, dan jumlah Outbox_Event pending.
5. WHEN terjadi error internal (HTTP 500), THE Finance_Service SHALL mencatat full stack trace pada log level ERROR dan tidak mengekspos detail tersebut ke response body.

---

### Persyaratan 22: Infrastruktur Service dan Konfigurasi

**User Story:** Sebagai developer yang mendeploy Finance_Service, saya ingin service dapat dikonfigurasi sepenuhnya via environment variables dan siap dijalankan di dalam container Docker, sehingga deployment ke berbagai environment (lokal, staging, produksi) menjadi mudah dan konsisten.

#### Kriteria Penerimaan

1. THE Finance_Service SHALL membaca semua konfigurasi dari environment variables, termasuk: `PORT`, `DATABASE_URL`, `JWKS_URL`, `RABBITMQ_URL`, dan konfigurasi secret per payment provider.
2. THE Finance_Service SHALL menjalankan database migration secara otomatis pada saat startup jika variabel `RUN_MIGRATION=true` diset.
3. THE Finance_Service SHALL mengekspos endpoint health check `GET /health` yang mengembalikan status koneksi database dan message broker tanpa memerlukan autentikasi.
4. THE Finance_Service SHALL siap dijalankan menggunakan Docker dengan file `Dockerfile` multi-stage build (builder + runner) untuk menghasilkan image yang ringan.
5. THE Finance_Service SHALL menyediakan file `.env.example` yang mendokumentasikan seluruh environment variable yang diperlukan dengan nilai contoh yang aman (tanpa secret sesungguhnya).
6. THE Finance_Service SHALL memvalidasi keberadaan semua environment variable wajib pada saat startup dan menghentikan proses dengan pesan error yang deskriptif jika ada yang hilang.

---

### Persyaratan 23: Penanganan Error dan Degraded Mode

**User Story:** Sebagai sistem yang bergantung pada Finance_Service, saya ingin service tetap beroperasi dengan aman ketika dependency eksternal tidak tersedia, sehingga kegagalan partial tidak menyebabkan kegagalan total pada alur keuangan.

#### Kriteria Penerimaan

1. WHEN JWKS endpoint dari core-service tidak dapat dijangkau, THE Finance_Service SHALL menggunakan JWKS yang terakhir berhasil di-cache, dengan TTL cache minimal 5 menit.
2. WHEN koneksi ke database tidak tersedia, THE Finance_Service SHALL mengembalikan HTTP 503 dengan error code `SERVICE_UNAVAILABLE` dan tidak memproses request apapun.
3. WHEN koneksi ke message broker (RabbitMQ) tidak tersedia, THE Finance_Service SHALL tetap dapat memproses operasi domain dan menyimpan event ke tabel `outbox_events`, dan mencoba mempublikasikan kembali event tersebut saat koneksi pulih.
4. WHEN pemanggilan HTTP ke reference-service untuk validasi `payment_component_id` gagal, THE Finance_Service SHALL mengembalikan HTTP 503 dengan error code `UPSTREAM_SERVICE_UNAVAILABLE` dan tidak melanjutkan pembuatan invoice.
5. THE Finance_Service SHALL mengimplementasikan retry dengan exponential backoff untuk pemanggilan HTTP ke service lain, dengan maksimal 3 kali percobaan.


---

### Persyaratan 24: Properti Correctness untuk Property-Based Testing

**User Story:** Sebagai QA Engineer, saya ingin ada properti-properti kebenaran yang dapat diuji menggunakan property-based testing (PBT), sehingga logika bisnis Finance_Service dapat diverifikasi secara komprehensif terhadap berbagai kombinasi input.

#### Kriteria Penerimaan — Invariansi Invoice

1. WHEN invoice berhasil dibuat atau diperbarui, THE Finance_Service SHALL memastikan bahwa `total_amount = SUM(final_amount untuk setiap invoice_item)` secara sinkron sebelum respons dikembalikan kepada pemanggil (properti invariansi agregasi).
2. WHEN `paid_amount` diperbarui pada invoice manapun, THE Finance_Service SHALL memastikan bahwa `paid_amount >= 0` dan `paid_amount <= total_amount` berlaku secara sinkron sebelum respons dikembalikan.
3. WHEN status invoice diperbarui, THE Finance_Service SHALL memastikan konsistensi: jika `paid_amount = 0` maka status tidak boleh `PAID` atau `PARTIALLY_PAID`; jika `0 < paid_amount < total_amount` maka status harus `PARTIALLY_PAID`; jika `paid_amount = total_amount` maka status harus `PAID`.

#### Kriteria Penerimaan — Idempotency Callback

4. WHEN payment gateway callback dengan `(provider, provider_event_id)` yang sudah pernah diproses diterima kembali, THE Finance_Service SHALL mengembalikan HTTP 200 dengan status `ignored` dan memastikan tidak ada Payment record baru yang dibuat (properti idempotency callback).
5. WHEN `POST /api/v1/finance/invoices` dipanggil berulang kali dengan `Idempotency-Key` yang sama dan request pertama telah berhasil, THE Finance_Service SHALL mengembalikan HTTP 200 dengan data invoice yang sama dan memastikan jumlah invoice yang terbuat tetap satu (properti idempotency creation).

#### Kriteria Penerimaan — Keseimbangan Jurnal

6. WHEN `POST /api/v1/finance/journals` dipanggil dengan Journal_Entries yang `SUM(debit) = SUM(credit)` dan nilai masing-masing entry antara 0.01 dan 999,999,999.99, THE Finance_Service SHALL menyimpan Journal dan mengembalikan HTTP 201 (properti penerimaan jurnal seimbang).
7. WHEN `POST /api/v1/finance/journals` dipanggil dengan Journal_Entries yang `SUM(debit) ≠ SUM(credit)`, THE Finance_Service SHALL menolak request dengan HTTP 422 dan error code `JOURNAL_NOT_BALANCED` tanpa menyimpan data apapun (properti penolakan jurnal tidak seimbang).

#### Kriteria Penerimaan — State Machine

8. WHEN serangkaian transisi status invoice valid diterapkan secara berurutan, THE Finance_Service SHALL memastikan bahwa status akhir invoice sesuai dengan transisi terakhir yang valid sesuai dengan definisi state machine di Persyaratan 3 (properti komposisi state machine).
9. WHEN operasi transisi status invoice yang tidak terdaftar dalam jalur valid pada Persyaratan 3 dicoba, THE Finance_Service SHALL mengembalikan HTTP 409 dengan error code `INVALID_STATUS_TRANSITION` dan memastikan status invoice tidak berubah.

#### Kriteria Penerimaan — Clearance Consistency

10. WHEN semua invoice untuk `student_id` dan `academic_period_id` tertentu berstatus `PAID` dan setidaknya terdapat satu invoice, THE Finance_Service SHALL memastikan bahwa evaluasi clearance menghasilkan status `CLEARED` untuk `service_scope` yang relevan jika tidak ada policy pemblokiran lain (properti konsistensi clearance-payment).
11. WHEN status clearance mahasiswa berubah secara valid, THE Finance_Service SHALL memastikan bahwa tepat satu Outbox_Event `finance.clearance_changed` ditulis dalam transaksi yang sama, berisi minimal field: `student_id`, `academic_period_id`, `service_scope`, `previous_status`, `new_status`, dan `changed_at` (properti exactly-once event emission).

#### Kriteria Penerimaan — Payroll Integrity

12. WHEN Payroll_Item dibuat atau divalidasi dengan `gross_amount >= 0.00` dan `0.00 <= deduction_amount <= gross_amount`, THE Finance_Service SHALL memastikan bahwa `net_amount = gross_amount - deduction_amount` secara sinkron sebelum respons dikembalikan (properti invariansi kalkulasi neto).
13. WHEN Payroll_Run berhasil disimpan, THE Finance_Service SHALL memastikan bahwa `total_amount = SUM(net_amount dari seluruh Payroll_Item yang terkait)` secara sinkron sebelum respons dikembalikan (properti agregasi payroll).


---

## Ringkasan Prioritas Endpoint

| Prioritas | Endpoint | Persyaratan |
|-----------|----------|-------------|
| **P0** | `POST /api/v1/finance/invoices` | Req. 2, 3 |
| **P0** | `GET /api/v1/finance/invoices/:id` | Req. 2 |
| **P0** | `GET /api/v1/finance/invoices` | Req. 2 |
| **P0** | `POST /api/v1/finance/payment-callbacks/:provider` | Req. 4 |
| **P0** | `POST /api/v1/finance/payment-verifications` | Req. 5 |
| **P0** | `GET /api/v1/finance/clearances` | Req. 7 |
| **P0** | `POST /api/v1/finance/clearances` | Req. 7, 8 |
| **P1** | `POST /api/v1/finance/invoices/:id/request-installment` | Req. 10 |
| **P1** | `GET /api/v1/finance/scholarships` | Req. 9 |
| **P1** | `POST /api/v1/finance/scholarships` | Req. 9 |
| **P2** | `GET /api/v1/finance/cash-accounts` | Req. 13 |
| **P2** | `POST /api/v1/finance/payroll-runs` | Req. 14 |
| **P2** | `GET /api/v1/finance/coa-accounts` | Req. 15 |
| **P2** | `POST /api/v1/finance/journals` | Req. 15 |
| **P2** | `GET /api/v1/finance/budgets` | Req. 16 |

## Ringkasan Event Contract

| Event | Trigger | Consumers |
|-------|---------|-----------|
| `finance.invoice_created` (outbound) | Invoice berhasil dibuat | PMB, Portal |
| `finance.payment_paid` (outbound) | Payment diverifikasi/dikonfirmasi | PMB, Academic, Portal |
| `finance.clearance_changed` (outbound) | Status clearance berubah | PMB, Academic, LMS, Portal |
| `pmb.applicant_created` (inbound) | Applicant baru di PMB | Finance (buat invoice otomatis) |
| `academic.student_created` (inbound) | Mahasiswa baru di Academic | Finance (update student_id di invoice) |

## Dependensi Eksternal

| Dependensi | Jenis | Keperluan |
|-----------|-------|-----------|
| `core-service` JWKS endpoint | HTTP (GET) | Validasi JWT RS256 |
| `reference-service` payment components | HTTP (GET) | Validasi `payment_component_id` saat buat invoice |
| RabbitMQ | Message broker | Publish/consume events via outbox/inbox |
| PostgreSQL `finance_db` | Database | Persistensi seluruh domain finance |
| Midtrans / Xendit | Payment Gateway | Menerima callback pembayaran |

