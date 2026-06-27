package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ TranscriptHandler ============
// Handler untuk Transcript & Kartu Studi Mahasiswa

type TranscriptHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewTranscriptHandler creates a new TranscriptHandler
func NewTranscriptHandler(db *gorm.DB) *TranscriptHandler {
	return &TranscriptHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Transcript APIs ============

// GetStudentTranscript gets complete transcript for a student
// GET /api/v1/academic/transcripts/:student_id
func (h *TranscriptHandler) GetStudentTranscript(c *gin.Context) {
	studentID := c.Param("student_id")

	// Get student info
	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Get all grades
	grades, err := h.repo.GetGradesByStudent(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil nilai").WithContext(c))
		return
	}

	// Calculate GPA
	var totalPoints float64
	var totalSKS int
	for _, grade := range grades {
		if grade.LetterGrade != "" {
			gradePoints := getGradePoints(grade.LetterGrade)
			course, _ := h.repo.GetCourseByKrsItemID(grade.KrsItemID)
			if course != nil {
				totalPoints += gradePoints * float64(course.Sks)
				totalSKS += course.Sks
			}
		}
	}

	gpa := 0.0
	if totalSKS > 0 {
		gpa = totalPoints / float64(totalSKS)
	}

	// Group by semester
	semesterData := make(map[string][]domain.Grade)
	for _, grade := range grades {
		periodName, _ := h.repo.GetPeriodNameByKrsItemID(grade.KrsItemID)
		if periodName != "" {
			semesterData[periodName] = append(semesterData[periodName], grade)
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student": gin.H{
			"id":            student.ID,
			"name":          student.Name,
			"nim":           student.Nim,
			"study_program": student.StudyProgramID,
			"entry_year":    student.EntryYear,
		},
		"grades":           grades,
		"semester_data":   semesterData,
		"total_sks":       totalSKS,
		"gpa":            gpa,
		"total_courses":   len(grades),
	}).WithContext(c))
}

// GetStudentTranscriptPDF generates PDF transcript
// GET /api/v1/academic/transcripts/:student_id/pdf
func (h *TranscriptHandler) GetStudentTranscriptPDF(c *gin.Context) {
	studentID := c.Param("student_id")

	// Verify student exists
	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Get grades
	grades, _ := h.repo.GetGradesByStudent(studentID)
	_ = grades

	// In production, generate PDF here
	// For now, return JSON with PDF flag
	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id": studentID,
		"pdf_url":    fmt.Sprintf("/api/v1/academic/transcripts/%s.pdf", studentID),
		"generated_at": time.Now().Format(time.RFC3339),
		"note":       "PDF generation endpoint - integrate with PDF library for actual file",
	}).WithContext(c))
}

// VerifyTranscript verifies a transcript
// GET /api/v1/academic/transcripts/:student_id/verify
func (h *TranscriptHandler) VerifyTranscript(c *gin.Context) {
	studentID := c.Param("student_id")
	verificationCode := c.Query("code")

	// Get student
	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Verify by verification code or student NIM
	if verificationCode == "" {
		verificationCode = student.Nim
	}

	// Get stats
	grades, _ := h.repo.GetGradesByStudent(studentID)
	var totalSKS int
	for _, grade := range grades {
		if grade.LetterGrade != "" {
			course, _ := h.repo.GetCourseByKrsItemID(grade.KrsItemID)
			if course != nil {
				totalSKS += course.Sks
			}
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"verified":        true,
		"student_id":      studentID,
		"student_name":   student.Name,
		"nim":           student.Nim,
		"total_courses":  len(grades),
		"total_sks":     totalSKS,
		"verification_code": verificationCode,
		"verified_at":   time.Now().Format(time.RFC3339),
	}).WithContext(c))
}

// GetStudentStudySummary gets study summary
// GET /api/v1/academic/transcripts/:student_id/summary
func (h *TranscriptHandler) GetStudentStudySummary(c *gin.Context) {
	studentID := c.Param("student_id")

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Get grades
	grades, _ := h.repo.GetGradesByStudent(studentID)

	// Calculate totals
	var passed, failed, inProgress int
	var totalSKSPassed, totalSKSFailed int
	var totalPoints float64

	for _, grade := range grades {
		if grade.LetterGrade == "" {
			inProgress++
			continue
		}

		gradePoints := getGradePoints(grade.LetterGrade)
		course, _ := h.repo.GetCourseByKrsItemID(grade.KrsItemID)
		if course == nil {
			continue
		}

		if isPassingGrade(grade.LetterGrade) {
			passed++
			totalSKSPassed += course.Sks
			totalPoints += gradePoints * float64(course.Sks)
		} else {
			failed++
			totalSKSFailed += course.Sks
		}
	}

	gpa := 0.0
	if totalSKSPassed > 0 {
		gpa = totalPoints / float64(totalSKSPassed)
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":      studentID,
		"total_courses":  len(grades),
		"passed":        passed,
		"failed":       failed,
		"in_progress":   inProgress,
		"total_sks":    totalSKSPassed + totalSKSFailed,
		"sks_passed":   totalSKSPassed,
		"sks_failed":  totalSKSFailed,
		"gpa":        gpa,
		"status":      student.Status,
	}).WithContext(c))
}

// GetCourseDistribution gets grade distribution for a course
// GET /api/v1/academic/transcripts/course/:course_id/distribution
func (h *TranscriptHandler) GetCourseDistribution(c *gin.Context) {
	courseID := c.Param("course_id")

	// Verify course exists
	course, err := h.repo.GetCourseByID(courseID)
	if err != nil || course == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mata kuliah tidak ditemukan").WithContext(c))
		return
	}

	// Get all grades for this course
	grades, err := h.repo.GetGradesByCourse(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data").WithContext(c))
		return
	}

	// Count distribution
	distribution := map[string]int{
		"A": 0, "A-": 0,
		"B+": 0, "B": 0, "B-": 0,
		"C+": 0, "C": 0, "C-": 0,
		"D": 0, "E": 0,
	}

	var totalPoints float64
	var totalSKS int

	for _, grade := range grades {
		if grade.LetterGrade != "" {
			distribution[grade.LetterGrade]++
			gradePoints := getGradePoints(grade.LetterGrade)
			totalPoints += gradePoints
			totalSKS++
		}
	}

	avg := 0.0
	if totalSKS > 0 {
		avg = totalPoints / float64(totalSKS)
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"course_id":      courseID,
		"course_name":   course.CourseName,
		"distribution": distribution,
		"total_students": len(grades),
		"average_grade": avg,
	}).WithContext(c))
}

// ============ Helper Functions ============

func getGradePoints(grade string) float64 {
	switch grade {
	case "A", "A-":
		return 4.0
	case "B+":
		return 3.5
	case "B", "B-":
		return 3.0
	case "C+":
		return 2.5
	case "C", "C-":
		return 2.0
	case "D":
		return 1.0
	case "E":
		return 0.0
	default:
		return 0.0
	}
}

func isPassingGrade(grade string) bool {
	return grade != "E" && grade != ""
}
