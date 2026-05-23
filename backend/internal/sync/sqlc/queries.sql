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
