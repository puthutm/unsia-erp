package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

// OutboxEvent represents an event in the outbox table
type OutboxEvent struct {
	ID              string    `gorm:"primaryKey;column:id"`
	EventName       string    `gorm:"column:event_name"`
	EventVersion   string    `gorm:"column:event_version"`
	EventKey       string    `gorm:"column:event_key;uniqueIndex"`
	PublisherService string  `gorm:"column:publisher_service"`
	AggregateType string   `gorm:"column:aggregate_type"`
	AggregateID   string   `gorm:"column:aggregate_id"`
	CorrelationID string   `gorm:"column:correlation_id"`
	CausationID    string   `gorm:"column:causation_id"`
	Payload       string    `gorm:"type:jsonb;column:payload"`
	Status        string   `gorm:"column:status"` // pending, published, failed
	RetryCount    int      `gorm:"column:retry_count;default:0"`
	ErrorMessage *string  `gorm:"column:error_message"`
	PublishedAt   *time.Time `gorm:"column:published_at"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}

// OutboxPoller polls and publishes outbox events to RabbitMQ
type OutboxPoller struct {
	db          *gorm.DB
	amqpConn   *amqp.Connection
	amqpChan   *amqp.Channel
	exchange   string
	queue      string
	interval   time.Duration
	maxRetries int
	stopCh     chan struct{}
}

// NewOutboxPoller creates a new outbox poller
func NewOutboxPoller(db *gorm.DB, amqpURL string, interval time.Duration) (*OutboxPoller, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		"unsia.events", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	maxRetries := 3
	if v := os.Getenv("OUTBOX_MAX_RETRIES"); v != "" {
		fmt.Sscanf(v, "%d", &maxRetries)
	}

	return &OutboxPoller{
		db:        db,
		amqpConn:  conn,
		amqpChan:  ch,
		exchange: "unsia.events",
		queue:    "finance.outbox",
		interval: interval,
		maxRetries: maxRetries,
		stopCh:   make(chan struct{}),
	}, nil
}

// Start starts the outbox poller
func (p *OutboxPoller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	log.Info().Str("interval", p.interval.String()).Msg("Outbox poller started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Outbox poller stopping (context done)")
			return
		case <-p.stopCh:
			log.Info().Msg("Outbox poller stopping")
			return
		case <-ticker.C:
			if err := p.pollAndPublish(ctx); err != nil {
				log.Error().Err(err).Msg("Error polling outbox")
			}
		}
	}
}

// Stop stops the outbox poller
func (p *OutboxPoller) Stop() {
	close(p.stopCh)
	if p.amqpChan != nil {
		p.amqpChan.Close()
	}
	if p.amqpConn != nil {
		p.amqpConn.Close()
	}
}

// pollAndPublish polls and publishes pending events
func (p *OutboxPoller) pollAndPublish(ctx context.Context) error {
	var events []OutboxEvent
	err := p.db.WithContext(ctx).
		Where("status = ? AND retry_count < ?", "pending", p.maxRetries).
		Order("created_at ASC").
		Limit(100).
		Find(&events).Error

	if err != nil {
		return fmt.Errorf("failed to fetch events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	log.Info().Int("count", len(events)).Msg("Processing outbox events")

	for _, event := range events {
		if err := p.publishEvent(ctx, &event); err != nil {
			log.Error().Err(err).Str("event_id", event.ID).Msg("Failed to publish event")
			
			// Update retry count
			retryCount := event.RetryCount + 1
			errMsg := err.Error()
			p.db.Model(&event).Updates(map[string]interface{}{
				"retry_count":   retryCount,
				"error_message": &errMsg,
			})
			continue
		}

		// Mark as published
		now := time.Now()
		p.db.Model(&event).Updates(map[string]interface{}{
			"status":        "published",
			"published_at":  &now,
		})

		log.Info().Str("event_id", event.ID).Str("event_name", event.EventName).Msg("Event published")
	}

	return nil
}

// publishEvent publishes a single event to RabbitMQ
func (p *OutboxPoller) publishEvent(ctx context.Context, event *OutboxEvent) error {
	// Parse payload
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Create message
	message := map[string]interface{}{
		"event_name":         event.EventName,
		"event_version":     event.EventVersion,
		"event_key":         event.EventKey,
		"publisher_service": event.PublisherService,
		"aggregate_type":    event.AggregateType,
		"aggregate_id":      event.AggregateID,
		"correlation_id":    event.CorrelationID,
		"causation_id":      event.CausationID,
		"occurred_at":       event.CreatedAt,
		"payload":           payload,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Routing key based on event name
	routingKey := strings.Replace(event.EventName, ".", ".", -1)

	// Publish to exchange
	err = p.amqpChan.PublishWithContext(ctx,
		p.exchange,  // exchange
		routingKey,  // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Headers: amqp.Table{
				"event_name":       event.EventName,
				"event_version":   event.EventVersion,
				"event_key":       event.EventKey,
				"publisher":       event.PublisherService,
				"correlation_id":   event.CorrelationID,
				"aggregate_type":  event.AggregateType,
				"aggregate_id":   event.AggregateID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// ProcessDeadLetter processes dead letter events
func (p *OutboxPoller) ProcessDeadLetter(ctx context.Context) error {
	var events []OutboxEvent
	err := p.db.WithContext(ctx).
		Where("status = ? AND retry_count >= ?", "pending", p.maxRetries).
		Order("created_at ASC").
		Limit(50).
		Find(&events).Error

	if err != nil {
		return fmt.Errorf("failed to fetch dead letters: %w", err)
	}

	for _, event := range events {
		log.Warn().Str("event_id", event.ID).Str("event_name", event.EventName).
			Int("retry_count", event.RetryCount).Msg("Marking event as dead letter")

		p.db.Model(&event).Update("status", "dead_letter")
	}

	return nil
}

// GetPendingCount returns the count of pending events
func (p *OutboxPoller) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := p.db.WithContext(ctx).
		Model(&OutboxEvent{}).
		Where("status = ?", "pending").
		Count(&count)

	return count, err
}
