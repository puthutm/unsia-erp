package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-pmb-service/internal/domain"
	"github.com/unsia-erp/unsia-pmb-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Selection Handlers ============

type SelectionCreateRequest struct {
	SelectionName   string `json:"selection_name" binding:"required"`
	SelectionType   string `json:"selection_type" binding:"required,oneof=written_test_health_check interview"`
	StartDate      string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate       string `json:"end_date" binding:"required"`   // YYYY-MM-DD
	Location      string `json:"location"`
	Capacity     *int   `json:"capacity"`
	IsActive     *bool  `json:"is_active"`
}

type SelectionUpdateRequest struct {
	SelectionName *string `json:"selection_name"`
	SelectionType *string `json:"selection_type"`
	StartDate     *string `json:"start_date"`
	EndDate      *string `json:"end_date"`
	Location     *string `json:"location"`
	Capacity    *int    `json:"capacity"`
	IsActive     *bool   `json:"is_active"`
}

type SelectionResultRequest struct {
	SelectionID   string  `json:"selection_id" binding:"required"`
	ApplicantID  string  `json:"applicant_id" binding:"required"`
	Score        float64 `json:"score"`
	ResultStatus string  `json:"result_status" binding:"required,oneof=passed failed absent"`
	Note        *string `json:"note"`
}

type SelectionHandler struct {
	repo *repository.PMBRepository
	db   *gorm.DB
}

func NewSelectionHandler(db *gorm.DB) *SelectionHandler {
	return &SelectionHandler{
		repo: repository.NewPMBRepository(db),
		db:   db,
	}
}

// CreateSelection handles POST /api/v1/selections
func (h *SelectionHandler) CreateSelection(c *gin.Context) {
	var req SelectionCreateRequest
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

	selection := domain.SelectionSchedule{
		SelectionName: req.SelectionName,
		SelectionType: req.SelectionType,
		StartDate:   startDate,
		EndDate:    endDate,
		Location:  req.Location,
		Capacity: req.Capacity,
		IsActive: isActive,
	}

	if err := h.db.Create(&selection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jadwal seleksi").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.selection.create",
		Module:      "pmb",
		ResourceType: "selection_schedule",
		ResourceID: selection.ID,
		NewValue:  selection,
	})

	c.JSON(http.StatusCreated, sharederr.Success(selection).WithContext(c))
}

// GetSelection handles GET /api/v1/selections/:id
func (h *SelectionHandler) GetSelection(c *gin.Context) {
	id := c.Param("id")
	selection, err := h.repo.GetSelectionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jadwal seleksi").WithContext(c))
		return
	}
	if selection == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal seleksi tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(selection).WithContext(c))
}

// GetSelections handles GET /api/v1/selections
func (h *SelectionHandler) GetSelections(c *gin.Context) {
	filter := repository.SelectionFilter{
		SelectionType: c.Query("selection_type"),
		IsActive:     c.Query("is_active") == "true",
		Page:         1,
		Limit:        20,
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

	selections, total, err := h.repo.GetSelections(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar jadwal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"data":  selections,
		"total": total,
		"page": filter.Page,
		"limit": filter.Limit,
	}).WithContext(c))
}

