package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-portal-service/internal/domain"
	"github.com/unsia-erp/unsia-portal-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type NotificationCreateRequest struct {
	RecipientUserID string `json:"recipient_user_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
	Message         string `json:"message"`
	SourceModule    string `json:"source_module"`
	LinkUrl         string `json:"link_url"`
}

type PreferenceUpdateRequest struct {
	PreferenceKey   string `json:"preference_key" binding:"required"`
	PreferenceValue string `json:"preference_value" binding:"required"`
}

type ShortcutSaveRequest struct {
	MenuCode  string `json:"menu_code" binding:"required"`
	MenuLabel string `json:"menu_label" binding:"required"`
	TargetUrl string `json:"target_url" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

type PortalHandler struct {
	repo *repository.PortalRepository
	db   *gorm.DB
}

func NewPortalHandler(db *gorm.DB) *PortalHandler {
	return &PortalHandler{
		repo: repository.NewPortalRepository(db),
		db:   db,
	}
}

func (h *PortalHandler) CreateNotification(c *gin.Context) {
	var req NotificationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	notif := domain.Notification{
		UserID:       req.RecipientUserID,
		Title:        req.Title,
		Message:      req.Message,
		ModuleSource: req.SourceModule,
		TargetUrl:    req.LinkUrl,
	}

	if err := h.repo.CreateNotification(&notif); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengirim notifikasi").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(notif).WithContext(c))
}

func (h *PortalHandler) ListNotifications(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	list, err := h.repo.ListNotificationsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil notifikasi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(list).WithContext(c))
}

func (h *PortalHandler) MarkRead(c *gin.Context) {
	id := c.Param("id")
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	if err := h.repo.MarkNotificationAsRead(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menandai notifikasi dibaca").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Notifikasi telah dibaca").WithContext(c))
}

func (h *PortalHandler) GetDashboard(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	role := "public"
	if claims != nil && claims.ActiveRole != "" {
		role = claims.ActiveRole
	}

	widgets := []map[string]interface{}{}

	switch role {
	case "super_admin":
		widgets = []map[string]interface{}{
			{"type": "stat", "title": "Total Users", "value": 1500},
			{"type": "stat", "title": "Active Sessions", "value": 45},
			{"type": "chart", "title": "System Load", "value": "healthy"},
		}
	case "mahasiswa":
		widgets = []map[string]interface{}{
			{"type": "profile", "title": "Status Akademik", "value": "Aktif"},
			{"type": "krs_status", "title": "Status KRS", "value": "Disetujui"},
			{"type": "announcement", "title": "Pengumuman", "value": "UAS akan dimulai tanggal 1 Juli 2026"},
		}
	case "admin_pmb":
		widgets = []map[string]interface{}{
			{"type": "stat", "title": "Pendaftar Baru Hari Ini", "value": 12},
			{"type": "stat", "title": "Verifikasi Dokumen Pending", "value": 8},
		}
	default:
		widgets = []map[string]interface{}{
			{"type": "announcement", "title": "Selamat Datang di Portal ERP UNSIA", "value": "Gunakan menu navigasi untuk mengakses SIAKAD."},
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"role":    role,
		"widgets": widgets,
	}).WithContext(c))
}

func (h *PortalHandler) UpdatePreference(c *gin.Context) {
	var req PreferenceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	up := domain.UserPreference{
		UserID:          userID,
		PreferenceKey:   req.PreferenceKey,
		PreferenceValue: req.PreferenceValue,
		UpdatedAt:       time.Now(),
	}

	if err := h.repo.SaveUserPreference(&up); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui preferensi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(up).WithContext(c))
}

func (h *PortalHandler) SaveShortcut(c *gin.Context) {
	var req ShortcutSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	ms := domain.MenuShortcut{
		UserID:    userID,
		MenuCode:  req.MenuCode,
		MenuLabel: req.MenuLabel,
		TargetUrl: req.TargetUrl,
		SortOrder: req.SortOrder,
	}

	if err := h.repo.SaveMenuShortcut(&ms); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan shortcut menu").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(ms).WithContext(c))
}

func (h *PortalHandler) DeleteShortcut(c *gin.Context) {
	menuCode := c.Param("menu_code")
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	userID := ""
	if claims != nil {
		userID = claims.Subject
	}

	if err := h.repo.DeleteMenuShortcut(userID, menuCode); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus shortcut menu").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Shortcut menu berhasil dihapus").WithContext(c))
}

func (h *PortalHandler) GetRoleMenus(c *gin.Context) {
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	role := "public"
	if claims != nil && claims.ActiveRole != "" {
		role = claims.ActiveRole
	}

	menus, err := h.repo.GetMenusByRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data menu navigasi").WithContext(c))
		return
	}

	var finalTree []domain.MenuResponse

	// Filter top-level menus
	for _, m := range menus {
		if m.ParentCode == nil || *m.ParentCode == "" {
			mr := domain.MenuResponse{
				MenuCode: m.Code,
				Label:    m.Label,
				Path:     m.Path,
				Icon:     m.Icon,
				Children: []domain.MenuResponse{},
			}

			// Find submenus
			for _, sub := range menus {
				if sub.ParentCode != nil && *sub.ParentCode == m.Code {
					mr.Children = append(mr.Children, domain.MenuResponse{
						MenuCode: sub.Code,
						Label:    sub.Label,
						Path:     sub.Path,
						Icon:     sub.Icon,
					})
				}
			}
			finalTree = append(finalTree, mr)
		}
	}

	c.JSON(http.StatusOK, sharederr.Success(finalTree).WithContext(c))
}

