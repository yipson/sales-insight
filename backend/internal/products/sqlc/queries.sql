-- ==========================================
-- Categories (native Clover categories)
-- ==========================================

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: GetCategoryByCloverID :one
SELECT * FROM categories
WHERE restaurant_id = $1 AND clover_category_id = $2
LIMIT 1;

-- name: ListCategoriesByRestaurant :many
SELECT * FROM categories
WHERE restaurant_id = $1
ORDER BY sort_order, name;

-- name: UpsertCategory :one
INSERT INTO categories (
    id, restaurant_id, clover_category_id, name, sort_order,
    created_at, updated_at, synced_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (restaurant_id, clover_category_id)
DO UPDATE SET
    name = EXCLUDED.name,
    sort_order = EXCLUDED.sort_order,
    updated_at = EXCLUDED.updated_at,
    synced_at = EXCLUDED.synced_at
RETURNING *;

-- ==========================================
-- Products
-- ==========================================

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: GetProductByCloverItemID :one
SELECT * FROM products
WHERE restaurant_id = $1 AND clover_item_id = $2
LIMIT 1;

-- name: ListProductsByRestaurant :many
SELECT * FROM products
WHERE restaurant_id = $1 AND is_deleted = false
ORDER BY name;

-- name: UpsertProduct :one
INSERT INTO products (
    id, restaurant_id, clover_item_id, name, description, price, cost, sku,
    category_id, is_available, is_deleted,
    created_at, updated_at, synced_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
ON CONFLICT (restaurant_id, clover_item_id)
DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    price = EXCLUDED.price,
    cost = EXCLUDED.cost,
    sku = EXCLUDED.sku,
    category_id = EXCLUDED.category_id,
    is_available = EXCLUDED.is_available,
    is_deleted = EXCLUDED.is_deleted,
    updated_at = EXCLUDED.updated_at,
    synced_at = EXCLUDED.synced_at
RETURNING *;

-- name: SoftDeleteProduct :exec
UPDATE products
SET is_deleted = true, updated_at = now()
WHERE id = $1;

-- ==========================================
-- Analytic Categories
-- ==========================================

-- name: GetAnalyticCategoryByID :one
SELECT * FROM analytic_categories
WHERE id = $1 LIMIT 1;

-- name: GetAnalyticCategoryBySlug :one
SELECT * FROM analytic_categories
WHERE restaurant_id = $1 AND slug = $2
LIMIT 1;

-- name: ListAnalyticCategoriesByRestaurant :many
SELECT * FROM analytic_categories
WHERE restaurant_id = $1 AND is_active = true
ORDER BY display_order, name;

-- name: CreateAnalyticCategory :one
INSERT INTO analytic_categories (
    id, restaurant_id, name, slug, display_order, color, is_active,
    created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateAnalyticCategory :exec
UPDATE analytic_categories SET
    name = $2,
    slug = $3,
    display_order = $4,
    color = $5,
    is_active = $6,
    updated_at = $7
WHERE id = $1;

-- name: DeleteAnalyticCategory :exec
DELETE FROM analytic_categories
WHERE id = $1;

-- ==========================================
-- Category Mappings
-- ==========================================

-- name: GetCategoryMappingByCategoryID :one
SELECT * FROM category_mappings
WHERE restaurant_id = $1 AND category_id = $2
LIMIT 1;

-- name: UpsertCategoryMapping :one
INSERT INTO category_mappings (
    id, restaurant_id, category_id, analytic_category_id, mapped_by, mapped_at,
    created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (restaurant_id, category_id)
DO UPDATE SET
    analytic_category_id = EXCLUDED.analytic_category_id,
    mapped_by = EXCLUDED.mapped_by,
    mapped_at = EXCLUDED.mapped_at,
    updated_at = EXCLUDED.updated_at
RETURNING *;
