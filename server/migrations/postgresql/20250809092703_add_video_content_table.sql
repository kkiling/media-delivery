-- +goose Up
-- +goose StatementBegin
CREATE TABLE video_content (
    -- Идентификатор записи (uuid.UUID как строка)
    id UUID PRIMARY KEY NOT NULL,
    -- Дата создания 
    created_at TIMESTAMPTZ NOT NULL,
    -- ContentID: вариант "фильм"
    movie_id BIGINT,
    -- ContentID: вариант "сериал"
    tvshow_id BIGINT,
    -- Номер сезона сериала
    season_number INTEGER,
    -- Статус доставки
    delivery_status INTEGER NOT NULL,
    -- Массив состояний State в JSON:
    states JSONB,
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
DROP TABLE video_content;
-- +goose StatementEnd
