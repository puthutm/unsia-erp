package repository

import (
	"errors"
	"time"

	"github.com/unsia-erp/unsia-pmb-service/internal/domain"
	"gorm.io/gorm"
)

type PMBRepository struct {
	db *gorm.DB
}

func NewPMBRepository(db *gorm.DB) *PMBRepository {
	return &PMBRepository{db: db}
}

// Applicant Operations
func (r *PMBRepository) CreateApplicant(app *domain.Applicant) error {
	return r.db.Create(app).Error
}

func (r *PMBRepository) GetApplicantByID(id string) (*domain.Applicant, error) {
	var app domain.Applicant
	err := r.db.Preload("Biodata").
		Preload("Addresses").
		Preload("EducationBackgrounds").
		Preload("FamilyMembers").
		Preload("FinancialProfiles").
		Preload("FacilityProfiles").
		Preload("Documents").
		Where("id = ?", id).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (r *PMBRepository) GetApplicantByPersonID(personID string) (*domain.Applicant, error) {
	var app domain.Applicant
	err := r.db.Where("person_id = ?", personID).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

func (r *PMBRepository) UpdateApplicantStatus(id string, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}
	if status == "submitted" {
		updates["submitted_at"] = &now
	} else if status == "accepted" {
		updates["accepted_at"] = &now
	}
	return r.db.Model(&domain.Applicant{}).Where("id = ?", id).Updates(updates).Error
}

// Document Operations
func (r *PMBRepository) CreateApplicantDocument(doc *domain.ApplicantDocument) error {
	return r.db.Create(doc).Error
}

func (r *PMBRepository) GetApplicantDocumentByID(id string) (*domain.ApplicantDocument, error) {
	var doc domain.ApplicantDocument
	err := r.db.Where("id = ?", id).First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *PMBRepository) UpdateApplicantDocumentVerification(id string, status string, verifiedBy string, rejectReason string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"verification_status": status,
		"verified_by":         &verifiedBy,
		"verified_at":         &now,
		"reject_reason":       rejectReason,
	}
	return r.db.Model(&domain.ApplicantDocument{}).Where("id = ?", id).Updates(updates).Error
}

// Loa Operations
func (r *PMBRepository) CreateLoaDocument(loa *domain.LoaDocument) error {
	return r.db.Create(loa).Error
}

// Handover Operations
func (r *PMBRepository) CreateHandoverLog(log *domain.HandoverLog) error {
	return r.db.Create(log).Error
}

func (r *PMBRepository) GetHandoverLogByIdempotencyKey(key string) (*domain.HandoverLog, error) {
	var log domain.HandoverLog
	err := r.db.Where("idempotency_key = ?", key).First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *PMBRepository) UpdateHandoverLog(id string, status string, studentRefID *string, nim *string, errMessage string) error {
	updates := map[string]interface{}{
		"status":        status,
		"error_message": errMessage,
	}
	if studentRefID != nil {
		updates["student_ref_id"] = studentRefID
	}
	if nim != nil {
		updates["nim"] = nim
	}
	return r.db.Model(&domain.HandoverLog{}).Where("id = ?", id).Updates(updates).Error
}

// Applicant List with Filtering
type ApplicantListFilter struct {
	Status          string
	StudyProgramID  string
	PmbWaveID       string
	AdmissionPathID string
	Search          string
	Page            int
	Limit           int
}

