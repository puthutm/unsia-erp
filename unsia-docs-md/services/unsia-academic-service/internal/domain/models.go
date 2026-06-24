package domain

import (
	"time"
)

type Student struct {
	ID                  string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PersonID            string     `gorm:"column:person_id;not null"`           // external_ref: core.persons.id
	UserID              *string    `gorm:"column:user_id"`                     // external_ref: core.users.id
	ApplicantID         *string    `gorm:"column:applicant_id;unique"`         // external_ref: pmb.applicants.id
	StudyProgramID      string     `gorm:"column:study_program_id;not null"`   // external_ref: ref.study_programs.id
	Nim                 string     `gorm:"column:nim;unique;not null"`
	StudentStatus       string     `gorm:"column:student_status;default:'active';not null"` // active, inactive, graduated, drop_out
	EntryAcademicYearID *string    `gorm:"column:entry_academic_year_id"`     // external_ref: ref.academic_years.id
	EntryPeriodID       *string    `gorm:"column:entry_period_id"`            // external_ref: ref.academic_periods.id
	CurriculumID        *string    `gorm:"column:curriculum_id"`
	AdvisorID          *string    `gorm:"column:advisor_id"`               // external_ref: hris.lecturers.id (Dosen PA)
	CurrentSemester     int        `gorm:"column:current_semester;default:1;not null"`
	ActiveDate          *time.Time `gorm:"column:active_date"`
	GraduationDate      *time.Time `gorm:"column:graduation_date"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
}

func (Student) TableName() string {
	return "students"
}

type NimFormatConfig struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code           string    `gorm:"column:code;unique;not null"`
	FormatTemplate string    `gorm:"column:format_template"`
	TokenOrder     string    `gorm:"type:jsonb;column:token_order"`
	IsActive       bool      `gorm:"column:is_active;default:true;not null"`
	CreatedBy      *string   `gorm:"column:created_by"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (NimFormatConfig) TableName() string {
	return "nim_format_configs"
}

type NimSequence struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudyProgramID string    `gorm:"column:study_program_id;not null"`
	EntryPeriodID  string    `gorm:"column:entry_period_id;not null"`
	SequenceYear   string    `gorm:"column:sequence_year;not null"`
	LastNumber     int       `gorm:"column:last_number;default:0;not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (NimSequence) TableName() string {
	return "nim_sequences"
}

type CourseOffering struct {
	ID               string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AcademicPeriodID string `gorm:"column:academic_period_id;not null"` // external_ref: ref.academic_periods.id
	CourseID         string `gorm:"column:course_id;not null"`          // external_ref: ref.courses.id
	IsActive         bool   `gorm:"column:is_active;default:true;not null"`
}

func (CourseOffering) TableName() string {
	return "course_offerings"
}

type Class struct {
	ID               string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	CourseOfferingID string `gorm:"column:course_offering_id;not null"`
	ClassCode        string `gorm:"column:class_code;not null"`
	Quota            int    `gorm:"column:quota;default:40;not null"`
	EnrolledCount    int    `gorm:"column:enrolled_count;default:0;not null"`
	ClassStatus      string `gorm:"column:class_status;default:'active';not null"`
	IsParallel       bool   `gorm:"column:is_parallel;default:false;not null"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (Class) TableName() string {
	return "classes"
}

type KRS struct {
	ID                 string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID          string     `gorm:"column:student_id;not null"`
	AcademicPeriodID   string     `gorm:"column:academic_period_id;not null"`
	Status             string     `gorm:"column:status;default:'draft';not null"` // draft, submitted, approved, rejected
	AdvisorID          *string    `gorm:"column:advisor_id"`
	IsPackage          bool       `gorm:"column:is_package;default:false;not null"`
	FinanceClearanceID *string    `gorm:"column:finance_clearance_id"`
	SubmittedAt        *time.Time `gorm:"column:submitted_at"`
	ApprovedAt         *time.Time `gorm:"column:approved_at"`
	Items              []KrsItem  `gorm:"foreignKey:KrsID;references:ID"`
}

func (KRS) TableName() string {
	return "krs"
}

type KrsItem struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	KrsID      string    `gorm:"column:krs_id;not null"`
	ClassID    string    `gorm:"column:class_id;not null"`
	Status     string    `gorm:"column:status;default:'selected';not null"` // selected, approved, dropped
	SelectedAt time.Time `gorm:"column:selected_at"`
}

func (KrsItem) TableName() string {
	return "krs_items"
}

