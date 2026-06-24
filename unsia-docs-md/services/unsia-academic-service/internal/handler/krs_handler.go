package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ KRS Request Types ============

type KrsItemRequest struct {
	ClassID string `json:"class_id" binding:"required"`
}

type KrsCreateRequest struct {
	StudentID        string           `json:"student_id" binding:"required"`
	AcademicPeriodID string           `json:"academic_period_id" binding:"required"`
	Items            []KrsItemRequest `json:"items" binding:"required,gt=0"`
}

type KrsUpdateItemRequest struct {
	ClassID string `json:"class_id"`
	Status string `json:"status" binding:"oneof=selected approved dropped"`
}

// KrsHandler handles KRS-related operations
type KrsHandler struct {
	repo          *repository.AcademicRepository
	db            *gorm.DB
	financeClient interface {
		Get(ctx context.Context, path string) (*http.Response, error)
	}
}

// NewKrsHandler creates a new KrsHandler
func NewKrsHandler(db *gorm.DB) *KrsHandler {
	return &KrsHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// CreateKrsDraft creates a new KRS draft
// POST /api/v1/academic/krs/draft
func (h *KrsHandler) CreateKrsDraft(c *gin.Context) {
	var req KrsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	krs := domain.KRS{
		StudentID:        req.StudentID,
		AcademicPeriodID: req.AcademicPeriodID,
		Status:           "draft",
	}

	var items []domain.KrsItem
	for _, itemReq := range req.Items {
		item := domain.KrsItem{
			ClassID:    itemReq.ClassID,
			Status:     "selected",
			SelectedAt: time.Now(),
		}
		items = append(items, item)
	}
	krs.Items = items

	if err := h.repo.CreateKRS(&krs); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan draft KRS").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(krs).WithContext(c))
}

// SubmitKrs submits a KRS for approval
// POST /api/v1/academic/krs/:id/submit
func (h *KrsHandler) SubmitKrs(c *gin.Context) {
	krsID := c.Param("id")

	krs, err := h.repo.GetKRSByID(krsID)
	if err != nil || krs == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "KRS tidak ditemukan").WithContext(c))
		return
	}

	if krs.Status != "draft" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Hanya KRS berstatus draft yang bisa disubmit").WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	// NOTE: Finance clearance check would be called here if financeClient is configured
	// For now, we'll skip clearing and proceed directly

	if err := h.repo.UpdateKRSStatus(krsID, "submitted"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate status KRS").WithContext(c))
		return
	}

	krs.Status = "submitted"

	// Publish event
	envelope := sharedevent.EventEnvelope{
		EventName:        "academic.krs_submitted",
		EventVersion:     "v1",
		PublisherService: "academic-service",
		AggregateType:    "krs",
		AggregateID:      krs.ID,
		CorrelationID:    cid,
		Payload: map[string]interface{}{
			"krs_id":             krs.ID,
			"student_id":         krs.StudentID,
			"academic_period_id": krs.AcademicPeriodID,
			"status":            krs.Status,
		},
	}

	conn := h.db.Statement.ConnPool
	sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")

	c.JSON(http.StatusOK, sharederr.Success(krs).WithContext(c))
}

// ApproveKrs approves a submitted KRS
// POST /api/v1/academic/krs/:id/approve
func (h *KrsHandler) ApproveKrs(c *gin.Context) {
	krsID := c.Param("id")

	krs, err := h.repo.GetKRSByID(krsID)
	if err != nil || krs == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "KRS tidak ditemukan").WithContext(c))
		return
	}

	if krs.Status != "submitted" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Hanya KRS berstatus submitted yang bisa diapprove").WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Update KRS status
		if err := tx.Model(krs).Updates(map[string]interface{}{
			"status":      "approved",
			"approved_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		// Update items status
		if err := tx.Model(&domain.KrsItem{}).Where("krs_id = ?", krs.ID).Update("status", "approved").Error; err != nil {
			return err
		}

		// Save active grades placeholders for each class
		var classIDs []string
		for _, item := range krs.Items {
			classIDs = append(classIDs, item.ClassID)

			grade := domain.Grade{
				KrsItemID: item.ID,
				Source:    "lms",
			}
			tx.Create(&grade)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses persetujuan KRS").WithContext(c))
		return
	}

	// Publish event
	envelope := sharedevent.EventEnvelope{
		EventName:        "academic.krs_approved",
		EventVersion:     "v1",
		PublisherService: "academic-service",
		AggregateType:    "krs",
		AggregateID:      krs.ID,
		CorrelationID:    cid,
		Payload: map[string]interface{}{
			"krs_id":             krs.ID,
			"student_id":         krs.StudentID,
			"academic_period_id": krs.AcademicPeriodID,
		},
	}

	conn := h.db.Statement.ConnPool
	sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.krs.approve",
		Module:       "academic",
		ResourceType: "krs",
		ResourceID:   krsID,
	})

	krs.Status = "approved"
	c.JSON(http.StatusOK, sharederr.Success(krs).WithContext(c))
}

