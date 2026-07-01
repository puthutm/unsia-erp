# ERD Final - Akademik Module

## 1. Skema Database
Database `academic_db` (SIAKAD) bertindak sebagai pusat penyimpanan transaksional data kemahasiswaan, penawaran kurikulum, jadwal perkuliahan, pengambilan KRS, penyerahan nilai KHS, dan yudisium kelulusan.

## 2. Tabel Utama
* **`academic.students`**: Data induk profil mahasiswa terdaftar.
* **`academic.curriculums`**: Data kurikulum pendidikan per prodi.
* **`academic.courses`**: Data daftar mata kuliah.
* **`academic.course_offerings`**: Konteks penawaran mata kuliah per periode semester.
* **`academic.classes`**: Pembukaan kelas kuliah paralel.
* **`academic.krs`**: Data induk Kartu Rencana Studi (KRS) per mahasiswa per semester.
* **`academic.krs_items`**: Pilihan kelas kuliah yang diambil mahasiswa dalam KRS.
* **`academic.grades`**: Transkrip nilai akhir mata kuliah mahasiswa.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table academic.students {
  id uuid [pk]
  person_id uuid [not null]
  user_id uuid
  applicant_id uuid
  study_program_id uuid [not null]
  nim varchar [unique, not null]
  student_status varchar
  entry_period_id uuid
  curriculum_id uuid
  current_semester int
}

Table academic.curriculums {
  id uuid [pk]
  study_program_id uuid [not null]
  code varchar [unique, not null]
  name varchar [not null]
  curriculum_year int [not null]
  is_active boolean
}

Table academic.courses {
  id uuid [pk]
  study_program_id uuid
  course_code varchar [unique, not null]
  course_name varchar [not null]
  sks int
  is_active boolean
}

Table academic.course_offerings {
  id uuid [pk]
  course_id uuid [not null]
  academic_period_id uuid [not null]
  curriculum_id uuid
  status varchar
}

Table academic.classes {
  id uuid [pk]
  course_offering_id uuid [not null]
  class_code varchar [not null]
  quota int
  enrolled_count int
}

Table academic.krs {
  id uuid [pk]
  student_id uuid [not null]
  academic_period_id uuid [not null]
  status varchar
  finance_clearance_id uuid
  submitted_at timestamp
  approved_at timestamp
}

Table academic.krs_items {
  id uuid [pk]
  krs_id uuid [not null]
  class_id uuid [not null]
  status varchar
}

Table academic.grades {
  id uuid [pk]
  krs_item_id uuid [not null]
  numeric_grade numeric
  letter_grade varchar
  grade_point numeric
  source varchar
}

Ref: academic.course_offerings.course_id > academic.courses.id
Ref: academic.classes.course_offering_id > academic.course_offerings.id
Ref: academic.krs.student_id > academic.students.id
Ref: academic.krs_items.krs_id > academic.krs.id
Ref: academic.krs_items.class_id > academic.classes.id
Ref: academic.grades.krs_item_id - academic.krs_items.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `academic.students.person_id` merujuk logis ke `core_db.persons.id`.
* `academic.students.study_program_id` merujuk logis ke `reference_db.study_programs.id`.
* `academic.krs.finance_clearance_id` merujuk logis ke `finance_db.student_clearances.id`.
