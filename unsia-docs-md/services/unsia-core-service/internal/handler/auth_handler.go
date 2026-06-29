package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-core-service/internal/domain"
	"github.com/unsia-erp/unsia-core-service/internal/infrastructure/keys"
	"github.com/unsia-erp/unsia-core-service/internal/infrastructure/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	CaptchaToken string `json:"captcha_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type SwitchRoleRequest struct {
	RoleCode   string `json:"role_code" binding:"required"`
	ScopeValue string `json:"scope_value"` // e.g. study_program UUID
}

type JWKKey struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKSResponse struct {
	Keys []JWKKey `json:"keys"`
}

type AuthHandler struct {
	userRepo *repository.UserRepository
	sessRepo *repository.SessionRepository
	appRepo  *repository.ApplicationRepository
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(db),
		sessRepo: repository.NewSessionRepository(db),
		appRepo:  repository.NewApplicationRepository(db),
	}
}

// Login verifies credentials, issues tokens, and saves session
// Supports login with NIM, NIP, Email, or Username
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Try to find user by NIM, NIK, or Email
	var user *domain.User
	var err error

	// Check if input contains @ (likely email)
	if contains := func(s string) bool {
		for _, r := range s {
			if r == '@' {
				return true
			}
		}
		return false
	}(req.Username); contains {
		// Try email first
		user, err = h.userRepo.GetByEmail(req.Username)
} else {
		// Try NIM or NIP - check which one matches
		user, err = h.userRepo.GetByNIM(req.Username)
		if err != nil {
			user, err = h.userRepo.GetByNIP(req.Username)
		}
	}

	// If still not found, try username as fallback
	if err != nil {
		user, err = h.userRepo.GetByUsername(req.Username)
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_CREDENTIALS", "NIM, NIK, Email, atau password salah").WithContext(c))
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusForbidden, sharederr.Error("USER_SUSPENDED", "Akun Anda dinonaktifkan").WithContext(c))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_CREDENTIALS", "Username atau password salah").WithContext(c))
		return
	}

	// Fetch user roles
	userRoles, err := h.userRepo.GetUserRoles(user.ID)
	if err != nil || len(userRoles) == 0 {
		c.JSON(http.StatusForbidden, sharederr.Error("NO_ROLES_ASSIGNED", "Akun Anda tidak memiliki role yang ditugaskan").WithContext(c))
		return
	}

	// Default to first assigned role
	defaultUserRole := userRoles[0]
	activeRole := defaultUserRole.Role.Code

	// Resolve scope
	scope := defaultUserRole.Role.ScopeType
	if defaultUserRole.Role.ScopeType == "study_program" && defaultUserRole.StudyProgramID != nil {
		scope = "study_program:" + *defaultUserRole.StudyProgramID
	}

	// Fetch permissions for this role
	permissions, err := h.userRepo.GetRolePermissions(defaultUserRole.RoleID)
	if err != nil {
		permissions = []domain.Permission{}
	}
	permissionCodes := make([]string, len(permissions))
	for i, p := range permissions {
		permissionCodes[i] = p.Code
	}

	// Generate JWT Access Token
	expireHours := 24
	if envHours := os.Getenv("JWT_EXPIRE_HOURS"); envHours != "" {
		if val, err := strconv.Atoi(envHours); err == nil {
			expireHours = val
		}
	}
	accessToken, err := keys.GenerateAccessToken(user.ID, activeRole, scope, permissionCodes, expireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("TOKEN_GENERATION_FAILED", "Gagal menerbitkan access token").WithContext(c))
		return
	}

	// Generate Refresh Token
	randBytes := make([]byte, 32)
	_, _ = rand.Read(randBytes)
	refreshToken := base64.URLEncoding.EncodeToString(randBytes)

	// Save Session
	expireDays := 7
	if envDays := os.Getenv("JWT_REFRESH_EXPIRE_DAYS"); envDays != "" {
		if val, err := strconv.Atoi(envDays); err == nil {
			expireDays = val
		}
	}
	hashAccess := sha256.Sum256([]byte(accessToken))
	accessTokenHash := hex.EncodeToString(hashAccess[:])

	hashRefresh := sha256.Sum256([]byte(refreshToken))
	refreshTokenHash := hex.EncodeToString(hashRefresh[:])

	session := domain.Session{
		UserID:           user.ID,
		TokenHash:        accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		ExpiredAt:        time.Now().AddDate(0, 0, expireDays),
	}
	if err := h.sessRepo.CreateSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SESSION_CREATION_FAILED", "Gagal menyimpan session").WithContext(c))
		return
	}

	// Save initial active role session
	activeSession := domain.ActiveRoleSession{
		UserID:    user.ID,
		RoleID:    defaultUserRole.RoleID,
		SessionID: session.ID,
		// Initially application is portal or none
		ApplicationID:  "", // Application registry linkable later
		StudyProgramID: defaultUserRole.StudyProgramID,
	}
	_ = h.sessRepo.SaveActiveRoleSession(&activeSession)

	c.JSON(http.StatusOK, sharederr.Success(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(expireHours * 3600),
		TokenType:    "Bearer",
	}).WithContext(c))
}

