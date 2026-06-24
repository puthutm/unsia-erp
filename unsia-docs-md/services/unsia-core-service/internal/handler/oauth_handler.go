package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type OAuthRegisterRequest struct {
	OwnerName             string   `json:"owner_name" binding:"required"`
	OwnerEmail            string   `json:"owner_email" binding:"required"`
	OwnerOrganization     string   `json:"owner_organization" binding:"required"`
	RequestedScopes       []string `json:"requested_scopes" binding:"required"`
	RequestedGrantTypes   []string `json:"requested_grant_types" binding:"required"`
	RequestedRedirectURIs []string `json:"requested_redirect_uris" binding:"required"`
}

type OAuthHandler struct {
	oauthRepo *repository.OAuthRepository
	userRepo  *repository.UserRepository
	db        *gorm.DB
}

func NewOAuthHandler(db *gorm.DB) *OAuthHandler {
	return &OAuthHandler{
		oauthRepo: repository.NewOAuthRepository(db),
		userRepo:  repository.NewUserRepository(db),
		db:        db,
	}
}

// Authorize handles OAuth 2.0 Authorization Code flow redirect with PKCE
func (h *OAuthHandler) Authorize(c *gin.Context) {
	responseType := c.Query("response_type")
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	scope := c.Query("scope")
	state := c.Query("state")
	codeChallenge := c.Query("code_challenge")
	codeChallengeMethod := c.Query("code_challenge_method")

	if responseType != "code" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_REQUEST", "response_type must be 'code'").WithContext(c))
		return
	}
	if clientID == "" || redirectURI == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_REQUEST", "Missing client_id or redirect_uri").WithContext(c))
		return
	}

	client, err := h.oauthRepo.GetClientByID(clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("CLIENT_NOT_FOUND", "Client OAuth tidak terdaftar").WithContext(c))
		return
	}

	if client.Status != "ACTIVE" {
		c.JSON(http.StatusForbidden, sharederr.Error("CLIENT_SUSPENDED", "Client OAuth dinonaktifkan").WithContext(c))
		return
	}

	// Validate Redirect URI
	valid, err := h.oauthRepo.ValidateRedirectURI(client.ID, redirectURI)
	if err != nil || !valid {
		c.JSON(http.StatusForbidden, sharederr.Error("INVALID_REDIRECT_URI", "Redirect URI tidak cocok dengan yang terdaftar").WithContext(c))
		return
	}

	// PKCE is mandatory
	if codeChallenge == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_REQUEST", "PKCE code_challenge is required").WithContext(c))
		return
	}
	if codeChallengeMethod == "" {
		codeChallengeMethod = "S256"
	}
	if codeChallengeMethod != "S256" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_REQUEST", "code_challenge_method must be S256").WithContext(c))
		return
	}

	// Authentication Check: User must be logged in. In a full web app, they have a session/cookie.
	// For API simulation, we check for a valid Bearer token in the request.
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		// Simulation: For SSO in browser, we would redirect to a login page.
		// Since this is a REST API, we reject with 401.
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "User session required to authorize client").WithContext(c))
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Invalid authorization format").WithContext(c))
		return
	}

	claims, err := sharedauth.ValidateJWT(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Session expired").WithContext(c))
		return
	}

	// Generate Auth Code
	randBytes := make([]byte, 32)
	_, _ = rand.Read(randBytes)
	code := base64.RawURLEncoding.EncodeToString(randBytes)
	codeHash := fmt.Sprintf("%x", sha256.Sum256([]byte(code)))

	authCode := domain.OAuthAuthorizationCode{
		CodeHash:            codeHash,
		ClientID:            client.ID,
		UserID:              claims.Subject,
		RedirectURI:         redirectURI,
		Scope:               scope,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		State:               state,
		ExpiresAt:           time.Now().Add(5 * time.Minute),
	}

	if err := h.oauthRepo.SaveAuthCode(&authCode); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal menerbitkan auth code").WithContext(c))
		return
	}

	// Standard Redirect
	redirectURL := fmt.Sprintf("%s?code=%s", redirectURI, code)
	if state != "" {
		redirectURL = redirectURL + "&state=" + state
	}

	c.Redirect(http.StatusFound, redirectURL)
}

