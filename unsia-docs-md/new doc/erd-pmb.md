# ERD Final - PMB Module

## 1. Skema Database
Database `pmb_db` menyimpan data profil pendaftar calon mahasiswa baru, berkas administrasi syarat masuk, hasil seleksi, dan status registrasi daftar ulang.

## 2. Tabel Utama
* **`pmb.applicants`**: Data induk pendaftaran PMB.
* **`pmb.applicant_biodata`**: Informasi biodata rinci pendaftar.
* **`pmb.applicant_documents`**: Unggahan berkas persyaratan pendaftaran.
* **`pmb.re_registrations`**: Pencatatan log daftar ulang calon mahasiswa.
* **`pmb.loa_documents`**: Penerbitan berkas LoA kelulusan.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table pmb.applicants {
  id uuid [pk]
  person_id uuid [not null]
  user_id uuid
  crm_lead_id uuid
  study_program_id uuid
  pmb_wave_id uuid
  admission_path_id uuid
  registration_number varchar [unique, not null]
  status varchar [not null]
  submitted_at timestamp
  accepted_at timestamp
  created_at timestamp
}

Table pmb.applicant_biodata {
  id uuid [pk]
  applicant_id uuid [not null]
  full_name varchar
  email varchar
  phone varchar
  nik varchar
  birth_place varchar
  birth_date date
  gender varchar
}

Table pmb.applicant_documents {
  id uuid [pk]
  applicant_id uuid [not null]
  document_type_id uuid [not null]
  file_path text [not null]
  verification_status varchar
}

Table pmb.re_registrations {
  id uuid [pk]
  applicant_id uuid [unique, not null]
  reregistration_number varchar [unique]
  invoice_ref_id uuid
  payment_status varchar
  status varchar
}

Table pmb.loa_documents {
  id uuid [pk]
  applicant_id uuid [unique, not null]
  loa_number varchar [unique]
  file_path text
}

Ref: pmb.applicant_biodata.applicant_id > pmb.applicants.id
Ref: pmb.applicant_documents.applicant_id > pmb.applicants.id
Ref: pmb.re_registrations.applicant_id > pmb.applicants.id
Ref: pmb.loa_documents.applicant_id > pmb.applicants.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `pmb.applicants.person_id` merujuk logis ke `core_db.persons.id`.
* `pmb.applicants.pmb_wave_id` merujuk logis ke `reference_db.pmb_waves.id`.
* `pmb.re_registrations.invoice_ref_id` merujuk logis ke `finance_db.invoices.id`.
