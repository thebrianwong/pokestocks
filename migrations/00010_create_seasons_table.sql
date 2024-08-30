-- +goose Up
CREATE TABLE seasons (
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  name VARCHAR NOT NULL,
  active BOOLEAN,
  PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE seasons;
