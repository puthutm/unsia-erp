package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/domain"
	"github.com/unsia-erp/unsia-core-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type RoleHandler struct {
	repo *repository.RoleRepository
}

type CreateRoleRequest struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ScopeType string `json:"scope_type" binding:"required,oneof=global study_program module self"`
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{
		repo: repository.NewRoleRepository(db),
	}
}

// Create handles POST /api/v1/admin/roles
func (h *RoleHandler) Create(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	// Validate actor active role: super_admin or admin_bppti
	if claims.ActiveRole != "super_admin" && claims.ActiveRole != "admin_bppti" {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya role super_admin atau admin_bppti yang diizinkan untuk membuat role").WithContext(c))
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if role code already exists
	existing, err := h.repo.GetByCode(req.Code)
	if err == nil && existing != nil {
		c.JSON(http.StatusConflict, sharederr.Error("ROLE_ALREADY_EXISTS", "Role dengan code tersebut sudah terdaftar").WithContext(c))
		return
	}

	role := domain.Role{
		Code:      req.Code,
		Name:      req.Name,
		ScopeType: req.ScopeType,
	}

	if err := h.repo.Create(&role); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan role baru").WithContext(c))
		return
	}

	// Audit Log
	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "core.role.create",
		Module:       "core",
		ResourceType: "role",
		ResourceID:   role.ID,
		NewValue:     role,
	})

	c.JSON(http.StatusCreated, sharederr.Success(role).WithContext(c))
}

// List handles GET /api/v1/admin/roles
func (h *RoleHandler) List(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	// Validate actor active role: super_admin or admin_bppti
	if claims.ActiveRole != "super_admin" && claims.ActiveRole != "admin_bppti" {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya role super_admin atau admin_bppti yang diizinkan untuk melihat daftar role").WithContext(c))
		return
	}

	roles, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar role").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(roles).WithContext(c))
}
