# BRD Final - Akademik Module

## 1. Visi & Kebutuhan Bisnis
Modul Akademik (SIAKAD) bertindak sebagai jantung operasional perkuliahan mahasiswa, kurikulum prodi, persetujuan KRS oleh dosen pembimbing akademik, serta pelaporan kelulusan nilai mahasiswa (KHS/Transkrip).

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-ACA-B-001** | NIM Generator Otomatis | Menghasilkan NIM mahasiswa baru secara sistematis berdasarkan kode prodi dan urutan pendaftaran. | P0 | Tidak ada duplikasi NIM pada angkatan yang sama (0% error). |
| **BRD-ACA-B-002** | Pengisian KRS Online | Mahasiswa dapat memilih kelas penawaran mata kuliah semester secara mandiri sesuai kapasitas SKS. | P0 | Validasi bentrok jadwal kuliah dan prasyarat MK berjalan otomatis. |
| **BRD-ACA-B-003** | Approval KRS Dosen PA | Dosen Pembimbing Akademik dapat meninjau dan menyetujui draf rencana studi KRS mahasiswa bimbingan. | P0 | SLA persetujuan KRS oleh dosen PA di bawah 3 hari kerja. |
| **BRD-ACA-B-004** | Kartu Hasil Studi (KHS) | Menerbitkan nilai akhir mata kuliah mahasiswa beserta perhitungan IPS dan IPK akumulatif. | P0 | Penguncian nilai terjamin aman; riwayat koreksi nilai terekam audit. |
| **BRD-ACA-B-005** | Yudisium & Alumni | Mengelola status kelayakan kelulusan mahasiswa untuk dipindahkan menjadi data alumni resmi. | P1 | Otomatisasi validasi kecukupan SKS kelulusan prodi. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-ACA-001**: Mahasiswa semester 1 dan 2 wajib menggunakan pengisian KRS Paket, sedangkan mahasiswa semester 3 ke atas menggunakan KRS Mandiri.
* **BR-ACA-002**: Pengisian KRS hanya diperbolehkan jika status clearance keuangan mahasiswa berstatus `CLEARED` atau `CONDITIONAL` (dispensasi aktif).
* **BR-ACA-003**: Dosen dilarang mengubah nilai akhir mata kuliah mahasiswa jika kelas semester tersebut sudah dinyatakan ditutup (*closed*), kecuali melalui pengajuan koreksi resmi yang disetujui Kaprodi.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Data handover pendaftar baru dari PMB, status clearance dari Finance, data dosen pengampu aktif dari HRIS.
* **Output**: Rencana studi disetujui, transkrip nilai KHS resmi, status kelulusan yudisium.
* **Integrasi Lintas Domain**: Event persetujuan KRS `academic.krs_approved` langsung memicu pembuatan kelas online dan pendaftaran peserta di LMS secara eventual konsisten.
