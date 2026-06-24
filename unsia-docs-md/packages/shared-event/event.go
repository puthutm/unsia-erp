package sharedevent

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type EventEnvelope struct {
	ID               string                 `json:"id"`
	EventName        string                 `json:"event_name"`
	EventVersion     string                 `json:"event_version"`
	EventKey         string                 `json:"event_key"`
	PublisherService string                 `json:"publisher_service"`
	AggregateType    string                 `json:"aggregate_type"`
	AggregateID      string                 `json:"aggregate_id"`
	CorrelationID    string                 `json:"correlation_id"`
	CausationID      string                 `json:"causation_id"`
	OccurredAt       time.Time              `json:"occurred_at"`
	Payload          interface{}            `json:"payload"`
}

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// BuildEventKey creates a deterministic event key
func BuildEventKey(name, id, version string) string {
	if version == "" {
		version = "v1"
	}
	return fmt.Sprintf("%s:%s:%s", name, id, version)
}

// WriteOutbox writes a domain/integration event to the outbox table in the same db transaction
func WriteOutbox(ctx context.Context, db DBExecutor, event EventEnvelope, eventType string) (string, error) {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}
	if event.EventVersion == "" {
		event.EventVersion = "v1"
	}
	if event.EventKey == "" {
		event.EventKey = BuildEventKey(event.EventName, event.AggregateID, event.EventVersion)
	}
	if eventType == "" {
		eventType = "INTEGRATION_EVENT"
	}

	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadHash := fmt.Sprintf("%x", sha256.Sum256(payloadBytes))

	var id string
	query := `
		INSERT INTO outbox_events (
			event_name, event_version, event_key, event_type, aggregate_type, aggregate_id,
			payload, payload_hash, correlation_id, causation_id, status, occurred_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6::uuid, $7, $8, $9, $10, 'PENDING', $11, NOW(), NOW()
		) RETURNING id
	`

	err = db.QueryRowContext(ctx, query,
		event.EventName,
		event.EventVersion,
		event.EventKey,
		eventType,
		event.AggregateType,
		event.AggregateID,
		payloadBytes,
		payloadHash,
		event.CorrelationID,
		event.CausationID,
		event.OccurredAt,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to insert outbox event: %w", err)
	}

	return id, nil
}

// ConsumeInbox checks if the incoming event is a duplicate.
// If it is new, it records it in `inbox_events` with status RECEIVED and returns (inboxID, true, nil).
// If it is a duplicate, it returns ("", false, nil).
func ConsumeInbox(ctx context.Context, db DBExecutor, consumerModule string, event EventEnvelope) (string, bool, error) {
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}

	payloadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		return "", false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadHash := fmt.Sprintf("%x", sha256.Sum256(payloadBytes))

	query := `
		INSERT INTO inbox_events (
			event_name, event_version, event_key, publisher_module, consumer_module,
			aggregate_type, aggregate_id, payload, payload_hash, correlation_id, causation_id,
			status, received_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7::uuid, $8, $9, $10, $11, 'RECEIVED', NOW(), NOW(), NOW()
		)
		ON CONFLICT (consumer_module, event_key) DO NOTHING
		RETURNING id
	`

	var id string
	err = db.QueryRowContext(ctx, query,
		event.EventName,
		event.EventVersion,
		event.EventKey,
		event.PublisherService,
		consumerModule,
		event.AggregateType,
		event.AggregateID,
		payloadBytes,
		payloadHash,
		event.CorrelationID,
		event.CausationID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		// Event already exists for this consumer module (ON CONFLICT DO NOTHING returned no row)
		return "", false, nil
	} else if err != nil {
		return "", false, fmt.Errorf("failed to record inbox event: %w", err)
	}

	return id, true, nil
}

// MarkPublished updates outbox event status to PUBLISHED
func MarkPublished(ctx context.Context, db DBExecutor, outboxID string) error {
	query := `
		UPDATE outbox_events 
		SET status = 'PUBLISHED', published_at = NOW(), updated_at = NOW() 
		WHERE id = $1::uuid
	`
	_, err := db.ExecContext(ctx, query, outboxID)
	return err
}

// MarkProcessed updates inbox event status to PROCESSED
func MarkProcessed(ctx context.Context, db DBExecutor, inboxID string) error {
	query := `
		UPDATE inbox_events 
		SET status = 'PROCESSED', processed_at = NOW(), updated_at = NOW() 
		WHERE id = $1::uuid
	`
	_, err := db.ExecContext(ctx, query, inboxID)
	return err
}

// MarkIgnored marks inbox as duplicate/ignored
func MarkIgnored(ctx context.Context, db DBExecutor, inboxID string) error {
	query := `
		UPDATE inbox_events 
		SET status = 'IGNORED_DUPLICATE', updated_at = NOW() 
		WHERE id = $1::uuid
	`
	_, err := db.ExecContext(ctx, query, inboxID)
	return err
}

// SendToDLQ moves the outbox/inbox event into the DLQ state due to continuous failures
func SendToDLQ(ctx context.Context, db DBExecutor, isOutbox bool, id string, lastError string) error {
	tableName := "inbox_events"
	if isOutbox {
		tableName = "outbox_events"
	}

	query := fmt.Sprintf(`
		UPDATE %s 
		SET status = 'DLQ', dead_letter_at = NOW(), last_error = $2, updated_at = NOW() 
		WHERE id = $1::uuid
	`, tableName)

	_, err := db.ExecContext(ctx, query, id, lastError)
	return err
}
