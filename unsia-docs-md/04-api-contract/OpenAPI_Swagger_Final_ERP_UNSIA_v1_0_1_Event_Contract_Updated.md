---
title: "UNSIA ERP / SIAKAD Terintegrasi API"
source_file: "OpenAPI_Swagger_Final_ERP_UNSIA_v1_0_1_Event_Contract_Updated.json"
format: markdown
---

# UNSIA ERP / SIAKAD Terintegrasi API

**Version:** `1.0.1-event-contract`

## Description

OpenAPI/Swagger final baseline untuk ERP Pendidikan / SIAKAD Terintegrasi UNSIA. Spesifikasi ini menurunkan API Contract dan Integration Contract v1.0 menjadi kontrak OpenAPI untuk Swagger UI, backend implementation, integration test, dan UAT lintas modul.

Prinsip: semua endpoint protected memakai Core Auth, active role, permission, application, dan data scope; endpoint kritis memakai Idempotency-Key; request lintas modul membawa X-Correlation-Id; cross-domain write dilakukan melalui API/service event resmi; response memakai envelope standar.

Update v1.0.1: menambahkan Event Contract API untuk event catalog, outbox, inbox, DLQ replay, reconciliation mismatch, event_key, event_version, retry, correlation_id, causation_id, dan observability endpoint.

## Servers

| URL | Description |
|---|---|
| `https://api.unsia.ac.id` | Production placeholder |
| `https://staging-api.unsia.ac.id` | Staging placeholder |
| `http://localhost:8000` | Local development |

## Tags / Module Area

| Tag | Description |
|---|---|
| `Core` | Identity, SSO, RBAC, audit, service token, idempotency, integration control. |
| `Referensi` | Master data lintas modul. |
| `CRM` | Lead, campaign, referral, dan konversi lead. |
| `PMB` | Applicant, dokumen, LoA, invoice request, dan handover akademik. |
| `Finance` | Invoice, payment callback, manual verification, clearance. |
| `Academic` | Mahasiswa, NIM, kelas, KRS, source grade, finalisasi nilai. |
| `HRIS` | Lecturer active read model. |
| `LMS` | Class sync, enrollment sync, grade sync. |
| `Assessment` | Session, attempt, scoring, result publish. |
| `Portal` | Notification dan dashboard role-based. |
| `Event Contract` | Event catalog, outbox/inbox monitoring, retry, DLQ replay, reconciliation mismatch, and event observability. |

## Endpoint Catalog

| Method | Path | Tags | Operation ID | Summary | Idempotent | Permission/Role |
|---|---|---|---|---|---|---|
| `POST` | `/api/v1/auth/login` | Core | `login` | Login user dan membuat session | `False` | Public |
| `POST` | `/api/v1/auth/refresh` | Core | `refreshToken` | Refresh access token | `False` | Authenticated |
| `GET` | `/api/v1/auth/me` | Core | `getAuthenticatedUser` | Mengambil profil user, active role, permission, dan scope | `False` | Authenticated |
| `POST` | `/api/v1/auth/switch-role` | Core | `switchActiveRole` | Mengubah active role session | `True` | Authenticated |
| `GET` | `/api/v1/applications` | Core | `listApplications` | Application launcher berdasarkan role | `False` | Authenticated |
| `POST` | `/api/v1/impersonations/start` | Core | `startImpersonation` | Memulai impersonation dengan reason | `True` | admin_bppti |
| `GET` | `/api/v1/ref/study-programs` | Referensi | `listStudyPrograms` | Daftar program studi | `False` | Authenticated |
| `GET` | `/api/v1/ref/academic-years` | Referensi | `listAcademicYears` | Daftar Tahun Ajaran | `False` | Authenticated |
| `POST` | `/api/v1/ref/academic-years` | Referensi | `createAcademicYear` | Membuat Tahun Ajaran operasional | `True` | admin_akademik_biro |
| `GET` | `/api/v1/ref/academic-periods` | Referensi | `listAcademicPeriods` | Daftar Periode Akademik | `False` | Authenticated |
| `POST` | `/api/v1/ref/academic-periods` | Referensi | `createAcademicPeriod` | Membuat Periode Akademik di bawah Tahun Ajaran | `True` | admin_akademik_biro |
| `GET` | `/api/v1/crm/leads` | CRM | `listLeads` | List lead sesuai scope | `False` | admin_crm \| agen_mitra |
| `POST` | `/api/v1/crm/leads` | CRM | `createLead` | Membuat lead | `True` | admin_crm \| agen_mitra |
| `POST` | `/api/v1/crm/leads/{lead_id}/convert-to-applicant` | CRM | `convertLeadToApplicant` | Convert lead qualified ke PMB applicant | `True` | admin_crm |
| `POST` | `/api/v1/pmb/applicants` | PMB | `createApplicant` | Membuat applicant | `True` | pendaftar \| admin_pmb |
| `POST` | `/api/v1/pmb/applicants/{applicant_id}/submit` | PMB | `submitApplicant` | Submit pendaftaran | `True` | pendaftar |
| `POST` | `/api/v1/pmb/applicants/{applicant_id}/documents` | PMB | `uploadApplicantDocument` | Upload dokumen applicant | `True` | pendaftar |
| `POST` | `/api/v1/pmb/applicants/{applicant_id}/documents/{document_id}/verify` | PMB | `verifyApplicantDocument` | Verifikasi dokumen | `True` | admin_pmb |
| `POST` | `/api/v1/pmb/applicants/{applicant_id}/request-invoice` | PMB | `requestApplicantInvoice` | Meminta invoice PMB ke Finance | `True` | admin_pmb |
| `POST` | `/api/v1/pmb/applicants/{applicant_id}/issue-loa` | PMB | `issueLoa` | Menerbitkan LoA setelah syarat terpenuhi | `True` | admin_pmb |
| `POST` | `/api/v1/pmb/applicants/{applicant_id}/handover-to-academic` | PMB | `handoverApplicantToAcademic` | Handover applicant ke Akademik | `True` | admin_pmb |
| `POST` | `/api/v1/pmb/applicants/{applicant_id}/selection-results` | PMB | `receiveAssessmentSelectionResult` | Menerima hasil seleksi/CBT dari Assessment | `True` | service:Assessment |
| `POST` | `/api/v1/finance/invoices` | Finance | `createInvoice` | Membuat invoice | `True` | service:PMB \| service:Academic |
| `GET` | `/api/v1/finance/invoices/{invoice_id}` | Finance | `getInvoice` | Membaca invoice | `False` | Authorized consumer |
| `POST` | `/api/v1/finance/payment-callbacks/{provider}` | Finance | `receivePaymentCallback` | Menerima callback payment gateway | `True` | payment provider |
| `POST` | `/api/v1/finance/payment-verifications` | Finance | `verifyManualPayment` | Verifikasi pembayaran manual | `True` | admin_keuangan |
| `GET` | `/api/v1/finance/clearances` | Finance | `checkClearance` | Cek clearance applicant/mahasiswa | `False` | Academic \| LMS \| Portal |
| `POST` | `/api/v1/academic/students/generate-from-applicant` | Academic | `generateStudentFromApplicant` | Generate mahasiswa dan NIM dari applicant | `True` | service:PMB \| admin_akademik_biro |
| `GET` | `/api/v1/academic/students` | Academic | `listStudents` | List mahasiswa sesuai scope | `False` | admin_akademik \| kaprodi |
| `POST` | `/api/v1/academic/classes` | Academic | `createAcademicClass` | Membuka kelas pada periode akademik | `True` | admin_akademik |
| `POST` | `/api/v1/academic/krs` | Academic | `createKrsDraft` | Membuat draft KRS | `True` | mahasiswa \| admin_akademik |
| `POST` | `/api/v1/academic/krs/{krs_id}/submit` | Academic | `submitKrs` | Submit KRS untuk approval | `True` | mahasiswa |
| `POST` | `/api/v1/academic/krs/{krs_id}/approve` | Academic | `approveKrs` | Approval KRS | `True` | dosen_pa |
| `POST` | `/api/v1/academic/grades/source-imports` | Academic | `importGradeSource` | Import source grade | `True` | service:LMS \| service:Assessment |
| `POST` | `/api/v1/academic/grades/{grade_id}/finalize` | Academic | `finalizeGrade` | Finalisasi nilai akademik | `True` | dosen \| admin_akademik |
| `GET` | `/api/v1/hris/lecturers` | HRIS | `listActiveLecturers` | Membaca dosen aktif | `False` | Academic \| LMS |
| `GET` | `/api/v1/hris/lecturers/{lecturer_id}` | HRIS | `getLecturer` | Detail dosen aktif | `False` | Academic \| LMS |
| `POST` | `/api/v1/lms/classes/sync-from-academic` | LMS | `syncClassFromAcademic` | Sync kelas akademik ke LMS | `True` | service:Academic |
| `POST` | `/api/v1/lms/enrollments/sync-from-krs` | LMS | `syncEnrollmentFromKrs` | Sync peserta kelas dari KRS valid | `True` | service:Academic |
| `POST` | `/api/v1/lms/grade-syncs` | LMS | `syncLmsGrade` | Kirim nilai aktivitas ke Academic | `True` | dosen \| service:LMS |
| `POST` | `/api/v1/assessment/sessions` | Assessment | `createAssessmentSession` | Membuat assessment session | `True` | PMB \| LMS \| admin_assessment |
| `POST` | `/api/v1/assessment/attempts` | Assessment | `createAssessmentAttempt` | Membuat attempt peserta | `True` | participant |
| `POST` | `/api/v1/assessment/results/publish` | Assessment | `publishAssessmentResult` | Publish hasil ke consumer | `True` | admin_assessment \| service:Assessment |
| `POST` | `/api/v1/portal/notifications` | Portal | `createPortalNotification` | Membuat notifikasi portal | `True` | service:any module |
| `GET` | `/api/v1/portal/dashboard` | Portal | `getRoleBasedDashboard` | Dashboard sesuai role | `False` | Authenticated |
| `GET` | `/api/v1/integration/event-contracts` | Event Contract | `listEventContracts` | List event contract catalog | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/event-contracts/{event_name}` | Event Contract | `getEventContractByName` | Detail event contract berdasarkan event_name | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/outbox-events` | Event Contract | `listOutboxEvents` | List outbox events | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/inbox-events` | Event Contract | `listInboxEvents` | List inbox events | `False` | technical_admin \| devops_sre \| auditor |
| `GET` | `/api/v1/integration/dlq-events` | Event Contract | `listDlqEvents` | List DLQ events | `False` | technical_admin \| devops_sre \| auditor |
| `POST` | `/api/v1/integration/dlq-events/{event_key}/replay` | Event Contract | `replayDlqEvent` | Replay event dari DLQ | `True` | devops_sre |
| `GET` | `/api/v1/integration/reconciliation-mismatches` | Event Contract | `listReconciliationMismatches` | List reconciliation mismatch logs | `False` | technical_admin \| devops_sre \| auditor |
| `POST` | `/api/v1/integration/reconciliation-mismatches/{mismatch_id}/resolve` | Event Contract | `resolveReconciliationMismatch` | Resolve reconciliation mismatch | `True` | technical_admin \| devops_sre \| owner_modul |

## Component Schemas

### `Meta`

| Field | Type | Description |
|---|---|---|
| `request_id` | `string` |  |
| `correlation_id` | `string` |  |
| `timestamp` | `string` |  |

### `PaginationMeta`

| Field | Type | Description |
|---|---|---|
| `page` | `integer` |  |
| `limit` | `integer` |  |
| `total` | `integer` |  |
| `total_page` | `integer` |  |
| `sort` | `string` |  |

### `SuccessEnvelope`

| Field | Type | Description |
|---|---|---|
| `success` | `boolean` |  |
| `code` | `string` |  |
| `message` | `string` |  |
| `data` | `object` |  |
| `meta` | `#/components/schemas/Meta` |  |