// Refresh issues a new access token given a valid refresh token
func (h *AuthHandler) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	session, err := h.sessRepo.GetByRefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_SESSION", "Refresh token tidak valid atau telah dicabut").WithContext(c))
		return
	}

	if session.ExpiredAt.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, sharederr.Error("EXPIRED_SESSION", "Sesi Anda telah kedaluwarsa").WithContext(c))
		return
	}

	// Fetch active role session to reissue correct role context
	activeSession, err := h.sessRepo.GetActiveRoleSession(session.ID)
	var activeRole string
	var scope string
	var permissionCodes []string

	if err == nil {
		activeRole = activeSession.Role.Code
		scope = activeSession.Role.ScopeType
		if activeSession.Role.ScopeType == "study_program" && activeSession.StudyProgramID != nil {
			scope = "study_program:" + *activeSession.StudyProgramID
		}

		permissions, _ := h.userRepo.GetRolePermissions(activeSession.RoleID)
		permissionCodes = make([]string, len(permissions))
		for i, p := range permissions {
			permissionCodes[i] = p.Code
		}
	} else {
		// Fallback to default first role
		userRoles, _ := h.userRepo.GetUserRoles(session.UserID)
		if len(userRoles) > 0 {
			activeRole = userRoles[0].Role.Code
			scope = userRoles[0].Role.ScopeType
			if userRoles[0].Role.ScopeType == "study_program" && userRoles[0].StudyProgramID != nil {
				scope = "study_program:" + *userRoles[0].StudyProgramID
			}
			permissions, _ := h.userRepo.GetRolePermissions(userRoles[0].RoleID)
			permissionCodes = make([]string, len(permissions))
			for i, p := range permissions {
				permissionCodes[i] = p.Code
			}
		}
	}

	expireHours := 24
	if envHours := os.Getenv("JWT_EXPIRE_HOURS"); envHours != "" {
		if val, err := strconv.Atoi(envHours); err == nil {
			expireHours = val
		}
	}

	accessToken, err := keys.GenerateAccessToken(session.UserID, activeRole, scope, permissionCodes, expireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("TOKEN_GENERATION_FAILED", "Gagal menerbitkan access token").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: body.RefreshToken,
		ExpiresIn:    int64(expireHours * 3600),
		TokenType:    "Bearer",
	}).WithContext(c))
}

// Me returns currently authenticated user details
func (h *AuthHandler) Me(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	user, err := h.userRepo.GetByID(claims.Subject)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("USER_NOT_FOUND", "User tidak ditemukan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"user_id":     user.ID,
		"person_id":   user.PersonID,
		"name":        user.Person.Name,
		"email":       user.Person.Email,
		"active_role": claims.ActiveRole,
		"permissions": claims.Permissions,
		"scope":       claims.Scope,
	}).WithContext(c))
}

