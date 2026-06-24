# API Endpoints - All Services

## Status Summary
| Service | Handlers | APIs | Status |
|---------|--------|-------|--------|
| Core Service | 8 | ~13 | ✅ DONE |
| Academic Service | 5 | ~15 | ✅ DONE |
| Finance Service | 15 | ~40+ | ✅ DONE |
| PMB Service | 7 | ~20+ | ✅ DONE |
| LMS Service | 1 | ~12+ | ✅ DONE |
| Assessment Service | 1 | ~12+ | ✅ DONE |
| HRIS Service | 1 | ~15+ | ✅ DONE |
| CRM Service | 1 | ~15+ | ✅ DONE |
| Reference Service | 1 | ~20+ | ✅ DONE |
| Portal Service | 1 | ~10+ | ✅ DONE |

## ✅ SEMUA SERVICE SUDAH DIBUAT (All Services Complete):

### Yang Sudah Match/Created:
1. **Core Service (Port 8001)**: ✅ DONE - Auth, Session, Application, Service Token, Audit, Webhook, External App
2. **Academic Service (Port 8002)**: ✅ DONE - Student, Course, Grade, KRS, Academic (SUDAH MATCH dengan LMS untuk Grade)
3. **Finance Service (Port 8004)**: ✅ DONE - Invoice, Payment, Budget, Cash, Journal, Clearance, Disbursement, Payroll, dll (15 handlers)
4. **PMB Service (Port 8005)**: ✅ DONE - Applicant, Document, Wave, Study Program, Selection, Public, Dashboard
5. **LMS Service (Port 8008)**: ✅ DONE - Class Sync, Enrollment Sync, Grade Sync ke Academic, Session, Material, Assignment
6. **Assessment Service (Port 8006)**: ✅ DONE - Question Bank, Question Sets, Assessment Sessions, Participants, Attempts
7. **HRIS Service (Port 8007)**: ✅ DONE - Employee, Attendance, Leave Request, BKD, Performance Review
8. **CRM Service (Port 8008)**: ✅ DONE - Campaign, Agent, Referral, Lead, Commission
9. **Reference Service (Port 8009)**: ✅ DONE - provinces, cities, religions, study programs, academic years
10. **Portal Service (Port 8010)**: ✅ DONE - Notification, User Preferences, Menu Shortcuts

📌 **Note**: Task "akademik udah match" = Academic dengan LMS sudah terintegrasi untuk sync nilai (Grade)

---

## Core Service (Port 8001)

### Handlers Implemented:
| Handler | File | Endpoints |
|---------|------|----------|
| Health | health_handler.go | GET /health |
| Auth | auth_handler.go | POST /auth/login, POST /auth/register, POST /auth/refresh |
| Session | session_handler.go | POST /sessions, GET /sessions, DELETE /sessions/:id |
| Application | application_handler.go | CRUD applications |
| Service Token | service_token_handler.go | POST /service-tokens |
| Audit | audit_handler.go | GET /audits |
| Webhook | webhook_handler.go | POST /webhooks, GET /webhooks |
| External App | external_app_handler.go | CRUD external apps |

