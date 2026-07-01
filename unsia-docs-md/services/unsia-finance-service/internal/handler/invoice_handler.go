package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sharedaudit "github.com/unsia-erp/shared-audit"
	sharedauth "github.com/unsia-erp/shared-auth"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedidempotency "github.com/unsia-erp/shared-idempotency"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
	"github.com/unsia-erp/unsia-finance-service/internal/service"
)

// CreateInvoice handles POST /api/v1/finance/invoices
func (h *FinanceHandler) CreateInvoice(c *gin.Context) {
	// Idempotency check
	idempotencyKey := c.GetHeader("Idempotency-Key")
	var useIdempotency bool
	if idempotencyKey != "" {
		useIdempotency = true
		idempotencyKey = "finance:invoice:create:" + idempotencyKey
		cachedResponse, exists, err := sharedidempotency.CheckAndLock(c.Request.Context(), idempotencyKey, 24*time.Hour)
		if err != nil {
			if err == sharedidempotency.ErrConcurrentRequest {
				c.JSON(http.StatusConflict, sharederr.Error("IDEMPOTENCY_KEY_IN_PROGRESS", "Request is currently being processed").WithContext(c))
				return
			}
			c.JSON(http.StatusInternalServerError, sharederr.Error("INTERNAL_ERROR", err.Error()).WithContext(c))
			return
		}
		if exists {
			var cached interface{}
			if json.Unmarshal([]byte(cachedResponse), &cached) == nil {
				c.JSON(http.StatusOK, cached)
				return
			}
		}
	}

	var req InvoiceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if useIdempotency {
			_ = sharedidempotency.SaveFailure(c.Request.Context(), idempotencyKey, err.Error())
		}
		c.JSON(http.StatusUnprocessableEntity, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	var dueTime *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			if useIdempotency {
				_ = sharedidempotency.SaveFailure(c.Request.Context(), idempotencyKey, "invalid due date format")
			}
			c.JSON(http.StatusUnprocessableEntity, sharederr.Error("INVALID_DUE_DATE", "due_date format must be YYYY-MM-DD").WithContext(c))
			return
		}
		dueTime = &parsed
	}

	// Validate items: minimal satu items dengan payment_component_id dan final_amount positif
	if len(req.Items) == 0 {
		if useIdempotency {
			_ = sharedidempotency.SaveFailure(c.Request.Context(), idempotencyKey, "at least one item is required")
		}
		c.JSON(http.StatusUnprocessableEntity, sharederr.Error("VALIDATION_ERROR", "At least one item is required").WithContext(c))
		return
	}

	var serviceItems []service.InvoiceItemRequest
	for _, it := range req.Items {
		finalAmt := it.FinalAmount
		if finalAmt == 0 {
			finalAmt = it.Amount
		}
		if finalAmt <= 0 {
			if useIdempotency {
				_ = sharedidempotency.SaveFailure(c.Request.Context(), idempotencyKey, "final_amount must be positive")
			}
			c.JSON(http.StatusUnprocessableEntity, sharederr.Error("VALIDATION_ERROR", "final_amount must be positive").WithContext(c))
			return
		}
		serviceItems = append(serviceItems, service.InvoiceItemRequest{
			PaymentComponentID: it.PaymentComponentID,
			Description:        it.Description,
			Amount:             it.Amount,
			DiscountAmount:     it.DiscountAmount,
			FinalAmount:        finalAmt,
		})
	}

	serviceReq := service.InvoiceCreateRequest{
		PayerType:        req.PayerType,
		PayerRefID:       req.PayerRefID,
		AcademicPeriodID: req.AcademicPeriodID,
		DueDate:          dueTime,
		Items:            serviceItems,
		SourceModule:     req.SourceModule,
		SourceRefID:      req.SourceRefID,
	}

	invoice, err := h.InvoiceService.CreateInvoice(&serviceReq, c)
	if err != nil {
		if useIdempotency {
			_ = sharedidempotency.SaveFailure(c.Request.Context(), idempotencyKey, err.Error())
		}
		if err.Error() == "PAYMENT_COMPONENT_NOT_FOUND" {
			c.JSON(http.StatusUnprocessableEntity, sharederr.Error("PAYMENT_COMPONENT_NOT_FOUND", "Payment component not found or inactive").WithContext(c))
			return
		}
		if strings.Contains(err.Error(), "UPSTREAM_SERVICE_UNAVAILABLE") {
			c.JSON(http.StatusServiceUnavailable, sharederr.Error("UPSTREAM_SERVICE_UNAVAILABLE", "Reference service is currently unavailable").WithContext(c))
			return
		}
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", err.Error()).WithContext(c))
		return
	}

	sharedaudit.Log(c, sharedaudit.AuditEntry{
		Action:       "finance.invoice.create",
		Module:       "finance",
		ResourceType: "invoice",
		ResourceID:   invoice.ID,
		NewValue:     invoice,
	})

	respPayload := sharederr.Success(mapInvoiceToDetailResponse(*invoice)).WithContext(c)

	if useIdempotency {
		respBytes, _ := json.Marshal(respPayload)
		_ = sharedidempotency.SaveResponse(c.Request.Context(), idempotencyKey, string(respBytes), 24*time.Hour)
	}

	c.JSON(http.StatusCreated, respPayload)
}

