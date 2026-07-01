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
	"github.com/unsia-erp/unsia-reference-service/internal/handler"
	"github.com/unsia-erp/unsia-reference-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-reference-service/internal/middleware"
	"gorm.io/gorm"
)

// DbAuditWriter redirects global shared audit logs to the reference database audit_logs table
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
	sharedobservability.InitLogger("reference-service")

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
	// If Core is offline, this will print error but proceed (degraded mode)
	if err := sharedauth.FetchJWKS(jwksURL); err != nil {
		sharedobservability.Logger.Warn().Err(err).Msg("Core Service JWKS not reachable on startup, entering degraded auth mode")
	}

	// Setup Gin Engine
	r := gin.New()

	// Attach shared CORS and observability middlewares
	r.Use(sharedobservability.CORSMiddleware())
	r.Use(sharedobservability.CorrelationIDMiddleware())
	r.Use(sharedobservability.RequestLoggerMiddleware())
	r.Use(sharedobservability.MetricsMiddleware())

	// Register public route for Prometheus metrics
	r.GET("/metrics", sharedobservability.MetricsHandler())

	refHandler := handler.NewReferenceHandler(db)

	// Public routes
	r.GET("/api/v1/ref/status-codes", refHandler.ListStatusCodes)

// Protected routes (require valid login)
	protected := r.Group("/api", middleware.AuthRequired())
	{
		protected.GET("/v1/ref/study-programs", refHandler.ListStudyPrograms)
		protected.GET("/v1/ref/academic-years", refHandler.ListAcademicYears)
		protected.GET("/v1/ref/academic-periods", refHandler.ListAcademicPeriods)
		protected.GET("/v1/ref/payment-components", refHandler.ListPaymentComponents)
		protected.GET("/v1/ref/payment-methods", refHandler.ListPaymentMethods)
		protected.GET("/v1/ref/document-types", refHandler.ListDocumentTypes)
		protected.GET("/v1/ref/pmb-waves", refHandler.ListPmbWaves)
		protected.GET("/v1/ref/religions", refHandler.ListReligions)
		protected.GET("/v1/ref/admission-paths", refHandler.ListAdmissionPaths)
		protected.GET("/v1/ref/provinces", refHandler.ListProvinces)
		protected.GET("/v1/ref/cities", refHandler.ListCities)
		protected.GET("/v1/ref/districts", refHandler.ListDistricts)
		protected.GET("/v1/ref/villages", refHandler.ListVillages)

		// /v1/reference Alias mapping for portal compatibility
		refGroup := protected.Group("/v1/reference")
		{
			refGroup.GET("/study-programs", refHandler.ListStudyPrograms)
			refGroup.GET("/academic-years", refHandler.ListAcademicYears)
			refGroup.GET("/academic-periods", refHandler.ListAcademicPeriods)
			refGroup.GET("/payment-components", refHandler.ListPaymentComponents)
			refGroup.GET("/payment-methods", refHandler.ListPaymentMethods)
			refGroup.GET("/document-types", refHandler.ListDocumentTypes)
			refGroup.GET("/pmb-waves", refHandler.ListPmbWaves)
			refGroup.GET("/religions", refHandler.ListReligions)
			refGroup.GET("/admission-paths", refHandler.ListAdmissionPaths)
			refGroup.GET("/provinces", refHandler.ListProvinces)
			refGroup.GET("/cities", refHandler.ListCities)
			refGroup.GET("/districts", refHandler.ListDistricts)
			refGroup.GET("/villages", refHandler.ListVillages)
		}

		// Prefix-free Aliases
		protected.GET("/v1/provinces", refHandler.ListProvinces)
		protected.GET("/v1/cities", refHandler.ListCities)
		protected.GET("/v1/districts", refHandler.ListDistricts)
		protected.GET("/v1/villages", refHandler.ListVillages)

		// Administrative creation routes (restricted to administrative/biro roles)
		// Role check is performed at application/scope checking logic or via permissions
		academicAdmin := protected.Group("", middleware.PermissionRequired("reference.master_data.manage"))
		{
			academicAdmin.POST("/v1/ref/study-programs", refHandler.CreateStudyProgram)
			academicAdmin.POST("/v1/ref/academic-years", refHandler.CreateAcademicYear)
			academicAdmin.POST("/v1/ref/academic-periods", refHandler.CreateAcademicPeriod)
			academicAdmin.POST("/v1/ref/pmb-waves", refHandler.CreatePmbWave)
			academicAdmin.PUT("/v1/ref/pmb-waves/:id", refHandler.UpdatePmbWave)
			academicAdmin.DELETE("/v1/ref/pmb-waves/:id", refHandler.DeletePmbWave)
			academicAdmin.POST("/v1/ref/payment-components", refHandler.CreatePaymentComponent)
			academicAdmin.POST("/v1/ref/payment-methods", refHandler.CreatePaymentMethod)
			academicAdmin.POST("/v1/ref/document-types", refHandler.CreateDocumentType)
			academicAdmin.POST("/v1/ref/religions", refHandler.CreateReligion)
			academicAdmin.POST("/v1/ref/countries", refHandler.CreateCountry)

			// /v1/reference Administrative creation aliases
			academicAdmin.POST("/v1/reference/study-programs", refHandler.CreateStudyProgram)
			academicAdmin.POST("/v1/reference/academic-years", refHandler.CreateAcademicYear)
			academicAdmin.POST("/v1/reference/academic-periods", refHandler.CreateAcademicPeriod)
			academicAdmin.POST("/v1/reference/pmb-waves", refHandler.CreatePmbWave)
			academicAdmin.PUT("/v1/reference/pmb-waves/:id", refHandler.UpdatePmbWave)
			academicAdmin.DELETE("/v1/reference/pmb-waves/:id", refHandler.DeletePmbWave)
			academicAdmin.POST("/v1/reference/payment-components", refHandler.CreatePaymentComponent)
			academicAdmin.POST("/v1/reference/payment-methods", refHandler.CreatePaymentMethod)
			academicAdmin.POST("/v1/reference/document-types", refHandler.CreateDocumentType)
			academicAdmin.POST("/v1/reference/religions", refHandler.CreateReligion)
			academicAdmin.POST("/v1/reference/countries", refHandler.CreateCountry)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	sharedobservability.Logger.Info().Msgf("Reference Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
