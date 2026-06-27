package sharedserviceclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BlackboxAI/unsia-docs-md/packages/shared-httpclient"
	"github.com/BlackboxAI/unsia-docs-md/packages/shared-errorenvelope"
)

// ServiceName represents available microservices
type ServiceName string

const (
	ServiceCore       ServiceName = "core"
	ServiceReference ServiceName = "reference"
	ServicePMB      ServiceName = "pmb"
	ServiceAcademic ServiceName = "academic"
	ServiceFinance  ServiceName = "finance"
	ServiceLMS     ServiceName = "lms"
	ServiceHRIS    ServiceName = "hris"
	ServiceAssessment ServiceName = "assessment"
	ServiceCRM     ServiceName = "crm"
	ServicePortal  ServiceName = "portal"
)

// Default service URLs from environment
var serviceURLs = map[ServiceName]string{
	ServiceCore:       "http://unsia-core-service:8001",
	ServiceReference: "http://unsia-reference-service:8002",
	ServicePMB:      "http://unsia-pmb-service:8003",
	ServiceAcademic: "http://unsia-academic-service:8004",
	ServiceFinance:  "http://unsia-finance-service:8005",
	ServiceLMS:     "http://unsia-lms-service:8006",
	ServiceHRIS:    "http://unsia-hris-service:8008",
	ServiceAssessment: "http://unsia-assessment-service:8007",
	ServiceCRM:     "http://unsia-crm-service:8009",
	ServicePortal:  "http://unsia-portal-service:8010",
}

// Config for service client
type Config struct {
	BaseURL          string        // Service base URL (overrides default from ServiceName)
	ServiceName     ServiceName   // Service identifier
	ServiceToken    string        // Service-to-service authentication token
	Timeout        time.Duration // Request timeout
	MaxRetries     int          // Max retry attempts
	FailureThreshold int         // Circuit breaker failure threshold
	Cooldown       time.Duration // Circuit breaker cooldown
}

// ServiceClient wraps HTTP client for specific service
type ServiceClient struct {
	name    ServiceName
	config Config
	client *sharedhttpclient.Client
	mu     sync.RWMutex
}

// New creates a new service client with configuration
func New(cfg Config) *ServiceClient {
	// Use default URL if not provided
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = serviceURLs[cfg.ServiceName]
	}

	// Override with environment variable if set
	envKey := fmt.Sprintf("%s_SERVICE_URL", strings.ToUpper(string(cfg.ServiceName)))
	if envURL := os.Getenv(envKey); envURL != "" {
		baseURL = envURL
	}

	serviceToken := cfg.ServiceToken
	if serviceToken == "" {
		serviceToken = os.Getenv("SERVICE_TOKEN")
	}

	httpCfg := sharedhttpclient.Config{
		BaseURL:          baseURL,
		ServiceToken:     serviceToken,
		SourceName:       string(cfg.ServiceName),
		Timeout:          cfg.Timeout,
		MaxRetries:       cfg.MaxRetries,
		FailureThreshold: cfg.FailureThreshold,
		Cooldown:         cfg.Cooldown,
	}

	sc := &ServiceClient{
		name:   cfg.ServiceName,
		config: cfg,
		client: sharedhttpclient.New(httpCfg),
	}

	return sc
}

// Factory provides service client creation with caching
type Factory struct {
	mu          sync.RWMutex
	clients    map[ServiceName]*ServiceClient
	defaultCfg Config
}

// NewFactory creates a new client factory
func NewFactory(defaultCfg Config) *Factory {
	return &Factory{
		clients: make(map[ServiceName]*ServiceClient),
		defaultCfg: defaultCfg,
	}
}

// GetClient returns or creates a service client
func (f *Factory) GetClient(name ServiceName, overrides ...Config) (*ServiceClient, error) {
	f.mu.RLock()
	if client, ok := f.clients[name]; ok {
		f.mu.RUnlock()
		return client, nil
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check after acquiring write lock
	if client, ok := f.clients[name]; ok {
		return client, nil
	}

	// Merge configs
	cfg := f.defaultCfg
	cfg.ServiceName = name

	for _, override := range overrides {
		if override.BaseURL != "" {
			cfg.BaseURL = override.BaseURL
		}
		if override.ServiceToken != "" {
			cfg.ServiceToken = override.ServiceToken
		}
		if override.Timeout > 0 {
			cfg.Timeout = override.Timeout
		}
	}

	client := New(cfg)
	f.clients[name] = client

	return client, nil
}

// Response types
type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *APIError      `json:"error,omitempty"`
	Meta   *ResponseMeta `json:"meta,omitempty"`
}

type APIError struct {
	Code    string      `json:"code"`
	Message string    `json:"message"`
	Details []string  `json:"details,omitempty"`
}

type ResponseMeta struct {
	RequestID  string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Page      int      `json:"page,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	Total     int64    `json:"total,omitempty"`
}

// Helper methods for ServiceClient

// Get performs GET request
func (sc *ServiceClient) Get(ctx context.Context, path string) (*APIResponse, error) {
	return sc.doRequest(ctx, http.MethodGet, path, nil)
}

// Post performs POST request
func (sc *ServiceClient) Post(ctx context.Context, path string, body interface{}) (*APIResponse, error) {
	return sc.doRequest(ctx, http.MethodPost, path, body)
}

// Put performs PUT request
func (sc *ServiceClient) Put(ctx context.Context, path string, body interface{}) (*APIResponse, error) {
	return sc.doRequest(ctx, http.MethodPut, path, body)
}

// Delete performs DELETE request
func (sc *ServiceClient) Delete(ctx context.Context, path string) (*APIResponse, error) {
	return sc.doRequest(ctx, http.MethodDelete, path, nil)
}

func (sc *ServiceClient) doRequest(ctx context.Context, method string, path string, body interface{}) (*APIResponse, error) {
	resp, err := sc.client.Request(ctx, method, path, body)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := &APIError{
			Code:    fmt.Sprintf("HTTP_%d", resp.StatusCode),
			Message: fmt.Sprintf("HTTP error: %d", resp.StatusCode),
		}
		return &APIResponse{
			Success: false,
			Error:  apiErr,
		}, nil
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResp, nil
}

// Generic request with typed response
func (sc *ServiceClient) GetTyped(ctx context.Context, path string, result interface{}) error {
	resp, err := sc.Get(ctx, path)
	if err != nil {
		return err
	}

	if !resp.Success {
		if resp.Error != nil {
			return fmt.Errorf("%s: %s", resp.Error.Code, resp.Error.Message)
		}
		return fmt.Errorf("unknown error")
	}

	if resp.Data != nil && result != nil {
		if err := json.Unmarshal(resp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}

func (sc *ServiceClient) PostTyped(ctx context.Context, path string, body, result interface{}) error {
	resp, err := sc.Post(ctx, path, body)
	if err != nil {
		return err
	}

	if !resp.Success {
		if resp.Error != nil {
			return fmt.Errorf("%s: %s", resp.Error.Code, resp.Error.Message)
		}
		return fmt.Errorf("unknown error")
	}

	if resp.Data != nil && result != nil {
		if err := json.Unmarshal(resp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}

// HealthCheck verifies service connectivity
func (sc *ServiceClient) HealthCheck(ctx context.Context) error {
	resp, err := sc.Get(ctx, "/api/v1/health")
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("health check failed")
	}
	return nil
}

// Use shared-errorenvelope for error responses
var _ = sharederrorenvelope.ErrorEnvelope{}
