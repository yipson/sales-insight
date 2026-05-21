package payments

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the contract for payment persistence.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	GetByCloverPaymentID(ctx context.Context, restaurantID uuid.UUID, cloverPaymentID string) (*Payment, error)
	ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Payment, error)
	Upsert(ctx context.Context, payment *Payment) error
	UpsertBatch(ctx context.Context, payments []Payment) error
}
