-- name: GetEmployeeByID :one
SELECT * FROM employees
WHERE id = $1 LIMIT 1;

-- name: GetEmployeeByCloverID :one
SELECT * FROM employees
WHERE restaurant_id = $1 AND clover_employee_id = $2
LIMIT 1;

-- name: ListEmployeesByRestaurant :many
SELECT * FROM employees
WHERE restaurant_id = $1
ORDER BY name;

-- name: UpsertEmployee :one
INSERT INTO employees (
    id,
    restaurant_id,
    clover_employee_id,
    name,
    email,
    phone,
    role,
    is_active,
    created_at,
    updated_at,
    synced_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
ON CONFLICT (restaurant_id, clover_employee_id)
DO UPDATE SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    phone = EXCLUDED.phone,
    role = EXCLUDED.role,
    is_active = EXCLUDED.is_active,
    updated_at = EXCLUDED.updated_at,
    synced_at = EXCLUDED.synced_at
RETURNING *;

-- name: DeactivateEmployee :exec
UPDATE employees
SET is_active = false, updated_at = now()
WHERE id = $1;
