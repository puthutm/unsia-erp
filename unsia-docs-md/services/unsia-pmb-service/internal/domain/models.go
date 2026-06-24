package domain

import (
	"time"
)

type Applicant struct {
	ID                   string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PersonID             string     `gorm:"column:person_id;not null"` // external_ref: core.persons.id
	UserID               *string    `gorm:"column:user_id"`           // external_ref: core.users.id
	CrmLeadID            *string    `gorm:"column:crm_lead_id"`       // external_ref: crm.leads.id
	StudyProgramID       *string    `gorm:"column:study_program_id"`  // external_ref: ref.study_programs.id
	PmbWaveID            *string    `gorm:"column:pmb_wave_id"`       // external_ref: ref.pmb_waves.id
	AdmissionPathID      *string    `gorm:"column:admission_path_id"` // external_ref: ref.admission_paths.id
	TargetEntryPeriodID  *string    `gorm:"column:target_entry_period_id"` // external_ref: ref.academic_periods.id
	RegistrationNumber   string     `gorm:"column:registration_number;unique;not null"`
	Status               string     `gorm:"column:status;default:'draft';not null"`
	SubmittedAt          *time.Time `gorm:"column:submitted_at"`
	AcceptedAt           *time.Time `gorm:"column:accepted_at"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at"`

	Biodata              *ApplicantBiodata            `gorm:"foreignKey:ApplicantID;references:ID"`
	Addresses            []ApplicantAddress           `gorm:"foreignKey:ApplicantID;references:ID"`
	EducationBackgrounds []ApplicantEducationBackground `gorm:"foreignKey:ApplicantID;references:ID"`
	FamilyMembers        []ApplicantFamilyMember      `gorm:"foreignKey:ApplicantID;references:ID"`
	FinancialProfiles    []ApplicantFinancialProfile  `gorm:"foreignKey:ApplicantID;references:ID"`
	FacilityProfiles     []ApplicantFacilityProfile   `gorm:"foreignKey:ApplicantID;references:ID"`
	Documents            []ApplicantDocument          `gorm:"foreignKey:ApplicantID;references:ID"`
}

func (Applicant) TableName() string {
	return "applicants"
}

type ApplicantBiodata struct {
	ID             string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID    string     `gorm:"column:applicant_id;not null;unique"`
	FullName       string     `gorm:"column:full_name"`
	Email          string     `gorm:"column:email"`
	Phone          string     `gorm:"column:phone"`
	Nik            string     `gorm:"column:nik"` // For camaba (calon mahasiswa baru)
	BirthPlace     string     `gorm:"column:birth_place"`
	BirthDate      *time.Time `gorm:"column:birth_date"`
	Gender         string     `gorm:"column:gender"`
	ReligionID     *string    `gorm:"column:religion_id"` // external_ref: ref.religions.id
	MaritalStatus  string     `gorm:"column:marital_status"`
	Citizenship    string     `gorm:"column:citizenship"`
	JacketSize     string     `gorm:"column:jacket_size"`
	CoreSyncStatus string     `gorm:"column:core_sync_status;default:'pending'"`
	CoreSyncedAt   *time.Time `gorm:"column:core_synced_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
}

func (ApplicantBiodata) TableName() string {
	return "applicant_biodata"
}

type ApplicantAddress struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID  string `gorm:"column:applicant_id;not null"`
	AddressType  string `gorm:"column:address_type"`
	Street       string `gorm:"column:street"`
	ProvinceID   *string `gorm:"column:province_id"` // external_ref: ref.provinces.id
	CityID       *string `gorm:"column:city_id"`     // external_ref: ref.cities.id
	DistrictID   *string `gorm:"column:district_id"` // external_ref: ref.districts.id
	VillageID    *string `gorm:"column:village_id"`  // external_ref: ref.villages.id
	PostalCode   string `gorm:"column:postal_code"`
	IsSameAsKtp  bool   `gorm:"column:is_same_as_ktp;default:false;not null"`
}

func (ApplicantAddress) TableName() string {
	return "applicant_addresses"
}

type ApplicantEducationBackground struct {
	ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID    string `gorm:"column:applicant_id;not null"`
	SchoolName     string `gorm:"column:school_name"`
	Major          string `gorm:"column:major"`
	GraduationYear string `gorm:"column:graduation_year"`
	Gpa            float64 `gorm:"column:gpa"`
}

func (ApplicantEducationBackground) TableName() string {
	return "applicant_education_backgrounds"
}

type ApplicantFamilyMember struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID  string `gorm:"column:applicant_id;not null"`
	Relationship string `gorm:"column:relationship"`
	FullName     string `gorm:"column:full_name"`
	Occupation   string `gorm:"column:occupation"`
	Income       float64 `gorm:"column:income"`
}

func (ApplicantFamilyMember) TableName() string {
	return "applicant_family_members"
}

type ApplicantFinancialProfile struct {
	ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID    string `gorm:"column:applicant_id;not null"`
	SponsorType    string `gorm:"column:sponsor_type"`
	SponsorName    string `gorm:"column:sponsor_name"`
	MonthlyIncome  float64 `gorm:"column:monthly_income"`
}

func (ApplicantFinancialProfile) TableName() string {
	return "applicant_financial_profiles"
}

type ApplicantFacilityProfile struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID  string `gorm:"column:applicant_id;not null"`
	FacilityType string `gorm:"column:facility_type"`
	Description  string `gorm:"column:description"`
}

func (ApplicantFacilityProfile) TableName() string {
	return "applicant_facility_profiles"
}

type ApplicantDocument struct {
	ID                 string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID        string     `gorm:"column:applicant_id;not null"`
	DocumentTypeCode   string     `gorm:"column:document_type_code;not null"` // external_ref: ref.document_types.code
	FileUrl            string     `gorm:"column:file_url;not null"`
	VerificationStatus string     `gorm:"column:verification_status;default:'pending';not null"` // pending, verified, rejected
	VerifiedBy         *string    `gorm:"column:verified_by"` // external_ref: core.users.id
	VerifiedAt         *time.Time `gorm:"column:verified_at"`
	RejectReason       string     `gorm:"column:reject_reason"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
}

