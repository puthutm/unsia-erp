# ERD Final - Core Module

## 1. Skema Database
Database `core_db` bertindak sebagai pusat otentikasi identitas, session management, dan pembatasan hak akses berbasis peran (RBAC).

## 2. Tabel Utama
* **`core.persons`**: Menyimpan data identitas pribadi detail fisik pengguna (nama, email, nomor identitas, tanggal lahir, alamat).
* **`core.users`**: Menyimpan username, email, password hash, dan status aktif akun pengguna.
* **`core.roles`**: Menyimpan daftar peran sistem (misal: `MAHASISWA`, `DOSEN`, `ADMIN`).
* **`core.permissions`**: Menyimpan daftar izin akses (misal: `academic.krs.approve`).
* **`core.user_roles`**: Menghubungkan pengguna dengan peran yang diaktifkan beserta lingkup program studinya.
* **`core.role_permissions`**: Menghubungkan peran dengan izin aksesnya.
* **`core.idempotency_keys`**: Tabel pengaman mutasi ganda pada API/Event.
* **`core.audit_logs`**: Tabel rekaman jejak audit global sensitif.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table core.persons {
  id uuid [pk]
  full_name varchar [not null]
  email varchar
  phone varchar
  identity_number varchar
  gender varchar
  birth_place varchar
  birth_date date
  religion_id uuid
  country_id uuid
  address text
  created_at timestamp
  updated_at timestamp
}

Table core.users {
  id uuid [pk]
  person_id uuid [not null]
  username varchar [unique, not null]
  email varchar [unique, not null]
  password_hash text [not null]
  status varchar [not null] // active, inactive, suspended
  last_login_at timestamp
  created_at timestamp
  updated_at timestamp
}

Table core.roles {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  scope_type varchar // global, prodi, module, self
  is_system boolean [default: false]
  is_active boolean [default: true]
}

Table core.permissions {
  id uuid [pk]
  code varchar [unique, not null] // module.resource.action
  module varchar [not null]
  resource varchar [not null]
  action varchar [not null]
  is_active boolean [default: true]
}

Table core.user_roles {
  id uuid [pk]
  user_id uuid [not null]
  role_id uuid [not null]
  study_program_id uuid
  assigned_at timestamp
}

Table core.role_permissions {
  id uuid [pk]
  role_id uuid [not null]
  permission_id uuid [not null]
  assigned_at timestamp
}

Ref: core.users.person_id > core.persons.id
Ref: core.user_roles.user_id > core.users.id
Ref: core.user_roles.role_id > core.roles.id
Ref: core.role_permissions.role_id > core.roles.id
Ref: core.role_permissions.permission_id > core.permissions.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `core.persons.id` dirujuk oleh tabel profil pendaftar (`pmb_db.applicants.person_ref_id`) dan tabel karyawan (`hris_db.employees.person_ref_id`) secara logis.
