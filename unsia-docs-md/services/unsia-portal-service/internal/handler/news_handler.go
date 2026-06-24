package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-portal-service/internal/domain"
	"github.com/unsia-erp/unsia-portal-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type NewsHandler struct {
	repo *repository.PortalRepository
	db   *gorm.DB
}

func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{
		repo: repository.NewPortalRepository(db),
		db:   db,
	}
}

// GET /api/v1/portal/news - List news
func (h *NewsHandler) ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	status := c.Query("status")

	news, total, err := h.repo.ListNews(page, limit, category, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data berita").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(news).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/portal/news - Create news
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content"`
		Summary  string `json:"summary"`
		Category string `json:"category"`
		ImageURL string `json:"image_url"`
		AuthorID string `json:"author_id"`
		Status   string `json:"status"` // draft, published, archived
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	news := domain.News{
		Title:    req.Title,
		Content:  req.Content,
		Summary:  req.Summary,
		Category: req.Category,
		ImageURL: req.ImageURL,
		AuthorID: req.AuthorID,
		Status:   req.Status,
	}

	if err := h.repo.CreateNews(&news); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat berita").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(news).WithContext(c))
}

// GET /api/v1/portal/news/:id - Get news
func (h *NewsHandler) GetNews(c *gin.Context) {
	newsID := c.Param("id")

	news, err := h.repo.GetNewsByID(newsID)
	if err != nil || news == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Berita tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(news).WithContext(c))
}

// PUT /api/v1/portal/news/:id - Update news
func (h *NewsHandler) UpdateNews(c *gin.Context) {
	newsID := c.Param("id")

	var req struct {
		Title    string `json:"title"`
		Content string `json:"content"`
		Summary string `json:"summary"`
		Category string `json:"category"`
		ImageURL string `json:"image_url"`
		Status  string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.repo.UpdateNews(newsID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui berita").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Berita berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/portal/news/:id - Delete news
func (h *NewsHandler) DeleteNews(c *gin.Context) {
	newsID := c.Param("id")

	if err := h.repo.DeleteNews(newsID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus berita").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Berita berhasil dihapus").WithContext(c))
}
