# ERD Final - Finance Module

## 1. Skema Database
Database `finance_db` menyimpan data penerbitan invoice tagihan, mutasi pembayaran masuk, jurnal akuntansi, dan status clearance keuangan mahasiswa.

## 2. Tabel Utama
* **`finance.invoices`**: Data induk tagihan biaya UKT/pendaftaran.
* **`finance.invoice_items`**: Rincian tagihan per komponen biaya.
* **`finance.payments`**: Data transaksi pembayaran yang dilakukan customer.
* **`finance.payment_gateway_callbacks`**: Log callback dari provider payment gateway luar (idempotent guard).
* **`finance.student_clearances`**: Status kelayakan finansial mahasiswa.
* **`finance.clearance_dispensations`**: Pengajuan dispensasi penangguhan blokir KRS.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table finance.invoices {
  id uuid [pk]
  invoice_number varchar [unique, not null]
  target_type varchar [not null]
  applicant_id uuid
  student_id uuid
  academic_period_id uuid
  total_amount numeric
  paid_amount numeric
  status varchar
  due_date date
  created_at timestamp
}

Table finance.invoice_items {
  id uuid [pk]
  invoice_id uuid [not null]
  payment_component_id uuid
  description text
  amount numeric
  discount_amount numeric
  final_amount numeric
}

Table finance.payments {
  id uuid [pk]
  invoice_id uuid [not null]
  payment_method_id uuid
  payment_number varchar [unique]
  amount numeric
  payment_status varchar
  paid_at timestamp
  external_reference varchar
  idempotency_key varchar
}

Table finance.payment_gateway_callbacks {
  id uuid [pk]
  payment_id uuid
  provider varchar [not null]
  provider_event_id varchar
  external_reference varchar
  idempotency_key varchar
  payload jsonb
  signature_valid boolean
  callback_status varchar
}

Table finance.student_clearances {
  id uuid [pk]
  student_id uuid [not null]
  academic_period_id uuid
  service_scope varchar
  status varchar
  reason text
  valid_until date
}

Table finance.clearance_dispensations {
  id uuid [pk]
  student_clearance_id uuid [not null]
  reason text
  approved_by uuid
  approved_at timestamp
  valid_until date
  status varchar
}

Ref: finance.invoice_items.invoice_id > finance.invoices.id
Ref: finance.payments.invoice_id > finance.invoices.id
Ref: finance.payment_gateway_callbacks.payment_id > finance.payments.id
Ref: finance.clearance_dispensations.student_clearance_id > finance.student_clearances.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `finance.invoices.applicant_id` merujuk logis ke `pmb_db.applicants.id`.
* `finance.invoices.student_id` merujuk logis ke `academic_db.students.id`.
* `finance.student_clearances.student_id` dirujuk oleh tabel penegak KRS (`academic_db.student_clearance_snapshots.student_id`) secara asinkron.
