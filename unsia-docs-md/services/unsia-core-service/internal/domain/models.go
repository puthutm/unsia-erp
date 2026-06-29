package domain

import (
	"time"
)

type Person struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	NIM       string    `gorm:"column:nim"`
	NIP       string    `gorm:"column:nip"`
	Phone     string    `gorm:"column:phone"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Person) TableName() string {
	return "persons"
}

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	PersonID     string    `gorm:"column:person_id"`
	Person       Person    `gorm:"foreignKey:PersonID"`
	Username     string    `gorm:"column:username"`
	PasswordHash string    `gorm:"column:password_hash"`
	Status       string    `gorm:"column:status;default:'active'"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Role struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code      string    `gorm:"column:code"`
	Name      string    `gorm:"column:name"`
	ScopeType string    `gorm:"column:scope_type"` // global, prodi, module, self
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Code      string    `gorm:"column:code"` // e.g. pmb.applicant.verify
	Name      string    `gorm:"column:name"`
	Module    string    `gorm:"column:module"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

type UserRole struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	UserID         string    `gorm:"column:user_id"`
	User           User      `gorm:"foreignKey:UserID"`
	RoleID         string    `gorm:"column:role_id"`
	Role           Role      `gorm:"foreignKey:RoleID"`
	StudyProgramID *string   `gorm:"column:study_program_id"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type RolePermission struct {
	ID           string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	RoleID       string     `gorm:"column:role_id"`
	Role         Role       `gorm:"foreignKey:RoleID"`
	PermissionID string     `gorm:"column:permission_id"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

type Application struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicationCode string    `gorm:"column:application_code"`
	Name            string    `gorm:"column:name"`
	URL             string    `gorm:"column:url"`
	Enabled         bool      `gorm:"column:enabled;default:true"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (Application) TableName() string {
	return "applications"
}

type Session struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	UserID           string    `gorm:"column:user_id"`
	User             User      `gorm:"foreignKey:UserID"`
	TokenHash        string    `gorm:"column:token_hash;uniqueIndex"`
	RefreshTokenHash string    `gorm:"column:refresh_token_hash"`
	IsRevoked        bool      `gorm:"column:is_revoked;default:false"`
	ExpiresAt        time.Time `gorm:"column:expires_at"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (Session) TableName() string {
	return "sessions"
}

type ActiveRoleSession struct {
	ID             string      `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	UserID         string      `gorm:"column:user_id"`
	User           User        `gorm:"foreignKey:UserID"`
	RoleID         string      `gorm:"column:role_id"`
	Role           Role        `gorm:"foreignKey:RoleID"`
	SessionID      string      `gorm:"column:session_id"`
	Session        Session     `gorm:"foreignKey:SessionID"`
	ApplicationID  string      `gorm:"column:application_id"`
	Application    Application `gorm:"foreignKey:ApplicationID"`
	StudyProgramID *string     `gorm:"column:study_program_id"`
	CreatedAt      time.Time   `gorm:"column:created_at"`
}

func (ActiveRoleSession) TableName() string {
	return "active_role_sessions"
}

type OAuthClient struct {
	ID                string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicationID     string     `gorm:"column:application_id"`
	ClientID          string     `gorm:"column:client_id"`
	ClientSecretHash  *string    `gorm:"column:client_secret_hash"`
	ClientName        string     `gorm:"column:client_name"`
	ClientType        string     `gorm:"column:client_type;default:'confidential'"`
	GrantTypes        string     `gorm:"column:grant_types;type:jsonb"`
	AllowedScopes     string     `gorm:"column:allowed_scopes;type:jsonb"`
	Status            string     `gorm:"column:status;default:'PENDING'"`
	OwnerName         string     `gorm:"column:owner_name"`
	OwnerEmail        string     `gorm:"column:owner_email"`
	OwnerOrganization string     `gorm:"column:owner_organization"`
	ApprovedAt        *time.Time `gorm:"column:approved_at"`
	ApprovedBy        *string    `gorm:"column:approved_by"`
	SuspendedAt       *time.Time `gorm:"column:suspended_at"`
	RevokedAt         *time.Time `gorm:"column:revoked_at"`
	IsActive          bool       `gorm:"column:is_active;default:true"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`
}

func (OAuthClient) TableName() string {
	return "oauth_clients"
}

