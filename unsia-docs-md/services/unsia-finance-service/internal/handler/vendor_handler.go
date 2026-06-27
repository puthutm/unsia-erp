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



// GetVendors handles GET /api/v1/finance/vendors
func (h *FinanceHandler) GetVendors(c *gin.Context) {
	filter := repository.VendorFilter{
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

	result, err := h.repo.GetVendors(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar vendor").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreateVendor handles POST /api/v1/finance/vendors
func (h *FinanceHandler) CreateVendor(c *gin.Context) {
	var req VendorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	vendor := domain.Vendor{
		ID:             fmt.Sprintf("VND-%s", time.Now().Format("20060102150405")),
		VendorCode:    "VND-" + time.Now().Format("20060102150405"),
		VendorName:    req.VendorName,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		TaxNumber:    req.TaxNumber,
		IsActive:      true,
	}

	if err := h.repo.CreateVendor(&vendor); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan vendor").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.vendor.create",
		Module:       "finance",
		ResourceType: "vendor",
		ResourceID:   vendor.ID,
		NewValue:     vendor,
	})

	c.JSON(http.StatusCreated, sharederr.Success(vendor).WithContext(c))
}
