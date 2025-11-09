-- name: SaveLabel :exec
INSERT INTO content_label (movie_id, tvshow_id, season_number, created_at, type_label)
VALUES ($1, $2, $3, $4, $5);

-- name: GetLabelsMovieID :many
SELECT created_at, type_label FROM content_label
WHERE movie_id=$1;

-- name: GetLabelsTVShow :many
SELECT created_at, type_label FROM content_label
WHERE tvshow_id=$1 AND season_number=$2;

-- name: DeleteLabelMovieID :one
DELETE FROM content_label
WHERE movie_id=$1
RETURNING movie_id;

-- name: DeleteLabelTVShow :one
DELETE FROM content_label
WHERE tvshow_id=$1 AND season_number=$2
RETURNING tvshow_id, season_number;