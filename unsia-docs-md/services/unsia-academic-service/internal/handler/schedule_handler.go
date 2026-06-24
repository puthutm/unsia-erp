package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-academic-service/internal/domain"
	"github.com/unsia-erp/unsia-academic-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

// ============ Schedule Handler ============

// ScheduleHandler handles schedule-related operations
type ScheduleHandler struct {
	repo *repository.AcademicRepository
	db   *gorm.DB
}

// NewScheduleHandler creates a new ScheduleHandler
func NewScheduleHandler(db *gorm.DB) *ScheduleHandler {
	return &ScheduleHandler{
		repo: repository.NewAcademicRepository(db),
		db:   db,
	}
}

// ============ Schedule Request Types ============

type ScheduleCreateRequest struct {
	ClassID       string  `json:"class_id" binding:"required"`
	DayOfWeek     int     `json:"day_of_week" binding:"required,min=1,max=7"`
	StartTime     string  `json:"start_time" binding:"required"`
	EndTime       string  `json:"end_time" binding:"required"`
	RoomID        *string `json:"room_id"`
	BuildingID   *string `json:"building_id"`
	ScheduleType string  `json:"schedule_type"`
	IsOnline     *bool   `json:"is_online"`
	MeetingLink  *string `json:"meeting_link"`
}

type ScheduleUpdateRequest struct {
	DayOfWeek     *int    `json:"day_of_week"`
	StartTime     *string `json:"start_time"`
	EndTime       *string `json:"end_time"`
	RoomID        *string `json:"room_id"`
	BuildingID   *string `json:"building_id"`
	ScheduleType *string `json:"schedule_type"`
	IsOnline     *bool   `json:"is_online"`
	MeetingLink  *string `json:"meeting_link"`
}

type BulkScheduleRequest struct {
	ClassID     string                  `json:"class_id" binding:"required"`
	Schedules   []ScheduleCreateRequest    `json:"schedules" binding:"required,gt=0"`
}

type ScheduleConflictRequest struct {
	DayOfWeek   int     `json:"day_of_week" binding:"required,min=1,max=7"`
	StartTime  string  `json:"start_time" binding:"required"`
	EndTime   string  `json:"end_time" binding:"required"`
	RoomID    *string `json:"room_id"`
	BuildingID *string `json:"building_id"`
}

// ============ Create Schedule ============

// CreateSchedule handles POST /api/v1/academic/schedules
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var req ScheduleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Validate time format
	if !isValidTimeFormat(req.StartTime) || !isValidTimeFormat(req.EndTime) {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError("Invalid time format. Use HH:MM").WithContext(c))
		return
	}

	// Validate time range
	if !isTimeRangeValid(req.StartTime, req.EndTime) {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError("End time must be after start time").WithContext(c))
		return
	}

	// Check if class exists
	class, err := h.repo.GetClassByID(req.ClassID)
	if err != nil || class == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kelas tidak ditemukan").WithContext(c))
		return
	}

	// Check for room conflict
	if req.RoomID != nil {
		conflict, _ := h.checkRoomConflict(req.ClassID, *req.RoomID, req.DayOfWeek, req.StartTime, req.EndTime, "")
		if conflict {
			c.JSON(http.StatusBadRequest, sharederr.Error("CONFLICT", "Ruangan sudah digunakan pada waktu tersebut").WithContext(c))
			return
		}
	}

	isOnline := false
	if req.IsOnline != nil {
		isOnline = *req.IsOnline
	}

	scheduleType := "regular"
	if req.ScheduleType != "" {
		scheduleType = req.ScheduleType
	}

	schedule := domain.ClassSchedule{
		ClassID:      req.ClassID,
		DayOfWeek:    req.DayOfWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		RoomID:       req.RoomID,
		BuildingID:   req.BuildingID,
		ScheduleType: scheduleType,
		IsOnline:     isOnline,
		MeetingLink: req.MeetingLink,
	}

	if err := h.repo.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jadwal").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.schedule.create",
		Module:       "academic",
		ResourceType: "schedule",
		ResourceID:   schedule.ID,
		NewValue:     schedule,
	})

	c.JSON(http.StatusCreated, sharederr.Success(schedule).WithContext(c))
}

