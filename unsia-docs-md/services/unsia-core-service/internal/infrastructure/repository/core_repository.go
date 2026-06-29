package repository

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/unsia-erp/unsia-core-service/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").Joins("JOIN persons ON persons.id = users.person_id").Where("persons.email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByNIM(nim string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").Joins("JOIN persons ON persons.id = users.person_id").Where("persons.nim = ?", nim).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByNIP(nip string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").Joins("JOIN persons ON persons.id = users.person_id").Where("persons.nip = ?", nip).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Person").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserRoles(userID string) ([]domain.UserRole, error) {
	var userRoles []domain.UserRole
	err := r.db.Preload("Role").Where("user_id = ?", userID).Find(&userRoles).Error
	return userRoles, err
}

func (r *UserRepository) GetRolePermissions(roleID string) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.Table("permissions").
		Select("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(session *domain.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) CreateImpersonationSession(session *domain.ImpersonationSession) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) GetByRefreshToken(token string) (*domain.Session, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	var session domain.Session
	err := r.db.Preload("User").Preload("User.Person").Where("refresh_token_hash = ? AND is_revoked = false", tokenHash).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) RevokeSession(sessionID string) error {
	return r.db.Model(&domain.Session{}).Where("id = ?", sessionID).Update("is_revoked", true).Error
}

func (r *SessionRepository) SaveActiveRoleSession(active *domain.ActiveRoleSession) error {
	// Delete any existing active role session for this session & user
	r.db.Where("session_id = ? AND user_id = ?", active.SessionID, active.UserID).Delete(&domain.ActiveRoleSession{})
	return r.db.Create(active).Error
}

func (r *SessionRepository) GetActiveRoleSession(sessionID string) (*domain.ActiveRoleSession, error) {
	var active domain.ActiveRoleSession
	err := r.db.Preload("Role").Preload("Application").Where("session_id = ?", sessionID).First(&active).Error
	if err != nil {
		return nil, err
	}
	return &active, nil
}

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) GetAllEnabled() ([]domain.Application, error) {
	var apps []domain.Application
	err := r.db.Where("enabled = true").Find(&apps).Error
	return apps, err
}

func (r *ApplicationRepository) GetByCode(code string) (*domain.Application, error) {
	var app domain.Application
	err := r.db.Where("application_code = ? AND enabled = true", code).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationRepository) GetApplicationByID(id string) (*domain.Application, error) {
	var app domain.Application
	err := r.db.Where("id = ?", id).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetByCode(code string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetAll() ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(auditLog *domain.AuditLog) error {
	return r.db.Create(auditLog).Error
}

func (r *AuditLogRepository) GetByID(id string) (*domain.AuditLog, error) {
	var auditLog domain.AuditLog
	err := r.db.First(&auditLog, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &auditLog, nil
}

func (r *AuditLogRepository) FindByEntity(entityName, entityID string, limit int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.Where("entity_name = ? AND entity_id = ?", entityName, entityID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) FindByUser(userID string, limit int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.Where("user_id = ? OR actor_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

