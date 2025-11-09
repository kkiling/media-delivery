-- name: CreateMkvMerge :exec
INSERT INTO mkv_merge (id, idempotency_key, params, status, created_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateMkvMerge :one
UPDATE mkv_merge
SET
    status = $1,
    error = COALESCE($2, error),
    completed_at = COALESCE($3, completed_at)
WHERE id = $4
RETURNING id;

-- name: UpdateProgress :one
UPDATE mkv_merge
SET progress = $1
WHERE id=$2
RETURNING id;

-- name: GetByID :one
SELECT id, idempotency_key, params, status, error, created_at, completed_at, progress FROM mkv_merge
WHERE id=$1;

-- name: GetByIdempotencyKey :one
SELECT id, idempotency_key, params, status, error, created_at, completed_at, progress FROM mkv_merge
WHERE idempotency_key=$1;

-- name: GetOldestUncompleted :one
SELECT id, idempotency_key, params, status, error, created_at, completed_at, progress FROM mkv_merge
WHERE completed_at is null
ORDER BY created_at
LIMIT 1;

-- name: AddMergeLogs :exec
INSERT INTO mkv_merge_logs (merge_id, created_at, type, content)
VALUES ($1, $2, $3, $4);

-- name: DeleteLogs :exec
DELETE from mkv_merge_logs where merge_id=$1;
