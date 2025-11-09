-- +goose Up
-- +goose StatementBegin
CREATE TABLE content_label (
    -- Дата создания
    created_at TIMESTAMPTZ NOT NULL,
    -- ContentID: вариант "фильм"
    movie_id BIGINT,
    -- ContentID: вариант "сериал"
    tvshow_id BIGINT,
    -- Номер сезона сериала
    season_number INTEGER,
    -- Тип лейбла
    type_label INTEGER NOT NULL,
    -- Ограничение для фильмов: один тип лейбла на фильм
    CONSTRAINT unique_movie_label UNIQUE (movie_id, type_label),
    -- Ограничение для сериалов: один тип лейбла на сезон сериала
    CONSTRAINT unique_tvshow_label UNIQUE (tvshow_id, season_number, type_label),
    -- Ровно один вариант ContentID: либо movie, либо tvshow+season
    CONSTRAINT content_choice_chk CHECK (
    (movie_id IS NOT NULL AND tvshow_id IS NULL AND season_number IS NULL)
        OR
    (movie_id IS NULL AND tvshow_id IS NOT NULL AND season_number IS NOT NULL)
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE content_label;
-- +goose StatementEnd
