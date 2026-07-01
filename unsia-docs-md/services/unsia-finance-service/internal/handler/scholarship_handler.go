package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)



// GetScholarships handles GET /api/v1/finance/scholarships
func (h *FinanceHandler) GetScholarships(c *gin.Context) {
	filter := repository.ScholarshipFilter{
		StudentID:       c.Query("student_id"),
		ScholarshipType: c.Query("scholarship_type"),
		Status:          c.Query("status"),
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

	result, err := h.repo.GetScholarships(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar scholarship").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreateScholarship handles POST /api/v1/finance/scholarships
func (h *FinanceHandler) CreateScholarship(c *gin.Context) {
	var req ScholarshipCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	scholarship := domain.Scholarship{
		StudentID:       req.StudentID,
		ScholarshipType: req.ScholarshipType,
		Amount:          req.Amount,
		Status:          status,
	}

	if err := h.repo.CreateScholarship(&scholarship); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan scholarship").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.scholarship.create",
		Module:       "finance",
		ResourceType: "scholarship",
		ResourceID:   scholarship.ID,
		NewValue:     scholarship,
	})

	c.JSON(http.StatusCreated, sharederr.Success(scholarship).WithContext(c))
}

// UpdateScholarship handles PUT /api/v1/finance/scholarships/:id
func (h *FinanceHandler) UpdateScholarship(c *gin.Context) {
	id := c.Param("id")
	var req ScholarshipUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	existing, err := h.repo.GetScholarshipByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Scholarship tidak ditemukan").WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.ScholarshipType != nil {
		updates["scholarship_type"] = *req.ScholarshipType
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.ApprovedBy != nil {
		updates["approved_by"] = *req.ApprovedBy
		now := time.Now()
		updates["approved_at"] = &now
	}

	if err := h.repo.UpdateScholarship(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui scholarship").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.scholarship.update",
		Module:       "finance",
		ResourceType: "scholarship",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Scholarship berhasil diperbarui").WithContext(c))
}

// DeleteScholarship handles DELETE /api/v1/finance/scholarships/:id
func (h *FinanceHandler) DeleteScholarship(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.repo.GetScholarshipByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Scholarship tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.DeleteScholarship(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus scholarship").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.scholarship.delete",
		Module:       "finance",
		ResourceType: "scholarship",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Scholarship berhasil dihapus").WithContext(c))
}