type Grade struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	KrsItemID    string     `gorm:"column:krs_item_id;not null;unique"`
	NumericGrade *float64   `gorm:"column:numeric_grade"`
	LetterGrade  string     `gorm:"column:letter_grade"`
	GradePoint   *float64   `gorm:"column:grade_point"`
	Status      string     `gorm:"column:status;default:'in_progress'"` // in_progress, submitted, final
	Source       string     `gorm:"column:source"`                      // lms, manual
	SubmittedAt  *time.Time `gorm:"column:submitted_at"`
	SubmittedBy  *string    `gorm:"column:submitted_by"`
	Components   []GradeComponent `gorm:"foreignKey:GradeID;references:ID"`
}

func (Grade) TableName() string {
	return "grades"
}

// GradeComponent represents a grading component (quiz, exam, assignment, etc.)
type GradeComponent struct {
	ID       string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	GradeID  string    `gorm:"column:grade_id;not null"`
	Name     string    `gorm:"column:name;not null"` // quiz, exam, assignment, etc.
	Weight  float64   `gorm:"column:weight;not null"`
	MaxScore float64  `gorm:"column:max_score;not null"`
}

func (GradeComponent) TableName() string {
	return "grade_components"
}

// GradeComponentScore stores individual component scores for a student
type GradeComponentScore struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	GradeEntryID string    `gorm:"column:grade_entry_id;not null"`
	Name        string    `gorm:"column:name;not null"`
	Score       float64   `gorm:"column:score"`
}

func (GradeComponentScore) TableName() string {
	return "grade_component_scores"
}

type Curriculum struct {
	ID                     string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudyProgramID         string     `gorm:"column:study_program_id;not null"`
	Code                   string     `gorm:"column:code;unique;not null"`
	Name                   string     `gorm:"column:name;not null"`
	CurriculumYear         int        `gorm:"column:curriculum_year;not null"`
	Status                 string     `gorm:"column:status;default:'draft';not null"`
	EffectiveStartPeriodID *string    `gorm:"column:effective_start_period_id"`
	EffectiveEndPeriodID   *string    `gorm:"column:effective_end_period_id"`
	IsActive               bool       `gorm:"column:is_active;default:true;not null"`
	IsDefaultForNewStudent bool       `gorm:"column:is_default_for_new_student;default:false;not null"`
	CreatedAt              time.Time  `gorm:"column:created_at"`
	UpdatedAt              time.Time  `gorm:"column:updated_at"`
}

func (Curriculum) TableName() string {
	return "curriculums"
}

type Course struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudyProgramID *string   `gorm:"column:study_program_id"`
	CourseCode     string    `gorm:"column:course_code;unique;not null"`
	CourseName     string    `gorm:"column:course_name;not null"`
	Sks            int       `gorm:"column:sks;default:2;not null"`
	CourseType     string    `gorm:"column:course_type"`
	MinimumGrade   *float64  `gorm:"column:minimum_grade"`
	IsActive       bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Course) TableName() string {
	return "courses"
}

type CurriculumCourse struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	CurriculumID string `gorm:"column:curriculum_id;not null"`
	CourseID     string `gorm:"column:course_id;not null"`
	Semester     int    `gorm:"column:semester;not null"`
	IsMandatory  bool   `gorm:"column:is_mandatory;default:true;not null"`
}

func (CurriculumCourse) TableName() string {
	return "curriculum_courses"
}

type ClassLecturer struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ClassID    string `gorm:"column:class_id;not null"`
	LecturerID string `gorm:"column:lecturer_id;not null"`
	RoleType   string `gorm:"column:role_type;default:'teacher';not null"`
}

func (ClassLecturer) TableName() string {
	return "class_lecturers"
}

