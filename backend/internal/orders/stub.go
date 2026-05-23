package orders

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
func (r *StubRepository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) { return nil, ErrNotImplemented }
func (r *StubRepository) GetByCloverOrderID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) (*Order, error) { return nil, ErrNotImplemented }
func (r *StubRepository) ListByRestaurantAndDate(ctx context.Context, restaurantID uuid.UUID, from, to time.Time) ([]Order, error) { return nil, ErrNotImplemented }
func (r *StubRepository) ListByEmployee(ctx context.Context, employeeID uuid.UUID, from, to time.Time) ([]Order, error) { return nil, ErrNotImplemented }
func (r *StubRepository) Upsert(ctx context.Context, order *Order) error { return ErrNotImplemented }
func (r *StubRepository) UpsertBatch(ctx context.Context, orders []Order) error { return ErrNotImplemented }
func (r *StubRepository) DeleteByCloverID(ctx context.Context, restaurantID uuid.UUID, cloverOrderID string) error { return ErrNotImplemented }

// StubOrderItemRepository is a placeholder implementation for Phase 4.
type StubOrderItemRepository struct{}

func NewStubOrderItemRepository() *StubOrderItemRepository { return &StubOrderItemRepository{} }
func (r *StubOrderItemRepository) UpsertBatch(ctx context.Context, items []OrderItem) error { return ErrNotImplemented }
func (r *StubOrderItemRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]OrderItem, error) { return nil, ErrNotImplemented }

// StubCategorySummaryRepository is a placeholder implementation for Phase 4.
type StubCategorySummaryRepository struct{}

func NewStubCategorySummaryRepository() *StubCategorySummaryRepository { return &StubCategorySummaryRepository{} }
func (r *StubCategorySummaryRepository) UpsertBatch(ctx context.Context, summaries []OrderCategorySummary) error { return ErrNotImplemented }
func (r *StubCategorySummaryRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]OrderCategorySummary, error) { return nil, ErrNotImplemented }
