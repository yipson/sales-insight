package sync

import (
	"time"

	"github.com/google/uuid"
)

// SyncEntity represents the type of entity being synchronized.
type SyncEntity string

const (
	SyncEntityOrders   SyncEntity = "orders"
	SyncEntityItems    SyncEntity = "items"
	SyncEntityEmployees SyncEntity = "employees"
	SyncEntityPayments SyncEntity = "payments"
	SyncEntityExport   SyncEntity = "export"
)

// Log represents a synchronization execution record.
type Log struct {
	ID               uuid.UUID  `json:"id"`
	RestaurantID     uuid.UUID  `json:"restaurant_id"`
	Entity           SyncEntity `json:"entity"`
	Status           string     `json:"status"`
	RecordsProcessed int32      `json:"records_processed"`
	RecordsInserted  int32      `json:"records_inserted"`
	RecordsUpdated   int32      `json:"records_updated"`
	RecordsSkipped   int32      `json:"records_skipped"`
	RecordsFailed    int32      `json:"records_failed"`
	CursorFrom       *time.Time `json:"cursor_from,omitempty"`
	CursorTo         *time.Time `json:"cursor_to,omitempty"`
	StartedAt        time.Time  `json:"started_at"`
	FinishedAt       *time.Time `json:"finished_at,omitempty"`
	TriggeredBy      string     `json:"triggered_by"`
	Details          []byte     `json:"details,omitempty"`
}

// Error represents a synchronization error record.
type Error struct {
	ID          uuid.UUID   `json:"id"`
	RestaurantID uuid.UUID  `json:"restaurant_id"`
	SyncLogID   *uuid.UUID  `json:"sync_log_id,omitempty"`
	Entity      SyncEntity  `json:"entity"`
	CloverID    string      `json:"clover_id,omitempty"`
	ErrorCode   string      `json:"error_code,omitempty"`
	ErrorMessage string    `json:"error_message"`
	HTTPStatus  int32       `json:"http_status,omitempty"`
	RetryAfter  int32       `json:"retry_after,omitempty"`
	IsResolved  bool        `json:"is_resolved"`
	ResolvedAt  *time.Time  `json:"resolved_at,omitempty"`
	RetryCount  int32       `json:"retry_count"`
	CreatedAt   time.Time   `json:"created_at"`
}
