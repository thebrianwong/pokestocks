-- +goose Up
ALTER TABLE holdings
  ADD CONSTRAINT no_duplicate_holdings
  UNIQUE (portfolio_id, pokemon_stock_pair_id);


-- +goose Down
ALTER TABLE holdings
  DROP CONSTRAINT no_duplicate_holdings;
