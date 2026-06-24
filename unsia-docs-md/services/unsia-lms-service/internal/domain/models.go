package domain

import (
	"time"
)

type Class struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AcademicClassID string    `gorm:"column:academic_class_id;unique;not null"` // external_ref: academic.classes.id
	LecturerID      *string   `gorm:"column:lecturer_id"`                       // external_ref: hris.lecturers.id
	Status          string    `gorm:"column:status;default:'active';not null"`
	SyncedAt        time.Time `gorm:"column:synced_at;default:now()"`

	Enrollments     []Enrollment `gorm:"foreignKey:LmsClassID;references:ID"`
}

func (Class) TableName() string {
	return "classes"
}

type Enrollment struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	LmsClassID       string    `gorm:"column:lms_class_id;not null"`
	StudentID        string    `gorm:"column:student_id;not null"` // external_ref: academic.students.id
	EnrollmentStatus string    `gorm:"column:enrollment_status;default:'active';not null"`
	EnrolledAt       time.Time `gorm:"column:enrolled_at;default:now()"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}

type Session struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	LmsClassID    string     `gorm:"column:lms_class_id;not null"`
	SessionNumber int        `gorm:"column:session_number;not null"`
	Title         string     `gorm:"column:title;not null"`
	SessionDate   *time.Time `gorm:"column:session_date"`
	StartTime     *string    `gorm:"column:start_time"`
	EndTime       *string    `gorm:"column:end_time"`
	Status        string     `gorm:"column:status;default:'draft';not null"`
}

func (Session) TableName() string {
	return "sessions"
}

type Material struct {
	ID                   string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	SessionID            string     `gorm:"column:session_id;not null"`
	AssessmentMaterialID *string    `gorm:"column:assessment_material_id"`
	Title                string     `gorm:"column:title;not null"`
	ContentType          string     `gorm:"column:content_type"`
	FileURL              string     `gorm:"column:file_url"`
	PublishedAt          *time.Time `gorm:"column:published_at"`
}

func (Material) TableName() string {
	return "materials"
}

type Assignment struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	SessionID   string     `gorm:"column:session_id;not null"`
	Title       string     `gorm:"column:title;not null"`
	Instruction string     `gorm:"column:instruction"`
	DueAt       *time.Time `gorm:"column:due_at"`
	Status      string     `gorm:"column:status;default:'active';not null"`
}

func (Assignment) TableName() string {
	return "assignments"
}

type AssignmentSubmission struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AssignmentID string     `gorm:"column:assignment_id;not null"`
	StudentID    string     `gorm:"column:student_id;not null"`
	FileURL      string     `gorm:"column:file_url"`
	SubmittedAt  time.Time  `gorm:"column:submitted_at;default:now()"`
	Score        *float64   `gorm:"column:score"`
	GradedBy     *string    `gorm:"column:graded_by"`
	GradedAt     *time.Time `gorm:"column:graded_at"`
}

func (AssignmentSubmission) TableName() string {
	return "assignment_submissions"
}

type Attendance struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	SessionID        string    `gorm:"column:session_id;not null"`
	StudentID        string    `gorm:"column:student_id;not null"`
	AttendanceStatus string    `gorm:"column:attendance_status;default:'present';not null"`
	SubmittedAt      time.Time `gorm:"column:submitted_at;default:now()"`
}

func (Attendance) TableName() string {
	return "attendances"
}

// ===========================================
// QUESTION BANK MODELS (Bank Soals)
// ===========================================

// QuestionBank - Core entity for question bank
type QuestionBank struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	CourseID    string    `gorm:"column:course_id;not null"` // external_ref: academic.courses.id
	Subject     string    `gorm:"column:subject"`
	Topic       string    `gorm:"column:topic"`
	CreatedBy   string    `gorm:"column:created_by"` // external_ref: hris.lecturers.id
	CreatedAt  time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:now()"`
	IsActive   bool      `gorm:"column:is_active;default:true"`
	Difficulty int       `gorm:"column:difficulty"` // 1-5 scale
}

func (QuestionBank) TableName() string {
	return "question_banks"
}

// Question - Main question entity
type Question struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID string   `gorm:"column:question_bank_id;not null"`
	QuestionType  string    `gorm:"column:question_type;not null"` // multiple_choice, true_false, short_answer, essay, matching, fill_blank, ordering
	QuestionText string    `gorm:"column:question_text;not null"`
	MediaURL     string    `gorm:"column:media_url"`
	MediaType    string    `gorm:"column:media_type"` // image, audio, video
	Points      float64   `gorm:"column:points;default:1.0"`
	Explanation string    `gorm:"column:explanation"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;default:now()"`
}

func (Question) TableName() string {
	return "questions"
}

