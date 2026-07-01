# Database Relation Diagram (ERD) Global - UNSIA ERP

## 1. Arsitektur Database Terdistribusi
Dalam desain modular terdistribusi v6.5, tidak ada foreign key tingkat database (*physical constraints*) lintas modul. Setiap modul mengisolasi datanya secara independen. Diagram di bawah mendefinisikan tabel-tabel utama di setiap modul dan hubungan logis antarmodul.

---

## 2. Katalog Skema DBML Global (Copy-paste ke dbdiagram.io)
```dbml
Project unsia_erp_global_master {
  database_type: 'PostgreSQL'
  Note: 'UNSIA ERP Global Master Database Schema. Mengisolasi 10 database terdistribusi secara fisik.'
}

// ==========================================
// 1. CORE DATABASE (core_db)
// ==========================================
Table core.persons {
  id uuid [pk]
  full_name varchar [not null]
  email varchar
  phone varchar
  identity_number varchar
  gender varchar
  address text
}

Table core.users {
  id uuid [pk]
  person_id uuid [not null]
  username varchar [unique, not null]
  email varchar [unique, not null]
  password_hash text [not null]
  status varchar [not null] // active, inactive, suspended
}

Table core.roles {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  scope_type varchar
  is_active boolean [default: true]
}

Table core.permissions {
  id uuid [pk]
  code varchar [unique, not null]
  module varchar [not null]
  resource varchar [not null]
  action varchar [not null]
}

Table core.user_roles {
  id uuid [pk]
  user_id uuid [not null]
  role_id uuid [not null]
  study_program_id uuid
}

Table core.role_permissions {
  id uuid [pk]
  role_id uuid [not null]
  permission_id uuid [not null]
}

Ref: core.users.person_id > core.persons.id
Ref: core.user_roles.user_id > core.users.id
Ref: core.user_roles.role_id > core.roles.id
Ref: core.role_permissions.role_id > core.roles.id
Ref: core.role_permissions.permission_id > core.permissions.id


// ==========================================
// 2. REFERENCE DATABASE (reference_db)
// ==========================================
Table ref.study_programs {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  degree_level varchar
  is_active boolean
}

Table ref.academic_years {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.academic_periods {
  id uuid [pk]
  academic_year_id uuid [not null]
  code varchar [unique, not null]
  name varchar [not null]
  status varchar
}

Table ref.pmb_waves {
  id uuid [pk]
  academic_year_id uuid
  target_entry_period_id uuid [not null]
  admission_path_id uuid
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table ref.payment_components {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
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


// ==========================================
// 3. CRM DATABASE (crm_db)
// ==========================================
Table crm.campaigns {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar [not null]
  status varchar
}

Table crm.agents {
  id uuid [pk]
  person_id uuid [not null]
  agent_code varchar [unique, not null]
  status varchar
}

Table crm.referrals {
  id uuid [pk]
  agent_id uuid
  referral_code varchar [unique, not null]
  is_valid boolean
}

Table crm.leads {
  id uuid [pk]
  person_id uuid [not null]
  campaign_id uuid
  referral_id uuid
  lead_number varchar [unique, not null]
  status varchar [not null]
}

Ref: crm.referrals.agent_id > crm.agents.id
Ref: crm.leads.campaign_id > crm.campaigns.id
Ref: crm.leads.referral_id > crm.referrals.id


// ==========================================
// 4. PMB DATABASE (pmb_db)
// ==========================================
Table pmb.applicants {
  id uuid [pk]
  person_id uuid [not null]
  user_id uuid
  crm_lead_id uuid
  study_program_id uuid
  pmb_wave_id uuid
  registration_number varchar [unique, not null]
  status varchar [not null]
}

Table pmb.applicant_biodata {
  id uuid [pk]
  applicant_id uuid [not null]
  full_name varchar
  nik varchar
}

Table pmb.applicant_documents {
  id uuid [pk]
  applicant_id uuid [not null]
  document_type_id uuid [not null]
  file_url text
  verification_status varchar
}

Table pmb.re_registrations {
  id uuid [pk]
  applicant_id uuid [not null]
  invoice_ref_id uuid
  status varchar
}

Ref: pmb.applicant_biodata.applicant_id > pmb.applicants.id
Ref: pmb.applicant_documents.applicant_id > pmb.applicants.id
Ref: pmb.re_registrations.applicant_id > pmb.applicants.id


// ==========================================
// 5. FINANCE DATABASE (finance_db)
// ==========================================
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
}

Table finance.invoice_items {
  id uuid [pk]
  invoice_id uuid [not null]
  payment_component_id uuid
  amount numeric
}

Table finance.payments {
  id uuid [pk]
  invoice_id uuid [not null]
  payment_method_id uuid
  payment_number varchar [unique]
  amount numeric
  payment_status varchar
  external_reference varchar
}

Table finance.student_clearances {
  id uuid [pk]
  student_id uuid [not null]
  academic_period_id uuid
  status varchar
}

Ref: finance.invoice_items.invoice_id > finance.invoices.id
Ref: finance.payments.invoice_id > finance.invoices.id


// ==========================================
// 6. ACADEMIC DATABASE (academic_db)
// ==========================================
Table academic.students {
  id uuid [pk]
  person_id uuid [not null]
  user_id uuid
  applicant_id uuid
  study_program_id uuid [not null]
  nim varchar [unique, not null]
  student_status varchar
  curriculum_id uuid
}

Table academic.curriculums {
  id uuid [pk]
  study_program_id uuid [not null]
  code varchar [unique, not null]
  name varchar [not null]
  is_active boolean
}

Table academic.courses {
  id uuid [pk]
  study_program_id uuid
  course_code varchar [unique, not null]
  course_name varchar [not null]
  sks int
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
}

Ref: academic.course_offerings.course_id > academic.courses.id
Ref: academic.classes.course_offering_id > academic.course_offerings.id
Ref: academic.krs.student_id > academic.students.id
Ref: academic.krs_items.krs_id > academic.krs.id
Ref: academic.krs_items.class_id > academic.classes.id
Ref: academic.grades.krs_item_id - academic.krs_items.id


// ==========================================
// 7. HRIS DATABASE (hris_db)
// ==========================================
Table hris.employees {
  id uuid [pk]
  person_id uuid [not null]
  nip varchar [unique]
  employment_status varchar
  is_active boolean
}

Table hris.lecturers {
  id uuid [pk]
  employee_id uuid [not null]
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
}

Table hris.leave_requests {
  id uuid [pk]
  employee_id uuid [not null]
  leave_type varchar
  start_date date
  end_date date
  status varchar
}

Ref: hris.lecturers.employee_id - hris.employees.id
Ref: hris.attendances.employee_id > hris.employees.id
Ref: hris.leave_requests.employee_id > hris.employees.id


// ==========================================
// 8. LMS DATABASE (lms_db)
// ==========================================
Table lms.classes {
  id uuid [pk]
  academic_class_id uuid [not null]
  lecturer_id uuid
  status varchar
}

Table lms.enrollments {
  id uuid [pk]
  lms_class_id uuid [not null]
  student_id uuid [not null]
  enrollment_status varchar
}

Table lms.sessions {
  id uuid [pk]
  lms_class_id uuid [not null]
  session_number int
  title varchar
}

Table lms.materials {
  id uuid [pk]
  session_id uuid [not null]
  title varchar
  file_url text
}

Table lms.assignments {
  id uuid [pk]
  session_id uuid [not null]
  title varchar
  due_at timestamp
}

Table lms.assignment_submissions {
  id uuid [pk]
  assignment_id uuid [not null]
  student_id uuid [not null]
  file_url text
  score numeric
}

Ref: lms.enrollments.lms_class_id > lms.classes.id
Ref: lms.sessions.lms_class_id > lms.classes.id
Ref: lms.materials.session_id > lms.sessions.id
Ref: lms.assignments.session_id > lms.sessions.id
Ref: lms.assignment_submissions.assignment_id > lms.assignments.id


// ==========================================
// 9. ASSESSMENT DATABASE (assessment_db)
// ==========================================
Table assessment.question_banks {
  id uuid [pk]
  code varchar [unique, not null]
  name varchar
  status varchar
}

Table assessment.questions {
  id uuid [pk]
  question_bank_id uuid [not null]
  question_type varchar
  question_text text
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
  session_type varchar
  context_module varchar
  context_id uuid
  title varchar
  duration_minutes int
}

Table assessment.assessment_participants {
  id uuid [pk]
  assessment_session_id uuid [not null]
  participant_type varchar
  applicant_id uuid
  student_id uuid
}

Table assessment.assessment_attempts {
  id uuid [pk]
  assessment_session_id uuid [not null]
  participant_id uuid [not null]
  attempt_number int
  status varchar
  total_score numeric
}

Table assessment.assessment_answers {
  id uuid [pk]
  attempt_id uuid [not null]
  question_id uuid [not null]
  selected_option_id uuid
  score numeric
}

Ref: assessment.questions.question_bank_id > assessment.question_banks.id
Ref: assessment.question_options.question_id > assessment.questions.id
Ref: assessment.assessment_participants.assessment_session_id > assessment.assessment_sessions.id
Ref: assessment.assessment_attempts.assessment_session_id > assessment.assessment_sessions.id
Ref: assessment.assessment_attempts.participant_id > assessment.assessment_participants.id
Ref: assessment.assessment_answers.attempt_id > assessment.assessment_attempts.id
Ref: assessment.assessment_answers.question_id > assessment.questions.id


// ==========================================
// 10. PORTAL DATABASE (portal_db)
// ==========================================
Table portal.notifications {
  id uuid [pk]
  user_id uuid [not null]
  title varchar
  message text
  module_source varchar
}

Table portal.notification_reads {
  id uuid [pk]
  notification_id uuid [not null]
  user_id uuid [not null]
  read_at timestamp
}

Table portal.user_preferences {
  id uuid [pk]
  user_id uuid [not null]
  preference_key varchar
  preference_value jsonb
}

Table portal.menu_shortcuts {
  id uuid [pk]
  user_id uuid [not null]
  menu_code varchar
  menu_label varchar
}

Ref: portal.notification_reads.notification_id > portal.notifications.id
```
