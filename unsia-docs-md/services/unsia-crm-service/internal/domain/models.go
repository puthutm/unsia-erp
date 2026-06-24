package domain

import (
	"time"
)

type Campaign struct {
	ID        string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code      string     `gorm:"column:code;unique;not null"`
	Name      string     `gorm:"column:name;not null"`
	Channel   string     `gorm:"column:channel"`
	StartDate *time.Time `gorm:"column:start_date"`
	EndDate   *time.Time `gorm:"column:end_date"`
	Status    string     `gorm:"column:status"`
	CreatedBy string     `gorm:"column:created_by"` // external_ref: core.users.id
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
}

func (Campaign) TableName() string {
	return "campaigns"
}

type Agent struct {
	ID               string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PersonID         string     `gorm:"column:person_id;not null"` // external_ref: core.persons.id
	AgentCode        string     `gorm:"column:agent_code;unique;not null"`
	OrganizationName string     `gorm:"column:organization_name"`
	Status           string     `gorm:"column:status;default:'active';not null"`
	ApprovalStatus   string     `gorm:"column:approval_status;default:'pending';not null"`
	ApprovedBy       *string    `gorm:"column:approved_by"` // external_ref: core.users.id
	ApprovedAt       *time.Time `gorm:"column:approved_at"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
}

func (Agent) TableName() string {
	return "agents"
}

type Referral struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ReferralType     string    `gorm:"column:referral_type;not null"`
	ReferrerPersonID *string   `gorm:"column:referrer_person_id"` // external_ref: core.persons.id
	AgentID          *string   `gorm:"column:agent_id"`
	ReferralCode     string    `gorm:"column:referral_code;unique;not null"`
	IsValid          bool      `gorm:"column:is_valid;default:true;not null"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (Referral) TableName() string {
	return "referrals"
}

type Lead struct {
	ID             string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PersonID       string     `gorm:"column:person_id;not null"` // external_ref: core.persons.id
	StudyProgramID *string    `gorm:"column:study_program_id"`  // external_ref: ref.study_programs.id
	LeadSourceID   *string    `gorm:"column:lead_source_id"`    // external_ref: ref.lead_sources.id
	CampaignID     *string    `gorm:"column:campaign_id"`
	ReferralID     *string    `gorm:"column:referral_id"`
	LeadNumber     string     `gorm:"column:lead_number;unique;not null"`
	Status         string     `gorm:"column:status;default:'new';not null"`
	OwnerUserID    *string    `gorm:"column:owner_user_id"` // external_ref: core.users.id
	ConvertedAt    *time.Time `gorm:"column:converted_at"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
}

func (Lead) TableName() string {
	return "leads"
}

type LeadActivity struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	LeadID       string    `gorm:"column:lead_id;not null"`
	UserID       *string   `gorm:"column:user_id"` // external_ref: core.users.id
	ActivityType string    `gorm:"column:activity_type"`
	Note         string    `gorm:"column:note"`
	ActivityAt   time.Time `gorm:"column:activity_at;default:now()"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (LeadActivity) TableName() string {
	return "lead_activities"
}

type LeadStatusHistory struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	LeadID    string    `gorm:"column:lead_id;not null"`
	OldStatus *string   `gorm:"column:old_status"`
	NewStatus string    `gorm:"column:new_status;not null"`
	ChangedBy *string   `gorm:"column:changed_by"` // external_ref: core.users.id
	Note      string    `gorm:"column:note"`
	ChangedAt time.Time `gorm:"column:changed_at;default:now()"`
}

func (LeadStatusHistory) TableName() string {
	return "lead_status_histories"
}

type CommissionRule struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ReferralType    string    `gorm:"column:referral_type;not null"`
	Amount          float64   `gorm:"column:amount;default:0.00;not null"`
	CalculationType string    `gorm:"column:calculation_type;default:'fixed';not null"`
	IsActive        bool      `gorm:"column:is_active;default:true;not null"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (CommissionRule) TableName() string {
	return "commission_rules"
}

type CommissionRecord struct {
	ID               string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	LeadID           string     `gorm:"column:lead_id;not null"`
	CommissionRuleID *string    `gorm:"column:commission_rule_id"`
	ReferrerPersonID *string    `gorm:"column:referrer_person_id"` // external_ref: core.persons.id
	Amount           float64    `gorm:"column:amount;default:0.00;not null"`
	Status           string     `gorm:"column:status;default:'draft';not null"`
	SentToFinanceAt  *time.Time `gorm:"column:sent_to_finance_at"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
}

func (CommissionRecord) TableName() string {
	return "commission_records"
}
