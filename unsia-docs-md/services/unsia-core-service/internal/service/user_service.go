package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type CreateUserInput struct {
	PersonID     string `json:"person_id" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	PersonID *string `json:"person_id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Status   *string `json:"status"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(input CreateUserInput) (*User, error) {
	// Check if username exists
	var existing User
	if err := s.db.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		return nil, errors.New("USERNAME_EXISTS: Username sudah digunakan")
	}

	// Check if email exists
	if err := s.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("EMAIL_EXISTS: Email sudah digunakan")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("PASSWORD_HASH_FAILED: Gagal mengenkripsi password")
	}

	user := User{
		ID:           uuid.New().String(),
		PersonID:     input.PersonID,
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal membuat user")
	}

	return &user, nil
}

// GetUserByID retrieves user by ID
func (s *UserService) GetUserByID(id string) (*User, error) {
	var user User
	if err := s.db.Preload("Person").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("USER_NOT_FOUND: User tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil user")
	}
	return &user, nil
}

// GetUserByUsername retrieves user by username
func (s *UserService) GetUserByUsername(username string) (*User, error) {
	var user User
	if err := s.db.Preload("Person").First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("USER_NOT_FOUND: User tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil user")
	}
	return &user, nil
}

// GetUserByEmail retrieves user by email
func (s *UserService) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := s.db.Preload("Person").First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("USER_NOT_FOUND: User tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal mengambil user")
	}
	return &user, nil
}

// UpdateUser updates user details
func (s *UserService) UpdateUser(id string, input UpdateUserInput) (*User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if input.PersonID != nil {
		user.PersonID = *input.PersonID
	}
	if input.Username != nil {
		// Check if new username exists
		var existing User
		if err := s.db.Where("username = ? AND id != ?", *input.Username, id).First(&existing).Error; err == nil {
			return nil, errors.New("USERNAME_EXISTS: Username sudah digunakan")
		}
		user.Username = *input.Username
	}
	if input.Email != nil {
		// Check if new email exists
		var existing User
		if err := s.db.Where("email = ? AND id != ?", *input.Email, id).First(&existing).Error; err == nil {
			return nil, errors.New("EMAIL_EXISTS: Email sudah digunakan")
		}
		user.Email = *input.Email
	}
	if input.Status != nil {
		user.Status = *input.Status
	}

	user.UpdatedAt = time.Now()

	if err := s.db.Save(user).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengupdate user")
	}

	return user, nil
}

// DeactivateUser soft deletes user by setting status to inactive
func (s *UserService) DeactivateUser(id string, reason string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Status = "inactive"
	user.UpdatedAt = time.Now()

	if err := s.db.Save(user).Error; err != nil {
		return errors.New("DB_ERROR: Gagal menonaktifkan user")
	}

	return nil
}

// ActivateUser activates user
func (s *UserService) ActivateUser(id string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Status = "active"
	user.UpdatedAt = time.Now()

	if err := s.db.Save(user).Error; err != nil {
		return errors.New("DB_ERROR: Gagal mengaktifkan user")
	}

	return nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(userID string, input ChangePasswordInput) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword)); err != nil {
		return errors.New("INVALID_OLD_PASSWORD: Password lama tidak cocok")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("PASSWORD_HASH_FAILED: Gagal mengenkripsi password baru")
	}

	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.db.Save(user).Error; err != nil {
		return errors.New("DB_ERROR: Gagal mengubah password")
	}

	return nil
}

// ListUsers returns paginated list of users
func (s *UserService) ListUsers(page, limit int, status string) ([]User, int64, error) {
	var users []User
	var total int64

	query := s.db.Model(&User{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("DB_ERROR: Gagal menghitung user")
	}

	offset := (page - 1) * limit
	if err := query.Preload("Person").Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, errors.New("DB_ERROR: Gagal mengambil daftar user")
	}

	return users, total, nil
}

// AssignRole assigns role to user
func (s *UserService) AssignRole(userID, roleID string, studyProgramID *string) (*UserRole, error) {
	// Verify user exists
	if _, err := s.GetUserByID(userID); err != nil {
		return nil, err
	}

	// Verify role exists
	var role Role
	if err := s.db.First(&role, "id = ?", roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ROLE_NOT_FOUND: Role tidak ditemukan")
		}
		return nil, errors.New("DB_ERROR: Gagal memverifikasi role")
	}

	// Check if assignment exists
	var existing UserRole
	if err := s.db.Where("user_id = ? AND role_id = ?", userID, roleID).First(&existing).Error; err == nil {
		return nil, errors.New("ROLE_ALREADY_ASSIGNED: Role sudah ditugaskan ke user")
	}

	userRole := UserRole{
		ID:             uuid.New().String(),
		UserID:         userID,
		RoleID:         roleID,
		StudyProgramID: studyProgramID,
		CreatedAt:      time.Now(),
	}

	if err := s.db.Create(&userRole).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal menugaskan role")
	}

	return &userRole, nil
}

// RemoveRole removes role from user
func (s *UserService) RemoveRole(userID, roleID string) error {
	result := s.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&UserRole{})
	if result.Error != nil {
		return errors.New("DB_ERROR: Gagal menghapus role")
	}
	if result.RowsAffected == 0 {
		return errors.New("ROLE_NOT_ASSIGNED: Role tidak ditugaskan ke user")
	}
	return nil
}

// GetUserRoles returns all roles assigned to user
func (s *UserService) GetUserRoles(userID string) ([]UserRole, error) {
	var userRoles []UserRole
	if err := s.db.Preload("Role").Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengambil role user")
	}
	return userRoles, nil
}

// GetRolePermissions returns all permissions for a role
func (s *UserService) GetRolePermissions(roleID string) ([]Permission, error) {
	var permissions []Permission
	if err := s.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, errors.New("DB_ERROR: Gagal mengambil permissions")
	}
	return permissions, nil
}
