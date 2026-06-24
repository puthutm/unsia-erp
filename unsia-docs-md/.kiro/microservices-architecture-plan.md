# Microservices Architecture Plan - ALL MODULES

## 📋 Overview

Rencana arsitektur microservices untuk proyek UNSIA ERP. Setiap menu/domain akan menjadi service terpisah dengan database sendiri dan berkomunikasi via RabbitMQ (event-driven).

---

## 📊 Current State vs Target

### Current State (Monolith per Modul)
```
services/
├── unsia-academic-service    (semua menu akademi)
├── unsia-assessment-service
├── unsia-core-service    (SSO, auth, users)
├── unsia-crm-service    (semua menu CRM)
├── unsia-finance-service (semua menu keuangan)
├── unsia-hris-service   (semua menu SDM/HRIS)
├── unsia-pmb-service   (PMB)
├── unsia-lms-service
├── unsia-portal-service
└── unsia-reference-service
```

### Target State (Microservices per Menu)
```
services/
├── unsia-sso-auth-service/         # SSO & Authentication (:8001)
├── unsia-sso-session-service/  # Session Management (:8002)
├── unsia-sso-token-service/   # Token/JWT (:8003)
│
├── unsia-academic-period-service/    # Tahun Akademik (:8101)
├── unsia-academic-faculty-service/  # Fakultas (:8102)
├── unsia-academic-major-service/   # Jurusan (:8103)
├── unsia-academic-student-service/ # Mahasiswa (:8104)
├── unsia-academic-registration-service/ # KRS (:8105)
├── unsia-acadianic-grade-service/  # Nilai (:8106)
├── unsia-academic-transcript-service/ # Transkrip (:8107)
├── unsia-academic-graduation-service/ # Kelulusan (:8108)
│
├── unsia-pmb-registration-service/ # Daftar (:8201)
├── unsia-pmb-verification-service/  # Verifikasi (:8202)
├── unsia-pmb-selection-service/   # Seleksi (:8203)
├── unsia-pmb-payment-service/ # Pembayaran PMB (:8204)
├── unsia-pmb-announcement-service/ # Pengumuman (:8205)
│
├── unsia-assessment-quiz-service/      # Kuis (:8301)
├── unsia-assessment-assignment-service/ # Tugas (:8302)
├── unsia-assessment-exam-service/     # Ujian (:8303)
├── unsia-assessment-result-service/ # Hasil (:8304)
│
├── unsia-lms-course-service/     # Mata Kuliah (:8401)
├── unsia-lms-material-service/ # Materi (:8402)
├── unsia-lms-discussion-service/ # Diskusi (:8403)
├── unsia-lms-attendance-service/ # Absensi (:8404)
├── unsia-lms-certificate-service/ # Sertifikat (:8405)
│
├── unsia-hris-employee-service/  # Karyawan (:8501)
├── unsia-hris-presence-service/ # Absensi (:8502)
├── unsia-hris-payroll-service/ # Gaji (:8503)
├── unsia-hris-recruitment-service/ # Rekrutmen (:8504)
├── unsia-hris-training-service/ # Training (:8505)
├── unsia-hris-performance-service/ # Kinerja (:8506)
│
├── unsia-crm-lead-service/   # Lead (:8601)
├── unsia-crm-opportunity-service/ # Opportunity (:8602)
├── unsia-crm-customer-service/ # Customer (:8603)
├── unsia-crm-commission-service/ # Komisi (:8604)
├── unsia-crm-activity-service/ # Aktivitas (:8605)
│
├── unsia-invoice-service/    # Invoice (:8701)
├── unsia-payment-service/   # Pembayaran (:8702)
├── unsia-clearance-service/ # Clearance (:8703)
├── unsia-scholarship-service/ # Beasiswa (:8704)
├── unsia-budget-service/  # RAB (:8705)
├── unsia-cashbook-service/ # Kas/Bank (:8706)
├── unsia-journal-service/ # Jurnal (:8707)
├── unsia-payroll-service/ # Payroll (:8708)
├── unsia-disbursement-service/ # Pencairan (:8709)
├── unsia-report-service/ # Laporan (:8710)
│
├── unsia-reference-region-service/  # Wilayah (:8801)
├── unsia-reference-religion-service/ # Agama (:8802)
├── unsia-reference-major-service/ # Jurusan (:8803)
├── unsia-reference-degree-service/ # Gelaran (:8804)
└── unsia-reference-general-service/ # General (:8805)
```

---

## 🎯 Detail per Kategori

