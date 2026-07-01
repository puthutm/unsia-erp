# BRD Final - HRIS Module

## 1. Visi & Kebutuhan Bisnis
Modul HRIS bertindak sebagai authority untuk mengelola database karyawan, dosen pembimbing akademik/pengampu, struktur jabatan organisasi unit, pencatatan absensi, pengajuan cuti, perhitungan BKD (Beban Kerja Dosen), dan parameter gaji pokok.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-HRIS-B-001** | Direktori Dosen & Karyawan | Mengelola profil kepegawaian, jabatan struktural, dan sertifikasi keahlian dosen secara terpusat. | P0 | Tidak ada duplikasi record NIDN dosen (0% error). |
| **BRD-HRIS-B-002** | Absensi Kerja Mandiri | Membantu karyawan mencatatkan jam masuk dan jam pulang kerja harian secara mandiri di portal. | P0 | Log kehadiran terekam aman tanpa celah manipulasi jam absen. |
| **BRD-HRIS-B-003** | Pengajuan Cuti Mandiri | Memfasilitasi alur pengajuan izin cuti tahunan staf beserta sistem approval atasan. | P1 | SLA approval persetujuan cuti di bawah 2 hari kerja. |
| **BRD-HRIS-B-004** | BKD Dosen per Semester | Memantau kinerja tridharma perguruan tinggi dosen pengampu per semester untuk pelaporan eksternal. | P1 | Otomatisasi rekap beban kerja mengajar dari jadwal SIAKAD. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-HRIS-001**: Setiap penunjukan dosen wali pembimbing KRS wajib menggunakan record NIDN dosen berstatus aktif di HRIS.
* **BR-HRIS-002**: Pengajuan cuti tahunan akan memotong sisa kuota cuti secara otomatis sesaat setelah admin HRIS menerbitkan persetujuan final.
* **BR-HRIS-003**: Pegawai yang berstatus nonaktif (pensiun/mengundurkan diri) otomatis langsung kehilangan akses SSO akun login portal ERP.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Data profil karyawan baru, log data check-in harian, berkas sertifikasi kompetensi dosen.
* **Output**: Status keaktifan mengajar dosen, slip rekap presensi bulanan, data usulan nominal dasar payroll.
* **Integrasi Lintas Domain**: HRIS menyebarkan status keaktifan dosen via event broker untuk membatasi penugasan dosen pengampu mata kuliah baru di modul SIAKAD.
