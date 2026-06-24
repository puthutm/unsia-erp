package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)

// PurchaseOrderHandler handles purchase order-related endpoints
type PurchaseOrderHandler struct {
	*FinanceHandler
}

// NewPurchaseOrderHandler creates a new PurchaseOrderHandler
func NewPurchaseOrderHandler(fh *FinanceHandler) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{FinanceHandler: fh}
}

// GetPurchaseOrders handles GET /api/v1/finance/purchase-orders
func (h *PurchaseOrderHandler) GetPurchaseOrders(c *gin.Context) {
	filter := repository.POFilter{
		Status:   c.Query("status"),
		VendorID: c.Query("vendor_id"),
		Search:  c.Query("search"),
	}

	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	filter.Page = page
	filter.Limit = limit

	result, err := h.repo.GetPOs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar purchase order").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// CreatePurchaseOrder handles POST /api/v1/finance/purchase-orders
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c *gin.Context) {
	var req POCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	po := domain.PurchaseOrder{
		ID:            fmt.Sprintf("PO-%s", time.Now().Format("20060102150405")),
		PONumber:     "PO-" + time.Now().Format("20060102150405"),
		VendorID:     req.VendorID,
		Description:  req.Description,
		TotalAmount:  req.TotalAmount,
		Status:       "draft",
	}

	if err := h.repo.CreatePO(&po); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan purchase order").WithContext(c))
		return
	}

	for _, item := range req.Items {
		poi := domain.PurchaseOrderItem{
			ID:          fmt.Sprintf("POI-%s", time.Now().Format("20060102150405")),
			POID:        po.ID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.Quantity * item.UnitPrice,
		}
		h.repo.CreatePOItem(&poi)
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.purchase_order.create",
		Module:       "finance",
		ResourceType: "purchase_order",
		ResourceID:   po.ID,
		NewValue:     po,
	})

	c.JSON(http.StatusCreated, sharederr.Success(po).WithContext(c))
}

// ApprovePurchaseOrder handles POST /api/v1/finance/purchase-orders/:id/approve
func (h *PurchaseOrderHandler) ApprovePurchaseOrder(c *gin.Context) {
	id := c.Param("id")

	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	if err := h.repo.UpdatePOStatus(id, "approved", actor); err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyetujui purchase order").WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.purchase_order.approve",
		Module:       "finance",
		ResourceType: "purchase_order",
		ResourceID:   id,
	})

	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Purchase order berhasil disetujui").WithContext(c))
}
