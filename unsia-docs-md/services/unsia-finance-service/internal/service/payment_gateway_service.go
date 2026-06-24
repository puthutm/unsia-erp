package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Provider represents a payment gateway provider
type Provider string

const (
	ProviderMidtrans Provider = "midtrans"
	ProviderXendit   Provider = "xendit"
)

// PaymentGatewayService handles payment gateway integration
type PaymentGatewayService struct {
	db *gorm.DB
	// Provider configurations
	providerSecrets map[Provider]string
	providerUrls     map[Provider]string
}

// NewPaymentGatewayService creates a new payment gateway service
func NewPaymentGatewayService(db *gorm.DB) *PaymentGatewayService {
	return &PaymentGatewayService{
		db: db,
		providerSecrets: map[Provider]string{
			ProviderMidtrans: os.Getenv("MIDTRANS_SERVER_KEY"),
			ProviderXendit:  os.Getenv("XENDIT_API_KEY"),
		},
		providerUrls: map[Provider]string{
			ProviderMidtrans: os.Getenv("MIDTRANS_CALLBACK_URL"),
			ProviderXendit:  os.Getenv("XENDIT_CALLBACK_URL"),
		},
	}
}

// CallbackPayload represents a callback payload from payment gateway
type CallbackPayload struct {
	ProviderEventID   string `json:"provider_event_id"`
	OrderID          string `json:"order_id"`
	Amount           int64  `json:"amount"`
	Status           string `json:"status"`
	ExternalReference string `json:"external_reference,omitempty"`
	Signature        string `json:"signature,omitempty"`
}

// ValidateSignature validates the callback signature from the provider
func (s *PaymentGatewayService) ValidateSignature(provider Provider, payload []byte, signature string) (bool, error) {
	secret, exists := s.providerSecrets[provider]
	if !exists || secret == "" {
		return false, fmt.Errorf("provider %s not configured", provider)
	}

	switch provider {
	case ProviderMidtrans:
		return s.validateMidtransSignature(payload, signature, secret)
	case ProviderXendit:
		return s.validateXenditSignature(payload, signature, secret)
	default:
		return false, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// validateMidtransSignature validates Midtrans signature
// Midtrans uses: SHA512(order_id + status_code + gross_amount + server_key)
func (s *PaymentGatewayService) validateMidtransSignature(payload []byte, signature, serverKey string) (bool, error) {
	var callback CallbackPayload
	if err := json.Unmarshal(payload, &callback); err != nil {
		return false, fmt.Errorf("invalid payload: %w", err)
	}

	// Construct the string to hash: order_id + status_code + gross_amount + server_key
	dataToHash := fmt.Sprintf("%s%s%d%s", callback.OrderID, callback.Status, callback.Amount, serverKey)
	
	// Calculate HMAC-SHA512
	h := hmac.New(sha256.New, []byte(serverKey))
	h.Write([]byte(dataToHash))
	expectedSignature := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

// validateXenditSignature validates Xendit signature
// Xendit uses: HMAC-SHA256 with callback token
func (s *PaymentGatewayService) validateXenditSignature(payload []byte, signature, apiKey string) (bool, error) {
	// Xendit uses the raw request body to validate
	h := hmac.New(sha256.New, []byte(apiKey))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

// ProcessCallback processes a payment gateway callback
func (s *PaymentGatewayService) ProcessCallback(provider Provider, payload CallbackPayload) (*CallbackResult, error) {
	// Validate required fields
	if payload.ProviderEventID == "" {
		return nil, fmt.Errorf("provider_event_id is required")
	}
	if payload.OrderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}
	if payload.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if payload.Status == "" {
		return nil, fmt.Errorf("status is required")
	}

	result := &CallbackResult{
		ProviderEventID:  payload.ProviderEventID,
		OrderID:      payload.OrderID,
		Amount:      payload.Amount,
		Status:      payload.Status,
		ProcessedAt: time.Now(),
	}

	// Map provider status to internal status
	switch strings.ToLower(payload.Status) {
	case "capture", "settlement", "success":
		result.InternalStatus = "success"
	case "expire", "expired":
		result.InternalStatus = "expired"
	case "cancel", "cancelled":
		result.InternalStatus = "cancelled"
	case "deny", "denied":
		result.InternalStatus = "failed"
	case "pending":
		result.InternalStatus = "pending"
	default:
		result.InternalStatus = "unknown"
	}

	return result, nil
}

// CallbackResult represents the result of processing a callback
type CallbackResult struct {
	ProviderEventID   string
	OrderID        string
	Amount        int64
	Status        string
	InternalStatus string
	ProcessedAt   time.Time
	ID            string
	PaymentID     string
}

// ProviderNotConfiguredError returns true if provider is not configured
func (s *PaymentGatewayService) ProviderNotConfigured(provider Provider) bool {
	secret, exists := s.providerSecrets[provider]
	return !exists || secret == ""
}

// GetProvider returns the provider from string
func GetProvider(provider string) (Provider, error) {
	switch strings.ToLower(provider) {
	case "midtrans":
		return ProviderMidtrans, nil
	case "xendit":
		return ProviderXendit, nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

// IsDuplicateCallback checks if callback with same provider_event_id already exists
func (s *PaymentGatewayService) IsDuplicateCallback(provider Provider, providerEventID string) (bool, error) {
	type Callback struct {
		ID              string `gorm:"primaryKey;column:id"`
		Provider        string `gorm:"column:provider"`
		ProviderEventID string `gorm:"column:provider_event_id"`
		CallbackStatus string `gorm:"column:callback_status"`
	}

	var existing Callback
	err := s.db.Table("payment_gateway_callbacks").
		Where("provider = ? AND provider_event_id = ?", provider, providerEventID).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
