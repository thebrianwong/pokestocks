-- +goose Up
CREATE TABLE pokemon_types (
  id bigint GENERATED ALWAYS AS IDENTITY,
  type varchar NOT NULL,
  sprite_url varchar NOT NULL,
  PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE pokemon_types;