package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type AttendanceClockIn struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	DeviceID  *string `json:"device_id"`
}

type AttendanceClockOut struct {
	EmployeeID string `json:"employee_id" binding:"required"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	DeviceID  *string `json:"device_id"`
}

type AttendanceFilter struct {
	EmployeeID *string `form:"employee_id"`
	StartDate  *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
	Status   *string `form:"status"` // present, absent, late, alpha
}

type AttendanceHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
	return &AttendanceHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

// ClockIn - POST /api/v1/attendance/clock-in
func (h *AttendanceHandler) ClockIn(c *gin.Context) {
	var req AttendanceClockIn
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Check if already clocked in today
	today := time.Now()
	existing, _ := h.repo.GetTodayAttendance(req.EmployeeID, today)
	if existing != nil && existing.ClockIn != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("ALREADY_CLOCKED_IN", "Anda sudah presensi masuk hari ini").WithContext(c))
		return
	}

	// Get work schedule to determine if late
	schedule, _ := h.repo.GetEmployeeWorkSchedule(req.EmployeeID, today.Weekday())
	isLate := false
	if schedule != nil {
		clockInTime, _ := time.Parse("15:04", schedule.ClockInTime)
		if today.After(clockInTime) {
			isLate = true
		}
	}

	att := domain.Attendance{
		EmployeeID: req.EmployeeID,
		Date:      today,
		ClockIn:   &today,
		IsLate:    isLate,
		Status:    "present",
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		DeviceID:  req.DeviceID,
	}

	if err := h.repo.CreateAttendance(&att); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan presensi").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(att).WithContext(c))
}

// ClockOut - POST /api/v1/attendance/clock-out
func (h *AttendanceHandler) ClockOut(c *gin.Context) {
	var req AttendanceClockOut
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	today := time.Now()
	existing, _ := h.repo.GetTodayAttendance(req.EmployeeID, today)
	if existing == nil || existing.ClockIn == nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("NOT_CLOCKED_IN", "Anda belum presensi masuk").WithContext(c))
		return
	}
	if existing.ClockOut != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("ALREADY_CLOCKED_OUT", "Anda sudah presensi pulang").WithContext(c))
		return
	}

	// Calculate work hours
	var workHours float64
	if existing.ClockIn != nil {
		workHours = today.Sub(*existing.ClockIn).Hours()
	}

	// Check if early leave
	isEarly := false
	schedule, _ := h.repo.GetEmployeeWorkSchedule(req.EmployeeID, today.Weekday())
	if schedule != nil {
		clockOutTime, _ := time.Parse("15:04", schedule.ClockOutTime)
		if today.Before(clockOutTime) {
			isEarly = true
		}
	}

	existing.ClockOut = &today
	existing.WorkHours = workHours
	existing.IsEarlyLeave = isEarly
	existing.Latitude = req.Latitude
	existing.Longitude = req.Longitude
	existing.DeviceID = req.DeviceID

	if err := h.repo.UpdateAttendance(existing); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan presensi").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(existing).WithContext(c))
}

// GetMyAttendance - GET /api/v1/attendance/my
func (h *AttendanceHandler) GetMyAttendance(c *gin.Context) {
	employeeID := c.GetHeader("X-Employee-ID")
	if employeeID == "" {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Employee ID diperlukan").WithContext(c))
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	attendances, err := h.repo.GetEmployeeAttendance(employeeID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data presensi").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(attendances).WithContext(c))
}

// GetTodayAttendance - GET /api/v1/attendance/today
func (h *AttendanceHandler) GetTodayAttendance(c *gin.Context) {
	employeeID := c.Query("employee_id")
	if employeeID == "" {
		employeeID = c.GetHeader("X-Employee-ID")
	}
	if employeeID == "" {
		c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Employee ID diperlukan").WithContext(c))
		return
	}

	today := time.Now()
	att, _ := h.repo.GetTodayAttendance(employeeID, today)
	c.JSON(http.StatusOK, sharederr.Success(att).WithContext(c))
}

// ListAttendance - GET /api/v1/attendance (Admin/Manager)
func (h *AttendanceHandler) ListAttendance(c *gin.Context) {
	var filter AttendanceFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	attendances, err := h.repo.ListAttendance(filter.EmployeeID, filter.StartDate, filter.EndDate, filter.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data presensi").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(attendances).WithContext(c))
}

// GetAttendanceReport - GET /api/v1/attendance/report
func (h *AttendanceHandler) GetAttendanceReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	workUnitID := c.Query("work_unit_id")

	report, err := h.repo.GetAttendanceReport(startDate, endDate, workUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal membuat laporan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(report).WithContext(c))
}

// ===== WORK SCHEDULE =====

type WorkScheduleCreateRequest struct {
	EmployeeID    string   `json:"employee_id" binding:"required"`
	DayOfWeek    int      `json:"day_of_week" binding:"required"` // 0=Sun, 1=Mon...
	ClockInTime  string   `json:"clock_in_time" binding:"required"`
	ClockOutTime string   `json:"clock_out_time" binding:"required"`
	IsOffDay    bool     `json:"is_off_day"`
}

// ListWorkSchedules - GET /api/v1/work-schedules
func (h *AttendanceHandler) ListWorkSchedules(c *gin.Context) {
	employeeID := c.Query("employee_id")
	schedules, err := h.repo.ListWorkSchedules(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil jadwal kerja").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(schedules).WithContext(c))
}

// CreateWorkSchedule - POST /api/v1/work-schedules
func (h *AttendanceHandler) CreateWorkSchedule(c *gin.Context) {
	var req WorkScheduleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	schedule := domain.WorkSchedule{
		EmployeeID:    req.EmployeeID,
		DayOfWeek:     req.DayOfWeek,
		ClockInTime:   req.ClockInTime,
		ClockOutTime:  req.ClockOutTime,
		IsOffDay:     req.IsOffDay,
	}

	if err := h.repo.CreateWorkSchedule(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jadwal kerja").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(schedule).WithContext(c))
}
