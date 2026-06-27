package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/service"
	"gorm.io/gorm"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
	db             *gorm.DB
}

func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	return &WebhookHandler{
		webhookService: service.NewWebhookService(db),
		db:             db,
	}
}

// CreateWebhook handles POST /api/v1/webhooks
func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	var req service.CreateWebhookInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	webhook, err := h.webhookService.CreateWebhook(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(webhook).WithContext(c))
}

// ListWebhooks handles GET /api/v1/webhooks?event=...
func (h *WebhookHandler) ListWebhooks(c *gin.Context) {
	event := c.Query("event")

	webhooks, err := h.webhookService.ListWebhooks(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(webhooks).WithContext(c))
}

// DeleteWebhook handles DELETE /api/v1/webhooks/:id
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	id := c.Param("id")

	if err := h.webhookService.DeleteWebhook(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Webhook deleted successfully",
	}).WithContext(c))
}

// GetWebhook handles GET /api/v1/webhooks/:id
func (h *WebhookHandler) GetWebhook(c *gin.Context) {
	id := c.Param("id")
	var webhook service.Webhook
	if err := h.db.Where("id = ?", id).First(&webhook).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Webhook not found").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(webhook).WithContext(c))
}

// UpdateWebhook handles PUT /api/v1/webhooks/:id
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		URL      string `json:"url" binding:"required,url"`
		Event    string `json:"event" binding:"required"`
		Secret   string `json:"secret"`
		IsActive bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var webhook service.Webhook
	if err := h.db.Where("id = ?", id).First(&webhook).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Webhook not found").WithContext(c))
		return
	}

	webhook.URL = req.URL
	webhook.Event = req.Event
	webhook.Secret = req.Secret
	webhook.IsActive = req.IsActive
	webhook.UpdatedAt = time.Now()

	if err := h.db.Save(&webhook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate webhook").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(webhook).WithContext(c))
}

// TestWebhook handles POST /api/v1/webhooks/:id/test
func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	// For now, just return success
	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Test webhook triggered",
	}).WithContext(c))
}

// TriggerEvent handles POST /api/v1/webhooks/trigger
func (h *WebhookHandler) TriggerEvent(c *gin.Context) {
	var req struct {
		Event string      `json:"event" binding:"required"`
		Data  interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if err := h.webhookService.TriggerEvent(req.Event, req.Data); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Event triggered to webhooks",
	}).WithContext(c))
}
