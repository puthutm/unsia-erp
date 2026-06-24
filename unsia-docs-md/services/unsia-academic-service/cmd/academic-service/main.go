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
	"github.com/unsia-erp/unsia-academic-service/internal/handler"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-academic-service/internal/middleware"
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

	sharedobservability.InitLogger("academic-service")

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed (academic_db): %v", err)
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

academicHandler := handler.NewAcademicHandler(db)
	studentHandler := handler.NewStudentHandler(db)
	gradeHandler := handler.NewGradeHandler(db)
	graduationHandler := handler.NewGraduationHandler(db)

	// Protected routes
	protected := r.Group("/api", middleware.AuthRequired())
	{
		// Students (existing)
		protected.POST("/v1/academic/students/generate-from-applicant", academicHandler.GenerateStudentFromApplicant) // PMB calls this
		protected.GET("/v1/academic/students", academicHandler.ListStudents)
		protected.GET("/v1/academic/students/:id/transcripts", academicHandler.GetStudentTranscript)

		// Student Advisor (PA) management - NEW
		protected.POST("/v1/academic/students/:id/advisor", middleware.PermissionRequired("academic.student.advisor.assign"), studentHandler.AssignAdvisor)
		protected.DELETE("/v1/academic/students/:id/advisor", middleware.PermissionRequired("academic.student.advisor.assign"), studentHandler.RemoveAdvisor)
		protected.GET("/v1/academic/advisors/:lecturer_id/students", studentHandler.GetStudentsByAdvisor)

		// Curriculums & Courses
		protected.POST("/v1/academic/curriculums", middleware.PermissionRequired("academic.curriculum.manage"), academicHandler.CreateCurriculum)
		protected.POST("/v1/academic/courses", middleware.PermissionRequired("academic.course.manage"), academicHandler.CreateCourse)
		protected.POST("/v1/academic/curriculum-courses", middleware.PermissionRequired("academic.curriculum.manage"), academicHandler.CreateCurriculumCourse)

		// Classes & Offerings
		protected.POST("/v1/academic/classes", middleware.PermissionRequired("academic.class.manage"), academicHandler.CreateClass)
		protected.POST("/v1/academic/classes/:id/lecturers", middleware.PermissionRequired("academic.class.manage"), academicHandler.PlotClassLecturer)
		protected.POST("/v1/academic/course-offerings", middleware.PermissionRequired("academic.course-offering.manage"), academicHandler.CreateCourseOffering)

		// KRS
		protected.POST("/v1/academic/krs", middleware.PermissionRequired("academic.krs.create"), academicHandler.CreateKrsDraft)
		protected.POST("/v1/academic/krs/:krs_id/submit", middleware.PermissionRequired("academic.krs.create"), academicHandler.SubmitKrs)
		protected.POST("/v1/academic/krs/:krs_id/approve", middleware.PermissionRequired("academic.krs.approve"), academicHandler.ApproveKrs)

// Grades
		protected.POST("/v1/academic/grades/source-imports", academicHandler.ImportGradeSource) // LMS or Assessment calls this
		protected.POST("/v1/academic/grades/:grade_id/finalize", middleware.PermissionRequired("academic.grade.finalize"), academicHandler.FinalizeGrade)
		protected.POST("/v1/academic/grades/:id/corrections", middleware.PermissionRequired("academic.grade.correct"), academicHandler.CorrectGrade)

		// Grade Management (new handlers)
		protected.POST("/v1/academic/grades", middleware.PermissionRequired("academic.grade.manage"), gradeHandler.CreateGrade)
		protected.GET("/v1/academic/grades/:id", gradeHandler.GetGrade)
		protected.GET("/v1/academic/grades/student/:student_id", gradeHandler.GetStudentGrades)
		protected.POST("/v1/academic/grades/:id/submit", gradeHandler.SubmitGrade)
		protected.POST("/v1/academic/grades/:id/finalize", gradeHandler.FinalizeGrade)
		protected.POST("/v1/academic/grades/:id/entries", gradeHandler.EnterStudentGrade)
		protected.POST("/v1/academic/grades/:id/entries/bulk", gradeHandler.BulkEnterGrades)
		protected.POST("/v1/academic/grades/conversion", middleware.PermissionRequired("academic.grade.manage"), gradeHandler.UpdateGradeConversion)
		protected.GET("/v1/academic/grades/conversion", gradeHandler.GetGradeConversions)
		protected.GET("/v1/academic/grades/transcript/:student_id", gradeHandler.GetTranscript)
		protected.GET("/v1/academic/grades/ipk/:student_id", gradeHandler.GetIPK)
		protected.GET("/v1/academic/grades/ips/:student_id", gradeHandler.GetIPS)

		// Graduation
		protected.GET("/v1/academic/graduation/eligibility/:student_id", graduationHandler.CheckGraduationEligibility)
		protected.POST("/v1/academic/graduation/apply", graduationHandler.ApplyGraduation)
		protected.POST("/v1/academic/graduation/approve", middleware.PermissionRequired("academic.graduation.approve"), graduationHandler.ApproveGraduation)
		protected.POST("/v1/academic/graduation/certificate", graduationHandler.GenerateCertificate)
		protected.GET("/v1/academic/graduation/alumni", graduationHandler.GetAlumni)
		protected.GET("/v1/academic/graduation/alumni/:year", graduationHandler.GetAlumniByYear)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8006"
	}

	sharedobservability.Logger.Info().Msgf("Academic Service started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