// GetInvoice handles GET /api/v1/finance/invoices/:id
func (h *FinanceHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	inv, err := h.repo.GetInvoiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil invoice").WithContext(c))
		return
	}
	if inv == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Invoice tidak ditemukan").WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.Success(mapInvoiceToDetailResponse(*inv)).WithContext(c))
}

// GetInvoices handles GET /api/v1/finance/invoices
func (h *FinanceHandler) GetInvoices(c *gin.Context) {
	filter := repository.InvoiceListFilter{
		Status:           c.Query("status"),
		TargetType:       c.Query("target_type"),
		StudentID:       c.Query("student_id"),
		ApplicantID:     c.Query("applicant_id"),
		AcademicPeriodID: c.Query("academic_period_id"),
		DueDateFrom:     c.Query("due_date_from"),
		DueDateTo:       c.Query("due_date_to"),
		Search:          c.Query("search"),
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

	result, err := h.repo.GetInvoices(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar invoice").WithContext(c))
		return
	}

	invoices, ok := result.Data.([]domain.Invoice)
	var mapped []InvoiceBriefResponse
	if ok {
		mapped = make([]InvoiceBriefResponse, len(invoices))
		for i, inv := range invoices {
			mapped[i] = mapInvoiceToBriefResponse(inv)
		}
	} else {
		mapped = []InvoiceBriefResponse{}
	}
	result.Data = mapped

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// Helper mapping functions
func mapInvoiceToBriefResponse(inv domain.Invoice) InvoiceBriefResponse {
	return InvoiceBriefResponse{
		ID:               inv.ID,
		InvoiceNumber:    inv.InvoiceNumber,
		TargetType:       inv.TargetType,
		ApplicantID:      inv.ApplicantID,
		StudentID:        inv.StudentID,
		AcademicPeriodID: inv.AcademicPeriodID,
		TotalAmount:      inv.TotalAmount,
		PaidAmount:       inv.PaidAmount,
		Status:           inv.Status,
		DueDate:          inv.DueDate,
		CreatedAt:        inv.CreatedAt,
	}
}

func mapInvoiceToDetailResponse(inv domain.Invoice) InvoiceDetailResponse {
	items := make([]InvoiceItemResponse, len(inv.Items))
	for i, item := range inv.Items {
		items[i] = InvoiceItemResponse{
			ID:                 item.ID,
			PaymentComponentID: item.PaymentComponentID,
			Description:        item.Description,
			Amount:             item.Amount,
			DiscountAmount:     item.DiscountAmount,
			FinalAmount:        item.FinalAmount,
		}
	}

	payments := make([]PaymentResponse, len(inv.Payments))
	for i, pay := range inv.Payments {
		payments[i] = PaymentResponse{
			ID:              pay.ID,
			PaymentMethodID: pay.PaymentMethodID,
			PaymentNumber:   pay.PaymentNumber,
			Amount:          pay.Amount,
			PaymentStatus:   pay.PaymentStatus,
			PaidAt:          pay.PaidAt,
		}
	}

	return InvoiceDetailResponse{
		ID:               inv.ID,
		InvoiceNumber:    inv.InvoiceNumber,
		TargetType:       inv.TargetType,
		ApplicantID:      inv.ApplicantID,
		StudentID:        inv.StudentID,
		AcademicPeriodID: inv.AcademicPeriodID,
		TotalAmount:      inv.TotalAmount,
		PaidAmount:       inv.PaidAmount,
		Status:           inv.Status,
		DueDate:          inv.DueDate,
		CreatedAt:        inv.CreatedAt,
		UpdatedAt:        inv.UpdatedAt,
		Items:            items,
		Payments:         payments,
	}
}

// IssueInvoice handles POST /api/v1/finance/invoices/:id/issue
func (h *FinanceHandler) IssueInvoice(c *gin.Context) {
	id := c.Param("id")
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	err := h.InvoiceService.UpdateInvoiceStatus(id, "ISSUED", actor, c)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Invoice tidak ditemukan").WithContext(c))
			return
		}
		c.JSON(http.StatusConflict, sharederr.Error("INVALID_STATUS_TRANSITION", err.Error()).WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Invoice status updated to ISSUED").WithContext(c))
}

// CancelInvoice handles POST /api/v1/finance/invoices/:id/cancel
func (h *FinanceHandler) CancelInvoice(c *gin.Context) {
	id := c.Param("id")
	claimsVal, _ := c.Get("claims")
	claims, _ := claimsVal.(*sharedauth.Claims)
	actor := ""
	if claims != nil {
		actor = claims.Subject
	}

	err := h.InvoiceService.UpdateInvoiceStatus(id, "CANCELLED", actor, c)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Invoice tidak ditemukan").WithContext(c))
			return
		}
		c.JSON(http.StatusConflict, sharederr.Error("INVALID_STATUS_TRANSITION", err.Error()).WithContext(c))
		return
	}
	c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Invoice status updated to CANCELLED").WithContext(c))
}
