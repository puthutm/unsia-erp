package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-lms-service/internal/domain"
	"github.com/unsia-erp/unsia-lms-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ForumHandler - Handler for forum/diskusi features
type ForumHandler struct {
	repo *repository.LMSRepository
	db   *gorm.DB
}

// NewForumHandler - Create new forum handler
func NewForumHandler(db *gorm.DB) *ForumHandler {
	return &ForumHandler{
		repo: repository.NewLMSRepository(db),
		db:   db,
	}
}

// CreateForumPostRequest - Request body for creating forum post
type CreateForumPostRequest struct {
	SessionID string  `json:"session_id" binding:"required"`
	Title    string  `json:"title"`
	Content  string  `json:"content" binding:"required"`
}

// CreateForumReplyRequest - Request body for creating reply
type CreateForumReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateForumPost - Create new forum post
func (h *ForumHandler) CreateForumPost(c *gin.Context) {
	var req CreateForumPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if session exists
	var session domain.Session
	if err := h.db.First(&session, "id = ?", req.SessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Sesi perkuliahan tidak ditemukan").WithContext(c))
		return
	}

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	post := domain.ForumPost{
		SessionID:  req.SessionID,
		AuthorID:  userStr,
		Title:    req.Title,
		Content:  req.Content,
		IsPinned: false,
		IsLocked: false,
	}

	if err := h.db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat posting forum").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(post).WithContext(c))
}

// GetForumPosts - Get forum posts for a session
func (h *ForumHandler) GetForumPosts(c *gin.Context) {
	sessionID := c.Query("session_id")

	query := h.db.Model(&domain.ForumPost{})

	if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	var posts []domain.ForumPost
	if err := query.Order("is_pinned DESC, created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data forum").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(posts).WithContext(c))
}

// GetForumPost - Get single forum post with replies
func (h *ForumHandler) GetForumPost(c *gin.Context) {
	postID := c.Param("id")

	var post domain.ForumPost
	if err := h.db.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Posting tidak ditemukan").WithContext(c))
		return
	}

	// Increment view count
	h.db.Model(&post).Update("view_count", post.ViewCount+1)

	// Get replies
	var replies []domain.ForumReply
	h.db.Where("post_id = ?", postID).Order("created_at ASC").Find(&replies)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"post":    post,
		"replies": replies,
	}).WithContext(c))
}

// UpdateForumPost - Update forum post
func (h *ForumHandler) UpdateForumPost(c *gin.Context) {
	postID := c.Param("id")

	var post domain.ForumPost
	if err := h.db.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Posting tidak ditemukan").WithContext(c))
		return
	}

	var req CreateForumPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	updates := map[string]interface{}{
		"title":   req.Title,
		"content": req.Content,
	}

	if err := h.db.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengubah posting").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(post).WithContext(c))
}

// DeleteForumPost - Delete forum post
func (h *ForumHandler) DeleteForumPost(c *gin.Context) {
	postID := c.Param("id")

	var post domain.ForumPost
	if err := h.db.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Posting tidak ditemukan").WithContext(c))
		return
	}

	// Delete all replies first
	h.db.Where("post_id = ?", postID).Delete(&domain.ForumReply{})

	// Delete post
	if err := h.db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus posting").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Posting berhasil dihapus").WithContext(c))
}

// PinForumPost - Pin/unpin forum post
func (h *ForumHandler) PinForumPost(c *gin.Context) {
	postID := c.Param("id")

	var post domain.ForumPost
	if err := h.db.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Posting tidak ditemukan").WithContext(c))
		return
	}

	newPinnedStatus := !post.IsPinned
	if err := h.db.Model(&post).Update("is_pinned", newPinnedStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengubah status pin").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"is_pinned": newPinnedStatus}).WithContext(c))
}

// LockForumPost - Lock/unlock forum post
func (h *ForumHandler) LockForumPost(c *gin.Context) {
	postID := c.Param("id")

	var post domain.ForumPost
	if err := h.db.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Posting tidak ditemukan").WithContext(c))
		return
	}

	newLockedStatus := !post.IsLocked
	if err := h.db.Model(&post).Update("is_locked", newLockedStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengubah status lock").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"is_locked": newLockedStatus}).WithContext(c))
}

// CreateReply - Create reply to forum post
func (h *ForumHandler) CreateReply(c *gin.Context) {
	postID := c.Param("id")

	// Check if post exists
	var post domain.ForumPost
	if err := h.db.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Posting tidak ditemukan").WithContext(c))
		return
	}

	if post.IsLocked {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Posting ini sudah dikunci").WithContext(c))
		return
	}

	var req CreateForumReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	reply := domain.ForumReply{
		PostID:   postID,
		AuthorID: userStr,
		Content:  req.Content,
	}

	if err := h.db.Create(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal Membuat balasan").WithContext(c))
		return
	}

	// Update reply count
	h.db.Model(&post).Update("reply_count", post.ReplyCount+1)

	c.JSON(http.StatusCreated, sharederr.Success(reply).WithContext(c))
}

// MarkAsAnswer - Mark reply as answer
func (h *ForumHandler) MarkAsAnswer(c *gin.Context) {
	replyID := c.Param("reply_id")

	var reply domain.ForumReply
	if err := h.db.First(&reply, "id = ?", replyID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Balasan tidak ditemukan").WithContext(c))
		return
	}

	// Unmark previous answer if any
	h.db.Model(&domain.ForumReply{}).Where("post_id = ? AND is_answer = ?", reply.PostID, true).Update("is_answer", false)

	// Mark this as answer
	if err := h.db.Model(&reply).Update("is_answer", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menandai sebagai jawaban").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"is_answer": true}).WithContext(c))
}

// UpvoteReply - Upvote a reply
func (h *ForumHandler) UpvoteReply(c *gin.Context) {
	replyID := c.Param("reply_id")

	var reply domain.ForumReply
	if err := h.db.First(&reply, "id = ?", replyID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Balasan tidak ditemukan").WithContext(c))
		return
	}

	newUpvotes := reply.Upvotes + 1
	if err := h.db.Model(&reply).Update("upvotes", newUpvotes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memberikan upvote").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{"upvotes": newUpvotes}).WithContext(c))
}

// DeleteReply - Delete a reply
func (h *ForumHandler) DeleteReply(c *gin.Context) {
	replyID := c.Param("reply_id")

	var reply domain.ForumReply
	if err := h.db.First(&reply, "id = ?", replyID).Error; err != nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Balasan tidak ditemukan").WithContext(c))
		return
	}

	// Get user ID
	userID, _ := c.Get("x-user-id")
	userStr, _ := userID.(string)

	// Check if user is author
	if reply.AuthorID != userStr {
		c.JSON(http.StatusForbidden, sharederr.Error("FORBIDDEN", "Tidak memiliki akses untuk menghapus balasan ini").WithContext(c))
		return
	}

	// Get post to update count
	var post domain.ForumPost
	h.db.First(&post, "id = ?", reply.PostID)

	// Delete reply
	if err := h.db.Delete(&reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus balasan").WithContext(c))
		return
	}

	// Update reply count
	if post.ID != "" {
		h.db.Model(&post).Update("reply_count", post.ReplyCount-1)
	}

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Balasan berhasil dihapus").WithContext(c))
}
