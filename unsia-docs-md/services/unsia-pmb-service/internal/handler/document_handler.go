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

type DocumentUploadRequest struct {
	DocumentType string `json:"document_type" binding:"required"` // ktp, ijazah, pas_foto, rapor, dll
	FileURL      string `json:"file_url" binding:"required"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
}

type DocumentVerifyRequest struct {
	VerificationStatus string  `json:"verification_status" binding:"required,oneof=approved rejected"`
	RejectionReason   *string `json:"rejection_reason"`
	Note             *string `json:"note"`
	VerifiedBy       string  `json:"verified_by" binding:"required"`
}

type DocumentHandler struct {
	repo *repository.PMBRepository
	db   *gorm.DB
}

func NewDocumentHandler(db *gorm.DB) *DocumentHandler {
	return &DocumentHandler{
		repo: repository.NewPMBRepository(db),
		db:   db,
	}
}

// UploadDocument handles POST /api/v1/applicants/:id/documents
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	applicantID := c.Param("id")

	var req DocumentUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Verify applicant exists
	app, err := h.repo.GetApplicantByID(applicantID)
	if err != nil || app == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Applicant tidak ditemukan").WithContext(c))
		return
	}

	doc := domain.ApplicantDocument{
		ApplicantID:    applicantID,
		DocumentType:  req.DocumentType,
		FileURL:       req.FileURL,
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		UploadStatus:  "uploaded",
		UploadedAt:    time.Now(),
	}

	if err := h.db.Create(&doc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengunggah dokumen").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.document.upload",
		Module:      "pmb",
		ResourceType: "applicant_document",
		ResourceID:  doc.ID,
		NewValue:    doc,
	})

	c.JSON(http.StatusCreated, sharederr.Success(doc).WithContext(c))
}

// VerifyDocument handles PUT /api/v1/documents/:id/verify
func (h *DocumentHandler) VerifyDocument(c *gin.Context) {
	docID := c.Param("id")

	var req DocumentVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var doc domain.ApplicantDocument
	if err := h.db.Where("id = ?", docID).First(&doc).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Dokumen tidak ditemukan").WithContext(c))
		return
	}

	now := time.Now()
	verification := domain.DocumentVerification{
		DocumentID:         doc.ID,
		VerifiedBy:        req.VerifiedBy,
		VerificationStatus: req.VerificationStatus,
		RejectionReason:    req.RejectionReason,
		Note:            req.Note,
		VerifiedAt:       &now,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Create verification record
		if err := tx.Create(&verification).Error; err != nil {
			return err
		}

		// Update document status
		status := "approved"
		if req.VerificationStatus == "rejected" {
			status = "rejected"
		}
		updateDoc := map[string]interface{}{
			"verification_status": status,
			"verified_at":         &now,
			"verified_by":       req.VerifiedBy,
		}
		if err := tx.Model(&doc).Updates(updateDoc).Error; err != nil {
			return err
		}

		// Update applicant biodata sync status (if all docs verified)
		if req.VerificationStatus == "approved" {
			h.updateBiodataSyncStatus(tx, doc.ApplicantID)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memverifikasi dokumen").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "pmb.document.verify",
		Module:      "pmb",
		ResourceType: "document_verification",
		ResourceID:  verification.ID,
		NewValue:    verification,
	})

	c.JSON(http.StatusOK, sharederr.Success(verification).WithContext(c))
}

// BatchVerifyDocuments handles PUT /api/v1/applicants/:id/documents/batch-verify
func (h *DocumentHandler) BatchVerifyDocuments(c *gin.Context) {
	applicantID := c.Param("id")

	var req struct {
		Documents []struct {
			DocumentID          string  `json:"document_id" binding:"required"`
			VerificationStatus string  `json:"verification_status" binding:"required,oneof=approved rejected"`
			RejectionReason    *string `json:"rejection_reason"`
		} `json:"documents" binding:"required,gt=0"`
		VerifiedBy string `json:"verified_by" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	now := time.Now()

	err := h.db.Transaction(func(tx *gorm.DB) error {
		for _, docReq := range req.Documents {
			var doc domain.ApplicantDocument
			if err := tx.Where("id = ? AND applicant_id = ?", docReq.DocumentID, applicantID).First(&doc).Error; err != nil {
				continue // skip invalid
			}

			verification := domain.DocumentVerification{
				DocumentID:          doc.ID,
				VerifiedBy:        req.VerifiedBy,
				VerificationStatus: docReq.VerificationStatus,
				RejectionReason:    docReq.RejectionReason,
				VerifiedAt:       &now,
			}
			if err := tx.Create(&verification).Error; err != nil {
				continue
			}

			status := "approved"
			if docReq.VerificationStatus == "rejected" {
				status = "rejected"
			}
			tx.Model(&doc).Updates(map[string]interface{}{
				"verification_status": status,
				"verified_at":         &now,
				"verified_by":         req.VerifiedBy,
			})
		}

		// Update biodata sync status
		h.updateBiodataSyncStatus(tx, applicantID)

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memverifikasi dokumen").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Dokumen berhasil diverifikasi").WithContext(c))
}

// GetDocuments handles GET /api/v1/applicants/:id/documents
func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	applicantID := c.Param("id")
	docType := c.Query("document_type")

	docs, err := h.repo.GetApplicantDocuments(applicantID, docType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil dokumen").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(docs).WithContext(c))
}

// GetPendingVerifications handles GET /api/v1/documents/pending-verification
func (h *DocumentHandler) GetPendingVerifications(c *gin.Context) {
	filter := repository.DocumentFilter{
		VerificationStatus: "pending",
		Page:            1,
		Limit:           20,
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

	docs, total, err := h.repo.GetPendingVerificationDocuments(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"data":  docs,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	}).WithContext(c))
}

func (h *DocumentHandler) updateBiodataSyncStatus(tx *gorm.DB, applicantID string) {
	var pendingDocs int64
	tx.Model(&domain.ApplicantDocument{}).
		Where("applicant_id = ? AND verification_status != 'approved'", applicantID).
		Count(&pendingDocs)

	if pendingDocs == 0 {
		tx.Model(&domain.ApplicantBiodata{}).
			Where("applicant_id = ?", applicantID).
			Updates(map[string]interface{}{
				"core_sync_status": "ready",
				"synced_at":       time.Now(),
			})

		// Publish event
		envelope := sharedevent.EventEnvelope{
			EventName:        "pmb.applicant_documents_completed",
			EventVersion:    "v1",
			PublisherService: "pmb-service",
			AggregateType:  "applicant",
			AggregateID:    applicantID,
			Payload: map[string]interface{}{
				"applicant_id": applicantID,
				"status":      "documents_verified",
			},
		}
		conn := tx.Statement.ConnPool
		sharedevent.WriteOutbox(nil, conn, envelope, "INTEGRATION_EVENT")
	}
}