// ClassSchedule represents the schedule for a class session
type ClassSchedule struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ClassID        string    `gorm:"column:class_id;not null"`
	DayOfWeek      int       `gorm:"column:day_of_week;not null"`        // 1=Monday, 7=Sunday
	StartTime      string    `gorm:"column:start_time;not null"`       // HH:MM format
	EndTime        string    `gorm:"column:end_time;not null"`         // HH:MM format
	RoomID         *string   `gorm:"column:room_id"`                 // references ref.ruangan.id
	BuildingID    *string   `gorm:"column:building_id"`              // references ref.gedung.id
	ScheduleType  string    `gorm:"column:schedule_type;default:'regular'"` // regular, praktikum, praktik
	IsOnline      bool      `gorm:"column:is_online;default:false;not null"`
	MeetingLink   *string   `gorm:"column:meeting_link"`             // for online classes
	CapacityUsed int       `gorm:"column:capacity_used;default:0"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (ClassSchedule) TableName() string {
	return "class_schedules"
}

// Room represents a classroom/venue
type Room struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	BuildingID   string    `gorm:"column:building_id;not null"`
	RoomCode    string    `gorm:"column:room_code;unique;not null"`
	RoomName    string    `gorm:"column:room_name;not null"`
	Capacity   int       `gorm:"column:capacity;default:40;not null"`
	RoomType   string    `gorm:"column:room_type"`              // lecture, lab, studio
	Floor      int       `gorm:"column:floor"`
	IsActive   bool     `gorm:"column:is_active;default:true;not null"`
	Equipment  string    `gorm:"type:jsonb;column:equipment"`  // json array of equipment
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Room) TableName() string {
	return "rooms"
}

// StudentAttendance represents student attendance for a class session
type StudentAttendance struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID  string    `gorm:"column:student_id;not null"`
	ClassID    string    `gorm:"column:class_id;not null"`
	SessionDate time.Time `gorm:"column:session_date;not null"`
	Status    string    `gorm:"column:status;not null"`         // present, absent, excused, sick
	Note       *string  `gorm:"column:note"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (StudentAttendance) TableName() string {
	return "student_attendances"
}

type GradeHistory struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	GradeID   string    `gorm:"column:grade_id;not null"`
	OldValue  string    `gorm:"type:jsonb;column:old_value"`
	NewValue  string    `gorm:"type:jsonb;column:new_value"`
	ChangedBy *string   `gorm:"column:changed_by"`
	Reason    string    `gorm:"column:reason"`
	ChangedAt time.Time `gorm:"column:changed_at"`
}

func (GradeHistory) TableName() string {
	return "grade_histories"
}

type KHS struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID        string    `gorm:"column:student_id;not null"`
	AcademicPeriodID string    `gorm:"column:academic_period_id;not null"`
	Ips              float64   `gorm:"column:ips"`
	TotalSks         int       `gorm:"column:total_sks"`
	FileURL          string    `gorm:"column:file_url"`
	IssuedAt         time.Time `gorm:"column:issued_at"`
}

func (KHS) TableName() string {
	return "khs"
}

type Transcript struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID string    `gorm:"column:student_id;not null"`
	Ipk       float64   `gorm:"column:ipk"`
	TotalSks  int       `gorm:"column:total_sks"`
	FileURL   string    `gorm:"column:file_url"`
	IssuedAt  time.Time `gorm:"column:issued_at"`
}

func (Transcript) TableName() string {
	return "transcripts"
}

// GradeConversion maps numeric scores to letter grades and grade points
type GradeConversion struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	GradeLetter string    `gorm:"column:grade_letter;unique;not null"`         // A, B, C, D, E
	MinScore   float64   `gorm:"column:min_score;not null"`              // Minimum score (e.g., 80 for A)
	MaxScore  float64   `gorm:"column:max_score;not null"`              // Maximum score (e.g., 100 for A)
	GradePoints float64   `gorm:"column:grade_points;not null"`            // Grade point (e.g., 4.0 for A)
	IsActive   bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (GradeConversion) TableName() string {
	return "grade_conversions"
}

// GradeEntry represents a student's grade entry for a course
type GradeEntry struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	GradeID     string     `gorm:"column:grade_id;not null"`
	StudentID   string     `gorm:"column:student_id;not null"`
	KrsItemID   string     `gorm:"column:krs_item_id;not null"`
	NumericGrade *float64   `gorm:"column:numeric_grade"`
	FinalGrade  *float64   `gorm:"column:final_grade"`      // Final numeric grade
	LetterGrade  *string    `gorm:"column:letter_grade"`
	GradePoint   *float64   `gorm:"column:grade_point"`
	Status      string     `gorm:"column:status;default:'draft'"` // draft, submitted, published
	EnteredAt   *time.Time `gorm:"column:entered_at"`
	SubmittedAt  *time.Time `gorm:"column:submitted_at"`
	SubmittedBy  *string    `gorm:"column:submitted_by"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
}

func (GradeEntry) TableName() string {
	return "grade_entries"
}

// GraduationRecord stores graduation data
type GraduationRecord struct {
	ID                 string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID          string     `gorm:"column:student_id;not null;unique"`
	GraduationPeriodID  string     `gorm:"column:graduation_period_id"`
	YudisiumNumber    string    `gorm:"column:yudisium_number;unique"`     // Ijazah number
	CertificateNumber string    `gorm:"column:certificate_number;unique"` // Certificate number
	Ipk              float64   `gorm:"column:ipk"`
	TotalSks          int       `gorm:"column:total_sks"`
	GraduatedAt       time.Time `gorm:"column:graduated_at"`
	Status           string    `gorm:"column:status;default:'pending'"` // pending, approved, published
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`
}

