package domain

import (
	"time"
)

type AssessmentSession struct {
	ID              string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionSetID   *string    `gorm:"column:question_set_id"`
	SessionType     string     `gorm:"column:session_type;not null"` // cbt, mid_exam, final_exam
	ContextModule   *string    `gorm:"column:context_module"`        // pmb, lms
	ContextID       *string    `gorm:"column:context_id"`
	Title           string     `gorm:"column:title;not null"`
	StartAt         *time.Time `gorm:"column:start_at"`
	EndAt           *time.Time `gorm:"column:end_at"`
	DurationMinutes *int       `gorm:"column:duration_minutes"`
	Status          string     `gorm:"column:status;default:'active';not null"` // active, closed
	PassingGrade    *float64   `gorm:"column:passing_grade"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`

	Participants    []AssessmentParticipant `gorm:"foreignKey:AssessmentSessionID;references:ID"`
	Attempts        []AssessmentAttempt     `gorm:"foreignKey:AssessmentSessionID;references:ID"`
}

func (AssessmentSession) TableName() string {
	return "assessment_sessions"
}

type AssessmentParticipant struct {
	ID                  string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AssessmentSessionID string `gorm:"column:assessment_session_id;not null"`
	ParticipantType     string `gorm:"column:participant_type;not null"` // applicant, student
	ApplicantID         *string `gorm:"column:applicant_id"`             // external_ref: pmb.applicants.id
	StudentID           *string `gorm:"column:student_id"`               // external_ref: academic.students.id
	UserID              *string `gorm:"column:user_id"`                  // external_ref: core.users.id
	Status              string `gorm:"column:status;default:'registered';not null"` // registered, attended, completed
}

func (AssessmentParticipant) TableName() string {
	return "assessment_participants"
}

type AssessmentAttempt struct {
	ID                  string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AssessmentSessionID string     `gorm:"column:assessment_session_id;not null"`
	ParticipantID       string     `gorm:"column:participant_id;not null"`
	AttemptNumber       int        `gorm:"column:attempt_number;default:1;not null"`
	IdempotencyKey      *string    `gorm:"column:idempotency_key;unique"`
	StartedAt           time.Time  `gorm:"column:started_at;default:now()"`
	SubmittedAt         *time.Time `gorm:"column:submitted_at"`
	Status              string     `gorm:"column:status;default:'started';not null"` // started, submitted, evaluated
	TotalScore          *float64   `gorm:"column:total_score"`
}

func (AssessmentAttempt) TableName() string {
	return "assessment_attempts"
}

type QuestionBank struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code        string    `gorm:"column:code;unique;not null"`
	Name        string    `gorm:"column:name;not null"`
	ModuleScope string    `gorm:"column:module_scope"`
	OwnerUserID *string   `gorm:"column:owner_user_id"`
	Status      string    `gorm:"column:status;default:'active';not null"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (QuestionBank) TableName() string {
	return "question_banks"
}

type Question struct {
	ID                string           `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID    string           `gorm:"column:question_bank_id;not null"`
	QuestionType      string           `gorm:"column:question_type;not null"`
	Difficulty        string           `gorm:"column:difficulty;default:'MEDIUM';not null"`
	QuestionText      string           `gorm:"column:question_text;not null"`
	AnswerExplanation string           `gorm:"column:answer_explanation"`
	Status            string           `gorm:"column:status;default:'active';not null"`
	CreatedBy         *string          `gorm:"column:created_by"`
	CreatedAt         time.Time        `gorm:"column:created_at"`
	UpdatedAt         time.Time        `gorm:"column:updated_at"`
	Options           []QuestionOption `gorm:"foreignKey:QuestionID;references:ID"`
}

func (Question) TableName() string {
	return "questions"
}

type QuestionVersion struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID        string    `gorm:"column:question_id;not null"`
	VersionNumber     int       `gorm:"column:version_number;not null"`
	QuestionType      string    `gorm:"column:question_type;not null"`
	Difficulty        string    `gorm:"column:difficulty;not null"`
	QuestionText      string    `gorm:"column:question_text;not null"`
	AnswerExplanation string    `gorm:"column:answer_explanation"`
	OptionsSnapshot   string    `gorm:"type:jsonb;column:options_snapshot"`
	Status            string    `gorm:"column:status;default:'draft';not null"`
	CreatedBy         *string   `gorm:"column:created_by"`
	CreatedAt         time.Time `gorm:"column:created_at"`
}

func (QuestionVersion) TableName() string {
	return "question_versions"
}

