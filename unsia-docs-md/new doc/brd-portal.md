# BRD Final - Portal Module

## 1. Visi & Kebutuhan Bisnis
Modul Portal menyajikan gerbang utama interaksi pengguna, mengumpulkan notifikasi aktivitas bisnis penting, menyimpan preferensi tampilan dashboard, serta menyajikan grafik metrik operasional terpadu (*KPI dashboard*) bagi jajaran pimpinan kampus.

## 2. Kebutuhan Bisnis Detail (Business Requirements)

| ID Kebutuhan | Nama Kebutuhan | Deskripsi Kebutuhan Bisnis | Prioritas | Indikator Keberhasilan (KPI) |
| --- | --- | --- | --- | --- |
| **BRD-POR-B-001** | Dashboard Peran Terpadu | Menyajikan halaman beranda personal yang dinamis sesuai peran aktif (Mahasiswa/Dosen/Pimpinan). | P0 | Desain responsif, waktu pemuatan halaman di bawah 1.5 detik. |
| **BRD-POR-B-002** | Notification Center | Mengirimkan push notification real-time untuk info tagihan, persetujuan KRS, dan nilai keluar. | P0 | Notifikasi tersalurkan ke WebSocket/SSE di bawah 2 detik setelah event terjadi. |
| **BRD-POR-B-003** | Menu Shortcut Mandiri | Pengguna dapat menyusun pintasan navigasi ke menu yang paling sering diakses. | P1 | Kemudahan reorder urutan pintasan menu di halaman dashboard. |
| **BRD-POR-B-004** | Dashboard KPI Pimpinan | Menyajikan ringkasan grafik statistik PMB, Keuangan, dan LMS untuk kebutuhan analisis eksekutif. | P1 | Angka ringkasan KPI sesuai dengan data transaksional di database modul asal. |

## 3. Aturan Bisnis Modul (Business Rules)
* **BR-POR-001**: Portal tidak diperkenankan menyimpan salinan kredensial password pengguna; seluruh sesi login memanfaatkan session auth Core.
* **BR-POR-002**: Pengiriman notifikasi ke inbox wajib mencatat status baca notifikasi per user secara unik.
* **BR-POR-003**: Metrik data dashboard pimpinan tidak boleh diperlakukan sebagai dasar keputusan audit final sebelum dilakukan proses rekonsiliasi data.

## 4. Aliran Informasi & Integrasi Bisnis
* **Input**: Event notifikasi dari seluruh modul transaksi ERP, pengaturan preferensi visual user.
* **Output**: Distribusi push notification, visualisasi grafik KPI operasional kampus.
* **Integrasi Lintas Domain**: Portal menangkap event notifikasi bisnis penting lewat antrean broker untuk langsung disalurkan sebagai pesan pop-up/inbox kepada pengguna yang berhak.
