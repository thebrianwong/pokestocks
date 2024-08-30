-- +goose Up
CREATE TABLE pokemon_stock_pairs (
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  pokemonId BIGINT NOT NULL,
  stockId BIGINT NOT NULL,
  PRIMARY KEY(id)
);

ALTER TABLE pokemon_stock_pairs
  ADD CONSTRAINT fk_pokemon_stock_pairs_pokemon
  FOREIGN KEY (pokemonId)
  REFERENCES pokemon (id);

ALTER TABLE pokemon_stock_pairs
  ADD CONSTRAINT fk_pokemon_stock_pairs_stocks
  FOREIGN KEY (stockId)
  REFERENCES stocks (id);

-- +goose Down
ALTER TABLE pokemon_stock_pairs
  DROP CONSTRAINT fk_pokemon_stock_pairs_stocks RESTRICT;

ALTER TABLE pokemon_stock_pairs
  DROP CONSTRAINT fk_pokemon_stock_pairs_pokemon RESTRICT;

DROP TABLE pokemon_stock_pairs;