# ERD Final - CRM Module

## 1. Skema Database
Database `crm_db` menyimpan data prospek pendaftar (*leads*), kampanye pemasaran, dan komisi kemitraan agen rujukan.

## 2. Tabel Utama
* **`crm.campaigns`**: Data kampanye pemasaran aktif.
* **`crm.agents`**: Registrasi agen kemitraan rujukan luar.
* **`crm.referrals`**: Kode referral yang di-generate agen rujukan.
* **`crm.leads`**: Pipeline prospek calon mahasiswa.
* **`crm.commission_records`**: Perhitungan komisi rujukan per kualifikasi leads.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table crm.campaigns {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  channel varchar
  start_date date
  end_date date
  status varchar
}

Table crm.agents {
  id uuid [pk]
  person_id uuid [not null]
  agent_code varchar [unique, not null]
  organization_name varchar
  status varchar
}

Table crm.referrals {
  id uuid [pk]
  referral_type varchar [not null]
  referrer_person_id uuid
  agent_id uuid
  referral_code varchar [unique, not null]
  is_valid boolean
}

Table crm.leads {
  id uuid [pk]
  person_id uuid [not null]
  study_program_id uuid
  lead_source_id uuid
  campaign_id uuid
  referral_id uuid
  lead_number varchar [unique, not null]
  status varchar [not null]
  converted_at timestamp
  created_at timestamp
}

Table crm.commission_records {
  id uuid [pk]
  lead_id uuid [not null]
  commission_rule_id uuid
  referrer_person_id uuid
  amount numeric
  status varchar
}

Ref: crm.referrals.agent_id > crm.agents.id
Ref: crm.leads.campaign_id > crm.campaigns.id
Ref: crm.leads.referral_id > crm.referrals.id
Ref: crm.commission_records.lead_id > crm.leads.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `crm.leads.person_id` secara logis merujuk ke `core_db.persons.id`.
* `crm.leads.study_program_id` merujuk ke `reference_db.study_programs.id`.
