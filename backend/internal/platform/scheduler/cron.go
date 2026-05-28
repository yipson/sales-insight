package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	"github.com/sales-insight/backend/internal/merchant"
	"github.com/sales-insight/backend/internal/sync"
)

// Scheduler wraps robfig/cron to run periodic sync tasks.
type Scheduler struct {
	cron         *cron.Cron
	engine       *sync.Engine
	merchantRepo   merchant.Repository
	syncLogRepo   sync.LogRepository
	logger       *slog.Logger
}

// NewScheduler creates a new scheduler with the given sync engine.
func NewScheduler(engine *sync.Engine, merchantRepo merchant.Repository, syncLogRepo sync.LogRepository, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		cron:         cron.New(),
		engine:       engine,
		merchantRepo: merchantRepo,
		syncLogRepo:  syncLogRepo,
		logger:       logger,
	}
}

// Start registers all sync jobs and starts the cron scheduler.
func (s *Scheduler) Start(ctx context.Context) {
	// Sync orders every 5 minutes
	_, _ = s.cron.AddFunc("*/5 * * * *", func() {
		s.runForAllMerchants(ctx, "orders", s.syncOrders)
	})

	// Sync items every 30 minutes
	_, _ = s.cron.AddFunc("*/30 * * * *", func() {
		s.runForAllMerchants(ctx, "items", s.syncItems)
	})

	// Sync employees every 1 hour
	_, _ = s.cron.AddFunc("0 * * * *", func() {
		s.runForAllMerchants(ctx, "employees", s.syncEmployees)
	})

	// Sync payments every 15 minutes
	_, _ = s.cron.AddFunc("*/15 * * * *", func() {
		s.runForAllMerchants(ctx, "payments", s.syncPayments)
	})

	s.cron.Start()
	s.logger.Info("scheduler started")
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("scheduler stopped")
}

// firstDayOfPreviousMonth returns the first day of the month before the given time.
func firstDayOfPreviousMonth(t time.Time) time.Time {
	firstOfCurrent := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return firstOfCurrent.AddDate(0, -1, 0)
}

// getCursorForEntity returns the cursor to use for the next sync of the given entity.
// It looks up the latest successful sync log. If none exists, it falls back to
// the first day of the previous month (ensuring the first sync is bounded).
func (s *Scheduler) getCursorForEntity(ctx context.Context, merchantID uuid.UUID, entity sync.SyncEntity) (time.Time, error) {
	logs, err := s.syncLogRepo.GetLatestByEntity(ctx, merchantID)
	if err != nil {
		return time.Time{}, err
	}

	// Look for the most recent log matching this entity
	for _, log := range logs {
		if log.Entity == entity && log.Status == "success" && log.CursorTo != nil {
			return *log.CursorTo, nil
		}
	}

	// No previous successful sync — fallback to first day of previous month
	return firstDayOfPreviousMonth(time.Now().UTC()), nil
}

// writeSyncLog persists a sync execution record.
func (s *Scheduler) writeSyncLog(ctx context.Context, merchantID uuid.UUID, entity sync.SyncEntity, status string, count int, cursorFrom, cursorTo time.Time, errMsg string) {
	log := &sync.Log{
		RestaurantID:     merchantID,
		Entity:             entity,
		Status:             status,
		RecordsProcessed:   int32(count),
		TriggeredBy:        "scheduler",
	}
	if !cursorFrom.IsZero() {
		log.CursorFrom = &cursorFrom
	}
	if !cursorTo.IsZero() {
		log.CursorTo = &cursorTo
	}
	if errMsg != "" {
		log.Details = []byte(`{"error": "` + errMsg + `"}`)
	}
	if err := s.syncLogRepo.CreateLog(ctx, log); err != nil {
		s.logger.Error("failed to write sync log",
			slog.String("error", err.Error()),
			slog.String("entity", string(entity)),
			slog.String("merchant_id", merchantID.String()),
		)
	}
}

// runForAllMerchants executes a sync function for every connected merchant.
func (s *Scheduler) runForAllMerchants(ctx context.Context, entity string, fn func(context.Context, uuid.UUID) error) {
	merchants, err := s.merchantRepo.List(ctx)
	if err != nil {
		s.logger.Error("failed to list merchants", slog.String("error", err.Error()))
		return
	}

	for _, m := range merchants {
		if !m.IsConnected {
			continue
		}

		if err := fn(ctx, m.ID); err != nil {
			s.logger.Error("sync failed",
				slog.String("entity", entity),
				slog.String("merchant_id", m.ID.String()),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (s *Scheduler) syncOrders(ctx context.Context, merchantID uuid.UUID) error {
	entity := sync.SyncEntityOrders
	cursor, err := s.getCursorForEntity(ctx, merchantID, entity)
	if err != nil {
		return err
	}

	newCursor, count, err := s.engine.SyncOrders(ctx, merchantID, cursor)
	if err != nil {
		s.writeSyncLog(ctx, merchantID, entity, "failed", 0, cursor, time.Time{}, err.Error())
		return err
	}

	s.writeSyncLog(ctx, merchantID, entity, "success", count, cursor, newCursor, "")
	s.logger.Info("synced orders",
		slog.String("merchant_id", merchantID.String()),
		slog.Int("count", count),
		slog.Time("cursor", newCursor),
	)
	return nil
}

func (s *Scheduler) syncItems(ctx context.Context, merchantID uuid.UUID) error {
	count, catCount, err := s.engine.SyncItems(ctx, merchantID)
	if err != nil {
		return err
	}
	s.logger.Info("synced items",
		slog.String("merchant_id", merchantID.String()),
		slog.Int("products", count),
		slog.Int("categories", catCount),
	)
	return nil
}

func (s *Scheduler) syncEmployees(ctx context.Context, merchantID uuid.UUID) error {
	count, err := s.engine.SyncEmployees(ctx, merchantID)
	if err != nil {
		return err
	}
	s.logger.Info("synced employees",
		slog.String("merchant_id", merchantID.String()),
		slog.Int("count", count),
	)
	return nil
}

func (s *Scheduler) syncPayments(ctx context.Context, merchantID uuid.UUID) error {
	entity := sync.SyncEntityPayments
	cursor, err := s.getCursorForEntity(ctx, merchantID, entity)
	if err != nil {
		return err
	}

	newCursor, count, err := s.engine.SyncPayments(ctx, merchantID, cursor)
	if err != nil {
		s.writeSyncLog(ctx, merchantID, entity, "failed", 0, cursor, time.Time{}, err.Error())
		return err
	}

	s.writeSyncLog(ctx, merchantID, entity, "success", count, cursor, newCursor, "")
	s.logger.Info("synced payments",
		slog.String("merchant_id", merchantID.String()),
		slog.Int("count", count),
		slog.Time("cursor", newCursor),
	)
	return nil
}
