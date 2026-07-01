package handler

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	"github.com/unsia-erp/unsia-finance-service/internal/domain"
	"github.com/unsia-erp/unsia-finance-service/internal/infrastructure/repository"
	"gorm.io/gorm"
)

type JournalCreateRequest struct {
	JournalDate string                `json:"journal_date" binding:"required"`
	Description string                `json:"description"`
	SourceType  string                `json:"source_type" binding:"required"`
	SourceID    *string               `json:"source_id"`
	Entries     []JournalEntryRequest `json:"entries" binding:"required,gt=1"`
}

type JournalEntryRequest struct {
	CoaAccountID string  `json:"coa_account_id" binding:"required"`
	Debit        float64 `json:"debit"`
	Credit       float64 `json:"credit"`
}

// GetJournals handles GET /api/v1/finance/journals
func (h *FinanceHandler) GetJournals(c *gin.Context) {
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
func (h *FinanceHandler) GetJournalDetail(c *gin.Context) {
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
	}).WithContext(c))
}

// CreateJournal handles POST /api/v1/finance/journals
func (h *FinanceHandler) CreateJournal(c *gin.Context) {
	var req JournalCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.JournalDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, sharederr.Error("INVALID_DATE", "journal_date format must be YYYY-MM-DD").WithContext(c))
		return
	}

	var totalDebit, totalCredit float64
	for _, entry := range req.Entries {
		totalDebit += entry.Debit
		totalCredit += entry.Credit
	}

	if math.Abs(totalDebit-totalCredit) > 0.001 {
		c.JSON(http.StatusUnprocessableEntity, sharederr.Error("JOURNAL_UNBALANCED", "Total debit must equal total credit").WithContext(c))
		return
	}

	journalNumber := "JV-" + time.Now().Format("20060102150405")
	journal := domain.Journal{
		JournalNumber: journalNumber,
		JournalDate:   parsedDate,
		Description:   req.Description,
		SourceType:    req.SourceType,
		SourceID:      req.SourceID,
	}

	errTx := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&journal).Error; err != nil {
			return err
		}

		for _, entryReq := range req.Entries {
			entry := domain.JournalEntry{
				JournalID:    journal.ID,
				CoaAccountID: &entryReq.CoaAccountID,
				Debit:        entryReq.Debit,
				Credit:       entryReq.Credit,
			}
			if err := tx.Create(&entry).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if errTx != nil {
		c.JSON(http.StatusInternalServerError, sharederr.Error("DB_ERROR", "Gagal menyimpan journal").WithContext(c))
		return
	}

	c.JSON(http.StatusCreated, sharederr.Success(journal).WithContext(c))
}
