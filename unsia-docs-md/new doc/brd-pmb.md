# BRD Final - PMB Module

## 1. Visi & Kebutuhan Bisnis
Modul Penerimaan Mahasiswa Baru (PMB) melayani proses seleksi penerimaan mahasiswa secara transparan dan teratur, mulai dari pengisian biodata diri, pemenuhan berkas administrasi, ujian CBT masuk, hingga handover data mahasiswa baru ke modul SIAKAD.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-PMB-B-001** | Pengisian Biodata Mandiri | Calon mahasiswa dapat mengisi biodata pribadi, riwayat sekolah, keluarga, dan profil finansial secara mandiri. | P0 | 100% kelengkapan data sebelum masuk ke tahap ujian seleksi. |
| **BRD-PMB-B-002** | Verifikasi Berkas Online | Admin PMB dapat memvalidasi dokumen syarat masuk pendaftaran pendaftar. | P0 | SLA verifikasi berkas di bawah 2 hari kerja. |
| **BRD-PMB-B-003** | Integrasi Pembayaran | Pendaftaran langsung terhubung dengan status tagihan keuangan pendaftaran secara aman. | P0 | Tidak ada pendaftar yang masuk ke tahap ujian CBT sebelum melunasi tagihan biaya formulir. |
| **BRD-PMB-B-004** | Penerbitan Surat Kelulusan | Sistem secara otomatis meng-issue Letter of Acceptance (LoA) setelah pendaftar dinyatakan lulus seleksi. | P1 | Kecepatan cetak surat LoA digital secara instan oleh pendaftar. |
| **BRD-PMB-B-005** | Serah Terima Mahasiswa | Handover mahasiswa lulus dari PMB ke SIAKAD dilakukan otomatis secara sistem. | P0 | Tidak ada duplikasi NIM atau kesalahan prodi saat proses handover mahasiswa baru. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-PMB-001**: Calon mahasiswa baru hanya dapat mengklik tombol "Daftar Ulang" jika status seleksi ujian masuknya dinyatakan lulus.
* **BR-PMB-002**: Pengiriman data handover mahasiswa ke SIAKAD hanya dapat dieksekusi apabila status registrasi dan verifikasi berkas daftar ulangnya berstatus disetujui.
* **BR-PMB-003**: Dilarang menghapus data pendaftar yang sudah melakukan transaksi keuangan/pembayaran di Finance.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Pendaftaran calon mahasiswa, unggahan dokumen scan ijazah/KTP, data skor kelulusan CBT dari Assessment.
* **Output**: Surat LoA digital, draf mahasiswa baru siap impor ke SIAKAD.
* **Integrasi Lintas Domain**: Event handover `pmb.ready_for_academic` ditangkap oleh SIAKAD untuk mendaftarkan NIM baru dan melacak rekam jejak histori PMB calon mahasiswa.
