package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebhookService struct {
	db        *gorm.DB
	httpClient *http.Client
}

func NewWebhookService(db *gorm.DB) *WebhookService {
	return &WebhookService{
		db: db,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type Webhook struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	URL         string    `gorm:"column:url"`
	Secret      string    `gorm:"column:secret"`
	Event      string    `gorm:"column:event"`
	IsActive   bool      `gorm:"column:is_active"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Webhook) TableName() string {
	return "webhooks"
}

type CreateWebhookInput struct {
	URL       string   `json:"url" binding:"required,url"`
	Event     string   `json:"event" binding:"required"`
	Secret    string   `json:"secret"`
	IsActive  bool     `json:"is_active"`
}

type WebhookPayload struct {
	Event      string      `json:"event"`
	ID         string      `json:"id"`
	Timestamp  time.Time   `json:"timestamp"`
	Data       interface{} `json:"data"`
}

// CreateWebhook creates a new webhook
func (s *WebhookService) CreateWebhook(input CreateWebhookInput) (*Webhook, error) {
	webhook := Webhook{
		ID:         uuid.New().String(),
		URL:        input.URL,
		Event:      input.Event,
		Secret:     input.Secret,
		IsActive:   input.IsActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(&webhook).Error; err != nil {
		return nil, err
	}

	return &webhook, nil
}

// ListWebhooks lists webhooks by event type
func (s *WebhookService) ListWebhooks(event string) ([]Webhook, error) {
	var webhooks []Webhook
	query := s.db.Where("is_active = ?", true)

	if event != "" {
		query = query.Where("event = ?", event)
	}

	if err := query.Find(&webhooks).Error; err != nil {
		return nil, err
	}

	return webhooks, nil
}

// DeleteWebhook deletes a webhook
func (s *WebhookService) DeleteWebhook(id string) error {
	return s.db.Where("id = ?", id).Delete(&Webhook{}).Error
}

// TriggerWebhookTrigger sends event data to webhook URL
func (s *WebhookService) TriggerWebhook(webhook Webhook, payload WebhookPayload) error {
	// Create signature
	var signature string
	if webhook.Secret != "" {
		signature = generateHMAC(payload.Data, webhook.Secret)
	}

	// Marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", payload.Event)
	req.Header.Set("X-Webhook-ID", payload.ID)
	if signature != "" {
		req.Header.Set("X-Webhook-Signature", signature)
	}

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// TriggerEvent triggers all webhooks for an event
func (s *WebhookService) TriggerEvent(event string, data interface{}) error {
	webhooks, err := s.ListWebhooks(event)
	if err != nil {
		return err
	}

	payload := WebhookPayload{
		Event:     event,
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Data:      data,
	}

	for _, webhook := range webhooks {
		if err := s.TriggerWebhook(webhook, payload); err != nil {
			// Log error but continue with other webhooks
			fmt.Printf("Webhook error for %s: %v\n", webhook.URL, err)
		}
	}

	return nil
}

// ProcessWebhookRetry retries failed webhook deliveries
func (s *WebhookService) ProcessWebhookRetry(ctx context.Context) error {
	// This would typically query a delivery_log table for failed webhooks
	// and retry them with exponential backoff
	return nil
}

func generateHMAC(data interface{}, secret string) string {
	// Simple HMAC generation - in production use proper HMAC
	jsonData, _ := json.Marshal(data)
	hash := fmt.Sprintf("%x", simpleHash(string(jsonData)+secret))
	return hash
}

func simpleHash(s string) int {
	h := 0
	for _, c := range s {
		h = h*31 + int(c)
	}
	return h
}
