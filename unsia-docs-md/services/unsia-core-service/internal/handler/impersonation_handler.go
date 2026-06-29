package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/domain"
	"github.com/unsia-erp/unsia-core-service/internal/infrastructure/keys"
	"github.com/unsia-erp/unsia-core-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type ImpersonationHandler struct {
	userRepo *repository.UserRepository
	sessRepo *repository.SessionRepository
	appRepo  *repository.ApplicationRepository
}

type ImpersonationStartRequest struct {
	TargetUserID    string `json:"target_user_id" binding:"required"`
	Reason          string `json:"reason" binding:"required"`
	DurationMinutes int    `json:"duration_minutes"`
}

func NewImpersonationHandler(db *gorm.DB) *ImpersonationHandler {
	return &ImpersonationHandler{
		userRepo: repository.NewUserRepository(db),
		sessRepo: repository.NewSessionRepository(db),
		appRepo:  repository.NewApplicationRepository(db),
	}
}

func (h *ImpersonationHandler) Start(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	// Validate actor has the required active role: admin_bppti
	if claims.ActiveRole != "admin_bppti" {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Hanya role admin_bppti yang diizinkan untuk memulai penyamaran (impersonation)").WithContext(c))
		return
	}

	var req ImpersonationStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// 1. Fetch and validate target user
	targetUser, err := h.userRepo.GetByID(req.TargetUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("USER_NOT_FOUND", "User target tidak ditemukan").WithContext(c))
		return
	}

	if targetUser.Status != "active" {
		c.JSON(http.StatusForbidden, sharederr.Error("USER_NOT_ACTIVE", "User target tidak aktif").WithContext(c))
		return
	}

	// 2. Fetch target user roles and permissions
	targetUserRoles, err := h.userRepo.GetUserRoles(targetUser.ID)
	if err != nil || len(targetUserRoles) == 0 {
		c.JSON(http.StatusForbidden, sharederr.Error("NO_ROLES_ASSIGNED", "User target tidak memiliki role yang ditugaskan").WithContext(c))
		return
	}

	// Choose first assigned role as default active role
	targetUserRole := targetUserRoles[0]
	targetRole := targetUserRole.Role

	// Resolve target scope
	targetScope := targetUserRole.Role.ScopeType
	if targetUserRole.Role.ScopeType == "study_program" && targetUserRole.StudyProgramID != nil {
		targetScope = "study_program:" + *targetUserRole.StudyProgramID
	}

	// Resolve target permissions
	permissions, err := h.userRepo.GetRolePermissions(targetUserRole.RoleID)
	if err != nil {
		permissions = []domain.Permission{}
	}
	targetPermissionCodes := make([]string, len(permissions))
	for i, p := range permissions {
		targetPermissionCodes[i] = p.Code
	}

	// Default duration is 30 minutes if not specified or invalid
	durationMinutes := req.DurationMinutes
	if durationMinutes <= 0 {
		durationMinutes = 30
	}

	// 3. Create active session record in database for the target user
	randBytes := make([]byte, 32)
	_, _ = rand.Read(randBytes)
	refreshToken := base64.URLEncoding.EncodeToString(randBytes)

	// Hash tokens for storage
	hashRefresh := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := hex.EncodeToString(hashRefresh[:])
	// Use a placeholder for access token hash; it will be set after JWT generation
	tokenHashPlaceholder := hex.EncodeToString(sha256.New().Sum(nil))

	session := domain.Session{
		UserID:           targetUser.ID,
		TokenHash:        tokenHashPlaceholder,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(time.Duration(durationMinutes) * time.Minute),
	}
	if err := h.sessRepo.CreateSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SESSION_CREATION_FAILED", "Gagal membuat sesi target").WithContext(c))
		return
	}

	// 4. Create impersonation session record
	impSession := domain.ImpersonationSession{
		ActorUserID:  claims.Subject,
		TargetUserID: targetUser.ID,
		TargetRoleID: targetRole.ID,
		SessionID:    session.ID,
		Reason:       req.Reason,
		StartedAt:    time.Now(),
		ExpiredAt:    time.Now().Add(time.Duration(durationMinutes) * time.Minute),
		Status:       "active",
	}
	if err := h.sessRepo.CreateImpersonationSession(&impSession); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("IMPERSONATION_SESSION_FAILED", "Gagal menyimpan data penyamaran").WithContext(c))
		return
	}

	// Save active role session to track active state
	activeSession := domain.ActiveRoleSession{
		UserID:         targetUser.ID,
		RoleID:         targetRole.ID,
		SessionID:      session.ID,
		StudyProgramID: targetUserRole.StudyProgramID,
	}
	_ = h.sessRepo.SaveActiveRoleSession(&activeSession)

	// 5. Generate Target Access Token using RSA keys
	claimsTarget := sharedauth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   targetUser.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(durationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "unsia-core-service",
		},
		ActiveRole:  targetRole.Code,
		Permissions: targetPermissionCodes,
		Scope:       targetScope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsTarget)
	token.Header["kid"] = keys.KeyID
	accessToken, err := token.SignedString(keys.SigningKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("TOKEN_GENERATION_FAILED", "Gagal menerbitkan access token penyamaran").WithContext(c))
		return
	}

	// 6. Write Audit Log
	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "core.impersonation.start",
		Module:       "core",
		ResourceType: "impersonation_session",
		ResourceID:   impSession.ID,
		NewValue:     impSession,
	})

	// 7. Return success with target user profile and token
	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"user_id":      targetUser.ID,
		"person_id":    targetUser.PersonID,
		"name":         targetUser.Person.Name,
		"email":        targetUser.Person.Email,
		"active_role":  targetRole.Code,
		"permissions":  targetPermissionCodes,
		"scope":        targetScope,
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int64(durationMinutes * 60),
	}).WithContext(c))
}
