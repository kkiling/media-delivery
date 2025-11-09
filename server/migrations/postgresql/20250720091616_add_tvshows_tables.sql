-- +goose Up
-- +goose StatementBegin

-- Таблица для хранения изображений
CREATE TABLE images (
    id TEXT PRIMARY KEY,
    w92 TEXT,
    w185 TEXT,
    w342 TEXT,
    original TEXT NOT NULL
);

-- Таблица для хранения TVShow
CREATE TABLE tv_shows (
    id BIGINT PRIMARY KEY CHECK (id >= 0), -- uint64
    name TEXT NOT NULL,
    original_name TEXT NOT NULL,
    overview TEXT NOT NULL,
    poster_id TEXT,
    first_air_date TIMESTAMPTZ NOT NULL,
    vote_average REAL NOT NULL,
    vote_count INTEGER NOT NULL,
    popularity REAL  NOT NULL,
    backdrop_id TEXT,
    genres TEXT[],
    last_air_date TIMESTAMPTZ NOT NULL,
    number_of_episodes INTEGER  NOT NULL CHECK (number_of_episodes >= 0),
    number_of_seasons INTEGER  NOT NULL CHECK (number_of_seasons >= 0),
    origin_country TEXT[],
    status TEXT,
    tagline TEXT NOT NULL,
    type TEXT
);

-- Таблица для хранения сезонов
CREATE TABLE seasons (
    tv_show_id BIGINT NOT NULL,
    season_number INTEGER NOT NULL CHECK (season_number >= 0),
    air_date TIMESTAMPTZ NOT NULL,
    episode_count INTEGER NOT NULL,
    name TEXT NOT NULL,
    overview TEXT NOT NULL,
    poster_id TEXT,
    vote_average REAL NOT NULL
);

CREATE UNIQUE INDEX idx_seasons_tv_show_id_season_number ON seasons (tv_show_id, season_number);


CREATE TABLE episodes (
    tv_show_id BIGINT NOT NULL,
    season_number INTEGER NOT NULL CHECK (season_number >= 0),
    air_date TIMESTAMPTZ NOT NULL,
    episode_number INTEGER NOT NULL,
    episode_type TEXT,
    name TEXT NOT NULL,
    overview TEXT NOT NULL,
    runtime INTEGER NOT NULL,
    still_id TEXT,
    vote_average REAL NOT NULL,
    vote_count INTEGER NOT NULL
);

CREATE UNIQUE INDEX idx_episodes_tv_show_id_season_number_episode_number ON episodes (tv_show_id, season_number, episode_number);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tv_shows;
DROP TABLE seasons;
DROP TABLE images;
-- +goose StatementEnd
