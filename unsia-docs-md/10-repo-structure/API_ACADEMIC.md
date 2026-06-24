# Academic Service API Documentation

**Base URL:** `http://localhost:8002`  
**Database:** `academic_db`  
**Status:** ✅ DONE (SUDAH MATCH dengan LMS untuk Grade)

---

## Table of Contents

1. [Student APIs](#1-student-apis)
2. [Course APIs](#2-course-apis)
3. [Grade APIs](#3-grade-apis)
4. [KRS APIs](#4-krs-apis)
5. [Advisor APIs](#5-advisor-apis)
6. [Graduation APIs](#6-graduation-apis)
7. [Transcript APIs](#7-transcript-apis)

---

## 1. Student APIs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/students` | List all students (with pagination, filter by prodi, angkatan) |
| POST | `/api/v1/students` | Create new student |
| GET | `/api/v1/students/:id` | Get student by ID |
| PUT | `/api/v1/students/:id` | Update student |
| DELETE | `/api/v1/students/:id` | Delete student |
| GET | `/api/v1/students/nim/:nim` | Get student by NIM |

### Request/Response Examples

**GET /api/v1/students**
```json
// Query Params: page, limit, study_program_id, enrollment_year
{
  "data": [
    {
      "id": "uuid",
      "nim": "2021001",
      "name": "John Doe",
      "study_program_id": "uuid",
      "enrollment_year": 2021,
      "status": "active"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

**POST /api/v1/students**
```json
{
  "nim": "2021001",
  "name": "John Doe",
  "study_program_id": "uuid",
  "enrollment_year": 2021,
  "email": "john.doe@unsia.ac.id",
  "phone": "+6281234567890"
}
```

---

## 2. Course APIs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/courses` | List all courses |
| POST | `/api/v1/courses` | Create new course |
| GET | `/api/v1/courses/:id` | Get course by ID |
| PUT | `/api/v1/courses/:id` | Update course |
| DELETE | `/api/v1/courses/:id` | Delete course |

### Request/Response Examples

**GET /api/v1/courses**
```json
// Query Params: page, limit, study_program_id, semester
{
  "data": [
    {
      "id": "uuid",
      "code": "CS101",
      "name": "Introduction to Computer Science",
      "credits": 4,
      "study_program_id": "uuid",
      "semester": 1
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

**POST /api/v1/courses**
```json
{
  "code": "CS101",
  "name": "Introduction to Computer Science",
  "credits": 4,
  "study_program_id": "uuid",
  "semester": 1,
  "description": "Basic concepts of computer science"
}
```

---

## 3. Grade APIs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/grades` | List all grades |
| POST | `/api/v1/grades` | Create new grade |
| GET | `/api/v1/grades/:id` | Get grade by ID |
| PUT | `/api/v1/grades/:id` | Update grade |
| POST | `/api/v1/grades/:id/submit` | Submit grade for review |
| POST | `/api/v1/grades/:id/finalize` | Finalize grade (lock) |
| POST | `/api/v1/grades/:id/entries` | Enter student grade |
| POST | `/api/v1/grades/:id/entries/bulk` | Bulk enter grades |
| GET | `/api/v1/grades/student/:student_id` | Get student grades |
| GET | `/api/v1/grades/transcript/:student_id` | Get transcript |
| GET | `/api/v1/grades/ipk/:student_id` | Get IPK (GPA) |
| GET | `/api/v1/grades/ips/:student_id` | Get IPS (semester GPA) |
| POST | `/api/v1/grades/conversion` | Update grade conversion |
| GET | `/api/v1/grades/conversion` | Get grade conversions |

### Note
Grade source can be: `lms`, `exam`, `quiz`, `assignment` - allows sync from LMS

### Request/Response Examples

**GET /api/v1/grades**
```json
// Query Params: page, limit, course_id, semester, study_program_id
{
  "data": [
    {
      "id": "uuid",
      "course_id": "uuid",
      "semester": "2021-1",
      "status": "draft",
      "source": "lms",
      "created_by": "uuid",
      "entries_count": 25,
      "submitted_at": null,
      "finalized_at": null
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 20
  }
}
```

**POST /api/v1/grades**
```json
{
  "course_id": "uuid",
  "semester": "2021-1",
  "source": "lms"
}
```

**POST /api/v1/grades/:id/entries**
```json
{
  "student_id": "uuid",
  "grade": "A",
  "numeric_grade": 4.0,
  "attendance": 100,
  "assignment_score": 85,
  "midterm_score": 80,
  "final_score": 88
}
```

**POST /api/v1/grades/:id/entries/bulk**
```json
{
  "entries": [
    {
      "student_id": "uuid",
      "grade": "A",
      "numeric_grade": 4.0,
      "attendance": 100,
      "assignment_score": 85,
      "midterm_score": 80,
      "final_score": 88
    }
  ]
}
```

**GET /api/v1/grades/transcript/:student_id**
```json
{
  "student_id": "uuid",
  "nim": "2021001",
  "name": "John Doe",
  "gpa": 3.75,
  "total_credits": 120,
  "grades": [
    {
      "semester": "2021-1",
      "course": "CS101",
      "course_name": "Introduction to Computer Science",
      "credits": 4,
      "grade": "A",
      "numeric_grade": 4.0
    }
  ]
}
```

---

## 4. KRS APIs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/krs` | List KRS |
| POST | `/api/v1/krs` | Create KRS |
| GET | `/api/v1/krs/:id` | Get KRS by ID |
| PUT | `/api/v1/krs/:id` | Update KRS |
| DELETE | `/api/v1/krs/:id` | Delete KRS |
| POST | `/api/v1/krs/:id/submit` | Submit KRS |
| POST | `/api/v1/krs/:id/approve` | Approve KRS (Dosen PA) |
| POST | `/api/v1/krs/:id/reject` | Reject KRS |

### Note
KRS submission requires clearance status from Finance Service (Port 8004)

### Request/Response Examples

**GET /api/v1/krs**
```json
// Query Params: page, limit, student_id, semester, status
{
  "data": [
    {
      "id": "uuid",
      "student_id": "uuid",
      "semester": "2021-1",
      "status": "approved",
      "total_credits": 18,
      "courses": [
        {
          "course_id": "uuid",
          "course_name": "CS101"
        }
      ],
      "submitted_at": "2021-01-15T10:00:00Z",
      "approved_at": "2021-01-16T14:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

**POST /api/v1/krs**
```json
{
  "student_id": "uuid",
  "semester": "2021-1",
  "course_ids": ["uuid1", "uuid2", "uuid3"]
}
```

**POST /api/v1/krs/:id/submit**
```json
{
  "notes": "KRS submitted for approval"
}
```

---

## 5. Advisor APIs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/advisors` | List student advisors |
| POST | `/api/v1/advisors` | Create advisor assignment |
| GET | `/api/v1/advisors/:id` | Get advisor by ID |
| PUT | `/api/v1/advisors/:id` | Update advisor |
| GET | `/api/v1/advisors/:id/mahasiswa-bimbingan` | Get students under this advisor |

### Request/Response Examples

**GET /api/v1/advisors**
```json
{
  "data": [
    {
      "id": "uuid",
      "lecturer_id": "uuid",
      "student_id": "uuid",
      "academic_year": "2021",
      "status": "active"
    }
  ]
}
```

**GET /api/v1/advisors/:id/mahasiswa-bimbingan**
```json
{
  "advisor_id": "uuid",
  "lecturer_name": "Dr. Jane Smith",
  "students": [
    {
      "nim": "2021001",
      "name": "John Doe",
      "study_program": "Computer Science",
      "semester": 5
    }
  ]
}
```

---

## 6. Graduation APIs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/graduations` | List graduation records |
| POST | `/api/v1/graduations` | Create graduation record |
| GET | `/api/v1/graduations/:id` | Get graduation by ID |
| PUT | `/api/v1/graduations/:id` | Update graduation |
| POST | `/api/v1/graduations/:id/approve` | Approve graduation |

### Request/Response Examples

**POST /api/v1/graduations**
```json
{
  "student_id": "uuid",
  "graduation_period": "2024-1",
  "gpa": 3.75,
  "total_credits": 144,
  "threshold_credits": 144,
  "status": "verified"
}
```

---

## 7. Transcript APIs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/transcripts` | List transcripts |
| GET | `/api/v1/transcripts/:student_id` | Get student transcript |
| POST | `/api/v1/transcripts/:student_id/download` | Generate transcript PDF |

### Request/Response Examples

**GET /api/v1/transcripts/:student_id**
```json
{
  "student_id": "uuid",
  "nim": "2021001",
  "name": "John Doe",
  "study_program": "Computer Science",
  "enrollment_date": "2021-08-01",
  "graduation_date": null,
  "gpa": 3.75,
  "cgpa": 3.75,
  "total_credits": 120,
  "required_credits": 144,
  "academic_standing": "good",
  "semesters": [
    {
      "semester": "2021-1",
      "gpa": 3.8,
      "credits": 18,
      "courses": [
        {
          "code": "CS101",
          "name": "Introduction to Computer Science",
          "credits": 4,
          "grade": "A",
          "points": 4.0
        }
      ]
    }
  ]
}
```

---

## Integration

### External Services

| Service | Endpoint | Purpose |
|---------|----------|---------|
| **LMS Service** (Port 8003) | Grade source sync | `source: lms` tracks grades from LMS |
| **Finance Service** (Port 8004) | Clearance status | KRS submission validation |
| **Core Service** (Port 8001) | Auth/Token | Authentication |

### Environment Variables

```env
DATABASE_URL=postgres://user:pass@localhost:5432/academic_db
PORT=8002
SSO_CORE_URL=http://localhost:8001
LMS_SERVICE_URL=http://localhost:8003
FINANCE_SERVICE_URL=http://localhost:8004
```

---

## Error Responses

All endpoints may return standard error responses:

```json
{
  "error": {
    "code": "ERR_CODE",
    "message": "Error description"
  }
}
```

Common error codes:
- `VALIDATION_ERROR` - Invalid request body
- `NOT_FOUND` - Resource not found
- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Insufficient permissions
- `CONFLICT` - Resource already exists

---

## Authentication

All endpoints require Bearer token authentication:

```
Authorization: Bearer <token>
```

Token can be obtained from Core Service (Port 8001).

---

*Generated: 2024*
