package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	sharederr "github.com/unsia-erp/shared-errorenvelope"
	sharedevent "github.com/unsia-erp/shared-event"
	sharedobservability "github.com/unsia-erp/shared-observability"
)

type OutboxRow struct {
	ID            string
	EventName     string
	EventVersion  string
	EventKey      string
	EventType     string
	AggregateType string
	AggregateID   string
	Payload       []byte
	CorrelationID string
	CausationID   string
	Status        string
	OccurredAt    time.Time
}

var eventConsumers = map[string][]string{
	"core.person_updated":              {"crm", "pmb", "hris", "academic", "portal"},
	"reference.study_program_updated":  {"pmb", "academic", "hris", "lms", "portal"},
	"reference.academic_period_updated": {"pmb", "finance", "academic", "lms", "assessment", "portal"},
	"crm.lead_qualified":               {"pmb"},
	"pmb.applicant_created":            {"finance", "assessment", "portal"},
	"finance.invoice_created":          {"pmb", "portal"},
	"finance.payment_paid":             {"pmb", "academic", "portal"},
	"finance.clearance_changed":        {"pmb", "academic", "lms", "portal"},
	"pmb.ready_for_academic":           {"academic"},
	"academic.student_created":          {"pmb", "finance", "lms", "portal"},
	"academic.class_opened":            {"lms", "portal"},
	"academic.krs_approved":            {"lms", "finance", "portal"},
	"lms.grade_input_submitted":        {"academic"},
	"assessment.result_calculated":     {"pmb", "lms", "academic", "portal"},
}

type ResolveMismatchRequest struct {
	Status string `json:"status" binding:"required,oneof=CORRECTED IGNORED"`
	Reason string `json:"reason" binding:"required"`
}

