package repository

import (
	"github.com/unsia-erp/unsia-reference-service/internal/domain"
	"gorm.io/gorm"
)

type ReferenceRepository struct {
	db *gorm.DB
}

func NewReferenceRepository(db *gorm.DB) *ReferenceRepository {
	return &ReferenceRepository{db: db}
}

// Study Program Operations
func (r *ReferenceRepository) ListStudyPrograms() ([]domain.StudyProgram, error) {
	var prodis []domain.StudyProgram
	err := r.db.Order("code asc").Find(&prodis).Error
	return prodis, err
}

func (r *ReferenceRepository) CreateStudyProgram(p *domain.StudyProgram) error {
	return r.db.Create(p).Error
}

func (r *ReferenceRepository) UpdateStudyProgramStatus(id string, status string) error {
	return r.db.Model(&domain.StudyProgram{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ReferenceRepository) GetStudyProgramByID(id string) (*domain.StudyProgram, error) {
	var p domain.StudyProgram
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Academic Year Operations
func (r *ReferenceRepository) ListAcademicYears() ([]domain.AcademicYear, error) {
	var years []domain.AcademicYear
	err := r.db.Order("code desc").Find(&years).Error
	return years, err
}

func (r *ReferenceRepository) CreateAcademicYear(y *domain.AcademicYear) error {
	return r.db.Create(y).Error
}

func (r *ReferenceRepository) UpdateAcademicYearStatus(id string, status string) error {
	return r.db.Model(&domain.AcademicYear{}).Where("id = ?", id).Update("status", status).Error
}

// Academic Period Operations
func (r *ReferenceRepository) ListAcademicPeriods() ([]domain.AcademicPeriod, error) {
	var periods []domain.AcademicPeriod
	err := r.db.Order("code desc").Find(&periods).Error
	return periods, err
}

func (r *ReferenceRepository) CreateAcademicPeriod(p *domain.AcademicPeriod) error {
	return r.db.Create(p).Error
}

func (r *ReferenceRepository) UpdateAcademicPeriodStatus(id string, status string) error {
	return r.db.Model(&domain.AcademicPeriod{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ReferenceRepository) GetAcademicPeriodByID(id string) (*domain.AcademicPeriod, error) {
	var p domain.AcademicPeriod
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Status Codes Operations
func (r *ReferenceRepository) ListStatusCodes() ([]domain.StatusCode, error) {
	var codes []domain.StatusCode
	err := r.db.Order("module asc, code asc").Find(&codes).Error
	return codes, err
}

// Payment Components Operations
func (r *ReferenceRepository) ListPaymentComponents() ([]domain.PaymentComponent, error) {
	var comps []domain.PaymentComponent
	err := r.db.Where("is_active = true").Order("code asc").Find(&comps).Error
	return comps, err
}

// Document Types Operations
func (r *ReferenceRepository) ListDocumentTypes() ([]domain.DocumentType, error) {
	var docs []domain.DocumentType
	err := r.db.Where("is_active = true").Order("code asc").Find(&docs).Error
	return docs, err
}

// Payment Methods Operations
func (r *ReferenceRepository) ListPaymentMethods() ([]domain.PaymentMethod, error) {
	var methods []domain.PaymentMethod
	err := r.db.Where("is_active = true").Order("code asc").Find(&methods).Error
	return methods, err
}

// PMB Wave Operations
func (r *ReferenceRepository) ListPmbWaves() ([]domain.PmbWave, error) {
	var waves []domain.PmbWave
	err := r.db.Where("is_active = true").Order("code desc").Find(&waves).Error
	return waves, err
}

func (r *ReferenceRepository) CreatePmbWave(w *domain.PmbWave) error {
	return r.db.Create(w).Error
}

func (r *ReferenceRepository) UpdatePmbWave(w *domain.PmbWave) error {
	return r.db.Save(w).Error
}

func (r *ReferenceRepository) GetPmbWaveByID(id string) (*domain.PmbWave, error) {
	var w domain.PmbWave
	err := r.db.Where("id = ?", id).First(&w).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *ReferenceRepository) GetPmbWaveByCode(code string) (*domain.PmbWave, error) {
	var w domain.PmbWave
	err := r.db.Where("code = ?", code).First(&w).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// Master Data CRUD Actions
func (r *ReferenceRepository) CreatePaymentComponent(comp *domain.PaymentComponent) error {
	return r.db.Create(comp).Error
}

func (r *ReferenceRepository) UpdatePaymentComponent(comp *domain.PaymentComponent) error {
	return r.db.Save(comp).Error
}

func (r *ReferenceRepository) CreatePaymentMethod(m *domain.PaymentMethod) error {
	return r.db.Create(m).Error
}

func (r *ReferenceRepository) UpdatePaymentMethod(m *domain.PaymentMethod) error {
	return r.db.Save(m).Error
}

func (r *ReferenceRepository) CreateDocumentType(doc *domain.DocumentType) error {
	return r.db.Create(doc).Error
}

func (r *ReferenceRepository) UpdateDocumentType(doc *domain.DocumentType) error {
	return r.db.Save(doc).Error
}

func (r *ReferenceRepository) CreateReligion(rel *domain.Religion) error {
	return r.db.Create(rel).Error
}

func (r *ReferenceRepository) CreateCountry(c *domain.Country) error {
	return r.db.Create(c).Error
}

// Religion Operations
func (r *ReferenceRepository) ListReligions() ([]domain.Religion, error) {
	var religions []domain.Religion
	err := r.db.Order("name asc").Find(&religions).Error
	return religions, err
}

// Admission Path Operations
func (r *ReferenceRepository) ListAdmissionPaths() ([]domain.AdmissionPath, error) {
	var paths []domain.AdmissionPath
	err := r.db.Where("is_active = true").Order("name asc").Find(&paths).Error
	return paths, err
}

func (r *ReferenceRepository) CreateAdmissionPath(p *domain.AdmissionPath) error {
	return r.db.Create(p).Error
}

// Province Operations
func (r *ReferenceRepository) ListProvinces() ([]domain.Province, error) {
	var provinces []domain.Province
	err := r.db.Order("name asc").Find(&provinces).Error
	return provinces, err
}

func (r *ReferenceRepository) CreateProvince(p *domain.Province) error {
	return r.db.Create(p).Error
}

// City Operations
func (r *ReferenceRepository) ListCities() ([]domain.City, error) {
	var cities []domain.City
	err := r.db.Order("name asc").Find(&cities).Error
	return cities, err
}

func (r *ReferenceRepository) ListCitiesByProvince(provincesID string) ([]domain.City, error) {
	var cities []domain.City
	err := r.db.Where("province_id = ?", provincesID).Order("name asc").Find(&cities).Error
	return cities, err
}

func (r *ReferenceRepository) CreateCity(c *domain.City) error {
	return r.db.Create(c).Error
}

// District Operations
func (r *ReferenceRepository) ListDistricts() ([]domain.District, error) {
	var districts []domain.District
	err := r.db.Order("name asc").Find(&districts).Error
	return districts, err
}

func (r *ReferenceRepository) ListDistrictsByCity(cityID string) ([]domain.District, error) {
	var districts []domain.District
	err := r.db.Where("city_id = ?", cityID).Order("name asc").Find(&districts).Error
	return districts, err
}

func (r *ReferenceRepository) CreateDistrict(d *domain.District) error {
	return r.db.Create(d).Error
}

// Village Operations
func (r *ReferenceRepository) ListVillages() ([]domain.Village, error) {
	var villages []domain.Village
	err := r.db.Order("name asc").Find(&villages).Error
	return villages, err
}

func (r *ReferenceRepository) ListVillagesByDistrict(districtsID string) ([]domain.Village, error) {
	var villages []domain.Village
	err := r.db.Where("district_id = ?", districtsID).Order("name asc").Find(&villages).Error
	return villages, err
}

func (r *ReferenceRepository) CreateVillage(v *domain.Village) error {
	return r.db.Create(v).Error
}
