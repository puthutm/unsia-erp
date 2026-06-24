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

// ScholarshipHandler handles scholarship-related endpoints
type ScholarshipHandler struct {
	*FinanceHandler
}

// NewScholarshipHandler creates a new ScholarshipHandler
func NewScholarshipHandler(fh *FinanceHandler) *ScholarshipHandler {
	return &ScholarshipHandler{FinanceHandler: fh}
}

// GetScholarships handles GET /api/v1/finance/scholarships
func (h *ScholarshipHandler) GetScholarships(c *gin.Context) {
	filter := repository.ScholarshipFilter{
		IsActive: c.Query("is_active") == "true",
		Search:  c.Query("search"),
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
func (h *ScholarshipHandler) CreateScholarship(c *gin.Context) {
	var req ScholarshipCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var startDate, endDate *time.Time
	if req.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err == nil {
			startDate = &parsed
		}
	}
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err == nil {
			endDate = &parsed
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	scholarship := domain.Scholarship{
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		MaxAmount:     req.MaxAmount,
		StartDate:     startDate,
		EndDate:       endDate,
		IsActive:      isActive,
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
func (h *ScholarshipHandler) UpdateScholarship(c *gin.Context) {
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
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.DiscountType != nil {
		updates["discount_type"] = *req.DiscountType
	}
	if req.DiscountValue != nil {
		updates["discount_value"] = *req.DiscountValue
	}
	if req.MaxAmount != nil {
		updates["max_amount"] = *req.MaxAmount
	}
	if req.StartDate != nil {
		if parsed, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			updates["start_date"] = parsed
		}
	}
	if req.EndDate != nil {
		if parsed, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			updates["end_date"] = parsed
		}
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
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
func (h *ScholarshipHandler) DeleteScholarship(c *gin.Context) {
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