func startHTTPServer(dbs map[string]*sql.DB, logger zerolog.Logger) {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())

	// Endpoint to get all event contracts
	r.GET("/api/v1/integration/event-contracts", func(c *gin.Context) {
		c.JSON(http.StatusOK, sharederr.Success(gin.H{
			"event_consumers": eventConsumers,
		}).WithContext(c))
	})

	// Endpoint to get a specific event contract
	r.GET("/api/v1/integration/event-contracts/:event_name", func(c *gin.Context) {
		eventName := c.Param("event_name")
		consumers, exists := eventConsumers[eventName]
		if !exists {
			c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Event contract not found").WithContext(c))
			return
		}
		c.JSON(http.StatusOK, sharederr.Success(gin.H{
			"event_name": eventName,
			"consumers":  consumers,
		}).WithContext(c))
	})

	// Endpoint to get outbox events from all databases
	r.GET("/api/v1/integration/outbox-events", func(c *gin.Context) {
		statusFilter := c.Query("status")
		serviceFilter := c.Query("service")

		var allEvents []gin.H
		for name, db := range dbs {
			if serviceFilter != "" && serviceFilter != name {
				continue
			}

			query := `
				SELECT id, event_name, event_version, event_key, event_type, aggregate_type, aggregate_id, status, occurred_at 
				FROM outbox_events
			`
			var args []interface{}
			if statusFilter != "" {
				query += " WHERE status = $1"
				args = append(args, statusFilter)
			}
			query += " ORDER BY occurred_at DESC LIMIT 50"

			rows, err := db.QueryContext(c.Request.Context(), query, args...)
			if err != nil {
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var id, nameStr, version, key, evtType, aggType, aggID, statusVal string
				var occurredAt time.Time
				if err := rows.Scan(&id, &nameStr, &version, &key, &evtType, &aggType, &aggID, &statusVal, &occurredAt); err == nil {
					allEvents = append(allEvents, gin.H{
						"service":        name,
						"id":             id,
						"event_name":     nameStr,
						"event_version":  version,
						"event_key":      key,
						"event_type":     evtType,
						"aggregate_type": aggType,
						"aggregate_id":   aggID,
						"status":         statusVal,
						"occurred_at":    occurredAt,
					})
				}
			}
		}

		sort.Slice(allEvents, func(i, j int) bool {
			t1 := allEvents[i]["occurred_at"].(time.Time)
			t2 := allEvents[j]["occurred_at"].(time.Time)
			return t1.After(t2)
		})

		c.JSON(http.StatusOK, sharederr.Success(allEvents).WithContext(c))
	})

	// Endpoint to get inbox events from all databases
	r.GET("/api/v1/integration/inbox-events", func(c *gin.Context) {
		statusFilter := c.Query("status")
		serviceFilter := c.Query("service")

		var allEvents []gin.H
		for name, db := range dbs {
			if serviceFilter != "" && serviceFilter != name {
				continue
			}

			query := `
				SELECT id, event_name, event_version, event_key, publisher_module, consumer_module, aggregate_type, aggregate_id, status, received_at 
				FROM inbox_events
			`
			var args []interface{}
			if statusFilter != "" {
				query += " WHERE status = $1"
				args = append(args, statusFilter)
			}
			query += " ORDER BY received_at DESC LIMIT 50"

			rows, err := db.QueryContext(c.Request.Context(), query, args...)
			if err != nil {
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var id, nameStr, version, key, pubMod, consMod, aggType, aggID, statusVal string
				var receivedAt time.Time
				if err := rows.Scan(&id, &nameStr, &version, &key, &pubMod, &consMod, &aggType, &aggID, &statusVal, &receivedAt); err == nil {
					allEvents = append(allEvents, gin.H{
						"service":          name,
						"id":               id,
						"event_name":       nameStr,
						"event_version":    version,
						"event_key":        key,
						"publisher_module": pubMod,
						"consumer_module":  consMod,
						"aggregate_type":   aggType,
						"aggregate_id":     aggID,
						"status":           statusVal,
						"received_at":      receivedAt,
					})
				}
			}
		}

		sort.Slice(allEvents, func(i, j int) bool {
			t1 := allEvents[i]["received_at"].(time.Time)
			t2 := allEvents[j]["received_at"].(time.Time)
			return t1.After(t2)
		})

		c.JSON(http.StatusOK, sharederr.Success(allEvents).WithContext(c))
	})

	// Endpoint to get DLQ events from all databases
	r.GET("/api/v1/integration/dlq-events", func(c *gin.Context) {
		var allDLQ []gin.H
		for name, db := range dbs {
			// Query outbox DLQ
			rows, err := db.QueryContext(c.Request.Context(), `
				SELECT id, event_name, event_version, event_key, event_type, aggregate_type, aggregate_id, status, COALESCE(last_error, ''), COALESCE(dead_letter_at, NOW()) 
				FROM outbox_events 
				WHERE status = 'DLQ'
			`)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var id, nameStr, version, key, evtType, aggType, aggID, statusVal, lastErr string
					var deadLetterAt time.Time
					if err := rows.Scan(&id, &nameStr, &version, &key, &evtType, &aggType, &aggID, &statusVal, &lastErr, &deadLetterAt); err == nil {
						allDLQ = append(allDLQ, gin.H{
							"service":        name,
							"type":           "outbox",
							"id":             id,
							"event_name":     nameStr,
							"event_version":  version,
							"event_key":      key,
							"event_type":     evtType,
							"aggregate_type": aggType,
							"aggregate_id":   aggID,
							"status":         statusVal,
							"last_error":     lastErr,
							"dead_letter_at": deadLetterAt,
						})
					}
				}
			}

			// Query inbox DLQ
			rows2, err := db.QueryContext(c.Request.Context(), `
				SELECT id, event_name, event_version, event_key, publisher_module, consumer_module, aggregate_type, aggregate_id, status, COALESCE(last_error, ''), COALESCE(dead_letter_at, NOW()) 
				FROM inbox_events 
				WHERE status = 'DLQ'
			`)
			if err == nil {
				defer rows2.Close()
				for rows2.Next() {
					var id, nameStr, version, key, pubMod, consMod, aggType, aggID, statusVal, lastErr string
					var deadLetterAt time.Time
					if err := rows2.Scan(&id, &nameStr, &version, &key, &pubMod, &consMod, &aggType, &aggID, &statusVal, &lastErr, &deadLetterAt); err == nil {
						allDLQ = append(allDLQ, gin.H{
							"service":          name,
							"type":             "inbox",
							"id":               id,
							"event_name":       nameStr,
							"event_version":    version,
							"event_key":        key,
							"publisher_module": pubMod,
							"consumer_module":  consMod,
							"aggregate_type":   aggType,
							"aggregate_id":     aggID,
							"status":           statusVal,
							"last_error":       lastErr,
							"dead_letter_at":   deadLetterAt,
						})
					}
				}
			}
		}

		sort.Slice(allDLQ, func(i, j int) bool {
			t1 := allDLQ[i]["dead_letter_at"].(time.Time)
			t2 := allDLQ[j]["dead_letter_at"].(time.Time)
			return t1.After(t2)
		})

		c.JSON(http.StatusOK, sharederr.Success(allDLQ).WithContext(c))
	})

	// Endpoint to replay a DLQ event by event_key
	r.POST("/api/v1/integration/dlq-events/:event_key/replay", func(c *gin.Context) {
		eventKey := c.Param("event_key")
		if eventKey == "" {
			c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Event key is required").WithContext(c))
			return
		}

		found := false
		for name, db := range dbs {
			// Replay outbox
			res, err := db.ExecContext(c.Request.Context(), `
				UPDATE outbox_events 
				SET status = 'PENDING', dead_letter_at = NULL, last_error = NULL 
				WHERE event_key = $1 AND status = 'DLQ'
			`, eventKey)
			if err == nil {
				if affected, _ := res.RowsAffected(); affected > 0 {
					found = true
					logger.Info().Msgf("Replayed outbox event %s in service %s", eventKey, name)
				}
			}

			// Replay inbox
			res2, err := db.ExecContext(c.Request.Context(), `
				UPDATE inbox_events 
				SET status = 'RECEIVED', dead_letter_at = NULL, last_error = NULL 
				WHERE event_key = $1 AND status = 'DLQ'
			`, eventKey)
			if err == nil {
				if affected, _ := res2.RowsAffected(); affected > 0 {
					found = true
					logger.Info().Msgf("Replayed inbox event %s in service %s", eventKey, name)
				}
			}
		}

		if found {
			c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Event successfully queued for retry").WithContext(c))
		} else {
			c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "DLQ event with the specified key was not found").WithContext(c))
		}
	})

	// Endpoint to get reconciliation mismatches
	r.GET("/api/v1/integration/reconciliation-mismatches", func(c *gin.Context) {
		var allMismatches []gin.H
		for name, db := range dbs {
			rows, err := db.QueryContext(c.Request.Context(), `
				SELECT id, source_module, source_table, source_ref_id, consumer_module, consumer_table, consumer_ref_id, source_event_key, mismatch_type, source_value, snapshot_value, status, COALESCE(reason, ''), detected_at, corrected_at, ignored_at 
				FROM reconciliation_mismatch_logs 
				ORDER BY detected_at DESC LIMIT 50
			`)
			if err != nil {
				continue
			}
			defer rows.Close()

			for rows.Next() {
				var id, srcMod, srcTab, srcRefID, consMod, consTab, consRefID, srcEventKey, mismatchType, statusVal, reason string
				var sourceValue, snapshotValue []byte
				var detectedAt time.Time
				var correctedAt, ignoredAt *time.Time

				var srcRefPtr, consRefPtr *string
				err := rows.Scan(
					&id, &srcMod, &srcTab, &srcRefPtr, &consMod, &consTab, &consRefPtr, &srcEventKey, &mismatchType,
					&sourceValue, &snapshotValue, &statusVal, &reason, &detectedAt, &correctedAt, &ignoredAt,
				)
				if err == nil {
					if srcRefPtr != nil {
						srcRefID = *srcRefPtr
					}
					if consRefPtr != nil {
						consRefID = *consRefPtr
					}

					var srcJSON, snapJSON interface{}
					_ = json.Unmarshal(sourceValue, &srcJSON)
					_ = json.Unmarshal(snapshotValue, &snapJSON)

					allMismatches = append(allMismatches, gin.H{
						"service":          name,
						"id":               id,
						"source_module":    srcMod,
						"source_table":     srcTab,
						"source_ref_id":    srcRefID,
						"consumer_module":  consMod,
						"consumer_table":   consTab,
						"consumer_ref_id":  consRefID,
						"source_event_key": srcEventKey,
						"mismatch_type":    mismatchType,
						"source_value":     srcJSON,
						"snapshot_value":   snapJSON,
						"status":           statusVal,
						"reason":           reason,
						"detected_at":      detectedAt,
						"corrected_at":     correctedAt,
						"ignored_at":       ignoredAt,
					})
				}
			}
		}

		sort.Slice(allMismatches, func(i, j int) bool {
			t1 := allMismatches[i]["detected_at"].(time.Time)
			t2 := allMismatches[j]["detected_at"].(time.Time)
			return t1.After(t2)
		})

		c.JSON(http.StatusOK, sharederr.Success(allMismatches).WithContext(c))
	})

	// Endpoint to resolve reconciliation mismatches
	r.POST("/api/v1/integration/reconciliation-mismatches/:id/resolve", func(c *gin.Context) {
		mismatchID := c.Param("id")
		if mismatchID == "" {
			c.JSON(http.StatusBadRequest, sharederr.Error("BAD_REQUEST", "Mismatch ID is required").WithContext(c))
			return
		}

		var req ResolveMismatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, sharederr.ValidationError(err.Error()).WithContext(c))
			return
		}

		found := false
		for name, db := range dbs {
			query := `
				UPDATE reconciliation_mismatch_logs 
				SET status = $1, reason = $2, updated_at = NOW(),
					corrected_at = CASE WHEN $1 = 'CORRECTED' THEN NOW() ELSE corrected_at END,
					ignored_at = CASE WHEN $1 = 'IGNORED' THEN NOW() ELSE ignored_at END
				WHERE id = $3::uuid
			`
			res, err := db.ExecContext(c.Request.Context(), query, req.Status, req.Reason, mismatchID)
			if err == nil {
				if affected, _ := res.RowsAffected(); affected > 0 {
					found = true
					logger.Info().Msgf("Resolved mismatch %s in service %s as %s", mismatchID, name, req.Status)
					break
				}
			}
		}

		if found {
			c.JSON(http.StatusOK, sharederr.SuccessWithMessage(nil, "Reconciliation mismatch resolved successfully").WithContext(c))
		} else {
			c.JSON(http.StatusNotFound, sharederr.Error("NOT_FOUND", "Reconciliation mismatch with the specified ID was not found").WithContext(c))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8020"
	}

	logger.Info().Msgf("Integration HTTP Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logger.Error().Err(err).Msg("Failed to start Integration HTTP Server")
	}
}

