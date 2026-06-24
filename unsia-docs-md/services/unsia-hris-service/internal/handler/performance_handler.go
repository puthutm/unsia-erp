package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// BKD = Beban Kerja Dosen (Teaching, Research, Service Load)
type BkdRecordCreateRequest struct {
	LecturerID        string  `json:"lecturer_id" binding:"required"`
	AcademicPeriodID  *string `json:"academic_period_id"`
	TeachingLoad     float64 `json:"teaching_load"`
	ResearchLoad     float64 `json:"research_load"`
	ServiceLoad      float64 `json:"service_load"`
	TeachingCredits  float64 `json:"teaching_credits"`  // SKS
	ResearchCredits  float64 `json:"research_credits"` // SKS
	ServiceCredits   float64 `json:"service_credits"`  // SKS
}

type BkdApprovalRequest struct {
	Status string `json:"status" binding:"required"` // approved, rejected
	Notes  string `json:"notes"`
}

// PerformanceReviewRequest struct
type PerformanceReviewRequest struct {
	EmployeeID      string  `json:"employee_id" binding:"required"`
	ReviewPeriodID  string  `json:"review_period_id" binding:"required"`
	Rating         float64 `json:"rating"` // 1-100
	Strengths      string  `json:"strengths"`
	Improvements   string  `json:"improvements"`
	Goals          string  `json:"goals"`
}

type PerformanceHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewPerformanceHandler(db *gorm.DB) *PerformanceHandler {
	return &PerformanceHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

// ===== BKD (BEBAN KERJA DOSEN) =====

// ListBkdRecords - GET /api/v1/bkd-records
func (h *PerformanceHandler) ListBkdRecords(c *gin.Context) {
	lecturerID := c.Query("lecturer_id")
	periodID := c.Query("academic_period_id")
	status := c.Query("status")

	records, err := h.repo.ListBkdRecords(lecturerID, periodID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data BKD").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(records).WithContext(c))
}

// GetBkdRecord - GET /api/v1/bkd-records/:id
func (h *PerformanceHandler) GetBkdRecord(c *gin.Context) {
	id := c.Param("id")
	record, err := h.repo.GetBkdRecordByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data BKD").WithContext(c))
		return
	}
	if record == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Record BKD tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(record).WithContext(c))
}

// CreateBkdRecord - POST /api/v1/bkd-records (Dosen submits)
func (h *PerformanceHandler) CreateBkdRecord(c *gin.Context) {
	var req BkdRecordCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Calculate total SKS
	totalSKS := req.TeachingCredits + req.ResearchCredits + req.ServiceCredits

	// Check if exceeds 12 SKS rule
	if totalSKS > 12 {
		c.JSON(http.StatusBadRequest, sharederr.Error("EXCEEDS_LIMIT", "Total SKS melebihi batas maksimal 12 SKS").WithContext(c))
		return
	}

	bkd := domain.BkdRecord{
		LecturerID:        req.LecturerID,
		AcademicPeriodID:  req.AcademicPeriodID,
		TeachingLoad:    req.TeachingLoad,
		ResearchLoad:    req.ResearchLoad,
		ServiceLoad:     req.ServiceLoad,
		TeachingCredits:  req.TeachingCredits,
		ResearchCredits:  req.ResearchCredits,
		ServiceCredits:  req.ServiceCredits,
		TotalCredits:    totalSKS,
		Status:          "draft", // draft, submitted, approved, rejected
	}

	if err := h.repo.CreateBkdRecord(&bkd); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan BKD").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(bkd).WithContext(c))
}

// SubmitBkdRecord - PUT /api/v1/bkd-records/:id/submit (Dosen submits for approval)
func (h *PerformanceHandler) SubmitBkdRecord(c *gin.Context) {
	id := c.Param("id")
	bkd, err := h.repo.GetBkdRecordByID(id)
	if err != nil || bkd == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Record BKD tidak ditemukan").WithContext(c))
		return
	}

	if bkd.Status != "draft" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_STATUS", "Hanya draft yang dapat diajukan").WithContext(c))
		return
	}

	bkd.Status = "submitted"
	if err := h.repo.UpdateBkdRecord(bkd); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengajukan BKD").WithContext(c))
		return
	}

	// TODO: Notify head of study program
	c.JSON(http.StatusOK, sharederr.Success(bkd).WithContext(c))
}

// ApproveBkdRecord - PUT /api/v1/bkd-records/:id/approve (Head approves)
func (h *PerformanceHandler) ApproveBkdRecord(c *gin.Context) {
	id := c.Param("id")
	var req BkdApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	bkd, err := h.repo.GetBkdRecordByID(id)
	if err != nil || bkd == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Record BKD tidak ditemukan").WithContext(c))
		return
	}

	bkd.Status = req.Status
	bkd.Notes = req.Notes

	if err := h.repo.UpdateBkdRecord(bkd); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyetujui BKD").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(bkd).WithContext(c))
}

// GetMyBkdRecords - GET /api/v1/bkd-records/my (Dosen view own)
func (h *PerformanceHandler) GetMyBkdRecords(c *gin.Context) {
	employeeID := c.GetHeader("X-Employee-ID")
	if employeeID == "" {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Employee ID diperlukan").WithContext(c))
		return
	}

	records, err := h.repo.ListBkdRecords(employeeID, "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data BKD").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(records).WithContext(c))
}

// ===== PERFORMANCE REVIEW =====

// ListPerformanceReviews - GET /api/v1/performance-reviews
func (h *PerformanceHandler) ListPerformanceReviews(c *gin.Context) {
	employeeID := c.Query("employee_id")
	periodID := c.Query("review_period_id")

	reviews, err := h.repo.ListPerformanceReviews(employeeID, periodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kinerja").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(reviews).WithContext(c))
}

// CreatePerformanceReview - POST /api/v1/performance-reviews (Manager creates)
func (h *PerformanceHandler) CreatePerformanceReview(c *gin.Context) {
	var req PerformanceReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	review := domain.PerformanceReview{
		EmployeeID:     req.EmployeeID,
		ReviewPeriodID:  req.ReviewPeriodID,
		Rating:        req.Rating,
		Strengths:     req.Strengths,
		Improvements:  req.Improvements,
		Goals:         req.Goals,
	}

	if err := h.repo.CreatePerformanceReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan kinerja").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(review).WithContext(c))
}

// GetPerformanceSummary - GET /api/v1/performance-reviews/summary
func (h *PerformanceHandler) GetPerformanceSummary(c *gin.Context) {
	periodID := c.Query("review_period_id")
	workUnitID := c.Query("work_unit_id")

	summary, err := h.repo.GetPerformanceSummary(periodID, workUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil ringkasan kinerja").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(summary).WithContext(c))
}
