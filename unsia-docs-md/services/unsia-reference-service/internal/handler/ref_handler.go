package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	"github.com/unsia-erp/unsia-reference-service/internal/domain"
	"github.com/unsia-erp/unsia-reference-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type StudyProgramCreateRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Degree string `json:"degree" binding:"required"`
}

type AcademicYearCreateRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"` // Format: YYYY-MM-DD
	EndDate   string `json:"end_date" binding:"required"`   // Format: YYYY-MM-DD
}

type AcademicPeriodCreateRequest struct {
	AcademicYearID string `json:"academic_year_id" binding:"required"`
	Code           string `json:"code" binding:"required"`
	Term           string `json:"term" binding:"required"`       // odd, even, intermediate
	StartDate      string `json:"start_date" binding:"required"` // Format: YYYY-MM-DD
	EndDate        string `json:"end_date" binding:"required"`   // Format: YYYY-MM-DD
}

type ReferenceHandler struct {
	repo *repository.ReferenceRepository
	db   *gorm.DB
}

func NewReferenceHandler(db *gorm.DB) *ReferenceHandler {
	return &ReferenceHandler{
		repo: repository.NewReferenceRepository(db),
		db:   db,
	}
}

// ListStudyPrograms returns all study programs
func (h *ReferenceHandler) ListStudyPrograms(c *gin.Context) {
	prodis, err := h.repo.ListStudyPrograms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data program studi").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(prodis).WithContext(c))
}

// CreateStudyProgram creates a new study program and publishes outbox event
func (h *ReferenceHandler) CreateStudyProgram(c *gin.Context) {
	var req StudyProgramCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	prodi := domain.StudyProgram{
		Code:   req.Code,
		Name:   req.Name,
		Degree: req.Degree,
		Status: "ACTIVE",
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&prodi).Error; err != nil {
			return err
		}

		// Write Outbox Event
		envelope := sharedevent.EventEnvelope{
			EventName:        "reference.study_program_updated",
			EventVersion:     "v1",
			PublisherService: "reference-service",
			AggregateType:    "study_program",
			AggregateID:      prodi.ID,
			CorrelationID:    cid,
			Payload:          prodi,
		}

		conn := tx.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(context.Background(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("TRANSACTION_FAILED", "Gagal menyimpan data dan menerbitkan event").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(prodi).WithContext(c))
}

// ListAcademicYears returns all academic years
func (h *ReferenceHandler) ListAcademicYears(c *gin.Context) {
	years, err := h.repo.ListAcademicYears()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data tahun ajaran").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(years).WithContext(c))
}

// CreateAcademicYear creates a new academic year
func (h *ReferenceHandler) CreateAcademicYear(c *gin.Context) {
	var req AcademicYearCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE_FORMAT", "start_date format must be YYYY-MM-DD").WithContext(c))
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE_FORMAT", "end_date format must be YYYY-MM-DD").WithContext(c))
		return
	}

	year := domain.AcademicYear{
		Code:      req.Code,
		Name:      req.Name,
		Status:    "INACTIVE", // initially inactive
		StartDate: start,
		EndDate:   end,
	}

	if err := h.repo.CreateAcademicYear(&year); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan data tahun ajaran").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(year).WithContext(c))
}

// ListAcademicPeriods returns all academic periods
func (h *ReferenceHandler) ListAcademicPeriods(c *gin.Context) {
	periods, err := h.repo.ListAcademicPeriods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data periode akademik").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(periods).WithContext(c))
}

// CreateAcademicPeriod creates a new academic period and publishes outbox event
func (h *ReferenceHandler) CreateAcademicPeriod(c *gin.Context) {
	var req AcademicPeriodCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE_FORMAT", "start_date format must be YYYY-MM-DD").WithContext(c))
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE_FORMAT", "end_date format must be YYYY-MM-DD").WithContext(c))
		return
	}

	period := domain.AcademicPeriod{
		AcademicYearID: req.AcademicYearID,
		Code:           req.Code,
		Term:           req.Term,
		Status:         "INACTIVE",
		StartDate:      start,
		EndDate:        end,
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	txErr := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&period).Error; err != nil {
			return err
		}

		// Write Outbox Event
		envelope := sharedevent.EventEnvelope{
			EventName:        "reference.academic_period_updated",
			EventVersion:     "v1",
			PublisherService: "reference-service",
			AggregateType:    "academic_period",
			AggregateID:      period.ID,
			CorrelationID:    cid,
			Payload:          period,
		}

		conn := tx.Statement.ConnPool
		_, err := sharedevent.WriteOutbox(context.Background(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if txErr != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("TRANSACTION_FAILED", "Gagal menyimpan data dan menerbitkan event").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(period).WithContext(c))
}

// ListStatusCodes returns managed business status codes
func (h *ReferenceHandler) ListStatusCodes(c *gin.Context) {
	codes, err := h.repo.ListStatusCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data status codes").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(codes).WithContext(c))
}

// ListPaymentComponents returns active payment components
func (h *ReferenceHandler) ListPaymentComponents(c *gin.Context) {
	comps, err := h.repo.ListPaymentComponents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data komponen pembayaran").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(comps).WithContext(c))
}