// SwitchRole changes user context role and returns new access/refresh token
func (h *AuthHandler) SwitchRole(c *gin.Context) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Missing authentication context").WithContext(c))
		return
	}
	claims := claimsVal.(*sharedauth.Claims)

	var req SwitchRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	userRoles, err := h.userRepo.GetUserRoles(claims.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memproses roles").WithContext(c))
		return
	}

	var targetRole *domain.UserRole
	for _, ur := range userRoles {
		if ur.Role.Code == req.RoleCode {
			// If it's a study program role, make sure program ID matches or is provided
			if ur.Role.ScopeType == "study_program" {
				if req.ScopeValue != "" && ur.StudyProgramID != nil && *ur.StudyProgramID == req.ScopeValue {
					targetRole = &ur
					break
				} else if req.ScopeValue == "" && ur.StudyProgramID != nil {
					targetRole = &ur // default if only one
					break
				}
			} else {
				targetRole = &ur
				break
			}
		}
	}

	if targetRole == nil {
		c.JSON(http.StatusForbidden, sharederr.Error("ROLE_ACCESS_DENIED", "Anda tidak memiliki akses ke role tersebut").WithContext(c))
		return
	}

	// Fetch permissions
	permissions, err := h.userRepo.GetRolePermissions(targetRole.RoleID)
	if err != nil {
		permissions = []domain.Permission{}
	}
	permissionCodes := make([]string, len(permissions))
	for i, p := range permissions {
		permissionCodes[i] = p.Code
	}

	scope := targetRole.Role.ScopeType
	if targetRole.Role.ScopeType == "study_program" && targetRole.StudyProgramID != nil {
		scope = "study_program:" + *targetRole.StudyProgramID
	}

	expireHours := 24
	if envHours := os.Getenv("JWT_EXPIRE_HOURS"); envHours != "" {
		if val, err := strconv.Atoi(envHours); err == nil {
			expireHours = val
		}
	}

	accessToken, err := keys.GenerateAccessToken(claims.Subject, targetRole.Role.Code, scope, permissionCodes, expireHours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("TOKEN_GENERATION_FAILED", "Gagal menerbitkan access token").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: "", // In switch-role we only replace access token context
		ExpiresIn:    int64(expireHours * 3600),
		TokenType:    "Bearer",
	}).WithContext(c))
}

// ListApplications retrieves launcher applications based on user
func (h *AuthHandler) ListApplications(c *gin.Context) {
	apps, err := h.appRepo.GetAllEnabled()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar aplikasi").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(apps).WithContext(c))
}

// JWKS serves the JWK public keys
func (h *AuthHandler) JWKS(c *gin.Context) {
	if keys.PublicKey == nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("KEYS_NOT_INITIALIZED", "Public key belum siap").WithContext(c))
		return
	}

	nBytes := keys.PublicKey.N.Bytes()
	nStr := base64.RawURLEncoding.EncodeToString(nBytes)

	eBytes := big.NewInt(int64(keys.PublicKey.E)).Bytes()
	eStr := base64.RawURLEncoding.EncodeToString(eBytes)

	jwks := JWKSResponse{
		Keys: []JWKKey{
			{
				Kty: "RSA",
				Use: "sig",
				Kid: keys.KeyID,
				Alg: "RS256",
				N:   nStr,
				E:   eStr,
			},
		},
	}

	c.JSON(http.StatusOK, jwks)
}

// OpenIDConfiguration serves the OIDC configuration metadata
func (h *AuthHandler) OpenIDConfiguration(c *gin.Context) {
	baseURL := os.Getenv("JWKS_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8001"
	}

	c.JSON(http.StatusOK, gin.H{
		"issuer":                                baseURL,
		"authorization_endpoint":                baseURL + "/api/v1/oauth/authorize",
		"token_endpoint":                        baseURL + "/api/v1/oauth/token",
		"jwks_uri":                              baseURL + "/.well-known/jwks.json",
		"response_types_supported":              []string{"code", "token"},
		"subject_types_supported":                []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_post", "client_secret_basic"},
	})
}


