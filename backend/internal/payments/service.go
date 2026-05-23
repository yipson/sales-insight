package payments

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service holds business logic for payments.
type Service struct {
	repo Repository
}

// NewService creates a new payment service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetByID returns a payment by UUID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Payment, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCloverPaymentID returns a payment by Clover payment ID.
func (s *Service) GetByCloverPaymentID(ctx context.Context, restaurantID uuid.UUID, cloverPaymentID string) (*Payment, error) {
	return s.repo.GetByCloverPaymentID(ctx, restaurantID, cloverPaymentID)
}

// ListByRestaurantAndDate returns payments for a restaurant in a date range.
func (s *Service) ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Payment, error) {
	return s.repo.ListByRestaurantAndDate(ctx, restaurantID, from, to)
}

// Upsert creates or updates a single payment (used by sync engine).
func (s *Service) Upsert(ctx context.Context, payment *Payment) error {
	return s.repo.Upsert(ctx, payment)
}

// UpsertBatch creates or updates multiple payments atomically (used by sync engine).
func (s *Service) UpsertBatch(ctx context.Context, payments []Payment) error {
	return s.repo.UpsertBatch(ctx, payments)
}
