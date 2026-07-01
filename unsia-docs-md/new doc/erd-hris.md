# ERD Final - HRIS Module

## 1. Skema Database
Database `hris_db` mengelola data karyawan, dosen, data presensi kerja, serta permohonan cuti pegawai.

## 2. Tabel Utama
* **`hris.employees`**: Profil karyawan utama institusi.
* **`hris.lecturers`**: Detail khusus fungsional dan kompetensi dosen.
* **`hris.attendances`**: Catatan presensi masuk dan keluar harian.
* **`hris.leave_requests`**: Log permohonan dispensasi cuti.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table hris.work_units {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table hris.employees {
  id uuid [pk]
  person_id uuid [not null]
  employee_type_id uuid
  work_unit_id uuid
  position_id uuid
  nip varchar [unique]
  employment_status varchar
  join_date date
  is_active boolean
}

Table hris.lecturers {
  id uuid [pk]
  employee_id uuid [not null]
  lecturer_status_id uuid
  functional_position_id uuid
  nidn varchar [unique]
  homebase_study_program_id uuid
  is_active boolean
}

Table hris.attendances {
  id uuid [pk]
  employee_id uuid [not null]
  attendance_date date
  check_in time
  check_out time
  status varchar
}

Table hris.leave_requests {
  id uuid [pk]
  employee_id uuid [not null]
  leave_type varchar
  start_date date
  end_date date
  status varchar
}

Ref: hris.employees.work_unit_id > hris.work_units.id
Ref: hris.lecturers.employee_id - hris.employees.id
Ref: hris.attendances.employee_id > hris.employees.id
Ref: hris.leave_requests.employee_id > hris.employees.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `hris.employees.person_id` merujuk logis ke `core_db.persons.id`.
* `hris.lecturers.homebase_study_program_id` merujuk logis ke `reference_db.study_programs.id`.
