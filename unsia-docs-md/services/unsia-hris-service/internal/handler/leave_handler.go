package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type LeaveTypeCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Code       string `json:"code" binding:"required"`
	QuotaDays  int    `json:"quota_days"`
	CanExceed  bool   `json:"can_exceed"`
	IsPaid     bool   `json:"is_paid"`
}

type LeaveRequestCreate struct {
	EmployeeID     string  `json:"employee_id" binding:"required"`
	LeaveTypeID   string  `json:"leave_type_id" binding:"required"`
	StartDate    string  `json:"start_date" binding:"required"` // RFC3339
	EndDate     string  `json:"end_date" binding:"required"`
	Reason      string  `json:"reason"`
}

type LeaveApprovalRequest struct {
	Status    string `json:"status" binding:"required"` // approved, rejected
	Notes    string `json:"notes"`
}

type LeaveHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewLeaveHandler(db *gorm.DB) *LeaveHandler {
	return &LeaveHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

// ===== LEAVE TYPES =====

// ListLeaveTypes - GET /api/v1/leave-types
func (h *LeaveHandler) ListLeaveTypes(c *gin.Context) {
	types, err := h.repo.ListLeaveTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data jenis cuti").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(types).WithContext(c))
}

// CreateLeaveType - POST /api/v1/leave-types
func (h *LeaveHandler) CreateLeaveType(c *gin.Context) {
	var req LeaveTypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	leaveType := domain.LeaveType{
		Name:       req.Name,
		Code:       req.Code,
		QuotaDays:  req.QuotaDays,
		CanExceed: req.CanExceed,
		IsPaid:    req.IsPaid,
	}

	if err := h.repo.CreateLeaveType(&leaveType); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan jenis cuti").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(leaveType).WithContext(c))
}

// ===== LEAVE REQUESTS (CUTI & IZIN) =====

// ListLeaveRequests - GET /api/v1/leave-requests
// Query: status, employee_id, leave_type_id
func (h *LeaveHandler) ListLeaveRequests(c *gin.Context) {
	status := c.Query("status")
	employeeID := c.Query("employee_id")
	leaveTypeID := c.Query("leave_type_id")

	requests, err := h.repo.ListLeaveRequests(status, employeeID, leaveTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pengajuan cuti").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(requests).WithContext(c))
}

// GetLeaveRequest - GET /api/v1/leave-requests/:id
func (h *LeaveHandler) GetLeaveRequest(c *gin.Context) {
	id := c.Param("id")
	req, err := h.repo.GetLeaveRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pengajuan cuti").WithContext(c))
		return
	}
	if req == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengajuan cuti tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(req).WithContext(c))
}

// CreateLeaveRequest - POST /api/v1/leave-requests (Employee submits)
func (h *LeaveHandler) CreateLeaveRequest(c *gin.Context) {
	var req LeaveRequestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	// Parse dates
	startDate, _ := time.Parse(time.RFC3339, req.StartDate)
	endDate, _ := time.Parse(time.RFC3339, req.EndDate)

	// Check quota
	quota, err := h.repo.GetLeaveQuota(req.EmployeeID, req.LeaveTypeID)
	if err == nil && !quota.CanExceed {
		days := int(endDate.Sub(startDate).Hours()/24) + 1
		if days > quota.RemainingDays {
			c.JSON(http.StatusBadRequest, sharederr.Error("QUOTA_EXCEED", "Kuota cuti tidak mencukupi").WithContext(c))
			return
		}
	}

	leaveReq := domain.LeaveRequest{
		EmployeeID:   req.EmployeeID,
		LeaveTypeID: req.LeaveTypeID,
		StartDate:   startDate,
		EndDate:    endDate,
		Reason:     req.Reason,
		Status:     "pending",
	}

	if err := h.repo.CreateLeaveRequest(&leaveReq); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan pengajuan cuti").WithContext(c))
		return
	}

	// TODO: Send notification to manager
	c.JSON(http.StatusCreated, sharederr.Success(leaveReq).WithContext(c))
}

// ApproveLeaveRequest - PUT /api/v1/leave-requests/:id/approve (Manager approves)
func (h *LeaveHandler) ApproveLeaveRequest(c *gin.Context) {
	id := c.Param("id")
	var req LeaveApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	leaveReq, err := h.repo.GetLeaveRequestByID(id)
	if err != nil || leaveReq == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Pengajuan cuti tidak ditemukan").WithContext(c))
		return
	}

	leaveReq.Status = req.Status
	leaveReq.Notes = req.Notes

	if err := h.repo.UpdateLeaveRequest(leaveReq); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui status cuti").WithContext(c))
		return
	}

	// TODO: Send notification to employee
	c.JSON(http.StatusOK, sharederr.Success(leaveReq).WithContext(c))
}

// GetMyLeaveQuota - GET /api/v1/leave-requests/quota (ESS view own quota)
func (h *LeaveHandler) GetMyLeaveQuota(c *gin.Context) {
	employeeID := c.GetHeader("X-Employee-ID")
	if employeeID == "" {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Employee ID diperlukan").WithContext(c))
		return
	}

	quotas, err := h.repo.GetEmployeeLeaveQuotas(employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil quota cuti").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(quotas).WithContext(c))
}

// GetMyLeaveRequests - GET /api/v1/leave-requests/my (ESS view own requests)
func (h *LeaveHandler) GetMyLeaveRequests(c *gin.Context) {
	employeeID := c.GetHeader("X-Employee-ID")
	if employeeID == "" {
		c.JSON(http.StatusUnauthorized, sharederr.Error("UNAUTHORIZED", "Employee ID diperlukan").WithContext(c))
		return
	}

	requests, err := h.repo.ListLeaveRequests("", employeeID, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data pengajuan cuti").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(requests).WithContext(c))
}
