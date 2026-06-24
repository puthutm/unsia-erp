package sharedidempotency

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrConcurrentRequest = errors.New("concurrent request in progress")
	ErrStoreNotSet       = errors.New("idempotency store is not registered")
)

type Store interface {
	CheckAndLock(ctx context.Context, module, key string, ttl time.Duration) (string, bool, error)
	SaveResponse(ctx context.Context, module, key, response string, ttl time.Duration) error
	SaveFailure(ctx context.Context, module, key, errMsg string) error
}

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

var (
	storeMutex    sync.RWMutex
	activeStore   Store
	defaultModule = "general"
)

// RegisterStore registers the active store
func RegisterStore(s Store) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	activeStore = s
}

// SetDefaultModule sets the default module name
func SetDefaultModule(mod string) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	defaultModule = mod
}

// splitKey splits a key like "pmb:handover:applicant-uuid" into ("pmb", "handover:applicant-uuid")
func splitKey(fullKey string) (string, string) {
	parts := strings.SplitN(fullKey, ":", 2)
	storeMutex.RLock()
	defMod := defaultModule
	storeMutex.RUnlock()

	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return defMod, fullKey
}

// CheckAndLock checks if the request is cached. If not, it locks the key.
func CheckAndLock(ctx context.Context, idempotencyKey string, ttl time.Duration) (string, bool, error) {
	storeMutex.RLock()
	store := activeStore
	storeMutex.RUnlock()

	if store == nil {
		return "", false, ErrStoreNotSet
	}

	module, key := splitKey(idempotencyKey)
	return store.CheckAndLock(ctx, module, key, ttl)
}

// SaveResponse saves the response to the store and unlocks the key
func SaveResponse(ctx context.Context, idempotencyKey string, response string, ttl time.Duration) error {
	storeMutex.RLock()
	store := activeStore
	storeMutex.RUnlock()

	if store == nil {
		return ErrStoreNotSet
	}

	module, key := splitKey(idempotencyKey)
	return store.SaveResponse(ctx, module, key, response, ttl)
}

// SaveFailure marks the key as failed so it can be retried immediately
func SaveFailure(ctx context.Context, idempotencyKey string, errMsg string) error {
	storeMutex.RLock()
	store := activeStore
	storeMutex.RUnlock()

	if store == nil {
		return ErrStoreNotSet
	}

	module, key := splitKey(idempotencyKey)
	return store.SaveFailure(ctx, module, key, errMsg)
}

// SQLStore implements Store interface using standard PostgreSQL database executor
type SQLStore struct {
	db      DBExecutor
	lockTTL time.Duration
}

func NewSQLStore(db DBExecutor, defaultLockTTL time.Duration) *SQLStore {
	if defaultLockTTL <= 0 {
		defaultLockTTL = 30 * time.Second // default lock TTL if request hangs
	}
	return &SQLStore{
		db:      db,
		lockTTL: defaultLockTTL,
	}
}

func (s *SQLStore) CheckAndLock(ctx context.Context, module, key string, ttl time.Duration) (string, bool, error) {
	now := time.Now()
	lockUntil := now.Add(s.lockTTL)
	expiresAt := now.Add(ttl)

	var status string
	var responseJSON sql.NullString
	var lockedUntil time.Time

	// Check existing key status
	querySelect := `
		SELECT status, response_json, locked_until 
		FROM idempotency_keys 
		WHERE module = $1 AND idempotency_key = $2
	`
	err := s.db.QueryRowContext(ctx, querySelect, module, key).Scan(&status, &responseJSON, &lockedUntil)

	if err == sql.ErrNoRows {
		// Insert new lock
		queryInsert := `
			INSERT INTO idempotency_keys (
				module, idempotency_key, status, locked_until, expires_at, created_at, updated_at
			) VALUES ($1, $2, 'processing', $3, $4, NOW(), NOW())
		`
		_, err = s.db.ExecContext(ctx, queryInsert, module, key, lockUntil, expiresAt)
		if err != nil {
			// Check for unique key violation (concurrency conflict)
			// PostgreSQL code 23505 is unique violation
			return "", false, ErrConcurrentRequest
		}
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}

	switch status {
	case "completed":
		return responseJSON.String, true, nil

	case "processing":
		if lockedUntil.After(now) {
			return "", false, ErrConcurrentRequest
		}
		// Fallthrough if lock has expired

	case "failed", "expired":
		// Lock expired or previous execution failed, try to reclaim it
	}

	// Update to claim/lock the record
	queryUpdate := `
		UPDATE idempotency_keys 
		SET status = 'processing', locked_until = $3, expires_at = $4, updated_at = NOW() 
		WHERE module = $1 AND idempotency_key = $2 
		  AND (status != 'processing' OR locked_until <= $5)
	`
	res, err := s.db.ExecContext(ctx, queryUpdate, module, key, lockUntil, expiresAt, now)
	if err != nil {
		return "", false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return "", false, err
	}
	if rows == 0 {
		return "", false, ErrConcurrentRequest
	}

	return "", false, nil
}

func (s *SQLStore) SaveResponse(ctx context.Context, module, key, response string, ttl time.Duration) error {
	query := `
		UPDATE idempotency_keys 
		SET status = 'completed', response_json = $3, completed_at = NOW(), expires_at = $4, updated_at = NOW() 
		WHERE module = $1 AND idempotency_key = $2
	`
	expiresAt := time.Now().Add(ttl)
	_, err := s.db.ExecContext(ctx, query, module, key, response, expiresAt)
	return err
}

func (s *SQLStore) SaveFailure(ctx context.Context, module, key, errMsg string) error {
	query := `
		UPDATE idempotency_keys 
		SET status = 'failed', last_error = $3, updated_at = NOW() 
		WHERE module = $1 AND idempotency_key = $2
	`
	_, err := s.db.ExecContext(ctx, query, module, key, errMsg)
	return err
}
