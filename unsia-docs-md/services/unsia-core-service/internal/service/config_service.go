package service

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppConfig struct {
	ID       string `gorm:"primaryKey;column:id"`
	Key      string `gorm:"uniqueIndex;column:key"`
	Value    string `gorm:"type:jsonb;column:value"`
	Type     string `gorm:"type:varchar(50);column:type"`
	IsPublic bool   `gorm:"default:false;column:is_public"`
}

func (AppConfig) TableName() string {
	return "app_configs"
}

// ConfigService manages application configuration
type ConfigService struct {
	db      *gorm.DB
	configs map[string]interface{}
	mu      sync.RWMutex
}

func NewConfigService(db *gorm.DB) (*ConfigService, error) {
	s := &ConfigService{
		db:      db,
		configs: make(map[string]interface{}),
	}

	// Load configs from database
	if err := s.loadConfigs(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ConfigService) loadConfigs() error {

	var configs []AppConfig
	if err := s.db.Find(&configs).Error; err != nil {
		return err
	}

	for _, c := range configs {
		var value interface{}
		json.Unmarshal([]byte(c.Value), &value)
		s.configs[c.Key] = value
	}

	return nil
}

// Get returns config value
func (s *ConfigService) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.configs[key]
	return val, ok
}

// GetString returns config as string
func (s *ConfigService) GetString(key string, defaultValue string) string {
	val, ok := s.configs[key]
	if !ok {
		return defaultValue
	}
	if str, ok := val.(string); ok {
		return str
	}
	return defaultValue
}

// GetInt returns config as int
func (s *ConfigService) GetInt(key string, defaultValue int) int {
	val, ok := s.configs[key]
	if !ok {
		return defaultValue
	}
	if n, ok := val.(float64); ok {
		return int(n)
	}
	return defaultValue
}

// GetBool returns config as bool
func (s *ConfigService) GetBool(key string, defaultValue bool) bool {
	val, ok := s.configs[key]
	if !ok {
		return defaultValue
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return defaultValue
}

// Set sets config value
func (s *ConfigService) Set(key string, value interface{}) error {

	jsonValue, _ := json.Marshal(value)
	config := AppConfig{
		Key:   key,
		Value: string(jsonValue),
		Type:  fmt.Sprintf("%T", value),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.db.Where("key = ?", key).First(&AppConfig{}).Error; err == nil {
		return s.db.Model(&AppConfig{}).Where("key = ?", key).Update("value", string(jsonValue)).Error
	}

	config.ID = uuid.New().String()
	return s.db.Create(&config).Error
}

// Delete deletes config
func (s *ConfigService) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.configs, key)
	return s.db.Where("key = ?", key).Delete(&AppConfig{}).Error
}

// FeatureFlags manages feature flags
type FeatureFlagService struct {
	db        *gorm.DB
	flags     map[string]bool
	flagCache map[string]FeatureFlag
	mu       sync.RWMutex
}

type FeatureFlag struct {
	ID          string    `gorm:"primaryKey"`
	Key         string    `gorm:"uniqueIndex"`
	Description string    `gorm:"column:description"`
	IsEnabled   bool      `gorm:"column:is_enabled"`
	Rollout    int       `gorm:"default:0"` // percentage 0-100
	TargetUser *string   `gorm:"column:target_user"`
	ValidFrom  *string   `gorm:"column:valid_from"`
	ValidTo    *string   `gorm:"column:valid_to"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (FeatureFlag) TableName() string {
	return "feature_flags"
}

func NewFeatureFlagService(db *gorm.DB) *FeatureFlagService {
	s := &FeatureFlagService{
		db:        db,
		flags:     make(map[string]bool),
		flagCache: make(map[string]FeatureFlag),
	}
	s.loadFlags()
	return s
}

func (s *FeatureFlagService) loadFlags() {
	var flags []FeatureFlag
	s.db.Find(&flags)

	for _, f := range flags {
		s.flags[f.Key] = f.IsEnabled
		s.flagCache[f.Key] = f
	}
}

// IsEnabled checks if feature is enabled
func (s *FeatureFlagService) IsEnabled(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	enabled, ok := s.flags[key]
	return ok && enabled
}

// IsEnabledForUser checks if feature is enabled for a specific user
func (s *FeatureFlagService) IsEnabledForUser(key string, userID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	flag, ok := s.flagCache[key]
	if !ok || !flag.IsEnabled {
		return false
	}

	// Check user-specific flag
	if flag.TargetUser != nil && *flag.TargetUser == userID {
		return true
	}

	// Check rollout percentage
	if flag.Rollout > 0 {
		// Simple hash-based rollout
		hash := simpleHash(userID + key)
		if hash%100 < flag.Rollout {
			return true
		}
	}

	return false
}

// Enable enables a feature
func (s *FeatureFlagService) Enable(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.flags[key] = true
	if flag, ok := s.flagCache[key]; ok {
		flag.IsEnabled = true
		return s.db.Model(&FeatureFlag{}).Where("key = ?", key).Update("is_enabled", true).Error
	}

	// Create new flag
	flag := FeatureFlag{
		ID:        uuid.New().String(),
		Key:       key,
		IsEnabled: true,
	}
	s.flagCache[key] = flag
	return s.db.Create(&flag).Error
}

// Disable disables a feature
func (s *FeatureFlagService) Disable(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.flags[key] = false
	if flag, ok := s.flagCache[key]; ok {
		flag.IsEnabled = false
		return s.db.Model(&FeatureFlag{}).Where("key = ?", key).Update("is_enabled", false).Error
	}
	return nil
}

// EnvConfig loads config from environment variables
type EnvConfig struct {
	// Server
	ServerPort string
	ServerHost string

	// Database
	DatabaseURL string

	// Redis
	RedisAddr string

	// JWT
	JWTSecret string

	// SMTP
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string

	// App
	AppName    string
	AppEnv    string
	AppDebug  bool
}

func LoadEnvConfig() *EnvConfig {
	return &EnvConfig{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/core_db"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:  getEnv("JWT_SECRET", "changeme"),
		SMTPHost:  getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:  getEnvInt("SMTP_PORT", 587),
		SMTPUser:  getEnv("SMTP_USER", ""),
		SMTPPass:  getEnv("SMTP_PASS", ""),
		AppName:   getEnv("APP_NAME", "UNSIA Core"),
		AppEnv:    getEnv("APP_ENV", "development"),
		AppDebug:  getEnvBool("APP_DEBUG", false),
	}
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		var n int
		fmt.Sscanf(val, "%d", &n)
		return n
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		return val == "true" || val == "1"
	}
	return defaultValue
}

func simpleHash(s string) int {
	h := 0
	for _, c := range s {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}
