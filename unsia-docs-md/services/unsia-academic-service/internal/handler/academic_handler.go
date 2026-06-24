package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	sharedhttpclient "github.com/unsia-erp/shared-httpclient"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type StudentGenerateRequest struct {
	ApplicantID         string  `json:"applicant_id" binding:"required"`
	CurriculumID        *string `json:"curriculum_id"`
	EntryAcademicYearID *string `json:"entry_academic_year_id"`
	EntryPeriodID       *string `json:"entry_period_id"`
	StudyProgramID      string  `json:"study_program_id" binding:"required"`
	Reason              string  `json:"reason"`
}

type ClassCreateRequest struct {
	CourseOfferingID string `json:"course_offering_id" binding:"required"`
	ClassCode        string `json:"class_code" binding:"required"`
	Quota            int    `json:"quota"`
}

type KrsItemRequest struct {
	ClassID string `json:"class_id" binding:"required"`
}

type KrsCreateRequest struct {
	StudentID        string           `json:"student_id" binding:"required"`
	AcademicPeriodID string           `json:"academic_period_id" binding:"required"`
	Items            []KrsItemRequest `json:"items" binding:"required,gt=0"`
}

type GradeImportRequest struct {
	SourceModule    string  `json:"source_module" binding:"required"`
	AcademicClassID string  `json:"academic_class_id" binding:"required"`
	StudentID       string  `json:"student_id" binding:"required"`
	Score           float64 `json:"score" binding:"required"`
}

type GradeFinalizeRequest struct {
	NumericGrade float64 `json:"numeric_grade" binding:"required"`
	LetterGrade  string  `json:"letter_grade" binding:"required"`
	GradePoint   float64 `json:"grade_point" binding:"required"`
}

type CurriculumCreateRequest struct {
	StudyProgramID         string  `json:"study_program_id" binding:"required"`
	Code                   string  `json:"code" binding:"required"`
	Name                   string  `json:"name" binding:"required"`
	CurriculumYear         int     `json:"curriculum_year" binding:"required"`
	EffectiveStartPeriodID *string `json:"effective_start_period_id"`
	EffectiveEndPeriodID   *string `json:"effective_end_period_id"`
	IsActive               *bool   `json:"is_active"`
	IsDefaultForNewStudent *bool   `json:"is_default_for_new_student"`
}

type CourseCreateRequest struct {
	StudyProgramID *string  `json:"study_program_id"`
	CourseCode     string   `json:"course_code" binding:"required"`
	CourseName     string   `json:"course_name" binding:"required"`
	Sks            int      `json:"sks"`
	CourseType     string   `json:"course_type"`
	MinimumGrade   *float64 `json:"minimum_grade"`
	IsActive       *bool    `json:"is_active"`
}

type CurriculumCourseCreateRequest struct {
	CurriculumID string `json:"curriculum_id" binding:"required"`
	CourseID     string `json:"course_id" binding:"required"`
	Semester     int    `json:"semester" binding:"required"`
	IsMandatory  *bool  `json:"is_mandatory"`
}

type CourseOfferingCreateRequest struct {
	AcademicPeriodID string `json:"academic_period_id" binding:"required"`
	CourseID         string `json:"course_id" binding:"required"`
	IsActive         *bool  `json:"is_active"`
}

type ClassLecturerPlotRequest struct {
	LecturerID string `json:"lecturer_id" binding:"required"`
	RoleType   string `json:"role_type"`
}

type GradeCorrectionRequest struct {
	NumericGrade float64 `json:"numeric_grade" binding:"required"`
	LetterGrade  string  `json:"letter_grade" binding:"required"`
	GradePoint   float64 `json:"grade_point" binding:"required"`
	Reason       string  `json:"reason" binding:"required"`
}

type AcademicHandler struct {
	repo          *repository.AcademicRepository
	db            *gorm.DB
	financeClient *sharedhttpclient.Client
}

func NewAcademicHandler(db *gorm.DB) *AcademicHandler {
	financeURL := os.Getenv("FINANCE_SERVICE_URL")
	if financeURL == "" {
		financeURL = "http://localhost:8005"
	}
	srvToken := os.Getenv("ACADEMIC_SERVICE_TOKEN")
	if srvToken == "" {
		srvToken = "academic_service_secret_token"
	}

	financeClient := sharedhttpclient.New(sharedhttpclient.Config{
		BaseURL:      financeURL,
		ServiceToken: srvToken,
		SourceName:   "academic-service",
		Timeout:      10 * time.Second,
	})

	return &AcademicHandler{
		repo:          repository.NewAcademicRepository(db),
		db:            db,
		financeClient: financeClient,
	}
}

