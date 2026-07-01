# BRD Final - Referensi Module

## 1. Visi & Kebutuhan Bisnis
Modul Referensi menyajikan data master standar baku yang digunakan bersama di seluruh sistem (seperti daftar program studi, periode akademik/semester aktif, dll.). Tujuannya adalah mencegah redudansi data master dan inkonsistensi penulisan nama entitas umum.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-REF-B-001** | Master Program Studi | Menyediakan satu daftar nama program studi resmi yang terintegrasi untuk PMB, SIAKAD, dan LMS. | P0 | Tidak ada perbedaan nama prodi antar modul; data sync terjamin. |
| **BRD-REF-B-002** | Kalender Akademik Terpadu | Mengatur periode pengisian KRS, pendaftaran PMB, dan rentang tanggal perkuliahan semester. | P0 | Semua modul mematuhi periode aktif yang sama untuk memvalidasi transaksi baru. |
| **BRD-REF-B-003** | Standarisasi Komponen Biaya | Mengelola daftar resmi komponen biaya pendidikan (UKT, uang gedung, biaya pendaftaran). | P0 | Seluruh invoice yang terbit hanya merujuk pada komponen biaya resmi dari Referensi. |
| **BRD-REF-B-004** | Master Kode Status | Menyediakan daftar kode status (draft, submitted, cleared, dll.) untuk standardisasi status bisnis di ERP. | P1 | Tidak ada status string bebas; transisi status mematuhi kamus kode status. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-REF-001**: Kalender operasional akademik (Tahun Ajaran/Periode) dipisahkan secara tegas dari versi kurikulum prodi.
* **BR-REF-002**: Pengaktifan periode akademik baru otomatis memicu status KRS semester sebelumnya menjadi terkunci (*closed*).
* **BR-REF-003**: Penghapusan fisik (*hard delete*) pada data master referensi dilarang keras jika data tersebut sudah memiliki relasi di database modul lain.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Pendaftaran prodi baru oleh BPPTI, pembukaan gelombang masuk PMB baru oleh Biro Akademik.
* **Output**: Dropdown data referensi terintegrasi, draf kalender perkuliahan semester.
* **Integrasi Lintas Domain**: Event perubahan prodi/periode disalurkan ke seluruh database modul ERP untuk sinkronisasi snapshot.
