package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu           sync.Mutex
	tokens       float64
	maxTokens    float64
	refillRate   float64 // tokens per second
	lastRefill   time.Time
}

func NewRateLimiter(tokensPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		tokens:       float64(burst),
		maxTokens:    float64(burst),
		refillRate:   tokensPerSecond,
		lastRefill:   time.Now(),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill).Seconds()
	r.tokens += elapsed * r.refillRate
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}
	r.lastRefill = now

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}

// RateLimitMiddleware creates rate limiting middleware
func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerSecond, burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, sharederr.Error("RATE_LIMITED", "Too many requests. Please try again later.").WithContext(c))
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimiter maps IP addresses to their own rate limiters
type IPRateLimiter struct {
	mu          sync.Mutex
	limiters   map[string]*RateLimiter
	window     time.Duration
	maxHits   int
}

func NewIPRateLimiter(requestsPerSecond float64, burst int, window time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
		window:   window,
		maxHits:  burst,
	}
}

func (r *IPRateLimiter) getLimiter(ip string) *RateLimiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limiter, exists := r.limiters[ip]; exists {
		return limiter
	}

	limiter := NewRateLimiter(10, 50) // default
	r.limiters[ip] = limiter
	return limiter
}

// IPRateLimitMiddleware rate limits by client IP
func IPRateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(requestsPerSecond, burst, time.Minute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, sharederr.Error("RATE_LIMITED", "Too many requests from your IP. Please try again later.").WithContext(c))
			c.Abort()
			return
		}
		c.Next()
	}
}
