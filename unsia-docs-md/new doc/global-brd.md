# Business Requirement Document (BRD) Global - UNSIA ERP (Developer & Business Edition)

## 1. Visi & KPI Bisnis Ekosistem ERP
Membangun platform sistem operasional akademik terintegrasi yang andal, aman, dan fleksibel untuk mendampingi siklus hidup mahasiswa dari calon pendaftar, mahasiswa aktif, proses belajar online, evaluasi kelayakan UKT, hingga alumni.

---

## 2. Matriks Kasus Penggunaan & Aktor (Use Case Matrix & Actors Mapping)

| ID Use Case | Aktor Utama | Nama Proses Bisnis | Prasyarat (Pre-Condition) | Pasca-Syarat (Post-Condition) |
| --- | --- | --- | --- | --- |
| **UC-G-001** | Calon Mahasiswa | Pendaftaran Akun PMB | Membuka portal pendaftaran publik | Draf akun applicant terbentuk di `pmb_db` |
| **UC-G-002** | Calon Mahasiswa | Pembayaran Formulir PMB | Tagihan invoice formulir terbit | Status tagihan lunas, hak seleksi terbuka |
| **UC-G-003** | Calon Mahasiswa | Ujian Seleksi CBT | Dokumen terverifikasi & tagihan formulir lunas | Skor kelulusan terhitung di Assessment |
| **UC-G-004** | Admin PMB | Handover Mahasiswa Baru | Calon mahasiswa melunasi daftar ulang & lulus CBT | NIM terbit, akun mahasiswa aktif di SIAKAD |
| **UC-G-005** | Mahasiswa | Pengisian KRS Mandiri | Status clearance keuangan lunas / dispensasi aktif | Draf KRS diajukan ke Dosen PA |
| **UC-G-006** | Dosen Wali (PA) | Persetujuan KRS | Antrean draf KRS masuk ke dosen wali | Status KRS disetujui, auto-provisioning LMS |
| **UC-G-007** | Dosen Pengampu | Finalisasi Nilai Semester | Nilai tugas LMS tersinkronisasi | Nilai akhir terkunci, KHS terbit untuk mahasiswa |

---

## 3. Aturan Bisnis Kritis Terperinci (Detailed Business Rules & Policies)

### A. Kebijakan Kelayakan Pengisian KRS (Academic Clearance Policy)
* **Aturan**: Mahasiswa hanya diperbolehkan masuk ke halaman pengisian KRS mandiri semester aktif jika status keuangannya berstatus `CLEARED` (tidak ada tunggakan jatuh tempo) atau `CONDITIONAL` (memiliki dispensasi aktif).
* **Validasi Sistem**:
  ```python
  if student.clearance_status not in ['CLEARED', 'CONDITIONAL']:
      raise AcademicClearanceBlockException("Akses pengisian KRS ditutup. Selesaikan tagihan UKT Anda atau ajukan dispensasi cicilan.")
  ```

### B. Aturan Pembuatan NIM Otomatis (NIM Issuance Rule)
* **Format NIM**: `[Kode_Jenjang][Kode_Tahun_Masuk][Kode_Prodi][Nomor_Urut_4_Digit]`
  * Contoh: `2601010045` (26 = angkatan 2026, 01 = jenjang S1, 01 = prodi Informatika, 0045 = urutan pendaftar ke-45 yang melunasi daftar ulang).
* **Aturan Kunci**: NIM tidak boleh dirilis sebelum status pembayaran daftar ulang calon mahasiswa dinyatakan sukses lunas oleh modul Finance.

### C. Batas Kuota SKS KRS Mandiri (SKS Credit Capacity Policy)
Batas beban studi yang boleh diambil mahasiswa dalam KRS Mandiri semester berjalan ditentukan secara mutlak berdasarkan perolehan Indeks Prestasi Semester (IPS) sebelumnya:
* **IPS $\ge$ 3.00**: Maksimal mengambil **24 SKS**.
* **IPS 2.50 - 2.99**: Maksimal mengambil **21 SKS**.
* **IPS 2.00 - 2.49**: Maksimal mengambil **18 SKS**.
* **IPS < 2.00**: Maksimal mengambil **15 SKS**.

---

## 4. Kamus Istilah Bisnis (Business Glossaries)
* **NIM (Nomor Induk Mahasiswa)**: Identitas unik berskala nasional sebagai bukti mahasiswa terdaftar aktif di perguruan tinggi.
* **NIDN (Nomor Induk Dosen Nasional)**: Nomor identitas tunggal bagi dosen yang terdaftar di Kementerian Pendidikan.
* **KRS (Kartu Rencana Studi)**: Formulir rencana pemilihan mata kuliah kelas perkuliahan yang wajib diisi mahasiswa di awal semester.
* **KHS (Kartu Hasil Studi)**: Lembar nilai akhir hasil pembelajaran mahasiswa yang diterbitkan di akhir semester.
* **LoA (Letter of Acceptance)**: Surat keputusan resmi pernyataan penerimaan calon mahasiswa baru dari institusi.
* **Dispensasi UKT**: Kelonggaran waktu pembayaran biaya kuliah semester yang disetujui biro keuangan melalui skema cicilan khusus.
* **Yudisium**: Ulasan kelayakan nilai akademik kumulatif mahasiswa sebagai syarat sah kelulusan sebelum wisuda.
* **BKD (Beban Kerja Dosen)**: Rekam jejak pemenuhan tridharma perguruan tinggi dosen mencakup mengajar, meneliti, dan mengabdi masyarakat.
