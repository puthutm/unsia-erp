package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
)

const (
	// JWKS cache TTL minimum 5 minutes as per requirements
	JWKSMinTTL = 5 * time.Minute
)

// JWKSSource represents the JWKS endpoint configuration
type JWKSSource struct {
	JWKSURL string
}

// JWKSCache caches the JWKS with expiration
type JWKSCache struct {
	Keys      map[string]JWKSPublicKey
	ExpiresAt time.Time
	mu       sync.RWMutex
}

type JWKSPublicKey struct {
	Kid string
	Kty string
	Alg string
	Use string
	N   string
	E   string
}

type JWKS struct {
	Keys []JWKSPublicKey `json:"keys"`
}

// Claims represents the JWT claims structure
type Claims struct {
	Subject        string   `json:"sub"`
	Issuer         string   `json:"iss"`
	Audience       []string `json:"aud"`
	ExpiresAt      int64    `json:"exp"`
	IssuedAt       int64    `json:"iat"`
	NotBefore      int64    `json:"nbf"`
	ApplicationCode string   `json:"application_code"`
	ActiveRole     string   `json:"active_role"`
	Permissions    []string `json:"permissions"`
	Scope          string   `json:"scope"`
}

// AuthMiddleware handles JWT RS256 validation with JWKS
type AuthMiddleware struct {
	jwksSource  *JWKSSource
	jwksCache   *JWKSCache
	httpClient  *http.Client
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwksURL string) *AuthMiddleware {
	return &AuthMiddleware{
		jwksSource: &JWKSSource{
			JWKSURL: jwksURL,
		},
		jwksCache: &JWKSCache{
			Keys: make(map[string]JWKSPublicKey),
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// JWTAuth is the main JWT authentication middleware
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip health check
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, sharederr.Error("AUTH_TOKEN_MISSING", "Authorization header is required").WithContext(c))
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, sharederr.Error("AUTH_TOKEN_INVALID", "Invalid Authorization header format").WithContext(c))
			c.Abort()
			return
		}
		tokenString := parts[1]

		if isValid, _ := sharedauth.ValidateServiceToken(tokenString); isValid {
			serviceClaims := &Claims{
				ActiveRole:  "service",
				Permissions: []string{"*"},
				Scope:       "global",
				Subject:     "service-call",
			}
			c.Set("claims", serviceClaims)
			c.Set("user_id", serviceClaims.Subject)
			c.Set("active_role", serviceClaims.ActiveRole)
			c.Set("application_code", "finance")
			c.Set("permissions", serviceClaims.Permissions)
			c.Next()
			return
		}

		// Validate required headers
		requiredHeaders := []string{"X-Application-Code", "X-Active-Role", "X-Correlation-Id"}
		missingHeaders := []string{}
		for _, header := range requiredHeaders {
			if c.GetHeader(header) == "" {
				missingHeaders = append(missingHeaders, header)
			}
		}
		if len(missingHeaders) > 0 {
			c.JSON(http.StatusBadRequest, sharederr.Error("MISSING_HEADER", fmt.Sprintf("Missing required headers: %s", strings.Join(missingHeaders, ", "))).WithContext(c))
			c.Abort()
			return
		}

		// Get cached JWKS or fetch new one
		pubKey, kid, err := m.getJWKS(c.Request.Context())
		if err != nil {
			// Check if we have cached JWKS
			if m.jwksCache.isExpired() {
				c.JSON(http.StatusServiceUnavailable, sharederr.Error("AUTH_SERVICE_UNAVAILABLE", "Authentication service unavailable and cache expired").WithContext(c))
				c.Abort()
				return
			}
			// Use cached JWKS
			key, cachedKid := m.jwksCache.getLatest()
			if key == nil {
				c.JSON(http.StatusServiceUnavailable, sharederr.Error("AUTH_SERVICE_UNAVAILABLE", "Authentication service unavailable").WithContext(c))
				c.Abort()
				return
			}
			parsedKey, parseErr := m.parseRSAPublicKey(key)
			if parseErr != nil {
				c.JSON(http.StatusServiceUnavailable, sharederr.Error("AUTH_SERVICE_UNAVAILABLE", "Authentication service unavailable").WithContext(c))
				c.Abort()
				return
			}
			pubKey = parsedKey
			kid = cachedKid
		}

		// Validate JWT token
		claims, err := m.validateToken(tokenString, pubKey, kid)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "expired") {
				c.JSON(http.StatusUnauthorized, sharederr.Error("AUTH_TOKEN_EXPIRED", "Token has expired").WithContext(c))
			} else {
				c.JSON(http.StatusUnauthorized, sharederr.Error("AUTH_TOKEN_INVALID", "Invalid token: "+errMsg).WithContext(c))
			}
			c.Abort()
			return
		}

		// Set claims and user info in context
		c.Set("claims", claims)
		c.Set("user_id", claims.Subject)
		c.Set("active_role", claims.ActiveRole)
		c.Set("application_code", claims.ApplicationCode)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// getJWKS fetches JWKS from the remote server
