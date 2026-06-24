package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Clearance Handler ============
// Handles student payment clearance checking (SPP/Tuition integration with Finance)

// ClearanceHandler handles student clearance/payment verification
type ClearanceHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewClearanceHandler creates a new ClearanceHandler
func NewClearanceHandler(db *gorm.DB) *ClearanceHandler {
	return &ClearanceHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Clearance Response Types ============

type ClearanceResponse struct {
	StudentID       string  `json:"student_id"`
	NIM             string  `json:"nim"`
	FullName        string  `json:"full_name"`
	StudyProgramID  string  `json:"study_program_id"`
	StudyProgram   string  `json:"study_program"`
	CurrentSemester int    `json:"current_semester"`
	PaymentStatus  string  `json:"payment_status"`
	AmountDue      float64 `json:"amount_due"`
	AmountPaid     float64 `json:"amount_paid"`
	Balance        float64 `json:"balance"`
	LastPaymentDate *string `json:"last_payment_date"`
	AcademicYear  string  `json:"academic_year"`
	ClearanceLevel string  `json:"clearance_level"` // "cleared", "partial", "overdue"
	Message       string  `json:"message"`
}

type PaymentDetail struct {
	StudentID      string  `json:"student_id"`
	InvoiceID     string  `json:"invoice_id"`
	AcademicYear  string  `json:"academic_year"`
	Amount        float64 `json:"amount"`
	PaidAmount    float64 `json:"paid_amount"`
	DueDate       string  `json:"due_date"`
	PaymentStatus string  `json:"status"`
	Installments  []InstallmentDetail `json:"installments,omitempty"`
}

type InstallmentDetail struct {
	InstallmentNo int     `json:"installment_no"`
	Amount       float64 `json:"amount"`
	DueDate      string  `json:"due_date"`
	PaidDate     *string `json:"paid_date,omitempty"`
	Status      string  `json:"status"` // "paid", "unpaid", "overdue"
}

// ============ Get Student Clearance ============

// GetStudentClearance handles GET /api/v1/academic/clearance/:student_id
// Returns the payment clearance status for a student
func (h *ClearanceHandler) GetStudentClearance(c *gin.Context) {
	studentID := c.Param("student_id")

	// Get student details
	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Get payment status from finance service (simulated integration)
	// In production, this would call finance-service via HTTP
	clearance := h.getPaymentStatusFromFinance(student)

	c.JSON(http.StatusOK, sharederr.Success(clearance).WithContext(c))
}

// GetMyClearance handles GET /api/v1/academic/clearance/me
// Returns the payment clearance for the logged-in student
func (h *ClearanceHandler) GetMyClearance(c *gin.Context) {
	studentID, _ := c.Get("x-user-id")
	if studentID == nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Unauthorized").WithContext(c))
		return
	}

	studentIDStr := studentID.(string)

	student, err := h.repo.GetStudentByID(studentIDStr)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	clearance := h.getPaymentStatusFromFinance(student)

	c.JSON(http.StatusOK, sharederr.Success(clearance).WithContext(c))
}

// ============ Get Payment History ============

// GetPaymentHistory handles GET /api/v1/academic/clearance/:student_id/payments
// Returns the payment history for a student
func (h *ClearanceHandler) GetPaymentHistory(c *gin.Context) {
	studentID := c.Param("student_id")
	academicYear := c.Query("academic_year")

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Get payment details from finance service (simulated)
	payments := h.getPaymentHistoryFromFinance(student, academicYear)

	c.JSON(http.StatusOK, sharederr.Success(payments).WithContext(c))
}

// ============ Check Clearance for KRS ============

// CheckKrsClearance handles GET /api/v1/academic/clearance/:student_id/krs-eligibility
// Returns whether a student is eligible to create KRS based on payment
func (h *ClearanceHandler) CheckKrsClearance(c *gin.Context) {
	studentID := c.Param("student_id")

	student, err := h.repo.GetStudentByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	clearance := h.getPaymentStatusFromFinance(student)

	eligible := clearance.ClearanceLevel == "cleared" || clearance.ClearanceLevel == "partial"

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":      studentID,
		"eligible":       eligible,
		"clearance_level": clearance.ClearanceLevel,
		"message":       clearance.Message,
		"balance":       clearance.Balance,
	}).WithContext(c))
}

// ============ Verify Clearance (for Admin) ============

// VerifyClearance handles POST /api/v1/academic/clearance/verify
// Admin manually verifies student clearance
func (h *ClearanceHandler) VerifyClearance(c *gin.Context) {
	var req struct {
		StudentID    string  `json:"student_id" binding:"required"`
		ClearanceLevel string `json:"clearance_level" binding:"required,oneof=cleared partial overdue"`
		Note        string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	student, err := h.repo.GetStudentByID(req.StudentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// In production, this would update the clearance status in finance-service
	// For now, we just return a success message
	
	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"student_id":      req.StudentID,
		"clearance_level": req.ClearanceLevel,
		"verified_at":    "2025-01-15T10:00:00Z",
		"note":          req.Note,
	}).WithContext(c))
}

// ============ Helper Functions ============

// getPaymentStatusFromFinance simulates calling finance service
// In production, this would make HTTP call to finance-service
func (h *ClearanceHandler) getPaymentStatusFromFinance(student *domain.Student) *ClearanceResponse {
	// Simulate payment status calculation
	// In production, this data comes from finance-service
	
	amountDue := 5000000.0 // SPP per semester
	amountPaid := amountDue // Assume paid for simulation
	balance := amountDue - amountPaid

	var clearanceLevel string
	var message string

	if balance <= 0 {
		clearanceLevel = "cleared"
		message = "Mahasiswa sudah membayar SPP secara penuh"
	} else if balance < amountDue {
		clearanceLevel = "partial"
		message = "Mahasiswa masih memiliki tunggakan"
	} else {
		clearanceLevel = "overdue"
		message = "Mahasiswa belum membayar SPP"
	}

	return &ClearanceResponse{
		StudentID:       student.ID,
		NIM:            student.Nim,
		FullName:       "Student Name", // Would fetch from person table
		StudyProgramID: student.StudyProgramID,
		CurrentSemester: student.CurrentSemester,
		PaymentStatus:  clearanceLevel,
		AmountDue:     amountDue,
		AmountPaid:    amountPaid,
		Balance:       balance,
		AcademicYear:  "2024/2025",
		ClearanceLevel: clearanceLevel,
		Message:       message,
	}
}

// getPaymentHistoryFromFinance simulates payment history from finance service
func (h *ClearanceHandler) getPaymentHistoryFromFinance(student *domain.Student, academicYear string) []PaymentDetail {
	// Simulate payment history
	// In production, this data comes from finance-service invoice records
	
	payments := []PaymentDetail{
		{
			StudentID:      student.ID,
			InvoiceID:    "INV-2024-001",
			AcademicYear: "2024/2025",
			Amount:       5000000,
			PaidAmount:   5000000,
			DueDate:      "2024-09-15",
			PaymentStatus: "paid",
			Installments: []InstallmentDetail{
				{InstallmentNo: 1, Amount: 2500000, DueDate: "2024-08-15", PaidDate: ptrStr("2024-08-10"), Status: "paid"},
				{InstallmentNo: 2, Amount: 2500000, DueDate: "2024-09-15", PaidDate: ptrStr("2024-09-12"), Status: "paid"},
			},
		},
	}

	return payments
}

// Helper to create string pointer
func ptrStr(s string) *string {
	return &s
}
