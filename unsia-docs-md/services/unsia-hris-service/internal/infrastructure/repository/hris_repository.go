package repository

import (
	"errors"

	"github.com/unsia-erp/unsia-hris-service/internal/domain"
	"gorm.io/gorm"
)

type HRISRepository struct {
	db *gorm.DB
}

func NewHRISRepository(db *gorm.DB) *HRISRepository {
	return &HRISRepository{db: db}
}

func (r *HRISRepository) ListActiveLecturers() ([]domain.Lecturer, error) {
	var lecturers []domain.Lecturer
	err := r.db.Where("is_active = true").Find(&lecturers).Error
	return lecturers, err
}

func (r *HRISRepository) GetLecturerByID(id string) (*domain.Lecturer, error) {
	var lec domain.Lecturer
	err := r.db.Where("id = ?", id).First(&lec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &lec, nil
}

func (r *HRISRepository) CreateEmployee(e *domain.Employee) error {
	return r.db.Create(e).Error
}

func (r *HRISRepository) GetEmployeeByID(id string) (*domain.Employee, error) {
	var emp domain.Employee
	err := r.db.Where("id = ?", id).First(&emp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &emp, nil
}

func (r *HRISRepository) CreateLecturer(l *domain.Lecturer) error {
	return r.db.Create(l).Error
}

func (r *HRISRepository) UpdateEmployeePosition(id string, positionID string) error {
	return r.db.Model(&domain.Employee{}).Where("id = ?", id).Update("position_id", positionID).Error
}

func (r *HRISRepository) CreateBkdRecord(br *domain.BkdRecord) error {
	return r.db.Create(br).Error
}
