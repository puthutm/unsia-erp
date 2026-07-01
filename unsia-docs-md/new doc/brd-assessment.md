# BRD Final - Assessment Module

## 1. Visi & Kebutuhan Bisnis
Modul Assessment menyajikan layanan ujian online (*computer-based testing*) yang reusable untuk berbagai konteks di lingkungan kampus, seperti ujian masuk PMB, kuis belajar LMS, maupun ujian semester UTS/UAS SIAKAD.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-ASM-B-001** | Bank Soal Terintegrasi | Mengelola koleksi soal ujian, kategori mata uji, bobot kesulitan, dan pengacakan butir soal. | P0 | Tidak ada kebocoran kunci jawaban ujian; enkripsi data terjamin. |
| **BRD-ASM-B-002** | Attempt CBT Engine | Menyediakan antarmuka pengerjaan ujian interaktif dengan fitur penanda ragu-ragu dan autosave. | P0 | autosave jawaban peserta per nomor di bawah 2 detik (anti data hilang). |
| **BRD-ASM-B-003** | Koreksi Nilai Ujian | Mengkalkulasikan skor ujian secara instan berdasarkan kunci jawaban terprogram. | P0 | Keakuratan kalkulasi nilai 100% akurat sesuai pembobotan. |
| **BRD-ASM-B-004** | Kuesioner & Survei | Mengakomodasi kebutuhan kuesioner evaluasi dosen oleh mahasiswa atau survei kepuasan mahasiswa. | P1 | Dukungan mode survei anonim untuk menjamin keaslian respon. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-ASM-001**: Peserta dilarang memulai sesi ujian CBT jika status pendaftarannya belum dinyatakan lolos verifikasi berkas di PMB atau KRS belum approved di SIAKAD.
* **BR-ASM-002**: Durasi pengerjaan kuis dihitung mundur secara ketat dari waktu klik mulai (*attempt start*), dan pengerjaan akan ter-submit paksa otomatis jika durasi habis.
* **BR-ASM-003**: Draf butir soal yang sedang aktif digunakan dalam sesi ujian berjalan dilarang keras diubah secara langsung (wajib membuat versi soal baru).

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Assign draf peserta seleksi PMB/kelas LMS, draf bank soal dosen, log klik pilihan jawaban peserta.
* **Output**: Skor total kelulusan ujian, detail lembar audit jawaban peserta.
* **Integrasi Lintas Domain**: Modul Assessment menyalurkan skor akhir kelulusan ujian peserta lewat event broker untuk langsung merubah status kelulusan PMB calon mahasiswa.
