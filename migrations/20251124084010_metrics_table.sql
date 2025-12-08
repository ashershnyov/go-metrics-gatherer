-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
