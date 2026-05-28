package sqlc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/sales-insight/backend/internal/sync"
)

// SQLCRepository implements sync.LogRepository and sync.ErrorRepository.
type SQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewSQLCRepository creates a new sqlc-backed sync repository.
func NewSQLCRepository(db *sql.DB) *SQLCRepository {
	return &SQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func logToDomain(l SyncLog) *sync.Log {
	log := &sync.Log{
		ID:               l.ID,
		RestaurantID:     l.RestaurantID,
		Entity:           sync.SyncEntity(l.Entity),
		Status:           l.Status,
		RecordsProcessed: l.RecordsProcessed,
		RecordsInserted:  l.RecordsInserted,
		RecordsUpdated:   l.RecordsUpdated,
		RecordsSkipped:   l.RecordsSkipped,
		RecordsFailed:    l.RecordsFailed,
		StartedAt:        l.StartedAt,
		TriggeredBy:      l.TriggeredBy,
	}
	if l.CursorFrom.Valid {
		t := l.CursorFrom.Time
		log.CursorFrom = &t
	}
	if l.CursorTo.Valid {
		t := l.CursorTo.Time
		log.CursorTo = &t
	}
	if l.FinishedAt.Valid {
		t := l.FinishedAt.Time
		log.FinishedAt = &t
	}
	if len(l.Details.RawMessage) > 0 {
		log.Details = l.Details.RawMessage
	}
	return log
}

func errorToDomain(e SyncError) *sync.Error {
	err := &sync.Error{
		ID:           e.ID,
		RestaurantID: e.RestaurantID,
		Entity:       sync.SyncEntity(e.Entity),
		ErrorMessage: e.ErrorMessage,
		IsResolved:   e.IsResolved,
		RetryCount:   e.RetryCount,
		CreatedAt:    e.CreatedAt,
	}
	if e.SyncLogID.Valid {
		id := e.SyncLogID.UUID
		err.SyncLogID = &id
	}
	if e.CloverID.Valid {
		err.CloverID = e.CloverID.String
	}
	if e.ErrorCode.Valid {
		err.ErrorCode = e.ErrorCode.String
	}
	if e.HttpStatus.Valid {
		err.HTTPStatus = e.HttpStatus.Int32
	}
	if e.RetryAfter.Valid {
		err.RetryAfter = e.RetryAfter.Int32
	}
	if e.ResolvedAt.Valid {
		t := e.ResolvedAt.Time
		err.ResolvedAt = &t
	}
	return err
}

// LogRepository implementation
func (r *SQLCRepository) GetLatestByEntity(ctx context.Context, restaurantID uuid.UUID) ([]sync.Log, error) {
	rows, err := r.queries.GetLatestSyncByEntity(ctx, restaurantID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []sync.Log
	for _, row := range rows {
		out = append(out, *logToDomain(row))
	}
	return out, nil
}

func (r *SQLCRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID, limit int32) ([]sync.Log, error) {
	rows, err := r.queries.ListSyncLogs(ctx, ListSyncLogsParams{
		RestaurantID: restaurantID,
		Limit:        limit,
	})
	if err != nil {
		return nil, err
	}
	var out []sync.Log
	for _, row := range rows {
		out = append(out, *logToDomain(row))
	}
	return out, nil
}

func (r *SQLCRepository) CreateLog(ctx context.Context, log *sync.Log) error {
	params := CreateSyncLogParams{
		RestaurantID:     log.RestaurantID,
		Entity:           SyncEntity(log.Entity),
		Status:           log.Status,
		RecordsProcessed: log.RecordsProcessed,
		RecordsInserted:  log.RecordsInserted,
		RecordsUpdated:   log.RecordsUpdated,
		RecordsSkipped:   log.RecordsSkipped,
		RecordsFailed:    log.RecordsFailed,
		TriggeredBy:      log.TriggeredBy,
	}
	if log.CursorFrom != nil {
		params.CursorFrom = sql.NullTime{Time: *log.CursorFrom, Valid: true}
	}
	if log.CursorTo != nil {
		params.CursorTo = sql.NullTime{Time: *log.CursorTo, Valid: true}
	}
	if len(log.Details) > 0 {
		params.Details = pqtype.NullRawMessage{
			RawMessage: log.Details,
			Valid:      true,
		}
	}
	row, err := r.queries.CreateSyncLog(ctx, params)
	if err != nil {
		return err
	}
	log.ID = row.ID
	log.StartedAt = row.StartedAt
	return nil
}

// ErrorRepository implementation
func (r *SQLCRepository) ListUnresolved(ctx context.Context, restaurantID uuid.UUID, limit int32) ([]sync.Error, error) {
	rows, err := r.queries.ListSyncErrorsUnresolved(ctx, ListSyncErrorsUnresolvedParams{
		RestaurantID: restaurantID,
		Limit:        limit,
	})
	if err != nil {
		return nil, err
	}
	var out []sync.Error
	for _, row := range rows {
		out = append(out, *errorToDomain(row))
	}
	return out, nil
}
