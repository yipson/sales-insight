package sqlc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sales-insight/backend/internal/products"
)

// ============================================================
// Product helpers
// ============================================================

func productToDomain(p Product) *products.Product {
	prod := &products.Product{
		ID:           p.ID,
		RestaurantID: p.RestaurantID,
		CloverItemID: p.CloverItemID,
		Name:         p.Name,
		Price:        p.Price,
		IsAvailable:  p.IsAvailable,
		IsDeleted:    p.IsDeleted,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		SyncedAt:     p.SyncedAt,
	}
	if p.Description.Valid {
		prod.Description = p.Description.String
	}
	if p.Cost.Valid {
		c := p.Cost.Int64
		prod.Cost = &c
	}
	if p.Sku.Valid {
		prod.SKU = p.Sku.String
	}
	if p.CategoryID.Valid {
		cid := p.CategoryID.UUID
		prod.CategoryID = &cid
	}
	return prod
}

func productFromDomain(p *products.Product) UpsertProductParams {
	params := UpsertProductParams{
		ID:           p.ID,
		RestaurantID: p.RestaurantID,
		CloverItemID: p.CloverItemID,
		Name:         p.Name,
		Price:        p.Price,
		IsAvailable:  p.IsAvailable,
		IsDeleted:    p.IsDeleted,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		SyncedAt:     p.SyncedAt,
	}
	params.Description = sql.NullString{String: p.Description, Valid: p.Description != ""}
	if p.Cost != nil {
		params.Cost = sql.NullInt64{Int64: *p.Cost, Valid: true}
	}
	params.Sku = sql.NullString{String: p.SKU, Valid: p.SKU != ""}
	if p.CategoryID != nil {
		params.CategoryID = uuid.NullUUID{UUID: *p.CategoryID, Valid: true}
	}
	return params
}

// ============================================================
// ProductSQLCRepository implements products.ProductRepository.
// ============================================================

type ProductSQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewProductSQLCRepository creates a new sqlc-backed product repository.
func NewProductSQLCRepository(db *sql.DB) *ProductSQLCRepository {
	return &ProductSQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func (r *ProductSQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*products.Product, error) {
	row, err := r.queries.GetProductByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return productToDomain(row), nil
}

func (r *ProductSQLCRepository) GetByCloverItemID(ctx context.Context, restaurantID uuid.UUID, cloverItemID string) (*products.Product, error) {
	row, err := r.queries.GetProductByCloverItemID(ctx, GetProductByCloverItemIDParams{
		RestaurantID: restaurantID,
		CloverItemID: cloverItemID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return productToDomain(row), nil
}

func (r *ProductSQLCRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]products.Product, error) {
	rows, err := r.queries.ListProductsByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	var out []products.Product
	for _, row := range rows {
		out = append(out, *productToDomain(row))
	}
	return out, nil
}

func (r *ProductSQLCRepository) Upsert(ctx context.Context, product *products.Product) error {
	_, err := r.queries.UpsertProduct(ctx, productFromDomain(product))
	return err
}

func (r *ProductSQLCRepository) UpsertBatch(ctx context.Context, products []products.Product) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for i := range products {
		if _, err := qtx.UpsertProduct(ctx, productFromDomain(&products[i])); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *ProductSQLCRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeleteProduct(ctx, id)
}

// ============================================================
// Category helpers
// ============================================================

func categoryToDomain(c Category) *products.Category {
	cat := &products.Category{
		ID:               c.ID,
		RestaurantID:     c.RestaurantID,
		CloverCategoryID: c.CloverCategoryID,
		Name:             c.Name,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		SyncedAt:         c.SyncedAt,
	}
	if c.SortOrder.Valid {
		so := c.SortOrder.Int32
		cat.SortOrder = &so
	}
	return cat
}

func categoryFromDomain(c *products.Category) UpsertCategoryParams {
	params := UpsertCategoryParams{
		ID:               c.ID,
		RestaurantID:     c.RestaurantID,
		CloverCategoryID: c.CloverCategoryID,
		Name:             c.Name,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		SyncedAt:         c.SyncedAt,
	}
	if c.SortOrder != nil {
		params.SortOrder = sql.NullInt32{Int32: *c.SortOrder, Valid: true}
	}
	return params
}

// ============================================================
// CategorySQLCRepository implements products.CategoryRepository.
// ============================================================

type CategorySQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewCategorySQLCRepository creates a new sqlc-backed category repository.
func NewCategorySQLCRepository(db *sql.DB) *CategorySQLCRepository {
	return &CategorySQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func (r *CategorySQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*products.Category, error) {
	row, err := r.queries.GetCategoryByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return categoryToDomain(row), nil
}

func (r *CategorySQLCRepository) GetByCloverCategoryID(ctx context.Context, restaurantID uuid.UUID, cloverCategoryID string) (*products.Category, error) {
	row, err := r.queries.GetCategoryByCloverID(ctx, GetCategoryByCloverIDParams{
		RestaurantID:     restaurantID,
		CloverCategoryID: cloverCategoryID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return categoryToDomain(row), nil
}

func (r *CategorySQLCRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]products.Category, error) {
	rows, err := r.queries.ListCategoriesByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	var out []products.Category
	for _, row := range rows {
		out = append(out, *categoryToDomain(row))
	}
	return out, nil
}

func (r *CategorySQLCRepository) Upsert(ctx context.Context, category *products.Category) error {
	_, err := r.queries.UpsertCategory(ctx, categoryFromDomain(category))
	return err
}

func (r *CategorySQLCRepository) UpsertBatch(ctx context.Context, categories []products.Category) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.queries.WithTx(tx)
	for i := range categories {
		if _, err := qtx.UpsertCategory(ctx, categoryFromDomain(&categories[i])); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// ============================================================
// AnalyticCategory helpers
// ============================================================

func analyticCategoryToDomain(ac AnalyticCategory) *products.AnalyticCategory {
	cat := &products.AnalyticCategory{
		ID:           ac.ID,
		RestaurantID: ac.RestaurantID,
		Name:         ac.Name,
		Slug:         ac.Slug,
		DisplayOrder: ac.DisplayOrder,
		IsActive:     ac.IsActive,
		CreatedAt:    ac.CreatedAt,
		UpdatedAt:    ac.UpdatedAt,
	}
	if ac.Color.Valid {
		cat.Color = ac.Color.String
	}
	return cat
}

func analyticCategoryFromDomain(ac *products.AnalyticCategory) CreateAnalyticCategoryParams {
	params := CreateAnalyticCategoryParams{
		ID:           ac.ID,
		RestaurantID: ac.RestaurantID,
		Name:         ac.Name,
		Slug:         ac.Slug,
		DisplayOrder: ac.DisplayOrder,
		IsActive:     ac.IsActive,
		CreatedAt:    ac.CreatedAt,
		UpdatedAt:    ac.UpdatedAt,
	}
	if ac.Color != "" {
		params.Color = sql.NullString{String: ac.Color, Valid: true}
	}
	return params
}

// ============================================================
// AnalyticCategorySQLCRepository implements products.AnalyticCategoryRepository.
// ============================================================

type AnalyticCategorySQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewAnalyticCategorySQLCRepository creates a new sqlc-backed analytic category repository.
func NewAnalyticCategorySQLCRepository(db *sql.DB) *AnalyticCategorySQLCRepository {
	return &AnalyticCategorySQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func (r *AnalyticCategorySQLCRepository) GetByID(ctx context.Context, id uuid.UUID) (*products.AnalyticCategory, error) {
	row, err := r.queries.GetAnalyticCategoryByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return analyticCategoryToDomain(row), nil
}

func (r *AnalyticCategorySQLCRepository) GetBySlug(ctx context.Context, restaurantID uuid.UUID, slug string) (*products.AnalyticCategory, error) {
	row, err := r.queries.GetAnalyticCategoryBySlug(ctx, GetAnalyticCategoryBySlugParams{
		RestaurantID: restaurantID,
		Slug:         slug,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return analyticCategoryToDomain(row), nil
}

func (r *AnalyticCategorySQLCRepository) ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]products.AnalyticCategory, error) {
	rows, err := r.queries.ListAnalyticCategoriesByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	var out []products.AnalyticCategory
	for _, row := range rows {
		out = append(out, *analyticCategoryToDomain(row))
	}
	return out, nil
}

func (r *AnalyticCategorySQLCRepository) Create(ctx context.Context, category *products.AnalyticCategory) error {
	row, err := r.queries.CreateAnalyticCategory(ctx, analyticCategoryFromDomain(category))
	if err != nil {
		return err
	}
	*category = *analyticCategoryToDomain(row)
	return nil
}

func (r *AnalyticCategorySQLCRepository) Update(ctx context.Context, category *products.AnalyticCategory) error {
	params := UpdateAnalyticCategoryParams{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		DisplayOrder: category.DisplayOrder,
		IsActive:     category.IsActive,
		UpdatedAt:    category.UpdatedAt,
	}
	if category.Color != "" {
		params.Color = sql.NullString{String: category.Color, Valid: true}
	}
	return r.queries.UpdateAnalyticCategory(ctx, params)
}

func (r *AnalyticCategorySQLCRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAnalyticCategory(ctx, id)
}

// ============================================================
// CategoryMapping helpers
// ============================================================

func categoryMappingToDomain(cm CategoryMapping) *products.CategoryMapping {
	m := &products.CategoryMapping{
		ID:                 cm.ID,
		RestaurantID:       cm.RestaurantID,
		CategoryID:         cm.CategoryID,
		AnalyticCategoryID: cm.AnalyticCategoryID,
		MappedAt:           cm.MappedAt,
		CreatedAt:          cm.CreatedAt,
		UpdatedAt:          cm.UpdatedAt,
	}
	if cm.MappedBy.Valid {
		m.MappedBy = cm.MappedBy.String
	}
	return m
}

func categoryMappingFromDomain(cm *products.CategoryMapping) UpsertCategoryMappingParams {
	params := UpsertCategoryMappingParams{
		ID:                 cm.ID,
		RestaurantID:       cm.RestaurantID,
		CategoryID:         cm.CategoryID,
		AnalyticCategoryID: cm.AnalyticCategoryID,
		MappedAt:           cm.MappedAt,
		CreatedAt:          cm.CreatedAt,
		UpdatedAt:          cm.UpdatedAt,
	}
	params.MappedBy = sql.NullString{String: cm.MappedBy, Valid: cm.MappedBy != ""}
	return params
}

// ============================================================
// CategoryMappingSQLCRepository implements products.CategoryMappingRepository.
// ============================================================

type CategoryMappingSQLCRepository struct {
	db      *sql.DB
	queries *Queries
}

// NewCategoryMappingSQLCRepository creates a new sqlc-backed category mapping repository.
func NewCategoryMappingSQLCRepository(db *sql.DB) *CategoryMappingSQLCRepository {
	return &CategoryMappingSQLCRepository{
		db:      db,
		queries: New(db),
	}
}

func (r *CategoryMappingSQLCRepository) GetByCategoryID(ctx context.Context, restaurantID uuid.UUID, categoryID uuid.UUID) (*products.CategoryMapping, error) {
	row, err := r.queries.GetCategoryMappingByCategoryID(ctx, GetCategoryMappingByCategoryIDParams{
		RestaurantID: restaurantID,
		CategoryID:   categoryID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return categoryMappingToDomain(row), nil
}

func (r *CategoryMappingSQLCRepository) Upsert(ctx context.Context, mapping *products.CategoryMapping) error {
	_, err := r.queries.UpsertCategoryMapping(ctx, categoryMappingFromDomain(mapping))
	return err
}