func (ApplicantDocument) TableName() string {
	return "applicant_documents"
}

type ApplicantStatusHistory struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID string    `gorm:"column:applicant_id;not null"`
	OldStatus *string   `gorm:"column:old_status"`
	NewStatus string    `gorm:"column:new_status;not null"`
	ChangedBy *string   `gorm:"column:changed_by"` // external_ref: core.users.id
	Note      string    `gorm:"column:note"`
	ChangedAt time.Time `gorm:"column:changed_at;default:now()"`
}

func (ApplicantStatusHistory) TableName() string {
	return "applicant_status_histories"
}

type ReRegistration struct {
	ID                 string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID        string     `gorm:"column:applicant_id;not null;unique"`
	ReRegistrationDate time.Time  `gorm:"column:re_registration_date;default:now()"`
	Status             string     `gorm:"column:status;default:'pending';not null"` // pending, completed
	VerifiedBy         *string    `gorm:"column:verified_by"` // external_ref: core.users.id
	VerifiedAt         *time.Time `gorm:"column:verified_at"`
}

func (ReRegistration) TableName() string {
	return "re_registrations"
}

type LoaDocument struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID string    `gorm:"column:applicant_id;not null"`
	LoaNumber   string    `gorm:"column:loa_number;unique;not null"`
	FileUrl     string    `gorm:"column:file_url;not null"`
	IssuedBy    string    `gorm:"column:issued_by;not null"` // external_ref: core.users.id
	IssuedAt    time.Time `gorm:"column:issued_at;default:now()"`
}

func (LoaDocument) TableName() string {
	return "loa_documents"
}

type HandoverLog struct {
	ID             string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicantID    string     `gorm:"column:applicant_id;not null"`
	HandoverDate   time.Time  `gorm:"column:handover_date;default:now()"`
	StudentRefID   *string    `gorm:"column:student_ref_id"` // external_ref: academic.students.id
	Nim            *string    `gorm:"column:nim"`
	Status         string     `gorm:"column:status;default:'pending';not null"` // pending, success, failed
	ErrorMessage   string     `gorm:"column:error_message"`
	IdempotencyKey string     `gorm:"column:idempotency_key;unique;not null"`
}

func (HandoverLog) TableName() string {
	return "handover_logs"
}
