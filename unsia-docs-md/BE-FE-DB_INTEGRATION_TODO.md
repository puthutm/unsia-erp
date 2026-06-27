# BE-FE-DB Integration TODO

Tracking progress untuk integrasi BE-FE-DB UNSIA ERP.

## Status Legend
- [ ] Todo
- 🔄 In Progress
- ✅ Done
- ❌ Blocked

---

## Phase 1: Database Infrastructure ✅

### 1.1 Database Setup
- [✅] Setup 10 databases terpisah di postgres
- [✅] Enable pgcrypto extension
- [✅] Configure docker-compose.yml

### 1.2 Cross-DB References
- [✅] Add student_id references di finance.student_clearances
- [✅] Add person_id references di semua services
- [✅] Add user_id references di PMB dan Academic

---

## Phase 2: HTTP Client Infrastructure ✅

### 2.1 shared-httpclient Package
- [✅] Circuit breaker implementation
- [✅] Auto-retry dengan exponential backoff
- [✅] Service token authentication
- [✅] Correlation ID propagation

---

## Phase 3: Service Clients 🔄

### 3.1 Shared Service Client Factory
- [ ] Create packages/shared-serviceclient/
- [ ] NewServiceClient factory function

### 3.2 Academic Service Clients
- [ ] client/academic_client.go - GetStudentByID
- [ ] client/academic_client.go - GetKRS
- [ ] client/academic_client.go - SubmitGrade
- [ ] client/finance_client.go - CheckClearance
- [ ] client/finance_client.go - GetInvoice
- [ ] client/reference_client.go - GetPeriod

### 3.3 Finance Service Clients  
- [ ] client/academic_client.go - GetStudentInfo
- [ ] client/reference_client.go - GetAcademicPeriod
- [ ] client/reference_client.go - GetStudyProgram

### 3.4 PMB Service Clients
- [ ] client/core_client.go - CreateUserPerson
- [ ] client/academic_client.go - CreateStudent
- [ ] client/academic_client.go - GetStudyPrograms

### 3.5 Service Token Configuration
- [ ] Add SERVICE_TOKEN env di docker-compose
- [ ] Implement token validation middleware

---

## Phase 4: Event Handlers 🔄

### 4.1 Outbox Events Table (Sudah Ada)
- [✅] academic_service: outbox_events table
- [✅] finance_service: outbox_events table
- [✅] pmb_service: outbox_events table

### 4.2 Inbox Events Table (Sudah Ada)
- [✅] academic_service: inbox_events table
- [✅] finance_service: inbox_events table  
- [✅] pmb_service: inbox_events table

### 4.3 Event Publisher Service
- [ ] services/[service]/internal/service/event_publisher.go
- [ ] Publish method implementation

### 4.4 Event Consumer / Handler
- [ ] academic_service: Handle applicant.registered
- [ ] academic_service: Handle clearance.updated
- [ ] finance_service: Handle student.created
- [ ] finance_service: Handle grade.submitted

---

## Phase 5: API Routes 🔄

### 5.1 Core Service Routes
- [✅] POST /api/v1/auth/login
- [✅] POST /api/v1/auth/register
- [✅] GET /api/v1/persons/:id
- [ ] GET /api/v1/users/me
- [ ] PUT /api/v1/users/me

### 5.2 Reference Service Routes
- [✅] GET /api/v1/academic-years
- [✅] GET /api/v1/academic-periods
- [✅] GET /api/v1/study-programs
- [✅] GET /api/v1/regions

### 5.3 Academic Service Routes
- [✅] GET /api/v1/students/:nim
- [✅] GET /api/v1/krs
- [✅] POST /api/v1/krs
- [✅] GET /api/v1/grades
- [✅] POST /api/v1/grades
- [ ] POST /api/v1/enrollment

### 5.4 Finance Service Routes
- [✅] GET /api/v1/invoices
- [✅] POST /api/v1/invoices
- [✅] POST /api/v1/payments/confirm
- [✅] GET /api/v1/clearances/:student_id

### 5.5 PMB Service Routes
- [✅] POST /api/v1/applicants/register
- [✅] GET /api/v1/applicants/:id
- [✅] PUT /api/v1/applicants/:id/status
- [ ] GET /api/v1/waves/:wave_id/applicants

### 5.6 Input Validation
- [ ] Add request validation middleware
- [ ] Add field-level validation

---

## Phase 6: Frontend Integration 🔄

### 6.1 Context Providers
- [✅] auth-context.tsx
- [✅] reference-context.tsx
- [ ] service-context.tsx

### 6.2 Custom Hooks
- [✅] use-pmb.ts
- [✅] use-academic.ts
- [✅] use-finance.ts
- [ ] use-core.ts

### 6.3 API Client Configuration
- [✅] lib/api.ts - Base fetch
- [✅] lib/constants.ts - API endpoints
- [ ] Add error handling
- [ ] Add retry logic

### 6.4 Pages Integration
- [✅] Login page → Core API
- [✅] PMB registration → PMB API  
- [✅] Dashboard → Multiple APIs
- [ ] Student profile → Academic + Core
- [ ] KRS page → Academic + Finance

### 6.5 Loading & Error States
- [ ] Add skeleton loaders
- [ ] Add error boundaries
- [ ] Add toast notifications

---

## Priority Tasks

### P0 - Must Have
- [ ] Service client factory (packages/shared-serviceclient/)
- [ ] Finance → Academic clearance check flow
- [ ] Academic → Finance KRS submission flow
- [ ] PMB → Academic student creation flow

### P1 - Should Have  
- [ ] Event handlers implementation
- [ ] Service token validation
- [ ] API validation middleware

### P2 - Nice to Have
- [ ] Full frontend integration
- [ ] Loading states
- [ ] Error boundaries

---

## Notes

### Current Blocker
- Service token belum dikonfigurasi di docker-compose
- Event consumerbelum diimplementasi

### Workaround
- Gunakan basic auth untuk development
- Manual testing via Postman

---
