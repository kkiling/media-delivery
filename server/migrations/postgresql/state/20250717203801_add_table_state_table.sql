-- +goose Up
-- +goose StatementBegin
CREATE TABLE state (
    id UUID PRIMARY KEY,
    idempotency_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    status INTEGER NOT NULL,
    step TEXT NOT NULL,
    type TEXT NOT NULL,
    error TEXT,
    data JSONB,
    fail_data JSONB,
    meta_data JSONB
);

CREATE UNIQUE INDEX idx_state_idempotency_key ON state(idempotency_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS state;
-- +goose StatementEnd