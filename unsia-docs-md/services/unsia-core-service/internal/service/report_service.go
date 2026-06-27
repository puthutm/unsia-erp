package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// ReportRequest represents report request
type ReportRequest struct {
	StartDate   time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Module     string    `json:"module"`
	Format     string    `json:"format"` // json, csv, xlsx
	GroupBy    string    `json:"group_by"`
}

// GenerateReport generates report based on filters
func (s *ReportService) GenerateReport(req ReportRequest) ([]map[string]interface{}, error) {
	// Based on module, generate different reports
	switch req.Module {
	case "user":
		return s.userReport(req)
	case "audit":
		return s.auditReport(req)
	case "finance":
		return s.financeReport(req)
	default:
		return nil, fmt.Errorf("unknown module: %s", req.Module)
	}
}

func (s *ReportService) userReport(req ReportRequest) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := s.db.Model(&gin.H{}).Table("users")

	if !req.StartDate.IsZero() {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if !req.EndDate.IsZero() {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		item := make(map[string]interface{})
		for i, col := range columns {
			item[col] = values[i]
		}
		results = append(results, item)
	}

	return results, nil
}

func (s *ReportService) auditReport(req ReportRequest) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := s.db.Model(&gin.H{}).Table("audit_logs")

	if !req.StartDate.IsZero() {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if !req.EndDate.IsZero() {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	rows, err := query.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		item := make(map[string]interface{})
		for i, col := range columns {
			item[col] = values[i]
		}
		results = append(results, item)
	}

	return results, nil
}

func (s *ReportService) financeReport(req ReportRequest) ([]map[string]interface{}, error) {
	// Simplified - would query finance tables
	return []map[string]interface{}{}, nil
}

// ExportToJSON exports data to JSON
func ExportToJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// ExportToCSV exports data to CSV
func ExportToCSV(data []map[string]interface{}) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	// Get headers from first row
	var headers []string
	for k := range data[0] {
		headers = append(headers, k)
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Write header
	writer.Write(headers)

	// Write rows
	for _, row := range data {
		var record []string
		for _, h := range headers {
			val := fmt.Sprintf("%v", row[h])
			record = append(record, val)
		}
		writer.Write(record)
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

// Statistics represents system statistics
type Statistics struct {
	TotalUsers      int64   `json:"total_users"`
	ActiveUsers    int64   `json:"active_users"`
	TotalRoles    int64   `json:"total_roles"`
	TotalApplications int64 `json:"total_applications"`
	TotalAuditLogs int64   `json:"total_audit_logs"`
	ReportDate    time.Time `json:"report_date"`
}

// GetStatistics gets system statistics
func (s *ReportService) GetStatistics() (Statistics, error) {
	stats := Statistics{ReportDate: time.Now()}

	// Count users
	s.db.Model(&gin.H{}).Table("users").Count(&stats.TotalUsers)

	// Count active users
	s.db.Model(&gin.H{}).Table("users").Where("is_active = ?", true).Count(&stats.ActiveUsers)

	// Count roles
	s.db.Model(&gin.H{}).Table("roles").Count(&stats.TotalRoles)

	// Count applications
	s.db.Model(&gin.H{}).Table("applications").Count(&stats.TotalApplications)

	// Count audit logs
	s.db.Model(&gin.H{}).Table("audit_logs").Count(&stats.TotalAuditLogs)

	return stats, nil
}

// DailyActivityReport gets daily activity report
type DailyActivityReport struct {
	Date          time.Time `json:"date"`
	LoginCount    int64     `json:"login_count"`
	ActionCount  int64     `json:"action_count"`
	NewUsers     int64     `json:"new_users"`
}

func (s *ReportService) GetDailyActivityReport(days int) ([]DailyActivityReport, error) {
	var reports []DailyActivityReport

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		report := DailyActivityReport{
			Date: date,
		}

		// Get login count for this day
		s.db.Model(&gin.H{}).Table("audit_logs").
			Where("action LIKE ? AND created_at::date = ?", "%login%", dateStr).
			Count(&report.LoginCount)

		// Get action count
		s.db.Model(&gin.H{}).Table("audit_logs").
			Where("created_at::date = ?", dateStr).
			Count(&report.ActionCount)

		// Get new users
		s.db.Model(&gin.H{}).Table("users").
			Where("created_at::date = ?", dateStr).
			Count(&report.NewUsers)

		reports = append(reports, report)
	}

	return reports, nil
}
