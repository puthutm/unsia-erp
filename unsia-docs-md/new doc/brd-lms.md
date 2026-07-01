# BRD Final - LMS Module

## 1. Visi & Kebutuhan Bisnis
Modul LMS memfasilitasi kegiatan belajar mengajar online secara interaktif, menyediakan ruang materi ajar, pengunggahan tugas mingguan, kuis mandiri, forum diskusi dosen-mahasiswa, serta tracking keaktifan belajar online.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-LMS-B-001** | Kelas Belajar Online | Menyediakan ruang belajar digital interaktif per kelas mata kuliah per semester berjalan. | P0 | Pembentukan ruang kelas otomatis sejalan dengan jadwal SIAKAD. |
| **BRD-LMS-B-002** | Pengunggahan Materi | Dosen dapat mengunggah file video, modul ajar, dan tautan video conference kelas kuliah. | P0 | Kemudahan akses materi ajar oleh mahasiswa secara responsif. |
| **BRD-LMS-B-003** | Tugas & Umpan Balik | Mengelola pengumpulan berkas jawaban tugas mahasiswa dan penilaian oleh dosen. | P0 | SLA koreksi nilai tugas oleh dosen sebelum batas semester berakhir. |
| **BRD-LMS-B-004** | Keaktifan Belajar Mhs | Mengukur tingkat keterlibatan mahasiswa berdasarkan materi yang dibaca dan kehadiran sesi. | P1 | Terbentuknya visual progress persentase ketuntasan belajar mahasiswa. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-LMS-001**: Mahasiswa dilarang mengakses kelas online LMS jika status KRS mata kuliah tersebut belum disetujui dosen PA di SIAKAD.
* **BR-LMS-002**: Pengiriman nilai tugas ke SIAKAD wajib mengunci input nilai agar tidak dapat dirubah setelah melewati batas penutupan semester aktif.
* **BR-LMS-003**: Dosen pembantu pengampu kelas LMS hanya memiliki izin untuk mengunggah materi dan dilarang mengubah bobot penilaian tugas utama kelas.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Data draf kelas dari SIAKAD, draf peserta dari KRS disetujui, berkas tugas mahasiswa, kuis dari Assessment.
* **Output**: Log keaktifan belajar mahasiswa, draf nilai tugas perkuliahan.
* **Integrasi Lintas Domain**: LMS menyalurkan hasil penilaian tugas mingguan ke modul SIAKAD via event broker untuk dimasukkan sebagai bagian komponen perhitungan nilai akhir perkuliahan.
