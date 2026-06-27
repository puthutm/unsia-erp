# BE-FE-DB Integration Plan

Rencana integrasi Backend (Go microservices), Frontend (Next.js), dan Database (PostgreSQL) untuk UNSIA ERP.

## Status Overview

| Phase | Component | Status | Priority |
|-------|-----------|--------|----------|
| 1 | Database Setup | ✅ Complete | - |
| 2 | HTTP Client Utility | ✅ Complete | - |
| 3 | Service Clients | 🔄 In Progress | High |
| 4 | Event Handlers | 🔄 In Progress | High |
| 5 | API Routes | 🔄 In Progress | High |
| 6 | Frontend Integration | 🔄 In Progress | High |

---

## Phase 1: Database Infrastructure ✅

### Independent Databases
 Semua service menggunakan database terpisah untuk isolation:

```
core_db          → unsia-core-service (port 8001)
reference_db    → unsia-reference-service (port 8002)  
pmb_db           → unsia-pmb-service (port 8003)
academic_db     → unsia-academic-service (port 8004)
finance_db       → unsia-finance-service (port 8005)
lms_db           → unsia-lms-service (port 8006)
hris_db          → unsia-hris-service (port 8008)
assessment_db   → unsia-assessment-service (port 8007)
crm_db          → unsia-crm-service (port 8009)
portal_db       → unsia-portal-service (port 8010)
```

### Cross-Service References
Field references antar databases menggunakan UUID dengan format: `external_ref: service.table.id`

## Phase 2: HTTP Communication ✅

### shared-httpclient Package
Location: `packages/shared-httpclient/client.go`

Features:
- Auto-retry with exponential backoff
- Circuit breaker pattern
- Service token authentication
- Correlation ID propagation

```go
// Usage example
client := sharedhttpclient.New(sharedhttpclient.Config{
    BaseURL:      "http://academic-service:8004",
    ServiceToken: "service-token-here",
    SourceName:   "finance-service", 
    Timeout:     10 * time.Second,
    MaxRetries:  3,
})
```

---

## Phase 3: Service Clients Implementation 🔄

### Architecture
Setiap service memiliki package `internal/client/` yang menyediakan typed client untuk interaksi dengan service lain:

```
services/[service]/internal/client/
├── academic_client.go    # calls to academic-service
├── core_client.go      # calls to core-service  
├── finance_client.go  # calls to finance-service
├── reference_client.go # calls to reference-service
└── ...
```

### Service URL Environment Variables
Diambil dari docker-compose.yml:

| Service | ENV Variable | Default URL |
|--------|------------|-------------|
| Core | CORE_SERVICE_URL | http://unsia-core-service:8001 |
| Reference | REFERENCE_SERVICE_URL | http://unsia-reference-service:8002 |
| PMB | PMB_SERVICE_URL | http://unsia-pmb-service:8003 |
| Academic | ACADEMIC_SERVICE_URL | http://unsia-academic-service:8004 |
| Finance | FINANCE_SERVICE_URL | http://unsia-finance-service:8005 |
| LMS | LMS_SERVICE_URL | http://unsia-lms-service:8006 |
| HRIS | HRIS_SERVICE_URL | http://unsia-hris-service:8008 |

### Implementation Tasks

#### 3.1 Academic Service Clients
Priority: HIGH
- `client/academic_client.go` - GetStudentByID, GetKRS, SubmitGrade
- `client/finance_client.go` - CheckClearance, GetInvoice  
- `client/reference_client.go` - GetPeriod, GetStudyProgram

#### 3.2 Finance Service Clients
Priority: HIGH
- `client/academic_client.go` - GetStudentInfo
- `client/reference_client.go` - GetAcademicPeriod

#### 3.3 PMB Service Clients
Priority: MEDIUM
- `client/core_client.go` - CreateUserPerson
- `client/academic_client.go` - CreateStudent

#### 3.4 Generic Service Client Factory
Location: `packages/shared-serviceclient/`

```go
// Factory pattern example
func NewServiceClient(serviceName string, cfg Config) (*ServiceClient, error)

client, err := NewServiceClient("academic", Config{
    BaseURL: os.Getenv("ACADEMIC_SERVICE_URL"),
    Token: os.Getenv("SERVICE_TOKEN"),
})
```

---

## Phase 4: Event-Driven Integration 🔄

### Outbox Pattern Implementation
Setiap service sudah punya `outbox_events` table. Perlu implement:

### 4.1 Event Publisher Service
Location: `services/[service]/internal/service/event_publisher.go`

```go
type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
}

type Event struct {
    ID        string    `json:"id"`
    Type     string    `json:"type"`        // "student.created", "clearance.updated"
    Source   string    `json:"source"`      // "academic-service"
    Payload  json.RawMessage `json:"payload"`
    Metadata Metadata  `json:"metadata"`
}
```

### 4.2 Event Subscriber / Inbox Handler
Location: `services/[service]/cmd/events/consumer.go`

Events yang perlu di-handle:

| Event Type | Producer | Consumers |
|-----------|----------|----------|
| `applicant.registered` | PMB | Academic, Core |
| `student.created` | Academic | Finance, LMS, Core |
| `clearance.updated` | Finance | Academic |
| `grade.submitted` | Academic | LMS |
| `enrollment.completed` | LMS | Academic, Finance |
| `employee.registered` | HRIS | Core, Finance |
| `user.activated` | Core | All services |

### 4.3 Inbox Event Handler Template

