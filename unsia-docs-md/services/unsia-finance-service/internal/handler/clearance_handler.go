package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
	"github.com/unsia-erp/unsia-finance-service/internal/state_machine"
	"gorm.io/gorm"
)

type ClearanceCreateRequest struct {
	StudentID        string  `json:"student_id" binding:"required"`
	AcademicPeriodID string  `json:"academic_period_id" binding:"required"`
	ServiceScope     string  `json:"service_scope" binding:"required,oneof=registration krs graduation"`
	Status           string  `json:"status" binding:"required,oneof=BLOCKED CONDITIONAL CLEARED REVOKED"`
	Reason           *string `json:"reason"`
}

// CheckClearance handles GET /api/v1/finance/clearances/check
func (h *FinanceHandler) CheckClearance(c *gin.Context) {
	studentID := c.Query("student_id")
	periodID := c.Query("academic_period_id")
	scope := c.Query("service_scope")

	if studentID == "" {
		c.JSON(http.StatusNotFound, sharederr.Error("STUDENT_NOT_FOUND", "student_id is required").WithContext(c))
		return
	}
	if periodID == "" {
		c.JSON(http.StatusUnprocessableEntity, sharederr.Error("NO_ACTIVE_ACADEMIC_PERIOD", "academic_period_id is required").WithContext(c))
		return
	}
	if scope == "" {
		scope = "registration"
	} else if scope != "registration" && scope != "krs" && scope != "graduation" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "invalid service_scope. Must be registration, krs, or graduation").WithContext(c))
		return
	}

	cl, err := h.repo.GetStudentClearance(studentID, periodID, scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengecek clearance").WithContext(c))
		return
	}

	status := "CLEARED"
	reasons := []string{}

	if cl != nil {
		status = cl.Status
		if cl.Reason != nil && *cl.Reason != "" {
			reasons = append(reasons, *cl.Reason)
		}

		if status == "BLOCKED" {
			var disp domain.ClearanceDispensation
			errDisp := h.db.Where("student_clearance_id = ? AND status = 'approved' AND valid_until > ?", cl.ID, time.Now()).First(&disp).Error
			if errDisp == nil {
				status = "CONDITIONAL"
			}
		}
	} else {
		// If no clearance record exists, let's verify if they have unpaid invoices for this period
		var unpaidCount int64
		h.db.Model(&domain.Invoice{}).
			Where("student_id = ? AND academic_period_id = ? AND status != 'PAID'", studentID, periodID).
			Count(&unpaidCount)

		if unpaidCount > 0 {
			status = "BLOCKED"
			reasons = append(reasons, "Mahasiswa memiliki tagihan aktif yang belum lunas pada periode ini.")
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(ClearanceStatusResponse{
		StudentID:        studentID,
		AcademicPeriodID: periodID,
		ServiceCode:      scope,
		ClearanceStatus:  status,
		BlockReasons:     reasons,
	}).WithContext(c))
}

