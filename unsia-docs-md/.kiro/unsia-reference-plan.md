# Reference Service Implementation Plan

## Service Overview
**Service Name**: unsia-reference-service (Master Data)
**Port**: 8001 (shared with Core)
**Database**: reference_db (PostgreSQL)

## Core Features

### 1. Master Data Management
- Academic periods (semester, tahun akademik)
- Study programs (prodi)
- Admission paths (jalur masuk)
- PMB waves (gelombang PMB)
- Payment components (biaya kuliah)
- Religions
- Citizienships

### 2. Configuration
- Document types (jenis dokumen)
- Facility types
- Relationship types (keluarga)
- Activity types

---

## Domain Models

```
AcademicPeriod
├── ID (UUID)
├── Code
├── Name
├── Type (SEMESTER, YEAR)
├── StartDate
├── EndDate
├── IsActive
└── CreatedAt

StudyProgram
├── ID (UUID)
├── Code
├── Name
├── Degree (S1, S2, S3, D4)
├── IsActive
└── CreatedAt

PmbWave
├── ID (UUID)
├── Name
├── Year
├── StartDate
├── EndDate
├── IsActive
└── CreatedAt

PaymentComponent
├── ID (UUID)
├── Code
├── Name
├── Amount
├── Type (TUITION, REGISTRATION, OTHER)
├── IsActive
└── CreatedAt
```

---

## Implementation Timeline: 4-6 weeks
