package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ChatHandler - Handler for chat/obrolan features
type ChatHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

// NewChatHandler - Create new chat handler
func NewChatHandler(db *gorm.DB) *ChatHandler {
	return &ChatHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// CreateRoomRequest - Request body for creating chat room
type CreateRoomRequest struct {
	LmsClassID *string `json:"lms_class_id"`
	SessionID *string `json:"session_id"`
	RoomType  string  `json:"room_type" binding:"required"` // class, session, group, private
	Name      string  `json:"name" binding:"required"`
}

// SendMessageRequest - Request body for sending message
type SendMessageRequest struct {
	Message string  `json:"message" binding:"required"`
	FileURL *string `json:"file_url"`
	ReplyTo *string `json:"reply_to"`
}

// CreateChatRoom - Create new chat room
func (h *ChatHandler) CreateChatRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	room := domain.ChatRoom{
		LmsClassID: req.LmsClassID,
		SessionID: req.SessionID,
		RoomType: req.RoomType,
		Name:     req.Name,
		CreatedBy: userStr,
		IsActive: true,
	}

	if err := h.db.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat ruang obrolan").WithContext(c))
		return
	}

	// Add creator as participant
	participant := domain.ChatParticipant{
		RoomID:   room.ID,
		UserID:   userStr,
		Role:    "admin",
		JoinedAt: time.Now(),
	}
	h.db.Create(&participant)

	c.JSON(http.StatusCreated, sharederr.Success(room).WithContext(c))
}

// ListChatRooms - List chat rooms
func (h *ChatHandler) ListChatRooms(c *gin.Context) {
	classID := c.Query("class_id")
	sessionID := c.Query("session_id")

	query := h.db.Model(&domain.ChatRoom{}).Where("is_active = ?", true)

	if classID != "" {
		query = query.Where("lms_class_id = ?", classID)
	}
	if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	var rooms []domain.ChatRoom
	if err := query.Order("created_at DESC").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data ruang obrolan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(rooms).WithContext(c))
}

// GetChatRoom - Get single chat room with messages
func (h *ChatHandler) GetChatRoom(c *gin.Context) {
	roomID := c.Param("id")

	var room domain.ChatRoom
	if err := h.db.First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Ruang obrolan tidak ditemukan").WithContext(c))
		return
	}

	// Get messages
	var messages []domain.ChatMessage
	h.db.Where("room_id = ?", roomID).Order("created_at DESC").Limit(100).Find(&messages)

	// Get participants
	var participants []domain.ChatParticipant
	h.db.Where("room_id = ?", roomID).Find(&participants)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"room":        room,
		"messages":    messages,
		"participants": participants,
	}).WithContext(c))
}

// SendMessage - Send message to chat room
func (h *ChatHandler) SendMessage(c *gin.Context) {
	roomID := c.Param("id")

	// Check if room exists
	var room domain.ChatRoom
	if err := h.db.First(&room, "id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Ruang obrolan tidak ditemukan").WithContext(c))
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	message := domain.ChatMessage{
		RoomID:   roomID,
		SenderID: userStr,
		Message: req.Message,
		FileURL: req.FileURL,
		ReplyTo: req.ReplyTo,
	}

	if err := h.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengirim pesan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(message).WithContext(c))
}

// ListMessages - List messages in a room
func (h *ChatHandler) ListMessages(c *gin.Context) {
	roomID := c.Param("id")

	var messages []domain.ChatMessage
	if err := h.db.Where("room_id = ?", roomID).Order("created_at DESC").Limit(100).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil pesan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(messages).WithContext(c))
}

// JoinRoom - Join a chat room
func (h *ChatHandler) JoinRoom(c *gin.Context) {
	roomID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	// Check if already a participant
	var existing domain.ChatParticipant
	if err := h.db.Where("room_id = ? AND user_id = ?", roomID, userStr).First(&existing).Error; err == nil {
		c.JSON(http.StatusOK, sharederr.Success(existing).WithContext(c))
		return
	}

	participant := domain.ChatParticipant{
		RoomID:   roomID,
		UserID:   userStr,
		Role:    "member",
		JoinedAt: time.Now(),
	}

	if err := h.db.Create(&participant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal bergabung ke ruang obrolan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(participant).WithContext(c))
}

// LeaveRoom - Leave a chat room
func (h *ChatHandler) LeaveRoom(c *gin.Context) {
	roomID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	if err := h.db.Where("room_id = ? AND user_id = ?", roomID, userStr).Delete(&domain.ChatParticipant{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal keluar dari ruang obrolan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Berhasil keluar dari ruang obrolan").WithContext(c))
}

// MarkAsRead - Mark messages as read
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	roomID := c.Param("id")

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	// Update last read timestamp
	if err := h.db.Model(&domain.ChatParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, userStr).
		Update("last_read_at", time.Now()).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menandai pesan sudah dibaca").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Pesan ditandai sudah dibaca").WithContext(c))
}
