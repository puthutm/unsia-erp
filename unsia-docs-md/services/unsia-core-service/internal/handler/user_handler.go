package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Status   string `json:"status"` // active, inactive, suspended
}

type UserUpdateRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
}

type UserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive suspended"`
	Reason string `json:"reason"`
}

type ScopeAssignmentRequest struct {
	StudyProgramID *string `json:"study_program_id"`
}

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// Create handles POST /api/v1/admin/users
func (h *UserHandler) Create(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	// Validate actor active role: super_admin or admin_bppti
	if claims.ActiveRole != "super_admin" && claims.ActiveRole != "admin_bppti" {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya role super_admin atau admin_bppti yang diizinkan untuk membuat user").WithContext(c))
		return
	}

	var req UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if username already exists
	var count int64
	h.db.Model(&domain.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, sharederr.Error("USERNAME_ALREADY_EXISTS", "Username sudah digunakan").WithContext(c))
		return
	}

	// Check if email already exists
	h.db.Model(&domain.Person{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, sharederr.Error("EMAIL_ALREADY_EXISTS", "Email sudah digunakan").WithContext(c))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal melakukan hash password").WithContext(c))
		return
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	person := domain.Person{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	user := domain.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Status:       status,
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&person).Error; err != nil {
			return err
		}
		user.PersonID = person.ID
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan user baru").WithContext(c))
		return
	}

	user.Person = person

	// Audit Log
	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "core.user.create",
		Module:       "core",
		ResourceType: "user",
		ResourceID:   user.ID,
		NewValue:     user,
	})

	c.JSON(http.StatusCreated, sharederr.Success(user).WithContext(c))
}

// Update handles PUT /api/v1/admin/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	if claims.ActiveRole != "super_admin" && claims.ActiveRole != "admin_bppti" {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya role super_admin atau admin_bppti yang diizinkan untuk mengubah user").WithContext(c))
		return
	}

	userID := c.Param("id")
	var user domain.User
	if err := h.db.Preload("Person").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "User tidak ditemukan").WithContext(c))
		return
	}

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check username conflict
	var count int64
	h.db.Model(&domain.User{}).Where("username = ? AND id != ?", req.Username, userID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, sharederr.Error("USERNAME_ALREADY_EXISTS", "Username sudah digunakan oleh user lain").WithContext(c))
		return
	}

	oldUser := user

	err := h.db.Transaction(func(tx *gorm.DB) error {
		user.Username = req.Username
		user.Person.Name = req.Name
		user.Person.Email = req.Email
		user.Person.Phone = req.Phone

		if err := tx.Save(&user.Person).Error; err != nil {
			return err
		}
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate user").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "core.user.update",
		Module:       "core",
		ResourceType: "user",
		ResourceID:   user.ID,
		OldValue:     oldUser,
		NewValue:     user,
	})

	c.JSON(http.StatusOK, sharederr.Success(user).WithContext(c))
}

// UpdateStatus handles PATCH /api/v1/admin/users/:id/status
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	if claims.ActiveRole != "super_admin" && claims.ActiveRole != "admin_bppti" {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya role super_admin atau admin_bppti yang diizinkan untuk mengubah status user").WithContext(c))
		return
	}

	userID := c.Param("id")
	var user domain.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "User tidak ditemukan").WithContext(c))
		return
	}

	var req UserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if req.Status != "active" && req.Reason == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("REASON_REQUIRED", "Alasan penonaktifan wajib diisi").WithContext(c))
		return
	}

	oldUser := user
	user.Status = req.Status
	user.UpdatedAt = time.Now()

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate status user").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "core.user.status_update",
		Module:       "core",
		ResourceType: "user",
		ResourceID:   user.ID,
		OldValue:     oldUser,
		NewValue:     user,
		Reason:       req.Reason,
	})

	c.JSON(http.StatusOK, sharederr.Success(user).WithContext(c))
}

// AssignScope handles POST /api/v1/admin/user-roles/:id/scopes
func (h *UserHandler) AssignScope(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	// Validate actor active role: super_admin or admin_bppti
	if claims.ActiveRole != "super_admin" && claims.ActiveRole != "admin_bppti" {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya role super_admin atau admin_bppti yang diizinkan untuk mengatur data scope").WithContext(c))
		return
	}

	userRoleID := c.Param("id")
	var userRole domain.UserRole
	if err := h.db.Where("id = ?", userRoleID).First(&userRole).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "User role assignment tidak ditemukan").WithContext(c))
		return
	}

	var req ScopeAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	oldUserRole := userRole
	userRole.StudyProgramID = req.StudyProgramID

	if err := h.db.Save(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengupdate scope prodi").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "core.user_role.scope_assign",
		Module:       "core",
		ResourceType: "user_role",
		ResourceID:   userRole.ID,
		OldValue:     oldUserRole,
		NewValue:     userRole,
	})

	c.JSON(http.StatusOK, sharederr.Success(userRole).WithContext(c))
}

// List handles GET /api/v1/users
func (h *UserHandler) List(c *gin.Context) {
	var users []domain.User
	if err := h.db.Preload("Person").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar user").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(users).WithContext(c))
}

// Get handles GET /api/v1/users/:id
func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var user domain.User
	if err := h.db.Preload("Person").Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "User tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(user).WithContext(c))
}

// Delete handles DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&domain.User{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus user").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(gin.H{"message": "User deleted"}).WithContext(c))
}

// Logout handles POST /api/v1/auth/logout
func (h *UserHandler) Logout(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if exists {
		claims := claimsVal.(*sharedauth.Claims)
		h.db.Model(&domain.Session{}).Where("user_id = ?", claims.Subject).Update("is_revoked", true)
	}
	c.JSON(http.StatusOK, sharederr.Success(gin.H{"message": "Logout successful"}).WithContext(c))
}

// ChangePassword handles POST /api/v1/auth/change-password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var user domain.User
	if err := h.db.Where("id = ?", claims.Subject).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("USER_NOT_FOUND", "User tidak ditemukan").WithContext(c))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_OLD_PASSWORD", "Password lama salah").WithContext(c))
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal mengenkripsi password baru").WithContext(c))
		return
	}

	if err := h.db.Model(&user).Update("password_hash", string(newHash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui password").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"message": "Password updated successfully"}).WithContext(c))
}

// ActivateUser handles POST /api/v1/users/:id/activate
func (h *UserHandler) ActivateUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Model(&domain.User{}).Where("id = ?", id).Update("status", "active").Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengaktifkan user").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(gin.H{"message": "User activated"}).WithContext(c))
}

// DeactivateUser handles POST /api/v1/users/:id/deactivate
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Model(&domain.User{}).Where("id = ?", id).Update("status", "inactive").Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menonaktifkan user").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(gin.H{"message": "User deactivated"}).WithContext(c))
}
