# Academic Service Implementation Plan

## Service Overview
**Service Name**: unsia-academic-service (Academic Management)
**Port**: 8006
**Database**: academic_db (PostgreSQL)

## Core Features

### 1. Student Management
- Create & manage student records
- Generate NIM (Nomor Induk Mahasiswa)
- Student status (active, graduated, dropped out, suspended)
- Student profiles and detail information

### 2. Academic Program
- Study programs management
- Curriculums (structured courses per study program)
- Course management
- Class schedules
- Academic years/semesters

### 3. Enrollment & KRS
- Course registration (KRS)
- Enrollment validation
- Course capacity management
- Schedule conflict detection

### 4. Grades & Transcripts
- Input grades
- Grade calculations (IPS, IPK)
- Transcripts generation
- GPA management

### 5. Graduation
- Graduation eligibility check
- Alumni transition
- Degree certificate management

---

## Domain Models

### Student Entity

```
Student
├── ID (UUID)
├── PersonID (FK → core.persons.id)
├── Nim (unique, generated)
├── StudyProgramID (FK → ref.study_programs.id)
├── CuriculumID (FK → curriculums.id)
├── EntryYear
├── EntryPeriodID (FK → ref.academic_periods.id)
├── Status (active, graduated, dropped, suspended)
├── GraduationDate
├── GraduationPeriodID (FK → ref.academic_periods.id)
├── TotalCreditsEarned
├── Gpa
├── CreatedAt
└── UpdatedAt
```

### Academic Entities

```
StudyProgram
├── ID (UUID)
├── Code (unique)
├── Name
├── Degree (S1, S2, S3, D4)
├── Capacity
├── Accreditation
├── HeadOfProgramID (FK → core.users.id)
├── Status
└── CreatedAt

Curriculum
├── ID (UUID)
├── StudyProgramID (FK)
├── Version
├── Year
├── IsActive
└── CreditsTotal

CurriculumCourse
├── CurriculumID (FK)
├── CourseID (FK)
├── Semester
├── Credits
└── Type (WAJIB, PILIHAN)
```

---

## API Endpoints

### Student Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/students` | Create student |
| GET | `/api/v1/students/:id` | Get student details |
| GET | `/api/v1/students/nim/:nim` | Get by NIM |
| PUT | `/api/v1/students/:id/status` | Update status |

### Academic Program
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/study-programs` | List study programs |
| GET | `/api/v1/curriculums` | List curriculums |
| GET | `/api/v1/courses` | List courses |

### Enrollment (KRS)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/enrollments` | Register courses |
| GET | `/api/v1/students/:id/krs` | Get KRS |
| DELETE | `/api/v1/enrollments/:id` | Drop course |

### Grades
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/grades` | Input grades |
| GET | `/api/v1/students/:id/transcript` | Get transcript |

---

## Service Integrations

### PMB Service (Port 8004)
- Receive accepted applicants → Create student records
- Generate NIM

### Core Service (Port 8001)
- Sync person data from core.persons
- Create user accounts for students

### Finance Service (Port 8008)
- Check tuition payment status for enrollment

---

## Implementation Steps

### Phase 1: Core (Week 1-2)
- [ ] Set up project
- [ ] Database migrations
- [ ] Domain models
- [ ] Repository

### Phase 2: Student Management (Week 3)
- [ ] CRUD operations
- [ ] NIM generation
- [ ] Status management

### Phase 3: Academic Programs (Week 4)
- [ ] Study programs
- [ ] Curriculums
- [ ] Courses

### Phase 4: Enrollment (Week 5-6)
- [ ] KRS functionality
- [ ] Schedule conflict detection
- [ ] Capacity management

### Phase 5: Grades (Week 7)
- [ ] Grade input
- [ ] IPS/IPK calculations
- [ ] Transcripts

### Phase 8: Graduation & Testing (Week 8)
- [ ] Graduation workflow
- [ ] Tests & documentation
