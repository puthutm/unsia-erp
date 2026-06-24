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
	"github.com/unsia-erp/unsia-crm-service/internal/handler"
	"github.com/unsia-erp/unsia-crm-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-crm-service/internal/middleware"
	"gorm.io/gorm"
)

// DbAuditWriter redirects global shared audit logs to the crm database audit_logs table
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
	sharedobservability.InitLogger("crm-service")

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

	// Configure local JWKS client endpoint to validate token signatures
	jwksURL := os.Getenv("JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://localhost:8001/.well-known/jwks.json"
	}
	sharedauth.Configure(jwksURL, 5*time.Minute)

	// Fetch JWKS keys from Core Service on startup
	if err := sharedauth.FetchJWKS(jwksURL); err != nil {
		sharedobservability.Logger.Warn().Err(err).Msg("Core Service JWKS not reachable on startup, entering degraded auth mode")
	}

	// Setup Gin Engine
	r := gin.New()

	// Attach shared observability middlewares
	r.Use(sharedobservability.CorrelationIDMiddleware())
	r.Use(sharedobservability.RequestLoggerMiddleware())
	r.Use(sharedobservability.MetricsMiddleware())

	// Register public route for Prometheus metrics
	r.GET("/metrics", sharedobservability.MetricsHandler())

	crmHandler := handler.NewCRMHandler(db)

	// Protected routes (require valid login)
	protected := r.Group("/api", middleware.AuthRequired())
	{
		// Campaigns
		protected.GET("/v1/crm/campaigns", middleware.PermissionRequired("crm.campaign.view"), crmHandler.ListCampaigns)
		protected.POST("/v1/crm/campaigns", middleware.PermissionRequired("crm.campaign.manage"), crmHandler.CreateCampaign)

		// Agents
		protected.GET("/v1/crm/agents", middleware.PermissionRequired("crm.agent.view"), crmHandler.ListAgents)
		protected.POST("/v1/crm/agents", crmHandler.CreateAgent) // Allow any authenticated user to register
		protected.POST("/v1/crm/agents/:id/approve", middleware.PermissionRequired("crm.agent.approve"), crmHandler.ApproveAgent)

		// Referrals
		protected.GET("/v1/crm/referrals", middleware.PermissionRequired("crm.referral.view"), crmHandler.ListReferrals)
		protected.POST("/v1/crm/referrals", middleware.PermissionRequired("crm.referral.manage"), crmHandler.CreateReferral)

		// Leads
		protected.GET("/v1/crm/leads", middleware.PermissionRequired("crm.lead.view"), crmHandler.ListLeads)
		protected.POST("/v1/crm/leads", middleware.PermissionRequired("crm.lead.create"), crmHandler.CreateLead)
		protected.POST("/v1/crm/leads/:id/activities", middleware.PermissionRequired("crm.lead.activity"), crmHandler.CreateLeadActivity)
		protected.POST("/v1/crm/leads/:id/convert-to-applicant", middleware.PermissionRequired("crm.lead.convert"), crmHandler.ConvertLeadToApplicant)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	sharedobservability.Logger.Info().Msgf("CRM Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
