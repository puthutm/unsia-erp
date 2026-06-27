package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/handler"
	"github.com/unsia-erp/unsia-core-service/internal/middleware"
	"github.com/unsia-erp/unsia-core-service/internal/service"
	"gorm.io/gorm"
)

// Config holds router configuration
type Config struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:         "8080",
		ReadTimeout:  15,
		WriteTimeout: 15,
		IdleTimeout:  60,
	}
}

// Setup sets up the router with all routes
func Setup(db *gorm.DB, cfg *Config) *gin.Engine {
	r := gin.New()
	serviceTokenService := service.NewServiceTokenService(db)
	
	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestIDMiddleware())

	// Health check (no auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, sharederr.Success(gin.H{
			"status": "healthy",
		}).WithContext(c))
	})
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, sharederr.Success(gin.H{
			"status": "ready",
		}).WithContext(c))
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Initialize handlers
		userHandler := handler.NewUserHandler(db)
		roleHandler := handler.NewRoleHandler(db)
		sessionHandler := handler.NewSessionHandler(db)
		applicationHandler := handler.NewApplicationHandler(db)
		tokenHandler := handler.NewServiceTokenHandler(db)
		auditHandler := handler.NewAuditHandler(db)
		webhookHandler := handler.NewWebhookHandler(db)
		externalAppHandler := handler.NewExternalAppHandler(db)
		authHandler := handler.NewAuthHandler(db)

		// Public routes (no auth required)
		public := v1.Group("")
		{
			public.POST("/auth/login", authHandler.Login)
			public.POST("/auth/refresh", sessionHandler.RefreshToken)
			
			// External app validation
			public.POST("/external-apps/validate", externalAppHandler.ValidateCredentials)
		}

		// Protected routes (JWT auth)
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// Auth
			protected.POST("/auth/logout", userHandler.Logout)
			protected.POST("/auth/change-password", userHandler.ChangePassword)
			protected.GET("/auth/me", authHandler.Me)

			// Users
			protected.GET("/users", userHandler.List)
			protected.POST("/users", userHandler.Create)
			protected.GET("/users/:id", userHandler.Get)
			protected.PUT("/users/:id", userHandler.Update)
			protected.DELETE("/users/:id", userHandler.Delete)
			protected.POST("/users/:id/activate", userHandler.ActivateUser)
			protected.POST("/users/:id/deactivate", userHandler.DeactivateUser)

			// Roles
			protected.GET("/roles", roleHandler.List)
			protected.POST("/roles", roleHandler.Create)
			protected.GET("/roles/:id", roleHandler.GetRole)
			protected.PUT("/roles/:id", roleHandler.UpdateRole)
			protected.DELETE("/roles/:id", roleHandler.DeleteRole)
			protected.POST("/roles/:id/assign", roleHandler.AssignRole)
			protected.POST("/roles/:id/revoke", roleHandler.RevokeRole)

			// Sessions
			protected.GET("/sessions", sessionHandler.ListSessions)
			protected.GET("/sessions/:id", sessionHandler.GetSession)
			protected.DELETE("/sessions/:id", sessionHandler.DeleteSession)
			protected.DELETE("/sessions", sessionHandler.RevokeAllSessions)

			// Applications
			protected.GET("/applications", applicationHandler.ListApplications)
			protected.POST("/applications", applicationHandler.CreateApplication)
			protected.GET("/applications/:id", applicationHandler.GetApplication)
			protected.PUT("/applications/:id", applicationHandler.UpdateApplication)
			protected.DELETE("/applications/:id", applicationHandler.DeleteApplication)

			// Service Tokens
			protected.GET("/service-tokens", tokenHandler.ListServiceTokens)
			protected.POST("/service-tokens", tokenHandler.CreateServiceToken)
			protected.GET("/service-tokens/:id", tokenHandler.GetServiceToken)
			protected.DELETE("/service-tokens/:id", tokenHandler.RevokeServiceToken)

			// Audit Logs
			protected.GET("/audit-logs", auditHandler.ListAuditLogs)
			protected.GET("/audit-logs/:id", auditHandler.GetAuditLog)

			// Webhooks
			protected.GET("/webhooks", webhookHandler.ListWebhooks)
			protected.POST("/webhooks", webhookHandler.CreateWebhook)
			protected.GET("/webhooks/:id", webhookHandler.GetWebhook)
			protected.PUT("/webhooks/:id", webhookHandler.UpdateWebhook)
			protected.DELETE("/webhooks/:id", webhookHandler.DeleteWebhook)
			protected.POST("/webhooks/:id/test", webhookHandler.TestWebhook)
			protected.POST("/webhooks/trigger", webhookHandler.TriggerEvent)

			// External Apps (admin)
			protected.GET("/external-apps", externalAppHandler.ListExternalApps)
			protected.POST("/external-apps", externalAppHandler.CreateExternalApp)
			protected.GET("/external-apps/:id", externalAppHandler.GetExternalApp)
			protected.PUT("/external-apps/:id", externalAppHandler.UpdateExternalApp)
			protected.DELETE("/external-apps/:id", externalAppHandler.DeleteExternalApp)
			protected.POST("/external-apps/:id/secret", externalAppHandler.RegenerateSecret)
		}
	}

	// Service-to-service routes (internal)
	internal := r.Group("/internal")
	internal.Use(middleware.ServiceTokenRequired(serviceTokenService))
	{
		internal.POST("/validate-token", func(c *gin.Context) {
			c.JSON(http.StatusOK, sharederr.Success(gin.H{
				"valid": true,
			}).WithContext(c))
		})
	}

	// No route matched
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Endpoint not found").WithContext(c))
	})

	return r
}

// Run runs the HTTP server
func Run(addr string, r *gin.Engine) error {
	return r.Run(addr)
}

// GetServer returns HTTP server with configuration
func GetServer(r *gin.Engine, cfg *Config) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}
}
