# Assessment Service Implementation Plan

## Service Overview
**Service Name**: unsia-assessment-service
**Port**: 8007
**Database**: assessment_db (PostgreSQL)

## Core Features

### 1. Exam Management
- Create & manage exams (uts, uas)
- Question banks
- Exam schedules

### 2. Exam Sessions
- Student exam slots
- Randomization
- Proctoring (optional)

### 3. Grading
- Auto-grading for objective
- Manual grading for essay
- Grade distribution

---

## Domain Models

```
Exam
├── ID (UUID)
├── ClassSectionID (FK)
├── Type (UTS, UAS)
├── Title
├── Duration (minutes)
├── StartDate
├── EndDate
├── Status (DRAFT, PUBLISHED, COMPLETED)
└── CreatedAt

Question
├── ID (UUID)
├── ExamID (FK)
├── QuestionText
├── QuestionType (MULTIPLE_CHOICE, ESSAY)
├── Options (JSON)
├── CorrectAnswer
├── Score
└── Order

ExamSession
├── ID (UUID)
├── ExamID (FK)
├── StudentID (FK)
├── StartTime
├── EndTime
├── Score
├── Answers (JSON)
└── Status (IN_PROGRESS, SUBMITTED, GRADED)
```

---

## Implementation Timeline: 6-8 weeks
