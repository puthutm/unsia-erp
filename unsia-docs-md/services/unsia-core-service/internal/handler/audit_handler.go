package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"gorm.io/gorm"
)

type AuditHandler struct {
	db *gorm.DB
}

func NewAuditHandler(db *gorm.DB) *AuditHandler {
	return &AuditHandler{db: db}
}

type AuditLog struct {
	ID                     string    `json:"id"`
	UserID                 *string   `json:"user_id"`
	ActorUserID            *string   `json:"actor_user_id"`
	TargetUserID            *string   `json:"target_user_id"`
	ActiveRoleID           *string   `json:"active_role_id"`
	ImpersonationSessionID  *string   `json:"impersonation_session_id"`
	ApplicationID          *string   `json:"application_id"`
	Module                string    `json:"module"`
	Action                string    `json:"action"`
	EntityName            *string   `json:"entity_name"`
	EntityID              *string   `json:"entity_id"`
	Reason                *string   `json:"reason"`
	OldValue              *string   `json:"old_value"`
	NewValue              *string   `json:"new_value"`
	RequestID             *string   `json:"request_id"`
	IPAddress             *string   `json:"ip_address"`
	UserAgent             *string   `json:"user_agent"`
	CreatedAt             time.Time `json:"created_at"`
}

// ListAuditLogs handles GET /api/v1/audit-logs
func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Parse query params
	module := c.Query("module")
	action := c.Query("action")
	entityName := c.Query("entity_name")
	entityID := c.Query("entity_id")

	// Build query
	query := h.db.Where("user_id = ? OR actor_user_id = ?", userID, userID)

	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if entityName != "" {
		query = query.Where("entity_name = ?", entityName)
	}
	if entityID != "" {
		query = query.Where("entity_id = ?", entityID)
	}

	var logs []AuditLog
	if err := query.Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("SERVER_ERROR", err.Error()).WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(logs).WithContext(c))
}

// GetAuditLog handles GET /api/v1/audit-logs/:id
func (h *AuditHandler) GetAuditLog(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var log AuditLog
	if err := h.db.Where("id = ? AND (user_id = ? OR actor_user_id = ?)", id, userID, userID).First(&log).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Audit log not found").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(log).WithContext(c))
}

// CreateAuditLog creates an audit log entry
func (h *AuditHandler) CreateAuditLog(log *AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.New().String()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	return h.db.Create(log).Error
}

// ListAuditLogsByEntity lists audit logs for a specific entity
func (h *AuditHandler) ListAuditLogsByEntity(entityName, entityID string) ([]AuditLog, error) {
	var logs []AuditLog
	err := h.db.Where("entity_name = ? AND entity_id = ?", entityName, entityID).
		Order("created_at DESC").
		Limit(50).
		Find(&logs).Error
	return logs, err
}
