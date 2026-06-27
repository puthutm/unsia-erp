package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExternalAppService struct {
	db *gorm.DB
}

func NewExternalAppService(db *gorm.DB) *ExternalAppService {
	return &ExternalAppService{db: db}
}

type ExternalApp struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"column:name"`
	Slug       string    `gorm:"column:slug;uniqueIndex"`
	Type       string    `gorm:"column:type"` // web, mobile, desktop, webhook, cron
	URL        *string   `gorm:"column:url"`
	CallbackURL *string `gorm:"column:callback_url"`
	LogoURL    *string   `gorm:"column:logo_url"`
	Description string  `gorm:"column:description"`
	ClientID   string    `gorm:"column:client_id;uniqueIndex"`
	ClientSecret string  `gorm:"column:client_secret"`
	IsActive   bool      `gorm:"column:is_active"`
	IsInternal bool      `gorm:"column:is_internal"`
	Scopes     string    `gorm:"column:scopes;type:jsonb"`
	IPWhitelist string  `gorm:"column:ip_whitelist;type:jsonb"`
	RateLimit  int       `gorm:"column:rate_limit"` // requests per minute
	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	ExpiredAt  *time.Time `gorm:"column:expired_at"`
	CreatedBy  string    `gorm:"column:created_by"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ExternalApp) TableName() string {
	return "external_apps"
}

