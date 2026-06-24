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

type OpportunityHandler struct {
	repo *repository.CrmRepository
	db   *gorm.DB
}

func NewOpportunityHandler(db *gorm.DB) *OpportunityHandler {
	return &OpportunityHandler{
		repo: repository.NewCrmRepository(db),
		db:   db,
	}
}

// GET /api/v1/opportunities - List opportunities
func (h *OpportunityHandler) ListOpportunities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	campaignID := c.Query("campaign_id")

	opportunities, total, err := h.repo.ListOpportunities(page, limit, status, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data opportunity").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(opportunities).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/opportunities - Create opportunity
func (h *OpportunityHandler) CreateOpportunity(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		ContactID   string  `json:"contact_id"`
		CampaignID string  `json:"campaign_id"`
		Value      float64 `json:"value"`
		Stage      string  `json:"stage"` // lead, qualified, proposal, negotiation, closed_won, closed_lost
		Note       string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	opportunity := domain.Opportunity{
		Name:       req.Name,
		ContactID:  req.ContactID,
		CampaignID: req.CampaignID,
		Value:     req.Value,
		Stage:     req.Stage,
		Note:      req.Note,
		Status:    "active",
	}

	if err := h.repo.CreateOpportunity(&opportunity); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat opportunity").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(opportunity).WithContext(c))
}

// GET /api/v1/opportunities/:id - Get opportunity
func (h *OpportunityHandler) GetOpportunity(c *gin.Context) {
	opportunityID := c.Param("id")

	opportunity, err := h.repo.GetOpportunityByID(opportunityID)
	if err != nil || opportunity == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Opportunity tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(opportunity).WithContext(c))
}

// PUT /api/v1/opportunities/:id - Update opportunity
func (h *OpportunityHandler) UpdateOpportunity(c *gin.Context) {
	opportunityID := c.Param("id")

	var req struct {
		Name    string  `json:"name"`
		Value  float64 `json:"value"`
		Stage string  `json:"stage"`
		Status string  `json:"status"`
		Note  string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Value > 0 {
		updates["value"] = req.Value
	}
	if req.Stage != "" {
		updates["stage"] = req.Stage
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Note != "" {
		updates["note"] = req.Note
	}

	if err := h.repo.UpdateOpportunity(opportunityID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui opportunity").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Opportunity berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/opportunities/:id - Delete opportunity
func (h *OpportunityHandler) DeleteOpportunity(c *gin.Context) {
	opportunityID := c.Param("id")

	if err := h.repo.DeleteOpportunity(opportunityID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus opportunity").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Opportunity berhasil dihapus").WithContext(c))
}
