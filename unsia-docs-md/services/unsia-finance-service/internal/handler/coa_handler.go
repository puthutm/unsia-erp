package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)

// GetCoaAccounts handles GET /api/v1/finance/coa-accounts
func (h *FinanceHandler) GetCoaAccounts(c *gin.Context) {
	var active *bool
	if act := c.Query("is_active"); act != "" {
		val := act == "true"
		active = &val
	}

	filter := repository.CoaAccountFilter{
		IsActive: active,
		Search:   c.Query("search"),
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

	result, err := h.repo.GetCoaAccounts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar COA").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}