### `ErrorDetail`

| Field | Type | Description |
|---|---|---|
| `field` | `string` |  |
| `message` | `string` |  |

### `ErrorEnvelope`

| Field | Type | Description |
|---|---|---|
| `success` | `boolean` |  |
| `code` | `string` |  |
| `message` | `string` |  |
| `errors` | `array` |  |
| `meta` | `#/components/schemas/Meta` |  |

### `TokenResponse`

| Field | Type | Description |
|---|---|---|
| `access_token` | `string` |  |
| `refresh_token` | `string` |  |
| `expires_in` | `integer` |  |
| `token_type` | `string` |  |

### `UserProfile`

| Field | Type | Description |
|---|---|---|
| `user_id` | `string` |  |
| `person_id` | `string` |  |
| `name` | `string` |  |
| `email` | `string` |  |
| `active_role` | `string` |  |
| `permissions` | `array` |  |
| `scope` | `object` |  |

### `Application`

| Field | Type | Description |
|---|---|---|
| `application_code` | `string` |  |
| `name` | `string` |  |
| `url` | `string` |  |
| `enabled` | `boolean` |  |

### `StudyProgram`

| Field | Type | Description |
|---|---|---|
| `study_program_id` | `string` |  |
| `code` | `string` |  |
| `name` | `string` |  |
| `degree` | `string` |  |
| `status` | `string` |  |

### `AcademicYear`

| Field | Type | Description |
|---|---|---|
| `academic_year_id` | `string` |  |
| `code` | `string` |  |
| `name` | `string` |  |
| `status` | `string` |  |

### `AcademicPeriod`

| Field | Type | Description |
|---|---|---|
| `academic_period_id` | `string` |  |
| `academic_year_id` | `string` |  |
| `code` | `string` |  |
| `term` | `string` |  |
| `status` | `string` |  |

### `Lead`

| Field | Type | Description |
|---|---|---|
| `lead_id` | `string` |  |
| `name` | `string` |  |
| `email` | `string` |  |
| `phone` | `string` |  |
| `status` | `string` |  |

### `ConvertLeadResponse`

| Field | Type | Description |
|---|---|---|
| `lead_id` | `string` |  |
| `applicant_id` | `string` |  |
| `applicant_number` | `string` |  |
| `status` | `string` |  |

### `Applicant`

| Field | Type | Description |
|---|---|---|
| `applicant_id` | `string` |  |
| `applicant_number` | `string` |  |
| `status` | `string` |  |
| `study_program_id` | `string` |  |

### `Invoice`

| Field | Type | Description |
|---|---|---|
| `invoice_id` | `string` |  |
| `invoice_number` | `string` |  |
| `payer_type` | `string` |  |
| `payer_ref_id` | `string` |  |
| `amount_total` | `number` |  |
| `status` | `string` |  |

### `Student`

| Field | Type | Description |
|---|---|---|
| `student_id` | `string` |  |
| `nim` | `string` |  |
| `person_id` | `string` |  |
| `study_program_id` | `string` |  |
| `entry_period_id` | `string` |  |
| `status` | `string` |  |

### `AcademicClass`

| Field | Type | Description |
|---|---|---|
| `academic_class_id` | `string` |  |
| `academic_period_id` | `string` |  |
| `course_id` | `string` |  |
| `class_code` | `string` |  |
| `status` | `string` |  |

### `Krs`

| Field | Type | Description |
|---|---|---|
| `krs_id` | `string` |  |
| `student_id` | `string` |  |
| `academic_period_id` | `string` |  |
| `status` | `string` |  |

### `Grade`

| Field | Type | Description |
|---|---|---|
| `grade_id` | `string` |  |
| `academic_class_id` | `string` |  |
| `student_id` | `string` |  |
| `final_grade` | `number` |  |
| `letter_grade` | `string` |  |
| `status` | `string` |  |

### `Lecturer`

| Field | Type | Description |
|---|---|---|
| `lecturer_id` | `string` |  |
| `person_id` | `string` |  |
| `name` | `string` |  |
| `homebase_study_program_id` | `string` |  |
| `employment_status` | `string` |  |

### `LmsSyncResponse`

| Field | Type | Description |
|---|---|---|
| `sync_status` | `string` |  |
| `lms_class_id` | `string` |  |
| `lms_enrollment_id` | `string` |  |

### `AssessmentSession`

| Field | Type | Description |
|---|---|---|
| `assessment_session_id` | `string` |  |
| `context` | `string` |  |
| `status` | `string` |  |

### `AssessmentAttempt`

| Field | Type | Description |
|---|---|---|
| `attempt_id` | `string` |  |
| `assessment_session_id` | `string` |  |
| `status` | `string` |  |

### `Notification`

| Field | Type | Description |
|---|---|---|
| `notification_id` | `string` |  |
| `recipient_user_id` | `string` |  |
| `event_type` | `string` |  |
| `read_status` | `boolean` |  |

### `DashboardResponse`

| Field | Type | Description |
|---|---|---|
| `role` | `string` |  |
| `widgets` | `array` |  |

### `ClearanceResponse`

| Field | Type | Description |
|---|---|---|
| `student_id` | `string` |  |
| `applicant_id` | `string` |  |
| `service_code` | `string` |  |
| `academic_period_id` | `string` |  |
| `clearance_status` | `string` |  |
| `block_reasons` | `array` |  |

### `LoginRequest`

| Field | Type | Description |
|---|---|---|
| `username` | `string` |  |
| `password` | `string` |  |
| `captcha_token` | `string` |  |

### `SwitchRoleRequest`

| Field | Type | Description |
|---|---|---|
| `role_code` | `string` |  |
| `scope_value` | `string` |  |

### `ImpersonationStartRequest`

| Field | Type | Description |
|---|---|---|
| `target_user_id` | `string` |  |
| `reason` | `string` |  |
| `duration_minutes` | `integer` |  |

### `AcademicYearCreateRequest`

| Field | Type | Description |
|---|---|---|
| `code` | `string` |  |
| `name` | `string` |  |
| `start_date` | `string` |  |
| `end_date` | `string` |  |
| `status` | `string` |  |

### `AcademicPeriodCreateRequest`

| Field | Type | Description |
|---|---|---|
| `academic_year_id` | `string` |  |
| `code` | `string` |  |
| `term` | `string` |  |
| `start_date` | `string` |  |
| `end_date` | `string` |  |

### `LeadCreateRequest`

| Field | Type | Description |
|---|---|---|
| `name` | `string` |  |
| `email` | `string` |  |
| `phone` | `string` |  |
| `campaign_id` | `string` |  |
| `source` | `string` |  |
| `target_entry_period_id` | `string` |  |
| `study_program_id` | `string` |  |

### `ConvertLeadRequest`

| Field | Type | Description |
|---|---|---|
| `target_entry_period_id` | `string` |  |
| `admission_path_id` | `string` |  |
| `study_program_id` | `string` |  |
| `reason` | `string` |  |

### `ApplicantCreateRequest`

| Field | Type | Description |
|---|---|---|
| `person` | `object` |  |
| `target_entry_period_id` | `string` |  |
| `study_program_id` | `string` |  |
| `admission_path_id` | `string` |  |

### `SubmitRequest`

| Field | Type | Description |
|---|---|---|
| `reason` | `string` |  |

### `DocumentUploadRequest`

| Field | Type | Description |
|---|---|---|
| `document_type_code` | `string` |  |
| `file` | `string` |  |

### `VerifyDocumentRequest`

| Field | Type | Description |
|---|---|---|
| `verification_status` | `string` |  |
| `reason` | `string` |  |

### `RequestInvoiceRequest`

| Field | Type | Description |
|---|---|---|
| `invoice_context` | `string` |  |
| `payment_component_codes` | `array` |  |
| `due_date` | `string` |  |
| `reason` | `string` |  |

### `IssueLoaRequest`

| Field | Type | Description |
|---|---|---|
| `reason` | `string` |  |
| `template_code` | `string` |  |

### `HandoverToAcademicRequest`

| Field | Type | Description |
|---|---|---|
| `entry_academic_year_id` | `string` |  |
| `entry_period_id` | `string` |  |
| `study_program_id` | `string` |  |
| `curriculum_id` | `string` |  |
| `reason` | `string` |  |

### `HandoverToAcademicResponse`

| Field | Type | Description |
|---|---|---|
| `student_id` | `string` |  |
| `nim` | `string` |  |
| `handover_status` | `string` |  |

### `InvoiceCreateRequest`

| Field | Type | Description |
|---|---|---|
| `payer_type` | `string` |  |
| `payer_ref_id` | `string` |  |
| `academic_period_id` | `string` |  |
| `items` | `array` |  |
| `source_module` | `string` |  |
| `source_ref_id` | `string` |  |
| `due_date` | `string` |  |

### `PaymentCallbackRequest`

| Field | Type | Description |
|---|---|---|
| `provider_event_id` | `string` |  |
| `invoice_id` | `string` |  |
| `payment_id` | `string` |  |
| `amount` | `number` |  |
| `payment_status` | `string` |  |
| `signature_status` | `string` |  |
| `payload_hash` | `string` |  |
| `raw_payload` | `object` |  |

### `PaymentVerificationRequest`

| Field | Type | Description |
|---|---|---|
| `invoice_id` | `string` |  |
| `payment_id` | `string` |  |
| `verification_status` | `string` |  |
| `amount` | `number` |  |
| `paid_at` | `string` |  |
| `attachment_ref` | `string` |  |
| `reason` | `string` |  |

