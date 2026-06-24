# LMS Service Implementation Plan

## Service Overview
**Service Name**: unsia-lms-service (Learning Management System)
**Port**: 8003
**Database**: lms_db (PostgreSQL)

## Core Features

### 1. Course Management
- Create & manage courses (matakuliah)
- Course materials (materi, files)
- Assignments (tugas)
- Quizzes & exams

### 2. Class & Enrollment
- Class sections (kelas)
- Student enrollment
- Teaching schedule

### 3. Learning Activities
- Discussion forums (diskusi)
- Announcements (pengumuman)
- Live sessions (sesi)

### 4. Assessments
- Grade submissions
- Feedback & rubrics

---

## Domain Models

```
Course
├── ID (UUID)
├── Code (unique)
├── Name
├── StudyProgramID (FK)
├── Credits
├── Description
└── Status

ClassSection
├── ID (UUID)
├── CourseID (FK)
├── SemesterID (FK)
├── LecturerID (FK → users.id)
├── ClassCode (A, B, C...)
└── Capacity

Enrollment
├── ClassSectionID (FK)
├── StudentID (FK)
└── Status

Assignment
├── ID (UUID)
├── ClassSectionID (FK)
├── Title
├── Description
├── DueDate
└── MaxScore
```

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/courses` | List courses |
| POST | `/api/v1/classes` | Create class |
| POST | `/api/v1/enrollments` | Enroll student |
| POST | `/api/v1/assignments` | Create assignment |

---

## Implementation Steps

### Phase 1: Core (Week 1-2)
- Project setup, migrations, models

### Phase 2: Course Management (Week 3-4)
- Courses, materials, assignments

### Phase 3: Class & Enrollment (Week 5)
- Class sections, student enrollment

### Phase 4: Learning Activities (Week 6)
- Forums, announcements, live sessions

### Phase 5: Assessment (Week 7)
- Grading, feedback

### Phase 6: Testing (Week 8)
- Tests & documentation
