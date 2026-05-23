-- name: GetPaymentByID :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: GetPaymentByCloverID :one
SELECT * FROM payments
WHERE restaurant_id = $1 AND clover_payment_id = $2
LIMIT 1;

-- name: ListPaymentsByRestaurantAndDate :many
SELECT * FROM payments
WHERE restaurant_id = $1
  AND created_time >= $2
  AND created_time < $3
ORDER BY created_time DESC;

-- name: UpsertPayment :one
INSERT INTO payments (
    id, restaurant_id, clover_payment_id, order_id, employee_id,
    amount, tip_amount, tax_amount, payment_type, card_type, result,
    external_payment_id, created_time, modified_time, synced_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
ON CONFLICT (restaurant_id, clover_payment_id)
DO UPDATE SET
    order_id = EXCLUDED.order_id,
    employee_id = EXCLUDED.employee_id,
    amount = EXCLUDED.amount,
    tip_amount = EXCLUDED.tip_amount,
    tax_amount = EXCLUDED.tax_amount,
    payment_type = EXCLUDED.payment_type,
    card_type = EXCLUDED.card_type,
    result = EXCLUDED.result,
    external_payment_id = EXCLUDED.external_payment_id,
    created_time = EXCLUDED.created_time,
    modified_time = EXCLUDED.modified_time,
    synced_at = EXCLUDED.synced_at
RETURNING *;
