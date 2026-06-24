---
title: "DBML Global ERP UNSIA"
source_file: "DBML_Global_UNSIA_v1_0_1_Event_Contract_Updated.dbml"
format: markdown
---

# DBML Global ERP UNSIA

Dokumen ini memuat definisi DBML baseline untuk ERP Pendidikan / SIAKAD Terintegrasi UNSIA.

## Daftar Table Terdeteksi

- `core.persons`
- `core.users`
- `core.roles`
- `core.permissions`
- `core.user_roles`
- `core.role_permissions`
- `core.applications`
- `core.oauth_clients`
- `core.redirect_uris`
- `core.service_tokens`
- `core.idempotency_keys`
- `core.integration_event_logs`
- `core.sessions`
- `core.active_role_sessions`
- `core.impersonation_sessions`
- `core.audit_logs`
- `ref.countries`
- `ref.provinces`
- `ref.cities`
- `ref.districts`
- `ref.villages`
- `ref.religions`
- `ref.study_programs`
- `ref.academic_years`
- `ref.academic_periods`
- `ref.admission_paths`
- `ref.pmb_waves`
- `ref.lead_sources`
- `ref.document_types`
- `ref.payment_components`
- `ref.payment_methods`
- `ref.employee_types`
- `ref.lecturer_statuses`
- `ref.status_codes`
- `crm.campaigns`
- `crm.agents`
- `crm.referrals`
- `crm.leads`
- `crm.lead_activities`
- `crm.lead_status_histories`
- `crm.commission_rules`
- `crm.commission_records`
- `pmb.applicants`
- `pmb.applicant_biodata`
- `pmb.applicant_addresses`
- `pmb.applicant_education_backgrounds`
- `pmb.applicant_family_members`
- `pmb.applicant_financial_profiles`
- `pmb.applicant_facility_profiles`
- `pmb.applicant_documents`
- `pmb.applicant_status_histories`
- `pmb.re_registrations`
- `pmb.loa_documents`
- `pmb.handover_logs`
- `finance.invoices`
- `finance.invoice_items`
- `finance.payments`
- `finance.payment_gateway_callbacks`
- `finance.payment_verifications`
- `finance.scholarships`
- `finance.installment_requests`
- `finance.clearance_policies`
- `finance.student_clearances`
- `finance.clearance_dispensations`
- `finance.cash_accounts`
- `finance.cash_transactions`
- `finance.payroll_runs`
- `finance.payroll_items`
- `finance.disbursements`
- `finance.tax_records`
- `finance.bpjs_records`
- `finance.coa_accounts`
- `finance.journals`
- `finance.journal_entries`
- `finance.budgets`
- `finance.budget_lines`
- `academic.students`
- `academic.student_advisors`
- `academic.nim_format_configs`
- `academic.nim_sequences`
- `academic.academic_period_study_program_settings`
- `academic.academic_settings`
- `academic.curriculums`
- `academic.courses`
- `academic.curriculum_courses`
- `academic.class_packages`
- `academic.class_package_items`
- `academic.course_offerings`
- `academic.classes`
- `academic.class_lecturers`
- `academic.class_schedules`
- `academic.krs`
- `academic.krs_items`
- `academic.grades`
- `academic.grade_histories`
- `academic.khs`
- `academic.transcripts`
- `academic.academic_letters`
- `academic.graduation_requirements`
- `academic.yudisium_records`
- `academic.alumni`
- `hris.work_units`
- `hris.positions`
- `hris.functional_positions`
- `hris.employees`
- `hris.lecturers`
- `hris.attendances`
- `hris.leave_requests`
- `hris.bkd_records`
- `hris.performance_reviews`
- `hris.certifications`
- `hris.payroll_sources`
- `lms.classes`
- `lms.enrollments`
- `lms.sessions`
- `lms.materials`
- `lms.videos`
- `lms.vicon_links`
- `lms.assignments`
- `lms.assignment_submissions`
- `lms.quiz_activities`
- `lms.discussions`
- `lms.discussion_comments`
- `lms.attendances`
- `lms.learning_progress`
- `lms.grade_syncs`
- `assessment.question_banks`
- `assessment.questions`
- `assessment.question_versions`
- `assessment.question_options`
- `assessment.material_banks`
- `assessment.materials`
- `assessment.question_sets`
- `assessment.question_set_items`
- `assessment.assessment_sessions`
- `assessment.assessment_participants`
- `assessment.assessment_attempts`
- `assessment.assessment_answers`
- `assessment.assessment_scores`
- `assessment.surveys`
- `assessment.survey_questions`
- `assessment.survey_responses`
- `portal.notifications`
- `portal.notification_reads`
- `portal.user_preferences`
- `portal.menu_shortcuts`
- `portal.portal_activity_logs`
- `core.event_contracts`
- `core.event_consumers`
- `core.event_replay_logs`
- `core.outbox_events`
- `core.inbox_events`
- `core.reconciliation_mismatch_logs`
- `ref.outbox_events`
- `ref.inbox_events`
- `ref.idempotency_keys`
- `ref.reconciliation_mismatch_logs`
- `crm.outbox_events`
- `crm.inbox_events`
- `crm.idempotency_keys`
- `crm.reconciliation_mismatch_logs`
- `pmb.outbox_events`
- `pmb.inbox_events`
- `pmb.idempotency_keys`
- `pmb.reconciliation_mismatch_logs`
- `finance.outbox_events`
- `finance.inbox_events`
- `finance.idempotency_keys`
- `finance.reconciliation_mismatch_logs`
- `academic.outbox_events`
- `academic.inbox_events`
- `academic.idempotency_keys`
- `academic.reconciliation_mismatch_logs`
- `hris.outbox_events`
- `hris.inbox_events`
- `hris.idempotency_keys`
- `hris.reconciliation_mismatch_logs`
- `lms.outbox_events`
- `lms.inbox_events`
- `lms.idempotency_keys`
- `lms.reconciliation_mismatch_logs`
- `assessment.outbox_events`
- `assessment.inbox_events`
- `assessment.idempotency_keys`
- `assessment.reconciliation_mismatch_logs`
- `portal.outbox_events`
- `portal.inbox_events`
- `portal.idempotency_keys`
- `portal.reconciliation_mismatch_logs`

## Raw DBML

