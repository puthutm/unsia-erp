# Complete API Endpoint Audit & Alignment Plan

This plan details all API path, port, and endpoint mismatches between the Go backend microservices and the Next.js frontend modules, and outlines the fixes to align them.

---

## User Review Required

> [!IMPORTANT]
> **Backend Port Shuffling**: Running the backend services locally outside Docker without specifying `PORT` environment variables leads to multiple port binding conflicts.
> We must align all backend service fallback port configurations in Go code with the planned ports mapped in `docker-compose.yml` to ensure seamless local debugging.

> [!IMPORTANT]
> **HRIS & Assessment Endpoint Mismatches**: The frontend hooks for HRIS and Assessment modules call endpoints that deviate from the backend paths defined in GORM handlers or contracts.
> We will add alias routing mapping in the backend to ensure backward compatibility and prevent `404 Not Found` errors in the frontend without breaking existing event/integration logic.

---

## Port Configurations Audit & Comparison

The table below outlines the ports mapped on the host machine in `docker-compose.yml`, the default fallbacks defined in the Go code, and the default fallback URLs in the Next.js frontend hooks.

| Service / Modul | Planned Port (Compose Host) | Backend Code Fallback | Frontend Default Hook | Port Mismatch Status & Impact |
| :--- | :---: | :---: | :---: | :--- |
| **Core Service** | `8001` | `8001` | `8001` | ✅ Aligned |
| **Reference Service** | `8002` | `8002` | `8002` | ✅ Aligned |
| **PMB Service** | `8003` | `8004` | `8003` | ⚠️ Mismatched (Code fallback `8004` conflicts with Academic Service) |
| **Academic Service** | `8004` | `8006` | `8004` | ⚠️ Mismatched (Code fallback `8006` conflicts with LMS Service) |
| **Finance Service** | `8005` | `8005` | `8005` | ✅ Aligned |
| **LMS Service** | `8006` | `8008` | `8006` (or `8081`) | ⚠️ Mismatched (Code fallback `8008` conflicts with HRIS; Hook may use `8081`) |
| **Assessment Service** | `8007` | `8009` | `8087` (or `8007`) | ⚠️ Mismatched (Code fallback `8009` conflicts with CRM; Hook uses `8087`!) |
| **HRIS Service** | `8008` | `8007` | `8008` | ⚠️ Mismatched (Code fallback `8007` conflicts with Assessment Service) |
| **CRM Service** | `8009` | `8003` | `8009` | ⚠️ Mismatched (Code fallback `8003` conflicts with PMB Service) |
| **Portal Service** | `8010` | `8010` | `8010` | ✅ Aligned |

---

## Endpoint Path & Routing Mismatches

| Modul Area | Frontend Expected Path | Backend Registered Path | Mismatch Impact | Resolution Action |
| :--- | :--- | :--- | :--- | :--- |
| **HRIS** | `GET/POST /api/v1/hris/attendance` | `GET/POST /api/v1/hris/attendances` | ❌ `404 Not Found` in Attendance page | Add alias route in HRIS service to map `/v1/hris/attendance` to `ListAttendances`/`RecordAttendance` |
| **HRIS** | `GET/POST /api/v1/hris/leave` | `GET/POST /api/v1/hris/leave-requests` | ❌ `404 Not Found` in Leave page | Add alias route in HRIS service to map `/v1/hris/leave` to `ListLeaveRequests`/`SubmitLeaveRequest` |
| **Assessment** | `POST /api/v1/assessment/sessions/:id/participants` | `POST /api/v1/assessment/participants` (body payload) | ❌ `404 Not Found` on registering participant | Add sub-resource route `POST /v1/assessment/sessions/:session_id/participants` in Assessment service |
| **Assessment** | `GET /api/v1/assessment/sessions/:id/participants` | `GET /api/v1/assessment/participants?assessment_session_id=...` | ❌ `404 Not Found` on viewing participant list | Add sub-resource route `GET /v1/assessment/sessions/:session_id/participants` in Assessment service |
| **LMS** | `GET /api/v1/lms/sessions` | `GET /api/v1/lms/classes/:id/sessions` | ❌ `404 Not Found` on general sessions view | Implement general session list handler `ListAllSessions` in LMS service |
| **Reference** | `GET /api/v1/reference/...` | `GET /api/v1/ref/...` | ❌ `404 Not Found` on central dropdowns | Map `/v1/reference` group as alias group in Reference service (Done) |

---

## Proposed Action Items

