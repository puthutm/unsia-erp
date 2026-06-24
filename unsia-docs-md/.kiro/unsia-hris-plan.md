# HRIS Service Implementation Plan

## Service Overview
**Service Name**: unsia-hris-service (Human Resources Information System)
**Port**: 8009
**Database**: hris_db (PostgreSQL)

## Core Features

### 1. Employee Management
- Create & manage employee records (Dosen, Tendik)
- Generate NIP (Nomor Induk Pegawai)
- Employee profiles (biodata, education, experience)
- Employee status (active, resigned, pensioned)

### 2. Organizational Structure
- Units/Departments
- Positions & Job Titles
- Organizational chart

### 3. Attendance & Presence
- Daily attendance recording
- Absence management
- Leave requests (cuti, sakit, izin)
- Overtime tracking

### 4. Payroll
- Salary components (gaji pokok, tunjangan)
- Payroll processing
- Slip gaji

### 5. Performance Appraisal
- Performance reviews
- KPIs assessment

---

## Domain Models

### Employee Entity

```
Employee
├── ID (UUID)
├── PersonID (FK → core.persons.id)
├── Nip (unique, generated)
├── EmployeeType (DOSEN, TENDIK)
├── UnitID (FK → units.id)
├── PositionID (FK → positions.id)
├── Status (active, resigned, pensioned)
├── JoinDate
├── ResignDate
├── CreatedAt
└── UpdatedAt
```

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/employees` | Create employee |
| GET | `/api/v1/employees/:id` | Get details |
| PUT | `/api/v1/employees/:id/status` | Update status |

---

## Implementation Steps

### Phase 1: Core (Week 1-2)
- Set up project, migrations, models

### Phase 2: Employee Management (Week 3)
- CRUD, NIP generation

### Phase 3: Attendance (Week 4-5)
- Presence & leave management

### Phase 4: Payroll (Week 6-7)
- Salary processing

### Phase 5: Performance (Week 8)
- Appraisal system
