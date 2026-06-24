package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

type CreateRoleInput struct {
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ScopeType string `json:"scope_type" binding:"required"` // global, prodi, module, self
}

type UpdateRoleInput struct {
	Name      *string `json:"name"`
	ScopeType *string `json:"scope_type"`
	IsActive  *bool   `json:"is_active"`
}

type PermissionInput struct {
	Code      string `json:"code" binding:"required"` // e.g. pmb.applicant.verify
	Name      string `json:"name" binding:"required"`
	Module   string `json:"module" binding:"required"`
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(input CreateRoleInput) (*Role, error) {
	// Check if code exists
	var existing Role
	if err := s.db.Where("code = ?", input.Code).First(&existing).Error; err == nil {
		return nil, errors.New("ROLE_CODE_EXISTS: Kode role sudah digunakan")
	}

	// Validate scope type
	validScopes := map[string]bool{
		"global":        true,
		"prodi":        true,
		"module":       true,
		"self":         true,
		"study_program": true,
	}
	if !validScopes[input.ScopeType] {
		return nil, errors.New("INVALID_SCOPE_TYPE: Tipe scope tidak valid")
	}

	role := Role{
		ID:        uuid.New().String(),
		Code:     input.Code,
		Name:     input.Name,
		ScopeType: input.ScopeType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal membuat role")
	}

	return &role, nil
}

// GetRoleByID retrieves role by ID
func (s *RoleService) GetRoleByID(id string) (*Role, error) {
	var role Role
	if err := s.db.First(&role, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ROLE_NOT_FOUND: Role tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil role")
	}
	return &role, nil
}

// GetRoleByCode retrieves role by code
func (s *RoleService) GetRoleByCode(code string) (*Role, error) {
	var role Role
	if err := s.db.Where("code = ?", code).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ROLE_NOT_FOUND: Role tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil role")
	}
	return &role, nil
}

// UpdateRole updates role details
func (s *RoleService) UpdateRole(id string, input UpdateRoleInput) (*Role, error) {
	role, err := s.GetRoleByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		role.Name = *input.Name
	}
	if input.ScopeType != nil {
		validScopes := map[string]bool{
			"global":        true,
			"prodi":         true,
			"module":        true,
			"self":          true,
			"study_program": true,
		}
		if !validScopes[*input.ScopeType] {
			return nil, errors.New("INVALID_SCOPE_TYPE: Tipe scope tidak valid")
		}
		role.ScopeType = *input.ScopeType
	}

	role.UpdatedAt = time.Now()

	if err := s.db.Save(role).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengupdate role")
	}

	return role, nil
}

// DeleteRole soft deletes role by setting is_active to false
func (s *RoleService) DeleteRole(id string) error {
	role, err := s.GetRoleByID(id)
	if err != nil {
		return err
	}

	role.UpdatedAt = time.Now()

	if err := s.db.Save(role).Error; err != nil {
		return errors.New("DB_ERROR: Gagal menghapus role")
	}

	return nil
}

// ListRoles returns paginated list of roles
func (s *RoleService) ListRoles(page, limit int) ([]Role, int64, error) {
	var roles []Role
	var total int64

	if err := s.db.Model(&Role{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("DB_ERROR: Gagal menghitung role")
	}

	offset := (page - 1) * limit
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, 0, errors.New("DB_ERROR: Gagal mengambil daftar role")
	}

	return roles, total, nil
}

// GetRolePermissions returns all permissions for a role
func (s *RoleService) GetRolePermissions(roleID string) ([]Permission, error) {
	var permissions []Permission
	if err := s.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengambil permissions")
	}
	return permissions, nil
}

// AssignPermission assigns permission to role
func (s *RoleService) AssignPermission(roleID, permissionID string) (*RolePermission, error) {
	// Verify role exists
	if _, err := s.GetRoleByID(roleID); err != nil {
		return nil, err
	}

	// Verify permission exists
	var perm Permission
	if err := s.db.First(&perm, "id = ?", permissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PERMISSION_NOT_FOUND: Permission tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal memverifikasi permission")
	}

	// Check if assignment exists
	var existing RolePermission
	if err := s.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&existing).Error; err == nil {
		return nil, errors.New("PERMISSION_ALREADY_ASSIGNED: Permission sudah ditugaskan ke role")
	}

	rolePerm := RolePermission{
		ID:           uuid.New().String(),
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(&rolePerm).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal menugaskan permission")
	}

	return &rolePerm, nil
}