// QuestionOption - Options for multiple choice/true_false/matching questions
type QuestionOption struct {
	ID         string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	OptionKey  string `gorm:"column:option_key;not null"` // A, B, C, D, E
	OptionText string `gorm:"column:option_text;not null"`
	MediaURL   string `gorm:"column:media_url"`
	IsCorrect bool   `gorm:"column:is_correct;default:false"`
	Order     int    `gorm:"column:order"`
}

func (QuestionOption) TableName() string {
	return "question_options"
}

// MatchingPair - Pairs for matching question type
type MatchingPair struct {
	ID            string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID    string `gorm:"column:question_id;not null"`
	LeftItem      string `gorm:"column:left_item;not null"`
	RightItem    string `gorm:"column:right_item;not null"`
	LeftMediaURL  string `gorm:"column:left_media_url"`
	RightMediaURL string `gorm:"column:right_media_url"`
}

func (MatchingPair) TableName() string {
	return "matching_pairs"
}

// FillInBlank - Fill in the blank segments
type FillInBlank struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	BlankKey   string `gorm:"column:blank_key"` // __1__, __2__
	Answer     string `gorm:"column:answer;not null"`
	AlternateAnswers string `gorm:"column:alternate_answers"` // JSON array for multiple accepted answers
	IsCaseSensitive bool `gorm:"column:is_case_sensitive;default:false"`
	Order       int    `gorm:"column:order"`
}

func (FillInBlank) TableName() string {
	return "fill_in_blanks"
}

// OrderingItem - Items for ordering/sequence question
type OrderingItem struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	ItemText   string `gorm:"column:item_text;not null"`
	CorrectOrder int   `gorm:"column:correct_order;not null"`
	MediaURL   string `gorm:"column:media_url"`
}

func (OrderingItem) TableName() string {
	return "ordering_items"
}

// EssayAnswer - Expected answer for short answer/essay questions
type EssayAnswer struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	ModelAnswer string `gorm:"column:model_answer"`
	KeywordAnswers string `gorm:"column:keyword_answers"` // JSON array of keywords
	MinWords int `gorm:"column:min_words"` // minimum word count for essay
	MaxWords int `gorm:"column:max_words"`
	Rubric   string `gorm:"column:rubric"` // JSON rubric for scoring
}

func (EssayAnswer) TableName() string {
	return "essay_answers"
}

// QuestionTag - Tags for categorizing questions
type QuestionTag struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	Tag        string `gorm:"column:tag;not null"`
}

func (QuestionTag) TableName() string {
	return "question_tags"
}

// QuestionHistory - Version history for questions
type QuestionHistory struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string    `gorm:"column:question_id;not null"`
	ChangedBy  string    `gorm:"column:changed_by"`
	ChangeType string   `gorm:"column:change_type"` // created, updated, deactivated
	Changes   string    `gorm:"column:changes"` // JSON diff
	CreatedAt time.Time `gorm:"column:created_at;default:now()"`
}

func (QuestionHistory) TableName() string {
	return "question_history"
}

// ===========================================
// ADDITIONAL QUESTION BANK MODELS
// ===========================================

// QuestionBlueprint - Templates for generating questions
type QuestionBlueprint struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID string  `gorm:"column:question_bank_id;not null"`
	Name        string    `gorm:"column:name;not null"`
	Template   string    `gorm:"column:template"` // question template with placeholders
	VariableSchema string `gorm:"column:variable_schema"` // JSON schema for variables
	CreatedAt  time.Time `gorm:"column:created_at;default:now()"`
}

func (QuestionBlueprint) TableName() string {
	return "question_blueprints"
}

// RandomQuestionConfig - Configuration for random question generation
type RandomQuestionConfig struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID string `gorm:"column:question_bank_id;not null"`
	TotalQuestions int    `gorm:"column:total_questions;not null"`
	DifficultyDistribution string `gorm:"column:difficulty_distribution"` // JSON: {"easy": 10, "medium": 20, "hard": 15}
	TopicDistribution string   `gorm:"column:topic_distribution"` // JSON
	RandomizeOptions bool     `gorm:"column:randomize_options;default:true"`
	RandomizeOrder bool      `gorm:"column:randomize_order;default:true"`
}

func (RandomQuestionConfig) TableName() string {
	return "random_question_configs"
}

// QuestionPool - Pool of questions for random selection
type QuestionPool struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID string  `gorm:"column:question_bank_id;not null"`
	SessionID    *string  `gorm:"column:session_id"` // linked exam session
	QuestionIDs string    `gorm:"column:question_ids"` // JSON array of selected question IDs
	SelectedAt  time.Time `gorm:"column:selected_at;default:now()"`
	Status     string    `gorm:"column:status;default:'pending'"` // pending, active, completed
}