// ============ Get Schedule ============

// GetSchedule handles GET /api/v1/academic/schedules/:id
func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, err := h.repo.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal tidak ditemukan").WithContext(c))
		return
	}

	// Get class details
	class, _ := h.repo.GetClassByID(schedule.ClassID)

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"schedule": schedule,
		"class":    class,
	}).WithContext(c))
}

// ============ List Schedules ============

// ListSchedules handles GET /api/v1/academic/schedules
func (h *ScheduleHandler) ListSchedules(c *gin.Context) {
	classID := c.Query("class_id")
	dayOfWeek := c.Query("day_of_week")
	roomID := c.Query("room_id")
	buildingID := c.Query("building_id")

	var schedules []domain.ClassSchedule
	query := h.db.Model(&domain.ClassSchedule{})

	if classID != "" {
		query = query.Where("class_id = ?", classID)
	}
	if dayOfWeek != "" {
		day, _ := strconv.Atoi(dayOfWeek)
		query = query.Where("day_of_week = ?", day)
	}
	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}
	if buildingID != "" {
		query = query.Where("building_id = ?", buildingID)
	}

	if err := query.Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar jadwal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(schedules).WithContext(c))
}

// ============ Update Schedule ============

// UpdateSchedule handles PUT /api/v1/academic/schedules/:id
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, err := h.repo.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal tidak ditemukan").WithContext(c))
		return
	}

	var req ScheduleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if req.StartTime != nil {
		if !isValidTimeFormat(*req.StartTime) {
			c.JSON(http.StatusBadRequest, sharederr.ValidationError("Invalid start time format").WithContext(c))
			return
		}
		schedule.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		if !isValidTimeFormat(*req.EndTime) {
			c.JSON(http.StatusBadRequest, sharederr.ValidationError("Invalid end time format").WithContext(c))
			return
		}
		schedule.EndTime = *req.EndTime
	}
	if req.DayOfWeek != nil {
		schedule.DayOfWeek = *req.DayOfWeek
	}
	if req.RoomID != nil {
		schedule.RoomID = req.RoomID
	}
	if req.BuildingID != nil {
		schedule.BuildingID = req.BuildingID
	}
	if req.ScheduleType != nil {
		schedule.ScheduleType = *req.ScheduleType
	}
	if req.IsOnline != nil {
		schedule.IsOnline = *req.IsOnline
	}
	if req.MeetingLink != nil {
		schedule.MeetingLink = req.MeetingLink
	}

	if err := h.repo.UpdateSchedule(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui jadwal").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.schedule.update",
		Module:       "academic",
		ResourceType: "schedule",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.Success(schedule).WithContext(c))
}

// DeleteSchedule handles DELETE /api/v1/academic/schedules/:id
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")

	schedule, err := h.repo.GetScheduleByID(id)
	if err != nil || schedule == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jadwal tidak ditemukan").WithContext(c))
		return
	}

	if err := h.repo.DeleteSchedule(id); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus jadwal").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.schedule.delete",
		Module:       "academic",
		ResourceType: "schedule",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Jadwal berhasil dihapus").WithContext(c))
}

// ============ Bulk Schedule ============

// CreateBulkSchedules handles POST /api/v1/academic/schedules/bulk
func (h *ScheduleHandler) CreateBulkSchedules(c *gin.Context) {
	var req BulkScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if class exists
	class, err := h.repo.GetClassByID(req.ClassID)
	if err != nil || class == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Kelas tidak ditemukan").WithContext(c))
		return
	}

	var created []domain.ClassSchedule
	var failed []string

	for i, schedReq := range req.Schedules {
		// Validate time format
		if !isValidTimeFormat(schedReq.StartTime) || !isValidTimeFormat(schedReq.EndTime) {
			failed = append(failed, fmt.Sprintf("Row %d: invalid time format", i+1))
			continue
		}

		// Validate time range
		if !isTimeRangeValid(schedReq.StartTime, schedReq.EndTime) {
			failed = append(failed, fmt.Sprintf("Row %d: end time before start time", i+1))
			continue
		}

		schedule := domain.ClassSchedule{
			ClassID:      req.ClassID,
			DayOfWeek:    schedReq.DayOfWeek,
			StartTime:   schedReq.StartTime,
			EndTime:     schedReq.EndTime,
			RoomID:      schedReq.RoomID,
			BuildingID:  schedReq.BuildingID,
			ScheduleType: schedReq.ScheduleType,
			IsOnline:     false,
			MeetingLink: schedReq.MeetingLink,
		}

		if schedReq.IsOnline != nil {
			schedule.IsOnline = *schedReq.IsOnline
		}

		if err := h.repo.CreateSchedule(&schedule); err != nil {
			failed = append(failed, fmt.Sprintf("Row %d: save error", i+1))
			continue
		}

		created = append(created, schedule)
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "academic.schedule.bulk_create",
		Module:       "academic",
		ResourceType: "schedule",
		ResourceID:   "bulk_" + req.ClassID,
	})

	c.JSON(http.StatusCreated, sharederr.Success(gin.H{
		"created": created,
		"failed":  failed,
		"total":  len(req.Schedules),
	}).WithContext(c))
}

