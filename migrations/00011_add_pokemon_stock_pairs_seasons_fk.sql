-- +goose Up
ALTER TABLE pokemon_stock_pairs
ADD COLUMN season_id BIGINT NOT NULL
  CONSTRAINT fk_pokemon_stock_pairs_seasons REFERENCES seasons (id)
  ON UPDATE CASCADE ON DELETE CASCADE; 

-- +goose Down
ALTER TABLE pokemon_stock_pairs
DROP CONSTRAINT fk_pokemon_stock_pairs_seasons RESTRICT;

ALTER TABLE pokemon_stock_pairs
DROP COLUMN season_id;