func (h *AcademicHandler) GenerateStudentFromApplicant(c *gin.Context) {
	var req StudentGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	correlationID, _ := c.Get("x-correlation-id")
	cid, _ := correlationID.(string)

	// Verify if student already exists for this applicant
	existing, err := h.repo.GetStudentByApplicantID(req.ApplicantID)
	if err == nil && existing != nil {
		c.JSON(http.StatusOK, sharederr.Success(gin.H{
			"student_id":      existing.ID,
			"nim":             existing.Nim,
			"handover_status": "success",
		}).WithContext(c))
		return
	}

	var student domain.Student
	err = h.db.Transaction(func(tx *gorm.DB) error {
		yearStr := time.Now().Format("2006")
		periodID := "default-period-id"
		if req.EntryPeriodID != nil {
			periodID = *req.EntryPeriodID
		}

		nim, err := h.repo.GenerateNIM(tx, req.StudyProgramID, periodID, yearStr)
		if err != nil {
			return err
		}

		// Person ID matches applicant's person ID. In a real scenario we could retrieve it.
		// For testing simplicity we will mock person ID from applicant ID or generate uuid
		personID := "person-generated-id"

		student = domain.Student{
			PersonID:            personID,
			ApplicantID:         &req.ApplicantID,
			StudyProgramID:      req.StudyProgramID,
			Nim:                 nim,
			StudentStatus:       "active",
			EntryAcademicYearID: req.EntryAcademicYearID,
			EntryPeriodID:       req.EntryPeriodID,
			CurriculumID:        req.CurriculumID,
			CurrentSemester:     1,
		}

		if err := tx.Create(&student).Error; err != nil {
			return err
		}

		envelope := sharedevent.EventEnvelope{
			EventName:        "academic.student_created",
			EventVersion:     "v1",
			PublisherService: "academic-service",
			AggregateType:    "student",
			AggregateID:      student.ID,
			CorrelationID:    cid,
			Payload: map[string]interface{}{
				"student_id":       student.ID,
				"person_id":        student.PersonID,
				"nim":              student.Nim,
				"study_program_id": student.StudyProgramID,
				"status":           student.StudentStatus,
			},
		}

		conn := tx.Statement.ConnPool
		_, err = sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal melakukan generate mahasiswa").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.student.create",
		Module:       "academic",
		ResourceType: "student",
		ResourceID:   student.ID,
		NewValue:     student,
	})

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"student_id":      student.ID,
		"nim":             student.Nim,
		"handover_status": "success",
	}).WithContext(c))
}

func (h *AcademicHandler) ListStudents(c *gin.Context) {
	studyProgramID := c.Query("study_program_id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	students, total, err := h.repo.ListStudents(studyProgramID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data mahasiswa").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(students).WithPagination(page, limit, int(total)).WithContext(c))
}

func (h *AcademicHandler) CreateClass(c *gin.Context) {
	var req ClassCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	quota := req.Quota
	if quota <= 0 {
		quota = 40
	}

	class := domain.Class{
		CourseOfferingID: req.CourseOfferingID,
		ClassCode:        req.ClassCode,
		Quota:            quota,
		EnrolledCount:    0,
		ClassStatus:      "active",
		IsParallel:       false,
	}

	if err := h.repo.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat kelas perkuliahan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(class).WithContext(c))
}

