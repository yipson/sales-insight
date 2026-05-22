package products

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mockProductRepository is a test double for ProductRepository.
type mockProductRepository struct {
	products map[uuid.UUID]*Product
	byClover map[string]*Product
	byRestaurant map[uuid.UUID][]Product
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[uuid.UUID]*Product),
		byClover: make(map[string]*Product),
		byRestaurant: make(map[uuid.UUID][]Product),
	}
}

func (m *mockProductRepository) GetByID(_ context.Context, id uuid.UUID) (*Product, error) {
	return m.products[id], nil
}

func (m *mockProductRepository) GetByCloverItemID(_ context.Context, restaurantID uuid.UUID, cloverItemID string) (*Product, error) {
	return m.byClover[cloverItemID], nil
}

func (m *mockProductRepository) ListByRestaurant(_ context.Context, restaurantID uuid.UUID) ([]Product, error) {
	return m.byRestaurant[restaurantID], nil
}

func (m *mockProductRepository) Upsert(_ context.Context, p *Product) error {
	m.products[p.ID] = p
	m.byClover[p.CloverItemID] = p
	m.byRestaurant[p.RestaurantID] = append(m.byRestaurant[p.RestaurantID], *p)
	return nil
}

func (m *mockProductRepository) UpsertBatch(_ context.Context, products []Product) error {
	for i := range products {
		if err := m.Upsert(nil, &products[i]); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockProductRepository) SoftDelete(_ context.Context, id uuid.UUID) error {
	if p, ok := m.products[id]; ok {
		p.IsDeleted = true
	}
	return nil
}

// mockCategoryRepository is a test double for CategoryRepository.
type mockCategoryRepository struct {
	categories map[uuid.UUID]*Category
	byClover   map[string]*Category
	byRestaurant map[uuid.UUID][]Category
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[uuid.UUID]*Category),
		byClover:   make(map[string]*Category),
		byRestaurant: make(map[uuid.UUID][]Category),
	}
}

func (m *mockCategoryRepository) GetByID(_ context.Context, id uuid.UUID) (*Category, error) {
	return m.categories[id], nil
}

func (m *mockCategoryRepository) GetByCloverCategoryID(_ context.Context, restaurantID uuid.UUID, cloverCategoryID string) (*Category, error) {
	return m.byClover[cloverCategoryID], nil
}

func (m *mockCategoryRepository) ListByRestaurant(_ context.Context, restaurantID uuid.UUID) ([]Category, error) {
	return m.byRestaurant[restaurantID], nil
}

func (m *mockCategoryRepository) Upsert(_ context.Context, c *Category) error {
	m.categories[c.ID] = c
	m.byClover[c.CloverCategoryID] = c
	m.byRestaurant[c.RestaurantID] = append(m.byRestaurant[c.RestaurantID], *c)
	return nil
}

func (m *mockCategoryRepository) UpsertBatch(_ context.Context, categories []Category) error {
	for i := range categories {
		if err := m.Upsert(nil, &categories[i]); err != nil {
			return err
		}
	}
	return nil
}

// mockAnalyticCategoryRepository is a test double for AnalyticCategoryRepository.
type mockAnalyticCategoryRepository struct {
	categories map[uuid.UUID]*AnalyticCategory
	bySlug     map[string]*AnalyticCategory
	byRestaurant map[uuid.UUID][]AnalyticCategory
}

func newMockAnalyticCategoryRepository() *mockAnalyticCategoryRepository {
	return &mockAnalyticCategoryRepository{
		categories: make(map[uuid.UUID]*AnalyticCategory),
		bySlug:     make(map[string]*AnalyticCategory),
		byRestaurant: make(map[uuid.UUID][]AnalyticCategory),
	}
}

func (m *mockAnalyticCategoryRepository) GetByID(_ context.Context, id uuid.UUID) (*AnalyticCategory, error) {
	return m.categories[id], nil
}

func (m *mockAnalyticCategoryRepository) GetBySlug(_ context.Context, restaurantID uuid.UUID, slug string) (*AnalyticCategory, error) {
	return m.bySlug[slug], nil
}

func (m *mockAnalyticCategoryRepository) ListByRestaurant(_ context.Context, restaurantID uuid.UUID) ([]AnalyticCategory, error) {
	return m.byRestaurant[restaurantID], nil
}

func (m *mockAnalyticCategoryRepository) Create(_ context.Context, ac *AnalyticCategory) error {
	m.categories[ac.ID] = ac
	m.bySlug[ac.Slug] = ac
	m.byRestaurant[ac.RestaurantID] = append(m.byRestaurant[ac.RestaurantID], *ac)
	return nil
}

func (m *mockAnalyticCategoryRepository) Update(_ context.Context, ac *AnalyticCategory) error {
	m.categories[ac.ID] = ac
	// Update slug index
	m.bySlug[ac.Slug] = ac
	return nil
}

func (m *mockAnalyticCategoryRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.categories, id)
	return nil
}

// mockCategoryMappingRepository is a test double for CategoryMappingRepository.
type mockCategoryMappingRepository struct {
	mappings map[uuid.UUID]*CategoryMapping // key: categoryID
}

