package sharederr

import (
	"context"
	"time"
)

type ResponseMeta struct {
	TraceID   string    `json:"trace_id"`
	Timestamp time.Time `json:"timestamp"`
}

type ErrorDetails struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

type Response struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
	Error      *ErrorDetails   `json:"error,omitempty"`
	Meta       ResponseMeta    `json:"meta"`
}

// GetTraceID retrieves trace ID from context
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val, ok := ctx.Value("trace_id").(string); ok {
		return val
	}
	if val, ok := ctx.Value("x-correlation-id").(string); ok {
		return val
	}
	// For gin.Context compatibility
	if gc, ok := ctx.(interface{ GetString(string) string }); ok {
		if tid := gc.GetString("trace_id"); tid != "" {
			return tid
		}
		if cid := gc.GetString("x-correlation-id"); cid != "" {
			return cid
		}
	}
	return ""
}

func Success(data interface{}) Response {
	return Response{
		Success: true,
		Message: "Request processed successfully",
		Data:    data,
		Meta: ResponseMeta{
			Timestamp: time.Now(),
		},
	}
}

func SuccessWithMessage(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta: ResponseMeta{
			Timestamp: time.Now(),
		},
	}
}

func Error(code string, message string) Response {
	return Response{
		Success: false,
		Error: &ErrorDetails{
			Code:    code,
			Message: message,
			Details: map[string]interface{}{},
		},
		Meta: ResponseMeta{
			Timestamp: time.Now(),
		},
	}
}

func ValidationError(details interface{}) Response {
	return Response{
		Success: false,
		Error: &ErrorDetails{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: details,
		},
		Meta: ResponseMeta{
			Timestamp: time.Now(),
		},
	}
}

func (r Response) WithTraceID(traceID string) Response {
	r.Meta.TraceID = traceID
	return r
}

func (r Response) WithContext(ctx context.Context) Response {
	if ctx != nil {
		r.Meta.TraceID = GetTraceID(ctx)
	}
	return r
}

func (r Response) WithPagination(page, limit, total int) Response {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	r.Pagination = &PaginationInfo{
		Page:       page,
		Limit:      limit,
		TotalRows:  int64(total),
		TotalPages: totalPages,
	}
	return r
}