// CreateExternalApp creates a new external app
func (s *ExternalAppService) CreateExternalApp(input CreateExternalAppInput) (*ExternalApp, error) {
	clientID := generateClientID()
	clientSecret := generateClientSecret()

	app := ExternalApp{
		ID:            uuid.New().String(),
		Name:          input.Name,
		Slug:          input.Slug,
		Type:          input.Type,
		URL:           input.URL,
		CallbackURL:   input.CallbackURL,
		Description:  input.Description,
		ClientID:      clientID,
		ClientSecret:  hashSecret(clientSecret),
		IsActive:     true,
		IsInternal:   input.IsInternal,
		Scopes:       serializeScopes(input.Scopes),
		IPWhitelist:  serializeIPs(input.IPWhitelist),
		RateLimit:   input.RateLimit,
		CreatedBy:   input.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(&app).Error; err != nil {
		return nil, err
	}

	// Return with unhashed secret only once
	app.ClientSecret = clientSecret
	return &app, nil
}

// GetExternalApp gets app by ID
func (s *ExternalAppService) GetExternalApp(id string) (*ExternalApp, error) {
	var app ExternalApp
	if err := s.db.Where("id = ?", id).First(&app).Error; err != nil {
		return nil, err
	}
	app.ClientSecret = "" // Never return secret
	return &app, nil
}

// GetExternalAppBySlug gets app by slug
func (s *ExternalAppService) GetExternalAppBySlug(slug string) (*ExternalApp, error) {
	var app ExternalApp
	if err := s.db.Where("slug = ?", slug).First(&app).Error; err != nil {
		return nil, err
	}
	app.ClientSecret = ""
	return &app, nil
}

// GetExternalAppByClientID gets app by client ID
func (s *ExternalAppService) GetExternalAppByClientID(clientID string) (*ExternalApp, error) {
	var app ExternalApp
	if err := s.db.Where("client_id = ?", clientID).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// ListExternalApps lists external apps
func (s *ExternalAppService) ListExternalApps(filter ExternalAppFilter) ([]ExternalApp, int64, error) {
	var apps []ExternalApp
	var total int64

	query := s.db.Model(&ExternalApp{})
	
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.IsInternal != nil {
		query = query.Where("is_internal = ?", *filter.IsInternal)
	}

	query.Count(&total)

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	query.Order("created_at DESC").Find(&apps)

	// Clear secrets
	for i := range apps {
		apps[i].ClientSecret = ""
	}

	return apps, total, nil
}

// UpdateExternalApp updates external app
func (s *ExternalAppService) UpdateExternalApp(id string, input UpdateExternalAppInput) (*ExternalApp, error) {
	var app ExternalApp
	if err := s.db.Where("id = ?", id).First(&app).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.Name != nil && *input.Name != "" {
		updates["name"] = *input.Name
	}
	if input.Description != nil && *input.Description != "" {
		updates["description"] = *input.Description
	}
	if input.URL != nil {
		updates["url"] = *input.URL
	}
	if input.CallbackURL != nil {
		updates["callback_url"] = *input.CallbackURL
	}
	if input.LogoURL != nil {
		updates["logo_url"] = *input.LogoURL
	}
	if input.Scopes != nil {
		updates["scopes"] = serializeScopes(input.Scopes)
	}
	if input.IPWhitelist != nil {
		updates["ip_whitelist"] = serializeIPs(input.IPWhitelist)
	}
	if input.RateLimit != nil && *input.RateLimit > 0 {
		updates["rate_limit"] = *input.RateLimit
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if err := s.db.Model(&ExternalApp{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetExternalApp(id)
}

// RegenerateSecret regenerates client secret
func (s *ExternalAppService) RegenerateSecret(id string) (string, error) {
	newSecret := generateClientSecret()

	if err := s.db.Model(&ExternalApp{}).Where("id = ?", id).Update("client_secret", hashSecret(newSecret)).Error; err != nil {
		return "", err
	}

	return newSecret, nil
}

// DeactivateExternalApp deactivates external app
func (s *ExternalAppService) DeactivateExternalApp(id string) error {
	return s.db.Model(&ExternalApp{}).Where("id = ?", id).Update("is_active", false).Error
}

// ActivateExternalApp activates external app
func (s *ExternalAppService) ActivateExternalApp(id string) error {
	return s.db.Model(&ExternalApp{}).Where("id = ?", id).Update("is_active", true).Error
}

// ValidateClientCredentials validates client ID and secret
func (s *ExternalAppService) ValidateClientCredentials(clientID, clientSecret string) (*ExternalApp, error) {
	var app ExternalApp
	if err := s.db.Where("client_id = ?", clientID).First(&app).Error; err != nil {
		return nil, err
	}

	if app.ClientSecret != hashSecret(clientSecret) {
		return nil, nil // Invalid credentials
	}

	if !app.IsActive {
		return nil, nil // App is inactive
	}

	if app.ExpiredAt != nil && app.ExpiredAt.Before(time.Now()) {
		return nil, nil // App has expired
	}

	// Update last login
	s.db.Model(&ExternalApp{}).Where("id = ?", app.ID).Update("last_login_at", time.Now())

	app.ClientSecret = ""
	return &app, nil
}

// CheckScope checks if app has required scope
func (s *ExternalAppService) CheckScope(app *ExternalApp, requiredScope string) bool {
	scopes := deserializeScopes(app.Scopes)
	for _, scope := range scopes {
		if scope == requiredScope || scope == "*" {
			return true
		}
	}
	return false
}

// CheckIP checks if IP is allowed
func (s *ExternalAppService) CheckIP(app *ExternalApp, ip string) bool {
	ips := deserializeIPs(app.IPWhitelist)
	if len(ips) == 0 {
		return true // No whitelist = allow all
	}

	for _, allowedIP := range ips {
		if allowedIP == ip {
			return true
		}
	}
	return false
}

// Input types
type CreateExternalAppInput struct {
	Name          string   `json:"name" binding:"required"`
	Slug         string   `json:"slug" binding:"required"`
	Type         string   `json:"type" binding:"required"`
	URL          *string  `json:"url"`
	CallbackURL  *string  `json:"callback_url"`
	LogoURL      *string  `json:"logo_url"`
	Description  string   `json:"description"`
	IsInternal   bool     `json:"is_internal"`
	Scopes      []string `json:"scopes"`
	IPWhitelist []string `json:"ip_whitelist"`
	RateLimit   int      `json:"rate_limit"`
	CreatedBy   string   `json:"-"`
}

type UpdateExternalAppInput struct {
	Name          *string  `json:"name"`
	Description  *string  `json:"description"`
	URL          *string  `json:"url"`
	CallbackURL  *string  `json:"callback_url"`
	LogoURL      *string  `json:"logo_url"`
	Scopes       []string `json:"scopes"`
	IPWhitelist  []string `json:"ip_whitelist"`
	RateLimit   *int     `json:"rate_limit"`
	IsActive     *bool   `json:"is_active"`
}

type ExternalAppFilter struct {
	Type       string
	IsActive   *bool
	IsInternal *bool
	Limit     int
	Offset    int
}

// Helper functions
func generateClientID() string {
	return "app_" + uuid.New().String()[:16]
}

func generateClientSecret() string {
	return uuid.New().String() + uuid.New().String()
}

func hashSecret(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}

func serializeScopes(scopes []string) string {
	if len(scopes) == 0 {
		return "[]"
	}
	result := "["
	for i, s := range scopes {
		if i > 0 {
			result += ","
		}
		result += `"` + s + `"`
	}
	result += "]"
	return result
}

func deserializeScopes(scopes string) []string {
	if scopes == "" || scopes == "[]" {
		return []string{}
	}
	// Simple parsing
	var result []string
	current := ""
	inQuote := false
	for _, c := range scopes {
		if c == '"' {
			inQuote = !inQuote
		} else if c == ',' && !inQuote {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else if inQuote {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func serializeIPs(ips []string) string {
	return serializeScopes(ips)
}

func deserializeIPs(ips string) []string {
	return deserializeScopes(ips)
}
