package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"gorm.io/gorm"
)

type MaterialCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	FileURL     string `json:"file_url" binding:"required"`
}

type MaterialHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

func NewMaterialHandler(db *gorm.DB) *MaterialHandler {
	return &MaterialHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// CreateMaterial - Membuat materi perkuliahan baru
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
	sessionID := c.Param("id")
	var req MaterialCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	now := time.Now()
	material := domain.Material{
		SessionID:   sessionID,
		Title:       req.Title,
		ContentType: req.ContentType,
		FileURL:     req.FileURL,
		PublishedAt: &now,
	}

	if err := h.repo.CreateMaterial(&material); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat materi perkuliahan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(material).WithContext(c))
}
