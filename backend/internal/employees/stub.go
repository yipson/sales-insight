package employees

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrNotImplemented is returned by stub repositories.
var ErrNotImplemented = errors.New("repository not implemented")

// StubRepository is a placeholder implementation for Phase 4.
type StubRepository struct{}

func NewStubRepository() *StubRepository { return &StubRepository{} }
func (r *StubRepository) GetByID(ctx context.Context, id uuid.UUID) (*Employee, error) { return nil, ErrNotImplemented }
func (r *StubRepository) GetByCloverEmployeeID(ctx context.Context, restaurantID uuid.UUID, cloverEmployeeID string) (*Employee, error) { return nil, ErrNotImplemented }
func (r *StubRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Employee, error) { return nil, ErrNotImplemented }
func (r *StubRepository) Upsert(ctx context.Context, employee *Employee) error { return ErrNotImplemented }
func (r *StubRepository) UpsertBatch(ctx context.Context, employees []Employee) error { return ErrNotImplemented }
func (r *StubRepository) Deactivate(ctx context.Context, id uuid.UUID) error { return ErrNotImplemented }
