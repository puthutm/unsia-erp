package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-portal-service/internal/domain"
	"github.com/unsia-erp/unsia-portal-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type PortalContentHandler struct {
	repo *repository.PortalRepository
}

func NewPortalContentHandler(db *gorm.DB) *PortalContentHandler {
	return &PortalContentHandler{
		repo: repository.NewPortalRepository(db),
	}
}

type NewsCreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author"`
}

type AnnouncementCreateRequest struct {
	Title      string `json:"title" binding:"required"`
	Message    string `json:"message" binding:"required"`
	TargetRole string `json:"target_role"`
}

type EventCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	EventDate   string `json:"event_date" binding:"required"`
}

func (h *PortalContentHandler) ListNews(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 10
	offset := 0

	if val, err := strconv.Atoi(limitStr); err == nil {
		limit = val
	}
	if val, err := strconv.Atoi(offsetStr); err == nil {
		offset = val
	}

	list, total, err := h.repo.ListNews(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data berita").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"news":  list,
		"total": total,
	}).WithContext(c))
}

func (h *PortalContentHandler) CreateNews(c *gin.Context) {
	var req NewsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	n := domain.News{
		Title:       req.Title,
		Content:     req.Content,
		Author:      req.Author,
		PublishedAt: time.Now(),
	}

	if err := h.repo.CreateNews(&n); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan berita").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(n).WithContext(c))
}

func (h *PortalContentHandler) ListAnnouncements(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 10
	offset := 0

	if val, err := strconv.Atoi(limitStr); err == nil {
		limit = val
	}
	if val, err := strconv.Atoi(offsetStr); err == nil {
		offset = val
	}

	// Filter based on active role in JWT claims
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	role := "all"
	if claims != nil && claims.ActiveRole != "" {
		role = claims.ActiveRole
	}

	list, total, err := h.repo.ListAnnouncements(role, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"announcements": list,
		"total":         total,
	}).WithContext(c))
}

func (h *PortalContentHandler) CreateAnnouncement(c *gin.Context) {
	var req AnnouncementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	target := req.TargetRole
	if target == "" {
		target = "all"
	}

	a := domain.Announcement{
		Title:      req.Title,
		Message:    req.Message,
		TargetRole: target,
		CreatedAt:  time.Now(),
	}

	if err := h.repo.CreateAnnouncement(&a); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan pengumuman").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(a).WithContext(c))
}

func (h *PortalContentHandler) ListEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 10
	offset := 0

	if val, err := strconv.Atoi(limitStr); err == nil {
		limit = val
	}
	if val, err := strconv.Atoi(offsetStr); err == nil {
		offset = val
	}

	list, total, err := h.repo.ListEvents(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil agenda event").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"events": list,
		"total":  total,
	}).WithContext(c))
}

func (h *PortalContentHandler) CreateEvent(c *gin.Context) {
	var req EventCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("VALIDATION_ERROR", "EventDate format must be YYYY-MM-DD").WithContext(c))
		return
	}

	e := domain.Event{
		Title:       req.Title,
		Description: req.Description,
		EventDate:   parsedDate,
		CreatedAt:   time.Now(),
	}

	if err := h.repo.CreateEvent(&e); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan event").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(e).WithContext(c))
}
