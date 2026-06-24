package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-crm-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type PipelineHandler struct {
	repo *repository.CRMRepository
	db   *gorm.DB
}

func NewPipelineHandler(db *gorm.DB) *PipelineHandler {
	return &PipelineHandler{
		repo: repository.NewCRMRepository(db),
		db:   db,
	}
}

// GetPipeline - Mengambil pipeline CRM (tampilan Kanban)
func (h *PipelineHandler) GetPipeline(c *gin.Context) {
	channel := c.Query("channel")

	pipeline, err := h.repo.GetPipeline(channel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pipeline").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(pipeline).WithContext(c))
}

// MoveLeads - Memindahkan leads antar stage di pipeline
func (h *PipelineHandler) MoveLeads(c *gin.Context) {
	leadsID := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Validasi status transition
	validTransitions := map[string][]string{
		"new":        {"contacted", "lost"},
		"contacted":  {"interested", "new", "lost"},
		"interested": {"pendaftar", "contacted", "lost"},
		"pendaftar": {"converted", "interested", "lost"},
	}

	currentStatus, _ := h.repo.GetLeadsStatus(leadsID)
	allowed := validTransitions[currentStatus]
	isValid := false
	for _, s := range allowed {
		if s == req.Status {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_TRANSITION", "Status transition tidak valid").WithContext(c))
		return
	}

	if err := h.repo.UpdateLeadsStatus(leadsID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate status").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message":    "Leads dipindahkan",
		"leads_id":   leadsID,
		"new_status": req.Status,
	}).WithContext(c))
}

// GetPipelineStats - Mengambil statistik pipeline
func (h *PipelineHandler) GetPipelineStats(c *gin.Context) {
	stats, err := h.repo.GetPipelineStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil statistik").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(stats).WithContext(c))
}

// GetFunnelConversion - Mengambil funnel konversi
func (h *PipelineHandler) GetFunnelConversion(c *gin.Context) {
	period := c.DefaultQuery("period", "2026")

	funnel, err := h.repo.GetFunnelConversion(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil funnel").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(funnel).WithContext(c))
}
