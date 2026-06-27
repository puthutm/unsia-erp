package sharedauth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	ActiveRole  string   `json:"active_role"`
	Permissions []string `json:"permissions"`
	Scope       string   `json:"scope"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

var (
	keysCache   = make(map[string]*rsa.PublicKey)
	cacheMutex  sync.RWMutex
	lastFetch   time.Time
	jwksURL     string
	jwksTTL     = 5 * time.Minute
	httpClient  = &http.Client{Timeout: 10 * time.Second}
)

// Configure sets the JWKS URL and TTL
func Configure(url string, ttl time.Duration) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	jwksURL = url
	if ttl > 0 {
		jwksTTL = ttl
	}
}

// FetchJWKS fetches the JWKS from the given URL and populates the cache
func FetchJWKS(url string) error {
	cacheMutex.Lock()
	if url != "" {
		jwksURL = url
	}
	targetURL := jwksURL
	cacheMutex.Unlock()

	if targetURL == "" {
		return fmt.Errorf("JWKS URL is not configured")
	}

	resp, err := httpClient.Get(targetURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS, status code: %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	newCache := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" && (key.Use == "sig" || key.Use == "") {
			pubKey, err := parseRSAPublicKey(key.N, key.E)
			if err != nil {
				continue // Skip malformed keys
			}
			newCache[key.Kid] = pubKey
		}
	}

	cacheMutex.Lock()
	keysCache = newCache
	lastFetch = time.Now()
	cacheMutex.Unlock()

	return nil
}

// parseRSAPublicKey reconstructs an RSA public key from JWK modulus (n) and exponent (e)
func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	decN, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}
	decE, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	var eVal int
	for _, b := range decE {
		eVal = (eVal << 8) | int(b)
	}

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(decN),
		E: eVal,
	}
	return pubKey, nil
}

// RefreshJWKSIfStale refreshes JWKS if TTL has expired
func RefreshJWKSIfStale() error {
	cacheMutex.RLock()
	stale := time.Since(lastFetch) > jwksTTL
	url := jwksURL
	cacheMutex.RUnlock()

	if stale && url != "" {
		return FetchJWKS(url)
	}
	return nil
}

// ValidateJWT validates a token string against the cached JWKS
func ValidateJWT(tokenStr string) (*Claims, error) {
	// Try refreshing if cache is stale
	_ = RefreshJWKSIfStale()

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in token header")
		}

		cacheMutex.RLock()
		pubKey, exists := keysCache[kid]
		cacheMutex.RUnlock()

		if !exists {
			// Try fetching JWKS again to check if it's a new key
			if err := FetchJWKS(""); err == nil {
				cacheMutex.RLock()
				pubKey, exists = keysCache[kid]
				cacheMutex.RUnlock()
			}
		}

		if !exists {
			return nil, fmt.Errorf("unknown key ID: %s", kid)
		}

		return pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse/validate token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ExtractClaims extracts claims from a token without validating the signature (useful for logging/tracing)
func ExtractClaims(tokenStr string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse unverified token: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// ValidateServiceToken validates a machine-to-machine service token
func ValidateServiceToken(token string) (bool, error) {
	// First check local environment variables for service token
	expectedToken := os.Getenv("SERVICE_TOKEN")
	if expectedToken != "" && token == expectedToken {
		return true, nil
	}

	// Also check individual service tokens
	for _, envKey := range []string{"CORE_SERVICE_TOKEN", "REFERENCE_SERVICE_TOKEN", "CRM_SERVICE_TOKEN", "PMB_SERVICE_TOKEN", "FINANCE_SERVICE_TOKEN", "ACADEMIC_SERVICE_TOKEN", "HRIS_SERVICE_TOKEN", "LMS_SERVICE_TOKEN", "ASSESSMENT_SERVICE_TOKEN", "PORTAL_SERVICE_TOKEN", "INTEGRATION_SERVICE_TOKEN"} {
		expected := os.Getenv(envKey)
		if expected != "" && token == expected {
			return true, nil
		}
	}

	// Alternatively, try validating as JWT and check for system active_role
	claims, err := ValidateJWT(token)
	if err == nil {
		if claims.ActiveRole == "service" || claims.ActiveRole == "system" {
			return true, nil
		}
	}

	return false, fmt.Errorf("invalid service token")
}

// ValidateJWTOrServiceToken checks if a token is a valid service token first. If it is, it returns a service-level claims. Otherwise, it falls back to standard ValidateJWT.
func ValidateJWTOrServiceToken(tokenStr string) (*Claims, error) {
	if isValid, _ := ValidateServiceToken(tokenStr); isValid {
		return &Claims{
			ActiveRole: "service",
			Permissions: []string{"*"},
			Scope: "global",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "service-call",
			},
		}, nil
	}
	return ValidateJWT(tokenStr)
}

