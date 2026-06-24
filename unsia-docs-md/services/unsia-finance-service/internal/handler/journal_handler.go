package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
)

// JournalHandler handles journal-related endpoints
type JournalHandler struct {
	*FinanceHandler
}

// NewJournalHandler creates a new JournalHandler
func NewJournalHandler(fh *FinanceHandler) *JournalHandler {
	return &JournalHandler{FinanceHandler: fh}
}

// GetJournals handles GET /api/v1/finance/journals
func (h *JournalHandler) GetJournals(c *gin.Context) {
	filter := repository.JournalFilter{
		SourceType: c.Query("source_type"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
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

	result, err := h.repo.GetJournals(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil daftar journal").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(result).WithContext(c))
}

// GetJournalDetail handles GET /api/v1/finance/journals/:id
func (h *JournalHandler) GetJournalDetail(c *gin.Context) {
	id := c.Param("id")

	journal, err := h.repo.GetJournalByID(id)
	if err != nil || journal == nil {
		c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Journal tidak ditemukan").WithContext(c))
		return
	}

	entries, err := h.repo.GetJournalEntries(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal mengambil journal entries").WithContext(c))
		return
	}

	c.JSON(http.StatusOK, sharederr.Success(gin.H{
		"journal": journal,
		"entries": entries,
	})).WithContext(c))
}
