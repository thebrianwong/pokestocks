-- +goose Up
ALTER TABLE pokemon_stock_pairs
RENAME COLUMN "pokemonid" TO "pokemon_id";

ALTER TABLE pokemon_stock_pairs
RENAME COLUMN "stockid" TO "stock_id";

-- +goose Down
ALTER TABLE pokemon_stock_pairs
RENAME COLUMN "stock_id" TO "stockid";

ALTER TABLE pokemon_stock_pairs
RENAME COLUMN "pokemon_id" TO "pokemonid";