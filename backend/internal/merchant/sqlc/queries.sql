-- name: GetMerchantByID :one
SELECT id, name, clover_merchant_id,
    access_token, refresh_token,
    access_token_expires_at, refresh_token_expires_at,
    clover_env, is_connected, is_active,
    ticket_completo_rules, created_at, updated_at
FROM restaurants
WHERE id = $1;

-- name: GetMerchantByCloverMerchantID :one
SELECT id, name, clover_merchant_id,
    access_token, refresh_token,
    access_token_expires_at, refresh_token_expires_at,
    clover_env, is_connected, is_active,
    ticket_completo_rules, created_at, updated_at
FROM restaurants
WHERE clover_merchant_id = $1;

-- name: ListMerchants :many
SELECT id, name, clover_merchant_id,
    access_token, refresh_token,
    access_token_expires_at, refresh_token_expires_at,
    clover_env, is_connected, is_active,
    ticket_completo_rules, created_at, updated_at
FROM restaurants
WHERE is_active = true;

-- name: CreateMerchant :exec
INSERT INTO restaurants (
    id, name, clover_merchant_id,
    access_token, refresh_token,
    access_token_expires_at, refresh_token_expires_at,
    clover_env, is_connected, is_active,
    ticket_completo_rules, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: UpdateMerchant :exec
UPDATE restaurants SET
    name = $2,
    clover_merchant_id = $3,
    clover_env = $4,
    is_connected = $5,
    is_active = $6,
    ticket_completo_rules = $7,
    updated_at = $8
WHERE id = $1;

-- name: UpdateMerchantTokens :exec
UPDATE restaurants SET
    access_token = $2,
    refresh_token = $3,
    access_token_expires_at = $4,
    refresh_token_expires_at = $5,
    updated_at = $6
WHERE id = $1;

-- name: DeleteMerchant :exec
UPDATE restaurants SET is_active = false, updated_at = $2 WHERE id = $1;
