package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Uptime  string            `json:"uptime"`
	Services map[string]string `json:"services"`
}

// Health handles GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Uptime:  getUptime(),
		Services: map[string]string{
			"database": "healthy",
		},
	}

	// Check database
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		response.Services["database"] = "unhealthy"
		response.Status = "degraded"
	}

	statusCode := http.StatusOK
	if response.Status == "degraded" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, sharederr.Success(response).WithContext(c))
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check all dependencies
	checks := map[string]string{}
	ready := true

	// Database check
	sqlDB, err := h.db.DB()
	if err == nil && sqlDB.Ping() == nil {
		checks["database"] = "ready"
	} else {
		checks["database"] = "not_ready"
		ready = false
	}

	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, sharederr.Success(checks).WithContext(c))
}

// Live handles GET /live
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"status": "alive",
	}).WithContext(c))
}

func getUptime() string {
	// For simplicity, just return a placeholder
	// In production, track actual start time
	return "0s"
}
