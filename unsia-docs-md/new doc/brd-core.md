# BRD Final - Core Module

## 1. Visi & Kebutuhan Bisnis
Modul Core bertindak sebagai otoritas utama otentikasi identitas (*single identity*) dan Single Sign-On (SSO) di seluruh ekosistem UNSIA ERP. Modul ini memastikan seluruh pengguna memiliki kredensial terpusat dan perizinan berbasis peran (RBAC) yang konsisten.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-CORE-B-001** | Otentikasi Terpusat (SSO) | Seluruh sistem (SIAKAD, Finance, PMB, LMS, HRIS) wajib menggunakan satu pintu gerbang otentikasi yang sama. | P0 | Tidak ada tabel credential/password mandiri di luar `core_db`. |
| **BRD-CORE-B-002** | Manajemen Hak Akses (RBAC) | Membatasi akses menu, tombol, dan data transaksi berdasarkan peran aktif pengguna beserta lingkup program studinya. | P0 | Role-based menu rendering dan backend-enforced API scope validation berjalan 100%. |
| **BRD-CORE-B-003** | Switch Active Role | Pengguna yang memiliki peran ganda (misal: Dosen sekaligus Kaprodi) dapat berpilih peran tanpa login ulang. | P0 | Waktu proses pergantian peran aktif di UI di bawah 1 detik. |
| **BRD-CORE-B-004** | Impersonation Sesi | Mengizinkan administrator BPPTI masuk sebagai pengguna lain untuk troubleshooting dengan audit yang ketat. | P1 | Setiap sesi impersonasi wajib meminta alasan tertulis dan dibatasi durasi aktifnya. |
| **BRD-CORE-B-005** | Global Audit Trail | Mencatat riwayat perubahan sensitif di seluruh sistem demi kepatuhan audit keamanan informasi. | P0 | Seluruh log mutasi data dapat ditelusuri per IP, waktu, lama/baru value, dan request ID. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-CORE-001**: Kredensial dan password hash dilarang keras diduplikasi atau disimpan dalam database modul lain di luar `core_db`.
* **BR-CORE-002**: Pengguna yang berstatus *Suspended* atau *Inactive* otomatis langsung kehilangan hak akses masuk ke seluruh portal ERP.
* **BR-CORE-003**: Impersonation hanya diperbolehkan jika disetujui melalui alasan operasional tertulis yang sah dan wajib memicu alert notifikasi ke email pengguna yang bersangkutan.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Pendaftaran pengguna baru, approval penugasan peran baru, data log masuk ip address.
* **Output**: Token JWT valid, data introspeksi otorisasi, dashboard app launcher.
* **Integrasi Lintas Domain**: Core mengirimkan pemberitahuan perubahan identitas diri lewat event broker ke seluruh modul ERP untuk penyelarasan profil lokal.
