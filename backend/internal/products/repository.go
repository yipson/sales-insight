package products

import (
	"context"

	"github.com/google/uuid"
)

// ProductRepository defines the contract for product persistence.
type ProductRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetByCloverItemID(ctx context.Context, restaurantID uuid.UUID, cloverItemID string) (*Product, error)
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Product, error)
	Upsert(ctx context.Context, product *Product) error
	UpsertBatch(ctx context.Context, products []Product) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

// CategoryRepository defines the contract for native category persistence.
type CategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
	GetByCloverCategoryID(ctx context.Context, restaurantID uuid.UUID, cloverCategoryID string) (*Category, error)
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Category, error)
	Upsert(ctx context.Context, category *Category) error
	UpsertBatch(ctx context.Context, categories []Category) error
}

// AnalyticCategoryRepository defines the contract for analytic category persistence.
type AnalyticCategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*AnalyticCategory, error)
	GetBySlug(ctx context.Context, restaurantID uuid.UUID, slug string) (*AnalyticCategory, error)
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]AnalyticCategory, error)
	Create(ctx context.Context, category *AnalyticCategory) error
	Update(ctx context.Context, category *AnalyticCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// CategoryMappingRepository defines the contract for category mapping persistence.
type CategoryMappingRepository interface {
	GetByCategoryID(ctx context.Context, restaurantID uuid.UUID, categoryID uuid.UUID) (*CategoryMapping, error)
	Upsert(ctx context.Context, mapping *CategoryMapping) error
}
