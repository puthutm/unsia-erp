package service

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page    int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// DefaultPagination returns default pagination
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:   1,
		Limit:  50,
		Offset: 0,
	}
}

// ParsePagination parses pagination from query params
func ParsePagination(c *gin.Context) PaginationParams {
	p := DefaultPagination()

	if page := c.Query("page"); page != "" {
		if pNum, err := strconv.Atoi(page); err == nil && pNum > 0 {
			p.Page = pNum
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			p.Limit = l
			if p.Limit > 100 {
				p.Limit = 100 // Max limit
			}
		}
	}

	// Calculate offset
	p.Offset = (p.Page - 1) * p.Limit

	return p
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page      int       `json:"page"`
	Limit    int       `json:"limit"`
	Total    int64     `json:"total"`
	TotalPage int     `json:"total_page"`
}

// BuildPaginationResponse builds paginated response
func BuildPaginationResponse(data interface{}, total int64, params PaginationParams) PaginationResponse {
	totalPage := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPage++
	}

	return PaginationResponse{
		Data:       data,
		Page:      params.Page,
		Limit:    params.Limit,
		Total:    total,
		TotalPage: totalPage,
	}
}

// CursorPagination for cursor-based pagination
type CursorPagination struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

// EncodeCursor encodes cursor from ID
func EncodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

// DecodeCursor decodes cursor to ID
func DecodeCursor(cursor string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	HasNext     bool   `json:"has_next"`
	HasPrev     bool   `json:"has_prev"`
	NextCursor  string `json:"next_cursor,omitempty"`
	PrevCursor  string `json:"prev_cursor,omitempty"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
}

// BuildPaginationMeta builds pagination metadata
func BuildPaginationMeta(items [] interface{}, total int64, params PaginationParams) PaginationMeta {
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}

	meta := PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}

	// Add next cursor
	if meta.HasNext && len(items) > 0 {
		lastItem := items[len(items)-1]
		if id := getID(lastItem); id != "" {
			meta.NextCursor = EncodeCursor(id)
		}
	}

	// Add prev cursor
	if meta.HasPrev && params.Page > 2 {
		firstItem := items[0]
		if id := getID(firstItem); id != "" {
			meta.PrevCursor = EncodeCursor(id)
		}
	}

	return meta
}

func getID(item interface{}) string {
	// Simple reflection to get ID field
	return fmt.Sprintf("%v", item)
}
