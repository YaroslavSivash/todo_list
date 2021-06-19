
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS notes(
    id          SERIAL PRIMARY KEY,
    body        VARCHAR,
    is_done     BOOLEAN,
    created_at  TIMESTAMP,
    updated_at  TIMESTAMP
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS notes;
