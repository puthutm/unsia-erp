package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-pmb-service/internal/domain"
	"github.com/unsia-erp/unsia-pmb-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Wave Handlers ============

type PmbWaveCreateRequest struct {
	WaveName         string  `json:"wave_name" binding:"required"`
	AcademicPeriodID string  `json:"academic_period_id" binding:"required"`
	StartDate       string  `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate        string  `json:"end_date" binding:"required"`   // YYYY-MM-DD
	RegistrationFee *int    `json:"registration_fee"`
	IsActive       *bool   `json:"is_active"`
}

type PmbWaveUpdateRequest struct {
	WaveName         *string `json:"wave_name"`
	AcademicPeriodID *string `json:"academic_period_id"`
	StartDate        *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
	RegistrationFee *int    `json:"registration_fee"`
	IsActive        *bool   `json:"is_active"`
}

type WaveHandler struct {
	repo *repository.PMBRepository
	db   *gorm.DB
}

func NewWaveHandler(db *gorm.DB) *WaveHandler {
	return &WaveHandler{
		repo: repository.NewPMBRepository(db),
		db:   db,
	}
}

// CreateWave handles POST /api/v1/pmb-waves
func (h *WaveHandler) CreateWave(c *gin.Context) {
	var req PmbWaveCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	wave := domain.PmbWave{
		WaveName:         req.WaveName,
		AcademicPeriodID: req.AcademicPeriodID,
		StartDate:       startDate,
		EndDate:        endDate,
		RegistrationFee: req.RegistrationFee,
		IsActive:       isActive,
	}

	if err := h.db.Create(&wave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan gelombang").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.wave.create",
		Module:      "pmb",
		ResourceType: "pmb_wave",
		ResourceID:  wave.ID,
		NewValue:    wave,
	})

	c.JSON(http.StatusCreated, sharederr.Success(wave).WithContext(c))
}

// GetWave handles GET /api/v1/pmb-waves/:id
func (h *WaveHandler) GetWave(c *gin.Context) {
	id := c.Param("id")
	wave, err := h.repo.GetWaveByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil gelombang").WithContext(c))
		return
	}
	if wave == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Gelombang tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(wave).WithContext(c))
}

// GetWaves handles GET /api/v1/pmb-waves
func (h *WaveHandler) GetWaves(c *gin.Context) {
	filter := repository.WaveFilter{
		AcademicPeriodID: c.Query("academic_period_id"),
		IsActive:        c.Query("is_active") == "true",
		Page:            1,
		Limit:           20,
	}

	var page, limit int
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	filter.Page = page
	filter.Limit = limit

	waves, total, err := h.repo.GetWaves(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar gelombang").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"data":  waves,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	}).WithContext(c))
}

// UpdateWave handles PUT /api/v1/pmb-waves/:id
func (h *WaveHandler) UpdateWave(c *gin.Context) {
	id := c.Param("id")
	var req PmbWaveUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	wave, err := h.repo.GetWaveByID(id)
	if err != nil || wave == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Gelombang tidak ditemukan").WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.WaveName != nil {
		updates["wave_name"] = *req.WaveName
	}
	if req.AcademicPeriodID != nil {
		updates["academic_period_id"] = *req.AcademicPeriodID
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
	if req.RegistrationFee != nil {
		updates["registration_fee"] = *req.RegistrationFee
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.repo.UpdateWave(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui gelombang").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.wave.update",
		Module:      "pmb",
		ResourceType: "pmb_wave",
		ResourceID:  id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Gelombang berhasil diperbarui").WithContext(c))
}

// DeleteWave handles DELETE /api/v1/pmb-waves/:id
func (h *WaveHandler) DeleteWave(c *gin.Context) {
	id := c.Param("id")

	wave, err := h.repo.GetWaveByID(id)
	if err != nil || wave == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Gelombang tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.DeleteWave(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus gelombang").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.wave.delete",
		Module:      "pmb",
		ResourceType: "pmb_wave",
		ResourceID:  id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Gelombang berhasil dihapus").WithContext(c))
}

// GetActiveWaves handles GET /api/v1/pmb-waves/active
func (h *WaveHandler) GetActiveWaves(c *gin.Context) {
	waves, err := h.repo.GetActiveWaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil gelombang aktif").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(waves).WithContext(c))
}

// GetWaveStatistics handles GET /api/v1/pmb-waves/:id/statistics
func (h *WaveHandler) GetWaveStatistics(c *gin.Context) {
	id := c.Param("id")

	stats, err := h.repo.GetWaveStatistics(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil statistik gelombang").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(stats).WithContext(c))
}
