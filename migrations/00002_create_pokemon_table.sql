-- +goose Up
CREATE TABLE pokemon (
  id bigint GENERATED ALWAYS AS IDENTITY,
  name varchar NOT NULL,
  pokedex_number int NOT NULL,
  PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE pokemon;