func (QuestionPool) TableName() string {
	return "question_pools"
}

// QuestionCategory - Categories for organizing questions
type QuestionCategory struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID string `gorm:"column:question_bank_id;not null"`
	Name        string `gorm:"column:name;not null"`
	ParentID   *string `gorm:"column:parent_id"` // hierarchical category
	Description string `gorm:"column:description"`
}

func (QuestionCategory) TableName() string {
	return "question_categories"
}

// QuestionMetadata - Extended metadata for questions
type QuestionMetadata struct {
	ID             string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID     string `gorm:"column:question_id;not null"`
	BloomLevel     string `gorm:"column:bloom_level"` // remember, understand, apply, analyze, evaluate, create
	LearningOutcome string `gorm:"column:learning_outcome"`
	CurriculumMapping string `gorm:"column:curriculum_mapping"`
	Semester      int    `gorm:"column:semester"`
	Year         int    `gorm:"column:year"`
	UsageCount    int    `gorm:"column:usage_count;default:0"`
	SuccessRate   float64 `gorm:"column:success_rate"`
	AvgTimeSpent  int    `gorm:"column:avg_time_spent"` // seconds
}

func (QuestionMetadata) TableName() string {
	return "question_metadata"
}

// QuestionAttachment - Attachments like images, files, audio
type QuestionAttachment struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	FileName   string `gorm:"column:file_name;not null"`
	FileURL   string `gorm:"column:file_url;not null"`
	FileType  string `gorm:"column:file_type"` // image, audio, video, document
	FileSize  int    `gorm:"column:file_size"`
	MimeType  string `gorm:"column:mime_type"`
	Order     int    `gorm:"column:order"`
}

func (QuestionAttachment) TableName() string {
	return "question_attachments"
}

// QuestionComment - Comments/feedback on questions
type QuestionComment struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string    `gorm:"column:question_id;not null"`
	UserID    string    `gorm:"column:user_id;not null"`
	Comment   string    `gorm:"column:comment;not null"`
	ParentID  *string   `gorm:"column:parent_id"` // threaded comments
	CreatedAt time.Time `gorm:"column:created_at;default:now()"`
}

func (QuestionComment) TableName() string {
	return "question_comments"
}

// QuestionFlag - Flag questions for review
type QuestionFlag struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string    `gorm:"column:question_id;not null"`
	UserID    string    `gorm:"column:user_id;not null"`
	Reason    string    `gorm:"column:reason"` // incorrect, ambiguous, typo, duplicate
	Status   string    `gorm:"column:status;default:'open'"` // open, resolved, rejected
	ResolvedBy *string `gorm:"column:resolved_by"`
	ResolvedAt *time.Time `gorm:"column:resolved_at"`
}

func (QuestionFlag) TableName() string {
	return "question_flags"
}

// QuestionTranslation - Multi-language question translations
type QuestionTranslation struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	Language  string `gorm:"column:language;not null"` // en, zh, ar, etc
	QuestionText string `gorm:"column:question_text"`
	Explanation string `gorm:"column:explanation"`
}

func (QuestionTranslation) TableName() string {
	return "question_translations"
}

// QuestionAnswerFeedback - Pre-defined feedback for answers
type QuestionAnswerFeedback struct {
	ID           string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID  string `gorm:"column:question_id;not null"`
	AnswerKey   string `gorm:"column:answer_key"` // A, B, C, D, E or __1__, __2__ for fill_blank
	Feedback    string `gorm:"column:feedback"`
	CorrectMsg  string `gorm:"column:correct_msg"`
	IncorrectMsg string `gorm:"column:incorrect_msg"`
	PartialMsg string `gorm:"column:partial_msg"` // for partial credit
}

func (QuestionAnswerFeedback) TableName() string {
	return "question_answer_feedback"
}

// QuestionStatistics - Real-time statistics for questions
type QuestionStatistics struct {
	ID              string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID      string `gorm:"column:question_id;not null;uniqueIndex"`
	TotalAttempts  int    `gorm:"column:total_attempts;default:0"`
	CorrectAttempts int   `gorm:"column:correct_attempts;default:0"`
	PartialAttempts int  `gorm:"column:partial_attempts;default:0"`
	SkippedAttempts int  `gorm:"column:skipped_attempts;default:0"`
	AvgScore      float64 `gorm:"column:avg_score;default:0"`
	AvgTimeSpent  int    `gorm:"column:avg_time_spent;default:0"` // seconds
	DiscriminationIndex float64 `gorm:"column:discrimination_index"`
	DifficultyIndex float64 `gorm:"column:difficulty_index"`
	LastUpdated  time.Time `gorm:"column:last_updated;default:now()"`
}

func (QuestionStatistics) TableName() string {
	return "question_statistics"
}