// ListDocumentTypes returns active document requirements
func (h *ReferenceHandler) ListDocumentTypes(c *gin.Context) {
	docs, err := h.repo.ListDocumentTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data jenis dokumen").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(docs).WithContext(c))
}

// ListPaymentMethods returns active payment methods
func (h *ReferenceHandler) ListPaymentMethods(c *gin.Context) {
	methods, err := h.repo.ListPaymentMethods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data metode pembayaran").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(methods).WithContext(c))
}

type PmbWaveRequest struct {
	AcademicYearID           *string `json:"academic_year_id"`
	TargetEntryPeriodID      string  `json:"target_entry_period_id" binding:"required"`
	AdmissionPathID          *string `json:"admission_path_id"`
	Code                     string  `json:"code" binding:"required"`
	Name                     string  `json:"name" binding:"required"`
	StartDate                string  `json:"start_date"` // YYYY-MM-DD
	EndDate                  string  `json:"end_date"`   // YYYY-MM-DD
	RegistrationStartAt      string  `json:"registration_start_at"`
	RegistrationEndAt        string  `json:"registration_end_at"`
	SelectionStartAt         string  `json:"selection_start_at"`
	SelectionEndAt           string  `json:"selection_end_at"`
	ReregistrationDeadlineAt string  `json:"reregistration_deadline_at"`
	Status                   string  `json:"status" binding:"oneof=draft open closed archived"`
}

type PaymentComponentRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	DefaultAmount float64 `json:"default_amount" binding:"required"`
}

type PaymentMethodRequest struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Provider string `json:"provider"`
}

type DocumentTypeRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	IsMandatory bool   `json:"is_mandatory"`
}

type ReligionRequest struct {
	Name string `json:"name" binding:"required"`
}

type CountryRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// ListPmbWaves returns active PMB Waves
func (h *ReferenceHandler) ListPmbWaves(c *gin.Context) {
	waves, err := h.repo.ListPmbWaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data gelombang PMB").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(waves).WithContext(c))
}

// CreatePmbWave creates a new PMB Wave
func (h *ReferenceHandler) CreatePmbWave(c *gin.Context) {
	var req PmbWaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	existing, err := h.repo.GetPmbWaveByCode(req.Code)
	if err == nil && existing != nil {
		c.JSON(http.StatusConflict, sharederr.Error("WAVE_ALREADY_EXISTS", "Gelombang PMB dengan code tersebut sudah terdaftar").WithContext(c))
		return
	}

	wave := domain.PmbWave{
		AcademicYearID:      req.AcademicYearID,
		TargetEntryPeriodID: req.TargetEntryPeriodID,
		AdmissionPathID:     req.AdmissionPathID,
		Code:                req.Code,
		Name:                req.Name,
		Status:              req.Status,
		IsActive:            true,
	}

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			wave.StartDate = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			wave.EndDate = &t
		}
	}
	if req.RegistrationStartAt != "" {
		if t, err := time.Parse(time.RFC3339, req.RegistrationStartAt); err == nil {
			wave.RegistrationStartAt = &t
		}
	}
	if req.RegistrationEndAt != "" {
		if t, err := time.Parse(time.RFC3339, req.RegistrationEndAt); err == nil {
			wave.RegistrationEndAt = &t
		}
	}
	if req.SelectionStartAt != "" {
		if t, err := time.Parse(time.RFC3339, req.SelectionStartAt); err == nil {
			wave.SelectionStartAt = &t
		}
	}
	if req.SelectionEndAt != "" {
		if t, err := time.Parse(time.RFC3339, req.SelectionEndAt); err == nil {
			wave.SelectionEndAt = &t
		}
	}
	if req.ReregistrationDeadlineAt != "" {
		if t, err := time.Parse(time.RFC3339, req.ReregistrationDeadlineAt); err == nil {
			wave.ReregistrationDeadlineAt = &t
		}
	}

	if err := h.repo.CreatePmbWave(&wave); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan gelombang PMB").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(wave).WithContext(c))
}

// UpdatePmbWave updates a PMB Wave
func (h *ReferenceHandler) UpdatePmbWave(c *gin.Context) {
	id := c.Param("id")
	wave, err := h.repo.GetPmbWaveByID(id)
	if err != nil || wave == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Gelombang PMB tidak ditemukan").WithContext(c))
		return
	}

	var req PmbWaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	wave.AcademicYearID = req.AcademicYearID
	wave.TargetEntryPeriodID = req.TargetEntryPeriodID
	wave.AdmissionPathID = req.AdmissionPathID
	wave.Code = req.Code
	wave.Name = req.Name
	wave.Status = req.Status

	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			wave.StartDate = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			wave.EndDate = &t
		}
	}
	if req.RegistrationStartAt != "" {
		if t, err := time.Parse(time.RFC3339, req.RegistrationStartAt); err == nil {
			wave.RegistrationStartAt = &t
		}
	}
	if req.RegistrationEndAt != "" {
		if t, err := time.Parse(time.RFC3339, req.RegistrationEndAt); err == nil {
			wave.RegistrationEndAt = &t
		}
	}
	if req.SelectionStartAt != "" {
		if t, err := time.Parse(time.RFC3339, req.SelectionStartAt); err == nil {
			wave.SelectionStartAt = &t
		}
	}
	if req.SelectionEndAt != "" {
		if t, err := time.Parse(time.RFC3339, req.SelectionEndAt); err == nil {
			wave.SelectionEndAt = &t
		}
	}
	if req.ReregistrationDeadlineAt != "" {
		if t, err := time.Parse(time.RFC3339, req.ReregistrationDeadlineAt); err == nil {
			wave.ReregistrationDeadlineAt = &t
		}
	}

	if err := h.repo.UpdatePmbWave(wave); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate gelombang PMB").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(wave).WithContext(c))
}

