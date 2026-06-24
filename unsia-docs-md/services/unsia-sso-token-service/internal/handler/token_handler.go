package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unsia-erp/unsia-sso-token-service/internal/service"
)

type TokenHandler struct {
	tokenService *service.TokenService
}

func NewTokenHandler(tokenService *service.TokenService) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *TokenHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request", "errors": []string{err.Error()}})
		return
	}

	accessToken, err := h.tokenService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token refreshed",
		"data": gin.H{
			"access_token": accessToken,
			"token_type": "Bearer",
		},
	})
}

type RevokeTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *TokenHandler) RevokeToken(c *gin.Context) {
	var req RevokeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request", "errors": []string{err.Error()}})
		return
	}

	err := h.tokenService.RevokeToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to revoke token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token revoked"})
}
