package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharedidempotency "github.com/unsia-erp/shared-idempotency"
	sharedobservability "github.com/unsia-erp/shared-observability"
	"github.com/unsia-erp/unsia-core-service/internal/handler"
	"github.com/unsia-erp/unsia-core-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-core-service/internal/infrastructure/keys"
	"github.com/unsia-erp/unsia-core-service/internal/middleware"
	"gorm.io/gorm"
)

// DbAuditWriter redirects global shared audit logs to the core database audit_logs table
type DbAuditWriter struct {
	db *gorm.DB
}

func (w *DbAuditWriter) Write(ctx context.Context, entry sharedaudit.AuditEntry) error {
	sqlDB, err := w.db.DB()
	if err != nil {
		return err
	}
	return sharedaudit.SaveToSQL(ctx, sqlDB, entry)
}

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize Logger
	sharedobservability.InitLogger("core-service")

	// Initialize Database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Retrieve underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve sql.DB: %v", err)
	}

	// Register shared idempotency store
	sharedidempotency.RegisterStore(sharedidempotency.NewSQLStore(sqlDB, 30*time.Second))

	// Register shared audit database writer
	sharedaudit.RegisterWriter(&DbAuditWriter{db: db})

	// Ensure RSA key pairs exist for RS256 token signing
	privKeyPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	pubKeyPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	privKey, err := keys.EnsureRSAKeys(privKeyPath, pubKeyPath)
	if err != nil {
		log.Fatalf("Failed to ensure RSA keypair: %v", err)
	}
	keys.SetSigningKey(privKey)

	// Configure local JWKS endpoint in shared auth cache
	jwksURL := os.Getenv("JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://localhost:8001/.well-known/jwks.json"
	}
	sharedauth.Configure(jwksURL, 5*time.Minute)

	// Configure JWKS Cache key initially
	_ = sharedauth.FetchJWKS(jwksURL)

	// Setup Gin Engine
	r := gin.New()

	// Attach shared CORS and observability middlewares
	r.Use(sharedobservability.CORSMiddleware())
	r.Use(sharedobservability.CorrelationIDMiddleware())
	r.Use(sharedobservability.RequestLoggerMiddleware())
	r.Use(sharedobservability.MetricsMiddleware())

	// Register public route for Prometheus metrics
	r.GET("/metrics", sharedobservability.MetricsHandler())

	authHandler := handler.NewAuthHandler(db)
	oauthHandler := handler.NewOAuthHandler(db)
	impersonationHandler := handler.NewImpersonationHandler(db)
	roleHandler := handler.NewRoleHandler(db)
	userHandler := handler.NewUserHandler(db)

	// Public OAuth / SSO endpoints
	r.POST("/api/v1/auth/login", authHandler.Login)
	r.POST("/api/v1/auth/refresh", authHandler.Refresh)
	r.GET("/.well-known/jwks.json", authHandler.JWKS)
	r.GET("/.well-known/openid-configuration", authHandler.OpenIDConfiguration)

	// Public OAuth 2.0 Server Endpoints
	r.GET("/api/v1/oauth/authorize", oauthHandler.Authorize)
	r.POST("/api/v1/oauth/token", oauthHandler.Token)
	r.POST("/api/v1/oauth/register", oauthHandler.Register)

	// Protected routes
	protected := r.Group("/api", middleware.AuthRequired())
	{
		protected.GET("/v1/auth/me", authHandler.Me)
		protected.POST("/v1/auth/switch-role", authHandler.SwitchRole)
		protected.GET("/v1/applications", authHandler.ListApplications)

		// Impersonation routes
		protected.POST("/v1/impersonations/start", impersonationHandler.Start)

		// Admin Role management routes
		protected.GET("/v1/admin/roles", roleHandler.List)
		protected.POST("/v1/admin/roles", roleHandler.Create)

		// Admin User management routes
		protected.POST("/v1/admin/users", userHandler.Create)
		protected.PUT("/v1/admin/users/:id", userHandler.Update)
		protected.PATCH("/v1/admin/users/:id/status", userHandler.UpdateStatus)
		protected.POST("/v1/admin/user-roles/:id/scopes", userHandler.AssignScope)

		// Admin OAuth Client management routes
		protected.GET("/v1/admin/oauth-clients", oauthHandler.AdminListRequests)
		protected.PATCH("/v1/admin/oauth-clients/:id/approve", oauthHandler.AdminApprove)
		protected.PATCH("/v1/admin/oauth-clients/:id/suspend", oauthHandler.AdminSuspend)
		protected.DELETE("/v1/admin/oauth-clients/:id/revoke", oauthHandler.AdminRevoke)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	sharedobservability.Logger.Info().Msgf("Core Auth Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
