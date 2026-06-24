package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-hris-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type WorkUnitCreateRequest struct {
	Code           string  `json:"code" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	ParentID       *string `json:"parent_id"`
	Type          string  `json:"type"` // faculty, department, unit
	HeadEmployeeID *string `json:"head_employee_id"`
}

type PositionCreateRequest struct {
	Code         string  `json:"code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	WorkUnitID   string  `json:"work_unit_id" binding:"required"`
	Level        int     `json:"level"`
	ParentID     *string `json:"parent_id"`
	IsHead       bool    `json:"is_head"`
}

type OrganizationHandler struct {
	repo *repository.HRISRepository
	db   *gorm.DB
}

func NewOrganizationHandler(db *gorm.DB) *OrganizationHandler {
	return &OrganizationHandler{
		repo: repository.NewHRISRepository(db),
		db:   db,
	}
}

// ===== WORK UNIT (STRUKTUR ORGANISASI) =====

// ListWorkUnits - GET /api/v1/work-units
func (h *OrganizationHandler) ListWorkUnits(c *gin.Context) {
	units, err := h.repo.ListWorkUnits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data unit kerja").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(units).WithContext(c))
}

// GetWorkUnit - GET /api/v1/work-units/:id
func (h *OrganizationHandler) GetWorkUnit(c *gin.Context) {
	id := c.Param("id")
	unit, err := h.repo.GetWorkUnitByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data unit kerja").WithContext(c))
		return
	}
	if unit == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Unit kerja tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(unit).WithContext(c))
}

// CreateWorkUnit - POST /api/v1/work-units
func (h *OrganizationHandler) CreateWorkUnit(c *gin.Context) {
	var req WorkUnitCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	unit := domain.WorkUnit{
		Code:            req.Code,
		Name:            req.Name,
		ParentID:        req.ParentID,
		Type:             req.Type,
		HeadEmployeeID:  req.HeadEmployeeID,
	}

	if err := h.repo.CreateWorkUnit(&unit); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan data unit kerja").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(unit).WithContext(c))
}

// UpdateWorkUnit - PUT /api/v1/work-units/:id
func (h *OrganizationHandler) UpdateWorkUnit(c *gin.Context) {
	id := c.Param("id")
	var req WorkUnitCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	unit, err := h.repo.GetWorkUnitByID(id)
	if err != nil || unit == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Unit kerja tidak ditemukan").WithContext(c))
		return
	}

	unit.Code = req.Code
	unit.Name = req.Name
	unit.ParentID = req.ParentID
	unit.Type = req.Type
	unit.HeadEmployeeID = req.HeadEmployeeID

	if err := h.repo.UpdateWorkUnit(unit); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal memperbarui unit kerja").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(unit).WithContext(c))
}

// DeleteWorkUnit - DELETE /api/v1/work-units/:id
func (h *OrganizationHandler) DeleteWorkUnit(c *gin.Context) {
	id := c.Param("id")
	err := h.repo.DeleteWorkUnit(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menghapus unit kerja").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(nil).WithContext(c))
}

// GetOrgStructure - GET /api/v1/work-units/structure
func (h *OrganizationHandler) GetOrgStructure(c *gin.Context) {
	structure, err := h.repo.GetOrgStructureTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil struktur organisasi").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(structure).WithContext(c))
}

// ===== POSITION (JABATAN) =====

// ListPositions - GET /api/v1/positions
func (h *OrganizationHandler) ListPositions(c *gin.Context) {
	workUnitID := c.Query("work_unit_id")
	positions, err := h.repo.ListPositions(workUnitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data jabatan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(positions).WithContext(c))
}

// GetPosition - GET /api/v1/positions/:id
func (h *OrganizationHandler) GetPosition(c *gin.Context) {
	id := c.Param("id")
	pos, err := h.repo.GetPositionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil data jabatan").WithContext(c))
		return
	}
	if pos == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Jabatan tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(pos).WithContext(c))
}

// CreatePosition - POST /api/v1/positions
func (h *OrganizationHandler) CreatePosition(c *gin.Context) {
	var req PositionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	pos := domain.Position{
		Code:       req.Code,
		Name:       req.Name,
		WorkUnitID: req.WorkUnitID,
		Level:      req.Level,
		ParentID:  req.ParentID,
		IsHead:     req.IsHead,
	}

	if err := h.repo.CreatePosition(&pos); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan data jabatan").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(pos).WithContext(c))
}

// PlotEmployeeToPosition - PUT /api/v1/positions/:id/plot
func (h *OrganizationHandler) PlotEmployeeToPosition(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		EmployeeID string `json:"employee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	if err := h.repo.PlotEmployeeToPosition(id, req.EmployeeID); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal ploting karyawan ke jabatan").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(nil).WithContext(c))
}