// ============ Schedule by Class ============

// GetSchedulesByClass handles GET /api/v1/academic/classes/:class_id/schedules
func (h *ScheduleHandler) GetSchedulesByClass(c *gin.Context) {
	classID := c.Param("class_id")

	schedules, err := h.repo.GetSchedulesByClassID(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jadwal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(schedules).WithContext(c))
}

// ============ Schedule Conflict Check ============

// CheckConflict handles POST /api/v1/academic/schedules/check-conflict
func (h *ScheduleHandler) CheckConflict(c *gin.Context) {
	var req ScheduleConflictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check room conflict
	hasConflict := false
	if req.RoomID != nil {
		conflict, _ := h.checkRoomConflict("", *req.RoomID, req.DayOfWeek, req.StartTime, req.EndTime, "")
		hasConflict = conflict
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"has_conflict":   hasConflict,
		"day_of_week":   req.DayOfWeek,
		"start_time":   req.StartTime,
		"end_time":     req.EndTime,
	}).WithContext(c))
}

// ============ Weekly Schedule ============

// GetWeeklySchedule handles GET /api/v1/academic/schedules/weekly
func (h *ScheduleHandler) GetWeeklySchedule(c *gin.Context) {
	studyProgramID := c.Query("study_program_id")
	academicPeriodID := c.Query("academic_period_id")

	schedules, err := h.repo.GetWeeklySchedule(studyProgramID, academicPeriodID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jadwal mingguan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(schedules).WithContext(c))
}

// ============ My Schedule (for Lecturer) ============

// GetMySchedule handles GET /api/v1/academic/schedules/my
func (h *ScheduleHandler) GetMySchedule(c *gin.Context) {
	lecturerID, _ := c.Get("x-user-id")
	if lecturerID == nil {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Unauthorized").WithContext(c))
		return
	}

	lecturerIDStr := lecturerID.(string)
	dayOfWeek := c.Query("day_of_week")

	schedules, err := h.repo.GetScheduleByLecturer(lecturerIDStr, dayOfWeek)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jadwal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(schedules).WithContext(c))
}

// ============ Student Schedule ============

// GetStudentSchedule handles GET /api/v1/academic/students/:student_id/schedule
func (h *ScheduleHandler) GetStudentSchedule(c *gin.Context) {
	studentID := c.Param("student_id")
	dayOfWeek := c.Query("day_of_week")

	schedules, err := h.repo.GetStudentSchedule(studentID, dayOfWeek)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jadwal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(schedules).WithContext(c))
}

// ============ Helper Functions ============

func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

func isTimeRangeValid(startTime, endTime string) bool {
	start, _ := time.Parse("15:04", startTime)
	end, _ := time.Parse("15:04", endTime)
	return end.After(start)
}

func (h *ScheduleHandler) checkRoomConflict(classID, roomID string, dayOfWeek int, startTime, endTime, excludeClassID string) (bool, error) {
	var count int64
	query := h.db.Model(&domain.ClassSchedule{}).
		Where("room_id = ?", roomID).
		Where("day_of_week = ?", dayOfWeek).
		Where("((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?))", 
			endTime, startTime, startTime, endTime)

	if excludeClassID != "" {
		query = query.Where("class_id != ?", excludeClassID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