func main() {
	_ = godotenv.Load()
	sharedobservability.InitLogger("integration-worker")
	logger := sharedobservability.Logger

	dbURLs := map[string]string{
		"core":       os.Getenv("CORE_DB_URL"),
		"reference":  os.Getenv("REFERENCE_DB_URL"),
		"crm":        os.Getenv("CRM_DB_URL"),
		"pmb":        os.Getenv("PMB_DB_URL"),
		"finance":    os.Getenv("FINANCE_DB_URL"),
		"academic":   os.Getenv("ACADEMIC_DB_URL"),
		"hris":       os.Getenv("HRIS_DB_URL"),
		"lms":        os.Getenv("LMS_DB_URL"),
		"assessment": os.Getenv("ASSESSMENT_DB_URL"),
		"portal":     os.Getenv("PORTAL_DB_URL"),
	}

	dbs := make(map[string]*sql.DB)
	for name, url := range dbURLs {
		if url == "" {
			url = "postgres://postgres:postgres@localhost:5432/" + name + "_db?sslmode=disable"
		}
		db, err := sql.Open("postgres", url)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to connect to DB: %s", name)
			continue
		}
		db.SetMaxOpenConns(5)
		db.SetMaxIdleConns(2)
		dbs[name] = db
		logger.Info().Msgf("Connected to database: %s", name)
	}

	// Start HTTP Server in a separate goroutine
	go startHTTPServer(dbs, logger)

	// Initialize RabbitMQ connection (Optional/Degraded fallback)
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	var rabbitConn *amqp.Connection
	var rabbitChan *amqp.Channel
	var amqpErr error

	rabbitConn, amqpErr = amqp.Dial(rabbitURL)
	if amqpErr != nil {
		logger.Warn().Err(amqpErr).Msg("RabbitMQ not reachable, operating in direct DB dispatch mode")
	} else {
		defer rabbitConn.Close()
		rabbitChan, amqpErr = rabbitConn.Channel()
		if amqpErr == nil {
			defer rabbitChan.Close()
			logger.Info().Msg("Successfully connected to RabbitMQ Broker")
		}
	}

	logger.Info().Msg("Integration Worker Daemon started successfully. Polling outboxes...")

	// Periodically poll all databases for outbox events
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	for range ticker.C {
		for sourceName, db := range dbs {
			pollAndProcessOutbox(ctx, sourceName, db, dbs, rabbitChan, logger)
		}
	}
}