// UpdateSelection handles PUT /api/v1/selections/:id
func (h *SelectionHandler) UpdateSelection(c *gin.Context) {
	id := c.Param("id")
	var req SelectionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	selection, err := h.repo.GetSelectionByID(id)
	if err != nil || selection == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal seleksi tidak ditemukan").WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.SelectionName != nil {
		updates["selection_name"] = *req.SelectionName
	}
	if req.SelectionType != nil {
		updates["selection_type"] = *req.SelectionType
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
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Capacity != nil {
		updates["capacity"] = *req.Capacity
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.repo.UpdateSelection(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui jadwal").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.selection.update",
		Module:      "pmb",
		ResourceType: "selection_schedule",
		ResourceID:  id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Jadwal seleksi berhasil diperbarui").WithContext(c))
}

// DeleteSelection handles DELETE /api/v1/selections/:id
func (h *SelectionHandler) DeleteSelection(c *gin.Context) {
	id := c.Param("id")

	selection, err := h.repo.GetSelectionByID(id)
	if err != nil || selection == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal seleksi tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.DeleteSelection(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus jadwal").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.selection.delete",
		Module:      "pmb",
		ResourceType: "selection_schedule",
		ResourceID:  id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Jadwal seleksi berhasil dihapus").WithContext(c))
}

// RegisterApplicantToSelection handles POST /api/v1/selections/:id/register
func (h *SelectionHandler) RegisterApplicantToSelection(c *gin.Context) {
	selectionID := c.Param("id")

	var req struct {
		ApplicantID string `json:"applicant_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify selection exists
	selection, err := h.repo.GetSelectionByID(selectionID)
	if err != nil || selection == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal seleksi tidak ditemukan").WithContext(c))
		return
	}

	// Verify applicant exists
	applicant, err := h.repo.GetApplicantByID(req.ApplicantID)
	if err != nil || applicant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Applicant tidak ditemukan").WithContext(c))
		return
	}

	registration := domain.SelectionRegistration{
		SelectionID:  selectionID,
		ApplicantID: req.ApplicantID,
		RegisterStatus: "registered",
		RegisteredAt:  time.Now(),
	}

	if err := h.db.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mendaftarkan applicant").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(registration).WithContext(c))
}

// RecordSelectionResult handles POST /api/v1/selections/:id/results
func (h *SelectionHandler) RecordSelectionResult(c *gin.Context) {
	selectionID := c.Param("id")

	var req SelectionResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify selection and applicant
	selection, err := h.repo.GetSelectionByID(selectionID)
	if err != nil || selection == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal seleksi tidak ditemukan").WithContext(c))
		return
	}

	applicant, err := h.repo.GetApplicantByID(req.ApplicantID)
	if err != nil || applicant == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Applicant tidak ditemukan").WithContext(c))
		return
	}

	result := domain.SelectionResult{
		SelectionID:  selectionID,
		ApplicantID: req.ApplicantID,
		Score:       req.Score,
		ResultStatus: req.ResultStatus,
		Note:       req.Note,
		RecordedAt:  time.Now(),
	}

	if err := h.db.Create(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan hasil seleksi").WithContext(c))
		return
	}

	// Update applicant status based on result
	newStatus := "selection_" + selection.SelectionType
	if req.ResultStatus == "passed" {
		newStatus = "passed_" + selection.SelectionType
	}

	history := domain.ApplicantStatusHistory{
		ApplicantID: req.ApplicantID,
		OldStatus:  applicant.Status,
		NewStatus:  newStatus,
		Note:      fmt.Sprintf("Selection %s: %s", selection.SelectionName, req.ResultStatus),
	}
	h.db.Create(&history)

	h.db.Model(applicant).Update("status", newStatus)

	// Publish event
	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	envelope := sharedevent.EventEnvelope{
		EventName:        "pmb.selection_completed",
		EventVersion:    "v1",
		PublisherService: "pmb-service",
		AggregateType:  "applicant",
		AggregateID:    req.ApplicantID,
		CorrelationID:  cid,
		Payload: map[string]interface{}{
			"applicant_id":  req.ApplicantID,
			"selection_id": selectionID,
			"result_status": req.ResultStatus,
			"score":       req.Score,
		},
	}
	conn := c.Request.Context().Value("db_conn").(interface{ ConnPool() interface{} })
	sharedevent.WriteOutbox(c.Request.Context(), conn.(interface{ ConnPool() interface{} }).(interface{ WriteOutbox(ctx interface{}, conn interface{}, envelope interface{}, topic string) (int64, error) }).(func(interface{}, interface{}, interface{}, string) (int64, error)), conn, envelope, "INTEGRATION_EVENT")

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.selection_result.record",
		Module:      "pmb",
		ResourceType: "selection_result",
		ResourceID: result.ID,
		NewValue:   result,
	})

	c.JSON(http.StatusCreated, sharederr.Success(result).WithContext(c))
}

// GetSelectionResults handles GET /api/v1/selections/:id/results
func (h *SelectionHandler) GetSelectionResults(c *gin.Context) {
	selectionID := c.Param("id")
	resultStatus := c.Query("result_status")

	results, err := h.repo.GetSelectionResults(selectionID, resultStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil hasil seleksi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(results).WithContext(c))
}

// GetApplicantSelections handles GET /api/v1/applicants/:id/selections
func (h *SelectionHandler) GetApplicantSelections(c *gin.Context) {
	applicantID := c.Param("id")

	results, err := h.repo.GetApplicantSelections(applicantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil riwayat seleksi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(results).WithContext(c))
}

// BatchRecordResults handles POST /api/v1/selections/:id/batch-results
func (h *SelectionHandler) BatchRecordResults(c *gin.Context) {
	selectionID := c.Param("id")

	var req struct {
		Results []SelectionResultRequest `json:"results" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		for _, r := range req.Results {
			result := domain.SelectionResult{
				SelectionID:  selectionID,
				ApplicantID: r.ApplicantID,
				Score:       r.Score,
				ResultStatus: r.ResultStatus,
				Note:        r.Note,
				RecordedAt:  time.Now(),
			}
			if err := tx.Create(&result).Error; err != nil {
				continue
			}

			// Update applicant status
			newStatus := "selection_" + r.ResultStatus
			if r.ResultStatus == "passed" {
				newStatus = "passed"
			}

			tx.Model(&domain.Applicant{}).Where("id = ?", r.ApplicantID).Update("status", newStatus)
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan hasil").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Hasil berhasil dicatat").WithContext(c))
}
