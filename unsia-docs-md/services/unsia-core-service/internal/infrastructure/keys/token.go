package keys

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sharedauth "github.com/unsia-erp/shared-auth"
)

var (
	SigningKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	KeyID      = "core-key-1"
)

// SetSigningKey configures the RSA private key for signing JWTs
func SetSigningKey(key *rsa.PrivateKey) {
	SigningKey = key
	PublicKey = &key.PublicKey
}

// GenerateAccessToken signs a new access token using RS256 algorithm
func GenerateAccessToken(userID string, activeRole string, scope string, permissions []string, expireHours int) (string, error) {
	if expireHours <= 0 {
		expireHours = 24 // default 24 hours
	}

	claims := sharedauth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "unsia-core-service",
		},
		ActiveRole:  activeRole,
		Permissions: permissions,
		Scope:       scope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = KeyID

	return token.SignedString(SigningKey)
}
