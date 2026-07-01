# ERD Final - LMS Module

## 1. Skema Database
Database `lms_db` menyimpan data delivery ruang perkuliahan online, forum interaksi, materi ajar digital, serta lembar pengumpulan tugas mandiri mahasiswa.

## 2. Tabel Utama
* **`lms.classes`**: Proyeksi kelas kuliah online.
* **`lms.enrollments`**: Daftar hak akses mahasiswa di dalam kelas belajar online.
* **`lms.sessions`**: Sesi tatap muka virtual atau kuliah mandiri per minggu.
* **`lms.materials`**: Unggahan slide/buku materi perkuliahan.
* **`lms.assignments`**: Draf instruksi tugas dari dosen.
* **`lms.assignment_submissions`**: File lembar jawaban tugas mahasiswa.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table lms.classes {
  id uuid [pk]
  academic_class_id uuid [not null]
  lecturer_id uuid
  status varchar
  synced_at timestamp
}

Table lms.enrollments {
  id uuid [pk]
  lms_class_id uuid [not null]
  student_id uuid [not null]
  enrollment_status varchar
  enrolled_at timestamp
}

Table lms.sessions {
  id uuid [pk]
  lms_class_id uuid [not null]
  session_number int
  title varchar
  session_date date
  status varchar
}

Table lms.materials {
  id uuid [pk]
  session_id uuid [not null]
  title varchar
  content_type varchar
  file_url text
}

Table lms.assignments {
  id uuid [pk]
  session_id uuid [not null]
  title varchar
  instruction text
  due_at timestamp
  status varchar
}

Table lms.assignment_submissions {
  id uuid [pk]
  assignment_id uuid [not null]
  student_id uuid [not null]
  file_url text
  submitted_at timestamp
  score numeric
}

Ref: lms.enrollments.lms_class_id > lms.classes.id
Ref: lms.sessions.lms_class_id > lms.classes.id
Ref: lms.materials.session_id > lms.sessions.id
Ref: lms.assignments.session_id > lms.sessions.id
Ref: lms.assignment_submissions.assignment_id > lms.assignments.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `lms.classes.academic_class_id` merujuk logis ke kelas fisik `academic_db.classes.id`.
* `lms.enrollments.student_id` merujuk logis ke `academic_db.students.id`.
* `lms.classes.lecturer_id` merujuk logis ke `hris_db.lecturers.id`.