// Token handles Token Issuance for both authorization_code (PKCE) and client_credentials grants
func (h *OAuthHandler) Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	// If variables are in request JSON instead of form data
	if grantType == "" {
		var reqBody struct {
			GrantType    string `json:"grant_type"`
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			Code         string `json:"code"`
			RedirectURI  string `json:"redirect_uri"`
			CodeVerifier string `json:"code_verifier"`
		}
		if err := c.ShouldBindJSON(&reqBody); err == nil {
			grantType = reqBody.GrantType
			clientID = reqBody.ClientID
			clientSecret = reqBody.ClientSecret
			c.Set("json_code", reqBody.Code)
			c.Set("json_redirect_uri", reqBody.RedirectURI)
			c.Set("json_code_verifier", reqBody.CodeVerifier)
		}
	}

	switch grantType {
	case "authorization_code":
		code := c.PostForm("code")
		redirectURI := c.PostForm("redirect_uri")
		codeVerifier := c.PostForm("code_verifier")

		if code == "" {
			if val, ok := c.Get("json_code"); ok {
				code = val.(string)
			}
		}
		if redirectURI == "" {
			if val, ok := c.Get("json_redirect_uri"); ok {
				redirectURI = val.(string)
			}
		}
		if codeVerifier == "" {
			if val, ok := c.Get("json_code_verifier"); ok {
				codeVerifier = val.(string)
			}
		}

		if code == "" || redirectURI == "" || codeVerifier == "" {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_REQUEST", "Missing code, redirect_uri or code_verifier").WithContext(c))
			return
		}

		client, err := h.oauthRepo.GetClientByID(clientID)
		if err != nil {
			c.JSON(http.StatusNotFound, sharederr.Error("CLIENT_NOT_FOUND", "Client tidak valid").WithContext(c))
			return
		}

		// Verify client secret if confidential client
		if client.ClientType == "confidential" {
			if client.ClientSecretHash == nil || bcrypt.CompareHashAndPassword([]byte(*client.ClientSecretHash), []byte(clientSecret)) != nil {
				c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_CLIENT", "Client secret salah").WithContext(c))
				return
			}
		}

		// Retrieve and check authorization code
		codeHash := fmt.Sprintf("%x", sha256.Sum256([]byte(code)))
		authCode, err := h.oauthRepo.GetAuthCode(codeHash)
		if err != nil {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_GRANT", "Authorization code tidak valid").WithContext(c))
			return
		}

		if authCode.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_GRANT", "Authorization code kedaluwarsa").WithContext(c))
			return
		}

		if authCode.UsedAt != nil {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_GRANT", "Authorization code sudah digunakan").WithContext(c))
			return
		}

		if authCode.RedirectURI != redirectURI {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_GRANT", "Redirect URI tidak cocok").WithContext(c))
			return
		}

		// Verify PKCE Verifier
		if authCode.CodeChallengeMethod == "S256" {
			hash := sha256.Sum256([]byte(codeVerifier))
			computedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
			if computedChallenge != authCode.CodeChallenge {
				c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_GRANT", "PKCE code_verifier tidak valid").WithContext(c))
				return
			}
		} else {
			if codeVerifier != authCode.CodeChallenge {
				c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_GRANT", "PKCE code_verifier tidak valid").WithContext(c))
				return
			}
		}

		// Mark auth code used
		_ = h.oauthRepo.MarkAuthCodeUsed(authCode.ID)

		// Fetch user details for scope & active role
		userRoles, _ := h.userRepo.GetUserRoles(authCode.UserID)
		activeRole := "guest"
		scope := "self"
		var permissions []string

		if len(userRoles) > 0 {
			activeRole = userRoles[0].Role.Code
			scope = userRoles[0].Role.ScopeType
			if userRoles[0].Role.ScopeType == "study_program" && userRoles[0].StudyProgramID != nil {
				scope = "study_program:" + *userRoles[0].StudyProgramID
			}
			rolePermissions, _ := h.userRepo.GetRolePermissions(userRoles[0].RoleID)
			permissions = make([]string, len(rolePermissions))
			for i, p := range rolePermissions {
				permissions[i] = p.Code
			}
		}

		// Generate access token
		accessToken, err := keys.GenerateAccessToken(authCode.UserID, activeRole, scope, permissions, 24)
		if err != nil {
			c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal merilis token").WithContext(c))
			return
		}

		c.JSON(http.StatusOK, TokenResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   86400,
		})

	case "client_credentials":
		if clientID == "" || clientSecret == "" {
			c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_REQUEST", "Missing client_id or client_secret").WithContext(c))
			return
		}

		client, err := h.oauthRepo.GetClientByID(clientID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_CLIENT", "Client OAuth tidak valid").WithContext(c))
			return
		}

		if client.ClientSecretHash == nil || bcrypt.CompareHashAndPassword([]byte(*client.ClientSecretHash), []byte(clientSecret)) != nil {
			c.JSON(http.StatusUnauthorized, sharederr.Error("INVALID_CLIENT", "Client secret salah").WithContext(c))
			return
		}

		// Client credentials produces machine-to-machine service token
		var scopes []string
		_ = json.Unmarshal([]byte(client.AllowedScopes), &scopes)

		accessToken, err := keys.GenerateAccessToken(client.ClientID, "service", "global", scopes, 24)
		if err != nil {
			c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal merilis token").WithContext(c))
			return
		}

		c.JSON(http.StatusOK, TokenResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   86400,
		})
	default:
		c.JSON(http.StatusBadRequest, sharederr.Error("UNSUPPORTED_GRANT_TYPE", "grant_type must be 'authorization_code' or 'client_credentials'").WithContext(c))
	}
}

