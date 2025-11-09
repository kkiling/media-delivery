-- +goose Up
-- +goose StatementBegin
CREATE TABLE step_execute_info (
    id SERIAL PRIMARY KEY,
    state_id UUID NOT NULL,
    start_executed_at TIMESTAMPTZ NOT NULL,
    complete_executed_at TIMESTAMPTZ NOT NULL,
    error TEXT,
    preview_step TEXT NOT NULL,
    next_step TEXT,
    FOREIGN KEY (state_id) REFERENCES state(id) ON DELETE CASCADE
);

CREATE INDEX idx_step_execute_state_id ON step_execute_info(state_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE step_execute_info;
-- +goose StatementEnd