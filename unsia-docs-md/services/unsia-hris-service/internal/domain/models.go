package domain

import (
	"time"
)

type Employee struct {
	ID               string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PersonID         string     `gorm:"column:person_id;not null"` // external_ref: core.persons.id
	EmployeeTypeID   *string    `gorm:"column:employee_type_id"`   // external_ref: ref.employee_types.id
	WorkUnitID       *string    `gorm:"column:work_unit_id"`
	PositionID       *string    `gorm:"column:position_id"`
	Nip              *string    `gorm:"column:nip;unique"`
	EmploymentStatus string     `gorm:"column:employment_status;default:'contract';not null"` // contract, permanent
	JoinDate         *time.Time `gorm:"column:join_date"`
	EndDate           *time.Time `gorm:"column:end_date"`
	IsActive         bool       `gorm:"column:is_active;default:true;not null"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`

	Lecturer         *Lecturer  `gorm:"foreignKey:EmployeeID;references:ID"`
}

func (Employee) TableName() string {
	return "employees"
}

type Lecturer struct {
	ID                     string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EmployeeID             string    `gorm:"column:employee_id;unique;not null"`
	LecturerStatusID       *string   `gorm:"column:lecturer_status_id"` // external_ref: ref.lecturer_statuses.id
	FunctionalPositionID   *string   `gorm:"column:functional_position_id"`
	Nidn                   *string   `gorm:"column:nidn;unique"`
	HomebaseStudyProgramID *string   `gorm:"column:homebase_study_program_id"` // external_ref: ref.study_programs.id
	CertificationStatus    *string   `gorm:"column:certification_status"`
	IsActive               bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt              time.Time `gorm:"column:created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at"`
}

func (Lecturer) TableName() string {
	return "lecturers"
}

type BkdRecord struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	LecturerID       string    `gorm:"column:lecturer_id;not null"`
	AcademicPeriodID *string   `gorm:"column:academic_period_id"`
	TeachingLoad     float64   `gorm:"column:teaching_load;default:0.00;not null"`
	ResearchLoad     float64   `gorm:"column:research_load;default:0.00;not null"`
	ServiceLoad      float64   `gorm:"column:service_load;default:0.00;not null"`
	Status           string    `gorm:"column:status;default:'draft';not null"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (BkdRecord) TableName() string {
	return "bkd_records"
}

type Attendance struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EmployeeID     string    `gorm:"column:employee_id;not null"`
	AttendanceDate time.Time `gorm:"column:attendance_date;not null"`
	CheckIn        *string   `gorm:"column:check_in"`
	CheckOut       *string   `gorm:"column:check_out"`
	Status         string    `gorm:"column:status;default:'present';not null"`
}

func (Attendance) TableName() string {
	return "attendances"
}

type LeaveRequest struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	EmployeeID  string     `gorm:"column:employee_id;not null"`
	LeaveType   string     `gorm:"column:leave_type;not null"`
	StartDate   time.Time  `gorm:"column:start_date;not null"`
	EndDate     time.Time  `gorm:"column:end_date;not null"`
	Status      string     `gorm:"column:status;default:'pending';not null"`
	ApprovedBy  *string    `gorm:"column:approved_by"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (LeaveRequest) TableName() string {
	return "leave_requests"
}