func (h *AcademicHandler) CreateKrsDraft(c *gin.Context) {
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

func (h *AcademicHandler) SubmitKrs(c *gin.Context) {
	krsID := c.Param("krs_id")

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

	ctx := context.WithValue(c.Request.Context(), "x-correlation-id", cid)

	// Call Finance Service to check clearance status
	clearanceURL := fmt.Sprintf("/api/v1/finance/clearances?student_id=%s&academic_period_id=%s", krs.StudentID, krs.AcademicPeriodID)
	resp, err := h.financeClient.Get(ctx, clearanceURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, sharederr.Error("FINANCE_SERVICE_UNAVAILABLE", "Layanan Finance tidak dapat diakses untuk pengecekan clearance").WithContext(c))
		return
	}
	defer resp.Body.Close()

	var clResp struct {
		Success bool `json:"success"`
		Data    struct {
			ClearanceStatus string   `json:"clearance_status"`
			BlockReasons    []string `json:"block_reasons"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&clResp); err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, sharederr.Error("FINANCE_SERVICE_ERROR", "Gagal mengambil data clearance dari Finance").WithContext(c))
		return
	}

	if clResp.Data.ClearanceStatus == "blocked" {
		reasons := "Mahasiswa diblokir secara finansial."
		if len(clResp.Data.BlockReasons) > 0 {
			reasons = clResp.Data.BlockReasons[0]
		}
		c.JSON(http.StatusForbidden, sharederr.Error("FINANCE_BLOCKED", fmt.Sprintf("Submit KRS ditolak: %s", reasons)).WithContext(c))
		return
	}

	if err := h.repo.UpdateKRSStatus(krsID, "submitted"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate status KRS").WithContext(c))
		return
	}

	krs.Status = "submitted"
	c.JSON(http.StatusOK, sharederr.Success(krs).WithContext(c))
}

func (h *AcademicHandler) ApproveKrs(c *gin.Context) {
	krsID := c.Param("krs_id")

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
			// Ignore conflict if exists
			tx.Create(&grade)
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
				"classes":            classIDs,
			},
		}

		conn := tx.Statement.ConnPool
		_, err = sharedevent.WriteOutbox(c.Request.Context(), conn, envelope, "INTEGRATION_EVENT")
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses persetujuan KRS").WithContext(c))
		return
	}

	krs.Status = "approved"
	c.JSON(http.StatusOK, sharederr.Success(krs).WithContext(c))
}

func (h *AcademicHandler) ImportGradeSource(c *gin.Context) {
	var req GradeImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Retrieve KRS Item matching student_id and academic_class_id
	var item domain.KrsItem
	err := h.db.Joins("JOIN krs ON krs.id = krs_items.krs_id").
		Where("krs.student_id = ? AND krs_items.class_id = ? AND krs.status = 'approved'", req.StudentID, req.AcademicClassID).
		First(&item).Error

	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kelas yang diikuti mahasiswa tersebut tidak ditemukan").WithContext(c))
		return
	}

	var grade domain.Grade
	err = h.db.Where("krs_item_id = ?", item.ID).First(&grade).Error
	if err != nil {
		// create a new one
		grade = domain.Grade{
			KrsItemID:    item.ID,
			NumericGrade: &req.Score,
			Source:       req.SourceModule,
		}
		h.db.Create(&grade)
	} else {
		h.db.Model(&grade).Updates(map[string]interface{}{
			"numeric_grade": &req.Score,
			"source":        req.SourceModule,
			"submitted_at":  time.Now(),
		})
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(grade, "Komponen nilai berhasil diimport").WithContext(c))
}

func (h *AcademicHandler) FinalizeGrade(c *gin.Context) {
	gradeID := c.Param("grade_id")
	var req GradeFinalizeRequest
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

	var grade domain.Grade
	err := h.db.Where("id = ?", gradeID).First(&grade).Error
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Data nilai tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.UpdateGrade(gradeID, req.NumericGrade, req.LetterGrade, req.GradePoint, "manual", actor); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memfinalisasi nilai").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Nilai akademik berhasil difinalisasi").WithContext(c))
}

type TranscriptResponse struct {
	StudentID string                   `json:"student_id"`
	IPK       float64                  `json:"ipk"`
	IPS       *float64                 `json:"ips,omitempty"`
	TotalSKS  int                      `json:"total_sks"`
	Grades    []StudentGradeDetailItem `json:"grades"`
}

type StudentGradeDetailItem struct {
	GradeID          string   `json:"grade_id"`
	ClassID          string   `json:"class_id"`
	CourseID         string   `json:"course_id"`
	CourseCode       string   `json:"course_code"`
	CourseName       string   `json:"course_name"`
	Sks              int      `json:"sks"`
	AcademicPeriodID string   `json:"academic_period_id"`
	NumericGrade     *float64 `json:"numeric_grade"`
	LetterGrade      string   `json:"letter_grade"`
	GradePoint       *float64 `json:"grade_point"`
}

func (h *AcademicHandler) CreateCurriculum(c *gin.Context) {
	var req CurriculumCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	isDefault := false
	if req.IsDefaultForNewStudent != nil {
		isDefault = *req.IsDefaultForNewStudent
	}

	curriculum := domain.Curriculum{
		StudyProgramID:         req.StudyProgramID,
		Code:                   req.Code,
		Name:                   req.Name,
		CurriculumYear:         req.CurriculumYear,
		Status:                 "draft",
		EffectiveStartPeriodID: req.EffectiveStartPeriodID,
		EffectiveEndPeriodID:   req.EffectiveEndPeriodID,
		IsActive:               isActive,
		IsDefaultForNewStudent: isDefault,
	}

	if err := h.repo.CreateCurriculum(&curriculum); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat kurikulum").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.curriculum.create",
		Module:       "academic",
		ResourceType: "curriculum",
		ResourceID:   curriculum.ID,
		NewValue:     curriculum,
	})

	c.JSON(http.StatusCreated, sharederr.Success(curriculum).WithContext(c))
}

func (h *AcademicHandler) CreateCourse(c *gin.Context) {
	var req CourseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	sks := 2
	if req.Sks > 0 {
		sks = req.Sks
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	course := domain.Course{
		StudyProgramID: req.StudyProgramID,
		CourseCode:     req.CourseCode,
		CourseName:     req.CourseName,
		Sks:            sks,
		CourseType:     req.CourseType,
		MinimumGrade:   req.MinimumGrade,
		IsActive:       isActive,
	}

	if err := h.repo.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat mata kuliah").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course.create",
		Module:       "academic",
		ResourceType: "course",
		ResourceID:   course.ID,
		NewValue:     course,
	})

	c.JSON(http.StatusCreated, sharederr.Success(course).WithContext(c))
}

func (h *AcademicHandler) CreateCurriculumCourse(c *gin.Context) {
	var req CurriculumCourseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	isMandatory := true
	if req.IsMandatory != nil {
		isMandatory = *req.IsMandatory
	}

	cc := domain.CurriculumCourse{
		CurriculumID: req.CurriculumID,
		CourseID:     req.CourseID,
		Semester:     req.Semester,
		IsMandatory:  isMandatory,
	}

	if err := h.repo.CreateCurriculumCourse(&cc); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memetakan mata kuliah ke kurikulum").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.curriculum_course.create",
		Module:       "academic",
		ResourceType: "curriculum_course",
		ResourceID:   cc.ID,
		NewValue:     cc,
	})

	c.JSON(http.StatusCreated, sharederr.Success(cc).WithContext(c))
}

func (h *AcademicHandler) CreateCourseOffering(c *gin.Context) {
	var req CourseOfferingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	co := domain.CourseOffering{
		AcademicPeriodID: req.AcademicPeriodID,
		CourseID:         req.CourseID,
		IsActive:         isActive,
	}

	if err := h.repo.CreateCourseOffering(&co); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat penawaran mata kuliah").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.course_offering.create",
		Module:       "academic",
		ResourceType: "course_offering",
		ResourceID:   co.ID,
		NewValue:     co,
	})

	c.JSON(http.StatusCreated, sharederr.Success(co).WithContext(c))
}

func (h *AcademicHandler) PlotClassLecturer(c *gin.Context) {
	classID := c.Param("id")
	var req ClassLecturerPlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	roleType := req.RoleType
	if roleType == "" {
		roleType = "teacher"
	}

	cl := domain.ClassLecturer{
		ClassID:    classID,
		LecturerID: req.LecturerID,
		RoleType:   roleType,
	}

	if err := h.repo.CreateClassLecturer(&cl); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memplot dosen ke kelas").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.class_lecturer.plot",
		Module:       "academic",
		ResourceType: "class_lecturer",
		ResourceID:   cl.ID,
		NewValue:     cl,
	})

	c.JSON(http.StatusCreated, sharederr.Success(cl).WithContext(c))
}

func (h *AcademicHandler) CorrectGrade(c *gin.Context) {
	gradeID := c.Param("id")
	var req GradeCorrectionRequest
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

	// Fetch current grade
	grade, err := h.repo.GetGradeByID(gradeID)
	if err != nil || grade == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Data nilai tidak ditemukan").WithContext(c))
		return
	}

	// Serialize old value
	oldValBytes, _ := json.Marshal(grade)
	oldValStr := string(oldValBytes)

	// Build new grade object for comparison/serialization
	updatedGrade := *grade
	updatedGrade.NumericGrade = &req.NumericGrade
	updatedGrade.LetterGrade = req.LetterGrade
	updatedGrade.GradePoint = &req.GradePoint
	updatedGrade.Source = "manual"
	now := time.Now()
	updatedGrade.SubmittedAt = &now
	updatedGrade.SubmittedBy = &actor

	newValBytes, _ := json.Marshal(updatedGrade)
	newValStr := string(newValBytes)

	// Perform transaction
	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Update grade
		if err := tx.Model(grade).Updates(map[string]interface{}{
			"numeric_grade": &req.NumericGrade,
			"letter_grade":  req.LetterGrade,
			"grade_point":   &req.GradePoint,
			"source":        "manual",
			"submitted_at":  &now,
			"submitted_by":  &actor,
		}).Error; err != nil {
			return err
		}

		// Create Grade History
		gh := domain.GradeHistory{
			GradeID:   gradeID,
			OldValue:  oldValStr,
			NewValue:  newValStr,
			ChangedBy: &actor,
			Reason:    req.Reason,
			ChangedAt: time.Now(),
		}
		if err := tx.Create(&gh).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses koreksi nilai").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.grade.correct",
		Module:       "academic",
		ResourceType: "grade",
		ResourceID:   gradeID,
		NewValue:     updatedGrade,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(updatedGrade, "Nilai berhasil dikoreksi").WithContext(c))
}

func (h *AcademicHandler) GetStudentTranscript(c *gin.Context) {
	studentID := c.Param("id")
	academicPeriodID := c.Query("academic_period_id")

	var items []StudentGradeDetailItem
	query := h.db.Table("grades").
		Select("grades.id as grade_id, krs_items.class_id, course_offerings.course_id, courses.course_code, courses.course_name, courses.sks, krs.academic_period_id, grades.numeric_grade, grades.letter_grade, grades.grade_point").
		Joins("JOIN krs_items ON krs_items.id = grades.krs_item_id").
		Joins("JOIN krs ON krs.id = krs_items.krs_id").
		Joins("JOIN classes ON classes.id = krs_items.class_id").
		Joins("JOIN course_offerings ON course_offerings.id = classes.course_offering_id").
		Joins("JOIN courses ON courses.id = course_offerings.course_id").
		Where("krs.student_id = ? AND krs.status = 'approved'", studentID)

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil riwayat nilai mahasiswa").WithContext(c))
		return
	}

	// Calculate Cumulative IPK & Total SKS using highest grade for each course
	bestGrades := make(map[string]StudentGradeDetailItem)
	for _, item := range items {
		if item.GradePoint == nil {
			continue
		}
		existing, exists := bestGrades[item.CourseID]
		if !exists || *item.GradePoint > *existing.GradePoint {
			bestGrades[item.CourseID] = item
		}
	}

	var totalPoints float64
	var totalSks int
	for _, item := range bestGrades {
		if item.GradePoint != nil {
			totalPoints += (*item.GradePoint) * float64(item.Sks)
			totalSks += item.Sks
		}
	}

	var ipk float64
	if totalSks > 0 {
		ipk = totalPoints / float64(totalSks)
	}

	var response TranscriptResponse
	response.StudentID = studentID
	response.IPK = ipk
	response.TotalSKS = totalSks
	response.Grades = items

	// If academic_period_id is supplied, calculate Semester IPS
	if academicPeriodID != "" {
		var totalPeriodPoints float64
		var totalPeriodSks int
		for _, item := range items {
			if item.AcademicPeriodID == academicPeriodID && item.GradePoint != nil {
				totalPeriodPoints += (*item.GradePoint) * float64(item.Sks)
				totalPeriodSks += item.Sks
			}
		}

		var ips float64
		if totalPeriodSks > 0 {
			ips = totalPeriodPoints / float64(totalPeriodSks)
		}
		response.IPS = &ips
	}

	c.JSON(http.StatusOK, sharederr.Success(response).WithContext(c))
}
