package sharedhttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	mutex              sync.Mutex
	state              State
	failureCount       int
	consecutiveSuccess int
	failureThreshold   int
	cooldownDuration   time.Duration
	lastStateChange    time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: threshold,
		cooldownDuration: cooldown,
		lastStateChange:  time.Now(),
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	if cb.state == StateOpen {
		if now.Sub(cb.lastStateChange) > cb.cooldownDuration {
			cb.state = StateHalfOpen
			cb.lastStateChange = now
			return true
		}
		return false
	}
	return true
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateHalfOpen {
		cb.consecutiveSuccess++
		if cb.consecutiveSuccess >= 1 {
			cb.state = StateClosed
			cb.failureCount = 0
			cb.consecutiveSuccess = 0
			cb.lastStateChange = time.Now()
		}
	} else if cb.state == StateClosed {
		cb.failureCount = 0
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failureCount++
	if cb.state == StateHalfOpen || cb.failureCount >= cb.failureThreshold {
		cb.state = StateOpen
		cb.consecutiveSuccess = 0
		cb.lastStateChange = time.Now()
	}
}

type Config struct {
	BaseURL          string
	ServiceToken     string
	SourceName       string
	Timeout          time.Duration
	MaxRetries       int
	FailureThreshold int
	Cooldown         time.Duration
}

type Client struct {
	config Config
	client *http.Client
	cb     *CircuitBreaker
}

func New(cfg Config) *Client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.Cooldown <= 0 {
		cfg.Cooldown = 30 * time.Second
	}
	return &Client{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		cb: NewCircuitBreaker(cfg.FailureThreshold, cfg.Cooldown),
	}
}

// Request executes an HTTP request with timeout, retries, and circuit breaker protection
func (c *Client) Request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			bodyReader = v
		default:
			bytesData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(bytesData)
		}
	}

	url := c.config.BaseURL + path
	var lastErr error
	delay := 1 * time.Second

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
				delay *= 2
			}
		}

		if !c.cb.Allow() {
			return nil, errors.New("circuit breaker is open")
		}

		// Recreate body reader for subsequent retries if it's seekable
		if attempt > 0 && bodyReader != nil {
			if seeker, ok := bodyReader.(io.ReadSeeker); ok {
				_, _ = seeker.Seek(0, io.SeekStart)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Inject headers
		if c.config.ServiceToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.config.ServiceToken)
		}
		if c.config.SourceName != "" {
			req.Header.Set("X-Source-Service", c.config.SourceName)
		}
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		correlationID := extractCorrelationID(ctx)
		if correlationID != "" {
			req.Header.Set("X-Correlation-Id", correlationID)
		}

		resp, err := c.client.Do(req)
		if err == nil {
			// Successful execution (either success status or handled application errors)
			// Status >= 500 triggers circuit breaker and retry
			if resp.StatusCode < 500 {
				c.cb.RecordSuccess()
				return resp, nil
			}
			lastErr = fmt.Errorf("server error: status code %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		c.cb.RecordFailure()
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.MaxRetries+1, lastErr)
}

func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, path, body)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, path, body)
}

func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil)
}

func extractCorrelationID(ctx context.Context) string {
	if val, ok := ctx.Value("correlation_id").(string); ok {
		return val
	}
	if val, ok := ctx.Value("x-correlation-id").(string); ok {
		return val
	}
	// For gin.Context compatibility
	if gc, ok := ctx.(interface{ GetString(string) string }); ok {
		if cid := gc.GetString("x-correlation-id"); cid != "" {
			return cid
		}
		if tid := gc.GetString("trace_id"); tid != "" {
			return tid
		}
	}
	return ""
}
