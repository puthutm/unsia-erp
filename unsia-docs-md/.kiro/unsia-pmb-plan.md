# PMB Service Implementation Plan

## Service Overview
**Service Name**: unsia-pmb-service (PMB - Penerimaan Mahasiswa Baru)
**Port**: 8004
**Database**: pmb_db (PostgreSQL)

## Two Portals

### 1. PMB Public Portal (pendaftar.unsia.ac.id)
- **Target Users**: Camaba (Calon Mahasiswa Baru)
- **Features**:
  - View PMB waves and admission paths
  - Register as new applicant
  - Fill biodata, education, family information
  - Upload required documents
  - View admission status
  - Download LoA (Letter of Acceptance)

### 2. PMB Admin Portal (admin-pmb.unsia.ac.id)
- **Target Users**: PMB Admin staff
- **Features**:
  - Dashboard & statistics
  - Manage applicants (view, filter, search)
  - Verify documents
  - Input selection test scores
  - Determine pass/fail results
  - Issue Letter of Acceptance (LoA)
  - Hand over to Academic Service
  - Reports & analytics
  - Configure PMB waves & requirements

## Core Features

### 1. Applicant Registration Management
- Create new applicant registration
- Manage biodata (personal information)
- Manage addresses (KTP, DOMISILI)
- Manage education background
- Manage family members
- Manage financial profile
- Manage facility profile

### 2. Document Management
- Upload required documents (KTP, IJAZAH, SKHU, PHOTO, etc.)
- Verify documents (pending → verified/rejected)
- Document type configuration

### 3. Payment Integration
- Request invoice to Finance Service
- Track payment status
- Generate payment confirmation

### 4. Selection & Admission
- Input selection test scores
- Determine pass/fail status
- Issue Letter of Acceptance (LoA)

### 5. Academic Handover
- Transfer accepted students to Academic Service
- Generate NIM for accepted students
- Track handover status

---

## Domain Models

### Main Entities

```
Applicant
├── ID (UUID)
├── PersonID (FK → core.persons.id)
├── UserID (FK → core.users.id)
├── CrmLeadID (FK → crm.leads.id)
├── StudyProgramID (FK → ref.study_programs.id)
├── PmbWaveID (FK → ref.pmb_waves.id)
├── AdmissionPathID (FK → ref.admission_paths.id)
├── TargetEntryPeriodID (FK → ref.academic_periods.id)
├── RegistrationNumber (unique)
├── Status (draft → submitted → verified → passed/failed → accepted → ready_for_academic)
├── SubmittedAt
├── AcceptedAt
├── CreatedAt
└── UpdatedAt

ApplicantBiodata
├── ApplicantID (FK)
├── FullName
├── Email
├── Phone
├── Nik (for camaba)
├── BirthPlace
├── BirthDate
├── Gender
├── ReligionID (FK)
├── MaritalStatus
├── Citizenship
├── JacketSize
├── CoreSyncStatus (pending → synced)
└── CoreSyncedAt

ApplicantAddress
├── ApplicantID (FK)
├── AddressType (KTPS, DOMISILI)
├── Street
├── ProvinceID (FK)
├── CityID (FK)
├── DistrictID (FK)
├── VillageID (FK)
├── PostalCode
└── IsSameAsKtp

ApplicantEducationBackground
├── ApplicantID (FK)
├── SchoolName
├── Major
├── GraduationYear
└── Gpa

ApplicantFamilyMember
├── ApplicantID (FK)
├── Relationship (AYAH, IBU, WALI, SAUDARA)
├── FullName
├── Occupation
└── Income

ApplicantFinancialProfile
├── ApplicantID (FK)
├── SponsorType (SISWA, ORANGTUA, KERJA, BEASISWA)
├── SponsorName
└── MonthlyIncome

ApplicantFacilityProfile
├── ApplicantID (FK)
├── FacilityType (LAPTOP, INTERNET)
└── Description

ApplicantDocument
├── ApplicantID (FK)
├── DocumentTypeCode (FK)
├── FileUrl
├── VerificationStatus (pending → verified, rejected)
├── VerifiedBy (FK)
├── VerifiedAt
└── RejectReason

ApplicantStatusHistory
├── ApplicantID (FK)
├── OldStatus
├── NewStatus
├── ChangedBy (FK)
├── Note
└── ChangedAt

ReRegistration
├── ApplicantID (FK, unique)
├── ReRegistrationDate
├── Status (pending, completed)
├── VerifiedBy (FK)
└── VerifiedAt

LoaDocument
├── ApplicantID (FK)
├── LoaNumber (unique)
├── FileUrl
├── IssuedBy (FK)
└── IssuedAt

HandoverLog
├── ApplicantID (FK)
├── HandoverDate
├── StudentRefID (FK → academic.students.id)
├── Nim
├── Status (pending, success, failed)
├── ErrorMessage
└── IdempotencyKey (unique)
```

---

## API Endpoints

