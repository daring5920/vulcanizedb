-- +goose Up
ALTER TABLE blocks
  ALTER COLUMN size TYPE VARCHAR USING size::VARCHAR;


-- +goose Down
ALTER TABLE blocks
  ALTER COLUMN size TYPE BIGINT USING size::BIGINT;