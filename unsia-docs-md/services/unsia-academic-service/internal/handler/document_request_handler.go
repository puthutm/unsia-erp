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

// ============ Document Request (Persuratan) Handler ============
// Handles student document request like: Transkrip, KTM, Ijazah Copy, etc.

// DocumentRequestHandler handles document request operations
type DocumentRequestHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewDocumentRequestHandler creates a new DocumentRequestHandler
func NewDocumentRequestHandler(db *gorm.DB) *DocumentRequestHandler {
	return &DocumentRequestHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Document Types ============

const (
	DocTypeTranskrip        = "transkrip"
	DocTypeIjazahCopy      = "ijazah_copy"
	DocTypeKTM             = "ktm"
	DocTypeKartuKelulusan   = "kartu_kelulusan"
	DocTypePrestasi        = "sertifikat_prestasi"
	DocTypePenghasilan     = "surat_penghasilan_ortu"
	DocTypeAktifKuliah     = "surat_aktif_kuliah"
	DocTypeLainnya         = "lainnya"
)

const (
	DocStatusPending   = "pending"
	DocStatusApproved = "approved"
	DocStatusRejected = "rejected"
	DocStatusReady    = "ready"
	DocStatusTaken   = "taken"
)

// ============ Request Types ============

type DocumentRequestCreate struct {
	StudentID    string  `json:"student_id" binding:"required"`
	DocumentType string  `json:"document_type" binding:"required"`
	Purpose     string  `json:"purpose"`
	Notes       string  `json:"notes"`
	Quantity    int     `json:"quantity"`
}

type DocumentRequestUpdate struct {
	Status     string `json:"status" binding:"required,oneof=pending approved rejected ready taken"`
	AdminNote  string `json:"admin_note"`
	ProcessedBy string `json:"processed_by"`
}

// ============ Document Request Endpoints ============

// CreateRequest creates a new document request
// POST /api/v1/academic/documents/requests
func (h *DocumentRequestHandler) CreateRequest(c *gin.Context) {
	var req DocumentRequestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify student exists
	student, err := h.repo.GetStudentByID(req.StudentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Mahasiswa tidak ditemukan").WithContext(c))
		return
	}

	// Check if student can request (not in debt)
	// In production, check with finance clearance first
	if student.StudentStatus != "active" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_STATUS", "Mahasiswa tidak aktif").WithContext(c))
		return
	}

	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	docRequest := domain.DocumentRequest{
		ID:            generateUUID(),
		StudentID:     req.StudentID,
		DocumentType:  req.DocumentType,
		Purpose:       req.Purpose,
		Notes:        req.Notes,
		Quantity:      quantity,
		Status:       DocStatusPending,
		RequestDate:   time.Now(),
		ProcessedDate: nil,
	}

	if err := h.db.Create(&docRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat permintaan").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.document_request.create",
		Module:       "academic",
		ResourceType: "document_request",
		ResourceID:   docRequest.ID,
		NewValue:     docRequest,
	})

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"request_id":     docRequest.ID,
		"document_type": docRequest.DocumentType,
		"request_date":  docRequest.RequestDate,
		"status":       docRequest.Status,
	}).WithContext(c))
}

// ListRequests lists document requests with filters
// GET /api/v1/academic/documents/requests
func (h *DocumentRequestHandler) ListRequests(c *gin.Context) {
	studentID := c.Query("student_id")
	documentType := c.Query("document_type")
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	query := h.db.Model(&domain.DocumentRequest{})

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if documentType != "" {
		query = query.Where("document_type = ?", documentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var requests []domain.DocumentRequest
	offset := (page - 1) * limit
	query.Offset(offset).Limit(limit).Order("request_date DESC").Find(&requests)

	c.JSON(http.StatusOK, sharederr.Success(requests).WithPagination(page, limit, int(total)).WithContext(c))
}

// GetRequest gets a specific document request
// GET /api/v1/academic/documents/requests/:id
func (h *DocumentRequestHandler) GetRequest(c *gin.Context) {
	requestID := c.Param("id")

	var docRequest domain.DocumentRequest
	if err := h.db.First(&docRequest, "id = ?", requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Permintaan tidak ditemukan").WithContext(c))
		return
	}

	// Load student info
	student, _ := h.repo.GetStudentByID(docRequest.StudentID)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"request":       docRequest,
		"student_nim":  student.Nim,
		"student_name": "Student Name", // Would fetch from person table
	}).WithContext(c))
}