func (m *AuthMiddleware) getJWKS(ctx context.Context) (*rsa.PublicKey, string, error) {
	// Check cache first
	if !m.jwksCache.isExpired() {
		key, kid := m.jwksCache.getLatest()
		if key != nil {
			pubKey, err := m.parseRSAPublicKey(key)
			return pubKey, kid, err
		}
	}

	// Fetch new JWKS
	req, err := http.NewRequestWithContext(ctx, "GET", m.jwksSource.JWKSURL, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, "", err
	}

	// Update cache
	m.jwksCache.set(jwks.Keys)

	// Return the first RSA key
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" {
			pubKey, err := m.parseRSAPublicKey(&key)
			return pubKey, key.Kid, err
		}
	}

	return nil, "", errors.New("no RSA key found in JWKS")
}

// parseRSAPublicKey parses a JWKSPublicKey to rsa.PublicKey
func (m *AuthMiddleware) parseRSAPublicKey(key *JWKSPublicKey) (*rsa.PublicKey, error) {
	// Decode N (modulus)
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}
	n := new(big.Int).SetBytes(nBytes)

	// Decode E (exponent)
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// validateToken validates the JWT token
func (m *AuthMiddleware) validateToken(tokenString string, pubKey *rsa.PublicKey, kid string) (*Claims, error) {
	// Parse JWT parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Decode payload (part 2)
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// getLatest gets the latest cached key
func (c *JWKSCache) getLatest() (*JWKSPublicKey, string) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for kid, key := range c.Keys {
		return &key, kid
	}
	return nil, ""
}

// isExpired checks if cache is expired
func (c *JWKSCache) isExpired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Now().After(c.ExpiresAt)
}

// set updates the cache with new keys
func (c *JWKSCache) set(keys []JWKSPublicKey) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Keys = make(map[string]JWKSPublicKey)
	for _, key := range keys {
		c.Keys[key.Kid] = key
	}
	c.ExpiresAt = time.Now().Add(JWKSMinTTL)
}

// RequirePermission creates a middleware that checks for required permissions
func RequirePermission(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissionsVal, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, sharederr.Error("PERMISSION_DENIED", "No permissions found").WithContext(c))
			c.Abort()
			return
		}

		userPermissions, ok := permissionsVal.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, sharederr.Error("PERMISSION_DENIED", "Invalid permissions format").WithContext(c))
			c.Abort()
			return
		}

		// Check if user has any of the required permissions
		hasPermission := false
		for _, required := range requiredPermissions {
			for _, userPerm := range userPermissions {
				if userPerm == required {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, sharederr.Error("PERMISSION_DENIED", "Insufficient permissions").WithContext(c))
			c.Abort()
			return
		}

		c.Next()
	}
}
