package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

type ApplicantCreateRequest struct {
	PersonID            string  `json:"person_id" binding:"required"`
	CrmLeadID           *string `json:"crm_lead_id"`
	StudyProgramID      *string `json:"study_program_id"`
	PmbWaveID           *string `json:"pmb_wave_id"`
	AdmissionPathID     *string `json:"admission_path_id"`
	TargetEntryPeriodID *string `json:"target_entry_period_id"`
}

type ApplicantHandler struct {
	repo *repository.PMBRepository
	db   *gorm.DB
}

func NewApplicantHandler(db *gorm.DB) *ApplicantHandler {
	return &ApplicantHandler{
		repo: repository.NewPMBRepository(db),
		db:   db,
	}
}

func generateRegNumber() string {
	now := time.Now().Format("2006")
	nBig, _ := rand.Int(rand.Reader, big.NewInt(90000))
	num := nBig.Int64() + 10000
	return fmt.Sprintf("PMB%s%d", now, num)
}

// CreateApplicant handles POST /api/v1/applicants
func (h *ApplicantHandler) CreateApplicant(c *gin.Context) {
	var req ApplicantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	applicant := domain.Applicant{
		PersonID:            req.PersonID,
		CrmLeadID:           req.CrmLeadID,
		StudyProgramID:      req.StudyProgramID,
		PmbWaveID:           req.PmbWaveID,
		AdmissionPathID:     req.AdmissionPathID,
		TargetEntryPeriodID: req.TargetEntryPeriodID,
		RegistrationNumber:  generateRegNumber(),
		Status:              "draft",
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&applicant).Error; err != nil {
			return err
		}

		// Initial empty biodata
		biodata := domain.ApplicantBiodata{
			ApplicantID:    applicant.ID,
			CoreSyncStatus: "pending",
		}
		if err := tx.Create(&biodata).Error; err != nil {
			return err
		}

		// Publish Event
		envelope := sharedevent.EventEnvelope{
			EventName:        "pmb.applicant_created",
			EventVersion:     "v1",
			PublisherService: "pmb-service",
			AggregateType:    "applicant",
			AggregateID:      applicant.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"applicant_id":        applicant.ID,
				"person_id":           applicant.PersonID,
				"registration_number": applicant.RegistrationNumber,
				"study_program_id":    applicant.StudyProgramID,
				"status":              applicant.Status,
			},
		}

		conn := tx.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan pendaftaran").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.applicant.create",
		Module:       "pmb",
		ResourceType: "applicant",
		ResourceID:   applicant.ID,
		NewValue:     applicant,
	})

	c.JSON(http.StatusCreated, sharederr.Success(applicant).WithContext(c))
}

// GetApplicant handles GET /api/v1/applicants/:id
func (h *ApplicantHandler) GetApplicant(c *gin.Context) {
	id := c.Param("id")
	app, err := h.repo.GetApplicantByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pendaftaran").WithContext(c))
		return
	}
	if app == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(app).WithContext(c))
}

// SubmitApplicant handles POST /api/v1/applicants/:id/submit
func (h *ApplicantHandler) SubmitApplicant(c *gin.Context) {
	id := c.Param("id")
	app, err := h.repo.GetApplicantByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pendaftaran").WithContext(c))
		return
	}
	if app == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	err = h.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		app.Status = "submitted"
		app.SubmittedAt = &now
		app.UpdatedAt = now

		if err := tx.Model(app).Updates(map[string]interface{}{
			"status":        "submitted",
			"submitted_at": &now,
			"updated_at":  now,
		}).Error; err != nil {
			return err
		}

		history := domain.ApplicantStatusHistory{
			ApplicantID: app.ID,
			NewStatus:   "submitted",
			Note:       "Applicant submitted forms and documents",
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		envelope := sharedevent.EventEnvelope{
			EventName:        "pmb.applicant_submitted",
			EventVersion:     "v1",
			PublisherService: "pmb-service",
			AggregateType:    "applicant",
			AggregateID:      app.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"applicant_id": app.ID,
				"submitted_at": now,
			},
		}

		conn := tx.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal melakukan submit pendaftaran").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(app).WithContext(c))
}

// GetApplicants handles GET /api/v1/applicants
func (h *ApplicantHandler) GetApplicants(c *gin.Context) {
	filter := repository.ApplicantListFilter{
		Status:          c.Query("status"),
		StudyProgramID:  c.Query("study_program_id"),
		PmbWaveID:     c.Query("pmb_wave_id"),
		AdmissionPathID: c.Query("admission_path_id"),
		Search:        c.Query("search"),
		Page:          1,
		Limit:         20,
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

	applicants, total, err := h.repo.GetApplicants(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar pendaftaran").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"data":  applicants,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	}).WithContext(c))
}

// ReceiveAssessmentSelectionResult handles POST /api/v1/applicants/:id/selection-result
func (h *ApplicantHandler) ReceiveAssessmentSelectionResult(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Score        float64 `json:"score" binding:"required"`
		ResultStatus string  `json:"result_status" binding:"required"` // pass, fail
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	app, err := h.repo.GetApplicantByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pendaftaran").WithContext(c))
		return
	}
	if app == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pendaftaran tidak ditemukan").WithContext(c))
		return
	}

	note := fmt.Sprintf("Assessment selection result: score %.2f, status %s", req.Score, req.ResultStatus)
	history := domain.ApplicantStatusHistory{
		ApplicantID: app.ID,
		NewStatus:   app.Status,
		Note:      note,
	}

	if err := h.db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mencatat hasil seleksi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Hasil seleksi berhasil dicatat").WithContext(c))
}
