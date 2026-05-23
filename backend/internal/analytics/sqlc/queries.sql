-- name: GetSalesSummary :one
SELECT
    COUNT(*)::int AS order_count,
    COALESCE(SUM(total_amount), 0)::bigint AS total_sales,
    COALESCE(AVG(total_amount), 0)::float8 AS avg_ticket,
    COALESCE(SUM(tax_amount), 0)::bigint AS total_tax,
    COALESCE(SUM(tip_amount), 0)::bigint AS total_tips,
    COALESCE(SUM(discount_amount), 0)::bigint AS total_discounts
FROM orders
WHERE restaurant_id = $1
  AND created_time >= $2
  AND created_time < $3;

-- name: GetSalesByEmployee :many
SELECT
    e.id AS employee_id,
    e.name AS employee_name,
    COUNT(o.id)::int AS order_count,
    COALESCE(SUM(o.total_amount), 0)::bigint AS total_sales,
    COALESCE(AVG(o.total_amount), 0)::float8 AS avg_ticket
FROM employees e
LEFT JOIN orders o ON e.id = o.employee_id AND o.created_time >= $2 AND o.created_time < $3
WHERE e.restaurant_id = $1 AND e.is_active = true
GROUP BY e.id, e.name
ORDER BY total_sales DESC;

-- name: GetTopProducts :many
SELECT
    p.id AS product_id,
    p.name AS product_name,
    SUM(oi.quantity)::bigint AS total_quantity,
    SUM(oi.total_price)::bigint AS total_revenue
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
JOIN products p ON oi.product_id = p.id
WHERE o.restaurant_id = $1
  AND o.created_time >= $2
  AND o.created_time < $3
  AND p.is_deleted = false
GROUP BY p.id, p.name
ORDER BY total_quantity DESC
LIMIT $4;

-- name: GetCategoryCoverage :many
SELECT
    e.id AS employee_id,
    e.name AS employee_name,
    COUNT(DISTINCT ocs.analytic_category_id)::int AS categories_covered
FROM employees e
JOIN orders o ON e.id = o.employee_id
JOIN order_category_summary ocs ON o.id = ocs.order_id
WHERE e.restaurant_id = $1
  AND o.created_time >= $2
  AND o.created_time < $3
GROUP BY e.id, e.name
ORDER BY categories_covered DESC;

-- name: CountActiveAnalyticCategories :one
SELECT COUNT(*)::int AS total_categories
FROM analytic_categories
WHERE restaurant_id = $1 AND is_active = true;

-- name: GetOrderCategoryBreakdown :many
SELECT
    ocs.order_id,
    ac.slug AS category_slug,
    ocs.total_quantity
FROM order_category_summary ocs
JOIN analytic_categories ac ON ocs.analytic_category_id = ac.id
JOIN orders o ON ocs.order_id = o.id
WHERE ocs.restaurant_id = $1
  AND o.created_time >= $2
  AND o.created_time < $3
  AND ac.is_active = true
ORDER BY ocs.order_id;
