# BRD Final - Finance Module

## 1. Visi & Kebutuhan Bisnis
Modul Finance menjamin transparansi dan keakuratan pengelolaan keuangan kampus, mencakup pembuatan invoice tagihan, penerimaan pembayaran dari payment gateway luar, rekonsiliasi kas bank, jurnal akuntansi pembukuan berpasangan, hingga evaluasi status pembebasan akademik mahasiswa.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-FIN-B-001** | Otomatisasi Invoice | Mampu membuat tagihan invoice UKT semesteran dan biaya pendaftaran PMB secara sistematis. | P0 | Tidak ada tagihan manual di luar skema komponen biaya resmi. |
| **BRD-FIN-B-002** | Payment Gateway Callback | Menerima instruksi lunas pembayaran dari gateway pembayaran VA/QRIS luar secara instan. | P0 | Verifikasi status pembayaran di bawah 5 detik setelah uang masuk. |
| **BRD-FIN-B-003** | Pembukuan Jurnal Akuntansi | Mencatatkan entri debit-kredit secara otomatis ke CoA untuk setiap tagihan dan pembayaran masuk. | P0 | Keseimbangan neraca saldo terjamin 100% akurat. |
| **BRD-FIN-B-004** | Clearance Finansial | Menyekat akses pengisian KRS dan cetak KHS mahasiswa jika memiliki tunggakan yang melewati batas toleransi. | P0 | Keamanan aset institusi terlindungi dari kebocoran piutang macet. |
| **BRD-FIN-B-005** | Cicilan UKT & Dispensasi | Mengelola pengajuan dispensasi cicilan dengan termin bayar yang disepakati. | P1 | Mahasiswa dengan dispensasi aktif dapat tetap mengisi KRS (status conditional). |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-FIN-001**: Nilai total invoice tagihan adalah penjumlahan seluruh komponen biaya dikurangi nominal beasiswa/potongan diskon yang sah.
* **BR-FIN-002**: Status kelayakan akademik mahasiswa dinyatakan `BLOCKED` jika terdapat tagihan jatuh tempo yang belum diselesaikan tanpa dispensasi aktif.
* **BR-FIN-003**: Log callback dari payment gateway luar dilarang diproses ulang apabila ID transaksi provider tersebut sudah tercatat lunas di database.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Request tagihan mahasiswa baru dari PMB, request tagihan UKT semester dari SIAKAD, data slip gaji dari HRIS.
* **Output**: Status clearance terbaru, kwitansi pembayaran resmi, laporan keuangan jurnal kas.
* **Integrasi Lintas Domain**: Finance mengirimkan status clearance terbaru via event broker untuk membatasi atau mengizinkan akses pengisian KRS di SIAKAD.
