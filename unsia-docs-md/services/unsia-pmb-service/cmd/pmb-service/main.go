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
	"github.com/unsia-erp/unsia-pmb-service/internal/handler"
	"github.com/unsia-erp/unsia-pmb-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-pmb-service/internal/middleware"
	"gorm.io/gorm"
)

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
	_ = godotenv.Load()

	sharedobservability.InitLogger("pmb-service")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed (pmb_db): %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve sql.DB: %v", err)
	}

	sharedidempotency.RegisterStore(sharedidempotency.NewSQLStore(sqlDB, 30*time.Second))
	sharedaudit.RegisterWriter(&DbAuditWriter{db: db})

	jwksURL := os.Getenv("JWKS_URL")
	if jwksURL == "" {
		jwksURL = "http://localhost:8001/.well-known/jwks.json"
	}
	sharedauth.Configure(jwksURL, 5*time.Minute)

	if err := sharedauth.FetchJWKS(jwksURL); err != nil {
		sharedobservability.Logger.Warn().Err(err).Msg("Core Service JWKS not reachable on startup, entering degraded auth mode")
	}

	r := gin.New()

	r.Use(sharedobservability.CORSMiddleware())
	r.Use(sharedobservability.CorrelationIDMiddleware())
	r.Use(sharedobservability.RequestLoggerMiddleware())
	r.Use(sharedobservability.MetricsMiddleware())

	r.GET("/metrics", sharedobservability.MetricsHandler())

	pmbHandler := handler.NewPMBHandler(db)

// Protected routes
	protected := r.Group("/api", middleware.AuthRequired())
	{
		// Applicants
		protected.POST("/v1/pmb/applicants", pmbHandler.CreateApplicant)
		protected.GET("/v1/pmb/applicants", middleware.PermissionRequired("pmb.applicant.view"), pmbHandler.GetApplicants)
		protected.GET("/v1/pmb/applicants/:id", pmbHandler.GetApplicant)
		protected.POST("/v1/pmb/applicants/:id/submit", middleware.PermissionRequired("pmb.applicant.submit"), pmbHandler.SubmitApplicant)

		// Biodata
		protected.GET("/v1/pmb/applicants/:id/biodata", pmbHandler.GetBiodata)
		protected.PUT("/v1/pmb/applicants/:id/biodata", pmbHandler.UpdateBiodata)

		// Addresses
		protected.GET("/v1/pmb/applicants/:id/addresses", pmbHandler.GetAddresses)
		protected.PUT("/v1/pmb/applicants/:id/addresses", pmbHandler.UpdateAddresses)

		// Education Background
		protected.GET("/v1/pmb/applicants/:id/education", pmbHandler.GetEducation)
		protected.PUT("/v1/pmb/applicants/:id/education", pmbHandler.UpdateEducation)

		// Family Members
		protected.GET("/v1/pmb/applicants/:id/family", pmbHandler.GetFamily)
		protected.PUT("/v1/pmb/applicants/:id/family", pmbHandler.UpdateFamily)

		// Financial Profile
		protected.GET("/v1/pmb/applicants/:id/financial", pmbHandler.GetFinancial)
		protected.PUT("/v1/pmb/applicants/:id/financial", pmbHandler.UpdateFinancial)

		// Facility Profile
		protected.GET("/v1/pmb/applicants/:id/facility", pmbHandler.GetFacility)
		protected.PUT("/v1/pmb/applicants/:id/facility", pmbHandler.UpdateFacility)

		// Documents
		protected.POST("/v1/pmb/applicants/:id/documents", middleware.PermissionRequired("pmb.applicant.document"), pmbHandler.UploadDocument)
		protected.GET("/v1/pmb/applicants/:id/documents", pmbHandler.GetDocuments)
		protected.POST("/v1/pmb/applicants/:id/documents/:doc_id/verify", middleware.PermissionRequired("pmb.applicant.verify"), pmbHandler.VerifyDocument)

		// Payment & Invoice
		protected.POST("/v1/pmb/applicants/:id/request-invoice", middleware.PermissionRequired("pmb.applicant.manage"), pmbHandler.RequestInvoice)

		// Selection & Admission
		protected.POST("/v1/pmb/applicants/:id/selection-results", pmbHandler.ReceiveAssessmentSelectionResult)
		protected.POST("/v1/pmb/applicants/:id/issue-loa", middleware.PermissionRequired("pmb.applicant.manage"), pmbHandler.IssueLoa)

		// Academic Handover
		protected.POST("/v1/pmb/applicants/:id/handover-to-academic", middleware.PermissionRequired("pmb.applicant.manage"), pmbHandler.HandoverToAcademic)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	sharedobservability.Logger.Info().Msgf("PMB Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
