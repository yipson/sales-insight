package employees

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the contract for employee persistence.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Employee, error)
	GetByCloverEmployeeID(ctx context.Context, restaurantID uuid.UUID, cloverEmployeeID string) (*Employee, error)
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Employee, error)
	Upsert(ctx context.Context, employee *Employee) error
	UpsertBatch(ctx context.Context, employees []Employee) error
	Deactivate(ctx context.Context, id uuid.UUID) error
}
