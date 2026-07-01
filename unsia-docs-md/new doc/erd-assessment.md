# ERD Final - Assessment Module

## 1. Skema Database
Database `assessment_db` menyimpan data bank soal, manajemen versi soal ujian, penugasan sesi kuis/CBT, penandaan lembar pengerjaan peserta (*attempts*), dan rincian lembar jawaban peserta.

## 2. Tabel Utama
* **`assessment.question_banks`**: Bank soal/kelompok mata uji.
* **`assessment.questions`**: Soal ujian utama beserta tipe dan bobot kesulitannya.
* **`assessment.question_options`**: Opsi pilihan ganda untuk jawaban soal.
* **`assessment.assessment_sessions`**: Pengaturan sesi ujian (durasi, modul pemanggil).
* **`assessment.assessment_participants`**: Daftar peserta terdaftar ujian.
* **`assessment.assessment_attempts`**: Lembar pengerjaan kuis/ujian per sesi.
* **`assessment.assessment_answers`**: Rekam jawaban peserta per soal per kuis.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table assessment.question_banks {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar
  module_scope varchar
  status varchar
}

Table assessment.questions {
  id uuid [pk]
  question_bank_id uuid [not null]
  question_type varchar
  difficulty varchar
  question_text text
  answer_explanation text
  status varchar
}

Table assessment.question_options {
  id uuid [pk]
  question_id uuid [not null]
  option_label varchar
  option_text text
  is_correct boolean
}

Table assessment.assessment_sessions {
  id uuid [pk]
  question_set_id uuid
  session_type varchar
  context_module varchar
  context_id uuid
  title varchar
  start_at timestamp
  end_at timestamp
  duration_minutes int
  passing_grade numeric
}

Table assessment.assessment_participants {
  id uuid [pk]
  assessment_session_id uuid [not null]
  participant_type varchar
  applicant_id uuid
  student_id uuid
  user_id uuid
  status varchar
}

Table assessment.assessment_attempts {
  id uuid [pk]
  assessment_session_id uuid [not null]
  participant_id uuid [not null]
  attempt_number int [default: 1]
  idempotency_key varchar
  started_at timestamp
  submitted_at timestamp
  status varchar
  total_score numeric
}

Table assessment.assessment_answers {
  id uuid [pk]
  attempt_id uuid [not null]
  question_id uuid [not null]
  selected_option_id uuid
  answer_text text
  score numeric
}

Ref: assessment.questions.question_bank_id > assessment.question_banks.id
Ref: assessment.question_options.question_id > assessment.questions.id
Ref: assessment.assessment_participants.assessment_session_id > assessment.assessment_sessions.id
Ref: assessment.assessment_attempts.assessment_session_id > assessment.assessment_sessions.id
Ref: assessment.assessment_attempts.participant_id > assessment.assessment_participants.id
Ref: assessment.assessment_answers.attempt_id > assessment.assessment_attempts.id
Ref: assessment.assessment_answers.question_id > assessment.questions.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `assessment.assessment_sessions.context_id` secara logis merujuk ke modul pemanggil context, seperti `pmb_db.applicants.id` atau `lms_db.classes.id`.
* `assessment.assessment_participants.student_id` merujuk logis ke `academic_db.students.id`.