// Register handles new OAuth Client registration requests (PENDING status)
func (h *OAuthHandler) Register(c *gin.Context) {
	var req OAuthRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	scopesJSON, _ := json.Marshal(req.RequestedScopes)
	grantsJSON, _ := json.Marshal(req.RequestedGrantTypes)
	urisJSON, _ := json.Marshal(req.RequestedRedirectURIs)

	regRequest := domain.ClientRegistrationRequest{
		OwnerName:             req.OwnerName,
		OwnerEmail:            req.OwnerEmail,
		OwnerOrganization:     req.OwnerOrganization,
		RequestedScopes:       string(scopesJSON),
		RequestedGrantTypes:   string(grantsJSON),
		RequestedRedirectURIs: string(urisJSON),
		Status:                "PENDING",
	}

	if err := h.oauthRepo.SaveRegistrationRequest(&regRequest); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal menyimpan request registrasi").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"request_id": regRequest.ID,
		"status":     regRequest.Status,
		"message":    "Permintaan registrasi klien OAuth berhasil disimpan dan menunggu persetujuan admin",
	}).WithContext(c))
}

// AdminListRequests returns all client registration requests
func (h *OAuthHandler) AdminListRequests(c *gin.Context) {
	reqs, err := h.oauthRepo.GetAllRegistrationRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(reqs).WithContext(c))
}

// AdminApprove approves registration, creates the OAuthClient record, and returns ClientID and ClientSecret
func (h *OAuthHandler) AdminApprove(c *gin.Context) {
	id := c.Param("id")

	req, err := h.oauthRepo.GetRegistrationRequest(id)
	if err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("REQUEST_NOT_FOUND", "Request tidak ditemukan").WithContext(c))
		return
	}

	if req.Status != "PENDING" {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_STATUS", "Hanya request status PENDING yang dapat disetujui").WithContext(c))
		return
	}

	// Generate Client ID
	randIDBytes := make([]byte, 16)
	_, _ = rand.Read(randIDBytes)
	clientID := "client_" + hex.EncodeToString(randIDBytes)

	// Generate Client Secret
	randSecretBytes := make([]byte, 24)
	_, _ = rand.Read(randSecretBytes)
	clientSecret := base64.RawURLEncoding.EncodeToString(randSecretBytes)

	secretHashBytes, _ := bcrypt.GenerateFromPassword([]byte(clientSecret), 10)
	secretHash := string(secretHashBytes)

	// Create application first in registry (registry model requirement)
	app := domain.Application{
		ApplicationCode: strings.ToUpper(req.OwnerOrganization) + "_" + strings.ReplaceAll(strings.ToUpper(req.OwnerName), " ", "_"),
		Name:            req.OwnerName,
		URL:             "http://localhost:3000", // placeholder
		Enabled:         true,
	}
	if err := h.db.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal mendaftarkan aplikasi ke registry").WithContext(c))
		return
	}

	client := domain.OAuthClient{
		ApplicationID:      app.ID,
		ClientID:           clientID,
		ClientSecretHash:   &secretHash,
		ClientName:         req.OwnerName,
		ClientType:         "confidential",
		GrantTypes:         req.RequestedGrantTypes,
		AllowedScopes:      req.RequestedScopes,
		Status:             "ACTIVE",
		OwnerName:          req.OwnerName,
		OwnerEmail:         req.OwnerEmail,
		OwnerOrganization:  req.OwnerOrganization,
		IsActive:           true,
	}

	var uris []string
	_ = json.Unmarshal([]byte(req.RequestedRedirectURIs), &uris)

	if err := h.oauthRepo.CreateClient(&client, uris); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", "Gagal menyimpan client").WithContext(c))
		return
	}

	// Update registration status
	req.Status = "APPROVED"
	now := time.Now()
	req.ReviewedAt = &now
	req.OAuthClientID = &client.ID
	h.db.Save(req)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"client_id":     clientID,
		"client_secret": clientSecret, // Returned ONLY ONCE!
		"status":        "ACTIVE",
		"message":       "Client OAuth berhasil dibuat. Harap catat client_secret ini karena tidak akan ditampilkan lagi.",
	}).WithContext(c))
}

// AdminSuspend suspends client
func (h *OAuthHandler) AdminSuspend(c *gin.Context) {
	id := c.Param("id")
	if err := h.oauthRepo.UpdateClientStatus(id, "SUSPENDED"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal suspend client").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success("Client status updated to SUSPENDED").WithContext(c))
}

// AdminRevoke revokes client
func (h *OAuthHandler) AdminRevoke(c *gin.Context) {
	id := c.Param("id")
	if err := h.oauthRepo.UpdateClientStatus(id, "REVOKED"); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal revoke client").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success("Client status updated to REVOKED").WithContext(c))
}