// RemovePermission removes permission from role
func (s *RoleService) RemovePermission(roleID, permissionID string) error {
	result := s.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermission{})
	if result.Error != nil {
		return errors.New("DB_ERROR: Gagal menghapus permission")
	}
	if result.RowsAffected == 0 {
		return errors.New("PERMISSION_NOT_ASSIGNED: Permission tidak ditugaskan ke role")
	}
	return nil
}

// BulkAssignPermissions assigns multiple permissions to role
func (s *RoleService) BulkAssignPermissions(roleID string, permissionIDs []string) ([]RolePermission, error) {
	rolePerms := make([]RolePermission, 0, len(permissionIDs))

	for _, permissionID := range permissionIDs {
		// Verify permission exists
		var perm Permission
		if err := s.db.First(&perm, "id = ?", permissionID).Error; err != nil {
			continue // Skip invalid permissions
		}

		// Check if assignment exists
		var existing RolePermission
		if err := s.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&existing).Error; err != nil {
			rolePerm := RolePermission{
				ID:           uuid.New().String(),
				RoleID:       roleID,
				PermissionID: permissionID,
				CreatedAt:    time.Now(),
			}
			rolePerms = append(rolePerms, rolePerm)
		}
	}

	if len(rolePerms) > 0 {
		if err := s.db.Create(&rolePerms).Error; err != nil {
			return nil, errors.New("DB_ERROR: Gagal menugaskan permissions")
		}
	}

	return rolePerms, nil
}

// =============================================================================
// PERMISSION SERVICE
// =============================================================================

type PermissionService struct {
	db *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// CreatePermission creates a new permission
func (s *PermissionService) CreatePermission(input PermissionInput) (*Permission, error) {
	// Check if code exists
	var existing Permission
	if err := s.db.Where("code = ?", input.Code).First(&existing).Error; err == nil {
		return nil, errors.New("PERMISSION_CODE_EXISTS: Kode permission sudah digunakan")
	}

	perm := Permission{
		ID:        uuid.New().String(),
		Code:      input.Code,
		Name:     input.Name,
		Module:   input.Module,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(&perm).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal membuat permission")
	}

	return &perm, nil
}

// GetPermissionByID retrieves permission by ID
func (s *PermissionService) GetPermissionByID(id string) (*Permission, error) {
	var perm Permission
	if err := s.db.First(&perm, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PERMISSION_NOT_FOUND: Permission tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil permission")
	}
	return &perm, nil
}

// GetPermissionByCode retrieves permission by code
func (s *PermissionService) GetPermissionByCode(code string) (*Permission, error) {
	var perm Permission
	if err := s.db.Where("code = ?", code).First(&perm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("PERMISSION_NOT_FOUND: Permission tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil permission")
	}
	return &perm, nil
}

// ListPermissions returns paginated list of permissions
func (s *PermissionService) ListPermissions(page, limit int, module string) ([]Permission, int64, error) {
	var permissions []Permission
	var total int64

	query := s.db.Model(&Permission{})
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("DB_ERROR: Gagal menghitung permission")
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&permissions).Error; err != nil {
		return nil, 0, errors.New("DB_ERROR: Gagal mengambil daftar permission")
	}

	return permissions, total, nil
}

// DeletePermission soft deletes permission
func (s *PermissionService) DeletePermission(id string) error {
	perm, err := s.GetPermissionByID(id)
	if err != nil {
		return err
	}

	perm.UpdatedAt = time.Now()

	if err := s.db.Save(perm).Error; err != nil {
		return errors.New("DB_ERROR: Gagal menghapus permission")
	}

	return nil
}

// GetPermissionsByCodes returns multiple permissions by codes
func (s *PermissionService) GetPermissionsByCodes(codes []string) ([]Permission, error) {
	var permissions []Permission
	if err := s.db.Where("code IN ?", codes).Find(&permissions).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengambil permissions")
	}
	return permissions, nil
}
