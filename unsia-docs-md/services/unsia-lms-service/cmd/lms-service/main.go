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
	"github.com/unsia-erp/unsia-lms-service/internal/handler"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-lms-service/internal/middleware"
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

	sharedobservability.InitLogger("lms-service")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed (lms_db): %v", err)
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

	r.Use(sharedobservability.CorrelationIDMiddleware())
	r.Use(sharedobservability.RequestLoggerMiddleware())
	r.Use(sharedobservability.MetricsMiddleware())

	r.GET("/metrics", sharedobservability.MetricsHandler())

	lmsHandler := handler.NewLMSHandler(db)

	// Protected routes
	protected := r.Group("/api", middleware.AuthRequired())
	{
		// Synced triggers (called by Academic Service or events)
		protected.POST("/v1/lms/classes/sync-from-academic", lmsHandler.SyncClassFromAcademic)
		protected.POST("/v1/lms/enrollments/sync-from-krs", lmsHandler.SyncEnrollmentFromKrs)

		// Grade Sync trigger
		protected.POST("/v1/lms/grade-syncs", lmsHandler.SyncLmsGrade)

		// Sessions & Materials
		protected.POST("/v1/lms/classes/:id/sessions", lmsHandler.CreateSession)
		protected.POST("/v1/lms/sessions/:id/materials", lmsHandler.CreateMaterial)
		protected.POST("/v1/lms/sessions/:id/attendance", lmsHandler.CreateAttendance)

		// Assignments & Submissions
		protected.POST("/v1/lms/assignments", lmsHandler.CreateAssignment)
		protected.POST("/v1/lms/assignments/:id/submissions", lmsHandler.CreateSubmission)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}

	sharedobservability.Logger.Info().Msgf("LMS Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
