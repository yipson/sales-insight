-- ==========================================
-- Orders
-- ==========================================

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrderByCloverOrderID :one
SELECT * FROM orders
WHERE restaurant_id = $1 AND clover_order_id = $2
LIMIT 1;

-- name: ListOrdersByRestaurantAndDate :many
SELECT * FROM orders
WHERE restaurant_id = $1
  AND created_time >= $2
  AND created_time < $3
ORDER BY created_time DESC;

-- name: UpsertOrder :one
INSERT INTO orders (
    id, restaurant_id, clover_order_id, employee_id, order_number, order_type, state,
    total_amount, tax_amount, discount_amount, tip_amount,
    created_time, modified_time, pay_type,
    item_count, unique_category_count, total_quantity, synced_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
)
ON CONFLICT (restaurant_id, clover_order_id)
DO UPDATE SET
    employee_id = EXCLUDED.employee_id,
    order_number = EXCLUDED.order_number,
    order_type = EXCLUDED.order_type,
    state = EXCLUDED.state,
    total_amount = EXCLUDED.total_amount,
    tax_amount = EXCLUDED.tax_amount,
    discount_amount = EXCLUDED.discount_amount,
    tip_amount = EXCLUDED.tip_amount,
    created_time = EXCLUDED.created_time,
    modified_time = EXCLUDED.modified_time,
    pay_type = EXCLUDED.pay_type,
    item_count = EXCLUDED.item_count,
    unique_category_count = EXCLUDED.unique_category_count,
    total_quantity = EXCLUDED.total_quantity,
    synced_at = EXCLUDED.synced_at
RETURNING *;

-- name: DeleteOrderByCloverID :exec
DELETE FROM orders
WHERE restaurant_id = $1 AND clover_order_id = $2;

-- name: ListOrdersByEmployeeAndDate :many
SELECT * FROM orders
WHERE employee_id = $1
  AND created_time >= $2
  AND created_time < $3
ORDER BY created_time DESC;

-- ==========================================
-- Order Items
-- ==========================================

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: UpsertOrderItem :one
INSERT INTO order_items (
    id, restaurant_id, order_id, product_id, clover_line_item_id, name,
    quantity, unit_price, total_price, analytic_category_id, created_at, synced_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
ON CONFLICT DO NOTHING
RETURNING *;

-- ==========================================
-- Order Category Summary
-- ==========================================

-- name: GetOrderCategorySummaryByOrderID :many
SELECT * FROM order_category_summary
WHERE order_id = $1;

-- name: UpsertOrderCategorySummary :one
INSERT INTO order_category_summary (
    id, restaurant_id, order_id, analytic_category_id, item_count, total_quantity, total_amount
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
ON CONFLICT (restaurant_id, order_id, analytic_category_id)
DO UPDATE SET
    item_count = EXCLUDED.item_count,
    total_quantity = EXCLUDED.total_quantity,
    total_amount = EXCLUDED.total_amount
RETURNING *;
