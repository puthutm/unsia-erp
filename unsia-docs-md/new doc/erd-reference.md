# ERD Final - Referensi Module

## 1. Skema Database
Database `reference_db` menyimpan master referensi data standar seperti program studi, kalender akademik (tahun ajaran & periode semester), gelombang PMB, komponen biaya tagihan, metode bayar, dan status kode global.

## 2. Tabel Utama
* **`ref.study_programs`**: Master data program studi di lingkungan UNSIA.
* **`ref.academic_years`**: Master tahun ajaran kalender operasional.
* **`ref.academic_periods`**: Master periode akademik semester aktif (ganjil, genap, pendek).
* **`ref.pmb_waves`**: Master data gelombang pendaftaran PMB.
* **`ref.payment_components`**: Komponen biaya tagihan (UKT, uang gedung, biaya pendaftaran).
* **`ref.payment_methods`**: Metode pembayaran yang didukung.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table ref.study_programs {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  degree_level varchar
  faculty_name varchar
  is_active boolean
}

Table ref.academic_years {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  start_year int
  end_year int
  status varchar
  is_active boolean
}

Table ref.academic_periods {
  id uuid [pk]
  academic_year_id uuid [not null]
  code varchar [unique, not null]
  name varchar [not null]
  semester_type varchar
  start_date date
  end_date date
  status varchar
  is_active boolean
}

Table ref.pmb_waves {
  id uuid [pk]
  academic_year_id uuid
  target_entry_period_id uuid [not null]
  admission_path_id uuid
  code varchar [unique, not null]
  name varchar [not null]
  registration_start_at timestamp
  registration_end_at timestamp
  status varchar
  is_active boolean
}

Table ref.payment_components {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  component_type varchar
  is_active boolean
}

Table ref.payment_methods {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  provider varchar
  is_active boolean
}

Ref: ref.academic_periods.academic_year_id > ref.academic_years.id
Ref: ref.pmb_waves.academic_year_id > ref.academic_years.id
Ref: ref.pmb_waves.target_entry_period_id > ref.academic_periods.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `ref.academic_periods.id` dan `ref.study_programs.id` disimpan sebagai external reference di database Akademik (`academic_db`), PMB (`pmb_db`), LMS (`lms_db`), dan Finance (`finance_db`).
