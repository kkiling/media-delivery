-- name: SaveVideoContent :exec
INSERT INTO video_content (id, created_at, movie_id, tvshow_id, season_number, delivery_status, states)
VALUES ($1, $2, $3, $4, $5, $6,$7);

-- name: GetVideoContentsMovieID :many
SELECT id, created_at,delivery_status, states FROM video_content
WHERE movie_id=$1;

-- name: GetVideoContentTVShow :many
SELECT id, created_at, delivery_status, states FROM video_content
WHERE tvshow_id=$1 AND season_number=$2;

-- name: GetVideoContentsByDeliveryStatus :many
SELECT id, created_at, movie_id, tvshow_id, season_number, delivery_status, states FROM video_content
WHERE delivery_status=ANY($1::int[]) ORDER BY created_at DESC limit $2;

-- name: UpdateVideoContent :one
UPDATE video_content
SET delivery_status=$1, states=$3
WHERE id=$2
RETURNING id;