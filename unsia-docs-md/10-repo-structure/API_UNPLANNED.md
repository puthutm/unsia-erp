# API Implementation Plan - Unplanned APIs

This document lists the APIs that still need implementation (marked as "Planned" in API_LIST.md).

## Status Summary

| Service | Port | Status | Priority |
|---------|------|--------|----------|
| LMS Service | 8003 | Partially Done | HIGH |
| Assessment Service | 8006 | Partially Done | HIGH |
| HRIS Service | 8007 | Partially Done | MEDIUM |
| CRM Service | 8008 | Partially Done | MEDIUM |
| Reference Service | 8009 | Partially Done | LOW |
| Portal Service | 8010 | Partially Done | LOW |

---

## 1. LMS Service (Port 8003) - HIGH PRIORITY

### Current Status: Partially Done
- Class handlers: Partially implemented  
- Enrollment handlers: Partially implemented
- Grade sync to Academic: DONE ✅

### APIs Still Needed:
- `GET /api/v1/lms/courses` - List LMS courses
- `POST /api/v1/lms/courses` - Create course
- `GET /api/v1/lms/classes` - List classes  
- `POST /api/v1/lms/classes` - Create class
- `POST /api/v1/lms/enrollments` - Enroll student
- `GET /api/v1/lms/enrollments` - List enrollments
- `POST /api/v1/lms/assignments` - Create assignment
- `GET /api/v1/lms/assignments` - List assignments
- `POST /api/v1/lms/sessions` - Create session
- `GET /api/v1/lms/sessions` - List sessions
- `POST /api/v1/lms/materials` - Upload material
- `GET /api/v1/lms/materials` - List materials

### Implementation Plan:
1. Create LMS course_handler.go
2. Create LMS class_handler.go  
3. Create LMS enrollment_handler.go
4. Create LMS assignment_handler.go
5. Create LMS session_handler.go
6. Create LMS material_handler.go
7. Add LMS router endpoints

---

## 2. Assessment Service (Port 8006) - HIGH PRIORITY

### Current Status: Partially Done
- Question Bank: Partially implemented

### APIs Still Needed:
- `GET /api/v1/assessment-sessions` - List sessions
- `POST /api/v1/assessment-sessions` - Create session
- `GET /api/v1/participants` - List participants
- `POST /api/v1/participants` - Register participant
- `GET /api/v1/attempts` - List attempts
- `POST /api/v1/attempts` - Start attempt
- `POST /api/v1/attempts/:id/submit` - Submit attempt

### Implementation Plan:
1. Create assessment_session_handler.go
2. Create participant_handler.go
3. Create attempt_handler.go
4. Add assessment router endpoints

---

## 3. HRIS Service (Port 8007) - MEDIUM PRIORITY

### Current Status: Partially Done
- Employee handlers: Partially implemented
- Attendance handlers: Partially implemented
- Leave handlers: Partially implemented

### APIs Still Needed:
- `GET /api/v1/employees` - List employees
- `POST /api/v1/employees` - Create employee
- `GET /api/v1/employees/:id` - Get employee
- `PUT /api/v1/employees/:id` - Update employee
- `GET /api/v1/attendances` - List attendances
- `POST /api/v1/attendances` - Record attendance
- `GET /api/v1/leave-requests` - List leave requests
- `POST /api/v1/leave-requests` - Submit leave request

### Implementation Plan:
1. Create HRIS router endpoints for employee
2. Create HRIS router endpoints for attendance  
3. Create HRIS router endpoints for leave
4. Add BKD handler
5. Add performance review handler

---

## 4. CRM Service (Port 8008) - MEDIUM PRIORITY

### Current Status: Partially Done
- Lead handlers: Partially implemented
- Pipeline handlers: Partially implemented

### APIs Still Needed:
- `GET /api/v1/contacts` - List contacts
- `POST /api/v1/contacts` - Create contact
- `GET /api/v1/opportunities` - List opportunities
- `POST /api/v1/opportunities` - Create opportunity
- `GET /api/v1/campaigns` - List campaigns
- `POST /api/v1/campaigns` - Create campaign

### Implementation Plan:
1. Create contact_handler.go
2. Create opportunity_handler.go  
3. Create campaign_handler.go
4. Add CRM router endpoints

---

## 5. Reference Service (Port 8009) - LOW PRIORITY

### Current Status: Partially Done
- Reference data: Partially implemented

### APIs Still Needed:
- `GET /api/v1/provinces` - List provinces
- `GET /api/v1/regencies` - List regencies
- `GET /api/v1/districts` - List districts
- `GET /api/v1/villages` - List villages

### Implementation Plan:
1. Add more reference data handlers
2. Create reference router endpoints

---

## 6. Portal Service (Port 8010) - LOW PRIORITY

### Current Status: Partially Done
- Basic handlers: Partially implemented

### APIs Still Needed:
- `GET /api/v1/portal/news` - List news
- `POST /api/v1/portal/news` - Create news
- `GET /api/v1/portal/announcements` - List announcements
- `POST /api/v1/portal/announcements` - Create announcement
- `GET /api/v1/portal/events` - List events
- `POST /api/v1/portal/events` - Create event

### Implementation Plan:
1. Create portal news_handler.go
2. Create portal announcement_handler.go  
3. Create portal event_handler.go
4. Add portal router endpoints

---

## Implementation Order

1. **Phase 1 (HIGH)**: LMS + Assessment Service
2. **Phase 2 (MEDIUM)**: HRIS + CRM Service  
3. **Phase 3 (LOW)**: Reference + Portal Service

## Notes

- All "✅ DONE" services in API_LIST.md are already implemented
- Only services marked as "Endpoints (Planned)" need implementation
- The Academic service "akademik" is already matched with LMS for grade sync ✅