### Applicant Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/applicants` | Create new applicant |
| GET | `/api/v1/applicants/:id` | Get applicant details |
| PUT | `/api/v1/applicants/:id/submit` | Submit application |
| GET | `/api/v1/applicants/:id/status` | Get application status |

### Biodata Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| PUT | `/api/v1/applicants/:id/biodata` | Update biodata |
| PUT | `/api/v1/applicants/:id/addresses` | Update addresses |
| PUT | `/api/v1/applicants/:id/education` | Update education background |
| PUT | `/api/v1/applicants/:id/family` | Update family members |
| PUT | `/api/v1/applicants/:id/financial` | Update financial profile |
| PUT | `/api/v1/applicants/:id/facility` | Update facility profile |

### Document Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/applicants/:id/documents` | Upload document |
| PUT | `/api/v1/applicants/:id/documents/:doc_id/verify` | Verify document |
| GET | `/api/v1/applicants/:id/documents` | List documents |

### Payment
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/applicants/:id/invoice` | Request invoice |
| GET | `/api/v1/applicants/:id/payment-status` | Check payment status |

### Selection & Admission
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/applicants/:id/selection-result` | Input selection result |
| POST | `/api/v1/applicants/:id/loa` | Issue Letter of Acceptance |

### Academic Handover
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/applicants/:id/handover` | Hand over to Academic |

---

## Service Integrations

### 1. Finance Service (Port 8005)
- **POST** `/api/v1/finance/invoices` - Create invoice
- **GET** `/api/v1/finance/invoices/{id}` - Check invoice status

### 2. Academic Service (Port 8006)
- **POST** `/api/v1/academic/students/generate-from-applicant` - Create student from applicant

### 3. Core Service (Port 8001)
- Sync person data to core.persons table
- Create user account for applicant

---

## Event Publishing

### Outbox Events
| Event Name | Version | Description |
|-----------|---------|-------------|
| `pmb.applicant_created` | v1 | Applicant registration created |
| `pmb.applicant_submitted` | v1 | Applicant submitted application |
| `pmb.document_verified` | v1 | Document verified |
| `pmb.payment_confirmed` | v1 | Payment confirmed |
| `pmb.accepted` | v1 | Applicant accepted |
| `pmb.ready_for_academic` | v1 | Ready for academic transfer |

### Inbox Events
| Event Name | Version | Description |
|-----------|---------|-------------|
| `finance.invoice_paid` | v1 | Invoice payment confirmed |

---

## Status Flow

```
draft → submitted → verified → passed / failed
                               ↓
                           accepted → ready_for_academic → academic_enrolled
                               ↓
                           failed / dropped
```

### Status Descriptions
- **draft**: Initial registration, can edit biodata
- **submitted**: Submitted, waiting for document verification
- **verified**: Documents verified, waiting for selection test
- **passed**: Passed selection test
- **failed**: Failed selection test
- **accepted**: Accepted, issued LoA
- **ready_for_academic**: Transferred to Academic Service
- **academic_enrolled**: Confirmed enrollment by Academic Service
- **dropped**: Withdrawn or cancelled

---

## Configuration

### Environment Variables
```
PORT=8004
DB_HOST=localhost
DB_PORT=5432
DB_NAME=pmb_db
DB_USER=postgres
DB_PASSWORD=postgres
FINANCE_SERVICE_URL=http://localhost:8005
ACADEMIC_SERVICE_URL=http://localhost:8006
PMB_SERVICE_TOKEN=pmb_service_secret_token
```

---

## Implementation Steps

### Phase 1: Core Domain (Week 1-2)
1. [ ] Set up project structure
2. [ ] Database migration (applicants, biodata, addresses, education, family, financial, facility)
3. [ ] Domain models
4. [ ] Repository layer

### Phase 2: Document Management (Week 3)
1. [ ] Document upload endpoint
2. [ ] Document verification workflow
3. [ ] Required documents configuration

### Phase 3: Payment Integration (Week 4)
1. [ ] Invoice request to Finance
2. [ ] Payment status callback
3. [ ] Payment verification

### Phase 4: Selection & Admission (Week 5)
1. [ ] Selection test score input
2. [ ] Pass/fail determination
3. [ ] LoA generation

### Phase 5: Academic Handover (Week 6)
1. [ ] Handover to Academic Service
2. [ ] NIM generation from Academic
3. [ ] Status sync

### Phase 6: Testing & Documentation (Week 7-8)
1. [ ] Unit tests
2. [ ] Integration tests
3. [ ] API documentation
4. [ ] UAT

---

## Dependencies

### Internal Packages
- `github.com/unsia-erp/shared-audit` - Audit logging
- `github.com/unsia-erp/shared-auth` - JWT handling
- `github.com/unsia-erp/shared-errorenvelope` - Error handling
- `github.com/unsia-erp/shared-event` - Event publishing
- `github.com/unsia-erp/shared-httpclient` - HTTP client for service calls

### External
- PostgreSQL (database)
- RabbitMQ (message queue)
