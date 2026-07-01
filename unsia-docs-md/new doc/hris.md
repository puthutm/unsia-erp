# PRD Final - HRIS Module

## 1. Pendahuluan & Latar Belakang
Modul HRIS bertindak sebagai authority untuk mengelola database karyawan, dosen pembimbing akademik/pengampu, struktur jabatan organisasi unit, pencatatan absensi, pengajuan cuti, perhitungan BKD (Beban Kerja Dosen), dan parameter gaji pokok.

## 2. Batasan Database (Database Boundary) & Ownership
* **Nama Database**: `hris_db`
* **Source of Truth (Kepemilikan Data)**:
  * `hris.work_units` (Unit organisasi / fakultas / program studi).
  * `hris.positions` & `hris.functional_positions` (Jabatan struktural & fungsional dosen/karyawan).
  * `hris.employees` (Biodata induk karyawan).
  * `hris.lecturers` (Data dosen pengampu aktif & NIDN).
  * `hris.attendances` (Catatan presensi masuk/keluar harian).
  * `hris.leave_requests` (Pengajuan dispensasi cuti tahunan).
  * `hris.bkd_records` (Pencatatan Beban Kerja Dosen per semester).
  * `hris.performance_reviews` & `hris.certifications` (Ulasan penilaian kinerja & sertifikasi).
  * `hris.payroll_sources` (Konfigurasi parameter gaji pokok kepegawaian).

## 3. Data Lintas Modul (Snapshot/Read Model)
* `person_snapshots`: Salinan profil diri staf dari Core untuk mempercepat lookup tanpa cross-DB query.

## 4. Persyaratan Produk & Acceptance Criteria (AC)

| ID Persyaratan | Prioritas | Deskripsi Persyaratan | Kriteria Penerimaan (Acceptance Criteria) |
| --- | --- | --- | --- |
| **PRD-HRIS-001** | P0 | Otoritas Tunggal SDM Kampus. | HRIS bertindak sebagai pemilik data tunggal dosen dan karyawan. Modul Akademik/LMS mengambil data lewat snapshot. |
| **PRD-HRIS-002** | P1 | Event Perubahan Status Dosen. | Penonaktifan dosen pengampu (karena cuti panjang/pensiun) mempublikasikan event status agar plotting dibatalkan secara sistem. |
| **PRD-HRIS-003** | P1 | Integrasi Kehadiran Absensi. | Pencatatan log check-in/check-out harus terjamin valid berdasarkan parameter jam kerja aktif. |

## 5. Event-Driven Integration (Event Catalog)
* **Event yang Diterbitkan (Publisher)**:
  * `hris.lecturer_status_updated`: Dikirim saat keaktifan mengajar dosen diubah.
    * *Payload Minimum*: `lecturer_id`, `nidn`, `is_active`, `occurred_at`.

## 6. Degraded Mode & Resilience Guardrails
* **HRIS Outage Resilience**: Jika database HRIS mengalami gangguan, SIAKAD dan LMS tetap dapat berjalan normal menggunakan data `lecturer_snapshots` terakhir yang tersimpan secara lokal. Penugasan dosen wali baru ditangguhkan sampai HRIS kembali online.
