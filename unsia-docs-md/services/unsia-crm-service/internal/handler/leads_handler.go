package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-crm-service/internal/infrastructure/repository"
	"github.com/unsia-erp/unsia-crm-service/internal/domain"
	"gorm.io/gorm"
)

type LeadsCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Channel     string `json:"channel" binding:"required"`
	ProgramID   string `json:"program_id"`
	ReferrerID  string `json:"referrer_id"`
	Source      string `json:"source"`
}

type LeadsUpdateRequest struct {
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	FollowUpAt  *time.Time `json:"follow_up_at"`
	AssignedTo  string `json:"assigned_to"`
}

type LeadsHandler struct {
	repo *repository.CRMRepository
	db   *gorm.DB
}

func NewLeadsHandler(db *gorm.DB) *LeadsHandler {
	return &LeadsHandler{
		repo: repository.NewCRMRepository(db),
		db:   db,
	}
}

// GetLeads - Mengambil semua leads dengan filter
func (h *LeadsHandler) GetLeads(c *gin.Context) {
	channel := c.Query("channel")
	status := c.Query("status")
	search := c.Query("search")

	leads, err := h.repo.GetLeads(channel, status, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data leads").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(leads).WithContext(c))
}

// GetLeadsByID - Mengambil leads berdasarkan ID
func (h *LeadsHandler) GetLeadsByID(c *gin.Context) {
	id := c.Param("id")

	leads, err := h.repo.GetLeadsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Leads tidak найден").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(leads).WithContext(c))
}

// CreateLeads - Membuat leads baru
func (h *LeadsHandler) CreateLeads(c *gin.Context) {
	var req LeadsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	leads := domain.Leads{
		ID:          generateLeadID(req.Channel),
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Channel:     req.Channel,
		ProgramID:   req.ProgramID,
		ReferrerID:  req.ReferrerID,
		Source:      req.Source,
		Status:      "new",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.CreateLeads(&leads); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat leads").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(leads).WithContext(c))
}

// UpdateLeads - Mengupdate status dan info leads
func (h *LeadsHandler) UpdateLeads(c *gin.Context) {
	id := c.Param("id")
	var req LeadsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	leads, err := h.repo.GetLeadsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Leads tidak найден").WithContext(c))
		return
	}

	if req.Status != "" {
		leads.Status = req.Status
	}
	if req.Notes != "" {
		leads.Notes = req.Notes
	}
	if req.FollowUpAt != nil {
		leads.FollowUpAt = req.FollowUpAt
	}
	if req.AssignedTo != "" {
		leads.AssignedTo = req.AssignedTo
	}
	leads.UpdatedAt = time.Now()

	if err := h.repo.UpdateLeads(leads); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate leads").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(leads).WithContext(c))
}

// ConvertToApplicant - Mengkonversi leads menjadi pendaftar PMB
func (h *LeadsHandler) ConvertToApplicant(c *gin.Context) {
	id := c.Param("id")

	leads, err := h.repo.GetLeadsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Leads tidak найден").WithContext(c))
		return
	}

	// Update status leads
	leads.Status = "converted"
	leads.UpdatedAt = time.Now()

	if err := h.repo.UpdateLeads(leads); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengkonversi leads").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message":   "Leads berhasil dikonversi ke pendaftar",
		"leads_id":  id,
		"applicant_id": generateApplicantID(),
	}).WithContext(c))
}

// GetLeadsStats - Mengambil statistik leads
func (h *LeadsHandler) GetLeadsStats(c *gin.Context) {
	stats, err := h.repo.GetLeadsStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil statistik").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(stats).WithContext(c))
}

func generateLeadID(channel string) string {
	prefix := map[string]string{
		"umum":      "LEAD-UM",
		"agen":      "LEAD-AR",
		"bts":       "LEAD-BTS",
		"sgs":       "LEAD-SGS",
		"egs":       "LEAD-EGS",
		"kerjasama": "LEAD-KS",
	}
	return prefix[channel] + "-" + time.Now().Format("2006-") + fmtRandom(4)
}

func generateApplicantID() string {
	return "APP-" + time.Now().Format("2006") + fmtRandom(4)
}

func fmtRandom(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += string(rune('0' + rand.Intn(10)))
	}
	return result
}
