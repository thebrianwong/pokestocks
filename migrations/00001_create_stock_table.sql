-- +goose Up
CREATE TABLE stock (
  id bigint GENERATED ALWAYS AS IDENTITY,
  symbol varchar NOT NULL,
  name varchar NOT NULL,
  PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE stock;