// UpdateRequest updates document request status (Admin)
// PUT /api/v1/academic/documents/requests/:id
func (h *DocumentRequestHandler) UpdateRequest(c *gin.Context) {
	requestID := c.Param("id")
	var req DocumentRequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var docRequest domain.DocumentRequest
	if err := h.db.First(&docRequest, "id = ?", requestID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Permintaan tidak ditemukan").WithContext(c))
		return
	}

	oldStatus := docRequest.Status

	updates := map[string]interface{}{
		"status": req.Status,
	}

	if req.Status == DocStatusApproved || req.Status == DocStatusRejected || req.Status == DocStatusReady {
		now := time.Now()
		updates["processed_date"] = now
	}

	if req.AdminNote != "" {
		updates["admin_note"] = req.AdminNote
	}

	if err := h.db.Model(&docRequest).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui permintaan").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.document_request.update",
		Module:       "academic",
		ResourceType: "document_request",
		ResourceID:  requestID,
		OldValue:    oldStatus,
		NewValue:    req.Status,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(gin.H{
		"request_id":     requestID,
		"old_status":   oldStatus,
		"new_status":  req.Status,
		"processed_date": time.Now(),
	}, "Permintaan berhasil diperbarui").WithContext(c))
}

// GetMyRequests gets current user's document requests
// GET /api/v1/academic/documents/requests/me
func (h *DocumentRequestHandler) GetMyRequests(c *gin.Context) {
	studentID, _ := c.Get("x-user-id")
	if studentID == nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Unauthorized").WithContext(c))
		return
	}

	studentIDStr := studentID.(string)

	var requests []domain.DocumentRequest
	h.db.Where("student_id = ?", studentIDStr).Order("request_date DESC").Find(&requests)

	c.JSON(http.StatusOK, sharederr.Success(requests).WithContext(c))
}

// GetDocumentTypes gets available document types
// GET /api/v1/academic/documents/types
func (h *DocumentRequestHandler) GetDocumentTypes(c *gin.Context) {
docTypes := []gin.H{
		{"id": DocTypeTranskrip, "name": "Transkrip Nilai", "description": "Transkrip akademik", "processing_days": 3},
		{"id": DocTypeIjazahCopy, "name": "Salinan Ijazah", "description": "Salinan ijazah legalisir", "processing_days": 7},
		{"id": DocTypeKTM, "name": "KTM Baru", "description": "Kartu Tanda Mahasiswa", "processing_days": 14},
		{"id": DocTypeKartuKelulusan, "name": "Kartu Kelulusan", "description": "Kartu kelulusan wisuda", "processing_days": 5},
		{"id": DocTypePrestasi, "name": "Sertifikat Prestasi", "description": "Sertifikat prestasi akademik", "processing_days": 3},
		{"id": DocTypePenghasilan, "name": "Surat Penghasilan Orang Tua", "description": "Untuk pengajuan scholarship", "processing_days": 2},
		{"id": DocTypeAktifKuliah, "name": "Surat Aktif Kuliah", "description": "Bukti mahasiswa aktif", "processing_days": 1},
		{"id": DocTypeLainnya, "name": "Surat Lainnya", "description": "Surat keterangan lain", "processing_days": 3},
	}

	c.JSON(http.StatusOK, sharederr.Success(docTypes).WithContext(c))
}

// GetStatistics gets document request statistics (Admin)
// GET /api/v1/academic/documents/statistics
func (h *DocumentRequestHandler) GetStatistics(c *gin.Context) {
	year := c.Query("year")
	month := c.Query("month")

	query := h.db.Model(&domain.DocumentRequest{})
	if year != "" {
		query = query.Where("EXTRACT(YEAR FROM request_date) = ?", year)
	}
	if month != "" {
		query = query.Where("EXTRACT(MONTH FROM request_date) = ?", month)
	}

	var stats struct {
		Total       int64 `json:"total"`
		Pending     int64 `json:"pending"`
		Approved   int64 `json:"approved"`
		Rejected   int64 `json:"rejected"`
		Ready      int64 `json:"ready"`
		Taken      int64 `json:"taken"`
	}

	query.Count(&stats.Total)
	h.db.Model(&domain.DocumentRequest{}).Where("status = ?", DocStatusPending).Count(&stats.Pending)
	h.db.Model(&domain.DocumentRequest{}).Where("status = ?", DocStatusApproved).Count(&stats.Approved)
	h.db.Model(&domain.DocumentRequest{}).Where("status = ?", DocStatusRejected).Count(&stats.Rejected)
	h.db.Model(&domain.DocumentRequest{}).Where("status = ?", DocStatusReady).Count(&stats.Ready)
	h.db.Model(&domain.DocumentRequest{}).Where("status = ?", DocStatusTaken).Count(&stats.Taken)

	c.JSON(http.StatusOK, sharederr.Success(stats).WithContext(c))
}

// ============ Helper Functions ============

func generateUUID() string {
	return "doc-" + time.Now().Format("20060102150405")
}