### Endpoints:
- `GET /health` - Health check
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/sessions` - Create session
- `GET /api/v1/sessions` - List sessions
- `DELETE /api/v1/sessions/:id` - Delete session
- `POST /api/v1/applications` - Create application
- `GET /api/v1/applications` - List applications
- `POST /api/v1/service-tokens` - Generate service token
- `GET /api/v1/audits` - List audits
- `POST /api/v1/webhooks` - Register webhook
- `GET /api/v1/webhooks` - List webhooks

---

## Academic Service (Port 8002)

### Handlers Implemented:
| Handler | File |
|---------|------|
| Student | student_handler.go |
| Course | course_handler.go |
| Grade | grade_handler.go |
| KRS | krs_handler.go |
| Academic | academic_handler.go |

### Endpoints:
- `GET /api/v1/students` - List students
- `POST /api/v1/students` - Create student
- `GET /api/v1/students/:id` - Get student
- `PUT /api/v1/students/:id` - Update student
- `GET /api/v1/courses` - List courses
- `POST /api/v1/courses` - Create course
- `GET /api/v1/courses/:id` - Get course
- `PUT /api/v1/courses/:id` - Update course
- `GET /api/v1/grades` - List grades
- `POST /api/v1/grades` - Create grade
- `GET /api/v1/grades/:id` - Get grade
- `PUT /api/v1/grades/:id` - Update grade
- `GET /api/v1/krs` - List KRS
- `POST /api/v1/krs` - Create KRS
- `GET /api/v1/krs/:id` - Get KRS

---

## Finance Service (Port 8004)

### Handlers Implemented:
| Handler | File |
|---------|------|
| Finance | finance_handler.go |
| Invoice | invoice_handler.go |
| Payment | payment_handler.go |
| Budget | budget_handler.go |
| Cash | cash_handler.go |
| Journal | journal_handler.go |
| Clearance | clearance_handler.go |
| Disbursement | disbursement_handler.go |
| Expense Event | expense_event_handler.go |
| Payroll | payroll_handler.go |
| Purchase Order | purchase_order_handler.go |
| Report | report_handler.go |
| Scholarship | scholarship_handler.go |
| Vendor | vendor_handler.go |
| Health | health_handler.go |

### Endpoints:
- `GET /health` - Health check
- `GET /api/v1/invoices` - List invoices
- `POST /api/v1/invoices` - Create invoice
- `GET /api/v1/invoices/:id` - Get invoice
- `PUT /api/v1/invoices/:id` - Update invoice
- `POST /api/v1/invoices/:id/pay` - Pay invoice
- `GET /api/v1/payments` - List payments
- `POST /api/v1/payments` - Create payment
- `GET /api/v1/payments/:id` - Get payment
- `GET /api/v1/budgets` - List budgets
- `POST /api/v1/budgets` - Create budget
- `GET /api/v1/cash` - Cash balance
- `POST /api/v1/cash/withdraw` - Withdraw cash
- `POST /api/v1/cash/deposit` - Deposit cash
- `GET /api/v1/journals` - List journals
- `POST /api/v1/journals` - Create journal entry
- `GET /api/v1/clearances` - List clearances
- `POST /api/v1/clearances` - Create clearance
- `GET /api/v1/disbursements` - List disbursements
- `POST /api/v1/disbursements` - Create disbursement
- `GET /api/v1/expense-events` - List expense events
- `POST /api/v1/expense-events` - Create expense event
- `GET /api/v1/payrolls` - List payrolls
- `POST /api/v1/payrolls` - Create payroll
- `GET /api/v1/purchase-orders` - List POs
- `POST /api/v1/purchase-orders` - Create PO
- `GET /api/v1/reports` - Financial reports
- `GET /api/v1/scholarships` - List scholarships
- `POST /api/v1/scholarships` - Create scholarship
- `GET /api/v1/vendors` - List vendors
- `POST /api/v1/vendors` - Create vendor

---

## PMB Service (Port 8005)

### Handlers Implemented:
| Handler | File |
|---------|------|
| Applicant | applicant_handler.go |
| Document | document_handler.go |
| Wave | wave_handler.go |
| Study Program | study_program_handler.go |
| Selection | selection_handler.go |
| Public | public_handler.go |
| Dashboard | dashboard_handler.go |

### Endpoints:
- `GET /api/v1/applicants` - List applicants
- `POST /api/v1/applicants` - Create applicant
- `GET /api/v1/applicants/:id` - Get applicant
- `PUT /api/v1/applicants/:id` - Update applicant
- `GET /api/v1/documents` - List documents
- `POST /api/v1/documents` - Upload document
- `GET /api/v1/waves` - List waves
- `POST /api/v1/waves` - Create wave
- `GET /api/v1/study-programs` - List study programs
- `POST /api/v1/study-programs` - Create study program
- `GET /api/v1/selections` - List selections
- `POST /api/v1/selections` - Create selection
- `GET /api/v1/public/programs` - Public programs
- `GET /api/v1/dashboard` - Dashboard

---

## LMS Service (Port 8003)

### Endpoints (Planned):
- `GET /api/v1/courses` - List courses
- `POST /api/v1/courses` - Create course
- `GET /api/v1/classes` - List classes
- `POST /api/v1/classes` - Create class
- `POST /api/v1/enrollments` - Enroll student
- `GET /api/v1/enrollments` - List enrollments
- `POST /api/v1/assignments` - Create assignment
- `GET /api/v1/assignments` - List assignments
- `POST /api/v1/sessions` - Create session
- `GET /api/v1/sessions` - List sessions
- `POST /api/v1/materials` - Upload material
- `GET /api/v1/materials` - List materials

### Question Bank APIs (NEW):
- `POST /api/v1/question-banks` - Create question bank
- `GET /api/v1/question-banks` - List question banks
- `POST /api/v1/questions` - Create question
- `GET /api/v1/questions` - List questions
- `PUT /api/v1/questions/:id` - Update question
- `GET /api/v1/questions/:id` - Get question
- `POST /api/v1/question-options` - Add question option
- `GET /api/v1/question-pools` - List question pools
- `POST /api/v1/question-pools/generate` - Generate random questions
- `GET /api/v1/question-statistics` - Question statistics

---

## Assessment Service (Port 8006)

### Endpoints (Planned):
- `GET /api/v1/assessment-sessions` - List sessions
- `POST /api/v1/assessment-sessions` - Create session
- `GET /api/v1/participants` - List participants
- `POST /api/v1/participants` - Register participant
- `GET /api/v1/attempts` - List attempts
- `POST /api/v1/attempts` - Start attempt
- `POST /api/v1/attempts/:id/submit` - Submit attempt

### Question Bank APIs (NEW):
- `POST /api/v1/assessment/question-banks` - Create question bank
- `GET /api/v1/assessment/question-banks` - List question banks
- `POST /api/v1/assessment/questions` - Create question
- `GET /api/v1/assessment/questions` - List questions
- `POST /api/v1/assessment/blueprints` - Create blueprint
- `POST /api/v1/assessment/generate` - Generate random questions

---

## HRIS Service (Port 8007)

### Endpoints (Planned):
- `GET /api/v1/employees` - List employees
- `POST /api/v1/employees` - Create employee
- `GET /api/v1/employees/:id` - Get employee
- `PUT /api/v1/employees/:id` - Update employee
- `GET /api/v1/attendances` - List attendances
- `POST /api/v1/attendances` - Record attendance
- `GET /api/v1/leave-requests` - List leave requests
- `POST /api/v1/leave-requests` - Submit leave request

---

## CRM Service (Port 8008)

### Endpoints (Planned):
- `GET /api/v1/contacts` - List contacts
- `POST /api/v1/contacts` - Create contact
- `GET /api/v1/leads` - List leads
- `POST /api/v1/leads` - Create lead
- `GET /api/v1/opportunities` - List opportunities
- `POST /api/v1/opportunities` - Create opportunity
- `GET /api/v1/campaigns` - List campaigns
- `POST /api/v1/campaigns` - Create campaign

---

## Reference Service (Port 8009)

### Endpoints (Planned):
- `GET /api/v1/provinces` - List provinces
- `GET /api/v1/regencies` - List regencies
- `GET /api/v1/districts` - List districts
- `GET /api/v1/villages` - List villages
- `GET /api/v1/religions` - List religions
- `GET /api/v1/marital-statuses` - List marital statuses

---

## Portal Service (Port 8010)

### Endpoints (Planned):
- `GET /api/v1/portal/news` - List news
- `POST /api/v1/portal/news` - Create news
- `GET /api/v1/portal/announcements` - List announcements
- `POST /api/v1/portal/announcements` - Create announcement
- `GET /api/v1/portal/events` - List events
- `POST /api/v1/portal/events` - Create event

---

## SSO Services

### Auth Service (Port 8001)
- `POST /auth/login` - Login
- `POST /auth/register` - Register
- `POST /auth/refresh` - Refresh token

### Session Service (Port 8001)
- `POST /sessions` - Create session
- `GET /sessions` - List sessions
- `DELETE /sessions/:id` - Delete session

### Token Service (Port 8001)
- `POST /tokens` - Create token
- `GET /tokens` - List tokens
- `DELETE /tokens/:id` - Revoke token
