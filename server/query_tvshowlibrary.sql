------------------------------------------------------------------------------------------------------------------------
-- name: SaveImage :exec
INSERT INTO images (id, w92, w185, w342, original)
VALUES ($1, $2, $3, $4, $5);

-- name: GetImage :one
SELECT id, w92, w185, w342, original FROM images WHERE id = $1 LIMIT 1;

-- name: SaveSeason :exec
INSERT INTO seasons (
    tv_show_id, season_number, air_date, episode_count,
    name, overview, poster_id, vote_average
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);


-- name: GetSeasons :many
SELECT air_date, episode_count, name, overview, poster_id, season_number, vote_average
FROM seasons
WHERE tv_show_id = $1
ORDER BY season_number;

-- name: SaveTVShow :exec
INSERT INTO tv_shows (
    id, name, original_name, overview, poster_id,
    first_air_date, vote_average, vote_count, popularity,
    backdrop_id, genres, last_air_date,
    number_of_episodes, number_of_seasons, origin_country,
    status, tagline, type
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8,
          $9, $10, $11, $12, $13, $14,
          $15, $16, $17, $18);


-- name: GetTVShow :one
SELECT
    id, name, original_name, overview, poster_id,
    first_air_date, vote_average, vote_count, popularity,
    backdrop_id, genres, last_air_date,
    number_of_episodes, number_of_seasons, origin_country,
    status, tagline, type
FROM tv_shows
WHERE id = $1;

-- name: GetTVShows :many
SELECT
    id, name, original_name, overview, poster_id,
    first_air_date, vote_average, vote_count, popularity
FROM tv_shows
ORDER BY popularity DESC;

-- name: GetEpisodes :many
SELECT air_date, episode_number,
       episode_type, name, overview, runtime,
       still_id, vote_average, vote_count
FROM episodes
WHERE tv_show_id = $1 AND season_number = $2;

-- name: SaveEpisode :exec
INSERT INTO episodes (
    tv_show_id, season_number, air_date, episode_number,
    episode_type, name, overview, runtime,
    still_id, vote_average, vote_count
) VALUES ($1, $2, $3, $4,
          $5, $6, $7, $8, $9,
          $10, $11);

