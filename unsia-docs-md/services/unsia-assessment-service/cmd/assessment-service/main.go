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
	"github.com/unsia-erp/unsia-assessment-service/internal/handler"
	"github.com/unsia-erp/unsia-assessment-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-assessment-service/internal/middleware"
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

	sharedobservability.InitLogger("assessment-service")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed (assessment_db): %v", err)
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

	assessHandler := handler.NewAssessmentHandler(db)

	// Protected routes
	protected := r.Group("/api", middleware.AuthRequired())
	{
		// Sessions
		protected.POST("/v1/assessment/sessions", middleware.PermissionRequired("assessment.session.manage"), assessHandler.CreateSession)
		protected.GET("/v1/assessment/sessions", assessHandler.ListSessions)

		// Participants
		protected.POST("/v1/assessment/participants", assessHandler.RegisterParticipant)
		protected.GET("/v1/assessment/participants", assessHandler.ListParticipants)

		// Question Bank & Questions
		protected.POST("/v1/assessment/question-banks", middleware.PermissionRequired("assessment.question-bank.manage"), assessHandler.CreateQuestionBank)
		protected.POST("/v1/assessment/questions", middleware.PermissionRequired("assessment.question.manage"), assessHandler.CreateQuestion)
		protected.POST("/v1/assessment/questions/:id/versions", middleware.PermissionRequired("assessment.question.manage"), assessHandler.CreateQuestionVersion)

		// Attempts & Exam taking
		protected.POST("/v1/assessment/attempts", assessHandler.CreateAttempt)
		protected.GET("/v1/assessment/attempts", assessHandler.ListAttempts)
		protected.POST("/v1/assessment/attempts/:id/answers", assessHandler.SaveAnswer)
		protected.POST("/v1/assessment/attempts/:id/submit", assessHandler.SubmitAttempt)

		// Result Publication
		protected.POST("/v1/assessment/results/publish", assessHandler.PublishResult)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8009"
	}

	sharedobservability.Logger.Info().Msgf("Assessment Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
