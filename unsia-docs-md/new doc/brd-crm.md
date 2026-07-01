# BRD Final - CRM Module

## 1. Visi & Kebutuhan Bisnis
Modul CRM mengoptimalkan pencatatan leads/calon pendaftar dari berbagai kanal pemasaran (sosial media, agen kemitraan, referral, event) dan memantau status tindak lanjut (*follow-up pipeline*) hingga dikonversi menjadi pendaftar PMB.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-CRM-B-001** | Lead Funnel Tracking | Memantau siklus hidup prospek dari kontak pertama hingga siap daftar PMB. | P0 | Terbentuknya laporan konversi leads di dashboard pimpinan secara real-time. |
| **BRD-CRM-B-002** | Kemitraan Agen Referral | Mengelola agen rujukan eksternal dan kode rujukan unik untuk promosi. | P0 | Otomatisasi generate kode referral agen yang valid. |
| **BRD-CRM-B-003** | Kalkulasi Komisi Rujukan | Menghitung komisi agen rujukan ketika calon mahasiswa yang dirujuk melunasi biaya pendaftaran. | P1 | Keakuratan perhitungan nominal komisi tanpa data ganda (100% akurat). |
| **BRD-CRM-B-004** | Log Aktivitas Follow-up | Menyimpan log aktivitas panggilan telepon, chat, atau email tindak lanjut marketing. | P1 | Riwayat interaksi prospek tercatat runut dan membantu alokasi tugas tim sales. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-CRM-001**: Leads hanya dapat dikonversi menjadi applicant jika statusnya sudah diubah menjadi `Qualified`.
* **BR-CRM-002**: Pengiriman data konversi harus membawa token unik yang diuji keunikannya di sisi PMB secara idempotent.
* **BR-CRM-003**: Komisi rujukan agen hanya akan diproses cair ke modul Keuangan jika status pendaftar yang dirujuk dinyatakan lunas UKT semester pertama.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Pendaftaran prospek baru, update status follow-up marketing, input skema komisi agen.
* **Output**: Draf pendaftar PMB, laporan komisi bulanan agen kemitraan.
* **Integrasi Lintas Domain**: CRM mengonsumsi event pembayaran lunas pendaftar PMB dari Finance untuk memicu rekonsiliasi dan perhitungan nominal komisi rujukan agen.
