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

func (r *HRISRepository) ListEmployees(limit, offset int) ([]domain.Employee, int64, error) {
	var list []domain.Employee
	var total int64
	r.db.Model(&domain.Employee{}).Count(&total)
	err := r.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

func (r *HRISRepository) UpdateEmployee(e *domain.Employee) error {
	return r.db.Save(e).Error
}

func (r *HRISRepository) ListAttendances(limit, offset int) ([]domain.Attendance, int64, error) {
	var list []domain.Attendance
	var total int64
	r.db.Model(&domain.Attendance{}).Count(&total)
	err := r.db.Limit(limit).Offset(offset).Order("attendance_date desc").Find(&list).Error
	return list, total, err
}

func (r *HRISRepository) RecordAttendance(a *domain.Attendance) error {
	// Check if already checked in for today
	var existing domain.Attendance
	err := r.db.Where("employee_id = ? AND attendance_date = ?", a.EmployeeID, a.AttendanceDate.Format("2006-01-02")).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(a).Error
		}
		return err
	}

	// Update check out if check in exists
	if a.CheckOut != nil {
		existing.CheckOut = a.CheckOut
	}
	existing.Status = a.Status
	return r.db.Save(&existing).Error
}

func (r *HRISRepository) ListLeaveRequests(limit, offset int) ([]domain.LeaveRequest, int64, error) {
	var list []domain.LeaveRequest
	var total int64
	r.db.Model(&domain.LeaveRequest{}).Count(&total)
	err := r.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&list).Error
	return list, total, err
}

func (r *HRISRepository) CreateLeaveRequest(l *domain.LeaveRequest) error {
	return r.db.Create(l).Error
}

