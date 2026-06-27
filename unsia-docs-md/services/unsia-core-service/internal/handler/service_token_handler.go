package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"gorm.io/gorm"
)

type ServiceTokenHandler struct {
	db *gorm.DB
}

func NewServiceTokenHandler(db *gorm.DB) *ServiceTokenHandler {
	return &ServiceTokenHandler{db: db}
}

type ServiceToken struct {
	ID            string    `json:"id"`
	ApplicationID string  `json:"application_id"`
	TokenHash    string   `json:"token_hash,omitempty"`
	Scopes       string   `json:"scopes"`
	ExpiresAt    time.Time `json:"expired_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateServiceToken handles POST /api/v1/service-tokens
func (h *ServiceTokenHandler) CreateServiceToken(c *gin.Context) {
	var req struct {
		ApplicationID string   `json:"application_id" binding:"required"`
		Scopes       []string `json:"scopes"`
		ExpiresIn    int      `json:"expires_in"` // hours
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Generate token
	token := uuid.New().String() + "." + uuid.New().String()
	tokenHash := hashToken(token)

	expiresIn := req.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 8760 // 1 year default
	}

	st := ServiceToken{
		ID:            uuid.New().String(),
		ApplicationID: req.ApplicationID,
		TokenHash:     tokenHash,
		Scopes:       toJSONArray(req.Scopes),
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Hour),
		CreatedAt:    time.Now(),
	}

	if err := h.db.Create(&st).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"token":       token,
		"expires_at":  st.ExpiresAt,
	}).WithContext(c))
}

// ListServiceTokens handles GET /api/v1/service-tokens
func (h *ServiceTokenHandler) ListServiceTokens(c *gin.Context) {
	appID, _ := c.Get("application_id")

	var tokens []ServiceToken
	if err := h.db.Where("application_id = ?", appID).Order("created_at DESC").Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	for i := range tokens {
		tokens[i].TokenHash = ""
	}

	c.JSON(http.StatusOK, sharederr.Success(tokens).WithContext(c))
}

// GetServiceToken handles GET /api/v1/service-tokens/:id
func (h *ServiceTokenHandler) GetServiceToken(c *gin.Context) {
	id := c.Param("id")
	appID, _ := c.Get("application_id")

	var token ServiceToken
	if err := h.db.Where("id = ? AND application_id = ?", id, appID).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Token not found").WithContext(c))
		return
	}

	token.TokenHash = ""
	c.JSON(http.StatusOK, sharederr.Success(token).WithContext(c))
}

// RevokeServiceToken handles DELETE /api/v1/service-tokens/:id
func (h *ServiceTokenHandler) RevokeServiceToken(c *gin.Context) {
	id := c.Param("id")
	appID, _ := c.Get("application_id")

	var token ServiceToken
	if err := h.db.Where("id = ? AND application_id = ?", id, appID).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Token not found").WithContext(c))
		return
	}

	now := time.Now()
	if err := h.db.Model(&token).Update("revoked_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"message": "Token revoked",
	}).WithContext(c))
}

// ValidateServiceToken validates a service token
func (h *ServiceTokenHandler) ValidateServiceToken(token string) (*ServiceToken, error) {
	tokenHash := hashToken(token)
	var st ServiceToken
	err := h.db.Where("token_hash = ?", tokenHash).First(&st).Error
	if err != nil {
		return nil, err
	}

	if st.RevokedAt != nil || st.ExpiresAt.Before(time.Now()) {
		return nil, nil // Token revoked or expired
	}

	return &st, nil
}

func hashToken(token string) string {
	// Simple hash - in production use proper hashing
	return token[:64]
}

