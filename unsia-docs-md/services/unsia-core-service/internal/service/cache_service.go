package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CacheService struct {
	client *redis.Client
	db     *gorm.DB
}

func NewCacheService(client *redis.Client, db *gorm.DB) *CacheService {
	return &CacheService{
		client: client,
		db:     db,
	}
}

// CacheOptions for cache operations
type CacheOptions struct {
	Expiration time.Duration
	Prefix    string
}

var defaultCacheOptions = CacheOptions{
	Expiration: 1 * time.Hour,
	Prefix:    "unsia",
}

// Get gets value from cache
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Set sets value to cache
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, opts ...CacheOption) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	options := defaultCacheOptions
	for _, opt := range opts {
		opt(&options)
	}

	return s.client.Set(ctx, key, data, options.Expiration).Err()
}

// Delete deletes key from cache
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// GetOrFetch gets from cache or fetches from database
func (s *CacheService) GetOrFetch(ctx context.Context, key string, dest interface{}, fetchFunc func() error) error {
	err := s.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	// Cache miss, fetch from database
	if err := fetchFunc(); err != nil {
		return err
	}

	// Set to cache
	return s.Set(ctx, key, dest)
}

// InvalidateByPattern invalidates keys by pattern
func (s *CacheService) InvalidateByPattern(ctx context.Context, pattern string) error {
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return s.client.Del(ctx, keys...).Err()
	}
	return nil
}

// CacheOption is a functional option
type CacheOption func(*CacheOptions)

// WithExpiration sets expiration
func WithExpiration(exp time.Duration) CacheOption {
	return func(o *CacheOptions) {
		o.Expiration = exp
	}
}

// WithPrefix sets prefix
func WithPrefix(prefix string) CacheOption {
	return func(o *CacheOptions) {
		o.Prefix = prefix
	}
}

// UserCacheService handles user caching
type UserCacheService struct {
	CacheService *CacheService
}

func NewUserCacheService(client *redis.Client, db *gorm.DB) *UserCacheService {
	return &UserCacheService{
		CacheService: NewCacheService(client, db),
	}
}

// CachedUser represents cached user
type CachedUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name     string `json:"name"`
	RoleID   string `json:"role_id"`
	Active   bool   `json:"active"`
	CachedAt int64  `json:"cached_at"`
}

// GetUser gets cached user
func (s *UserCacheService) GetUser(ctx context.Context, userID string) (*CachedUser, error) {
	key := fmt.Sprintf("user:%s", userID)
	var user CachedUser
	err := s.CacheService.Get(ctx, key, &user)
	return &user, err
}

// SetUser caches user
func (s *UserCacheService) SetUser(ctx context.Context, user *CachedUser) error {
	key := fmt.Sprintf("user:%s", user.ID)
	return s.CacheService.Set(ctx, key, user)
}

// InvalidateUser invalidates user cache
func (s *UserCacheService) InvalidateUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:%s", userID)
	return s.CacheService.Delete(ctx, key)
}

// RoleCacheService handles role caching
type RoleCacheService struct {
	CacheService *CacheService
}

func NewRoleCacheService(client *redis.Client, db *gorm.DB) *RoleCacheService {
	return &RoleCacheService{
		CacheService: NewCacheService(client, db),
	}
}

// GetRole gets cached role
func (s *RoleCacheService) GetRole(ctx context.Context, roleID string) (map[string]interface{}, error) {
	key := fmt.Sprintf("role:%s", roleID)
	var role map[string]interface{}
	err := s.CacheService.Get(ctx, key, &role)
	return role, err
}

// SetRole caches role
func (s *RoleCacheService) SetRole(ctx context.Context, roleID string, role map[string]interface{}) error {
	key := fmt.Sprintf("role:%s", roleID)
	return s.CacheService.Set(ctx, key, role, WithExpiration(24*time.Hour))
}

// InvalidateRole invalidates role cache
func (s *RoleCacheService) InvalidateRole(ctx context.Context, roleID string) error {
	key := fmt.Sprintf("role:%s", roleID)
	return s.CacheService.Delete(ctx, key)
}