### 1. SSO/AUTH MODULE (3 Services)
**Base**: `unsia-core-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-sso-auth-service | :8001 | Login, password,认证 | `sso_auth_db` |
| unsia-sso-session-service | :8002 | Session management | `sso_session_db` |
| unsia-sso-token-service | :8003 | JWT, token refresh | `sso_token_db` |

**Endpoints:**
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/logout`
- `POST /api/v1/auth/refresh-token`
- `POST /api/v1/auth/change-password`
- `GET /api/v1/auth/me`

---

### 2. ACADEMIC MODULE (8 Services)
**Base**: `unsia-academic-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-academic-period-service | :8101 | Tahun akademik | `academic_period_db` |
| unsia-academic-faculty-service | :8102 | Fakultas | `academic_faculty_db` |
| unsia-academic-major-service | :8103 | Jurusan/Prodi | `academic_major_db` |
| unsia-academic-student-service | :8104 | Data mahasiswa | `academic_student_db` |
| unsia-academic-registration-service | :8105 | KRS/KRS Online | `academic_registration_db` |
| unsia-academic-grade-service | :8106 | Nilai & IP | `academic_grade_db` |
| unsia-academic-transcript-service | :8107 | Transkrip | `academic_transcript_db` |
| unsia-academic-graduation-service | :8108 | Kelulusan | `academic_graduation_db` |

---

### 3. PMB MODULE (5 Services)
**Base**: `unsia-pmb-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-pmb-registration-service | :8201 | Pendaftaran | `pmb_registration_db` |
| unsia-pmb-verification-service | :8202 | Verifikasi data | `pmb_verification_db` |
| unsia-pmb-selection-service | :8203 | Seleksi masuk | `pmb_selection_db` |
| unsia-pmb-payment-service | :8204 | Pembayaran daftar | `pmb_payment_db` |
| unsia-pmb-announcement-service | :8205 | Pengumuman hasil | `pmb_announcement_db` |

---

### 4. ASSESSMENT MODULE (4 Services)
**Base**: `unsia-assessment-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-assessment-quiz-service | :8301 | Kuis | `assessment_quiz_db` |
| unsia-assessment-assignment-service | :8302 | Tugas | `assessment_assignment_db` |
| unsia-assessment-exam-service | :8303 | Ujian | `assessment_exam_db` |
| unsia-assessment-result-service | :8304 | Hasil tes | `assessment_result_db` |

---

### 5. LMS MODULE (5 Services)
**Base**: `unsia-lms-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-lms-course-service | :8401 | Mata kuliah | `lms_course_db` |
| unsia-lms-material-service | :8402 | Materi belajar | `lms_material_db` |
| unsia-lms-discussion-service | :8403 | Forum diskusi | `lms_discussion_db` |
| unsia-lms-attendance-service | :8404 | Kehadiran | `lms_attendance_db` |
| unsia-lms-certificate-service | :8405 | Sertifikat | `lms_certificate_db` |

---

### 6. HRIS/SDM MODULE (6 Services)
**Base**: `unsia-hris-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-hris-employee-service | :8501 | Data karyawan | `hris_employee_db` |
| unsia-hris-presence-service | :8502 | Absensi | `hris_presence_db` |
| unsia-hris-payroll-service | :8503 | Gaji & slip | `hris_payroll_db` |
| unsia-hris-recruitment-service | :8504 | Rekrutmen | `hris_recruitment_db` |
| unsia-hris-training-service | :8505 | Training | `hris_training_db` |
| unsia-hris-performance-service | :8506 | Kinerja | `hris_performance_db` |

---

### 7. CRM MODULE (5 Services)
**Base**: `unsia-crm-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-crm-lead-service | :8601 | Lead/Prospek | `crm_lead_db` |
| unsia-crm-opportunity-service | :8602 | Opportunity | `crm_opportunity_db` |
| unsia-crm-customer-service | :8603 | Customer | `crm_customer_db` |
| unsia-crm-commission-service | :8604 | Komisiesel | `crm_commission_db` |
| unsia-crm-activity-service | :8605 | Aktivitas | `crm_activity_db` |

---

