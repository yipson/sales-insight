package products

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service holds business logic for products, categories and mappings.
type Service struct {
	productRepo           ProductRepository
	categoryRepo          CategoryRepository
	analyticCategoryRepo  AnalyticCategoryRepository
	mappingRepo           CategoryMappingRepository
}

// NewService creates a new product service.
func NewService(
	productRepo ProductRepository,
	categoryRepo CategoryRepository,
	analyticCategoryRepo AnalyticCategoryRepository,
	mappingRepo CategoryMappingRepository,
) *Service {
	return &Service{
		productRepo:          productRepo,
		categoryRepo:         categoryRepo,
		analyticCategoryRepo: analyticCategoryRepo,
		mappingRepo:          mappingRepo,
	}
}

// ============================================================
// Products
// ============================================================

func (s *Service) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *Service) ListProductsByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Product, error) {
	return s.productRepo.ListByRestaurant(ctx, restaurantID)
}

func (s *Service) UpsertProduct(ctx context.Context, product *Product) error {
	return s.productRepo.Upsert(ctx, product)
}

func (s *Service) UpsertProductBatch(ctx context.Context, products []Product) error {
	return s.productRepo.UpsertBatch(ctx, products)
}

func (s *Service) SoftDeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.productRepo.SoftDelete(ctx, id)
}

// ============================================================
// Native Categories
// ============================================================

func (s *Service) GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *Service) ListCategoriesByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]Category, error) {
	return s.categoryRepo.ListByRestaurant(ctx, restaurantID)
}

func (s *Service) UpsertCategory(ctx context.Context, category *Category) error {
	return s.categoryRepo.Upsert(ctx, category)
}

func (s *Service) UpsertCategoryBatch(ctx context.Context, categories []Category) error {
	return s.categoryRepo.UpsertBatch(ctx, categories)
}

// ============================================================
// Analytic Categories
// ============================================================

func (s *Service) GetAnalyticCategoryByID(ctx context.Context, id uuid.UUID) (*AnalyticCategory, error) {
	return s.analyticCategoryRepo.GetByID(ctx, id)
}

func (s *Service) GetAnalyticCategoryBySlug(ctx context.Context, restaurantID uuid.UUID, slug string) (*AnalyticCategory, error) {
	return s.analyticCategoryRepo.GetBySlug(ctx, restaurantID, slug)
}

func (s *Service) ListAnalyticCategoriesByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]AnalyticCategory, error) {
	return s.analyticCategoryRepo.ListByRestaurant(ctx, restaurantID)
}

func (s *Service) CreateAnalyticCategory(ctx context.Context, restaurantID uuid.UUID, name, slug string, displayOrder int32, color string) (*AnalyticCategory, error) {
	if name == "" || slug == "" {
		return nil, fmt.Errorf("name and slug are required")
	}
	existing, err := s.analyticCategoryRepo.GetBySlug(ctx, restaurantID, slug)
	if err != nil {
		return nil, fmt.Errorf("lookup analytic category: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("analytic category with slug %q already exists", slug)
	}

	ac := &AnalyticCategory{
		ID:           uuid.New(),
		RestaurantID: restaurantID,
		Name:         name,
		Slug:         slug,
		DisplayOrder: displayOrder,
		Color:        color,
		IsActive:     true,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := s.analyticCategoryRepo.Create(ctx, ac); err != nil {
		return nil, fmt.Errorf("create analytic category: %w", err)
	}
	return ac, nil
}

func (s *Service) UpdateAnalyticCategory(ctx context.Context, id uuid.UUID, name, slug string, displayOrder int32, color string, isActive bool) (*AnalyticCategory, error) {
	existing, err := s.analyticCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lookup analytic category: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("analytic category not found")
	}

	existing.Name = name
	existing.Slug = slug
	existing.DisplayOrder = displayOrder
	existing.Color = color
	existing.IsActive = isActive
	existing.UpdatedAt = time.Now().UTC()

	if err := s.analyticCategoryRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update analytic category: %w", err)
	}
	return existing, nil
}

func (s *Service) DeleteAnalyticCategory(ctx context.Context, id uuid.UUID) error {
	existing, err := s.analyticCategoryRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("lookup analytic category: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("analytic category not found")
	}
	return s.analyticCategoryRepo.Delete(ctx, id)
}

// ============================================================
// Category Mappings
// ============================================================

func (s *Service) GetCategoryMapping(ctx context.Context, restaurantID, categoryID uuid.UUID) (*CategoryMapping, error) {
	return s.mappingRepo.GetByCategoryID(ctx, restaurantID, categoryID)
}

func (s *Service) UpsertCategoryMapping(ctx context.Context, mapping *CategoryMapping) error {
	return s.mappingRepo.Upsert(ctx, mapping)
}