func (GraduationRecord) TableName() string {
	return "graduation_records"
}

// GradeAppeal represents a student's appeal for grade review
type GradeAppeal struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	GradeEntryID string    `gorm:"column:grade_entry_id;not null"`      // references grade_entries.id
	StudentID   string    `gorm:"column:student_id;not null"`
	KrsItemID   string    `gorm:"column:krs_item_id;not null"`
	Reason      string    `gorm:"type:text;column:reason;not null"` // Reason for appeal
	RequestedGrade *string `gorm:"column:requested_grade"`       // The grade student is requesting
	Status     string    `gorm:"column:status;default:'pending'"` // pending, approved, rejected
	ReviewedBy *string   `gorm:"column:reviewed_by"`
	ReviewedAt *time.Time `gorm:"column:reviewed_at"`
	ReviewNote *string   `gorm:"column:review_note"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (GradeAppeal) TableName() string {
	return "grade_appeals"
}

// AcademicPeriod represents academic year/semester periods
type AcademicPeriod struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AcademicYearID string   `gorm:"column:academic_year_id;not null"`  // references ref.academic_years.id
	PeriodName    string    `gorm:"column:period_name;not null"`    // e.g., "Semester Ganjil 2024/2025"
	PeriodType    string    `gorm:"column:period_type;not null"`   // ganjil, genap,短期
	StartDate     time.Time `gorm:"column:start_date;not null"`
	EndDate       time.Time `gorm:"column:end_date;not null"`
	IsActive     bool      `gorm:"column:is_active;default:true;not null"`
	CKrsOpen     bool      `gorm:"column:ckrs_open;default:false;not null"` // KRS registration open
	CKrsClose    bool      `gorm:"column:ckrs_close;default:false;not null"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (AcademicPeriod) TableName() string {
	return "academic_periods"
}

// StudyProgram represents a study program/major
type StudyProgram struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code         string    `gorm:"column:code;unique;not null"`        // e.g., "SI", "TI"
	Name        string    `gorm:"column:name;not null"`         // e.g., "Sistem Informasi"
	Degree      string    `gorm:"column:degree;not null"`        // S1, S2, S3, D4
	Capacity    int       `gorm:"column:capacity;default:60;not null"`
	Accreditation string  `gorm:"column:accreditation"`       // A, B, C
	IsActive    bool     `gorm:"column:is_active;default:true;not null"`
	HeadLecturerID *string `gorm:"column:head_lecturer_id"`      // references hris.lecturers.id
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (StudyProgram) TableName() string {
	return "study_programs"
}

// StudentAdvisor represents the student-advisor assignment (PA) relationship
type StudentAdvisor struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	StudentID    string    `gorm:"column:student_id;not null"`
	AdvisorID    string    `gorm:"column:advisor_id;not null"`     // references hris.lecturers.id
	AcademicYear string   `gorm:"column:academic_year;not null"` // e.g., "2024/2025"
	Semester     int       `gorm:"column:semester;not null"`      // 1 or 2
	AssignedAt   time.Time `gorm:"column:assigned_at"`
	IsActive     bool     `gorm:"column:is_active;default:true;not null"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (StudentAdvisor) TableName() string {
	return "student_advisors"
}

// Lecturer represents a lecturer/teacher
type Lecturer struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PersonID       string    `gorm:"column:person_id;not null"`
	Nik            string    `gorm:"column:nik;unique;not null"`        // Lecturer ID
	StudyProgramID *string  `gorm:"column:study_program_id"`          // Optional - for PA assignment
	AcademicTitle  string    `gorm:"column:academic_title"`          // Prof, Dr, etc.
	Position       string    `gorm:"column:position"`               // lecturer, head, etc.
	Status        string    `gorm:"column:status;default:'active'"` // active, inactive, retired
	IsActive       bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Lecturer) TableName() string {
	return "lecturers"
}
