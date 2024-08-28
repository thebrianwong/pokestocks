-- +goose Up
ALTER TABLE stock
ADD COLUMN active BOOLEAN DEFAULT TRUE;

-- +goose Down
ALTER TABLE stock
DROP COLUMN active;
