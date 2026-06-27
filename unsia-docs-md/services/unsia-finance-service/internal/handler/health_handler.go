package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status        string            `json:"status"`
	Service      string            `json:"service"`
	Version      string            `json:"version"`
	Timestamp    time.Time         `json:"timestamp"`
	Uptime       time.Duration    `json:"uptime"`
	Dependencies DependencyStatus `json:"dependencies"`
}

// DependencyStatus shows the status of dependencies
type DependencyStatus struct {
	Database  ComponentHealth `json:"database"`
	MessageQueue ComponentHealth `json:"message_queue,omitempty"`
}

// ComponentHealth represents the health of a component
type ComponentHealth struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

var startTime = time.Now()

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()

	// Check database
	dbStatus := h.checkDatabase(ctx)

	response := HealthResponse{
		Status:    "healthy",
		Service:  "unsia-finance-service",
		Version:  "v1.0.0",
		Timestamp: time.Now(),
		Uptime:   time.Since(startTime),
		Dependencies: DependencyStatus{
			Database: dbStatus,
		},
	}

	// Determine overall status
	if dbStatus.Status != "healthy" {
		response.Status = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, sharederr.Error("SERVICE_UNAVAILABLE", "Service is unavailable").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(response).WithContext(c))
}

// LivenessCheck handles GET /health/live (simple liveness probe)
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// ReadinessCheck handles GET /health/ready (readiness probe)
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ctx := c.Request.Context()

	// Check database connection
	dbStatus := h.checkDatabase(ctx)

	if dbStatus.Status != "healthy" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"reason": "database not available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// checkDatabase checks database connectivity
func (h *HealthHandler) checkDatabase(ctx context.Context) ComponentHealth {
	start := time.Now()

	sqlDB, err := h.db.DB()
	if err != nil {
		return ComponentHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}

	// Check connection
	if err := sqlDB.PingContext(ctx); err != nil {
		return ComponentHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}

	// Get database stats
	var result int64
	if err := sqlDB.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return ComponentHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}

	latency := time.Since(start)

	return ComponentHealth{
		Status:  "healthy",
		Latency: latency.String(),
	}
}

// MetricsHandler handles Prometheus metrics
func (h *HealthHandler) MetricsHandler(c *gin.Context) {
	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	// Basic metrics
	metrics := `# HELP finance_service_up Whether the finance service is up
# TYPE finance_service_up gauge
finance_service_up 1

# HELP finance_db_connections Active database connections
# TYPE finance_db_connections gauge`

	sqlDB, err := h.db.DB()
	if err == nil {
		stats := sqlDB.Stats()
		metrics += fmt.Sprintf("\nfinance_db_connections %d", stats.InUse)
	}

	c.String(http.StatusOK, metrics)
}

// InitDB initializes the database connection
func InitDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Invoice{},
		&domain.InvoiceItem{},
		&domain.Payment{},
		&domain.PaymentGatewayCallback{},
		&domain.PaymentVerification{},
		&domain.StudentClearance{},
		&domain.ClearancePolicy{},
		&domain.InstallmentRequest{},
		&domain.CoaAccount{},
		&domain.Journal{},
		&domain.JournalEntry{},
	)
}

// RunMigrations runs golang-migrate migrations
func RunMigrations(databaseURL string) error {
	// This is a placeholder - typically you'd use golang-migrate CLI
	// or a migration tool like goose
	// For now, we'll use GORM's AutoMigrate
	db, err := InitDB(databaseURL)
	if err != nil {
		return err
	}

	return AutoMigrate(db)
}

// CheckPendingInvoices checks for expired invoices
func CheckPendingInvoices(db *gorm.DB) error {
	now := time.Now()
	
	// Update status to EXPIRED for issued invoices past due date
	result := db.Model(&domain.Invoice{}).
		Where("status = ? AND due_date < ?", "ISSUED", now).
		Updates(map[string]interface{}{
			"status":     "EXPIRED",
			"updated_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
