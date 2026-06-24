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
	"github.com/unsia-erp/unsia-hris-service/internal/handler"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-hris-service/internal/middleware"
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

	sharedobservability.InitLogger("hris-service")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed (hris_db): %v", err)
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

	hrisHandler := handler.NewHRISHandler(db)

	// Protected routes
	protected := r.Group("/api", middleware.AuthRequired())
		protected.GET("/v1/hris/lecturers", hrisHandler.ListActiveLecturers)
		protected.GET("/v1/hris/lecturers/:id", hrisHandler.GetLecturer)
		protected.POST("/v1/hris/employees", hrisHandler.CreateEmployee)
		protected.POST("/v1/hris/lecturers", hrisHandler.CreateLecturer)
		protected.PUT("/v1/hris/employees/:id/positions", hrisHandler.PlotPosition)
		protected.POST("/v1/hris/bkd-records", hrisHandler.CreateBkdRecord)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8007"
	}

	sharedobservability.Logger.Info().Msgf("HRIS Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
