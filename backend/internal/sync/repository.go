package sync

import (
	"context"

	"github.com/google/uuid"
)

// LogRepository defines the contract for sync log persistence.
type LogRepository interface {
	GetLatestByEntity(ctx context.Context, restaurantID uuid.UUID) ([]Log, error)
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID, limit int32) ([]Log, error)
	CreateLog(ctx context.Context, log *Log) error
}

// ErrorRepository defines the contract for sync error persistence.
type ErrorRepository interface {
	ListUnresolved(ctx context.Context, restaurantID uuid.UUID, limit int32) ([]Error, error)
}
