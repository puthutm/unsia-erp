package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/service"
	"gorm.io/gorm"
)

type ExternalAppHandler struct {
	externalAppService *service.ExternalAppService
}

func NewExternalAppHandler(db *gorm.DB) *ExternalAppHandler {
	return &ExternalAppHandler{
		externalAppService: service.NewExternalAppService(db),
	}
}

// CreateExternalApp handles POST /api/v1/external-apps
func (h *ExternalAppHandler) CreateExternalApp(c *gin.Context) {
	var req service.CreateExternalAppInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get created by from context
	if userID, exists := c.Get("user_id"); exists {
		req.CreatedBy = userID.(string)
	}

	app, err := h.externalAppService.CreateExternalApp(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	// Return with client_secret
	c.JSON(http.StatusCreated, sharederr.Success(app).WithContext(c))
}

// GetExternalApp handles GET /api/v1/external-apps/:id
func (h *ExternalAppHandler) GetExternalApp(c *gin.Context) {
	id := c.Param("id")

	app, err := h.externalAppService.GetExternalApp(id)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "External app not found").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(app).WithContext(c))
}

// ListExternalApps handles GET /api/v1/external-apps
func (h *ExternalAppHandler) ListExternalApps(c *gin.Context) {
	filter := service.ExternalAppFilter{
		Type:    c.Query("type"),
		Limit:   getQueryInt(c, "limit", 50),
		Offset:  getQueryInt(c, "offset", 0),
	}

	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filter.IsActive = &active
	}

	if isInternal := c.Query("is_internal"); isInternal != "" {
		internal := isInternal == "true"
		filter.IsInternal = &internal
	}

	apps, total, err := h.externalAppService.ListExternalApps(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"data":  apps,
		"total": total,
	}).WithContext(c))
}

// UpdateExternalApp handles PUT /api/v1/external-apps/:id
func (h *ExternalAppHandler) UpdateExternalApp(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateExternalAppInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	app, err := h.externalAppService.UpdateExternalApp(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(app).WithContext(c))
}

// DeleteExternalApp handles DELETE /api/v1/external-apps/:id
func (h *ExternalAppHandler) DeleteExternalApp(c *gin.Context) {
	id := c.Param("id")

	if err := h.externalAppService.DeactivateExternalApp(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "External app deactivated",
	}).WithContext(c))
}

// RegenerateSecret handles POST /api/v1/external-apps/:id/secret
func (h *ExternalAppHandler) RegenerateSecret(c *gin.Context) {
	id := c.Param("id")

	newSecret, err := h.externalAppService.RegenerateSecret(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"client_secret": newSecret,
	}).WithContext(c))
}

// ValidateCredentials handles POST /api/v1/external-apps/validate
func (h *ExternalAppHandler) ValidateCredentials(c *gin.Context) {
	var req struct {
		ClientID     string `json:"client_id" binding:"required"`
		ClientSecret string `json:"client_secret" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	app, err := h.externalAppService.ValidateClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	if app == nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_CREDENTIALS", "Invalid client credentials").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"valid": true,
		"app":   app,
	}).WithContext(c))
}

func getQueryInt(c *gin.Context, key string, defaultValue int) int {
	if val := c.Query(key); val != "" {
		var n int
		if _, err := fmt.Sscanf(val, "%d", &n); err == nil {
			return n
		}
	}
	return defaultValue
}