```dbml
// DBML Global UNSIA v1.0.1 Event Contract Updated
// Dirapikan dari DBML Global UNSIA v6.4 ERD Detailed FULL.
// Catatan: file ini mempertahankan struktur tabel dan relasi baseline. Update v1.0.1 menambahkan standar Event Contract: outbox_events, inbox_events, idempotency_keys per modul, event_contract catalog, retry, DLQ, correlation/causation id, source_event_key, synced_at, dan reconciliation_mismatch_logs.


Project unsia_erp_global_v6_dba_reviewed {
  database_type: 'PostgreSQL'
  Note: 'UNSIA ERP / SIAKAD Terintegrasi v6 - DBA Reviewed. Update v1.0.1 menambahkan kontrol Event Contract: event catalog, transactional outbox, inbox idempotent consumer, retry, DLQ, correlation_id, causation_id, source_event_key, synced_at, dan reconciliation mismatch log untuk mendukung distributed modular database.'
}

/* =========================
   CORE / SSO / RBAC
========================= */

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
  province_id uuid
  city_id uuid
  district_id uuid
  village_id uuid
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
  status varchar [not null, note: 'active, inactive, suspended']
  last_login_at timestamp
  created_at timestamp
  updated_at timestamp
}

Table core.roles {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  scope_type varchar [note: 'global, prodi, module, self']
  is_system boolean [default: false]
  is_active boolean [default: true]
}

Table core.permissions {
  id uuid [pk]
  code varchar [unique, not null, note: 'module.resource.action']
  module varchar [not null]
  resource varchar [not null]
  action varchar [not null]
  is_active boolean [default: true]
}

Table core.user_roles {
  id uuid [pk]
  user_id uuid [not null]
  role_id uuid [not null]
  study_program_id uuid [note: 'nullable; untuk scope admin prodi/kaprodi']
  assigned_at timestamp
  indexes {
    (user_id, role_id, study_program_id) [unique]
  }
}

Table core.role_permissions {
  id uuid [pk]
  role_id uuid [not null]
  permission_id uuid [not null]
  assigned_at timestamp

  indexes {
    (role_id, permission_id) [unique]
    role_id
    permission_id
  }
}

Table core.applications {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  launch_url text [not null]
  sso_protocol varchar
  is_active boolean [default: true]
}

Table core.oauth_clients {
  id uuid [pk]
  application_id uuid [not null]
  client_id varchar [unique, not null]
  client_secret_hash text
  client_type varchar [not null, note: 'confidential, public']
  grant_types jsonb
  scopes jsonb
  is_active boolean [default: true]
  created_at timestamp
  updated_at timestamp

  indexes {
    application_id
    client_id
  }
}

Table core.redirect_uris {
  id uuid [pk]
  oauth_client_id uuid [not null]
  redirect_uri text [not null]
  is_active boolean [default: true]
  created_at timestamp

  indexes {
    oauth_client_id
    (oauth_client_id, redirect_uri) [unique]
  }
}

Table core.service_tokens {
  id uuid [pk]
  application_id uuid [not null]
  token_hash text [not null]
  scopes jsonb
  expired_at timestamp
  revoked_at timestamp
  created_at timestamp

  indexes {
    application_id
    token_hash [unique]
    expired_at
  }
}

Table core.idempotency_keys {
  id uuid [pk]
  module varchar [not null]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar
  request_hash text
  response_json jsonb
  response_payload jsonb
  status varchar [not null, note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    (module, idempotency_key) [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table core.integration_event_logs {
  id uuid [pk]
  source_module varchar [not null]
  target_module varchar [not null]
  event_type varchar [not null]
  event_key varchar [not null]
  idempotency_key varchar
  correlation_id varchar
  payload jsonb
  status varchar [not null, note: 'pending, success, failed, ignored']
  error_message text
  created_at timestamp
  processed_at timestamp

  indexes {
    (source_module, target_module, event_type, event_key) [unique]
    idempotency_key
    correlation_id
    status
    created_at
  }
}


Table core.sessions {
  id uuid [pk]
  user_id uuid [not null]
  token_hash text [not null]
  refresh_token_hash text
  expired_at timestamp [not null]
  revoked_at timestamp
  created_at timestamp

  indexes {
    user_id
    token_hash [unique]
    refresh_token_hash
    expired_at
  }
}

Table core.active_role_sessions {
  id uuid [pk]
  user_id uuid [not null]
  role_id uuid [not null]
  session_id uuid [not null]
  application_id uuid
  activated_at timestamp

  indexes {
    user_id
    role_id
    session_id
    application_id
    (session_id, application_id) [unique]
  }
}

Table core.impersonation_sessions {
  id uuid [pk]
  actor_user_id uuid [not null]
  target_user_id uuid [not null]
  target_role_id uuid [not null]
  application_id uuid
  session_id uuid [not null]
  reason text [not null]
  started_at timestamp
  ended_at timestamp
  expired_at timestamp [not null]
  status varchar [not null]
}

Table core.audit_logs {
  id uuid [pk]
  user_id uuid
  actor_user_id uuid
  target_user_id uuid
  active_role_id uuid
  impersonation_session_id uuid
  application_id uuid
  module varchar [not null]
  action varchar [not null]
  entity_name varchar
  entity_id uuid
  reason text
  old_value jsonb
  new_value jsonb
  request_id varchar
  ip_address varchar
  user_agent text
  created_at timestamp

  indexes {
    user_id
    actor_user_id
    target_user_id
    active_role_id
    impersonation_session_id
    application_id
    (module, entity_name, entity_id)
    request_id
    created_at
  }
}

/* =========================
   REFERENSI
========================= */

Table ref.countries {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.provinces {
  id uuid [pk]
  country_id uuid [not null]
  code varchar [not null]
  name varchar [not null]
  is_active boolean
}

Table ref.cities {
  id uuid [pk]
  province_id uuid [not null]
  code varchar [not null]
  name varchar [not null]
  type varchar
  is_active boolean
}

Table ref.districts {
  id uuid [pk]
  city_id uuid [not null]
  code varchar [not null]
  name varchar [not null]
  is_active boolean
}

Table ref.villages {
  id uuid [pk]
  district_id uuid [not null]
  code varchar [not null]
  name varchar [not null]
  postal_code varchar
  is_active boolean
}

Table ref.religions {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.study_programs {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  degree_level varchar
  faculty_name varchar
  mode varchar
  is_active boolean
}

Table ref.academic_years {
  id uuid [pk]
  code varchar [unique, not null, note: 'Contoh: 2026/2027']
  name varchar [not null, note: 'Contoh: Tahun Ajaran 2026/2027']
  start_year int
  end_year int
  start_date date
  end_date date
  status varchar [note: 'draft, active, closed, archived']
  is_active boolean
}

Table ref.academic_periods {
  id uuid [pk]
  academic_year_id uuid [not null]
  code varchar [unique, not null, note: 'Contoh: 2026/2027-GANJIL']
  name varchar [not null, note: 'Contoh: Ganjil 2026/2027']
  semester_type varchar [note: 'ganjil, genap, pendek']
  start_date date
  end_date date
  class_start_date date
  class_end_date date
  total_meetings int
  min_attendance_percentage numeric
  status varchar [note: 'draft, open, active, closed, archived']
  is_active boolean

  indexes {
    academic_year_id
    semester_type
    status
    (academic_year_id, semester_type) [unique]
  }
}

Table ref.admission_paths {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.pmb_waves {
  id uuid [pk]
  academic_year_id uuid [note: 'Tahun ajaran kalender operasional, bukan tahun kurikulum']
  target_entry_period_id uuid [not null, note: 'Periode akademik target masuk pendaftar']
  admission_path_id uuid
  code varchar [unique, not null]
  name varchar [not null]
  start_date date [note: 'Legacy alias untuk registration_start_at']
  end_date date [note: 'Legacy alias untuk registration_end_at']
  registration_start_at timestamp
  registration_end_at timestamp
  selection_start_at timestamp
  selection_end_at timestamp
  reregistration_deadline_at timestamp
  status varchar [note: 'draft, open, closed, archived']
  is_active boolean

  indexes {
    academic_year_id
    target_entry_period_id
    admission_path_id
    status
  }
}

Table ref.lead_sources {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.document_types {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  module_scope varchar
  is_required boolean
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

Table ref.employee_types {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.lecturer_statuses {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.status_codes {
  id uuid [pk]
  module varchar [not null]
  code varchar [not null]
  name varchar [not null]
  description text
  is_active boolean

  indexes {
    (module, code) [unique]
    module
    code
  }
}

/* =========================
   CRM
========================= */

Table crm.campaigns {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  channel varchar
  start_date date
  end_date date
  status varchar
  created_by uuid
}

Table crm.agents {
  id uuid [pk]
  person_id uuid [not null]
  agent_code varchar [unique, not null]
  organization_name varchar
  status varchar
  approval_status varchar
  approved_by uuid
  approved_at timestamp
}

Table crm.referrals {
  id uuid [pk]
  referral_type varchar [not null]
  referrer_person_id uuid
  agent_id uuid
  referral_code varchar [unique, not null]
  is_valid boolean
  created_at timestamp
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
  owner_user_id uuid
  converted_at timestamp
  created_at timestamp
}

Table crm.lead_activities {
  id uuid [pk]
  lead_id uuid [not null]
  user_id uuid
  activity_type varchar
  note text
  activity_at timestamp
}

Table crm.lead_status_histories {
  id uuid [pk]
  lead_id uuid [not null]
  old_status varchar
  new_status varchar
  changed_by uuid
  note text
  changed_at timestamp
}

Table crm.commission_rules {
  id uuid [pk]
  referral_type varchar
  amount numeric
  calculation_type varchar
  is_active boolean
}

Table crm.commission_records {
  id uuid [pk]
  lead_id uuid [not null]
  commission_rule_id uuid
  referrer_person_id uuid
  amount numeric
  status varchar
  sent_to_finance_at timestamp
}

/* =========================
   PMB
========================= */
Table pmb.applicants {
  id uuid [pk]
  person_id uuid [not null]
  user_id uuid
  crm_lead_id uuid
  study_program_id uuid
  pmb_wave_id uuid
  admission_path_id uuid
  target_entry_period_id uuid [note: 'Snapshot dari ref.pmb_waves.target_entry_period_id']
  registration_number varchar [unique, not null]
  status varchar [not null, note: 'draft, submitted, verified, accepted, reregistration_completed, ready_for_academic']
  submitted_at timestamp
  accepted_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    person_id
    user_id
    crm_lead_id
    study_program_id
    pmb_wave_id
    admission_path_id
    target_entry_period_id
    status
  }
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
  religion_id uuid
  marital_status varchar
  citizenship varchar
  jacket_size varchar
  core_sync_status varchar
  core_synced_at timestamp
  updated_at timestamp
}

Table pmb.applicant_addresses {
  id uuid [pk]
  applicant_id uuid [not null]
  address_type varchar
  street text
  province_id uuid
  city_id uuid
  district_id uuid
  village_id uuid
  postal_code varchar
  is_same_as_ktp boolean
}

Table pmb.applicant_education_backgrounds {
  id uuid [pk]
  applicant_id uuid [not null]
  education_level_id uuid
  institution_name varchar
  npsn_or_pt_code varchar
  nisn_or_previous_nim varchar
  graduation_year int
  average_score numeric
}

Table pmb.applicant_family_members {
  id uuid [pk]
  applicant_id uuid [not null]
  relation varchar
  nik varchar
  full_name varchar
  education_level_id uuid
  occupation varchar
  income_range varchar
  phone varchar
  dependent_count int
}

Table pmb.applicant_financial_profiles {
  id uuid [pk]
  applicant_id uuid [not null]
  personal_income_range varchar
  bank_name varchar
  bank_account_name varchar
  bank_account_number varchar
  scholarship_interest boolean
}

Table pmb.applicant_facility_profiles {
  id uuid [pk]
  applicant_id uuid [not null]
  employment_status varchar
  has_vehicle boolean
  has_pjj_device boolean
  internet_access varchar
  special_need_status varchar
}

Table pmb.applicant_documents {
  id uuid [pk]
  applicant_id uuid [not null]
  document_type_id uuid [not null]
  file_url text
  verification_status varchar
  verification_note text
  verified_by uuid
  verified_at timestamp
  uploaded_at timestamp
}

Table pmb.applicant_status_histories {
  id uuid [pk]
  applicant_id uuid [not null]
  old_status varchar
  new_status varchar
  changed_by uuid
  note text
  changed_at timestamp
}

Table pmb.re_registrations {
  id uuid [pk]
  applicant_id uuid [not null]
  status varchar
  submitted_at timestamp
  verified_at timestamp
  verified_by uuid
}

Table pmb.loa_documents {
  id uuid [pk]
  applicant_id uuid [not null]
  loa_number varchar [unique, not null]
  file_url text
  issued_at timestamp
  issued_by uuid
}

Table pmb.handover_logs {
  id uuid [pk]
  applicant_id uuid [not null]
  target_module varchar [not null, note: 'academic']
  handover_status varchar [not null, note: 'pending, success, failed, ignored']
  idempotency_key varchar [not null]
  correlation_id varchar
  payload jsonb
  response_json jsonb
  error_message text
  handed_over_by uuid
  handed_over_at timestamp

  indexes {
    applicant_id
    idempotency_key [unique]
    correlation_id
    handover_status
  }
}


/* =========================
   FINANCE
========================= */

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
  updated_at timestamp
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
  created_at timestamp

  indexes {
    invoice_id
    payment_method_id
    payment_status
    external_reference
    idempotency_key [unique]
    (invoice_id, payment_status)
  }
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
  callback_status varchar [note: 'received, processed, ignored, failed']
  received_at timestamp
  processed_at timestamp

  indexes {
    payment_id
    provider
    provider_event_id
    external_reference
    idempotency_key [unique]
    (provider, provider_event_id) [unique]
  }
}

Table finance.payment_verifications {
  id uuid [pk]
  payment_id uuid [not null]
  verified_by uuid
  verification_status varchar
  rejection_reason varchar
  note text
  verified_at timestamp
}

Table finance.scholarships {
  id uuid [pk]
  student_id uuid [not null]
  scholarship_type varchar
  amount numeric
  status varchar
  approved_by uuid
  approved_at timestamp
}

Table finance.installment_requests {
  id uuid [pk]
  invoice_id uuid [not null]
  student_id uuid
  status varchar
  reason text
  requested_at timestamp
  approved_by uuid
  approved_at timestamp
}

Table finance.clearance_policies {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar
  service_scope varchar
  rule_json jsonb
  is_active boolean
}

Table finance.student_clearances {
  id uuid [pk]
  student_id uuid [not null]
  academic_period_id uuid
  service_scope varchar
  status varchar
  reason text
  valid_until date
  updated_by uuid
  updated_at timestamp
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

Table finance.cash_accounts {
  id uuid [pk]
  account_code varchar [unique, not null]
  account_name varchar
  bank_name varchar
  account_number varchar
  is_active boolean
}

Table finance.cash_transactions {
  id uuid [pk]
  cash_account_id uuid [not null]
  transaction_type varchar
  source_type varchar
  source_id uuid
  amount numeric
  description text
  transaction_at timestamp
}

Table finance.payroll_runs {
  id uuid [pk]
  payroll_period varchar
  run_date date
  total_amount numeric
  status varchar
  approved_by uuid
}

Table finance.payroll_items {
  id uuid [pk]
  payroll_run_id uuid [not null]
  employee_id uuid
  gross_amount numeric
  deduction_amount numeric
  net_amount numeric
  status varchar
}

Table finance.disbursements {
  id uuid [pk]
  disbursement_type varchar
  commission_record_id uuid
  recipient_person_id uuid
  amount numeric
  status varchar
  disbursed_at timestamp
}

Table finance.tax_records {
  id uuid [pk]
  tax_type varchar
  source_type varchar
  source_id uuid
  amount numeric
  status varchar
  tax_period date
}

Table finance.bpjs_records {
  id uuid [pk]
  employee_id uuid
  amount numeric
  period varchar
  status varchar
}

Table finance.coa_accounts {
  id uuid [pk]
  account_code varchar [unique, not null]
  account_name varchar
  normal_balance varchar
  is_active boolean
}

Table finance.journals {
  id uuid [pk]
  journal_number varchar [unique, not null]
  journal_date date
  source_type varchar
  source_id uuid
  description text
  created_by uuid
}

Table finance.journal_entries {
  id uuid [pk]
  journal_id uuid [not null]
  coa_account_id uuid
  debit numeric
  credit numeric
  description text
}

Table finance.budgets {
  id uuid [pk]
  budget_code varchar [unique, not null]
  name varchar
  fiscal_year varchar
  total_amount numeric
  status varchar
}

Table finance.budget_lines {
  id uuid [pk]
  budget_id uuid [not null]
  coa_account_id uuid
  description text
  amount numeric
  realized_amount numeric
}

/* =========================
   ACADEMIC
========================= */
Table academic.students {
  id uuid [pk]
  person_id uuid [not null]
  user_id uuid
  applicant_id uuid
  study_program_id uuid [not null]
  nim varchar [unique, not null]
  student_status varchar
  entry_academic_year_id uuid
  entry_period_id uuid
  curriculum_id uuid
  current_semester int
  active_date date
  created_at timestamp
  updated_at timestamp

  indexes {
    person_id
    user_id
    applicant_id [unique]
    study_program_id
    entry_academic_year_id
    entry_period_id
    curriculum_id
    student_status
    (study_program_id, student_status)
  }
}

Table academic.student_advisors {
  id uuid [pk]
  student_id uuid [not null]
  lecturer_id uuid [not null]
  academic_period_id uuid
  is_active boolean
  assigned_at timestamp

  indexes {
    student_id
    lecturer_id
    academic_period_id
    (student_id, academic_period_id) [unique]
  }
}

Table academic.nim_format_configs {
  id uuid [pk]
  code varchar [unique, not null]
  format_template text
  token_order jsonb
  is_active boolean
  created_by uuid
}

Table academic.nim_sequences {
  id uuid [pk]
  study_program_id uuid [not null]
  entry_period_id uuid [not null]
  sequence_year varchar [not null]
  last_number int [not null, default: 0]
  updated_at timestamp

  indexes {
    (study_program_id, entry_period_id, sequence_year) [unique]
    study_program_id
    entry_period_id
  }
}


Table academic.academic_period_study_program_settings {
  id uuid [pk]
  academic_period_id uuid [not null]
  study_program_id uuid [not null]
  class_start_date date
  class_end_date date
  total_meetings int
  min_attendance_percentage numeric
  pddikti_start_date date
  pddikti_end_date date
  is_active boolean
  created_at timestamp
  updated_at timestamp

  indexes {
    academic_period_id
    study_program_id
    (academic_period_id, study_program_id) [unique]
  }
}

Table academic.academic_settings {
  id uuid [pk]
  setting_key varchar [unique, not null]
  setting_value jsonb
  updated_by uuid
  updated_at timestamp
}

Table academic.curriculums {
  id uuid [pk]
  study_program_id uuid [not null]
  code varchar [unique, not null, note: 'Contoh: SI-KUR-2021']
  name varchar [not null, note: 'Contoh: Kurikulum Sistem Informasi 2021']
  curriculum_year int [not null, note: 'Tahun Kurikulum, bukan Tahun Ajaran']
  year int [note: 'Deprecated alias. Gunakan curriculum_year.']
  status varchar [note: 'draft, active, inactive, archived']
  effective_at timestamp [note: 'Legacy alias. Gunakan effective_start_period_id.']
  effective_start_period_id uuid
  effective_end_period_id uuid
  is_active boolean
  is_default_for_new_student boolean
  created_at timestamp
  updated_at timestamp

  indexes {
    study_program_id
    curriculum_year
    status
    effective_start_period_id
    effective_end_period_id
    is_default_for_new_student
    (study_program_id, curriculum_year) [unique]
  }
}

Table academic.courses {
  id uuid [pk]
  study_program_id uuid
  course_code varchar [unique, not null]
  course_name varchar [not null]
  sks int
  course_type varchar
  minimum_grade numeric
  is_active boolean
}

Table academic.curriculum_courses {
  id uuid [pk]
  curriculum_id uuid [not null]
  course_id uuid [not null]
  semester int
  is_mandatory boolean

  indexes {
    curriculum_id
    course_id
    (curriculum_id, course_id) [unique]
  }
}

Table academic.class_packages {
  id uuid [pk]
  study_program_id uuid [not null]
  curriculum_id uuid [not null]
  semester int
  package_name varchar
  status varchar
}

Table academic.class_package_items {
  id uuid [pk]
  class_package_id uuid [not null]
  course_id uuid [not null]
  recommended_class_id uuid
}

Table academic.course_offerings {
  id uuid [pk]
  course_id uuid [not null]
  academic_period_id uuid [not null]
  curriculum_id uuid [note: 'Kurikulum Prodi yang menjadi konteks penawaran, bila berlaku']
  status varchar
  opened_at timestamp

  indexes {
    course_id
    academic_period_id
    curriculum_id
    status
    (course_id, academic_period_id, curriculum_id) [unique]
  }
}

Table academic.classes {
  id uuid [pk]
  course_offering_id uuid [not null]
  class_code varchar [not null]
  quota int
  enrolled_count int
  class_status varchar
  is_parallel boolean
  created_at timestamp

  indexes {
    course_offering_id
    class_code
    (course_offering_id, class_code) [unique]
  }
}

Table academic.class_lecturers {
  id uuid [pk]
  class_id uuid [not null]
  lecturer_id uuid [not null]
  role_type varchar [note: 'coordinator, teacher, assistant']

  indexes {
    class_id
    lecturer_id
    (class_id, lecturer_id, role_type) [unique]
  }
}

Table academic.class_schedules {
  id uuid [pk]
  class_id uuid [not null]
  day varchar
  start_time time
  end_time time
  room_or_link text
  session_type varchar
}

Table academic.krs {
  id uuid [pk]
  student_id uuid [not null]
  academic_period_id uuid [not null]
  status varchar
  advisor_id uuid
  is_package boolean
  finance_clearance_id uuid
  submitted_at timestamp
  approved_at timestamp

  indexes {
    student_id
    academic_period_id
    advisor_id
    finance_clearance_id
    status
    (student_id, academic_period_id) [unique]
  }
}

Table academic.krs_items {
  id uuid [pk]
  krs_id uuid [not null]
  class_id uuid [not null]
  status varchar
  selected_at timestamp

  indexes {
    krs_id
    class_id
    status
    (krs_id, class_id) [unique]
  }
}

Table academic.grades {
  id uuid [pk]
  krs_item_id uuid [not null]
  numeric_grade numeric
  letter_grade varchar
  grade_point numeric
  source varchar
  submitted_at timestamp
  submitted_by uuid

  indexes {
    krs_item_id [unique]
    submitted_by
    source
  }
}

Table academic.grade_histories {
  id uuid [pk]
  grade_id uuid [not null]
  old_value jsonb
  new_value jsonb
  changed_by uuid
  reason text
  changed_at timestamp
}

Table academic.khs {
  id uuid [pk]
  student_id uuid [not null]
  academic_period_id uuid [not null]
  ips numeric
  total_sks int
  file_url text
  issued_at timestamp
}

Table academic.transcripts {
  id uuid [pk]
  student_id uuid [not null]
  ipk numeric
  total_sks int
  file_url text
  issued_at timestamp
}

Table academic.academic_letters {
  id uuid [pk]
  student_id uuid [not null]
  letter_type varchar
  status varchar
  file_url text
  requested_at timestamp
  issued_at timestamp
}

Table academic.graduation_requirements {
  id uuid [pk]
  study_program_id uuid
  degree_level varchar
  minimum_sks int
  minimum_gpa numeric
  requirement_json jsonb
  is_active boolean
}

Table academic.yudisium_records {
  id uuid [pk]
  student_id uuid [not null]
  yudisium_date date
  graduation_status varchar
  final_gpa numeric
  transcript_number varchar
}

Table academic.alumni {
  id uuid [pk]
  student_id uuid [not null]
  person_id uuid [not null]
  graduation_date date
  alumni_number varchar
  status varchar
  created_at timestamp
}

/* =========================
   HRIS / SDM
========================= */

Table hris.work_units {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  parent_unit_id uuid
  is_active boolean
}

Table hris.positions {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  level varchar
  is_active boolean
}

Table hris.functional_positions {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  rank varchar
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
  end_date date
  is_active boolean

  indexes {
    person_id [unique]
    employee_type_id
    work_unit_id
    position_id
    employment_status
    is_active
  }
}

Table hris.lecturers {
  id uuid [pk]
  employee_id uuid [not null]
  lecturer_status_id uuid
  functional_position_id uuid
  nidn varchar
  homebase_study_program_id uuid
  certification_status varchar
  is_active boolean

  indexes {
    employee_id [unique]
    nidn [unique]
    lecturer_status_id
    functional_position_id
    homebase_study_program_id
    is_active
  }
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
  approved_by uuid
}

Table hris.bkd_records {
  id uuid [pk]
  lecturer_id uuid [not null]
  academic_period_id uuid
  teaching_load numeric
  research_load numeric
  service_load numeric
  status varchar
}

Table hris.performance_reviews {
  id uuid [pk]
  employee_id uuid [not null]
  review_period varchar
  score numeric
  status varchar
  reviewed_by uuid
}

Table hris.certifications {
  id uuid [pk]
  employee_id uuid [not null]
  certification_name varchar
  issuer varchar
  issued_date date
  expired_date date
  file_url text
}

Table hris.payroll_sources {
  id uuid [pk]
  employee_id uuid [not null]
  payroll_period varchar
  base_salary numeric
  allowance_amount numeric
  deduction_amount numeric
  status varchar
}

/* =========================
   LMS / ICEMS
========================= */

Table lms.classes {
  id uuid [pk]
  academic_class_id uuid [not null]
  lecturer_id uuid
  status varchar
  synced_at timestamp

  indexes {
    academic_class_id [unique]
    lecturer_id
    status
  }
}

Table lms.enrollments {
  id uuid [pk]
  lms_class_id uuid [not null]
  student_id uuid [not null]
  enrollment_status varchar
  enrolled_at timestamp

  indexes {
    lms_class_id
    student_id
    enrollment_status
    (lms_class_id, student_id) [unique]
  }
}

Table lms.sessions {
  id uuid [pk]
  lms_class_id uuid [not null]
  session_number int
  title varchar
  session_date date
  start_time time
  end_time time
  status varchar

  indexes {
    lms_class_id
    session_date
    status
    (lms_class_id, session_number) [unique]
  }
}

Table lms.materials {
  id uuid [pk]
  session_id uuid [not null]
  assessment_material_id uuid
  title varchar
  content_type varchar
  file_url text
  published_at timestamp
}

Table lms.videos {
  id uuid [pk]
  session_id uuid [not null]
  title varchar
  video_url text
  duration_minutes int
}

Table lms.vicon_links {
  id uuid [pk]
  session_id uuid [not null]
  provider varchar
  join_url text
  start_at timestamp
  end_at timestamp
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
  graded_by uuid
  graded_at timestamp

  indexes {
    assignment_id
    student_id
    graded_by
    (assignment_id, student_id) [unique]
  }
}

Table lms.quiz_activities {
  id uuid [pk]
  session_id uuid [not null]
  assessment_session_id uuid [not null]
  title varchar
  status varchar
}

Table lms.discussions {
  id uuid [pk]
  session_id uuid [not null]
  title varchar
  created_by uuid
  created_at timestamp
}

Table lms.discussion_comments {
  id uuid [pk]
  discussion_id uuid [not null]
  user_id uuid
  content text
  parent_comment_id uuid
  created_at timestamp
}

Table lms.attendances {
  id uuid [pk]
  session_id uuid [not null]
  student_id uuid [not null]
  attendance_status varchar
  submitted_at timestamp

  indexes {
    session_id
    student_id
    attendance_status
    (session_id, student_id) [unique]
  }
}

Table lms.learning_progress {
  id uuid [pk]
  enrollment_id uuid [not null]
  progress_percent numeric
  last_access_at timestamp
  completed_at timestamp

  indexes {
    enrollment_id [unique]
    last_access_at
  }
}

Table lms.grade_syncs {
  id uuid [pk]
  lms_class_id uuid [not null]
  academic_class_id uuid [not null]
  sync_status varchar
  synced_at timestamp
  payload jsonb
}

/* =========================
   ASSESSMENT
========================= */

Table assessment.question_banks {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar
  module_scope varchar
  owner_user_id uuid
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
  created_by uuid
}

Table assessment.question_versions {
  id uuid [pk]
  question_id uuid [not null]
  version_number int [not null]
  question_type varchar
  difficulty varchar
  question_text text
  answer_explanation text
  options_snapshot jsonb
  status varchar [note: 'draft, approved, archived']
  created_by uuid
  created_at timestamp

  indexes {
    question_id
    created_by
    (question_id, version_number) [unique]
  }
}

Table assessment.question_options {
  id uuid [pk]
  question_id uuid [not null]
  option_label varchar
  option_text text
  is_correct boolean
  sort_order int
}

Table assessment.material_banks {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar
  module_scope varchar
  owner_user_id uuid
  status varchar
}

Table assessment.materials {
  id uuid [pk]
  material_bank_id uuid [not null]
  title varchar
  material_type varchar
  file_url text
  content_text text
  status varchar
  created_by uuid
}

Table assessment.question_sets {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar
  randomize_questions boolean
  randomize_options boolean
  status varchar
}

Table assessment.question_set_items {
  id uuid [pk]
  question_set_id uuid [not null]
  question_id uuid [not null]
  score_weight numeric
  sort_order int

  indexes {
    question_set_id
    question_id
    (question_set_id, question_id) [unique]
  }
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
  status varchar
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

  indexes {
    assessment_session_id
    participant_id
    idempotency_key [unique]
    (assessment_session_id, participant_id, attempt_number) [unique]
  }
}

Table assessment.assessment_answers {
  id uuid [pk]
  attempt_id uuid [not null]
  question_id uuid [not null]
  selected_option_id uuid
  answer_text text
  score numeric
  graded_by uuid
  graded_at timestamp

  indexes {
    attempt_id
    question_id
    selected_option_id
    graded_by
    (attempt_id, question_id) [unique]
  }
}

Table assessment.assessment_scores {
  id uuid [pk]
  attempt_id uuid [not null]
  total_score numeric
  result_status varchar
  published_at timestamp
  sent_to_context_at timestamp
}

Table assessment.surveys {
  id uuid [pk]
  title varchar
  target_type varchar
  is_anonymous boolean
  start_at timestamp
  end_at timestamp
  status varchar
  created_by uuid
}

Table assessment.survey_questions {
  id uuid [pk]
  survey_id uuid [not null]
  question_type varchar
  question_text text
  sort_order int
}

Table assessment.survey_responses {
  id uuid [pk]
  survey_id uuid [not null]
  respondent_user_id uuid
  submitted_at timestamp
  response_json jsonb
}

/* =========================
   PORTAL
========================= */

Table portal.notifications {
  id uuid [pk]
  user_id uuid [not null]
  title varchar
  message text
  module_source varchar
  target_url text
  sent_at timestamp
}

Table portal.notification_reads {
  id uuid [pk]
  notification_id uuid [not null]
  user_id uuid [not null]
  read_at timestamp

  indexes {
    notification_id
    user_id
    (notification_id, user_id) [unique]
  }
}

Table portal.user_preferences {
  id uuid [pk]
  user_id uuid [not null]
  preference_key varchar
  preference_value jsonb
  updated_at timestamp

  indexes {
    user_id
    (user_id, preference_key) [unique]
  }
}

Table portal.menu_shortcuts {
  id uuid [pk]
  user_id uuid [not null]
  menu_code varchar
  menu_label varchar
  target_url text
  sort_order int

  indexes {
    user_id
    (user_id, menu_code) [unique]
  }
}

Table portal.portal_activity_logs {
  id uuid [pk]
  user_id uuid [not null]
  activity_type varchar
  module_target varchar
  description text
  created_at timestamp
}

/* =========================
   RELATIONS
========================= */

Ref: core.users.person_id > core.persons.id
Ref: core.user_roles.user_id > core.users.id
Ref: core.user_roles.role_id > core.roles.id
Ref: core.user_roles.study_program_id > ref.study_programs.id
Ref: core.role_permissions.role_id > core.roles.id
Ref: core.role_permissions.permission_id > core.permissions.id
Ref: core.sessions.user_id > core.users.id
Ref: core.active_role_sessions.user_id > core.users.id
Ref: core.active_role_sessions.role_id > core.roles.id
Ref: core.active_role_sessions.session_id > core.sessions.id
Ref: core.active_role_sessions.application_id > core.applications.id
Ref: core.impersonation_sessions.actor_user_id > core.users.id
Ref: core.impersonation_sessions.target_user_id > core.users.id
Ref: core.impersonation_sessions.target_role_id > core.roles.id
Ref: core.impersonation_sessions.application_id > core.applications.id
Ref: core.audit_logs.user_id > core.users.id

Ref: core.oauth_clients.application_id > core.applications.id
Ref: core.redirect_uris.oauth_client_id > core.oauth_clients.id
Ref: core.service_tokens.application_id > core.applications.id
Ref: core.audit_logs.actor_user_id > core.users.id
Ref: core.audit_logs.target_user_id > core.users.id
Ref: core.audit_logs.active_role_id > core.roles.id
Ref: core.audit_logs.impersonation_session_id > core.impersonation_sessions.id
Ref: core.audit_logs.application_id > core.applications.id

Ref: core.persons.religion_id > ref.religions.id
Ref: core.persons.country_id > ref.countries.id
Ref: core.persons.province_id > ref.provinces.id
Ref: core.persons.city_id > ref.cities.id
Ref: core.persons.district_id > ref.districts.id
Ref: core.persons.village_id > ref.villages.id
Ref: ref.provinces.country_id > ref.countries.id
Ref: ref.cities.province_id > ref.provinces.id
Ref: ref.districts.city_id > ref.cities.id
Ref: ref.villages.district_id > ref.districts.id
Ref: ref.academic_periods.academic_year_id > ref.academic_years.id

Ref: crm.campaigns.created_by > core.users.id
Ref: crm.agents.person_id > core.persons.id
Ref: crm.agents.approved_by > core.users.id
Ref: crm.referrals.referrer_person_id > core.persons.id
Ref: crm.referrals.agent_id > crm.agents.id
Ref: crm.leads.person_id > core.persons.id
Ref: crm.leads.study_program_id > ref.study_programs.id
Ref: crm.leads.lead_source_id > ref.lead_sources.id
Ref: crm.leads.campaign_id > crm.campaigns.id
Ref: crm.leads.referral_id > crm.referrals.id
Ref: crm.leads.owner_user_id > core.users.id
Ref: crm.lead_activities.lead_id > crm.leads.id
Ref: crm.lead_activities.user_id > core.users.id
Ref: crm.lead_status_histories.lead_id > crm.leads.id
Ref: crm.lead_status_histories.changed_by > core.users.id
Ref: crm.commission_records.lead_id > crm.leads.id
Ref: crm.commission_records.commission_rule_id > crm.commission_rules.id
Ref: crm.commission_records.referrer_person_id > core.persons.id

Ref: pmb.applicants.person_id > core.persons.id
Ref: pmb.applicants.user_id > core.users.id
Ref: pmb.applicants.crm_lead_id > crm.leads.id
Ref: pmb.applicants.study_program_id > ref.study_programs.id
Ref: pmb.applicants.pmb_wave_id > ref.pmb_waves.id
Ref: pmb.applicants.admission_path_id > ref.admission_paths.id
Ref: pmb.applicant_biodata.applicant_id > pmb.applicants.id
Ref: pmb.applicant_biodata.religion_id > ref.religions.id
Ref: pmb.applicant_addresses.applicant_id > pmb.applicants.id
Ref: pmb.applicant_addresses.province_id > ref.provinces.id
Ref: pmb.applicant_addresses.city_id > ref.cities.id
Ref: pmb.applicant_addresses.district_id > ref.districts.id
Ref: pmb.applicant_addresses.village_id > ref.villages.id
Ref: pmb.applicant_education_backgrounds.applicant_id > pmb.applicants.id
Ref: pmb.applicant_education_backgrounds.education_level_id > ref.status_codes.id
Ref: pmb.applicant_family_members.applicant_id > pmb.applicants.id
Ref: pmb.applicant_financial_profiles.applicant_id > pmb.applicants.id
Ref: pmb.applicant_facility_profiles.applicant_id > pmb.applicants.id
Ref: pmb.applicant_documents.applicant_id > pmb.applicants.id
Ref: pmb.applicant_documents.document_type_id > ref.document_types.id
Ref: pmb.applicant_documents.verified_by > core.users.id
Ref: pmb.applicant_status_histories.applicant_id > pmb.applicants.id
Ref: pmb.applicant_status_histories.changed_by > core.users.id
Ref: pmb.re_registrations.applicant_id > pmb.applicants.id
Ref: pmb.re_registrations.verified_by > core.users.id
Ref: pmb.loa_documents.applicant_id > pmb.applicants.id
Ref: pmb.loa_documents.issued_by > core.users.id
Ref: pmb.handover_logs.applicant_id > pmb.applicants.id
Ref: pmb.handover_logs.handed_over_by > core.users.id

Ref: finance.invoices.applicant_id > pmb.applicants.id
Ref: finance.invoices.student_id > academic.students.id
Ref: finance.invoices.academic_period_id > ref.academic_periods.id
Ref: finance.invoice_items.invoice_id > finance.invoices.id
Ref: finance.invoice_items.payment_component_id > ref.payment_components.id
Ref: finance.payments.invoice_id > finance.invoices.id
Ref: finance.payments.payment_method_id > ref.payment_methods.id
Ref: finance.payment_gateway_callbacks.payment_id > finance.payments.id
Ref: finance.payment_verifications.payment_id > finance.payments.id
Ref: finance.payment_verifications.verified_by > core.users.id
Ref: finance.scholarships.student_id > academic.students.id
Ref: finance.scholarships.approved_by > core.users.id
Ref: finance.installment_requests.invoice_id > finance.invoices.id
Ref: finance.installment_requests.student_id > academic.students.id
Ref: finance.installment_requests.approved_by > core.users.id
Ref: finance.student_clearances.student_id > academic.students.id
Ref: finance.student_clearances.academic_period_id > ref.academic_periods.id
Ref: finance.student_clearances.updated_by > core.users.id
Ref: finance.clearance_dispensations.student_clearance_id > finance.student_clearances.id
Ref: finance.clearance_dispensations.approved_by > core.users.id
Ref: finance.cash_transactions.cash_account_id > finance.cash_accounts.id
Ref: finance.payroll_items.payroll_run_id > finance.payroll_runs.id
Ref: finance.payroll_items.employee_id > hris.employees.id
Ref: finance.disbursements.commission_record_id > crm.commission_records.id
Ref: finance.disbursements.recipient_person_id > core.persons.id
Ref: finance.bpjs_records.employee_id > hris.employees.id
Ref: finance.journal_entries.journal_id > finance.journals.id
Ref: finance.journal_entries.coa_account_id > finance.coa_accounts.id
Ref: finance.budget_lines.budget_id > finance.budgets.id
Ref: finance.budget_lines.coa_account_id > finance.coa_accounts.id

Ref: academic.students.person_id > core.persons.id
Ref: academic.students.user_id > core.users.id
Ref: academic.students.applicant_id > pmb.applicants.id
Ref: academic.students.study_program_id > ref.study_programs.id
Ref: academic.students.entry_period_id > ref.academic_periods.id
Ref: academic.student_advisors.student_id > academic.students.id
Ref: academic.student_advisors.lecturer_id > hris.lecturers.id
Ref: academic.student_advisors.academic_period_id > ref.academic_periods.id
Ref: academic.nim_format_configs.created_by > core.users.id
Ref: academic.nim_sequences.study_program_id > ref.study_programs.id
Ref: academic.nim_sequences.entry_period_id > ref.academic_periods.id
Ref: academic.academic_settings.updated_by > core.users.id
Ref: academic.curriculums.study_program_id > ref.study_programs.id
Ref: academic.courses.study_program_id > ref.study_programs.id
Ref: academic.curriculum_courses.curriculum_id > academic.curriculums.id
Ref: academic.curriculum_courses.course_id > academic.courses.id
Ref: academic.class_packages.study_program_id > ref.study_programs.id
Ref: academic.class_packages.curriculum_id > academic.curriculums.id
Ref: academic.class_package_items.class_package_id > academic.class_packages.id
Ref: academic.class_package_items.course_id > academic.courses.id
Ref: academic.class_package_items.recommended_class_id > academic.classes.id
Ref: academic.course_offerings.course_id > academic.courses.id
Ref: academic.course_offerings.academic_period_id > ref.academic_periods.id
Ref: academic.classes.course_offering_id > academic.course_offerings.id
Ref: academic.class_lecturers.class_id > academic.classes.id
Ref: academic.class_lecturers.lecturer_id > hris.lecturers.id
Ref: academic.class_schedules.class_id > academic.classes.id
Ref: academic.krs.student_id > academic.students.id
Ref: academic.krs.academic_period_id > ref.academic_periods.id
Ref: academic.krs.advisor_id > hris.lecturers.id
Ref: academic.krs.finance_clearance_id > finance.student_clearances.id
Ref: academic.krs_items.krs_id > academic.krs.id
Ref: academic.krs_items.class_id > academic.classes.id
Ref: academic.grades.krs_item_id > academic.krs_items.id
Ref: academic.grades.submitted_by > core.users.id
Ref: academic.grade_histories.grade_id > academic.grades.id
Ref: academic.grade_histories.changed_by > core.users.id
Ref: academic.khs.student_id > academic.students.id
Ref: academic.khs.academic_period_id > ref.academic_periods.id
Ref: academic.transcripts.student_id > academic.students.id
Ref: academic.academic_letters.student_id > academic.students.id
Ref: academic.graduation_requirements.study_program_id > ref.study_programs.id
Ref: academic.yudisium_records.student_id > academic.students.id
Ref: academic.alumni.student_id > academic.students.id
Ref: academic.alumni.person_id > core.persons.id

Ref: hris.work_units.parent_unit_id > hris.work_units.id
Ref: hris.employees.person_id > core.persons.id
Ref: hris.employees.employee_type_id > ref.employee_types.id
Ref: hris.employees.work_unit_id > hris.work_units.id
Ref: hris.employees.position_id > hris.positions.id
Ref: hris.lecturers.employee_id > hris.employees.id
Ref: hris.lecturers.lecturer_status_id > ref.lecturer_statuses.id
Ref: hris.lecturers.functional_position_id > hris.functional_positions.id
Ref: hris.lecturers.homebase_study_program_id > ref.study_programs.id
Ref: hris.attendances.employee_id > hris.employees.id
Ref: hris.leave_requests.employee_id > hris.employees.id
Ref: hris.leave_requests.approved_by > core.users.id
Ref: hris.bkd_records.lecturer_id > hris.lecturers.id
Ref: hris.bkd_records.academic_period_id > ref.academic_periods.id
Ref: hris.performance_reviews.employee_id > hris.employees.id
Ref: hris.performance_reviews.reviewed_by > core.users.id
Ref: hris.certifications.employee_id > hris.employees.id
Ref: hris.payroll_sources.employee_id > hris.employees.id

Ref: lms.classes.academic_class_id > academic.classes.id
Ref: lms.classes.lecturer_id > hris.lecturers.id
Ref: lms.enrollments.lms_class_id > lms.classes.id
Ref: lms.enrollments.student_id > academic.students.id
Ref: lms.sessions.lms_class_id > lms.classes.id
Ref: lms.materials.session_id > lms.sessions.id
Ref: lms.materials.assessment_material_id > assessment.materials.id
Ref: lms.videos.session_id > lms.sessions.id
Ref: lms.vicon_links.session_id > lms.sessions.id
Ref: lms.assignments.session_id > lms.sessions.id
Ref: lms.assignment_submissions.assignment_id > lms.assignments.id
Ref: lms.assignment_submissions.student_id > academic.students.id
Ref: lms.assignment_submissions.graded_by > core.users.id
Ref: lms.quiz_activities.session_id > lms.sessions.id
Ref: lms.quiz_activities.assessment_session_id > assessment.assessment_sessions.id
Ref: lms.discussions.session_id > lms.sessions.id
Ref: lms.discussions.created_by > core.users.id
Ref: lms.discussion_comments.discussion_id > lms.discussions.id
Ref: lms.discussion_comments.user_id > core.users.id
Ref: lms.discussion_comments.parent_comment_id > lms.discussion_comments.id
Ref: lms.attendances.session_id > lms.sessions.id
Ref: lms.attendances.student_id > academic.students.id
Ref: lms.learning_progress.enrollment_id > lms.enrollments.id
Ref: lms.grade_syncs.lms_class_id > lms.classes.id
Ref: lms.grade_syncs.academic_class_id > academic.classes.id

Ref: assessment.question_banks.owner_user_id > core.users.id
Ref: assessment.questions.question_bank_id > assessment.question_banks.id
Ref: assessment.questions.created_by > core.users.id
Ref: assessment.question_versions.question_id > assessment.questions.id
Ref: assessment.question_versions.created_by > core.users.id
Ref: assessment.question_options.question_id > assessment.questions.id
Ref: assessment.material_banks.owner_user_id > core.users.id
Ref: assessment.materials.material_bank_id > assessment.material_banks.id
Ref: assessment.materials.created_by > core.users.id
Ref: assessment.question_set_items.question_set_id > assessment.question_sets.id
Ref: assessment.question_set_items.question_id > assessment.questions.id
Ref: assessment.assessment_sessions.question_set_id > assessment.question_sets.id
Ref: assessment.assessment_participants.assessment_session_id > assessment.assessment_sessions.id
Ref: assessment.assessment_participants.applicant_id > pmb.applicants.id
Ref: assessment.assessment_participants.student_id > academic.students.id
Ref: assessment.assessment_participants.user_id > core.users.id
Ref: assessment.assessment_attempts.assessment_session_id > assessment.assessment_sessions.id
Ref: assessment.assessment_attempts.participant_id > assessment.assessment_participants.id
Ref: assessment.assessment_answers.attempt_id > assessment.assessment_attempts.id
Ref: assessment.assessment_answers.question_id > assessment.questions.id
Ref: assessment.assessment_answers.selected_option_id > assessment.question_options.id
Ref: assessment.assessment_answers.graded_by > core.users.id
Ref: assessment.assessment_scores.attempt_id > assessment.assessment_attempts.id
Ref: assessment.surveys.created_by > core.users.id
Ref: assessment.survey_questions.survey_id > assessment.surveys.id
Ref: assessment.survey_responses.survey_id > assessment.surveys.id
Ref: assessment.survey_responses.respondent_user_id > core.users.id

Ref: portal.notifications.user_id > core.users.id
Ref: portal.notification_reads.notification_id > portal.notifications.id
Ref: portal.notification_reads.user_id > core.users.id
Ref: portal.user_preferences.user_id > core.users.id
Ref: portal.menu_shortcuts.user_id > core.users.id
Ref: portal.portal_activity_logs.user_id > core.users.id

// v6.1 Academic calendar, PMB entry period, and curriculum-year separation
Ref: ref.pmb_waves.academic_year_id > ref.academic_years.id
Ref: ref.pmb_waves.target_entry_period_id > ref.academic_periods.id
Ref: ref.pmb_waves.admission_path_id > ref.admission_paths.id
Ref: pmb.applicants.target_entry_period_id > ref.academic_periods.id
Ref: academic.students.entry_academic_year_id > ref.academic_years.id
Ref: academic.students.curriculum_id > academic.curriculums.id
Ref: academic.curriculums.effective_start_period_id > ref.academic_periods.id
Ref: academic.curriculums.effective_end_period_id > ref.academic_periods.id
Ref: academic.academic_period_study_program_settings.academic_period_id > ref.academic_periods.id
Ref: academic.academic_period_study_program_settings.study_program_id > ref.study_programs.id
Ref: academic.course_offerings.curriculum_id > academic.curriculums.id

/* =========================
   EVENT CONTRACT CATALOG v1.0.1
   Berlaku sebagai katalog teknis event lintas modul.
========================= */

Table core.event_contracts {
  id uuid [pk]
  event_name varchar [not null, note: 'format domain.action, contoh finance.payment_paid']
  event_version varchar [not null, default: 'v1']
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  publisher_module varchar [not null]
  publisher_database varchar
  aggregate_type varchar [not null]
  payload_schema jsonb [not null]
  validation_schema jsonb
  status varchar [not null, default: 'active', note: 'draft, active, deprecated, retired']
  backward_compatible boolean [default: true]
  description text
  created_at timestamp
  updated_at timestamp

  indexes {
    (event_name, event_version) [unique]
    publisher_module
    status
  }
}

Table core.event_consumers {
  id uuid [pk]
  event_contract_id uuid [not null]
  consumer_module varchar [not null]
  handler_name varchar
  retry_policy jsonb
  dlq_enabled boolean [default: true]
  max_retry int [default: 10]
  is_active boolean [default: true]
  created_at timestamp
  updated_at timestamp

  indexes {
    (event_contract_id, consumer_module) [unique]
    consumer_module
    is_active
  }
}

Table core.event_replay_logs {
  id uuid [pk]
  event_key varchar [not null]
  event_name varchar [not null]
  event_version varchar
  source_module varchar [not null]
  consumer_module varchar
  replay_reason text [not null]
  replayed_by uuid
  replayed_at timestamp
  replay_status varchar [not null, note: 'queued, success, failed']
  last_error text
  audit_ref_id uuid

  indexes {
    event_key
    event_name
    source_module
    consumer_module
    replay_status
    replayed_at
  }
}

/* --- CORE EVENT INFRASTRUCTURE --- */
Table core.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table core.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'core']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table core.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'core']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- REF EVENT INFRASTRUCTURE --- */
Table ref.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table ref.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'ref']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table ref.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'ref']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table ref.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'ref']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- CRM EVENT INFRASTRUCTURE --- */
Table crm.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table crm.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'crm']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table crm.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'crm']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table crm.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'crm']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- PMB EVENT INFRASTRUCTURE --- */
Table pmb.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table pmb.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'pmb']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table pmb.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'pmb']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table pmb.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'pmb']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- FINANCE EVENT INFRASTRUCTURE --- */
Table finance.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table finance.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'finance']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table finance.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'finance']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table finance.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'finance']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- ACADEMIC EVENT INFRASTRUCTURE --- */
Table academic.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table academic.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'academic']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table academic.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'academic']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table academic.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'academic']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- HRIS EVENT INFRASTRUCTURE --- */
Table hris.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table hris.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'hris']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table hris.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'hris']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table hris.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'hris']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- LMS EVENT INFRASTRUCTURE --- */
Table lms.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table lms.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'lms']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table lms.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'lms']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table lms.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'lms']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- ASSESSMENT EVENT INFRASTRUCTURE --- */
Table assessment.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table assessment.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'assessment']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table assessment.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'assessment']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table assessment.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'assessment']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}

/* --- PORTAL EVENT INFRASTRUCTURE --- */
Table portal.outbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null, default: 'v1']
  event_key varchar [not null]
  event_type varchar [not null, note: 'DOMAIN_EVENT, INTEGRATION_EVENT, NOTIFICATION_EVENT, SNAPSHOT_EVENT']
  aggregate_type varchar [not null]
  aggregate_id uuid [not null]
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  idempotency_key varchar
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'PENDING', note: 'PENDING, PUBLISHED, RETRYING, FAILED, DLQ']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  occurred_at timestamp [not null]
  published_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    event_key [unique]
    event_name
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table portal.inbox_events {
  id uuid [pk]
  event_name varchar [not null]
  event_version varchar [not null]
  event_key varchar [not null]
  publisher_module varchar [not null]
  publisher_database varchar
  consumer_module varchar [not null, default: 'portal']
  aggregate_type varchar
  aggregate_id uuid
  payload jsonb [not null]
  payload_hash text
  headers jsonb
  correlation_id varchar
  causation_id varchar
  status varchar [not null, default: 'RECEIVED', note: 'RECEIVED, PROCESSED, RETRYING, FAILED, DLQ, IGNORED_DUPLICATE']
  retry_count int [default: 0]
  max_retry int [default: 10]
  next_retry_at timestamp
  locked_at timestamp
  locked_by varchar
  last_error text
  received_at timestamp
  processed_at timestamp
  dead_letter_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    (consumer_module, event_key) [unique]
    event_name
    publisher_module
    status
    next_retry_at
    (aggregate_type, aggregate_id)
    correlation_id
    dead_letter_at
  }
}

Table portal.idempotency_keys {
  id uuid [pk]
  idempotency_key varchar [not null]
  source_module varchar
  target_module varchar [default: 'portal']
  request_hash text
  response_payload jsonb
  status varchar [not null, default: 'processing', note: 'processing, completed, failed, expired']
  locked_until timestamp
  trace_id varchar
  correlation_id varchar
  last_error text
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
  completed_at timestamp

  indexes {
    idempotency_key [unique]
    status
    locked_until
    expires_at
    correlation_id
  }
}

Table portal.reconciliation_mismatch_logs {
  id uuid [pk]
  source_module varchar [not null]
  source_table varchar [not null]
  source_ref_id uuid
  consumer_module varchar [not null, default: 'portal']
  consumer_table varchar
  consumer_ref_id uuid
  source_event_key varchar
  mismatch_type varchar [not null, note: 'missing_source, missing_snapshot, value_mismatch, stale_snapshot, duplicate_projection']
  source_value jsonb
  snapshot_value jsonb
  status varchar [not null, default: 'OPEN', note: 'OPEN, CORRECTED, IGNORED, PENDING_REVIEW']
  reason text
  detected_at timestamp
  corrected_at timestamp
  ignored_at timestamp
  created_at timestamp
  updated_at timestamp

  indexes {
    status
    mismatch_type
    source_module
    consumer_module
    source_event_key
    detected_at
  }
}


/* =========================
   STANDARD SNAPSHOT / READ MODEL FIELDS v1.0.1
   Setiap tabel snapshot/read model pada modul consumer wajib menambahkan field berikut pada migration SQL:
   - source_event_key varchar
   - source_event_name varchar
   - source_updated_at timestamp
   - synced_at timestamp
   - sync_status varchar note: 'fresh, stale, failed, pending_review'
   - reconciliation_status varchar note: 'matched, mismatch, corrected, ignored'

   Standar ini berlaku untuk contoh projection:
   - pmb applicant invoice status snapshot dari Finance
   - academic student clearance snapshot dari Finance
   - lms academic class snapshot dari Academic
   - lms enrollment projection dari Academic KRS
   - portal dashboard read model dari semua modul
========================= */

/* =========================
   EVENT CONTRACT INTERNAL REFERENCES
========================= */
Ref: core.event_consumers.event_contract_id > core.event_contracts.id
Ref: core.event_replay_logs.audit_ref_id > core.audit_logs.id


```