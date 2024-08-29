-- +goose Up
ALTER TABLE stock
    RENAME TO stocks;

-- +goose Down
ALTER TABLE stocks
    RENAME TO stock;
