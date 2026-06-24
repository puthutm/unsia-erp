package domain

import (
	"time"
)

type StudyProgram struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code      string    `gorm:"column:code"`
	Name      string    `gorm:"column:name"`
	Degree    string    `gorm:"column:degree"`
	Status    string    `gorm:"column:status;default:'ACTIVE'"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (StudyProgram) TableName() string {
	return "study_programs"
}

type AcademicYear struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code      string    `gorm:"column:code"`
	Name      string    `gorm:"column:name"`
	Status    string    `gorm:"column:status;default:'INACTIVE'"`
	StartDate time.Time `gorm:"column:start_date"`
	EndDate   time.Time `gorm:"column:end_date"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (AcademicYear) TableName() string {
	return "academic_years"
}

type AcademicPeriod struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AcademicYearID string    `gorm:"column:academic_year_id"`
	Code           string    `gorm:"column:code"`
	Term           string    `gorm:"column:term"` // odd, even, intermediate
	Status         string    `gorm:"column:status;default:'INACTIVE'"`
	StartDate      time.Time `gorm:"column:start_date"`
	EndDate        time.Time `gorm:"column:end_date"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (AcademicPeriod) TableName() string {
	return "academic_periods"
}

type StatusCode struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Module      string `gorm:"column:module"`
	Code        string `gorm:"column:code"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (StatusCode) TableName() string {
	return "status_codes"
}

type PaymentComponent struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code          string    `gorm:"column:code"`
	Name          string    `gorm:"column:name"`
	DefaultAmount float64   `gorm:"column:default_amount"`
	IsActive      bool      `gorm:"column:is_active;default:true"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (PaymentComponent) TableName() string {
	return "payment_components"
}

type PaymentMethod struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code      string    `gorm:"column:code"`
	Name      string    `gorm:"column:name"`
	Provider  string    `gorm:"column:provider"`
	IsActive  bool      `gorm:"column:is_active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}

type DocumentType struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code        string    `gorm:"column:code"`
	Name        string    `gorm:"column:name"`
	IsMandatory bool      `gorm:"column:is_mandatory;default:false"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (DocumentType) TableName() string {
	return "document_types"
}

type Country struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`
}

func (Country) TableName() string {
	return "countries"
}

type Province struct {
	ID        string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	CountryID string `gorm:"column:country_id"`
	Name      string `gorm:"column:name"`
}

func (Province) TableName() string {
	return "provinces"
}

type City struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ProvinceID string `gorm:"column:province_id"`
	Name       string `gorm:"column:name"`
}

func (City) TableName() string {
	return "cities"
}

type District struct {
	ID     string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	CityID string `gorm:"column:city_id"`
	Name   string `gorm:"column:name"`
}

func (District) TableName() string {
	return "districts"
}

type Village struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	DistrictID string `gorm:"column:district_id"`
	Name       string `gorm:"column:name"`
}

func (Village) TableName() string {
	return "villages"
}

type Religion struct {
	ID   string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Name string `gorm:"column:name"`
}

func (Religion) TableName() string {
	return "religions"
}

type AdmissionPath struct {
	ID       string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code     string `gorm:"column:code"`
	Name     string `gorm:"column:name"`
	IsActive bool   `gorm:"column:is_active;default:true"`
}

func (AdmissionPath) TableName() string {
	return "admission_paths"
}

type PmbWave struct {
	ID                       string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AcademicYearID           *string    `gorm:"column:academic_year_id"`
	TargetEntryPeriodID      string     `gorm:"column:target_entry_period_id"`
	AdmissionPathID          *string    `gorm:"column:admission_path_id"`
	Code                     string     `gorm:"column:code"`
	Name                     string     `gorm:"column:name"`
	StartDate                *time.Time `gorm:"column:start_date"`
	EndDate                  *time.Time `gorm:"column:end_date"`
	RegistrationStartAt      *time.Time `gorm:"column:registration_start_at"`
	RegistrationEndAt        *time.Time `gorm:"column:registration_end_at"`
	SelectionStartAt         *time.Time `gorm:"column:selection_start_at"`
	SelectionEndAt           *time.Time `gorm:"column:selection_end_at"`
	ReregistrationDeadlineAt *time.Time `gorm:"column:reregistration_deadline_at"`
	Status                   string     `gorm:"column:status;default:'draft'"` // draft, open, closed, archived
	IsActive                 bool       `gorm:"column:is_active;default:true"`
	CreatedAt                time.Time  `gorm:"column:created_at"`
	UpdatedAt                time.Time  `gorm:"column:updated_at"`
}

func (PmbWave) TableName() string {
	return "pmb_waves"
}
