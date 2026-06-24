package repository

import (
	"errors"
	"time"

	sharedrbac "github.com/unsia-erp/shared-rbac"
	"github.com/unsia-erp/unsia-crm-service/internal/domain"
	"gorm.io/gorm"
)

type CRMRepository struct {
	db *gorm.DB
}

func NewCRMRepository(db *gorm.DB) *CRMRepository {
	return &CRMRepository{db: db}
}

// Campaign Operations
func (r *CRMRepository) ListCampaigns() ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.db.Order("created_at desc").Find(&campaigns).Error
	return campaigns, err
}

func (r *CRMRepository) CreateCampaign(c *domain.Campaign) error {
	return r.db.Create(c).Error
}

func (r *CRMRepository) GetCampaignByID(id string) (*domain.Campaign, error) {
	var c domain.Campaign
	err := r.db.Where("id = ?", id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Agent Operations
func (r *CRMRepository) ListAgents() ([]domain.Agent, error) {
	var agents []domain.Agent
	err := r.db.Order("created_at desc").Find(&agents).Error
	return agents, err
}

func (r *CRMRepository) GetAgentByID(id string) (*domain.Agent, error) {
	var a domain.Agent
	err := r.db.Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *CRMRepository) GetAgentByPersonID(personID string) (*domain.Agent, error) {
	var a domain.Agent
	err := r.db.Where("person_id = ?", personID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *CRMRepository) CreateAgent(a *domain.Agent) error {
	return r.db.Create(a).Error
}

func (r *CRMRepository) UpdateAgentApproval(id string, approvalStatus string, approvedBy string, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"approval_status": approvalStatus,
		"approved_by":     approvedBy,
		"approved_at":     &now,
		"status":          status,
		"updated_at":      now,
	}
	return r.db.Model(&domain.Agent{}).Where("id = ?", id).Updates(updates).Error
}

// Referral Operations
func (r *CRMRepository) ListReferrals() ([]domain.Referral, error) {
	var referrals []domain.Referral
	err := r.db.Order("created_at desc").Find(&referrals).Error
	return referrals, err
}

func (r *CRMRepository) GetReferralByID(id string) (*domain.Referral, error) {
	var ref domain.Referral
	err := r.db.Where("id = ?", id).First(&ref).Error
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *CRMRepository) GetReferralByCode(code string) (*domain.Referral, error) {
	var ref domain.Referral
	err := r.db.Where("referral_code = ? AND is_valid = true", code).First(&ref).Error
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *CRMRepository) CreateReferral(ref *domain.Referral) error {
	return r.db.Create(ref).Error
}

// Lead Operations
func (r *CRMRepository) ListLeads(scope sharedrbac.UserScope) ([]domain.Lead, error) {
	var leads []domain.Lead
	query := r.db.Order("created_at desc")

	if scope.Type == "own_lead" || scope.Type == "self" {
		// Attempt to see if this user is also an agent
		var agent domain.Agent
		err := r.db.Where("person_id = ?", scope.UserID).First(&agent).Error
		if err == nil {
			// If they are an agent, they can see leads owned by them OR leads referred by them
			query = query.Where("owner_user_id = ? OR referral_id IN (SELECT id FROM referrals WHERE agent_id = ?)", scope.UserID, agent.ID)
		} else {
			// Otherwise just filter by owner_user_id
			query = query.Where("owner_user_id = ?", scope.UserID)
		}
	}

	err := query.Find(&leads).Error
	return leads, err
}

func (r *CRMRepository) GetLeadByID(id string) (*domain.Lead, error) {
	var l domain.Lead
	err := r.db.Where("id = ?", id).First(&l).Error
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *CRMRepository) CreateLead(l *domain.Lead) error {
	return r.db.Create(l).Error
}

func (r *CRMRepository) UpdateLeadStatus(id string, status string, convertedAt *time.Time) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if convertedAt != nil {
		updates["converted_at"] = convertedAt
	}
	return r.db.Model(&domain.Lead{}).Where("id = ?", id).Updates(updates).Error
}

// Lead Activity Operations
func (r *CRMRepository) ListLeadActivities(leadID string) ([]domain.LeadActivity, error) {
	var activities []domain.LeadActivity
	err := r.db.Where("lead_id = ?", leadID).Order("activity_at desc").Find(&activities).Error
	return activities, err
}

func (r *CRMRepository) CreateLeadActivity(act *domain.LeadActivity) error {
	return r.db.Create(act).Error
}

// Commission Rule Operations
func (r *CRMRepository) ListCommissionRules() ([]domain.CommissionRule, error) {
	var rules []domain.CommissionRule
	err := r.db.Where("is_active = true").Find(&rules).Error
	return rules, err
}

func (r *CRMRepository) GetCommissionRuleByReferralType(refType string) (*domain.CommissionRule, error) {
	var rule domain.CommissionRule
	err := r.db.Where("referral_type = ? AND is_active = true", refType).First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

// Commission Record Operations
func (r *CRMRepository) CreateCommissionRecord(cr *domain.CommissionRecord) error {
	return r.db.Create(cr).Error
}

func (r *CRMRepository) ListCommissionRecords() ([]domain.CommissionRecord, error) {
	var records []domain.CommissionRecord
	err := r.db.Order("created_at desc").Find(&records).Error
	return records, err
}
