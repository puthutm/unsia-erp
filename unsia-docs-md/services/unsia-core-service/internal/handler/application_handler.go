package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"gorm.io/gorm"
)

type ApplicationHandler struct {
	db *gorm.DB
}

func NewApplicationHandler(db *gorm.DB) *ApplicationHandler {
	return &ApplicationHandler{db: db}
}

type Application struct {
	ID          string `json:"id"`
	Name       string `json:"name"`
	ClientID   string `json:"client_id"`
	ClientSecret string `json:"client_secret,omitempty"`
	RedirectURIs string `json:"redirect_uris"`
	GrantTypes string `json:"grant_types"`
	Scopes    string `json:"scopes"`
	IsActive  bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// CreateApplication handles POST /api/v1/applications
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req struct {
		Name         string   `json:"name" binding:"required"`
		RedirectURIs []string `json:"redirect_uris" binding:"required"`
		GrantTypes  []string `json:"grant_types"`
		Scopes     []string `json:"scopes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	app := Application{
		ID:            uuid.New().String(),
		Name:          req.Name,
		ClientID:      "app_" + uuid.New().String()[:16],
		ClientSecret:  uuid.New().String(),
		RedirectURIs:  toJSONArray(req.RedirectURIs),
		GrantTypes:   toJSONArray(req.GrantTypes),
		Scopes:       toJSONArray(req.Scopes),
		IsActive:     true,
		CreatedAt:    getTimestamp(),
		UpdatedAt:    getTimestamp(),
	}

	if err := h.db.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(app).WithContext(c))
}

// GetApplication handles GET /api/v1/applications/:id
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	id := c.Param("id")

	var app Application
	if err := h.db.Where("id = ?", id).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Application not found").WithContext(c))
		return
	}

	app.ClientSecret = ""
	c.JSON(http.StatusOK, sharederr.Success(app).WithContext(c))
}

// ListApplications handles GET /api/v1/applications
func (h *ApplicationHandler) ListApplications(c *gin.Context) {
	var apps []Application
	if err := h.db.Order("created_at DESC").Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	// Clear secrets
	for i := range apps {
		apps[i].ClientSecret = ""
	}

	c.JSON(http.StatusOK, sharederr.Success(apps).WithContext(c))
}

// UpdateApplication handles PUT /api/v1/applications/:id
func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name         *string  `json:"name"`
		RedirectURIs []string `json:"redirect_uris"`
		GrantTypes  []string `json:"grant_types"`
		Scopes     []string `json:"scopes"`
		IsActive   *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := map[string]interface{}{
		"updated_at": getTimestamp(),
	}

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if len(req.RedirectURIs) > 0 {
		updates["redirect_uris"] = toJSONArray(req.RedirectURIs)
	}
	if len(req.GrantTypes) > 0 {
		updates["grant_types"] = toJSONArray(req.GrantTypes)
	}
	if len(req.Scopes) > 0 {
		updates["scopes"] = toJSONArray(req.Scopes)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.db.Model(&Application{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Application updated",
	}).WithContext(c))
}

// DeleteApplication handles DELETE /api/v1/applications/:id
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&Application{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Application deleted",
	}).WithContext(c))
}

func toJSONArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	result := "["
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += `"` + s + `"`
	}
	result += "]"
	return result
}

func getTimestamp() string {
	return now().Format("2006-01-02T15:04:05Z07:00")
}
