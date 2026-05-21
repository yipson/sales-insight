package products

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrNotImplemented is returned by stub repositories.
var ErrNotImplemented = errors.New("repository not implemented")

// StubProductRepository is a placeholder implementation for Phase 4.
type StubProductRepository struct{}

func NewStubProductRepository() *StubProductRepository { return &StubProductRepository{} }
func (r *StubProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) { return nil, ErrNotImplemented }
func (r *StubProductRepository) GetByCloverItemID(ctx context.Context, restaurantID uuid.UUID, cloverItemID string) (*Product, error) { return nil, ErrNotImplemented }
func (r *StubProductRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Product, error) { return nil, ErrNotImplemented }
func (r *StubProductRepository) Upsert(ctx context.Context, product *Product) error { return ErrNotImplemented }
func (r *StubProductRepository) UpsertBatch(ctx context.Context, products []Product) error { return ErrNotImplemented }
func (r *StubProductRepository) SoftDelete(ctx context.Context, id uuid.UUID) error { return ErrNotImplemented }

// StubCategoryRepository is a placeholder implementation for Phase 4.
type StubCategoryRepository struct{}

func NewStubCategoryRepository() *StubCategoryRepository { return &StubCategoryRepository{} }
func (r *StubCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*Category, error) { return nil, ErrNotImplemented }
func (r *StubCategoryRepository) GetByCloverCategoryID(ctx context.Context, restaurantID uuid.UUID, cloverCategoryID string) (*Category, error) { return nil, ErrNotImplemented }
func (r *StubCategoryRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Category, error) { return nil, ErrNotImplemented }
func (r *StubCategoryRepository) Upsert(ctx context.Context, category *Category) error { return ErrNotImplemented }
func (r *StubCategoryRepository) UpsertBatch(ctx context.Context, categories []Category) error { return ErrNotImplemented }

// StubAnalyticCategoryRepository is a placeholder implementation for Phase 4.
type StubAnalyticCategoryRepository struct{}

func NewStubAnalyticCategoryRepository() *StubAnalyticCategoryRepository { return &StubAnalyticCategoryRepository{} }
func (r *StubAnalyticCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*AnalyticCategory, error) { return nil, ErrNotImplemented }
func (r *StubAnalyticCategoryRepository) GetBySlug(ctx context.Context, restaurantID uuid.UUID, slug string) (*AnalyticCategory, error) { return nil, ErrNotImplemented }
func (r *StubAnalyticCategoryRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]AnalyticCategory, error) { return nil, ErrNotImplemented }
func (r *StubAnalyticCategoryRepository) Create(ctx context.Context, category *AnalyticCategory) error { return ErrNotImplemented }
func (r *StubAnalyticCategoryRepository) Update(ctx context.Context, category *AnalyticCategory) error { return ErrNotImplemented }
func (r *StubAnalyticCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error { return ErrNotImplemented }

// StubCategoryMappingRepository is a placeholder implementation for Phase 4.
type StubCategoryMappingRepository struct{}

func NewStubCategoryMappingRepository() *StubCategoryMappingRepository { return &StubCategoryMappingRepository{} }
func (r *StubCategoryMappingRepository) GetByCategoryID(ctx context.Context, restaurantID uuid.UUID, categoryID uuid.UUID) (*CategoryMapping, error) { return nil, ErrNotImplemented }
func (r *StubCategoryMappingRepository) Upsert(ctx context.Context, mapping *CategoryMapping) error { return ErrNotImplemented }