### `GenerateStudentFromApplicantRequest`

| Field | Type | Description |
|---|---|---|
| `applicant_id` | `string` |  |
| `curriculum_id` | `string` |  |
| `entry_academic_year_id` | `string` |  |
| `entry_period_id` | `string` |  |
| `study_program_id` | `string` |  |
| `reason` | `string` |  |

### `ClassCreateRequest`

| Field | Type | Description |
|---|---|---|
| `academic_period_id` | `string` |  |
| `course_id` | `string` |  |
| `class_code` | `string` |  |
| `lecturer_ids` | `array` |  |
| `capacity` | `integer` |  |

### `KrsCreateRequest`

| Field | Type | Description |
|---|---|---|
| `student_id` | `string` |  |
| `academic_period_id` | `string` |  |
| `items` | `array` |  |

### `GradeSourceImportRequest`

| Field | Type | Description |
|---|---|---|
| `source_module` | `string` |  |
| `source_ref_id` | `string` |  |
| `academic_class_id` | `string` |  |
| `student_id` | `string` |  |
| `component_code` | `string` |  |
| `score` | `number` |  |
| `max_score` | `number` |  |
| `submitted_at` | `string` |  |

### `GradeFinalizeRequest`

| Field | Type | Description |
|---|---|---|
| `final_grade` | `number` |  |
| `letter_grade` | `string` |  |
| `reason` | `string` |  |

### `LmsClassSyncRequest`

| Field | Type | Description |
|---|---|---|
| `academic_period_id` | `string` |  |
| `academic_class_id` | `string` |  |
| `course_id` | `string` |  |
| `class_code` | `string` |  |
| `lecturer_ids` | `array` |  |
| `schedule` | `array` |  |

### `EnrollmentSyncRequest`

| Field | Type | Description |
|---|---|---|
| `academic_class_id` | `string` |  |
| `krs_item_id` | `string` |  |
| `student_id` | `string` |  |
| `enrollment_status` | `string` |  |

### `LmsGradeSyncRequest`

| Field | Type | Description |
|---|---|---|
| `source_ref_id` | `string` |  |
| `academic_class_id` | `string` |  |
| `student_id` | `string` |  |
| `component_code` | `string` |  |
| `score` | `number` |  |
| `max_score` | `number` |  |
| `submitted_at` | `string` |  |

### `AssessmentSessionRequest`

| Field | Type | Description |
|---|---|---|
| `context` | `string` |  |
| `title` | `string` |  |
| `consumer_module` | `string` |  |
| `wave_id` | `string` |  |
| `academic_class_id` | `string` |  |
| `question_set_id` | `string` |  |
| `start_at` | `string` |  |
| `end_at` | `string` |  |

### `AssessmentAttemptRequest`

| Field | Type | Description |
|---|---|---|
| `assessment_session_id` | `string` |  |
| `participant_ref_id` | `string` |  |
| `started_at` | `string` |  |

### `AssessmentResultPublishRequest`

| Field | Type | Description |
|---|---|---|
| `attempt_id` | `string` |  |
| `consumer_module` | `string` |  |
| `score` | `number` |  |
| `result_status` | `string` |  |
| `target_ref_id` | `string` |  |

### `NotificationRequest`

| Field | Type | Description |
|---|---|---|
| `recipient_user_id` | `string` |  |
| `event_type` | `string` |  |
| `source_module` | `string` |  |
| `event_key` | `string` |  |
| `severity` | `string` |  |
| `title` | `string` |  |
| `message` | `string` |  |
| `link_url` | `string` |  |
| `expires_at` | `string` |  |

### `EventContract`

| Field | Type | Description |
|---|---|---|
| `id` | `string` |  |
| `event_name` | `string` |  |
| `event_version` | `string` |  |
| `event_type` | `string` |  |
| `publisher_module` | `string` |  |
| `publisher_database` | `string` |  |
| `aggregate_type` | `string` |  |
| `payload_schema` | `object` |  |
| `validation_schema` | `object` |  |
| `status` | `string` |  |
| `description` | `string` |  |
| `created_at` | `string` |  |
| `updated_at` | `string` |  |

### `EventEnvelope`

| Field | Type | Description |
|---|---|---|
| `event_name` | `string` |  |
| `event_version` | `string` |  |
| `event_key` | `string` |  |
| `event_type` | `string` |  |
| `publisher_service` | `string` |  |
| `publisher_database` | `string` |  |
| `aggregate_type` | `string` |  |
| `aggregate_id` | `string` |  |
| `correlation_id` | `string` |  |
| `causation_id` | `string` |  |
| `occurred_at` | `string` |  |
| `published_at` | `string` |  |
| `payload` | `object` |  |

### `OutboxEvent`
### `InboxEvent`
### `DlqEvent`

| Field | Type | Description |
|---|---|---|
| `event_key` | `string` |  |
| `event_name` | `string` |  |
| `event_version` | `string` |  |
| `publisher_module` | `string` |  |
| `consumer_module` | `string` |  |
| `retry_count` | `integer` |  |
| `last_error` | `string` |  |
| `dead_letter_at` | `string` |  |
| `payload` | `object` |  |

### `EventReplayRequest`

| Field | Type | Description |
|---|---|---|
| `reason` | `string` |  |
| `target_consumer` | `string` |  |
| `force` | `boolean` |  |

### `EventReplayResponse`

| Field | Type | Description |
|---|---|---|
| `event_key` | `string` |  |
| `replay_status` | `string` |  |
| `audit_ref_id` | `string` |  |
| `message` | `string` |  |

### `ReconciliationMismatch`

| Field | Type | Description |
|---|---|---|
| `id` | `string` |  |
| `source_module` | `string` |  |
| `source_table` | `string` |  |
| `source_ref_id` | `string` |  |
| `consumer_module` | `string` |  |
| `consumer_table` | `string` |  |
| `consumer_ref_id` | `string` |  |
| `source_event_key` | `string` |  |
| `mismatch_type` | `string` |  |
| `source_value` | `object` |  |
| `snapshot_value` | `object` |  |
| `status` | `string` |  |
| `reason` | `string` |  |
| `detected_at` | `string` |  |
| `corrected_at` | `string` |  |
| `ignored_at` | `string` |  |

### `ReconciliationResolveRequest`

| Field | Type | Description |
|---|---|---|
| `action` | `string` |  |
| `reason` | `string` |  |
| `correction_payload` | `object` |  |

### `ReconciliationResolveResponse`

| Field | Type | Description |
|---|---|---|
| `mismatch_id` | `string` |  |
| `status` | `string` |  |
| `audit_ref_id` | `string` |  |


