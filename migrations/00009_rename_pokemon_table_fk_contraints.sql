-- +goose Up
ALTER TABLE pokemon
RENAME CONSTRAINT "pokemon_pokemon_types_fk_type_1_id"
TO "fk_pokemon_pokemon_types_type_1_id";

ALTER TABLE pokemon
RENAME CONSTRAINT "pokemon_pokemon_types_fk_type_2_id"
TO "fk_pokemon_pokemon_types_type_2_id";

-- +goose Down
ALTER TABLE pokemon
RENAME CONSTRAINT "fk_pokemon_pokemon_types_type_2_id"
TO "pokemon_pokemon_types_fk_type_2_id";

ALTER TABLE pokemon
RENAME CONSTRAINT "fk_pokemon_pokemon_types_type_1_id"
TO "pokemon_pokemon_types_fk_type_1_id";