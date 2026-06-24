package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharedidempotency "github.com/unsia-erp/shared-idempotency"
	sharedobservability "github.com/unsia-erp/shared-observability"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/handler"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/database"
	"github.com/unsia-erp/unsia-finance-service/internal/middleware"
	"github.com/unsia-erp/unsia-finance-service/internal/service"
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

	sharedobservability.InitLogger("finance-service")

	// Database initialization
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Database connection failed (finance_db): %v", err)
	}

	// Auto migrate new domain models
	if err := db.AutoMigrate(
		&domain.InboxEvent{},
		&domain.OutboxEvent{},
		&domain.IdempotencyKey{},
	); err != nil {
		log.Printf("Warning: AutoMigrate failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve sql.DB: %v", err)
	}

	sharedidempotency.RegisterStore(sharedidempotency.NewSQLStore(sqlDB, 30*time.Second))
	sharedaudit.RegisterWriter(&DbAuditWriter{db: db})

	// JWKS Configuration
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

// Health check endpoint (public - no auth required)
	healthHandler := handler.NewHealthHandler(db)
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/health/live", healthHandler.LivenessCheck)
	r.GET("/health/ready", healthHandler.ReadinessCheck)

	// Initialize services
	invoiceService := service.NewInvoiceService(db)
	paymentGatewayService := service.NewPaymentGatewayService(db)
	financeHandler := handler.NewFinanceHandler(db)
	financeHandler.InvoiceService = invoiceService

	// Auth middleware - use NewAuthMiddleware (new implementation)
	jwtMiddleware := middleware.NewAuthMiddleware(jwksURL)

	// Protected routes
	protected := r.Group("/api", jwtMiddleware.JWTAuth())
	{
		// Invoices
		protected.POST("/v1/finance/invoices", financeHandler.CreateInvoice)
		protected.GET("/v1/finance/invoices", financeHandler.GetInvoices)
		protected.GET("/v1/finance/invoices/:id", financeHandler.GetInvoice)

		// Payments
		protected.GET("/v1/finance/payments", financeHandler.GetPayments)

		// Webhook callback (no permission required - uses provider signature)
		protected.POST("/v1/finance/payment-callbacks/:provider", financeHandler.ReceivePaymentCallback)

		// Verification
		protected.POST("/v1/finance/payment-verifications", middleware.RequirePermission("finance.payment.verify"), financeHandler.VerifyManualPayment)

		// Clearances/List
		protected.GET("/v1/finance/clearances", financeHandler.GetClearances)
		protected.GET("/v1/finance/clearances/check", financeHandler.CheckClearance)
		protected.POST("/v1/finance/clearance-policies", middleware.RequirePermission("finance.clearance.manage"), financeHandler.CreateClearancePolicy)
		protected.PUT("/v1/finance/clearance-policies/:id", middleware.RequirePermission("finance.clearance-policy.manage"), financeHandler.UpdateClearancePolicy)

// Installment/dispensation requests
		protected.POST("/v1/finance/installment-requests", middleware.RequirePermission("finance.installment.request"), financeHandler.CreateInstallmentRequest)
		protected.PATCH("/v1/finance/installment-requests/:id/approve", middleware.RequirePermission("finance.installment.approve"), financeHandler.ApproveInstallmentRequest)

		// Scholarships
		protected.GET("/v1/finance/scholarships", financeHandler.GetScholarships)
		protected.POST("/v1/finance/scholarships", middleware.RequirePermission("finance.scholarship.manage"), financeHandler.CreateScholarship)
		protected.PUT("/v1/finance/scholarships/:id", middleware.RequirePermission("finance.scholarship.manage"), financeHandler.UpdateScholarship)
		protected.DELETE("/v1/finance/scholarships/:id", middleware.RequirePermission("finance.scholarship.manage"), financeHandler.DeleteScholarship)

		// Cash Accounts (Kas & Bank)
		protected.GET("/v1/finance/cash-accounts", financeHandler.GetCashAccounts)
		protected.POST("/v1/finance/cash-accounts", middleware.RequirePermission("finance.cashaccount.manage"), financeHandler.CreateCashAccount)
		protected.GET("/v1/finance/cash-accounts/:id/mutations", financeHandler.GetCashMutations)
		protected.POST("/v1/finance/cash-accounts/:id/mutations", middleware.RequirePermission("finance.cashaccount.manage"), financeHandler.CreateCashMutation)

// Journals & Buku Besar
		protected.GET("/v1/finance/journals", financeHandler.GetJournals)
		protected.GET("/v1/finance/journals/:id", financeHandler.GetJournalDetail)

		// Budgets (RAB)
		protected.GET("/v1/finance/budgets", financeHandler.GetBudgets)
		protected.POST("/v1/finance/budgets", middleware.RequirePermission("finance.budget.manage"), financeHandler.CreateBudget)
		protected.GET("/v1/finance/budgets/:id", financeHandler.GetBudgetDetail)

		// Vendors
		protected.GET("/v1/finance/vendors", financeHandler.GetVendors)
		protected.POST("/v1/finance/vendors", middleware.RequirePermission("finance.vendor.manage"), financeHandler.CreateVendor)

		// Purchase Orders
		protected.GET("/v1/finance/purchase-orders", financeHandler.GetPurchaseOrders)
		protected.POST("/v1/finance/purchase-orders", middleware.RequirePermission("finance.po.manage"), financeHandler.CreatePurchaseOrder)
		protected.POST("/v1/finance/purchase-orders/:id/approve", middleware.RequirePermission("finance.po.approve"), financeHandler.ApprovePurchaseOrder)

// Expense Events (Wisuda & Kegiatan)
		protected.GET("/v1/finance/events", financeHandler.GetExpenseEvents)
		protected.POST("/v1/finance/events", middleware.RequirePermission("finance.event.manage"), financeHandler.CreateExpenseEvent)

		// Payroll Runs
		protected.GET("/v1/finance/payroll-runs", financeHandler.GetPayrollRuns)
		protected.POST("/v1/finance/payroll-runs/:id/approve", middleware.RequirePermission("finance.payroll.approve"), financeHandler.ApprovePayrollRun)

		// Disbursements (CRM Commission)
		protected.GET("/v1/finance/disbursements", financeHandler.GetDisbursements)
		protected.POST("/v1/finance/disbursements", middleware.RequirePermission("finance.disbursement.manage"), financeHandler.CreateDisbursement)
		protected.POST("/v1/finance/disbursements/:id/approve", middleware.RequirePermission("finance.disbursement.approve"), financeHandler.ApproveDisbursement)
		protected.POST("/v1/finance/disbursements/:id/process", middleware.RequirePermission("finance.disbursement.process"), financeHandler.ProcessDisbursement)

// Reports (Keuangan)
		protected.GET("/v1/finance/reports/position", financeHandler.GetBalanceSheet)
		protected.GET("/v1/finance/reports/activity", financeHandler.GetIncomeStatement)
		protected.GET("/v1/finance/reports/cashflow", financeHandler.GetCashFlow)
	}

	// Outbox poller (background)
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	outboxPoller, err := infrastructure.NewOutboxPoller(db, amqpURL, 5*time.Second)
	if err != nil {
		sharedobservability.Logger.Warn().Err(err).Msg("Failed to initialize outbox poller - running without event publishing")
	} else {
		go outboxPoller.Start(context.Background())
		defer outboxPoller.Stop()
	}

	// Start event consumer in background
	go func() {
		if err := infrastructure.RunEventConsumer(db); err != nil {
			sharedobservability.Logger.Error().Err(err).Msg("Event consumer stopped with error")
		}
	}()

	// Graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8005"
	}

	go func() {
		sharedobservability.Logger.Info().Msgf("Finance Service started on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Server failed to run: %v", err)
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	sharedobservability.Logger.Info().Msg("Finance Service shutting down...")
}
