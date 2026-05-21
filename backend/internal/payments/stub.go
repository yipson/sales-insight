package payments

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrNotImplemented is returned by stub repositories.
var ErrNotImplemented = errors.New("repository not implemented")

// StubRepository is a placeholder implementation for Phase 4.
type StubRepository struct{}

func NewStubRepository() *StubRepository { return &StubRepository{} }
func (r *StubRepository) GetByID(ctx context.Context, id uuid.UUID) (*Payment, error) { return nil, ErrNotImplemented }
func (r *StubRepository) GetByCloverPaymentID(ctx context.Context, restaurantID uuid.UUID, cloverPaymentID string) (*Payment, error) { return nil, ErrNotImplemented }
func (r *StubRepository) ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Payment, error) { return nil, ErrNotImplemented }
func (r *StubRepository) Upsert(ctx context.Context, payment *Payment) error { return ErrNotImplemented }
func (r *StubRepository) UpsertBatch(ctx context.Context, payments []Payment) error { return ErrNotImplemented }
