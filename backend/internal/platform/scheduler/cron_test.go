package scheduler

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/sync"
)

// mockSyncLogRepository is a test double for sync.LogRepository.
type mockSyncLogRepository struct {
	logs []sync.Log
}

func (m *mockSyncLogRepository) GetLatestByEntity(_ context.Context, _ uuid.UUID) ([]sync.Log, error) {
	return m.logs, nil
}

func (m *mockSyncLogRepository) ListByRestaurant(_ context.Context, _ uuid.UUID, _ int32) ([]sync.Log, error) {
	return m.logs, nil
}

func (m *mockSyncLogRepository) CreateLog(_ context.Context, _ *sync.Log) error {
	return nil
}

func newTestScheduler(logs []sync.Log) *Scheduler {
	return &Scheduler{
		syncLogRepo: &mockSyncLogRepository{logs: logs},
		logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func TestFirstDayOfPreviousMonth(t *testing.T) {
	cases := []struct {
		input    time.Time
		expected time.Time
	}{
		{
			input:    time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, c := range cases {
		got := firstDayOfPreviousMonth(c.input)
		if !got.Equal(c.expected) {
			t.Errorf("firstDayOfPreviousMonth(%v) = %v, want %v", c.input, got, c.expected)
		}
	}
}

func TestGetCursorForEntity_WithPreviousLog(t *testing.T) {
	prevCursor := time.Date(2026, 5, 20, 14, 30, 0, 0, time.UTC)
	logs := []sync.Log{
		{
			Entity:   sync.SyncEntityOrders,
			Status:   "success",
			CursorTo: &prevCursor,
		},
	}

	s := newTestScheduler(logs)
	ctx := context.Background()
	merchantID := uuid.New()

	cursor, err := s.getCursorForEntity(ctx, merchantID, sync.SyncEntityOrders)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cursor.Equal(prevCursor) {
		t.Errorf("cursor = %v, want %v", cursor, prevCursor)
	}
}

func TestGetCursorForEntity_WithoutPreviousLog(t *testing.T) {
	s := newTestScheduler(nil)
	ctx := context.Background()
	merchantID := uuid.New()

	cursor, err := s.getCursorForEntity(ctx, merchantID, sync.SyncEntityOrders)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := firstDayOfPreviousMonth(time.Now().UTC())
	if !cursor.Equal(expected) {
		t.Errorf("cursor = %v, want %v (first day of previous month)", cursor, expected)
	}
}

func TestGetCursorForEntity_SkipsFailedLogs(t *testing.T) {
	failedCursor := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	logs := []sync.Log{
		{
			Entity:   sync.SyncEntityOrders,
			Status:   "failed",
			CursorTo: &failedCursor,
		},
	}

	s := newTestScheduler(logs)
	ctx := context.Background()
	merchantID := uuid.New()

	cursor, err := s.getCursorForEntity(ctx, merchantID, sync.SyncEntityOrders)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should fall back to first day of previous month, not use failed log
	expected := firstDayOfPreviousMonth(time.Now().UTC())
	if !cursor.Equal(expected) {
		t.Errorf("cursor = %v, want %v (should skip failed log)", cursor, expected)
	}
}