func (r *PMBRepository) GetApplicants(filter ApplicantListFilter) ([]domain.Applicant, int64, error) {
	query := r.db.Model(&domain.Applicant{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.StudyProgramID != "" {
		query = query.Where("study_program_id = ?", filter.StudyProgramID)
	}
	if filter.PmbWaveID != "" {
		query = query.Where("pmb_wave_id = ?", filter.PmbWaveID)
	}
	if filter.AdmissionPathID != "" {
		query = query.Where("admission_path_id = ?", filter.AdmissionPathID)
	}
	if filter.Search != "" {
		query = query.Where("registration_number ILIKE ?", "%"+filter.Search+"%")
	}

	var total int64
	query.Count(&total)

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var applicants []domain.Applicant
	err := query.Preload("Biodata").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&applicants).Error

	return applicants, total, err
}

// Biodata Operations
func (r *PMBRepository) GetBiodataByApplicantID(applicantID string) (*domain.ApplicantBiodata, error) {
	var biodata domain.ApplicantBiodata
	err := r.db.Where("applicant_id = ?", applicantID).First(&biodata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &biodata, nil
}

func (r *PMBRepository) UpdateBiodata(applicantID string, updates map[string]interface{}) error {
	return r.db.Model(&domain.ApplicantBiodata{}).Where("applicant_id = ?", applicantID).Updates(updates).Error
}

// Address Operations
func (r *PMBRepository) GetAddressesByApplicantID(applicantID string) ([]domain.ApplicantAddress, error) {
	var addresses []domain.ApplicantAddress
	err := r.db.Where("applicant_id = ?", applicantID).Find(&addresses).Error
	return addresses, err
}

func (r *PMBRepository) UpsertAddress(addr *domain.ApplicantAddress) error {
	return r.db.Where("applicant_id = ? AND address_type = ?", addr.ApplicantID, addr.AddressType).
		Assign(*addr).
		FirstOrCreate(addr).Error
}

// Education Operations
func (r *PMBRepository) GetEducationBackgroundsByApplicantID(applicantID string) ([]domain.ApplicantEducationBackground, error) {
	var educations []domain.ApplicantEducationBackground
	err := r.db.Where("applicant_id = ?", applicantID).Find(&educations).Error
	return educations, err
}

func (r *PMBRepository) UpsertEducationBackground(edu *domain.ApplicantEducationBackground) error {
	return r.db.Where("applicant_id = ?", edu.ApplicantID).
		Assign(*edu).
		FirstOrCreate(edu).Error
}

// Family Operations
func (r *PMBRepository) GetFamilyMembersByApplicantID(applicantID string) ([]domain.ApplicantFamilyMember, error) {
	var members []domain.ApplicantFamilyMember
	err := r.db.Where("applicant_id = ?", applicantID).Find(&members).Error
	return members, err
}

func (r *PMBRepository) UpsertFamilyMember(member *domain.ApplicantFamilyMember) error {
	return r.db.Where("applicant_id = ? AND relationship = ?", member.ApplicantID, member.Relationship).
		Assign(*member).
		FirstOrCreate(member).Error
}

// Financial Operations
func (r *PMBRepository) GetFinancialProfilesByApplicantID(applicantID string) ([]domain.ApplicantFinancialProfile, error) {
	var profiles []domain.ApplicantFinancialProfile
	err := r.db.Where("applicant_id = ?", applicantID).Find(&profiles).Error
	return profiles, err
}

func (r *PMBRepository) UpsertFinancialProfile(profile *domain.ApplicantFinancialProfile) error {
	return r.db.Where("applicant_id = ? AND sponsor_type = ?", profile.ApplicantID, profile.SponsorType).
		Assign(*profile).
		FirstOrCreate(profile).Error
}

// Facility Operations
func (r *PMBRepository) GetFacilityProfilesByApplicantID(applicantID string) ([]domain.ApplicantFacilityProfile, error) {
	var profiles []domain.ApplicantFacilityProfile
	err := r.db.Where("applicant_id = ?", applicantID).Find(&profiles).Error
	return profiles, err
}

func (r *PMBRepository) UpsertFacilityProfile(profile *domain.ApplicantFacilityProfile) error {
	return r.db.Where("applicant_id = ? AND facility_type = ?", profile.ApplicantID, profile.FacilityType).
		Assign(*profile).
		FirstOrCreate(profile).Error
}

// Document Operations
func (r *PMBRepository) GetDocumentsByApplicantID(applicantID string) ([]domain.ApplicantDocument, error) {
	var documents []domain.ApplicantDocument
	err := r.db.Where("applicant_id = ?", applicantID).Find(&documents).Error
	return documents, err
}