func pollAndProcessOutbox(ctx context.Context, sourceName string, db *sql.DB, allDBs map[string]*sql.DB, amqpChan *amqp.Channel, logger zerolog.Logger) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, event_name, event_version, event_key, event_type, aggregate_type, aggregate_id, 
		       payload, correlation_id, causation_id, status, occurred_at 
		FROM outbox_events 
		WHERE status = 'PENDING' 
		ORDER BY occurred_at ASC 
		LIMIT 10
	`)
	if err != nil {
		// Log error but don't crash, the table might not exist yet if migrations haven't run
		return
	}
	defer rows.Close()

	var events []OutboxRow
	for rows.Next() {
		var r OutboxRow
		err := rows.Scan(
			&r.ID, &r.EventName, &r.EventVersion, &r.EventKey, &r.EventType, &r.AggregateType, &r.AggregateID,
			&r.Payload, &r.CorrelationID, &r.CausationID, &r.Status, &r.OccurredAt,
		)
		if err == nil {
			events = append(events, r)
		}
	}

	for _, event := range events {
		logger.Info().Msgf("[%s] Processing pending outbox event: %s", sourceName, event.EventName)

		// 1. Publish to RabbitMQ exchange if active
		if amqpChan != nil {
			err = amqpChan.PublishWithContext(ctx,
				"unsia-erp-exchange",
				event.EventName,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        event.Payload,
					MessageId:   event.ID,
					Headers: amqp.Table{
						"x-correlation-id": event.CorrelationID,
						"x-causation-id":   event.CausationID,
					},
				},
			)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to publish event to RabbitMQ exchange")
			}
		}

		// 2. Direct dispatch to target databases (for decoupling & resilience)
		consumers := eventConsumers[event.EventName]
		for _, consumerName := range consumers {
			targetDB, ok := allDBs[consumerName]
			if !ok {
				continue
			}

			logger.Info().Msgf("Forwarding event [%s] -> inbox of [%s]", event.EventName, consumerName)

			var payloadObj interface{}
			_ = json.Unmarshal(event.Payload, &payloadObj)

			env := sharedevent.EventEnvelope{
				ID:               event.ID,
				EventName:        event.EventName,
				EventVersion:     event.EventVersion,
				EventKey:         event.EventKey,
				PublisherService: sourceName + "-service",
				AggregateType:    event.AggregateType,
				AggregateID:      event.AggregateID,
				CorrelationID:    event.CorrelationID,
				CausationID:      event.CausationID,
				OccurredAt:       event.OccurredAt,
				Payload:          payloadObj,
			}

			// Insert into inbox table of consumer DB
			inboxID, isNew, err := sharedevent.ConsumeInbox(ctx, targetDB, consumerName, env)
			if err != nil {
				logger.Error().Err(err).Msgf("Failed to write to inbox of %s", consumerName)
				continue
			}

			if isNew {
				// Execute side-effect reactions
				err = executeInboxSideEffect(ctx, consumerName, targetDB, env, logger)
				if err != nil {
					logger.Error().Err(err).Msgf("Side effect execution failed for %s", consumerName)
					_ = sharedevent.SendToDLQ(ctx, targetDB, false, inboxID, err.Error())
				} else {
					_ = sharedevent.MarkProcessed(ctx, targetDB, inboxID)
				}
			} else {
				logger.Info().Msgf("Duplicate event %s ignored for %s", event.EventKey, consumerName)
			}
		}

		// 3. Mark outbox as published
		_ = sharedevent.MarkPublished(ctx, db, event.ID)
		logger.Info().Msgf("[%s] Outbox event marked as published: %s", sourceName, event.ID)
	}
}

// executeInboxSideEffect processes local database adjustments for each consumer module upon event arrival.
func executeInboxSideEffect(ctx context.Context, consumerName string, db *sql.DB, env sharedevent.EventEnvelope, logger zerolog.Logger) error {
	logger.Info().Msgf("Running inbox side effect for consumer %s on event %s", consumerName, env.EventName)
	switch consumerName {
	case "portal":
		// Auto write notifications to portal_db
		title := "Notifikasi Baru"
		message := fmt.Sprintf("Event %s telah diterima.", env.EventName)

		payloadMap, ok := env.Payload.(map[string]interface{})
		userID := ""

		if ok {
			switch env.EventName {
			case "finance.payment_paid":
				title = "Pembayaran Berhasil"
				message = fmt.Sprintf("Tagihan Anda sebesar Rp%v telah sukses dibayar.", payloadMap["amount"])
				if u, exists := payloadMap["payer_ref_id"]; exists {
					userID, _ = u.(string)
				}
			case "academic.student_created":
				title = "NIM Berhasil Dibuat"
				if nimVal, exists := payloadMap["nim"]; exists {
					nim, _ := nimVal.(string)
					message = fmt.Sprintf("Selamat! NIM Anda telah terdaftar: %s.", nim)
				}
				if u, exists := payloadMap["person_id"]; exists {
					userID, _ = u.(string)
				}
			case "academic.krs_approved":
				title = "KRS Disetujui"
				message = "Rencana Studi Semester (KRS) Anda telah disetujui oleh Dosen PA."
				if u, exists := payloadMap["student_id"]; exists {
					userID, _ = u.(string)
				}
			}
		}

		if userID != "" && userID != "person-generated-id" {
			_, err := db.ExecContext(ctx, `
				INSERT INTO notifications (user_id, title, message, module_source, sent_at)
				VALUES ($1, $2, $3, $4, NOW())
			`, userID, title, message, env.PublisherService)
			return err
		}

	case "lms":
		// Auto sync KRS approvals to LMS enrollments
		if env.EventName == "academic.krs_approved" {
			payloadMap, ok := env.Payload.(map[string]interface{})
			if ok {
				studentID, _ := payloadMap["student_id"].(string)
				classesVal, exists := payloadMap["classes"]
				if exists {
					if classSlice, ok := classesVal.([]interface{}); ok {
						for _, classVal := range classSlice {
							classID, _ := classVal.(string)
							// Get LMS class id
							var lmsClassID string
							err := db.QueryRowContext(ctx, `SELECT id FROM classes WHERE academic_class_id = $1`, classID).Scan(&lmsClassID)
							if err == sql.ErrNoRows {
								// Class doesn't exist in LMS yet, sync it
								err = db.QueryRowContext(ctx, `
									INSERT INTO classes (academic_class_id, status, synced_at)
									VALUES ($1, 'active', NOW()) RETURNING id
								`, classID).Scan(&lmsClassID)
							}
							if err != nil {
								return err
							}

							// Upsert enrollment
							_, err = db.ExecContext(ctx, `
								INSERT INTO enrollments (lms_class_id, student_id, enrollment_status, enrolled_at)
								VALUES ($1, $2, 'active', NOW())
								ON CONFLICT (lms_class_id, student_id) DO UPDATE SET enrollment_status = 'active'
							`, lmsClassID, studentID)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return nil
}
