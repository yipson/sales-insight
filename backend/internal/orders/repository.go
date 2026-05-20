package orders

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the contract for order persistence.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	GetByCloverOrderID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) (*Order, error)
	ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Order, error)
	Upsert(ctx context.Context, order *Order) error
	UpsertBatch(ctx context.Context, orders []Order) error
	DeleteByCloverID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) error
}

// OrderItemRepository defines the contract for order item persistence.
type OrderItemRepository interface {
	UpsertBatch(ctx context.Context, items []OrderItem) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]OrderItem, error)
}

// CategorySummaryRepository defines the contract for order category summary persistence.
type CategorySummaryRepository interface {
	UpsertBatch(ctx context.Context, summaries []OrderCategorySummary) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]OrderCategorySummary, error)
}
