# PRD Final - Assessment Module

## 1. Pendahuluan & Latar Belakang
Modul Assessment bertanggung jawab menyediakan engine ujian berbasis komputer (CBT), bank soal ujian (*question bank*), pencatatan versi soal, penjadwalan sesi ujian CBT (UTS/UAS/Ujian Masuk PMB), penyimpanan lembar jawaban peserta, kalkulasi otomatis skor ujian, dan penerbitan hasil ujian kelulusan.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `assessment_db`
* **Source of Truth (Kepemilikan Data)**:
  * `assessment.question_banks` (Bank soal).
  * `assessment.questions`, `assessment.question_options` & `assessment.question_versions` (Butir soal, opsi jawaban, dan histori revisi versi soal).
  * `assessment.material_banks` & `assessment.materials` (Kumpulan materi referensi soal).
  * `assessment.question_sets` & `assessment.question_set_items` (Paket kompilasi soal ujian).
  * `assessment.assessment_sessions` (Sesi pelaksanaan ujian CBT).
  * `assessment.assessment_participants` (Daftar peserta terdaftar sesi).
  * `assessment.assessment_attempts` (Log pengerjaan lembar ujian).
  * `assessment.assessment_answers` (Rekam jawaban peserta per nomor).
  * `assessment.assessment_scores` (Nilai total kelulusan ujian).
  * `assessment.surveys`, `assessment.survey_questions` & `assessment.survey_responses` (Mesin kuesioner/survei).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `participant_snapshots`: Salinan profil pendaftar/mahasiswa dari PMB/Akademik untuk keperluan visualisasi kartu ujian di layar CBT.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-ASM-001** | P0 | Desain CBT Engine Mandiri. | Assessment menyimpan attempt, answer, dan score secara mandiri tanpa bergantung pada database modul pemanggil. |
| **PRD-ASM-002** | P0 | Penyelarasan Hasil Ujian Aman. | Skor ujian dikirim ke consumer via event `assessment.result_calculated` dengan jaminan retry policy jika target down. |
| **PRD-ASM-003** | P1 | Keamanan Ujian Shuffling Server. | Pengacakan urutan soal dan durasi hitung mundur waktu ujian wajib dikendalikan secara mutlak di sisi backend/server. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `assessment.result_calculated`: Dikirim setelah peserta menekan tombol submit ujian atau durasi pengerjaan habis.
    * *Payload Minimum*: `assessment_session_id`, `participant_type`, `participant_ref_id`, `total_score`, `passed`.

## 6. Degraded Mode & Resilience Guardrails
* **Context Outage Resilience**: Jika modul pemanggil (PMB/LMS) down, peserta yang sudah terdaftar tetap dapat melangsungkan ujian CBT secara normal. Ekspor data hasil ujian akan antre di outbox Assessment dan otomatis terkirim setelah modul tujuan kembali online.