## Raw OpenAPI JSON

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "UNSIA ERP / SIAKAD Terintegrasi API",
    "version": "1.0.1-event-contract",
    "description": "OpenAPI/Swagger final baseline untuk ERP Pendidikan / SIAKAD Terintegrasi UNSIA. Spesifikasi ini menurunkan API Contract dan Integration Contract v1.0 menjadi kontrak OpenAPI untuk Swagger UI, backend implementation, integration test, dan UAT lintas modul.\n\nPrinsip: semua endpoint protected memakai Core Auth, active role, permission, application, dan data scope; endpoint kritis memakai Idempotency-Key; request lintas modul membawa X-Correlation-Id; cross-domain write dilakukan melalui API/service event resmi; response memakai envelope standar.\n\nUpdate v1.0.1: menambahkan Event Contract API untuk event catalog, outbox, inbox, DLQ replay, reconciliation mismatch, event_key, event_version, retry, correlation_id, causation_id, dan observability endpoint."
  },
  "servers": [
    {
      "url": "https://api.unsia.ac.id",
      "description": "Production placeholder"
    },
    {
      "url": "https://staging-api.unsia.ac.id",
      "description": "Staging placeholder"
    },
    {
      "url": "http://localhost:8000",
      "description": "Local development"
    }
  ],
  "tags": [
    {
      "name": "Core",
      "description": "Identity, SSO, RBAC, audit, service token, idempotency, integration control."
    },
    {
      "name": "Referensi",
      "description": "Master data lintas modul."
    },
    {
      "name": "CRM",
      "description": "Lead, campaign, referral, dan konversi lead."
    },
    {
      "name": "PMB",
      "description": "Applicant, dokumen, LoA, invoice request, dan handover akademik."
    },
    {
      "name": "Finance",
      "description": "Invoice, payment callback, manual verification, clearance."
    },
    {
      "name": "Academic",
      "description": "Mahasiswa, NIM, kelas, KRS, source grade, finalisasi nilai."
    },
    {
      "name": "HRIS",
      "description": "Lecturer active read model."
    },
    {
      "name": "LMS",
      "description": "Class sync, enrollment sync, grade sync."
    },
    {
      "name": "Assessment",
      "description": "Session, attempt, scoring, result publish."
    },
    {
      "name": "Portal",
      "description": "Notification dan dashboard role-based."
    },
    {
      "name": "Event Contract",
      "description": "Event catalog, outbox/inbox monitoring, retry, DLQ replay, reconciliation mismatch, and event observability."
    }
  ],
  "security": [
    {
      "bearerAuth": []
    }
  ],
  "paths": {
    "/api/v1/auth/login": {
      "post": {
        "tags": [
          "Core"
        ],
        "operationId": "login",
        "summary": "Login user dan membuat session",
        "description": "Public endpoint untuk login. Tidak memakai bearer token.",
        "parameters": [],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/TokenResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Public",
        "x-idempotent": false,
        "x-audit-log": "audit login",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/LoginRequest"
              }
            }
          }
        },
        "security": []
      }
    },
    "/api/v1/auth/refresh": {
      "post": {
        "tags": [
          "Core"
        ],
        "operationId": "refreshToken",
        "summary": "Refresh access token",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/TokenResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": false,
        "x-audit-log": "audit session",
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/auth/me": {
      "get": {
        "tags": [
          "Core"
        ],
        "operationId": "getAuthenticatedUser",
        "summary": "Mengambil profil user, active role, permission, dan scope",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/UserProfile"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/auth/switch-role": {
      "post": {
        "tags": [
          "Core"
        ],
        "operationId": "switchActiveRole",
        "summary": "Mengubah active role session",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/UserProfile"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": true,
        "x-audit-log": "active_role_sessions",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SwitchRoleRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/applications": {
      "get": {
        "tags": [
          "Core"
        ],
        "operationId": "listApplications",
        "summary": "Application launcher berdasarkan role",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/Application"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/impersonations/start": {
      "post": {
        "tags": [
          "Core"
        ],
        "operationId": "startImpersonation",
        "summary": "Memulai impersonation dengan reason",
        "description": "Role/permission: admin_bppti.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/UserProfile"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_bppti",
        "x-idempotent": true,
        "x-audit-log": "impersonation audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ImpersonationStartRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/ref/study-programs": {
      "get": {
        "tags": [
          "Referensi"
        ],
        "operationId": "listStudyPrograms",
        "summary": "Daftar program studi",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/StudyProgram"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/ref/academic-years": {
      "get": {
        "tags": [
          "Referensi"
        ],
        "operationId": "listAcademicYears",
        "summary": "Daftar Tahun Ajaran",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/AcademicYear"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      },
      "post": {
        "tags": [
          "Referensi"
        ],
        "operationId": "createAcademicYear",
        "summary": "Membuat Tahun Ajaran operasional",
        "description": "Role/permission: admin_akademik_biro.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/AcademicYear"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_akademik_biro",
        "x-idempotent": true,
        "x-audit-log": "audit master",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/AcademicYearCreateRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/ref/academic-periods": {
      "get": {
        "tags": [
          "Referensi"
        ],
        "operationId": "listAcademicPeriods",
        "summary": "Daftar Periode Akademik",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/AcademicPeriod"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      },
      "post": {
        "tags": [
          "Referensi"
        ],
        "operationId": "createAcademicPeriod",
        "summary": "Membuat Periode Akademik di bawah Tahun Ajaran",
        "description": "Role/permission: admin_akademik_biro.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/AcademicPeriod"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_akademik_biro",
        "x-idempotent": true,
        "x-audit-log": "audit master",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/AcademicPeriodCreateRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/crm/leads": {
      "get": {
        "tags": [
          "CRM"
        ],
        "operationId": "listLeads",
        "summary": "List lead sesuai scope",
        "description": "Role/permission: admin_crm | agen_mitra.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Sort"
          },
          {
            "$ref": "#/components/parameters/Search"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/Lead"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_crm | agen_mitra",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      },
      "post": {
        "tags": [
          "CRM"
        ],
        "operationId": "createLead",
        "summary": "Membuat lead",
        "description": "Role/permission: admin_crm | agen_mitra.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Lead"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_crm | agen_mitra",
        "x-idempotent": true,
        "x-audit-log": "lead history",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/LeadCreateRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/crm/leads/{lead_id}/convert-to-applicant": {
      "post": {
        "tags": [
          "CRM"
        ],
        "operationId": "convertLeadToApplicant",
        "summary": "Convert lead qualified ke PMB applicant",
        "description": "Role/permission: admin_crm.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "lead_id",
            "in": "path",
            "required": true,
            "description": "lead_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/ConvertLeadResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_crm",
        "x-idempotent": true,
        "x-audit-log": "lead status history, integration log",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ConvertLeadRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "createApplicant",
        "summary": "Membuat applicant",
        "description": "Role/permission: pendaftar | admin_pmb.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Applicant"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "pendaftar | admin_pmb",
        "x-idempotent": true,
        "x-audit-log": "applicant status history",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ApplicantCreateRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants/{applicant_id}/submit": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "submitApplicant",
        "summary": "Submit pendaftaran",
        "description": "Role/permission: pendaftar.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "applicant_id",
            "in": "path",
            "required": true,
            "description": "applicant_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Applicant"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "pendaftar",
        "x-idempotent": true,
        "x-audit-log": "status history",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SubmitRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants/{applicant_id}/documents": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "uploadApplicantDocument",
        "summary": "Upload dokumen applicant",
        "description": "Role/permission: pendaftar.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "applicant_id",
            "in": "path",
            "required": true,
            "description": "applicant_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Applicant"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "pendaftar",
        "x-idempotent": true,
        "x-audit-log": "document audit",
        "requestBody": {
          "required": true,
          "content": {
            "multipart/form-data": {
              "schema": {
                "$ref": "#/components/schemas/DocumentUploadRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants/{applicant_id}/documents/{document_id}/verify": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "verifyApplicantDocument",
        "summary": "Verifikasi dokumen",
        "description": "Role/permission: admin_pmb.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "applicant_id",
            "in": "path",
            "required": true,
            "description": "applicant_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "name": "document_id",
            "in": "path",
            "required": true,
            "description": "document_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Applicant"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_pmb",
        "x-idempotent": true,
        "x-audit-log": "document audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/VerifyDocumentRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants/{applicant_id}/request-invoice": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "requestApplicantInvoice",
        "summary": "Meminta invoice PMB ke Finance",
        "description": "Role/permission: admin_pmb.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "applicant_id",
            "in": "path",
            "required": true,
            "description": "applicant_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Invoice"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_pmb",
        "x-idempotent": true,
        "x-audit-log": "integration log",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/RequestInvoiceRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants/{applicant_id}/issue-loa": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "issueLoa",
        "summary": "Menerbitkan LoA setelah syarat terpenuhi",
        "description": "Role/permission: admin_pmb.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "applicant_id",
            "in": "path",
            "required": true,
            "description": "applicant_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Applicant"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_pmb",
        "x-idempotent": true,
        "x-audit-log": "LoA audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/IssueLoaRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants/{applicant_id}/handover-to-academic": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "handoverApplicantToAcademic",
        "summary": "Handover applicant ke Akademik",
        "description": "Role/permission: admin_pmb.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "applicant_id",
            "in": "path",
            "required": true,
            "description": "applicant_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/HandoverToAcademicResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_pmb",
        "x-idempotent": true,
        "x-audit-log": "handover log",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/HandoverToAcademicRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/pmb/applicants/{applicant_id}/selection-results": {
      "post": {
        "tags": [
          "PMB"
        ],
        "operationId": "receiveAssessmentSelectionResult",
        "summary": "Menerima hasil seleksi/CBT dari Assessment",
        "description": "Role/permission: service:Assessment.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "applicant_id",
            "in": "path",
            "required": true,
            "description": "applicant_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Applicant"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "service:Assessment",
        "x-idempotent": true,
        "x-audit-log": "assessment score audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/AssessmentResultPublishRequest"
              }
            }
          }
        },
        "security": [
          {
            "serviceTokenAuth": []
          }
        ]
      }
    },
    "/api/v1/finance/invoices": {
      "post": {
        "tags": [
          "Finance"
        ],
        "operationId": "createInvoice",
        "summary": "Membuat invoice",
        "description": "Role/permission: service:PMB | service:Academic.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Invoice"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "service:PMB | service:Academic",
        "x-idempotent": true,
        "x-audit-log": "invoice audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/InvoiceCreateRequest"
              }
            }
          }
        },
        "security": [
          {
            "serviceTokenAuth": []
          }
        ]
      }
    },
    "/api/v1/finance/invoices/{invoice_id}": {
      "get": {
        "tags": [
          "Finance"
        ],
        "operationId": "getInvoice",
        "summary": "Membaca invoice",
        "description": "Role/permission: Authorized consumer.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "invoice_id",
            "in": "path",
            "required": true,
            "description": "invoice_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Invoice"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authorized consumer",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/finance/payment-callbacks/{provider}": {
      "post": {
        "tags": [
          "Finance"
        ],
        "operationId": "receivePaymentCallback",
        "summary": "Menerima callback payment gateway",
        "description": "Role/permission: payment provider.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "provider",
            "in": "path",
            "required": true,
            "description": "Payment gateway provider code",
            "schema": {
              "type": "string"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Invoice"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "payment provider",
        "x-idempotent": true,
        "x-audit-log": "callback log, payment audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/PaymentCallbackRequest"
              }
            }
          }
        },
        "security": [
          {
            "paymentProviderSignature": []
          }
        ]
      }
    },
    "/api/v1/finance/payment-verifications": {
      "post": {
        "tags": [
          "Finance"
        ],
        "operationId": "verifyManualPayment",
        "summary": "Verifikasi pembayaran manual",
        "description": "Role/permission: admin_keuangan.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Invoice"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_keuangan",
        "x-idempotent": true,
        "x-audit-log": "payment verification audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/PaymentVerificationRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/finance/clearances": {
      "get": {
        "tags": [
          "Finance"
        ],
        "operationId": "checkClearance",
        "summary": "Cek clearance applicant/mahasiswa",
        "description": "Role/permission: Academic | LMS | Portal.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "student_id",
            "in": "query",
            "required": false,
            "description": "student_id",
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "applicant_id",
            "in": "query",
            "required": false,
            "description": "applicant_id",
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "service_code",
            "in": "query",
            "required": true,
            "description": "service_code",
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "academic_period_id",
            "in": "query",
            "required": false,
            "description": "academic_period_id",
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/ClearanceResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Academic | LMS | Portal",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/students/generate-from-applicant": {
      "post": {
        "tags": [
          "Academic"
        ],
        "operationId": "generateStudentFromApplicant",
        "summary": "Generate mahasiswa dan NIM dari applicant",
        "description": "Role/permission: service:PMB | admin_akademik_biro.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Student"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "service:PMB | admin_akademik_biro",
        "x-idempotent": true,
        "x-audit-log": "NIM sequence, audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/GenerateStudentFromApplicantRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/students": {
      "get": {
        "tags": [
          "Academic"
        ],
        "operationId": "listStudents",
        "summary": "List mahasiswa sesuai scope",
        "description": "Role/permission: admin_akademik | kaprodi.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Sort"
          },
          {
            "$ref": "#/components/parameters/Search"
          },
          {
            "name": "filter[study_program_id]",
            "in": "query",
            "required": false,
            "description": "filter[study_program_id]",
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/Student"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_akademik | kaprodi",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/classes": {
      "post": {
        "tags": [
          "Academic"
        ],
        "operationId": "createAcademicClass",
        "summary": "Membuka kelas pada periode akademik",
        "description": "Role/permission: admin_akademik.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/AcademicClass"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_akademik",
        "x-idempotent": true,
        "x-audit-log": "class audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ClassCreateRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/krs": {
      "post": {
        "tags": [
          "Academic"
        ],
        "operationId": "createKrsDraft",
        "summary": "Membuat draft KRS",
        "description": "Role/permission: mahasiswa | admin_akademik.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Krs"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "mahasiswa | admin_akademik",
        "x-idempotent": true,
        "x-audit-log": "KRS audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/KrsCreateRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/krs/{krs_id}/submit": {
      "post": {
        "tags": [
          "Academic"
        ],
        "operationId": "submitKrs",
        "summary": "Submit KRS untuk approval",
        "description": "Role/permission: mahasiswa.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "krs_id",
            "in": "path",
            "required": true,
            "description": "krs_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Krs"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "mahasiswa",
        "x-idempotent": true,
        "x-audit-log": "KRS status history",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SubmitRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/krs/{krs_id}/approve": {
      "post": {
        "tags": [
          "Academic"
        ],
        "operationId": "approveKrs",
        "summary": "Approval KRS",
        "description": "Role/permission: dosen_pa.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "krs_id",
            "in": "path",
            "required": true,
            "description": "krs_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Krs"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "dosen_pa",
        "x-idempotent": true,
        "x-audit-log": "KRS status history",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SubmitRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/grades/source-imports": {
      "post": {
        "tags": [
          "Academic"
        ],
        "operationId": "importGradeSource",
        "summary": "Import source grade",
        "description": "Role/permission: service:LMS | service:Assessment.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Grade"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "service:LMS | service:Assessment",
        "x-idempotent": true,
        "x-audit-log": "grade source audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/GradeSourceImportRequest"
              }
            }
          }
        },
        "security": [
          {
            "serviceTokenAuth": []
          }
        ]
      }
    },
    "/api/v1/academic/grades/{grade_id}/finalize": {
      "post": {
        "tags": [
          "Academic"
        ],
        "operationId": "finalizeGrade",
        "summary": "Finalisasi nilai akademik",
        "description": "Role/permission: dosen | admin_akademik.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "grade_id",
            "in": "path",
            "required": true,
            "description": "grade_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Grade"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "dosen | admin_akademik",
        "x-idempotent": true,
        "x-audit-log": "grade history",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/GradeFinalizeRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/hris/lecturers": {
      "get": {
        "tags": [
          "HRIS"
        ],
        "operationId": "listActiveLecturers",
        "summary": "Membaca dosen aktif",
        "description": "Role/permission: Academic | LMS.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Search"
          },
          {
            "name": "status",
            "in": "query",
            "required": false,
            "description": "status",
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/Lecturer"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Academic | LMS",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/hris/lecturers/{lecturer_id}": {
      "get": {
        "tags": [
          "HRIS"
        ],
        "operationId": "getLecturer",
        "summary": "Detail dosen aktif",
        "description": "Role/permission: Academic | LMS.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "name": "lecturer_id",
            "in": "path",
            "required": true,
            "description": "lecturer_id identifier",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Lecturer"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Academic | LMS",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/lms/classes/sync-from-academic": {
      "post": {
        "tags": [
          "LMS"
        ],
        "operationId": "syncClassFromAcademic",
        "summary": "Sync kelas akademik ke LMS",
        "description": "Role/permission: service:Academic.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/LmsSyncResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "service:Academic",
        "x-idempotent": true,
        "x-audit-log": "sync log",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/LmsClassSyncRequest"
              }
            }
          }
        },
        "security": [
          {
            "serviceTokenAuth": []
          }
        ]
      }
    },
    "/api/v1/lms/enrollments/sync-from-krs": {
      "post": {
        "tags": [
          "LMS"
        ],
        "operationId": "syncEnrollmentFromKrs",
        "summary": "Sync peserta kelas dari KRS valid",
        "description": "Role/permission: service:Academic.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/LmsSyncResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "service:Academic",
        "x-idempotent": true,
        "x-audit-log": "sync log",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/EnrollmentSyncRequest"
              }
            }
          }
        },
        "security": [
          {
            "serviceTokenAuth": []
          }
        ]
      }
    },
    "/api/v1/lms/grade-syncs": {
      "post": {
        "tags": [
          "LMS"
        ],
        "operationId": "syncLmsGrade",
        "summary": "Kirim nilai aktivitas ke Academic",
        "description": "Role/permission: dosen | service:LMS.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Grade"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "dosen | service:LMS",
        "x-idempotent": true,
        "x-audit-log": "grade sync audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/LmsGradeSyncRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/assessment/sessions": {
      "post": {
        "tags": [
          "Assessment"
        ],
        "operationId": "createAssessmentSession",
        "summary": "Membuat assessment session",
        "description": "Role/permission: PMB | LMS | admin_assessment.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/AssessmentSession"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "PMB | LMS | admin_assessment",
        "x-idempotent": true,
        "x-audit-log": "session audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/AssessmentSessionRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/assessment/attempts": {
      "post": {
        "tags": [
          "Assessment"
        ],
        "operationId": "createAssessmentAttempt",
        "summary": "Membuat attempt peserta",
        "description": "Role/permission: participant.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/AssessmentAttempt"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "participant",
        "x-idempotent": true,
        "x-audit-log": "attempt log",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/AssessmentAttemptRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/assessment/results/publish": {
      "post": {
        "tags": [
          "Assessment"
        ],
        "operationId": "publishAssessmentResult",
        "summary": "Publish hasil ke consumer",
        "description": "Role/permission: admin_assessment | service:Assessment.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/AssessmentAttempt"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "admin_assessment | service:Assessment",
        "x-idempotent": true,
        "x-audit-log": "result audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/AssessmentResultPublishRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/portal/notifications": {
      "post": {
        "tags": [
          "Portal"
        ],
        "operationId": "createPortalNotification",
        "summary": "Membuat notifikasi portal",
        "description": "Role/permission: service:any module.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/Notification"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "service:any module",
        "x-idempotent": true,
        "x-audit-log": "notification log",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/NotificationRequest"
              }
            }
          }
        },
        "security": [
          {
            "serviceTokenAuth": []
          }
        ]
      }
    },
    "/api/v1/portal/dashboard": {
      "get": {
        "tags": [
          "Portal"
        ],
        "operationId": "getRoleBasedDashboard",
        "summary": "Dashboard sesuai role",
        "description": "Role/permission: Authenticated.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/DashboardResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "Authenticated",
        "x-idempotent": false,
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/event-contracts": {
      "get": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "listEventContracts",
        "summary": "List event contract catalog",
        "description": "Role/permission: technical_admin | devops_sre | auditor.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Sort"
          },
          {
            "$ref": "#/components/parameters/Search"
          },
          {
            "$ref": "#/components/parameters/EventStatusQuery"
          },
          {
            "$ref": "#/components/parameters/EventNameQuery"
          },
          {
            "$ref": "#/components/parameters/ModuleQuery"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/EventContract"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "technical_admin | devops_sre | auditor",
        "x-idempotent": false,
        "x-audit-log": "event observability read",
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/event-contracts/{event_name}": {
      "get": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "getEventContractByName",
        "summary": "Detail event contract berdasarkan event_name",
        "description": "Role/permission: technical_admin | devops_sre | auditor.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/EventNamePath"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/EventContract"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "technical_admin | devops_sre | auditor",
        "x-idempotent": false,
        "x-audit-log": "event contract read",
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/outbox-events": {
      "get": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "listOutboxEvents",
        "summary": "List outbox events",
        "description": "Role/permission: technical_admin | devops_sre | auditor.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Sort"
          },
          {
            "$ref": "#/components/parameters/Search"
          },
          {
            "$ref": "#/components/parameters/EventStatusQuery"
          },
          {
            "$ref": "#/components/parameters/EventNameQuery"
          },
          {
            "$ref": "#/components/parameters/ModuleQuery"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/OutboxEvent"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "technical_admin | devops_sre | auditor",
        "x-idempotent": false,
        "x-audit-log": "event observability read",
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/inbox-events": {
      "get": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "listInboxEvents",
        "summary": "List inbox events",
        "description": "Role/permission: technical_admin | devops_sre | auditor.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Sort"
          },
          {
            "$ref": "#/components/parameters/Search"
          },
          {
            "$ref": "#/components/parameters/EventStatusQuery"
          },
          {
            "$ref": "#/components/parameters/EventNameQuery"
          },
          {
            "$ref": "#/components/parameters/ModuleQuery"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/InboxEvent"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "technical_admin | devops_sre | auditor",
        "x-idempotent": false,
        "x-audit-log": "event observability read",
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/dlq-events": {
      "get": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "listDlqEvents",
        "summary": "List DLQ events",
        "description": "Role/permission: technical_admin | devops_sre | auditor.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Sort"
          },
          {
            "$ref": "#/components/parameters/Search"
          },
          {
            "$ref": "#/components/parameters/EventStatusQuery"
          },
          {
            "$ref": "#/components/parameters/EventNameQuery"
          },
          {
            "$ref": "#/components/parameters/ModuleQuery"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/DlqEvent"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "technical_admin | devops_sre | auditor",
        "x-idempotent": false,
        "x-audit-log": "event observability read",
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/dlq-events/{event_key}/replay": {
      "post": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "replayDlqEvent",
        "summary": "Replay event dari DLQ",
        "description": "Role/permission: devops_sre. Replay wajib memakai reason dan audit trail.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/EventKeyPath"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/EventReplayResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "devops_sre",
        "x-idempotent": true,
        "x-audit-log": "event DLQ replay audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/EventReplayRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/reconciliation-mismatches": {
      "get": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "listReconciliationMismatches",
        "summary": "List reconciliation mismatch logs",
        "description": "Role/permission: technical_admin | devops_sre | auditor.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/Page"
          },
          {
            "$ref": "#/components/parameters/Limit"
          },
          {
            "$ref": "#/components/parameters/Sort"
          },
          {
            "$ref": "#/components/parameters/Search"
          },
          {
            "$ref": "#/components/parameters/EventStatusQuery"
          },
          {
            "$ref": "#/components/parameters/EventNameQuery"
          },
          {
            "$ref": "#/components/parameters/ModuleQuery"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/ReconciliationMismatch"
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "technical_admin | devops_sre | auditor",
        "x-idempotent": false,
        "x-audit-log": "event observability read",
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    },
    "/api/v1/integration/reconciliation-mismatches/{mismatch_id}/resolve": {
      "post": {
        "tags": [
          "Event Contract"
        ],
        "operationId": "resolveReconciliationMismatch",
        "summary": "Resolve reconciliation mismatch",
        "description": "Role/permission: technical_admin | devops_sre | owner_modul. Resolusi wajib mencatat reason.",
        "parameters": [
          {
            "$ref": "#/components/parameters/XApplicationCode"
          },
          {
            "$ref": "#/components/parameters/XActiveRole"
          },
          {
            "$ref": "#/components/parameters/XCorrelationId"
          },
          {
            "$ref": "#/components/parameters/MismatchIdPath"
          },
          {
            "$ref": "#/components/parameters/IdempotencyKey"
          }
        ],
        "responses": {
          "200": {
            "description": "Request processed successfully",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessEnvelope"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "data": {
                          "$ref": "#/components/schemas/ReconciliationResolveResponse"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "401": {
            "description": "Authentication required or token expired",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "403": {
            "description": "Forbidden or scope denied",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "404": {
            "description": "Data not found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "409": {
            "description": "Business rule violation, duplicate request, or conflict state",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "422": {
            "description": "Validation error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "500": {
            "description": "Internal server error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          },
          "502": {
            "description": "Integration failed",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorEnvelope"
                }
              }
            }
          }
        },
        "x-permission-roles": "technical_admin | devops_sre | owner_modul",
        "x-idempotent": true,
        "x-audit-log": "reconciliation resolve audit",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/ReconciliationResolveRequest"
              }
            }
          }
        },
        "security": [
          {
            "bearerAuth": []
          }
        ]
      }
    }
  },
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT",
        "description": "Access token dari Core SSO/OAuth/OIDC-style."
      },
      "serviceTokenAuth": {
        "type": "apiKey",
        "in": "header",
        "name": "X-Service-Token",
        "description": "Service token komunikasi internal antar modul."
      },
      "paymentProviderSignature": {
        "type": "apiKey",
        "in": "header",
        "name": "X-Provider-Signature",
        "description": "Signature callback dari payment provider."
      }
    },
    "parameters": {
      "XApplicationCode": {
        "name": "X-Application-Code",
        "in": "header",
        "required": true,
        "schema": {
          "type": "string"
        },
        "description": "Kode aplikasi pemanggil, misalnya portal, pmb, academic, finance, lms."
      },
      "XActiveRole": {
        "name": "X-Active-Role",
        "in": "header",
        "required": true,
        "schema": {
          "type": "string"
        },
        "description": "Role aktif dari Core active role session."
      },
      "XCorrelationId": {
        "name": "X-Correlation-Id",
        "in": "header",
        "required": true,
        "schema": {
          "type": "string"
        },
        "description": "ID korelasi untuk tracing lintas modul."
      },
      "IdempotencyKey": {
        "name": "Idempotency-Key",
        "in": "header",
        "required": true,
        "schema": {
          "type": "string"
        },
        "description": "Key deterministik untuk mencegah duplicate processing."
      },
      "Page": {
        "name": "page",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer",
          "default": 1
        }
      },
      "Limit": {
        "name": "limit",
        "in": "query",
        "required": false,
        "schema": {
          "type": "integer",
          "default": 20,
          "maximum": 100
        }
      },
      "Sort": {
        "name": "sort",
        "in": "query",
        "required": false,
        "schema": {
          "type": "string",
          "example": "created_at:desc"
        }
      },
      "Search": {
        "name": "search",
        "in": "query",
        "required": false,
        "schema": {
          "type": "string"
        }
      },
      "EventNamePath": {
        "name": "event_name",
        "in": "path",
        "required": true,
        "schema": {
          "type": "string"
        },
        "description": "Event name, contoh finance.payment_paid."
      },
      "EventKeyPath": {
        "name": "event_key",
        "in": "path",
        "required": true,
        "schema": {
          "type": "string"
        },
        "description": "Event key unik."
      },
      "MismatchIdPath": {
        "name": "mismatch_id",
        "in": "path",
        "required": true,
        "schema": {
          "type": "string",
          "format": "uuid"
        },
        "description": "ID mismatch reconciliation."
      },
      "EventStatusQuery": {
        "name": "status",
        "in": "query",
        "required": false,
        "schema": {
          "type": "string"
        },
        "description": "Filter status event."
      },
      "EventNameQuery": {
        "name": "event_name",
        "in": "query",
        "required": false,
        "schema": {
          "type": "string"
        },
        "description": "Filter nama event."
      },
      "ModuleQuery": {
        "name": "module",
        "in": "query",
        "required": false,
        "schema": {
          "type": "string"
        },
        "description": "Filter modul publisher/consumer."
      }
    },
    "schemas": {
      "Meta": {
        "type": "object",
        "properties": {
          "request_id": {
            "type": "string",
            "example": "req_20260618_000001"
          },
          "correlation_id": {
            "type": "string",
            "example": "corr-uuid"
          },
          "timestamp": {
            "type": "string",
            "format": "date-time",
            "example": "2026-06-18T10:00:00+07:00"
          }
        },
        "required": [
          "request_id",
          "correlation_id",
          "timestamp"
        ]
      },
      "PaginationMeta": {
        "type": "object",
        "properties": {
          "page": {
            "type": "integer",
            "example": 1
          },
          "limit": {
            "type": "integer",
            "example": 20
          },
          "total": {
            "type": "integer",
            "example": 250
          },
          "total_page": {
            "type": "integer",
            "example": 13
          },
          "sort": {
            "type": "string",
            "example": "created_at:desc"
          }
        }
      },
      "SuccessEnvelope": {
        "type": "object",
        "properties": {
          "success": {
            "type": "boolean",
            "example": true
          },
          "code": {
            "type": "string",
            "example": "OK"
          },
          "message": {
            "type": "string",
            "example": "Request processed successfully"
          },
          "data": {
            "type": "object",
            "additionalProperties": true
          },
          "meta": {
            "$ref": "#/components/schemas/Meta"
          }
        },
        "required": [
          "success",
          "code",
          "message",
          "data",
          "meta"
        ]
      },
      "ErrorDetail": {
        "type": "object",
        "properties": {
          "field": {
            "type": "string",
            "example": "academic_period_id"
          },
          "message": {
            "type": "string",
            "example": "academic_period_id is required"
          }
        }
      },
      "ErrorEnvelope": {
        "type": "object",
        "properties": {
          "success": {
            "type": "boolean",
            "example": false
          },
          "code": {
            "type": "string",
            "enum": [
              "AUTH_REQUIRED",
              "TOKEN_EXPIRED",
              "FORBIDDEN",
              "SCOPE_DENIED",
              "NOT_FOUND",
              "VALIDATION_ERROR",
              "BUSINESS_RULE_VIOLATION",
              "DUPLICATE_REQUEST",
              "CONFLICT_STATE",
              "INTEGRATION_FAILED",
              "INTERNAL_ERROR",
              "EVENT_SCHEMA_INVALID",
              "EVENT_DUPLICATE",
              "EVENT_VERSION_UNSUPPORTED",
              "SOURCE_REF_NOT_FOUND",
              "SNAPSHOT_UPDATE_FAILED",
              "RECONCILIATION_REQUIRED",
              "CONSUMER_TEMPORARY_FAILURE",
              "CONSUMER_PERMANENT_FAILURE",
              "DLQ_REPLAY_DENIED"
            ]
          },
          "message": {
            "type": "string",
            "example": "Request validation failed"
          },
          "errors": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/ErrorDetail"
            }
          },
          "meta": {
            "$ref": "#/components/schemas/Meta"
          }
        },
        "required": [
          "success",
          "code",
          "message",
          "meta"
        ]
      },
      "TokenResponse": {
        "type": "object",
        "properties": {
          "access_token": {
            "type": "string"
          },
          "refresh_token": {
            "type": "string"
          },
          "expires_in": {
            "type": "integer",
            "example": 3600
          },
          "token_type": {
            "type": "string",
            "example": "Bearer"
          }
        }
      },
      "UserProfile": {
        "type": "object",
        "properties": {
          "user_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "person_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "name": {
            "type": "string"
          },
          "email": {
            "type": "string"
          },
          "active_role": {
            "type": "string"
          },
          "permissions": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "scope": {
            "type": "object",
            "additionalProperties": true
          }
        }
      },
      "Application": {
        "type": "object",
        "properties": {
          "application_code": {
            "type": "string"
          },
          "name": {
            "type": "string"
          },
          "url": {
            "type": "string"
          },
          "enabled": {
            "type": "boolean"
          }
        }
      },
      "StudyProgram": {
        "type": "object",
        "properties": {
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "code": {
            "type": "string"
          },
          "name": {
            "type": "string"
          },
          "degree": {
            "type": "string"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "AcademicYear": {
        "type": "object",
        "properties": {
          "academic_year_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "code": {
            "type": "string",
            "example": "2026/2027"
          },
          "name": {
            "type": "string"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "AcademicPeriod": {
        "type": "object",
        "properties": {
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_year_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "code": {
            "type": "string",
            "example": "2026-GANJIL"
          },
          "term": {
            "type": "string"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "Lead": {
        "type": "object",
        "properties": {
          "lead_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "name": {
            "type": "string"
          },
          "email": {
            "type": "string"
          },
          "phone": {
            "type": "string"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "ConvertLeadResponse": {
        "type": "object",
        "properties": {
          "lead_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "applicant_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "applicant_number": {
            "type": "string",
            "example": "PMB-2026-000001"
          },
          "status": {
            "type": "string",
            "example": "draft"
          }
        }
      },
      "Applicant": {
        "type": "object",
        "properties": {
          "applicant_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "applicant_number": {
            "type": "string"
          },
          "status": {
            "type": "string"
          },
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          }
        }
      },
      "Invoice": {
        "type": "object",
        "properties": {
          "invoice_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "invoice_number": {
            "type": "string"
          },
          "payer_type": {
            "type": "string"
          },
          "payer_ref_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "amount_total": {
            "type": "number"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "Student": {
        "type": "object",
        "properties": {
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "nim": {
            "type": "string"
          },
          "person_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "entry_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "AcademicClass": {
        "type": "object",
        "properties": {
          "academic_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "course_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "class_code": {
            "type": "string"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "Krs": {
        "type": "object",
        "properties": {
          "krs_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "Grade": {
        "type": "object",
        "properties": {
          "grade_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "final_grade": {
            "type": "number"
          },
          "letter_grade": {
            "type": "string"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "Lecturer": {
        "type": "object",
        "properties": {
          "lecturer_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "person_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "name": {
            "type": "string"
          },
          "homebase_study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "employment_status": {
            "type": "string"
          }
        }
      },
      "LmsSyncResponse": {
        "type": "object",
        "properties": {
          "sync_status": {
            "type": "string",
            "example": "upserted"
          },
          "lms_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "lms_enrollment_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          }
        }
      },
      "AssessmentSession": {
        "type": "object",
        "properties": {
          "assessment_session_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "context": {
            "type": "string"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "AssessmentAttempt": {
        "type": "object",
        "properties": {
          "attempt_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "assessment_session_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "status": {
            "type": "string"
          }
        }
      },
      "Notification": {
        "type": "object",
        "properties": {
          "notification_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "recipient_user_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "event_type": {
            "type": "string"
          },
          "read_status": {
            "type": "boolean"
          }
        }
      },
      "DashboardResponse": {
        "type": "object",
        "properties": {
          "role": {
            "type": "string"
          },
          "widgets": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "widget_code": {
                  "type": "string"
                },
                "title": {
                  "type": "string"
                },
                "summary": {
                  "type": "object",
                  "additionalProperties": true
                }
              }
            }
          }
        }
      },
      "ClearanceResponse": {
        "type": "object",
        "properties": {
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "applicant_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "service_code": {
            "type": "string"
          },
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "clearance_status": {
            "type": "string",
            "enum": [
              "clear",
              "conditional",
              "blocked"
            ]
          },
          "block_reasons": {
            "type": "array",
            "items": {
              "type": "string"
            }
          }
        }
      },
      "LoginRequest": {
        "type": "object",
        "properties": {
          "username": {
            "type": "string",
            "example": "admin@unsia.ac.id"
          },
          "password": {
            "type": "string",
            "format": "password"
          },
          "captcha_token": {
            "type": "string",
            "nullable": true
          }
        },
        "required": [
          "username",
          "password"
        ]
      },
      "SwitchRoleRequest": {
        "type": "object",
        "properties": {
          "role_code": {
            "type": "string",
            "example": "admin_pmb"
          },
          "scope_value": {
            "type": "string",
            "nullable": true
          }
        },
        "required": [
          "role_code"
        ]
      },
      "ImpersonationStartRequest": {
        "type": "object",
        "properties": {
          "target_user_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "reason": {
            "type": "string"
          },
          "duration_minutes": {
            "type": "integer",
            "default": 30
          }
        },
        "required": [
          "target_user_id",
          "reason"
        ]
      },
      "AcademicYearCreateRequest": {
        "type": "object",
        "properties": {
          "code": {
            "type": "string",
            "example": "2026/2027"
          },
          "name": {
            "type": "string"
          },
          "start_date": {
            "type": "string",
            "format": "date"
          },
          "end_date": {
            "type": "string",
            "format": "date"
          },
          "status": {
            "type": "string",
            "default": "active"
          }
        },
        "required": [
          "code",
          "name"
        ]
      },
      "AcademicPeriodCreateRequest": {
        "type": "object",
        "properties": {
          "academic_year_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "code": {
            "type": "string",
            "example": "2026-GANJIL"
          },
          "term": {
            "type": "string",
            "enum": [
              "GANJIL",
              "GENAP",
              "PENDEK"
            ]
          },
          "start_date": {
            "type": "string",
            "format": "date"
          },
          "end_date": {
            "type": "string",
            "format": "date"
          }
        },
        "required": [
          "academic_year_id",
          "code",
          "term"
        ]
      },
      "LeadCreateRequest": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string"
          },
          "email": {
            "type": "string"
          },
          "phone": {
            "type": "string"
          },
          "campaign_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "source": {
            "type": "string"
          },
          "target_entry_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          }
        },
        "required": [
          "name",
          "phone"
        ]
      },
      "ConvertLeadRequest": {
        "type": "object",
        "properties": {
          "target_entry_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "admission_path_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "reason": {
            "type": "string"
          }
        },
        "required": [
          "target_entry_period_id",
          "admission_path_id",
          "study_program_id",
          "reason"
        ]
      },
      "ApplicantCreateRequest": {
        "type": "object",
        "properties": {
          "person": {
            "type": "object",
            "properties": {
              "name": {
                "type": "string"
              },
              "email": {
                "type": "string"
              },
              "phone": {
                "type": "string"
              }
            }
          },
          "target_entry_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "admission_path_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          }
        },
        "required": [
          "person",
          "target_entry_period_id",
          "study_program_id"
        ]
      },
      "SubmitRequest": {
        "type": "object",
        "properties": {
          "reason": {
            "type": "string"
          }
        }
      },
      "DocumentUploadRequest": {
        "type": "object",
        "properties": {
          "document_type_code": {
            "type": "string",
            "example": "KTP"
          },
          "file": {
            "type": "string",
            "format": "binary"
          }
        },
        "required": [
          "document_type_code",
          "file"
        ]
      },
      "VerifyDocumentRequest": {
        "type": "object",
        "properties": {
          "verification_status": {
            "type": "string",
            "enum": [
              "verified",
              "rejected"
            ]
          },
          "reason": {
            "type": "string"
          }
        },
        "required": [
          "verification_status"
        ]
      },
      "RequestInvoiceRequest": {
        "type": "object",
        "properties": {
          "invoice_context": {
            "type": "string",
            "enum": [
              "PMB_FORM",
              "RE_REGISTRATION"
            ]
          },
          "payment_component_codes": {
            "type": "array",
            "items": {
              "type": "string"
            }
          },
          "due_date": {
            "type": "string",
            "format": "date"
          },
          "reason": {
            "type": "string"
          }
        },
        "required": [
          "invoice_context",
          "payment_component_codes",
          "due_date"
        ]
      },
      "IssueLoaRequest": {
        "type": "object",
        "properties": {
          "reason": {
            "type": "string"
          },
          "template_code": {
            "type": "string"
          }
        },
        "required": [
          "reason"
        ]
      },
      "HandoverToAcademicRequest": {
        "type": "object",
        "properties": {
          "entry_academic_year_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "entry_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "curriculum_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "reason": {
            "type": "string"
          }
        },
        "required": [
          "entry_academic_year_id",
          "entry_period_id",
          "study_program_id",
          "curriculum_id",
          "reason"
        ]
      },
      "HandoverToAcademicResponse": {
        "type": "object",
        "properties": {
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "nim": {
            "type": "string",
            "example": "20261010001"
          },
          "handover_status": {
            "type": "string",
            "example": "completed"
          }
        }
      },
      "InvoiceCreateRequest": {
        "type": "object",
        "properties": {
          "payer_type": {
            "type": "string",
            "enum": [
              "applicant",
              "student"
            ]
          },
          "payer_ref_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "items": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "payment_component_code": {
                  "type": "string"
                },
                "amount": {
                  "type": "number"
                }
              }
            }
          },
          "source_module": {
            "type": "string"
          },
          "source_ref_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "due_date": {
            "type": "string",
            "format": "date"
          }
        },
        "required": [
          "payer_type",
          "payer_ref_id",
          "academic_period_id",
          "items",
          "source_module",
          "source_ref_id"
        ]
      },
      "PaymentCallbackRequest": {
        "type": "object",
        "properties": {
          "provider_event_id": {
            "type": "string"
          },
          "invoice_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "payment_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "amount": {
            "type": "number"
          },
          "payment_status": {
            "type": "string",
            "enum": [
              "paid",
              "partial",
              "failed",
              "expired"
            ]
          },
          "signature_status": {
            "type": "string",
            "enum": [
              "valid",
              "invalid"
            ]
          },
          "payload_hash": {
            "type": "string"
          },
          "raw_payload": {
            "type": "object",
            "additionalProperties": true
          }
        },
        "required": [
          "provider_event_id",
          "invoice_id",
          "amount",
          "payment_status",
          "payload_hash"
        ]
      },
      "PaymentVerificationRequest": {
        "type": "object",
        "properties": {
          "invoice_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "payment_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "verification_status": {
            "type": "string",
            "enum": [
              "verified",
              "rejected"
            ]
          },
          "amount": {
            "type": "number"
          },
          "paid_at": {
            "type": "string",
            "format": "date-time"
          },
          "attachment_ref": {
            "type": "string"
          },
          "reason": {
            "type": "string"
          }
        },
        "required": [
          "invoice_id",
          "verification_status",
          "reason"
        ]
      },
      "GenerateStudentFromApplicantRequest": {
        "type": "object",
        "properties": {
          "applicant_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "curriculum_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "entry_academic_year_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "entry_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "study_program_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "reason": {
            "type": "string"
          }
        },
        "required": [
          "applicant_id",
          "curriculum_id",
          "entry_period_id"
        ]
      },
      "ClassCreateRequest": {
        "type": "object",
        "properties": {
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "course_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "class_code": {
            "type": "string"
          },
          "lecturer_ids": {
            "type": "array",
            "items": {
              "type": "string",
              "format": "uuid",
              "example": "00000000-0000-0000-0000-000000000000"
            }
          },
          "capacity": {
            "type": "integer"
          }
        },
        "required": [
          "academic_period_id",
          "course_id",
          "class_code"
        ]
      },
      "KrsCreateRequest": {
        "type": "object",
        "properties": {
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "items": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "academic_class_id": {
                  "type": "string",
                  "format": "uuid",
                  "example": "00000000-0000-0000-0000-000000000000"
                }
              }
            }
          }
        },
        "required": [
          "student_id",
          "academic_period_id"
        ]
      },
      "GradeSourceImportRequest": {
        "type": "object",
        "properties": {
          "source_module": {
            "type": "string",
            "enum": [
              "lms",
              "assessment"
            ]
          },
          "source_ref_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "component_code": {
            "type": "string",
            "enum": [
              "QUIZ",
              "ASSIGNMENT",
              "CBT"
            ]
          },
          "score": {
            "type": "number"
          },
          "max_score": {
            "type": "number"
          },
          "submitted_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "source_module",
          "source_ref_id",
          "academic_class_id",
          "student_id",
          "component_code",
          "score",
          "max_score"
        ]
      },
      "GradeFinalizeRequest": {
        "type": "object",
        "properties": {
          "final_grade": {
            "type": "number"
          },
          "letter_grade": {
            "type": "string"
          },
          "reason": {
            "type": "string"
          }
        },
        "required": [
          "reason"
        ]
      },
      "LmsClassSyncRequest": {
        "type": "object",
        "properties": {
          "academic_period_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "course_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "class_code": {
            "type": "string"
          },
          "lecturer_ids": {
            "type": "array",
            "items": {
              "type": "string",
              "format": "uuid",
              "example": "00000000-0000-0000-0000-000000000000"
            }
          },
          "schedule": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "day": {
                  "type": "string"
                },
                "start_time": {
                  "type": "string"
                },
                "end_time": {
                  "type": "string"
                }
              }
            }
          }
        },
        "required": [
          "academic_period_id",
          "academic_class_id",
          "course_id",
          "class_code"
        ]
      },
      "EnrollmentSyncRequest": {
        "type": "object",
        "properties": {
          "academic_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "krs_item_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "enrollment_status": {
            "type": "string",
            "default": "active"
          }
        },
        "required": [
          "academic_class_id",
          "krs_item_id",
          "student_id"
        ]
      },
      "LmsGradeSyncRequest": {
        "type": "object",
        "properties": {
          "source_ref_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "student_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "component_code": {
            "type": "string"
          },
          "score": {
            "type": "number"
          },
          "max_score": {
            "type": "number"
          },
          "submitted_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "source_ref_id",
          "academic_class_id",
          "student_id",
          "component_code",
          "score",
          "max_score"
        ]
      },
      "AssessmentSessionRequest": {
        "type": "object",
        "properties": {
          "context": {
            "type": "string",
            "enum": [
              "PMB",
              "LMS",
              "ACADEMIC",
              "SURVEY"
            ]
          },
          "title": {
            "type": "string"
          },
          "consumer_module": {
            "type": "string"
          },
          "wave_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "academic_class_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "question_set_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "start_at": {
            "type": "string",
            "format": "date-time"
          },
          "end_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "context",
          "title"
        ]
      },
      "AssessmentAttemptRequest": {
        "type": "object",
        "properties": {
          "assessment_session_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "participant_ref_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "started_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "assessment_session_id",
          "participant_ref_id"
        ]
      },
      "AssessmentResultPublishRequest": {
        "type": "object",
        "properties": {
          "attempt_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "consumer_module": {
            "type": "string",
            "enum": [
              "pmb",
              "lms",
              "academic"
            ]
          },
          "score": {
            "type": "number"
          },
          "result_status": {
            "type": "string"
          },
          "target_ref_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          }
        },
        "required": [
          "attempt_id",
          "consumer_module",
          "score",
          "result_status"
        ]
      },
      "NotificationRequest": {
        "type": "object",
        "properties": {
          "recipient_user_id": {
            "type": "string",
            "format": "uuid",
            "example": "00000000-0000-0000-0000-000000000000"
          },
          "event_type": {
            "type": "string"
          },
          "source_module": {
            "type": "string"
          },
          "event_key": {
            "type": "string"
          },
          "severity": {
            "type": "string",
            "enum": [
              "info",
              "warning",
              "urgent"
            ]
          },
          "title": {
            "type": "string"
          },
          "message": {
            "type": "string"
          },
          "link_url": {
            "type": "string"
          },
          "expires_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "recipient_user_id",
          "event_type",
          "source_module",
          "event_key",
          "message"
        ]
      },
      "EventContract": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid"
          },
          "event_name": {
            "type": "string",
            "example": "finance.payment_paid"
          },
          "event_version": {
            "type": "string",
            "example": "v1"
          },
          "event_type": {
            "type": "string",
            "enum": [
              "DOMAIN_EVENT",
              "INTEGRATION_EVENT",
              "NOTIFICATION_EVENT",
              "SNAPSHOT_EVENT"
            ]
          },
          "publisher_module": {
            "type": "string",
            "example": "finance"
          },
          "publisher_database": {
            "type": "string",
            "example": "finance_db"
          },
          "aggregate_type": {
            "type": "string",
            "example": "payment"
          },
          "payload_schema": {
            "type": "object",
            "additionalProperties": true
          },
          "validation_schema": {
            "type": "object",
            "additionalProperties": true
          },
          "status": {
            "type": "string",
            "enum": [
              "draft",
              "active",
              "deprecated",
              "retired"
            ]
          },
          "description": {
            "type": "string"
          },
          "created_at": {
            "type": "string",
            "format": "date-time"
          },
          "updated_at": {
            "type": "string",
            "format": "date-time"
          }
        },
        "required": [
          "event_name",
          "event_version",
          "event_type",
          "publisher_module",
          "aggregate_type",
          "payload_schema",
          "status"
        ]
      },
      "EventEnvelope": {
        "type": "object",
        "properties": {
          "event_name": {
            "type": "string",
            "example": "finance.payment_paid"
          },
          "event_version": {
            "type": "string",
            "example": "v1"
          },
          "event_key": {
            "type": "string",
            "example": "finance.payment_paid:payment_id:8f2c:v1"
          },
          "event_type": {
            "type": "string",
            "enum": [
              "DOMAIN_EVENT",
              "INTEGRATION_EVENT",
              "NOTIFICATION_EVENT",
              "SNAPSHOT_EVENT"
            ]
          },
          "publisher_service": {
            "type": "string",
            "example": "finance-service"
          },
          "publisher_database": {
            "type": "string",
            "example": "finance_db"
          },
          "aggregate_type": {
            "type": "string",
            "example": "payment"
          },
          "aggregate_id": {
            "type": "string",
            "format": "uuid"
          },
          "correlation_id": {
            "type": "string"
          },
          "causation_id": {
            "type": "string"
          },
          "occurred_at": {
            "type": "string",
            "format": "date-time"
          },
          "published_at": {
            "type": "string",
            "format": "date-time"
          },
          "payload": {
            "type": "object",
            "additionalProperties": true
          }
        },
        "required": [
          "event_name",
          "event_version",
          "event_key",
          "publisher_service",
          "aggregate_type",
          "aggregate_id",
          "occurred_at",
          "payload"
        ]
      },
      "OutboxEvent": {
        "allOf": [
          {
            "$ref": "#/components/schemas/EventEnvelope"
          },
          {
            "type": "object",
            "properties": {
              "id": {
                "type": "string",
                "format": "uuid"
              },
              "status": {
                "type": "string",
                "enum": [
                  "PENDING",
                  "PUBLISHED",
                  "RETRYING",
                  "FAILED",
                  "DLQ"
                ]
              },
              "retry_count": {
                "type": "integer"
              },
              "max_retry": {
                "type": "integer"
              },
              "next_retry_at": {
                "type": "string",
                "format": "date-time",
                "nullable": true
              },
              "last_error": {
                "type": "string",
                "nullable": true
              },
              "dead_letter_at": {
                "type": "string",
                "format": "date-time",
                "nullable": true
              },
              "created_at": {
                "type": "string",
                "format": "date-time"
              },
              "updated_at": {
                "type": "string",
                "format": "date-time"
              }
            }
          }
        ]
      },
      "InboxEvent": {
        "allOf": [
          {
            "$ref": "#/components/schemas/EventEnvelope"
          },
          {
            "type": "object",
            "properties": {
              "id": {
                "type": "string",
                "format": "uuid"
              },
              "consumer_module": {
                "type": "string",
                "example": "pmb"
              },
              "status": {
                "type": "string",
                "enum": [
                  "RECEIVED",
                  "PROCESSED",
                  "RETRYING",
                  "FAILED",
                  "DLQ",
                  "IGNORED_DUPLICATE"
                ]
              },
              "retry_count": {
                "type": "integer"
              },
              "next_retry_at": {
                "type": "string",
                "format": "date-time",
                "nullable": true
              },
              "received_at": {
                "type": "string",
                "format": "date-time"
              },
              "processed_at": {
                "type": "string",
                "format": "date-time",
                "nullable": true
              },
              "last_error": {
                "type": "string",
                "nullable": true
              },
              "dead_letter_at": {
                "type": "string",
                "format": "date-time",
                "nullable": true
              }
            }
          }
        ]
      },
      "DlqEvent": {
        "type": "object",
        "properties": {
          "event_key": {
            "type": "string"
          },
          "event_name": {
            "type": "string"
          },
          "event_version": {
            "type": "string"
          },
          "publisher_module": {
            "type": "string"
          },
          "consumer_module": {
            "type": "string"
          },
          "retry_count": {
            "type": "integer"
          },
          "last_error": {
            "type": "string"
          },
          "dead_letter_at": {
            "type": "string",
            "format": "date-time"
          },
          "payload": {
            "type": "object",
            "additionalProperties": true
          }
        },
        "required": [
          "event_key",
          "event_name",
          "publisher_module",
          "consumer_module",
          "last_error",
          "dead_letter_at"
        ]
      },
      "EventReplayRequest": {
        "type": "object",
        "properties": {
          "reason": {
            "type": "string",
            "example": "Replay after consumer fix"
          },
          "target_consumer": {
            "type": "string",
            "example": "pmb"
          },
          "force": {
            "type": "boolean",
            "default": false
          }
        },
        "required": [
          "reason"
        ]
      },
      "EventReplayResponse": {
        "type": "object",
        "properties": {
          "event_key": {
            "type": "string"
          },
          "replay_status": {
            "type": "string",
            "enum": [
              "queued",
              "success",
              "failed"
            ]
          },
          "audit_ref_id": {
            "type": "string",
            "format": "uuid"
          },
          "message": {
            "type": "string"
          }
        }
      },
      "ReconciliationMismatch": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid"
          },
          "source_module": {
            "type": "string"
          },
          "source_table": {
            "type": "string"
          },
          "source_ref_id": {
            "type": "string",
            "format": "uuid"
          },
          "consumer_module": {
            "type": "string"
          },
          "consumer_table": {
            "type": "string"
          },
          "consumer_ref_id": {
            "type": "string",
            "format": "uuid"
          },
          "source_event_key": {
            "type": "string"
          },
          "mismatch_type": {
            "type": "string",
            "enum": [
              "missing_source",
              "missing_snapshot",
              "value_mismatch",
              "stale_snapshot",
              "duplicate_projection"
            ]
          },
          "source_value": {
            "type": "object",
            "additionalProperties": true
          },
          "snapshot_value": {
            "type": "object",
            "additionalProperties": true
          },
          "status": {
            "type": "string",
            "enum": [
              "OPEN",
              "CORRECTED",
              "IGNORED",
              "PENDING_REVIEW"
            ]
          },
          "reason": {
            "type": "string",
            "nullable": true
          },
          "detected_at": {
            "type": "string",
            "format": "date-time"
          },
          "corrected_at": {
            "type": "string",
            "format": "date-time",
            "nullable": true
          },
          "ignored_at": {
            "type": "string",
            "format": "date-time",
            "nullable": true
          }
        }
      },
      "ReconciliationResolveRequest": {
        "type": "object",
        "properties": {
          "action": {
            "type": "string",
            "enum": [
              "mark_corrected",
              "mark_ignored",
              "set_pending_review"
            ]
          },
          "reason": {
            "type": "string"
          },
          "correction_payload": {
            "type": "object",
            "additionalProperties": true
          }
        },
        "required": [
          "action",
          "reason"
        ]
      },
      "ReconciliationResolveResponse": {
        "type": "object",
        "properties": {
          "mismatch_id": {
            "type": "string",
            "format": "uuid"
          },
          "status": {
            "type": "string",
            "enum": [
              "CORRECTED",
              "IGNORED",
              "PENDING_REVIEW"
            ]
          },
          "audit_ref_id": {
            "type": "string",
            "format": "uuid"
          }
        }
      }
    }
  },
  "x-api-governance": {
    "api_versioning": "All endpoints use /api/v1. Breaking changes must use /api/v2 or versioned contract.",
    "backend_scope": "Permission and data scope must be enforced in backend, not only UI.",
    "audit": "Sensitive actions log actor, active role, application, request_id, reason, IP, user agent, old value, and new value when relevant.",
    "idempotency": "Payment callback, handover, NIM generation, class sync, enrollment sync, and grade sync must use deterministic idempotency key.",
    "domain_ownership": "Module only writes to its own domain. Cross-domain write uses API/service event."
  }
}
```