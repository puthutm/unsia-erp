package handler

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	sharedhttpclient "github.com/unsia-erp/shared-httpclient"
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

type DocumentUploadRequest struct {
	DocumentTypeCode string `json:"document_type_code" binding:"required"`
	FileUrl          string `json:"file_url" binding:"required"`
}

type DocumentVerifyRequest struct {
	VerificationStatus string  `json:"verification_status" binding:"required,oneof=verified rejected"`
	RejectReason       *string `json:"reject_reason"`
}

type SelectionResultRequest struct {
	Score        float64 `json:"score" binding:"required"`
	ResultStatus string  `json:"result_status" binding:"required"` // pass, fail
}

type PMBHandler struct {
	repo           *repository.PMBRepository
	db             *gorm.DB
	financeClient  *sharedhttpclient.Client
	academicClient *sharedhttpclient.Client
}

func NewPMBHandler(db *gorm.DB) *PMBHandler {
	financeURL := os.Getenv("FINANCE_SERVICE_URL")
	if financeURL == "" {
		financeURL = "http://localhost:8005"
	}
	academicURL := os.Getenv("ACADEMIC_SERVICE_URL")
	if academicURL == "" {
		academicURL = "http://localhost:8006"
	}
	srvToken := os.Getenv("PMB_SERVICE_TOKEN")
	if srvToken == "" {
		srvToken = "pmb_service_secret_token"
	}

	financeClient := sharedhttpclient.New(sharedhttpclient.Config{
		BaseURL:      financeURL,
		ServiceToken: srvToken,
		SourceName:   "pmb-service",
		Timeout:      10 * time.Second,
	})

	academicClient := sharedhttpclient.New(sharedhttpclient.Config{
		BaseURL:      academicURL,
		ServiceToken: srvToken,
		SourceName:   "pmb-service",
		Timeout:      10 * time.Second,
	})

	return &PMBHandler{
		repo:           repository.NewPMBRepository(db),
		db:             db,
		financeClient:  financeClient,
		academicClient: academicClient,
	}
}

func generateRegNumber() string {
	now := time.Now().Format("2006")
	nBig, _ := rand.Int(rand.Reader, big.NewInt(90000))
	num := nBig.Int64() + 10000
	return fmt.Sprintf("PMB%s%d", now, num)
}

