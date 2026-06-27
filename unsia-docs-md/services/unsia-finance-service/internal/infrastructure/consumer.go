package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"github.com/unsia-erp/unsia-finance-service/internal/domain"
)

// EventConsumer consumes inbound events from RabbitMQ
type EventConsumer struct {
	db        *gorm.DB
	amqpConn *amqp.Connection
	amqpChan *amqp.Channel
	queue    string
	stopCh   chan struct{}
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(db *gorm.DB, amqpURL string) (*EventConsumer, error) {
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

	// Declare queue
	_, err = ch.QueueDeclare(
		"finance.events", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queues to exchange
	bindings := []string{
		"pmb.applicant_created",
		"academic.student_created",
		"academic.student_updated",
	}

	for _, binding := range bindings {
		err = ch.QueueBind(
			"finance.events", // queue
			binding,         // routing key
			"unsia.events",  // exchange
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	return &EventConsumer{
		db:        db,
		amqpConn:  conn,
		amqpChan:  ch,
		queue:    "finance.events",
		stopCh:   make(chan struct{}),
	}, nil
}

// Start starts consuming events
func (c *EventConsumer) Start(ctx context.Context) error {
	msgs, err := c.amqpChan.Consume(
		c.queue,    // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Info().Str("queue", c.queue).Msg("Event consumer started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Event consumer stopping (context done)")
			return nil
		case <-c.stopCh:
			log.Info().Msg("Event consumer stopping")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				log.Error().Msg("Channel closed")
				return nil
			}

			if err := c.handleMessage(ctx, msg); err != nil {
				log.Error().Err(err).Str("message_id", msg.MessageId).Msg("Failed to handle message")
				// Negative acknowledge - requeue
				msg.Nack(false, true)
			} else {
				// Acknowledge
				msg.Ack(false)
			}
		}
	}
}

// Stop stops the event consumer
func (c *EventConsumer) Stop() error {
	close(c.stopCh)
	if c.amqpChan != nil {
		c.amqpChan.Close()
	}
	if c.amqpConn != nil {
		c.amqpConn.Close()
	}
	return nil
}

// handleMessage handles an incoming message
func (c *EventConsumer) handleMessage(ctx context.Context, msg amqp.Delivery) error {
	// Get event name from headers
	eventName, ok := msg.Headers["event_name"].(string)
	if !ok {
		return fmt.Errorf("missing event_name in headers")
	}

	eventKey, _ := msg.Headers["event_key"].(string)
	correlationID, _ := msg.Headers["correlation_id"].(string)

	log.Info().
		Str("event_name", eventName).
		Str("event_key", eventKey).
		Msg("Received event")

	// Record inbox event
	inbox := domain.InboxEvent{
		EventName:     eventName,
		EventKey:      eventKey,
		CorrelationID: correlationID,
		Payload:       string(msg.Body),
		Status:        "received",
		ReceivedAt:    time.Now(),
	}

	if err := c.db.Create(&inbox).Error; err != nil {
		log.Error().Err(err).Msg("Failed to record inbox event")
	}

	// Check for duplicate
	if eventKey != "" {
		var existing domain.InboxEvent
		err := c.db.First(&existing, "event_key = ? AND status = ?", eventKey, "processed").Error
		if err == nil {
			log.Info().Str("event_key", eventKey).Msg("Duplicate event - skipping")
			inbox.Status = "duplicate"
			c.db.Model(&inbox).Update("status", "duplicate")
			return nil
		}
	}

	// Process based on event name
	var processErr error
	switch {
	case strings.HasPrefix(eventName, "pmb.applicant_created"):
		processErr = c.handleApplicantCreated(ctx, msg.Body)
	case strings.HasPrefix(eventName, "academic.student_created"):
		processErr = c.handleStudentCreated(ctx, msg.Body)
	case strings.HasPrefix(eventName, "academic.student_updated"):
		processErr = c.handleStudentUpdated(ctx, msg.Body)
	default:
		log.Warn().Str("event_name", eventName).Msg("Unknown event type")
		processErr = nil
	}

	// Update inbox status
	if processErr != nil {
		inbox.Status = "failed"
		errMsg := processErr.Error()
		inbox.ErrorMessage = &errMsg
		log.Error().Err(processErr).Msg("Failed to process event")
	} else {
		inbox.Status = "processed"
		now := time.Now()
		inbox.ProcessedAt = &now
	}

	c.db.Model(&inbox).Updates(map[string]interface{}{
		"status":        inbox.Status,
		"error_message": inbox.ErrorMessage,
		"processed_at":  inbox.ProcessedAt,
	})

	return processErr
}

// handleApplicantCreated handles pmb.applicant_created event
func (c *EventConsumer) handleApplicantCreated(ctx context.Context, payload []byte) error {
	var event struct {
		ApplicantID       string `json:"applicant_id"`
		FullName          string `json:"full_name"`
		Email             string `json:"email"`
		PaymentComponents []struct {
			PaymentComponentID string  `json:"payment_component_id"`
			Amount             float64 `json:"amount"`
		} `json:"payment_components"`
		AutoCreateInvoice bool `json:"auto_create_invoice"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Only create invoice if auto_create_invoice is true
	if !event.AutoCreateInvoice {
		log.Info().Str("applicant_id", event.ApplicantID).Msg("Skipping invoice creation (auto_create_invoice=false)")
		return nil
	}

	if len(event.PaymentComponents) == 0 {
		log.Warn().Str("applicant_id", event.ApplicantID).Msg("No payment components - skipping invoice")
		return nil
	}

	// Create invoice
	invoice := domain.Invoice{
		InvoiceNumber: generateInvoiceNumber(),
		TargetType:    "applicant",
		ApplicantID:   &event.ApplicantID,
		Status:        "DRAFT",
	}

	var totalAmount float64
	for _, pc := range event.PaymentComponents {
		item := domain.InvoiceItem{
			PaymentComponentID: &pc.PaymentComponentID,
			Amount:             pc.Amount,
			FinalAmount:        pc.Amount,
		}
		invoice.Items = append(invoice.Items, item)
		totalAmount += pc.Amount
	}
	invoice.TotalAmount = totalAmount

	// Save invoice
	if err := c.db.Create(&invoice).Error; err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	log.Info().
		Str("applicant_id", event.ApplicantID).
		Str("invoice_id", invoice.ID).
		Msg("Invoice created from applicant event")

	return nil
}

// handleStudentCreated handles academic.student_created event
func (c *EventConsumer) handleStudentCreated(ctx context.Context, payload []byte) error {
	var event struct {
		StudentID   string `json:"student_id"`
		ApplicantID string `json:"applicant_id"`
		FullName    string `json:"full_name"`
		Email       string `json:"email"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if event.ApplicantID == "" {
		log.Info().Str("student_id", event.StudentID).Msg("No applicant_id - skipping")
		return nil
	}

	// Update invoice: link student_id to invoice with applicant_id
	result := c.db.Model(&domain.Invoice{}).
		Where("applicant_id = ? AND student_id IS NULL", event.ApplicantID).
		Update("student_id", event.StudentID)

	if result.Error != nil {
		return fmt.Errorf("failed to update invoice: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		log.Info().
			Str("student_id", event.StudentID).
			Str("applicant_id", event.ApplicantID).
			Int64("updated", result.RowsAffected).
			Msg("Updated invoices with student_id")
	}

	return nil
}

// handleStudentUpdated handles academic.student_updated event
func (c *EventConsumer) handleStudentUpdated(ctx context.Context, payload []byte) error {
	// Similar to student_created but for updates
	var event struct {
		StudentID string `json:"student_id"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Info().Str("student_id", event.StudentID).Msg("Student updated event received")
	return nil
}

func generateInvoiceNumber() string {
	now := time.Now().Format("20060102")
	return fmt.Sprintf("INV%s%04d", now, time.Now().UnixNano()%10000)
}

// InboxEvent represents inbox events table
type InboxEvent struct {
	ID            string     `gorm:"primaryKey;column:id"`
	EventName     string     `gorm:"column:event_name"`
	EventKey      string     `gorm:"column:event_key"`
	CorrelationID string     `gorm:"column:correlation_id"`
	Payload       string     `gorm:"type:jsonb;column:payload"`
	Status        string     `gorm:"column:status"` // received, processed, failed, duplicate
	ErrorMessage  *string    `gorm:"column:error_message"`
	ReceivedAt    time.Time  `gorm:"column:received_at"`
	ProcessedAt   *time.Time `gorm:"column:processed_at"`
}

func (InboxEvent) TableName() string {
	return "inbox_events"
}

// RunEventConsumer runs the event consumer
func RunEventConsumer(db *gorm.DB) error {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	consumer, err := NewEventConsumer(db, amqpURL)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Info().Msg("Received shutdown signal")
		cancel()
	}()

	err = consumer.Start(ctx)
	consumer.Stop()

	return err
}
