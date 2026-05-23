-- name: GetLatestSyncByEntity :many
SELECT DISTINCT ON (entity) *
FROM sync_logs
WHERE restaurant_id = $1
ORDER BY entity, started_at DESC;

-- name: ListSyncLogs :many
SELECT *
FROM sync_logs
WHERE restaurant_id = $1
ORDER BY started_at DESC
LIMIT $2;

-- name: ListSyncErrorsUnresolved :many
SELECT *
FROM sync_errors
WHERE restaurant_id = $1 AND is_resolved = false
ORDER BY created_at DESC
LIMIT $2;

-- name: CreateSyncLog :one
INSERT INTO sync_logs (
    restaurant_id, entity, status, records_processed,
    records_inserted, records_updated, records_skipped, records_failed,
    cursor_from, cursor_to, triggered_by, details
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, started_at;
