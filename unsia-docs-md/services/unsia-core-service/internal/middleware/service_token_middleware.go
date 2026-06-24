package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/service"
)

// ServiceTokenRequired verifies ServiceToken for service-to-service communication
func ServiceTokenRequired(serviceTokenService *service.ServiceTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("X-Service-Token")
		if authHeader == "" {
			// Also accept from Authorization header with "ServiceToken" scheme
			authHeader = c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(401, sharederr.Error("UNAUTHORIZED", "Missing X-Service-Token header").WithContext(c))
				c.Abort()
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "servicetoken" {
				c.JSON(401, sharederr.Error("UNAUTHORIZED", "Authorization header format must be 'ServiceToken <token>'").WithContext(c))
				c.Abort()
				return
			}
			authHeader = parts[1]
		}

		// Validate service token
		st, err := serviceTokenService.ValidateServiceToken(authHeader)
		if err != nil {
			c.JSON(401, sharederr.Error("UNAUTHORIZED", "Invalid or expired service token").WithContext(c))
			c.Abort()
			return
		}

		// Save service token info to context for downstream handlers
		c.Set("service_token", st)
		c.Set("application_id", st.ApplicationID)

		c.Next()
	}
}

// OptionalServiceToken validates ServiceToken if provided, but does not require it
func OptionalServiceToken(serviceTokenService *service.ServiceTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("X-Service-Token")
		if authHeader == "" {
			authHeader = c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "servicetoken" {
					authHeader = parts[1]
				} else {
					authHeader = ""
				}
			}
		}

		if authHeader != "" {
			st, err := serviceTokenService.ValidateServiceToken(authHeader)
			if err == nil {
				c.Set("service_token", st)
				c.Set("application_id", st.ApplicationID)
			}
		}

		c.Next()
	}
}