type RedirectURI struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	OAuthClientID string    `gorm:"column:oauth_client_id"`
	RedirectURI   string    `gorm:"column:redirect_uri"`
	IsActive      bool      `gorm:"column:is_active;default:true"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (RedirectURI) TableName() string {
	return "redirect_uris"
}

type OAuthAuthorizationCode struct {
	ID                  string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	CodeHash            string     `gorm:"column:code_hash"`
	ClientID            string     `gorm:"column:client_id"`
	UserID              string     `gorm:"column:user_id"`
	RedirectURI         string     `gorm:"column:redirect_uri"`
	Scope               string     `gorm:"column:scope"`
	CodeChallenge       string     `gorm:"column:code_challenge"`
	CodeChallengeMethod string     `gorm:"column:code_challenge_method;default:'S256'"`
	State               string     `gorm:"column:state"`
	ExpiresAt           time.Time  `gorm:"column:expires_at"`
	UsedAt              *time.Time `gorm:"column:used_at"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
}

func (OAuthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}

type ClientRegistrationRequest struct {
	ID                    string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	OAuthClientID         *string    `gorm:"column:oauth_client_id"`
	OwnerName             string     `gorm:"column:owner_name"`
	OwnerEmail            string     `gorm:"column:owner_email"`
	OwnerOrganization     string     `gorm:"column:owner_organization"`
	RequestedScopes       string     `gorm:"column:requested_scopes;type:jsonb"`
	RequestedGrantTypes   string     `gorm:"column:requested_grant_types;type:jsonb"`
	RequestedRedirectURIs string     `gorm:"column:requested_redirect_uris;type:jsonb"`
	Status                string     `gorm:"column:status;default:'PENDING'"`
	AdminNotes            string     `gorm:"column:admin_notes"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
	ReviewedAt            *time.Time `gorm:"column:reviewed_at"`
	ReviewedBy            *string    `gorm:"column:reviewed_by"`
}

func (ClientRegistrationRequest) TableName() string {
	return "client_registration_requests"
}

type ImpersonationSession struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ActorUserID   string     `gorm:"column:actor_user_id"`
	ActorUser     User       `gorm:"foreignKey:ActorUserID"`
	TargetUserID  string     `gorm:"column:target_user_id"`
	TargetUser    User       `gorm:"foreignKey:TargetUserID"`
	TargetRoleID  string     `gorm:"column:target_role_id"`
	TargetRole    Role       `gorm:"foreignKey:TargetRoleID"`
	ApplicationID *string    `gorm:"column:application_id"`
	Application   *Application `gorm:"foreignKey:ApplicationID"`
	SessionID     string     `gorm:"column:session_id"`
	Session       Session    `gorm:"foreignKey:SessionID"`
	Reason        string     `gorm:"column:reason"`
	StartedAt     time.Time  `gorm:"column:started_at;default:now()"`
	EndedAt       *time.Time `gorm:"column:ended_at"`
	ExpiredAt     time.Time  `gorm:"column:expired_at"`
	Status        string     `gorm:"column:status;default:'active'"`
}

func (ImpersonationSession) TableName() string {
	return "impersonation_sessions"
}

type ServiceToken struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	ApplicationID string    `gorm:"column:application_id"`
	TokenHash     string    `gorm:"column:token_hash"`
	Scopes        string    `gorm:"column:scopes;type:jsonb"`
	ExpiresAt     time.Time `gorm:"column:expired_at"`
	RevokedAt     *time.Time `gorm:"column:revoked_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (ServiceToken) TableName() string {
	return "service_tokens"
}

type AuditLog struct {
	ID                     string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	UserID                *string    `gorm:"column:user_id"`
	ActorUserID           *string    `gorm:"column:actor_user_id"`
	TargetUserID          *string    `gorm:"column:target_user_id"`
	ActiveRoleID          *string    `gorm:"column:active_role_id"`
	ImpersonationSessionID *string    `gorm:"column:impersonation_session_id"`
	ApplicationID         *string    `gorm:"column:application_id"`
	Module                string     `gorm:"column:module"`
	Action                string     `gorm:"column:action"`
	EntityName            *string    `gorm:"column:entity_name"`
	EntityID              *string    `gorm:"column:entity_id"`
	Reason                *string    `gorm:"column:reason"`
	OldValue              *string    `gorm:"column:old_value;type:jsonb"`
	NewValue              *string    `gorm:"column:new_value;type:jsonb"`
	RequestID             *string    `gorm:"column:request_id"`
	IPAddress             *string    `gorm:"column:ip_address"`
	UserAgent             *string    `gorm:"column:user_agent"`
	CreatedAt             time.Time  `gorm:"column:created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