type QuestionOption struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID  string `gorm:"column:question_id;not null"`
	OptionLabel string `gorm:"column:option_label;not null"`
	OptionText  string `gorm:"column:option_text;not null"`
	IsCorrect   bool   `gorm:"column:is_correct;default:false;not null"`
	SortOrder   int    `gorm:"column:sort_order;default:0;not null"`
}

func (QuestionOption) TableName() string {
	return "question_options"
}

type AssessmentAnswer struct {
	ID               string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	AttemptID        string     `gorm:"column:attempt_id;not null"`
	QuestionID       string     `gorm:"column:question_id;not null"`
	SelectedOptionID *string    `gorm:"column:selected_option_id"`
	AnswerText       string     `gorm:"column:answer_text"`
	Score            *float64   `gorm:"column:score"`
	GradedBy         *string    `gorm:"column:graded_by"`
	GradedAt         *time.Time `gorm:"column:graded_at"`
}

func (AssessmentAnswer) TableName() string {
	return "assessment_answers"
}

// ===========================================
// QUESTION BANK (BANK SOAL) MODELS
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

func (AssessmentQuestionBank) TableName() string {
	return "assessment_question_banks"
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

func (AssessmentQuestion) TableName() string {
	return "assessment_questions"
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

func (AssessmentQuestionOption) TableName() string {
	return "assessment_question_options"
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

func (AssessmentMatchingPair) TableName() string {
	return "assessment_matching_pairs"
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

func (AssessmentFillInBlank) TableName() string {
	return "assessment_fill_in_blanks"
}

// OrderingItem - Items for ordering/sequence question
type OrderingItem struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	ItemText   string `gorm:"column:item_text;not null"`
	CorrectOrder int   `gorm:"column:correct_order;not null"`
	MediaURL   string `gorm:"column:media_url"`
}

func (AssessmentOrderingItem) TableName() string {
	return "assessment_ordering_items"
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

func (AssessmentEssayAnswer) TableName() string {
	return "assessment_essay_answers"
}

// QuestionTag - Tags for categorizing questions
type QuestionTag struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	Tag        string `gorm:"column:tag;not null"`
}

func (AssessmentQuestionTag) TableName() string {
	return "assessment_question_tags"
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

func (AssessmentQuestionHistory) TableName() string {
	return "assessment_question_history"
}

// QuestionBlueprint - Templates for generating questions
type QuestionBlueprint struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID string  `gorm:"column:question_bank_id;not null"`
	Name        string    `gorm:"column:name;not null"`
	Template   string    `gorm:"column:template"` // question template with placeholders
	VariableSchema string `gorm:"column:variable_schema"` // JSON schema for variables
	CreatedAt  time.Time `gorm:"column:created_at;default:now()"`
}

func (AssessmentQuestionBlueprint) TableName() string {
	return "assessment_question_blueprints"
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

func (AssessmentRandomQuestionConfig) TableName() string {
	return "assessment_random_question_configs"
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

func (AssessmentQuestionPool) TableName() string {
	return "assessment_question_pools"
}

// QuestionCategory - Categories for organizing questions
type QuestionCategory struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionBankID string `gorm:"column:question_bank_id;not null"`
	Name        string `gorm:"column:name;not null"`
	ParentID   *string `gorm:"column:parent_id"` // hierarchical category
	Description string `gorm:"column:description"`
}

func (AssessmentQuestionCategory) TableName() string {
	return "assessment_question_categories"
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

func (AssessmentQuestionMetadata) TableName() string {
	return "assessment_question_metadata"
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

func (AssessmentQuestionAttachment) TableName() string {
	return "assessment_question_attachments"
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

func (AssessmentQuestionComment) TableName() string {
	return "assessment_question_comments"
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

func (AssessmentQuestionFlag) TableName() string {
	return "assessment_question_flags"
}

// QuestionTranslation - Multi-language question translations
type QuestionTranslation struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	QuestionID string `gorm:"column:question_id;not null"`
	Language  string `gorm:"column:language;not null"` // en, zh, ar, etc
	QuestionText string `gorm:"column:question_text"`
	Explanation string `gorm:"column:explanation"`
}

func (AssessmentQuestionTranslation) TableName() string {
	return "assessment_question_translations"
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

func (AssessmentQuestionAnswerFeedback) TableName() string {
	return "assessment_question_answer_feedback"
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

func (AssessmentQuestionStatistics) TableName() string {
	return "assessment_question_statistics"
}
