package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Grade Handler ============

type GradeHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

func NewGradeHandler(db *gorm.DB) *GradeHandler {
	return &GradeHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Grade Request Types ============

type GradeComponent struct {
	Name       string  `json:"name" binding:"required"`
	Weight    float64 `json:"weight" binding:"required"`
	MaxScore  float64 `json:"max_score" binding:"required"`
}

type GradeCreateRequest struct {
	KrsItemID   string          `json:"krs_item_id" binding:"required"`
	Source     string          `json:"source" binding:"required,oneof=lms exam quiz assignment"`
	Components []GradeComponent `json:"components" binding:"gt=0"`
}

type GradeEntryRequest struct {
	StudentID   string                 `json:"student_id" binding:"required"`
	KrsItemID   string                 `json:"krs_item_id" binding:"required"`
	Components []GradeComponentScore    `json:"components" binding:"gt=0"`
	FinalGrade *string                 `json:"final_grade"`
	LetterGrade *string                `json:"letter_grade"`
	Status     string                 `json:"status" binding:"required,oneof=in_progress submitted final"`
}

type GradeComponentScore struct {
	Name   string  `json:"name" binding:"required"`
	Score float64 `json:"score"`
}

type GradeConversionRequest struct {
	GradeLetter string  `json:"grade_letter" binding:"required"`
	MinScore    float64 `json:"min_score"`
	MaxScore    float64 `json:"max_score"`
	GradePoints float64 `json:"grade_points" binding:"required"`
}

type BulkGradeEntryRequest struct {
	KrsItemID   string                   `json:"krs_item_id" binding:"required"`
	Entries    []GradeEntryRequestBulk  `json:"entries" binding:"gt=0"`
}

type GradeEntryRequestBulk struct {
	StudentID   string            `json:"student_id" binding:"required"`
	Scores     map[string]float64 `json:"scores"`
	FinalGrade *string           `json:"final_grade"`
	LetterGrade *string          `json:"letter_grade"`
}

// CreateGrade handles POST /api/v1/academic/grades
func (h *GradeHandler) CreateGrade(c *gin.Context) {
	var req GradeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	grade := domain.Grade{
		KrsItemID: req.KrsItemID,
		Source:   req.Source,
		Status:   "in_progress",
	}

	var comps []domain.GradeComponent
	for _, comp := range req.Components {
		c := domain.GradeComponent{
			GradeID:  grade.ID,
			Name:     comp.Name,
			Weight:  comp.Weight,
			MaxScore: comp.MaxScore,
		}
		comps = append(comps, c)
	}
	grade.Components = comps

	if err := h.repo.CreateGrade(&grade); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan grade").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(grade).WithContext(c))
}

// GetGrade handles GET /api/v1/academic/grades/:id
func (h *GradeHandler) GetGrade(c *gin.Context) {
	id := c.Param("id")

	grade, err := h.repo.GetGradeByID(id)
	if err != nil || grade == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Grade tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(grade).WithContext(c))
}

// GetStudentGrades handles GET /api/v1/academic/grades/student/:student_id
func (h *GradeHandler) GetStudentGrades(c *gin.Context) {
	studentID := c.Param("student_id")
	academicPeriodID := c.Query("academic_period_id")

	var grades []domain.Grade
	query := h.db.Table("grades").
		Select("grades.id, grades.krs_item_id, grades.source, grades.status, grades.final_grade, grades.letter_grade").
		Joins("JOIN krs_items ON krs_items.id = grades.krs_item_id").
		Joins("JOIN krs ON krs.id = krs_items.krs_id").
		Where("krs.student_id = ?", studentID)

	if academicPeriodID != "" {
		query = query.Where("krs.academic_period_id = ?", academicPeriodID)
	}

	if err := query.Find(&grades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil nilai").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(grades).WithContext(c))
}

// SubmitGrade handles POST /api/v1/academic/grades/:id/submit
func (h *GradeHandler) SubmitGrade(c *gin.Context) {
	id := c.Param("id")

	grade, err := h.repo.GetGradeByID(id)
	if err != nil || grade == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Grade tidak ditemukan").WithContext(c))
		return
	}

	if grade.Status != "in_progress" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Hanya grade berstatus in_progress yang bisa disubmit").WithContext(c))
		return
	}

	grade.Status = "submitted"
	if err := h.repo.UpdateGradeStatus(id, "submitted"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate status grade").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(grade).WithContext(c))
}

// FinalizeGrade handles POST /api/v1/academic/grades/:id/finalize
func (h *GradeHandler) FinalizeGrade(c *gin.Context) {
	id := c.Param("id")

	grade, err := h.repo.GetGradeByID(id)
	if err != nil || grade == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Grade tidak ditemukan").WithContext(c))
		return
	}

	if grade.Status != "submitted" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Hanya grade berstatus submitted yang bisa difinalisasi").WithContext(c))
		return
	}

	grade.Status = "final"
	if err := h.repo.UpdateGradeStatus(id, "final"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memfinalisasi grade").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.grade.finalize",
		Module:      "academic",
		ResourceType: "grade",
		ResourceID:  id,
	})

	c.JSON(http.StatusOK, sharederr.Success(grade).WithContext(c))
}