### 1. Fix Default Fallback Ports in Backend services
Update the port logic in the `main.go` file of each service to use the correct planned ports when the `PORT` environment variable is not defined:
- **PMB Service**: Change fallback port in `services/unsia-pmb-service/cmd/pmb-service/main.go` from `8004` to `8003`.
- **Academic Service**: Change fallback port in `services/unsia-academic-service/cmd/academic-service/main.go` from `8006` to `8004`.
- **LMS Service**: Change fallback port in `services/unsia-lms-service/cmd/lms-service/main.go` from `8008` to `8006`.
- **Assessment Service**: Change fallback port in `services/unsia-assessment-service/cmd/assessment-service/main.go` from `8009` to `8007`.
- **HRIS Service**: Change fallback port in `services/unsia-hris-service/cmd/hris-service/main.go` from `8007` to `8008`.
- **CRM Service**: Change fallback port in `services/unsia-crm-service/cmd/crm-service/main.go` from `8003` to `8009`.

### 2. Fix Frontend Hooks Configuration Ports
- **Assessment Frontend**: Update `frontend/unsia-assessment/hooks/use-assessment.ts` fallback port to `8007`.
- **LMS Frontend**: Ensure `frontend/unsia-lms/hooks/use-lms.ts` fallback port is `8006`.

### 3. Add Routing Aliases in HRIS Service (Port 8008)
Update `services/unsia-hris-service/cmd/hris-service/main.go` to add:
```go
// Attendance alias route (singular vs plural)
protected.GET("/v1/hris/attendance", middleware.PermissionRequired("hris.attendance.view"), hrisHandler.ListAttendances)
protected.POST("/v1/hris/attendance", hrisHandler.RecordAttendance)

// Leave alias route (singular vs plural leave-requests)
protected.GET("/v1/hris/leave", middleware.PermissionRequired("hris.leave.view"), hrisHandler.ListLeaveRequests)
protected.POST("/v1/hris/leave", hrisHandler.SubmitLeaveRequest)
```

### 4. Implement Sessions Sub-Resource routes in Assessment Service (Port 8007)
Update `services/unsia-assessment-service/cmd/assessment-service/main.go` to register:
```go
// Session-specific sub-resource participant handlers
protected.POST("/v1/assessment/sessions/:session_id/participants", assessHandler.RegisterSessionParticipant)
protected.GET("/v1/assessment/sessions/:session_id/participants", assessHandler.ListSessionParticipants)
```

Implement the corresponding wrappers in `assessment_handler.go`:
```go
func (h *AssessmentHandler) RegisterSessionParticipant(c *gin.Context) {
	sessionID := c.Param("session_id")
	var req ParticipantRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}
	req.AssessmentSessionID = sessionID
	
	// Delegate to original register logic
	part := domain.AssessmentParticipant{
		AssessmentSessionID: req.AssessmentSessionID,
		ParticipantType:     req.ParticipantType,
		ApplicantID:         req.ApplicantID,
		StudentID:           req.StudentID,
		UserID:              req.UserID,
		Status:              "registered",
	}

	if err := h.repo.RegisterParticipant(&part); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mendaftarkan peserta").WithContext(c))
		return
	}
	c.JSON(http.StatusCreated, sharederr.Success(part).WithContext(c))
}

func (h *AssessmentHandler) ListSessionParticipants(c *gin.Context) {
	sessionID := c.Param("session_id")
	// Query GORM repository directly using sessionID
	var participants []domain.AssessmentParticipant
	err := h.db.Where("assessment_session_id = ?", sessionID).Find(&participants).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data peserta").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(participants).WithContext(c))
}
```

### 5. Add Missing LMS General Sessions Endpoint (Port 8006)
Update `services/unsia-lms-service/cmd/lms-service/main.go` to map `/v1/lms/sessions` (Done in `plan_fikk_bug.md` list).

---

## Verification Plan

### Automated Tests
Run build validation across all backend services to confirm syntax correctness:
- `cd services/unsia-pmb-service && go build ./...`
- `cd services/unsia-academic-service && go build ./...`
- `cd services/unsia-lms-service && go build ./...`
- `cd services/unsia-assessment-service && go build ./...`
- `cd services/unsia-hris-service && go build ./...`
- `cd services/unsia-crm-service && go build ./...`

### Manual Verification
1. Spin up all containers using `docker-compose up -d`.
2. Access the LMS page, HRIS Attendance page, and Assessment page.
3. Validate that requests succeed without any `404 Not Found` or `Connection Refused` errors.
