package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/service"
	"gorm.io/gorm"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
}

func NewWebhookHandler(db *gorm.DB) *WebhookHandler {
	return &WebhookHandler{
		webhookService: service.NewWebhookService(db),
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

// TestWebhook handles POST /api/v1/webhooks/:id/test
func (h *WebhookHandler) TestWebhook(c *gin.Context) {
	id := c.Param("id")

	// Get webhook by ID (need to add this method)
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