```go
// Example: services/unsia-academic-service/cmd/events/consumer.go
func (c *Consumer) HandleInboxEvents(ctx context.Context) error {
    events, err := c.repo.GetPendingInboxEvents(ctx, "academic-service")
    if err != nil {
        return err
    }
    
    for _, event := range events {
        switch event.Type {
        case "applicant.registered":
            return c.handleApplicantRegistered(ctx, event)
        case "clearance.updated":
            return c.handleClearanceUpdated(ctx, event)
        }
    }
    return nil
}
```

---

## Phase 5: API Routes 🔄

### Frontend API Integration
Frontend menggunakan environment variables dari docker-compose.yml:

```typescript
// unsia-portal-web/lib/constants.ts
export const API_ENDPOINTS = {
  CORE: process.env.NEXT_PUBLIC_API_CORE_URL || 'http://localhost:8001',
  REFERENCE: process.env.NEXT_PUBLIC_API_REFERENCE_URL || 'http://localhost:8002',
  ACADEMIC: process.env.NEXT_PUBLIC_API_ACADEMIC_URL || 'http://localhost:8004',
  FINANCE: process.env.NEXT_PUBLIC_API_FINANCE_URL || 'http://localhost:8005',
  // ...
}
```

### API Route Structure

#### 5.1 Core Service Routes
- `POST /api/v1/auth/login` - Local auth
- `POST /api/v1/auth/refresh` - Refresh token  
- `GET /api/v1/persons/:id` - Get person by ID
- `GET /api/v1/roles` - List roles
- `GET /api/v1/permissions` - List permissions

#### 5.2 Reference Service Routes
- `GET /api/v1/academic-years` - List akademik years
- `GET /api/v1/academic-periods` - List periods
- `GET /api/v1/study-programs` - List prodi

#### 5.3 Academic Service Routes
- `GET /api/v1/students/:nim` - Get student
- `GET /api/v1/krs` - List KRS
- `POST /api/v1/krs` - Create KRS
- `GET /api/v1/grades` - List grades

#### 5.4 Finance Service Routes
- `GET /api/v1/invoices` - List invoices
- `POST /api/v1/invoices` - Create invoice
- `POST /api/v1/payments/confirm` - Confirm payment
- `GET /api/v1/clearances/:studentId` - Get clearance

### 5.5 API Response Format

```json
// Standard success response
{
  "success": true,
  "data": {},
  "meta": {
    "request_id": "uuid",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}

// Standard error response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": []
  }
}
```

---

## Phase 6: Frontend Integration 🔄

### 6.1 React Contexts
Location: `frontend/unsia-portal-web/contexts/`

```typescript
// contexts/auth-context.tsx
// - Authentication state management
// - Login/logout functions
// - Token refresh

// contexts/service-context.tsx  
// - Service client initialization
// - Provides typed API clients
```

### 6.2 Custom Hooks
Location: `frontend/unsia-portal-web/hooks/`

```typescript
// hooks/use-pmb.ts
// - registerApplicant()
// - getApplicantStatus()

// hooks/use-academic.ts
// - getStudentProfile()  
// - getKRS()
// - submitKRS()

// hooks/use-finance.ts
// - getInvoices()
// - getClearance()
// - processPayment()
```

### 6.3 Data Fetching Pattern

```typescript
// Using React Query + custom hooks
function StudentProfile({ nim }: { nim: string }) {
  const { student, isLoading, error } = useAcademic((api) => api.getStudent(nim))
  
  if (isLoading) return <Skeleton />
  if (error) return <ErrorDisplay error={error} />
  
  return (
    <div>
      <h1>{student.name}</h1>
      <p>{student.program}</p>
    </div>
  )
}
```

---

## Implementation Order

### Priority 1: Core Service Clients (Week 1)
- [ ] Create `packages/shared-serviceclient/` factory
- [ ] Academic service → finance client (clearance check)
- [ ] Finance service → academic client (student info)
- [ ] PMB service → core client (user creation)
- [ ] Update docker-compose with SERVICE_TOKEN env vars

### Priority 2: Event Handlers (Week 2)  
- [ ] Implement inbox event consumer in academic-service
- [ ] Implement inbox event consumer in finance-service
- [ ] Add event publisher to PMB (applicant.registered)
- [ ] Test cross-service event flow

### Priority 3: API Routes (Week 3)
- [ ] Verify all service GET routes
- [ ] Add POST/PUT routes as needed
- [ ] Add input validation
- [ ] Add rate limiting

### Priority 4: Frontend Integration (Week 4)
- [ ] Update context providers
- [ ] Add all custom hooks
- [ ] Connect CRUD pages to API
- [ ] Add loading/error states

---

## Testing Strategy

### Unit Tests
- Service client methods
- Event handler logic
- Repository methods

### Integration Tests
- Service-to-service HTTP calls
- Database cross-references

### E2E Tests
- Login → Dashboard flow
- Registration → Student creation flow
- KRS submission → Finance clearance flow

---

## Rollback Plan

If integration fails:

1. **Database**: Each service has independent DB - no rollback needed
2. **HTTP**: Circuit breaker provides degradation - services remain operational
3. **Events**: Failed events remain in outbox for retry
4. **Frontend**: Can fall back to internal service calls

---

## Future Enhancements

1. **gRPC** - For faster internal communication
2. **GraphQL** - For flexible frontend queries  
3. **Graph Database** - For complex relationship queries
4. **CQRS** - For read/write separation
