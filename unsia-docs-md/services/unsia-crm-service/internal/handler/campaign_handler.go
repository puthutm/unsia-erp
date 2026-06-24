package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-crm-service/internal/domain"
	"github.com/unsia-erp/unsia-crm-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type CampaignHandler struct {
	repo *repository.CrmRepository
	db   *gorm.DB
}

func NewCampaignHandler(db *gorm.DB) *CampaignHandler {
	return &CampaignHandler{
		repo: repository.NewCrmRepository(db),
		db:   db,
	}
}

// GET /api/v1/campaigns - List campaigns
func (h *CampaignHandler) ListCampaigns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	campaigns, total, err := h.repo.ListCampaigns(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kampanye").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(campaigns).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/campaigns - Create campaign
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		StartDate  string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Budget    float64 `json:"budget"`
		Status    string `json:"status"` // draft, active, paused, completed
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	campaign := domain.Campaign{
		Name:         req.Name,
		Description:  req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Budget:      req.Budget,
		Status:      req.Status,
	}

	if err := h.repo.CreateCampaign(&campaign); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat kampanye").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(campaign).WithContext(c))
}

// GET /api/v1/campaigns/:id - Get campaign
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	campaignID := c.Param("id")

	campaign, err := h.repo.GetCampaignByID(campaignID)
	if err != nil || campaign == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kampanye tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(campaign).WithContext(c))
}

// PUT /api/v1/campaigns/:id - Update campaign
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	campaignID := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		StartDate  string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Budget    float64 `json:"budget"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.StartDate != "" {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		updates["end_date"] = req.EndDate
	}
	if req.Budget > 0 {
		updates["budget"] = req.Budget
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.repo.UpdateCampaign(campaignID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui kampanye").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Kampanye berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/campaigns/:id - Delete campaign
func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	campaignID := c.Param("id")

	if err := h.repo.DeleteCampaign(campaignID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus kampanye").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Kampanye berhasil dihapus").WithContext(c))
}
