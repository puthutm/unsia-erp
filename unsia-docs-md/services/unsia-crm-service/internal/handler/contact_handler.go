package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-crm-service/internal/domain"
	"github.com/unsia-erp/unsia-crm-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type ContactHandler struct {
	repo *repository.CrmRepository
	db   *gorm.DB
}

func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{
		repo: repository.NewCrmRepository(db),
		db:   db,
	}
}

// GET /api/v1/contacts - List contacts
func (h *ContactHandler) ListContacts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	source := c.Query("source")

	contacts, total, err := h.repo.ListContacts(page, limit, status, source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data kontak").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(contacts).WithPagination(page, limit, int(total)).WithContext(c))
}

// POST /api/v1/contacts - Create contact
func (h *ContactHandler) CreateContact(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Email       string  `json:"email"`
		Phone       string  `json:"phone"`
		Company     string  `json:"company"`
		Position    string  `json:"position"`
		Source     string  `json:"source"`
		Note       string  `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	contact := domain.Contact{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Company:  req.Company,
		Position: req.Position,
		Source:   req.Source,
		Note:     req.Note,
		Status:   "active",
	}

	if err := h.repo.CreateContact(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat kontak").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(contact).WithContext(c))
}

// GET /api/v1/contacts/:id - Get contact
func (h *ContactHandler) GetContact(c *gin.Context) {
	contactID := c.Param("id")

	contact, err := h.repo.GetContactByID(contactID)
	if err != nil || contact == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kontak tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(contact).WithContext(c))
}

// PUT /api/v1/contacts/:id - Update contact
func (h *ContactHandler) UpdateContact(c *gin.Context) {
	contactID := c.Param("id")

	var req struct {
		Name    string `json:"name"`
		Email  string `json:"email"`
		Phone string `json:"phone"`
		Company string `json:"company"`
		Position string `json:"position"`
		Status string `json:"status"`
		Note   string `json:"note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if req.Position != "" {
		updates["position"] = req.Position
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Note != "" {
		updates["note"] = req.Note
	}

	if err := h.repo.UpdateContact(contactID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui kontak").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Kontak berhasil diperbarui").WithContext(c))
}

// DELETE /api/v1/contacts/:id - Delete contact
func (h *ContactHandler) DeleteContact(c *gin.Context) {
	contactID := c.Param("id")

	if err := h.repo.DeleteContact(contactID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus kontak").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Kontak berhasil dihapus").WithContext(c))
}