// CreateOrUpdateClearance handles POST /api/v1/finance/clearances
func (h *FinanceHandler) CreateOrUpdateClearance(c *gin.Context) {
	var req ClearanceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if req.Status == "REVOKED" {
		if req.Reason == nil || *req.Reason == "" {
			c.JSON(http.StatusUnprocessableEntity, sharederr.Error("REASON_REQUIRED", "Reason is required when revoking clearance").WithContext(c))
			return
		}
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	var existing domain.StudentClearance
	err := h.db.Where("student_id = ? AND academic_period_id = ? AND service_scope = ?", req.StudentID, req.AcademicPeriodID, req.ServiceScope).First(&existing).Error

	now := time.Now()
	statusChanged := false
	var oldStatus string

	if err == gorm.ErrRecordNotFound {
		statusChanged = true
		cl := domain.StudentClearance{
			StudentID:        req.StudentID,
			AcademicPeriodID: &req.AcademicPeriodID,
			ServiceScope:     req.ServiceScope,
			Status:           req.Status,
			Reason:           req.Reason,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		errSave := h.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&cl).Error; err != nil {
				return err
			}
			envelope := sharedevent.EventEnvelope{
				EventName:        "finance.clearance_changed",
				EventVersion:     "v1",
				PublisherService: "finance-service",
				AggregateType:    "clearance",
				AggregateID:      cl.ID,
				CorrelationID:    cid,
				Payload: map[string]interface{}{
					"student_id":         cl.StudentID,
					"academic_period_id": cl.AcademicPeriodID,
					"service_scope":      cl.ServiceScope,
					"previous_status":    "",
					"new_status":         cl.Status,
					"changed_at":         now,
				},
			}
			conn := tx.Statement.ConnPool
			_, err := sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
			return err
		})

		if errSave != nil {
			c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan clearance").WithContext(c))
			return
		}

		sharedaudit.Log(c, sharedaudit.AuditEntry{
			Action:       "finance.clearance.create",
			Module:       "finance",
			ResourceType: "clearance",
			ResourceID:   cl.ID,
			NewValue:     cl,
		})

		c.JSON(http.StatusCreated, sharederr.Success(cl).WithContext(c))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", err.Error()).WithContext(c))
		return
	}

	oldStatus = existing.Status
	if oldStatus != req.Status {
		statusChanged = true
	}

	if statusChanged {
		sm := state_machine.NewClearanceStateMachine()
		if errSM := sm.ValidateTransitionWithReason(state_machine.ClearanceStatus(oldStatus), state_machine.ClearanceStatus(req.Status), req.Reason); errSM != nil {
			c.JSON(http.StatusConflict, sharederr.Error("INVALID_CLEARANCE_STATUS_TRANSITION", errSM.Error()).WithContext(c))
			return
		}
	}

	existing.Status = req.Status
	existing.Reason = req.Reason
	existing.UpdatedAt = now

	errUpdate := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		if statusChanged {
			envelope := sharedevent.EventEnvelope{
				EventName:        "finance.clearance_changed",
				EventVersion:     "v1",
				PublisherService: "finance-service",
				AggregateType:    "clearance",
				AggregateID:      existing.ID,
				CorrelationID:    cid,
				Payload: map[string]interface{}{
					"student_id":         existing.StudentID,
					"academic_period_id": existing.AcademicPeriodID,
					"service_scope":      existing.ServiceScope,
					"previous_status":    oldStatus,
					"new_status":         existing.Status,
					"changed_at":         now,
				},
			}
			conn := tx.Statement.ConnPool
			_, err := sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
			return err
		}
		return nil
	})

	if errUpdate != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui clearance").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.clearance.update",
		Module:       "finance",
		ResourceType: "clearance",
		ResourceID:   existing.ID,
		OldValue:     oldStatus,
		NewValue:     existing,
	})

	c.JSON(http.StatusOK, sharederr.Success(existing).WithContext(c))
}

// CreateClearancePolicy handles POST /api/v1/finance/clearance-policies
func (h *FinanceHandler) CreateClearancePolicy(c *gin.Context) {
	var req ClearancePolicyCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	policy := domain.ClearancePolicy{
		Code:         req.Code,
		Name:         req.Name,
		ServiceScope: req.ServiceScope,
		RuleJson:     req.RuleJson,
		IsActive:     isActive,
	}

	if err := h.repo.CreateClearancePolicy(&policy); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat kebijakan clearance").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.clearance_policy.create",
		Module:       "finance",
		ResourceType: "clearance_policy",
		ResourceID:   policy.ID,
		NewValue:     policy,
	})

	c.JSON(http.StatusCreated, sharederr.Success(policy).WithContext(c))
}

// UpdateClearancePolicy handles PUT /api/v1/finance/clearance-policies/:id
func (h *FinanceHandler) UpdateClearancePolicy(c *gin.Context) {
	id := c.Param("id")
	var req ClearancePolicyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	policy, err := h.repo.GetClearancePolicyByID(id)
	if err != nil || policy == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kebijakan clearance tidak ditemukan").WithContext(c))
		return
	}

	if req.Name != "" {
		policy.Name = req.Name
	}
	if req.ServiceScope != "" {
		policy.ServiceScope = req.ServiceScope
	}
	if req.RuleJson != "" {
		policy.RuleJson = req.RuleJson
	}
	if req.IsActive != nil {
		policy.IsActive = *req.IsActive
	}

	if err := h.repo.UpdateClearancePolicy(policy); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui kebijakan clearance").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.clearance_policy.update",
		Module:       "finance",
		ResourceType: "clearance_policy",
		ResourceID:   policy.ID,
		NewValue:     policy,
	})

	c.JSON(http.StatusOK, sharederr.Success(policy).WithContext(c))
}

// GetClearances handles GET /api/v1/finance/clearances
func (h *FinanceHandler) GetClearances(c *gin.Context) {
	filter := repository.ClearanceListFilter{
		Status:           c.Query("status"),
		StudentID:       c.Query("student_id"),
		AcademicPeriodID: c.Query("academic_period_id"),
		ServiceScope:    c.Query("service_scope"),
	}

	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	filter.Page = page
	filter.Limit = limit

	result, err := h.repo.GetClearances(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar clearance").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}
