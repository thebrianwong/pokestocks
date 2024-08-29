-- +goose Up
ALTER TABLE pokemon
ADD COLUMN type_1_id BIGINT NOT NULL
  CONSTRAINT pokemon_pokemon_types_fk_type_1_id REFERENCES pokemon_types (id)
  ON UPDATE CASCADE ON DELETE CASCADE; 

ALTER TABLE pokemon
ADD COLUMN type_2_id BIGINT NOT NULL
  CONSTRAINT pokemon_pokemon_types_fk_type_2_id REFERENCES pokemon_types (id)
  ON UPDATE CASCADE ON DELETE CASCADE; 

ALTER TABLE pokemon
ADD COLUMN sprite_url VARCHAR;

-- +goose Down
ALTER TABLE pokemon
DROP COLUMN sprite_url;

ALTER TABLE pokemon
DROP CONSTRAINT pokemon_pokemon_types_fk_type_2_id RESTRICT;

ALTER TABLE pokemon
DROP COLUMN type_2_id;

ALTER TABLE pokemon
DROP CONSTRAINT pokemon_pokemon_types_fk_type_1_id RESTRICT;

ALTER TABLE pokemon
DROP COLUMN type_1_id;