func (h *PMBHandler) CreateApplicant(c *gin.Context) {
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

func (h *PMBHandler) GetApplicant(c *gin.Context) {
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

func (h *PMBHandler) SubmitApplicant(c *gin.Context) {
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
			"status":       "submitted",
			"submitted_at": &now,
			"updated_at":   now,
		}).Error; err != nil {
			return err
		}

		history := domain.ApplicantStatusHistory{
			ApplicantID: app.ID,
			NewStatus:   "submitted",
			Note:        "Applicant submitted forms and documents",
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

func (h *PMBHandler) UploadDocument(c *gin.Context) {
	id := c.Param("id")
	var req DocumentUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	doc := domain.ApplicantDocument{
		ApplicantID:        id,
		DocumentTypeCode:   req.DocumentTypeCode,
		FileUrl:            req.FileUrl,
		VerificationStatus: "pending",
	}

	if err := h.repo.CreateApplicantDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan dokumen").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(doc).WithContext(c))
}

func (h *PMBHandler) VerifyDocument(c *gin.Context) {
	docID := c.Param("doc_id")
	var req DocumentVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	reason := ""
	if req.RejectReason != nil {
		reason = *req.RejectReason
	}

	if err := h.repo.UpdateApplicantDocumentVerification(docID, req.VerificationStatus, actor, reason); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate verifikasi dokumen").WithContext(c))
		return
	}

	doc, _ := h.repo.GetApplicantDocumentByID(docID)
	c.JSON(http.StatusOK, sharederr.Success(doc).WithContext(c))
}

func (h *PMBHandler) RequestInvoice(c *gin.Context) {
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

	// Prepare payload for Finance Invoice
	invoicePayload := map[string]interface{}{
		"payer_type":         "applicant",
		"payer_ref_id":       app.ID,
		"academic_period_id": app.TargetEntryPeriodID,
		"source_module":      "pmb",
		"source_ref_id":      app.ID,
		"due_date":           time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
		"items": []map[string]interface{}{
			{
				"payment_component_code": "PMB_REG_FEE",
				"amount":                 250000.00,
			},
		},
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

	resp, err := h.financeClient.Post(ctx, "/api/v1/finance/invoices", invoicePayload)
	if err != nil {
		c.JSON(http.StatusBadGateway, sharederr.Error("FINANCE_SERVICE_UNAVAILABLE", fmt.Sprintf("Gagal menghubungi layanan Finance: %v", err)).WithContext(c))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.JSON(http.StatusBadGateway, sharederr.Error("FINANCE_SERVICE_ERROR", fmt.Sprintf("Layanan Finance mengembalikan error: %d", resp.StatusCode)).WithContext(c))
		return
	}

	var invoiceResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID            string  `json:"invoice_id"`
			InvoiceNumber string  `json:"invoice_number"`
			AmountTotal   float64 `json:"amount_total"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&invoiceResp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("PARSE_ERROR", "Gagal membaca response Finance").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(invoiceResp.Data).WithContext(c))
}

func (h *PMBHandler) IssueLoa(c *gin.Context) {
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

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	loa := domain.LoaDocument{
		ApplicantID: app.ID,
		LoaNumber:   fmt.Sprintf("LOA/%s/%s", app.RegistrationNumber, time.Now().Format("20060102")),
		FileUrl:     fmt.Sprintf("http://storage.unsia.ac.id/loa/%s.pdf", app.RegistrationNumber),
		IssuedBy:    actor,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&loa).Error; err != nil {
			return err
		}

		app.Status = "accepted"
		if err := tx.Model(app).Updates(map[string]interface{}{
			"status":      "accepted",
			"accepted_at": time.Now(),
			"updated_at":  time.Now(),
		}).Error; err != nil {
			return err
		}

		history := domain.ApplicantStatusHistory{
			ApplicantID: app.ID,
			NewStatus:   "accepted",
			Note:        "Letter of Acceptance issued. Status updated to accepted.",
		}
		return tx.Create(&history).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menerbitkan LoA").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(loa).WithContext(c))
}

func (h *PMBHandler) HandoverToAcademic(c *gin.Context) {
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

	idempotencyKey := fmt.Sprintf("handover-%s", app.ID)

	// Check if already handed over
	existing, err := h.repo.GetHandoverLogByIdempotencyKey(idempotencyKey)
	if err == nil && existing != nil && existing.Status == "success" {
		c.JSON(http.StatusOK, sharederr.Success(gin.H{
			"student_id":      existing.StudentRefID,
			"nim":             existing.Nim,
			"handover_status": "success",
		}).WithContext(c))
		return
	}

	logEntry := domain.HandoverLog{
		ApplicantID:    app.ID,
		Status:         "pending",
		IdempotencyKey: idempotencyKey,
	}
	_ = h.repo.CreateHandoverLog(&logEntry)

	// Call Academic Service
	academicPayload := map[string]interface{}{
		"applicant_id":           app.ID,
		"curriculum_id":          "curriculum-default-id", // placeholder
		"entry_academic_year_id": app.PmbWaveID,           // mapping placeholders
		"entry_period_id":        app.TargetEntryPeriodID,
		"study_program_id":       app.StudyProgramID,
		"reason":                 "PMB registration completed",
	}

	ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

	resp, err := h.academicClient.Post(ctx, "/api/v1/academic/students/generate-from-applicant", academicPayload)
	if err != nil {
		_ = h.repo.UpdateHandoverLog(logEntry.ID, "failed", nil, nil, err.Error())
		c.JSON(http.StatusBadGateway, sharederr.Error("ACADEMIC_SERVICE_UNAVAILABLE", fmt.Sprintf("Gagal menghubungi layanan Akademik: %v", err)).WithContext(c))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		_ = h.repo.UpdateHandoverLog(logEntry.ID, "failed", nil, nil, fmt.Sprintf("Status code: %d", resp.StatusCode))
		c.JSON(http.StatusBadGateway, sharederr.Error("ACADEMIC_SERVICE_ERROR", fmt.Sprintf("Layanan Akademik mengembalikan error: %d", resp.StatusCode)).WithContext(c))
		return
	}

	var academicResp struct {
		Success bool `json:"success"`
		Data    struct {
			StudentID      string `json:"student_id"`
			Nim            string `json:"nim"`
			HandoverStatus string `json:"handover_status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&academicResp); err != nil {
		_ = h.repo.UpdateHandoverLog(logEntry.ID, "failed", nil, nil, err.Error())
		c.JSON(http.StatusInternalServerError, sharederr.Error("PARSE_ERROR", "Gagal membaca response Akademik").WithContext(c))
		return
	}

	// Update Log and applicant status in transaction
	err = h.db.Transaction(func(tx *gorm.DB) error {
		app.Status = "ready_for_academic"
		if err := tx.Model(app).Updates(map[string]interface{}{
			"status":     "ready_for_academic",
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		history := domain.ApplicantStatusHistory{
			ApplicantID: app.ID,
			NewStatus:   "ready_for_academic",
			Note:        "Handed over to academic, assigned NIM: " + academicResp.Data.Nim,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		if err := tx.Model(&domain.HandoverLog{}).Where("id = ?", logEntry.ID).Updates(map[string]interface{}{
			"status":         "success",
			"student_ref_id": academicResp.Data.StudentID,
			"nim":            academicResp.Data.Nim,
		}).Error; err != nil {
			return err
		}

		envelope := sharedevent.EventEnvelope{
			EventName:        "pmb.ready_for_academic",
			EventVersion:     "v1",
			PublisherService: "pmb-service",
			AggregateType:    "applicant",
			AggregateID:      app.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"applicant_id": app.ID,
				"student_id":   academicResp.Data.StudentID,
				"nim":          academicResp.Data.Nim,
			},
		}

		conn := tx.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(ctx, conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses transaksi handover").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(academicResp.Data).WithContext(c))
}

func (h *PMBHandler) ReceiveAssessmentSelectionResult(c *gin.Context) {
	id := c.Param("id")
	var req SelectionResultRequest
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
		Note:        note,
	}

	if err := h.db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mencatat hasil seleksi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Hasil seleksi berhasil dicatat").WithContext(c))
}

// Applicant List with Filtering and Pagination
func (h *PMBHandler) GetApplicants(c *gin.Context) {
	filter := repository.ApplicantListFilter{
		Status:          c.Query("status"),
		StudyProgramID:  c.Query("study_program_id"),
		PmbWaveID:       c.Query("pmb_wave_id"),
		AdmissionPathID: c.Query("admission_path_id"),
		Search:          c.Query("search"),
		Page:            1,
		Limit:           20,
	}

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &filter.Page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &filter.Limit)
	}

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

// Get Applicant Biodata
func (h *PMBHandler) GetBiodata(c *gin.Context) {
	id := c.Param("id")

	biodata, err := h.repo.GetBiodataByApplicantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil biodata").WithContext(c))
		return
	}
	if biodata == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Biodata tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(biodata).WithContext(c))
}

// Update Applicant Biodata
type BiodataUpdateRequest struct {
	FullName        string  `json:"full_name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	Nik             string  `json:"nik"`
	Nip             string  `json:"nip"`
	BirthPlace     string  `json:"birth_place"`
	BirthDate      string  `json:"birth_date"`
	Gender          string  `json:"gender"`
	ReligionID      string  `json:"religion_id"`
	MaritalStatus   string  `json:"marital_status"`
	Citizenship    string  `json:"citizenship"`
	JacketSize     string  `json:"jacket_size"`
}

func (h *PMBHandler) UpdateBiodata(c *gin.Context) {
	id := c.Param("id")
	var req BiodataUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := map[string]interface{}{}
	if req.FullName != "" {
		updates["full_name"] = req.FullName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Nik != "" {
		updates["nik"] = req.Nik
	}
	if req.Nip != "" {
		updates["nip"] = req.Nip
	}
	if req.BirthPlace != "" {
		updates["birth_place"] = req.BirthPlace
	}
	if req.BirthDate != "" {
		updates["birth_date"] = req.BirthDate
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.ReligionID != "" {
		updates["religion_id"] = req.ReligionID
	}
	if req.MaritalStatus != "" {
		updates["marital_status"] = req.MaritalStatus
	}
	if req.Citizenship != "" {
		updates["citizenship"] = req.Citizenship
	}
	if req.JacketSize != "" {
		updates["jacket_size"] = req.JacketSize
	}

	if err := h.repo.UpdateBiodata(id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate biodata").WithContext(c))
		return
	}

	biodata, _ := h.repo.GetBiodataByApplicantID(id)
	c.JSON(http.StatusOK, sharederr.Success(biodata).WithContext(c))
}

// Get Applicant Addresses
func (h *PMBHandler) GetAddresses(c *gin.Context) {
	id := c.Param("id")

	addresses, err := h.repo.GetAddressesByApplicantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil alamat").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(addresses).WithContext(c))
}

// Update Applicant Address
type AddressUpdateRequest struct {
	AddressType   string  `json:"address_type" binding:"required,oneof=KTPS DOMISILI"`
	Street       string  `json:"street"`
	ProvinceID  string  `json:"province_id"`
	CityID       string  `json:"city_id"`
	DistrictID  string  `json:"district_id"`
	VillageID   string  `json:"village_id"`
	PostalCode  string  `json:"postal_code"`
	IsSameAsKtp bool    `json:"is_same_as_ktp"`
}

func (h *PMBHandler) UpdateAddresses(c *gin.Context) {
	id := c.Param("id")
	var req AddressUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var provincePtr, cityPtr, districtPtr, villagePtr *string
	if req.ProvinceID != "" {
		provincePtr = &req.ProvinceID
	}
	if req.CityID != "" {
		cityPtr = &req.CityID
	}
	if req.DistrictID != "" {
		districtPtr = &req.DistrictID
	}
	if req.VillageID != "" {
		villagePtr = &req.VillageID
	}

	addr := domain.ApplicantAddress{
		ApplicantID:  id,
		AddressType:  req.AddressType,
		Street:       req.Street,
		ProvinceID:  provincePtr,
		CityID:       cityPtr,
		DistrictID:   districtPtr,
		VillageID:    villagePtr,
		PostalCode:   req.PostalCode,
		IsSameAsKtp:  req.IsSameAsKtp,
	}

	if err := h.repo.UpsertAddress(&addr); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan alamat").WithContext(c))
		return
	}

	addresses, _ := h.repo.GetAddressesByApplicantID(id)
	c.JSON(http.StatusOK, sharederr.Success(addresses).WithContext(c))
}

// Get Applicant Education Background
func (h *PMBHandler) GetEducation(c *gin.Context) {
	id := c.Param("id")

	educations, err := h.repo.GetEducationBackgroundsByApplicantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pendidikan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(educations).WithContext(c))
}

// Update Education Background
type EducationUpdateRequest struct {
	SchoolName     string `json:"school_name"`
	Major          string `json:"major"`
	GraduationYear int    `json:"graduation_year"`
	Gpa           string `json:"gpa"`
}

func (h *PMBHandler) UpdateEducation(c *gin.Context) {
	id := c.Param("id")
	var req EducationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var gpaVal float64
	if req.Gpa != "" {
		_, _ = fmt.Sscanf(req.Gpa, "%f", &gpaVal)
	}

	var gradYearStr string
	if req.GraduationYear > 0 {
		gradYearStr = fmt.Sprintf("%d", req.GraduationYear)
	}

	edu := domain.ApplicantEducationBackground{
		ApplicantID:     id,
		SchoolName:      req.SchoolName,
		Major:           req.Major,
		GraduationYear:  gradYearStr,
		Gpa:             gpaVal,
	}

	if err := h.repo.UpsertEducationBackground(&edu); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan pendidikan").WithContext(c))
		return
	}

	educations, _ := h.repo.GetEducationBackgroundsByApplicantID(id)
	c.JSON(http.StatusOK, sharederr.Success(educations).WithContext(c))
}

// Get Applicant Family Members
func (h *PMBHandler) GetFamily(c *gin.Context) {
	id := c.Param("id")

	members, err := h.repo.GetFamilyMembersByApplicantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil keluarga").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(members).WithContext(c))
}

// Update Family Member
type FamilyUpdateRequest struct {
	Relationship string  `json:"relationship" binding:"required,oneof=AYAH IBU WALI SAUDARA"`
	FullName     string  `json:"full_name"`
	Occupation  string  `json:"occupation"`
	Income      float64 `json:"income"`
}

func (h *PMBHandler) UpdateFamily(c *gin.Context) {
	id := c.Param("id")
	var req FamilyUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	member := domain.ApplicantFamilyMember{
		ApplicantID: id,
		Relationship: req.Relationship,
		FullName:     req.FullName,
		Occupation:  req.Occupation,
		Income:       req.Income,
	}

	if err := h.repo.UpsertFamilyMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan keluarga").WithContext(c))
		return
	}

	members, _ := h.repo.GetFamilyMembersByApplicantID(id)
	c.JSON(http.StatusOK, sharederr.Success(members).WithContext(c))
}

// Get Applicant Financial Profile
func (h *PMBHandler) GetFinancial(c *gin.Context) {
	id := c.Param("id")

	profiles, err := h.repo.GetFinancialProfilesByApplicantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil finansial").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(profiles).WithContext(c))
}

// Update Financial Profile
type FinancialUpdateRequest struct {
	SponsorType   string  `json:"sponsor_type" binding:"required,oneof=SISWA ORANGTUA KERJA BEASISWA"`
	SponsorName   string  `json:"sponsor_name"`
	MonthlyIncome float64 `json:"monthly_income"`
}

func (h *PMBHandler) UpdateFinancial(c *gin.Context) {
	id := c.Param("id")
	var req FinancialUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	profile := domain.ApplicantFinancialProfile{
		ApplicantID:    id,
		SponsorType:    req.SponsorType,
		SponsorName:    req.SponsorName,
		MonthlyIncome:  req.MonthlyIncome,
	}

	if err := h.repo.UpsertFinancialProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan finansial").WithContext(c))
		return
	}

	profiles, _ := h.repo.GetFinancialProfilesByApplicantID(id)
	c.JSON(http.StatusOK, sharederr.Success(profiles).WithContext(c))
}

// Get Applicant Facility Profile
func (h *PMBHandler) GetFacility(c *gin.Context) {
	id := c.Param("id")

	profiles, err := h.repo.GetFacilityProfilesByApplicantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil fasilitas").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(profiles).WithContext(c))
}

// Update Facility Profile
type FacilityUpdateRequest struct {
	FacilityType string `json:"facility_type" binding:"required,oneof=LAPTOP INTERNET"`
	Description  string `json:"description"`
}

func (h *PMBHandler) UpdateFacility(c *gin.Context) {
	id := c.Param("id")
	var req FacilityUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	profile := domain.ApplicantFacilityProfile{
		ApplicantID: id,
		FacilityType: req.FacilityType,
		Description:  req.Description,
	}

	if err := h.repo.UpsertFacilityProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan fasilitas").WithContext(c))
		return
	}

	profiles, _ := h.repo.GetFacilityProfilesByApplicantID(id)
	c.JSON(http.StatusOK, sharederr.Success(profiles).WithContext(c))
}

// Get Applicant Documents
func (h *PMBHandler) GetDocuments(c *gin.Context) {
	id := c.Param("id")

	documents, err := h.repo.GetDocumentsByApplicantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil dokumen").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(documents).WithContext(c))
}