### 8. FINANCE MODULE (10 Services)
**Base**: `unsia-finance-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-invoice-service | :8701 | Invoice | `finance_invoice_db` |
| unsia-payment-service | :8702 | Pembayaran | `finance_payment_db` |
| unsia-clearance-service | :8703 | Clearance | `finance_clearance_db` |
| unsia-scholarship-service | :8704 | Beasiswa | `finance_scholarship_db` |
| unsia-budget-service | :8705 | RAB | `finance_budget_db` |
| unsia-cashbook-service | :8706 | Kas/Bank | `finance_cashbook_db` |
| unsia-journal-service | :8707 | Jurnal | `finance_journal_db` |
| unsia-payroll-service | :8708 | Payroll | `finance_payroll_db` |
| unsia-disbursement-service | :8709 | Pencairan | `finance_disbursement_db` |
| unsia-report-service | :8710 | Laporan | (read-only) |

---

### 9. REFERENCE MODULE (5 Services)
**Base**: `unsia-reference-service`

| Service | Port | Deskripsi | Database |
|---------|------|----------|----------|
| unsia-reference-region-service | :8801 | Wilayah/Provinsi | `reference_region_db` |
| unsia-reference-religion-service | :8802 | Agama | `reference_religion_db` |
| unsia-reference-major-service | :8803 | Jurusan SMU | `reference_major_db` |
| unsia-reference-degree-service | :8804 | Gelaran | `reference_degree_db` |
| unsia-reference-general-service | :8805 | General | `reference_general_db` |

---

## 📊 Rekapitulasi

| Modul | Jumlah Service | Port Range |
|------|---------------|------------|
| SSO/Auth | 3 | 8001-8003 |
| Academic | 8 | 8101-8108 |
| PMB | 5 | 8201-8205 |
| Assessment | 4 | 8301-8304 |
| LMS | 5 | 8401-8405 |
| HRIS/SDM | 6 | 8501-8506 |
| CRM | 5 | 8601-8605 |
| Finance | 10 | 8701-8710 |
| Reference | 5 | 8801-8805 |
| **TOTAL** | **51 Services** | |

---

## 🔗 Event Flow

### Student Registration Flow
```
unsia-pmb-registration (8201)
    │
    └── event: applicant.registered ──▶ unsia-pmb-verification (8202)
                                            │
                                            └── event: applicant.verified ──▶ unsia-academic-student (8104)
                                                                            │
                                                                            └── event: student.created ──▶ unsia-clearance (8703)
                                                                                                     │
                                                                                                     └── event: clearance.created ──▶ unsia-invoice (8701)
```

### Payment Flow
```
unsia-invoice (8701)
    │
    └── event: invoice.paid ──▶ unsia-payment (8702)
                                   │
                                   └── event: payment.completed ──▶ unsia-clearance (8703)
                                                                         │
                                                                         └── event: clearance.cleared ──▶ unsia-cashbook (8706)
                                                                                                        │
                                                                                                        └── event: cash.updated ──▶ unsia-journal (8707)
```

---

## 📁 Struktur Folder Services

```
services/
├── unsia-sso-auth-service/
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── domain/models.go
│   │   ├── handler/auth_handler.go
│   │   ├── service/auth_service.go
│   │   └── middleware/auth_middleware.go
│   ├── migrations/
│   ├── .env
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── unsia-sso-session-service/
│   └── ...
│
├── unsia-academic-period-service/
│   └── ...
│
├── (dan seterusnya untuk 51 services)
```

---

## ✅ Steps to Implement

### Phase 1: SSO/Auth (Week 1-2)
- [ ] unsia-sso-auth-service
- [ ] unsia-sso-session-service
- [ ] unsia-sso-token-service

### Phase 2: Academic + PMB (Week 3-4)
- [ ] 8 academic services
- [ ] 5 pmb services

### Phase 3: Finance (Week 5-6)
- [ ] 10 finance services

### Phase 4: HRIS + CRM (Week 7-8)
- [ ] 6 hris services
- [ ] 5 crm services

### Phase 5: LMS + Assessment (Week 9)
- [ ] 5 lms services
- [ ] 4 assessment services

### Phase 6: Reference (Week 10)
- [ ] 5 reference services

---

## 🔧 Shared Packages (Dependencies)

Setiap service menggunakan:
- `shared-auth` - JWT validation
- `shared-rbac` - Permission check
- `shared-errorenvelope` - Response format
- `shared-audit` - Audit logging
- `shared-idempotency` - Idempotency
- `shared-event` - Event handling
- `shared-observability` - Logging & metrics

---

## 📝 Catatan

1. **Database per Service**: Setiap service memiliki schema/database sendiri
2. **Event-Driven**: Kommunikasi antar service via RabbitMQ
3. **Port Allocation**: services menggunakan port berbeda untuk kemudahan development
4. **API Gateway**: Dianjurkan menggunakan API gateway (nginx/traefik) untuk production

---

*Plan ini dibuat untuk memisahkan setiap modul monolith menjadi microservices kecil per menu/domain.*
