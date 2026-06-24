package repository

import (
	"errors"
	"time"

	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"gorm.io/gorm"
)

type LMSRepository struct {
	db *gorm.DB
}

func NewLMSRepository(db *gorm.DB) *LMSRepository {
	return &LMSRepository{db: db}
}

func (r *LMSRepository) GetClassByAcademicID(academicClassID string) (*domain.Class, error) {
	var c domain.Class
	err := r.db.Where("academic_class_id = ?", academicClassID).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *LMSRepository) SyncClass(academicClassID string, lecturerID *string) (*domain.Class, error) {
	var c domain.Class
	err := r.db.Where("academic_class_id = ?", academicClassID).First(&c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c = domain.Class{
				AcademicClassID: academicClassID,
				LecturerID:      lecturerID,
				Status:          "active",
				SyncedAt:        time.Now(),
			}
			if err := r.db.Create(&c).Error; err != nil {
				return nil, err
			}
			return &c, nil
		}
		return nil, err
	}

	c.LecturerID = lecturerID
	c.SyncedAt = time.Now()
	if err := r.db.Save(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *LMSRepository) SyncEnrollment(lmsClassID string, studentID string, status string) (*domain.Enrollment, error) {
	var en domain.Enrollment
	err := r.db.Where("lms_class_id = ? AND student_id = ?", lmsClassID, studentID).First(&en).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			en = domain.Enrollment{
				LmsClassID:       lmsClassID,
				StudentID:        studentID,
				EnrollmentStatus: status,
				EnrolledAt:       time.Now(),
			}
			if err := r.db.Create(&en).Error; err != nil {
				return nil, err
			}
			return &en, nil
		}
		return nil, err
	}

	en.EnrollmentStatus = status
	if err := r.db.Save(&en).Error; err != nil {
		return nil, err
	}
	return &en, nil
}

func (r *LMSRepository) CreateSession(s *domain.Session) error {
	return r.db.Create(s).Error
}

func (r *LMSRepository) CreateMaterial(m *domain.Material) error {
	return r.db.Create(m).Error
}

func (r *LMSRepository) CreateAssignment(a *domain.Assignment) error {
	return r.db.Create(a).Error
}

func (r *LMSRepository) CreateAssignmentSubmission(as *domain.AssignmentSubmission) error {
	return r.db.Create(as).Error
}

func (r *LMSRepository) CreateAttendance(at *domain.Attendance) error {
	return r.db.Create(at).Error
}
