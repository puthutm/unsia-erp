package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Attendance Handler ============

// AttendanceHandler handles student attendance operations
type AttendanceHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewAttendanceHandler creates a new AttendanceHandler
func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
	return &AttendanceHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Attendance Request Types ============

type AttendanceRecordRequest struct {
	StudentID    string  `json:"student_id" binding:"required"`
	ClassID      string  `json:"class_id" binding:"required"`
	SessionDate  string  `json:"session_date" binding:"required"`
	Status      string  `json:"status" binding:"required,oneof=present absent excused sick"`
	Note        *string `json:"note"`
}

type BulkAttendanceRequest struct {
	ClassID     string                       `json:"class_id" binding:"required"`
	SessionDate string                       `json:"session_date" binding:"required"`
	Attendances []AttendanceRecordRequest     `json:"attendances" binding:"required,gt=0"`
}

type AttendanceStatsRequest struct {
	StudentID        string `json:"student_id"`
	ClassID          string `json:"class_id"`
	AcademicPeriodID string `json:"academic_period_id"`
}

// ============ Record Attendance ============

// RecordAttendance handles POST /api/v1/academic/attendances
func (h *AttendanceHandler) RecordAttendance(c *gin.Context) {
	var req AttendanceRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Parse date
	sessionDate, err := time.Parse("2006-01-02", req.SessionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError("Invalid date format. Use YYYY-MM-DD").WithContext(c))
		return
	}

	// Check if student is enrolled
	enrolled, _ := h.repo.IsStudentEnrolledInClass(req.StudentID, req.ClassID)
	if !enrolled {
		c.JSON(http.StatusBadRequest, sharederr.Error("NOT_ENROLLED", "Mahasiswa tidak terdaftar di kelas ini").WithContext(c))
		return
	}

	// Check if attendance already exists
	existing, _ := h.repo.GetAttendance(req.StudentID, req.ClassID, sessionDate)
	if existing != nil {
		// Update existing
		existing.Status = req.Status
		if req.Note != nil {
			existing.Note = req.Note
		}
		h.repo.UpdateAttendance(existing)
		
		sharedaudit.Log(c, sharedaudit.AuditEntry{
			Action:       "academic.attendance.update",
			Module:       "academic",
			ResourceType: "attendance",
			ResourceID:   existing.ID,
		})
		
		c.JSON(http.StatusOK, sharederr.Success(existing).WithContext(c))
		return
	}

	// Create new attendance
	attendance := domain.StudentAttendance{
		StudentID:   req.StudentID,
		ClassID:     req.ClassID,
		SessionDate: sessionDate,
		Status:      req.Status,
		Note:        req.Note,
	}

	if err := h.repo.CreateAttendance(&attendance); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan absensi").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.attendance.create",
		Module:       "academic",
		ResourceType: "attendance",
		ResourceID:   attendance.ID,
		NewValue:     attendance,
	})

	c.JSON(http.StatusCreated, sharederr.Success(attendance).WithContext(c))
}

// ============ Bulk Record Attendance ============

// RecordBulkAttendance handles POST /api/v1/academic/attendances/bulk
func (h *AttendanceHandler) RecordBulkAttendance(c *gin.Context) {
	var req BulkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Parse date
	sessionDate, err := time.Parse("2006-01-02", req.SessionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError("Invalid date format").WithContext(c))
		return
	}

	var created []domain.StudentAttendance
	var failed []string

	for i, att := range req.Attendances {
		// Verify enrollment
		enrolled, _ := h.repo.IsStudentEnrolledInClass(att.StudentID, req.ClassID)
		if !enrolled {
			failed = append(failed, "Row "+strconv.Itoa(i+1)+": not enrolled")
			continue
		}

		// Check if exists
		existing, _ := h.repo.GetAttendance(att.StudentID, req.ClassID, sessionDate)
		if existing != nil {
			existing.Status = att.Status
			if att.Note != nil {
				existing.Note = att.Note
			}
			h.repo.UpdateAttendance(existing)
			created = append(created, *existing)
			continue
		}

		attendance := domain.StudentAttendance{
			StudentID:   att.StudentID,
			ClassID:     req.ClassID,
			SessionDate: sessionDate,
			Status:      att.Status,
			Note:        att.Note,
		}

		if err := h.repo.CreateAttendance(&attendance); err != nil {
			failed = append(failed, "Row "+strconv.Itoa(i+1)+": save error")
			continue
		}

		created = append(created, attendance)
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.attendance.bulk_create",
		Module:       "academic",
		ResourceType: "attendance",
		ResourceID:  "bulk_" + req.ClassID,
	})

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"created": created,
		"failed":  failed,
		"total":  len(req.Attendances),
	}).WithContext(c))
}

// ============ Get Attendance ============

// GetStudentAttendance handles GET /api/v1/academic/students/:student_id/attendances
func (h *AttendanceHandler) GetStudentAttendance(c *gin.Context) {
	studentID := c.Param("student_id")
	classID := c.Query("class_id")
	academicPeriodID := c.Query("academic_period_id")

	attendances, err := h.repo.GetStudentAttendance(studentID, classID, academicPeriodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data absensi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(attendances).WithContext(c))
}

// GetClassAttendance handles GET /api/v1/academic/classes/:class_id/attendances
func (h *AttendanceHandler) GetClassAttendance(c *gin.Context) {
	classID := c.Param("class_id")
	sessionDate := c.Query("session_date")

	attendances, err := h.repo.GetClassAttendance(classID, sessionDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data absensi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(attendances).WithContext(c))
}

// ============ Attendance Statistics ============

// GetAttendanceStats handles GET /api/v1/academic/attendances/stats
func (h *AttendanceHandler) GetAttendanceStats(c *gin.Context) {
	studentID := c.Query("student_id")
	classID := c.Query("class_id")

	if studentID == "" && classID == "" {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError("student_id or class_id required").WithContext(c))
		return
	}

	stats, err := h.repo.GetAttendanceStats(studentID, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil statistik absensi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(stats).WithContext(c))
}

// ============ My Attendance (for Student) ============

// GetMyAttendance handles GET /api/v1/academic/attendances/me
func (h *AttendanceHandler) GetMyAttendance(c *gin.Context) {
	studentID, _ := c.Get("x-user-id")
	if studentID == nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Unauthorized").WithContext(c))
		return
	}

	studentIDStr := studentID.(string)
	classID := c.Query("class_id")
	academicPeriodID := c.Query("academic_period_id")

	attendances, err := h.repo.GetStudentAttendance(studentIDStr, classID, academicPeriodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data absensi").WithContext(c))
		return
	}

	// Calculate summary
	stats, _ := h.repo.GetAttendanceStats(studentIDStr, classID)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"details": attendances,
		"stats":   stats,
	}).WithContext(c))
}
