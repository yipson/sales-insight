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
	cron       *cron.Cron
	engine     *sync.Engine
	merchantRepo merchant.Repository
	logger     *slog.Logger
}

// NewScheduler creates a new scheduler with the given sync engine.
func NewScheduler(engine *sync.Engine, merchantRepo merchant.Repository, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		cron:       cron.New(),
		engine:     engine,
		merchantRepo: merchantRepo,
		logger:     logger,
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
	// For MVP, use zero time as initial cursor
	// In production, this should be loaded from sync_logs
	cursor := time.Time{}
	newCursor, count, err := s.engine.SyncOrders(ctx, merchantID, cursor)
	if err != nil {
		return err
	}
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
	// For MVP, use zero time as initial cursor
	cursor := time.Time{}
	newCursor, count, err := s.engine.SyncPayments(ctx, merchantID, cursor)
	if err != nil {
		return err
	}
	s.logger.Info("synced payments",
		slog.String("merchant_id", merchantID.String()),
		slog.Int("count", count),
		slog.Time("cursor", newCursor),
	)
	return nil
}
