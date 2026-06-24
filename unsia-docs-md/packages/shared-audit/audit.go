package sharedaudit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	sharedauth "github.com/unsia-erp/shared-auth"
)

type AuditEntry struct {
	ID            string    `json:"id"`             // UUID
	Actor         string    `json:"actor"`          // user_id
	ActiveRole    string    `json:"active_role"`    // role yang digunakan saat aksi
	Action        string    `json:"action"`         // contoh: "pmb.applicant.handover"
	Module        string    `json:"module"`         // contoh: "pmb"
	ResourceType  string    `json:"resource_type"`  // contoh: "applicant"
	ResourceID    string    `json:"resource_id"`    // UUID resource
	OldValue      interface{} `json:"old_value"`      // JSON sebelum perubahan (nullable)
	NewValue      interface{} `json:"new_value"`      // JSON setelah perubahan (nullable)
	Reason        string    `json:"reason"`         // alasan jika aksi sensitif (nullable)
	CorrelationID string    `json:"correlation_id"` // X-Correlation-Id dari request
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	OccurredAt    time.Time `json:"occurred_at"`
}

type Writer interface {
	Write(ctx context.Context, entry AuditEntry) error
}

type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// StdoutWriter is the default writer that prints logs to standard output
type StdoutWriter struct {
	logger *log.Logger
}

func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{
		logger: log.New(os.Stdout, "[AUDIT] ", log.LstdFlags),
	}
}

func (s *StdoutWriter) Write(ctx context.Context, entry AuditEntry) error {
	bytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	s.logger.Println(string(bytes))
	return nil
}

var (
	writerMutex sync.RWMutex
	activeWriter Writer = NewStdoutWriter()
)

// RegisterWriter overrides the default writer
func RegisterWriter(w Writer) {
	writerMutex.Lock()
	defer writerMutex.Unlock()
	activeWriter = w
}

// GetWriter gets the active writer
func GetWriter() Writer {
	writerMutex.RLock()
	defer writerMutex.RUnlock()
	return activeWriter
}

// Log writes an audit log, automatically extracting metadata from context if missing
func Log(ctx context.Context, entry AuditEntry) {
	if entry.OccurredAt.IsZero() {
		entry.OccurredAt = time.Now()
	}

	// Try extracting trace / correlation ID
	if entry.CorrelationID == "" {
		if cid, ok := ctx.Value("correlation_id").(string); ok {
			entry.CorrelationID = cid
		} else if cid, ok := ctx.Value("x-correlation-id").(string); ok {
			entry.CorrelationID = cid
		} else if gc, ok := ctx.(interface{ GetString(string) string }); ok {
			if cid := gc.GetString("x-correlation-id"); cid != "" {
				entry.CorrelationID = cid
			} else if tid := gc.GetString("trace_id"); tid != "" {
				entry.CorrelationID = tid
			}
		}
	}

	// Try extracting actor / role from auth claims
	if entry.Actor == "" || entry.ActiveRole == "" {
		var claims *sharedauth.Claims
		if val, ok := ctx.Value("claims").(*sharedauth.Claims); ok {
			claims = val
		} else if val, ok := ctx.Value("user_claims").(*sharedauth.Claims); ok {
			claims = val
		}

		if claims != nil {
			if entry.Actor == "" {
				entry.Actor = claims.Subject
			}
			if entry.ActiveRole == "" {
				entry.ActiveRole = claims.ActiveRole
			}
		}
	}

	// Try extracting IP Address and User Agent
	if entry.IPAddress == "" {
		if ip, ok := ctx.Value("ip_address").(string); ok {
			entry.IPAddress = ip
		}
	}
	if entry.UserAgent == "" {
		if ua, ok := ctx.Value("user_agent").(string); ok {
			entry.UserAgent = ua
		}
	}

	writer := GetWriter()
	_ = writer.Write(ctx, entry)
}

// SaveToSQL is a helper to insert an audit entry into a SQL database
func SaveToSQL(ctx context.Context, exec DBExecutor, entry AuditEntry) error {
	var oldJSON, newJSON []byte
	var err error

	if entry.OldValue != nil {
		oldJSON, err = json.Marshal(entry.OldValue)
		if err != nil {
			return fmt.Errorf("failed to marshal OldValue: %w", err)
		}
	}

	if entry.NewValue != nil {
		newJSON, err = json.Marshal(entry.NewValue)
		if err != nil {
			return fmt.Errorf("failed to marshal NewValue: %w", err)
		}
	}

	query := `
		INSERT INTO audit_logs (
			id, actor, active_role, action, module, resource_type, resource_id,
			old_value, new_value, reason, correlation_id, ip_address, user_agent, occurred_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14
		)
	`

	_, err = exec.ExecContext(
		ctx, query,
		entry.ID,
		entry.Actor,
		entry.ActiveRole,
		entry.Action,
		entry.Module,
		entry.ResourceType,
		entry.ResourceID,
		sql.NullString{String: string(oldJSON), Valid: oldJSON != nil},
		sql.NullString{String: string(newJSON), Valid: newJSON != nil},
		sql.NullString{String: entry.Reason, Valid: entry.Reason != ""},
		entry.CorrelationID,
		entry.IPAddress,
		entry.UserAgent,
		entry.OccurredAt,
	)

	return err
}