// RejectKrs rejects a submitted KRS
// POST /api/v1/academic/krs/:id/reject
func (h *KrsHandler) RejectKrs(c *gin.Context) {
	krsID := c.Param("id")

	krs, err := h.repo.GetKRSByID(krsID)
	if err != nil || krs == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "KRS tidak ditemukan").WithContext(c))
		return
	}

	if krs.Status != "submitted" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Hanya KRS berstatus submitted yang bisa ditolak").WithContext(c))
		return
	}

	if err := h.repo.UpdateKRSStatus(krsID, "rejected"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menolak KRS").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.krs.reject",
		Module:       "academic",
		ResourceType: "krs",
		ResourceID:   krsID,
	})

	krs.Status = "rejected"
	c.JSON(http.StatusOK, sharederr.Success(krs).WithContext(c))
}

// GetKrsDetail retrieves KRS details
// GET /api/v1/academic/krs/:id
func (h *KrsHandler) GetKrsDetail(c *gin.Context) {
	krsID := c.Param("id")

	krs, err := h.repo.GetKRSByID(krsID)
	if err != nil || krs == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "KRS tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(krs).WithContext(c))
}

// ListKrs lists KRS records
// GET /api/v1/academic/krs
func (h *KrsHandler) ListKrs(c *gin.Context) {
	studentID := c.Query("student_id")
	academicPeriodID := c.Query("academic_period_id")
	status := c.Query("status")

	var krsList []domain.KRS
	query := h.db.Model(&domain.KRS{})

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if academicPeriodID != "" {
		query = query.Where("academic_period_id = ?", academicPeriodID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&krsList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar KRS").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(krsList).WithContext(c))
}

// UpdateKrsItem updates a KRS item
// PUT /api/v1/academic/krs/:id/items/:item_id
func (h *KrsHandler) UpdateKrsItem(c *gin.Context) {
	krsID := c.Param("id")
	itemID := c.Param("item_id")

	var req KrsUpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify KRS exists
	krs, err := h.repo.GetKRSByID(krsID)
	if err != nil || krs == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "KRS tidak ditemukan").WithContext(c))
		return
	}

	// Verify item belongs to KRS
	item, err := h.repo.GetKrsItemByID(itemID)
	if err != nil || item == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Item KRS tidak ditemukan").WithContext(c))
		return
	}

	if item.KrsID != krsID {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Item tidak belonging to this KRS").WithContext(c))
		return
	}

	if req.Status != "" {
		if err := h.repo.UpdateKrsItemStatus(itemID, req.Status); err != nil {
			c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate item KRS").WithContext(c))
			return
		}
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Item KRS berhasil diperbarui").WithContext(c))
}

// AddKrsItem adds an item to existing KRS
// POST /api/v1/academic/krs/:id/items
func (h *KrsHandler) AddKrsItem(c *gin.Context) {
	krsID := c.Param("id")

	var req KrsItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify KRS exists
	krs, err := h.repo.GetKRSByID(krsID)
	if err != nil || krs == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "KRS tidak ditemukan").WithContext(c))
		return
	}

	if krs.Status != "draft" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Hanya bisa menambahkan item ke KRS berstatus draft").WithContext(c))
		return
	}

	item := domain.KrsItem{
		KrsID:      krsID,
		ClassID:    req.ClassID,
		Status:     "selected",
		SelectedAt: time.Now(),
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menambahkan item KRS").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(item).WithContext(c))
}

// DeleteKrsItem removes an item from KRS
// DELETE /api/v1/academic/krs/:id/items/:item_id
func (h *KrsHandler) DeleteKrsItem(c *gin.Context) {
	krsID := c.Param("id")
	itemID := c.Param("item_id")

	// Verify KRS exists
	krs, err := h.repo.GetKRSByID(krsID)
	if err != nil || krs == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "KRS tidak ditemukan").WithContext(c))
		return
	}

	if krs.Status != "draft" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Hanya bisa menghapus item dari KRS berstatus draft").WithContext(c))
		return
	}

	if err := h.db.Where("id = ?", itemID).Delete(&domain.KrsItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus item KRS").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Item KRS berhasil dihapus").WithContext(c))
}

// GetAvailableClasses retrieves classes available for a student in a period
// GET /api/v1/academic/krs/available-classes
func (h *KrsHandler) GetAvailableClasses(c *gin.Context) {
	studentID := c.Query("student_id")
	academicPeriodID := c.Query("academic_period_id")
	studyProgramID := c.Query("study_program_id")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "student_id is required").WithContext(c))
		return
	}

	if academicPeriodID == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "academic_period_id is required").WithContext(c))
		return
	}

	// Get classes that are available for the student's study program
	var classes []domain.Class
	query := h.db.Table("classes").
		Select("classes.id, classes.class_code, classes.quota, classes.enrolled_count, courses.course_code, courses.course_name, courses.sks").
		Joins("JOIN course_offerings ON course_offerings.id = classes.course_offering_id").
		Joins("JOIN courses ON courses.id = course_offerings.course_id").
		Where("course_offerings.academic_period_id = ?", academicPeriodID).
		Where("classes.class_status = 'active'").
		Where("classes.quota > classes.enrolled_count")

	if studyProgramID != "" {
		query = query.Where("courses.study_program_id = ?", studyProgramID)
	}

	if err := query.Find(&classes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil kelas tersedia").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(classes).WithContext(c))
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
