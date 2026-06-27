package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)



// CheckClearance handles GET /api/v1/finance/clearances/check
func (h *FinanceHandler) CheckClearance(c *gin.Context) {
	studentID := c.Query("student_id")
	periodID := c.Query("academic_period_id")
	scope := c.Query("service_scope")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "student_id is required").WithContext(c))
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

	status := "cleared"
	reasons := []string{}

	if cl != nil {
		status = cl.Status
		if cl.Reason != nil && *cl.Reason != "" {
			reasons = append(reasons, *cl.Reason)
		}
	} else {
		// If no clearance record exists, let's verify if they have unpaid invoices for this period
		var unpaidCount int64
		h.db.Model(&domain.Invoice{}).
			Where("student_id = ? AND academic_period_id = ? AND status != 'paid'", studentID, periodID).
			Count(&unpaidCount)

		if unpaidCount > 0 {
			status = "blocked"
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