func newMockCategoryMappingRepository() *mockCategoryMappingRepository {
	return &mockCategoryMappingRepository{
		mappings: make(map[uuid.UUID]*CategoryMapping),
	}
}

func (m *mockCategoryMappingRepository) GetByCategoryID(_ context.Context, restaurantID uuid.UUID, categoryID uuid.UUID) (*CategoryMapping, error) {
	return m.mappings[categoryID], nil
}

func (m *mockCategoryMappingRepository) Upsert(_ context.Context, cm *CategoryMapping) error {
	m.mappings[cm.CategoryID] = cm
	return nil
}

func newTestService() (*Service, *mockProductRepository, *mockCategoryRepository, *mockAnalyticCategoryRepository, *mockCategoryMappingRepository) {
	prodRepo := newMockProductRepository()
	catRepo := newMockCategoryRepository()
	analyticRepo := newMockAnalyticCategoryRepository()
	mappingRepo := newMockCategoryMappingRepository()
	svc := NewService(prodRepo, catRepo, analyticRepo, mappingRepo)
	return svc, prodRepo, catRepo, analyticRepo, mappingRepo
}

func TestService_GetProductByID(t *testing.T) {
	svc, prodRepo, _, _, _ := newTestService()
	ctx := context.Background()

	p := &Product{ID: uuid.New(), Name: "Taco"}
	prodRepo.products[p.ID] = p

	got, err := svc.GetProductByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Name != "Taco" {
		t.Errorf("expected Taco, got %v", got)
	}
}

func TestService_UpsertProductBatch(t *testing.T) {
	svc, prodRepo, _, _, _ := newTestService()
	ctx := context.Background()
	restID := uuid.New()

	products := []Product{
		{ID: uuid.New(), Name: "Taco", RestaurantID: restID, CloverItemID: "clv_1"},
		{ID: uuid.New(), Name: "Burrito", RestaurantID: restID, CloverItemID: "clv_2"},
	}

	if err := svc.UpsertProductBatch(ctx, products); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(prodRepo.products) != 2 {
		t.Errorf("expected 2 products, got %d", len(prodRepo.products))
	}
}

func TestService_CreateAnalyticCategory(t *testing.T) {
	svc, _, _, analyticRepo, _ := newTestService()
	ctx := context.Background()
	restID := uuid.New()

	ac, err := svc.CreateAnalyticCategory(ctx, restID, "Tacos", "tacos", 1, "#FF0000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ac.Name != "Tacos" {
		t.Errorf("name = %q, want %q", ac.Name, "Tacos")
	}
	if !ac.IsActive {
		t.Error("expected active")
	}
	if len(analyticRepo.categories) != 1 {
		t.Errorf("expected 1 category, got %d", len(analyticRepo.categories))
	}
}

func TestService_CreateAnalyticCategory_DuplicateSlug(t *testing.T) {
	svc, _, _, analyticRepo, _ := newTestService()
	ctx := context.Background()
	restID := uuid.New()

	_ = analyticRepo.Create(nil, &AnalyticCategory{
		ID: uuid.New(), RestaurantID: restID, Name: "Tacos", Slug: "tacos",
	})

	_, err := svc.CreateAnalyticCategory(ctx, restID, "Tacos 2", "tacos", 2, "#00FF00")
	if err == nil {
		t.Error("expected error for duplicate slug")
	}
}

func TestService_UpdateAnalyticCategory(t *testing.T) {
	svc, _, _, analyticRepo, _ := newTestService()
	ctx := context.Background()

	existing := &AnalyticCategory{
		ID: uuid.New(), Name: "Old", Slug: "old", DisplayOrder: 1,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	_ = analyticRepo.Create(nil, existing)

	updated, err := svc.UpdateAnalyticCategory(ctx, existing.ID, "New", "new", 2, "#0000FF", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "New" {
		t.Errorf("name = %q, want %q", updated.Name, "New")
	}
}

func TestService_SoftDeleteProduct(t *testing.T) {
	svc, prodRepo, _, _, _ := newTestService()
	ctx := context.Background()

	p := &Product{ID: uuid.New(), Name: "Taco", IsDeleted: false}
	prodRepo.products[p.ID] = p

	if err := svc.SoftDeleteProduct(ctx, p.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !prodRepo.products[p.ID].IsDeleted {
		t.Error("expected product to be soft deleted")
	}
}

func TestService_UpsertCategoryMapping(t *testing.T) {
	svc, _, _, _, mappingRepo := newTestService()
	ctx := context.Background()

	mapping := &CategoryMapping{
		ID:                 uuid.New(),
		RestaurantID:       uuid.New(),
		CategoryID:         uuid.New(),
		AnalyticCategoryID: uuid.New(),
		MappedBy:           "admin",
	}

	if err := svc.UpsertCategoryMapping(ctx, mapping); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mappingRepo.mappings) != 1 {
		t.Errorf("expected 1 mapping, got %d", len(mappingRepo.mappings))
	}
}