// EnterStudentGrade handles POST /api/v1/academic/grades/:id/entries
func (h *GradeHandler) EnterStudentGrade(c *gin.Context) {
	gradeID := c.Param("id")

	var req GradeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	grade, err := h.repo.GetGradeByID(gradeID)
	if err != nil || grade == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Grade tidak ditemukan").WithContext(c))
		return
	}

	entry := domain.GradeEntry{
		GradeID:     gradeID,
		StudentID:  req.StudentID,
		KrsItemID:  req.KrsItemID,
		FinalGrade: req.FinalGrade,
		LetterGrade: req.LetterGrade,
		Status:     req.Status,
		EnteredAt:  time.Now(),
	}

	if err := h.db.Create(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan entry").WithContext(c))
		return
	}

	// Create component scores
	for _, comp := range req.Components {
		score := domain.GradeComponentScore{
			GradeEntryID: entry.ID,
			Name:         comp.Name,
			Score:        comp.Score,
		}
		h.db.Create(&score)
	}

	c.JSON(http.StatusCreated, sharederr.Success(entry).WithContext(c))
}

// BulkEnterGrades handles POST /api/v1/academic/grades/:id/entries/bulk
func (h *GradeHandler) BulkEnterGrades(c *gin.Context) {
	gradeID := c.Param("id")

	var req BulkGradeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	grade, err := h.repo.GetGradeByID(gradeID)
	if err != nil || grade == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Grade tidak ditemukan").WithContext(c))
		return
	}

	var created []domain.GradeEntry
	var failed []string

	for _, entryReq := range req.Entries {
		entry := domain.GradeEntry{
			GradeID:     gradeID,
			StudentID:   entryReq.StudentID,
			KrsItemID:  req.KrsItemID,
			FinalGrade: entryReq.FinalGrade,
			LetterGrade: entryReq.LetterGrade,
			Status:     "in_progress",
			EnteredAt:  time.Now(),
		}

		if err := h.db.Create(&entry).Error; err != nil {
			failed = append(failed, entryReq.StudentID+" (save error)")
			continue
		}

		// Create component scores from scores map
		for name, score := range entryReq.Scores {
			compScore := domain.GradeComponentScore{
				GradeEntryID: entry.ID,
				Name:         name,
				Score:        score,
			}
			h.db.Create(&compScore)
		}

		created = append(created, entry)
	}

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"created": created,
		"failed": failed,
	}).WithContext(c))
}

// UpdateGradeConversion handles POST /api/v1/academic/grades/conversion
func (h *GradeHandler) UpdateGradeConversion(c *gin.Context) {
	var req GradeConversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get grade by letter
	existing, _ := h.repo.GetGradeConversionByLetter(req.GradeLetter)
	
	if existing != nil {
		existing.MinScore = req.MinScore
		existing.MaxScore = req.MaxScore
		existing.GradePoints = req.GradePoints
		h.repo.UpdateGradeConversion(existing)
	} else {
		conv := domain.GradeConversion{
			GradeLetter: req.GradeLetter,
			MinScore:   req.MinScore,
			MaxScore:   req.MaxScore,
			GradePoints: req.GradePoints,
		}
		h.repo.CreateGradeConversion(&conv)
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Grade conversion updated").WithContext(c))
}

// GetGradeConversions handles GET /api/v1/academic/grades/conversion
func (h *GradeHandler) GetGradeConversions(c *gin.Context) {
	conversions, err := h.repo.GetAllGradeConversions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil konversi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(conversions).WithContext(c))
}

// GetTranscript handles GET /api/v1/academic/grades/transcript/:student_id
func (h *GradeHandler) GetTranscript(c *gin.Context) {
	studentID := c.Param("student_id")
	academicPeriodID := c.Query("academic_period_id")

	transcript, err := h.repo.GenerateTranscript(studentID, academicPeriodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghasilkan transkrip").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(transcript).WithContext(c))
}

// GetIPK handles GET /api/v1/academic/grades/ipk/:student_id
func (h *GradeHandler) GetIPK(c *gin.Context) {
	studentID := c.Param("student_id")

	ipk, err := h.repo.CalculateIPK(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghitung IPK").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(ipk).WithContext(c))
}

// GetIPS handles GET /api/v1/academic/grades/ips/:student_id
func (h *GradeHandler) GetIPS(c *gin.Context) {
	studentID := c.Param("student_id")
	academicPeriodID := c.Query("academic_period_id")

	if academicPeriodID == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "academic_period_id is required").WithContext(c))
		return
	}

	ips, err := h.repo.CalculateIPS(studentID, academicPeriodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghitung IPS").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(ips).WithContext(c))
}

// CalculateLetterGrade converts numeric score to letter grade
func (h *GradeHandler) CalculateLetterGrade(score float64) string {
	conversions, err := h.repo.GetAllGradeConversions()
	if err != nil || len(conversions) == 0 {
		if score >= 85 {
			return "A"
		} else if score >= 75 {
			return "B"
		} else if score >= 65 {
			return "C"
		} else if score >= 55 {
			return "D"
		}
		return "E"
	}

	for _, conv := range conversions {
		if score >= conv.MinScore {
			return conv.GradeLetter
		}
	}

	return "E"
}
