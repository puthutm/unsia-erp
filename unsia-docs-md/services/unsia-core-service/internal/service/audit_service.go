package service

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

type CreateAuditLogInput struct {
	UserID                *string `json:"user_id"`
	ActorUserID           *string `json:"actor_user_id"`
	TargetUserID          *string `json:"target_user_id"`
	ActiveRoleID          *string `json:"active_role_id"`
	ImpersonationSessionID *string `json:"impersonation_session_id"`
	ApplicationID         *string `json:"application_id"`
	Module               string  `json:"module" binding:"required"`
	Action               string  `json:"action" binding:"required"`
	EntityName            *string `json:"entity_name"`
	EntityID              *string `json:"entity_id"`
	Reason                *string `json:"reason"`
	OldValue              *string `json:"old_value"`
	NewValue              *string `json:"new_value"`
	RequestID             *string `json:"request_id"`
	IPAddress             *string `json:"ip_address"`
	UserAgent             *string `json:"user_agent"`
}

// CreateAuditLog creates a new audit log entry
func (s *AuditService) CreateAuditLog(input CreateAuditLogInput) (*AuditLog, error) {
	auditLog := AuditLog{
		ID:                      uuid.New().String(),
		UserID:                  input.UserID,
		ActorUserID:              input.ActorUserID,
		TargetUserID:             input.TargetUserID,
		ActiveRoleID:             input.ActiveRoleID,
		ImpersonationSessionID:  input.ImpersonationSessionID,
		ApplicationID:          input.ApplicationID,
		Module:                input.Module,
		Action:                input.Action,
		EntityName:             input.EntityName,
		EntityID:              input.EntityID,
		Reason:                input.Reason,
		OldValue:              input.OldValue,
		NewValue:              input.NewValue,
		RequestID:             input.RequestID,
		IPAddress:             input.IPAddress,
		UserAgent:              input.UserAgent,
		CreatedAt:             time.Now(),
	}

	if err := s.db.Create(&auditLog).Error; err != nil {
		return nil, err
	}

	return &auditLog, nil
}

// GetAuditLogByID retrieves audit log by ID
func (s *AuditService) GetAuditLogByID(id string) (*AuditLog, error) {
	var auditLog AuditLog
	if err := s.db.First(&auditLog, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &auditLog, nil
}

// ListAuditLogs lists audit logs with filters
type ListAuditLogsFilter struct {
	UserID       *string `json:"user_id"`
	ActorUserID  *string `json:"actor_user_id"`
	Module       *string `json:"module"`
	Action       *string `json:"action"`
	EntityName   *string `json:"entity_name"`
	EntityID     *string `json:"entity_id"`
	ApplicationID *string `json:"application_id"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Limit       int      `json:"limit"`
	Offset      int      `json:"offset"`
}

func (s *AuditService) ListAuditLogs(filter ListAuditLogsFilter) ([]AuditLog, int64, error) {
	query := s.db.Model(&AuditLog{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.ActorUserID != nil {
		query = query.Where("actor_user_id = ?", *filter.ActorUserID)
	}
	if filter.Module != nil {
		query = query.Where("module = ?", *filter.Module)
	}
	if filter.Action != nil {
		query = query.Where("action = ?", *filter.Action)
	}
	if filter.EntityName != nil {
		query = query.Where("entity_name = ?", *filter.EntityName)
	}
	if filter.EntityID != nil {
		query = query.Where("entity_id = ?", *filter.EntityID)
	}
	if filter.ApplicationID != nil {
		query = query.Where("application_id = ?", *filter.ApplicationID)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Apply pagination
	if filter.Limit == 0 {
		filter.Limit = 50
	}
	if filter.Offset == 0 {
		filter.Offset = 0
	}

	var auditLogs []AuditLog
	err := query.Order("created_at DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&auditLogs).Error

	return auditLogs, total, err
}

// GetAuditLogsByEntity retrieves audit logs for a specific entity
func (s *AuditService) GetAuditLogsByEntity(entityName, entityID string) ([]AuditLog, error) {
	var auditLogs []AuditLog
	err := s.db.Where("entity_name = ? AND entity_id = ?", entityName, entityID).
		Order("created_at DESC").
		Find(&auditLogs).Error
	return auditLogs, err
}

// GetAuditLogsByUser retrieves audit logs for a user
func (s *AuditService) GetAuditLogsByUser(userID string, limit int) ([]AuditLog, error) {
	var auditLogs []AuditLog
	err := s.db.Where("user_id = ? OR actor_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&auditLogs).Error
	return auditLogs, err
}