// DeletePmbWave soft deletes/deactivates a PMB Wave
func (h *ReferenceHandler) DeletePmbWave(c *gin.Context) {
	id := c.Param("id")
	wave, err := h.repo.GetPmbWaveByID(id)
	if err != nil || wave == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Gelombang PMB tidak ditemukan").WithContext(c))
		return
	}

	wave.IsActive = false
	if err := h.repo.UpdatePmbWave(wave); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menonaktifkan gelombang PMB").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Gelombang PMB berhasil dinonaktifkan").WithContext(c))
}

// CreatePaymentComponent creates a payment component
func (h *ReferenceHandler) CreatePaymentComponent(c *gin.Context) {
	var req PaymentComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	comp := domain.PaymentComponent{
		Code:          req.Code,
		Name:          req.Name,
		DefaultAmount: req.DefaultAmount,
		IsActive:      true,
	}

	if err := h.repo.CreatePaymentComponent(&comp); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan komponen pembayaran").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(comp).WithContext(c))
}

// CreatePaymentMethod creates a payment method
func (h *ReferenceHandler) CreatePaymentMethod(c *gin.Context) {
	var req PaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	method := domain.PaymentMethod{
		Code:     req.Code,
		Name:     req.Name,
		Provider: req.Provider,
		IsActive: true,
	}

	if err := h.repo.CreatePaymentMethod(&method); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan metode pembayaran").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(method).WithContext(c))
}

// CreateDocumentType creates a document type
func (h *ReferenceHandler) CreateDocumentType(c *gin.Context) {
	var req DocumentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	doc := domain.DocumentType{
		Code:        req.Code,
		Name:        req.Name,
		IsMandatory: req.IsMandatory,
		IsActive:    true,
	}

	if err := h.repo.CreateDocumentType(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jenis dokumen").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(doc).WithContext(c))
}

// CreateReligion creates a religion
func (h *ReferenceHandler) CreateReligion(c *gin.Context) {
	var req ReligionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	rel := domain.Religion{
		Name: req.Name,
	}

	if err := h.repo.CreateReligion(&rel); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan agama").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(rel).WithContext(c))
}

// CreateCountry creates a country
func (h *ReferenceHandler) CreateCountry(c *gin.Context) {
	var req CountryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	country := domain.Country{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.CreateCountry(&country); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan negara").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(country).WithContext(c))
}

// ListReligions returns all religions
func (h *ReferenceHandler) ListReligions(c *gin.Context) {
	religions, err := h.repo.ListReligions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data agama").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(religions).WithContext(c))
}

// ListAdmissionPaths returns all admission paths
func (h *ReferenceHandler) ListAdmissionPaths(c *gin.Context) {
	paths, err := h.repo.ListAdmissionPaths()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data jalur masuk").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(paths).WithContext(c))
}

// ListProvinces returns all provinces
func (h *ReferenceHandler) ListProvinces(c *gin.Context) {
	provinces, err := h.repo.ListProvinces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data provinsi").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(provinces).WithContext(c))
}

// ListCities returns all cities or cities by province
func (h *ReferenceHandler) ListCities(c *gin.Context) {
	provinceID := c.Query("province_id")
	var cities []domain.City
	var err error

	if provinceID != "" {
		cities, err = h.repo.ListCitiesByProvince(provinceID)
	} else {
		cities, err = h.repo.ListCities()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kota").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(cities).WithContext(c))
}

// ListDistricts returns all districts or districts by city
func (h *ReferenceHandler) ListDistricts(c *gin.Context) {
	cityID := c.Query("city_id")
	var districts []domain.District
	var err error

	if cityID != "" {
		districts, err = h.repo.ListDistrictsByCity(cityID)
	} else {
		districts, err = h.repo.ListDistricts()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kecamatan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(districts).WithContext(c))
}

// ListVillages returns all villages or villages by district
func (h *ReferenceHandler) ListVillages(c *gin.Context) {
	districtID := c.Query("district_id")
	var villages []domain.Village
	var err error

	if districtID != "" {
		villages, err = h.repo.ListVillagesByDistrict(districtID)
	} else {
		villages, err = h.repo.ListVillages()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data desa").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(villages).WithContext(c))
}
