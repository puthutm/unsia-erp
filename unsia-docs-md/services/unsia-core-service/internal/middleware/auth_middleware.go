package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
)

// AuthRequired verifies JWT token on incoming requests and populates context with claims
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing Authorization header").WithContext(c))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Authorization header format must be 'Bearer <token>'").WithContext(c))
			c.Abort()
			return
		}

		claims, err := sharedauth.ValidateJWTOrServiceToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Invalid or expired access token").WithContext(c))
			c.Abort()
			return
		}

		// Save claims to context for downstream handlers
		c.Set("claims", claims)
		c.Set("user_claims", claims)
		c.Set("user_id", claims.Subject)
		c.Set("userID", claims.Subject)

		c.Next()
	}
}
