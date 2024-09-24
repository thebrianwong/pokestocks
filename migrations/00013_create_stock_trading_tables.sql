-- +goose Up
CREATE TABLE portfolios (
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  description VARCHAR NOT NULL,
  cash NUMERIC NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE transactions (
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  portfolio_id BIGINT NOT NULL,
  pokemon_stock_pair_id BIGINT NOT NULL,
  quantity INT NOT NULL,
  price NUMERIC NOT NULL,
  buy BOOLEAN NOT NULL,
  date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(id)
);

ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_portfolios
  FOREIGN KEY (portfolio_id)
  REFERENCES portfolios (id);

ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_pokemon_stock_pairs
  FOREIGN KEY (pokemon_stock_pair_id)
  REFERENCES pokemon_stock_pairs (id);

CREATE TABLE holdings (
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  portfolio_id BIGINT NOT NULL,
  pokemon_stock_pair_id BIGINT NOT NULL,
  quantity INT NOT NULL,
  PRIMARY KEY(id)
);

ALTER TABLE holdings
  ADD CONSTRAINT fk_holdings_portfolios
  FOREIGN KEY (portfolio_id)
  REFERENCES portfolios (id);

ALTER TABLE holdings
  ADD CONSTRAINT fk_holdings_pokemon_stock_pairs
  FOREIGN KEY (pokemon_stock_pair_id)
  REFERENCES pokemon_stock_pairs (id);

-- +goose Down
ALTER TABLE holdings
  DROP CONSTRAINT fk_holdings_pokemon_stock_pairs RESTRICT;

ALTER TABLE holdings
  DROP CONSTRAINT fk_holdings_portfolios RESTRICT;

DROP TABLE holdings;

ALTER TABLE transactions
  DROP CONSTRAINT fk_transactions_pokemon_stock_pairs RESTRICT;

ALTER TABLE transactions
  DROP CONSTRAINT fk_transactions_portfolios RESTRICT;

DROP TABLE transactions;

DROP TABLE portfolios;