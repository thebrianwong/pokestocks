package main

import (
	"context"
	"encoding/json"
	"os"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func readNameJson(file string) []string {
	data := []string{}

	f, err := os.Open(file)
	if err != nil {
		utils.LogFailureError("Error opening json file", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&data)
	if err != nil {
		utils.LogFailureError("Error reading name data into memory", err)
	}

	return data
}

type typeSpriteData struct {
	id     int
	type1  string
	type2  string
	sprite string
}

type transformedData struct {
	pokemonName string
	*typeSpriteData
}

type pokemonTypeName struct {
	Name string `json:"name"`
}

type pokemonType struct {
	Slot            int             `json:"slot"`
	PokemonTypeName pokemonTypeName `json:"pokemon_v2_type"`
}

type sprite struct {
	SpriteUrl string `json:"sprites"`
}

type spritesAggregate struct {
	Sprites []sprite `json:"nodes"`
}

type pokemon struct {
	Id               int              `json:"id"`
	Name             string           `json:"name"`
	PokemonTypes     []pokemonType    `json:"pokemon_v2_pokemontypes"`
	SpritesAggregate spritesAggregate `json:"pokemon_v2_pokemonsprites_aggregate"`
}

type pokemonAggregate struct {
	Pokemon []pokemon `json:"nodes"`
}

type rawData struct {
	PokemonAggregate pokemonAggregate `json:"pokemon_v2_pokemons_aggregate"`
}

func readTypeSpriteJson(file string) []typeSpriteData {
	data := []typeSpriteData{}

	f, err := os.Open(file)
	if err != nil {
		utils.LogFailureError("Error opening json file", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	caser := cases.Title(language.English, cases.NoLower)

	rawData := []rawData{}
	err = dec.Decode(&rawData)
	if err != nil {
		utils.LogFailureError("Error reading type and sprite data into memory", err)
	}
	for _, obj := range rawData {
		var id int
		var type1 string
		var type2 string
		var sprite string

		pokemon := obj.PokemonAggregate.Pokemon[0]

		id = pokemon.Id
		types := pokemon.PokemonTypes
		type1 = types[0].PokemonTypeName.Name
		type1 = caser.String(type1)
		if len(types) == 2 {
			type2 = types[1].PokemonTypeName.Name
			type2 = caser.String(type2)
		}
		sprite = pokemon.SpritesAggregate.Sprites[0].SpriteUrl

		objData := typeSpriteData{id, type1, type2, sprite}
		data = append(data, objData)
	}

	return data
}

func combineData(nameData []string, typeSpriteData []typeSpriteData) []transformedData {
	data := []transformedData{}

	for i := 0; i < len(nameData); i++ {
		tsData := typeSpriteData[i]
		pokemonName := nameData[i]

		combinedData := transformedData{
			typeSpriteData: &tsData,
			pokemonName:    pokemonName,
		}

		data = append(data, combinedData)
	}

	return data
}

func insertIntoDb(ctx context.Context, db *pgxpool.Pool, pokemonData []transformedData) {
	options := pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		utils.LogFailureError("Error starting db transaction", err)
	}

	defer tx.Rollback(ctx)

	batch := pgx.Batch{}
	for _, pokemon := range pokemonData {
		if pokemon.type2 != "" {
			query := `
				INSERT INTO pokemon (name, pokedex_number, type_1_id, type_2_id, sprite_url)
				SELECT name, pokedex_number, type_1_id, type_2_id, sprite_url
				FROM (VALUES (
					$1, 
					$2::INTEGER,
					(SELECT id FROM pokemon_types WHERE type=$3),
					(SELECT id FROM pokemon_types WHERE type=$4),
					$5
				))
				AS data(name, pokedex_number, type_1_id, type_2_id, sprite_url)
				WHERE NOT EXISTS (
					SELECT 1
					FROM pokemon
					WHERE pokemon.name = data.name
					AND pokemon.pokedex_number = data.pokedex_number
				)
			`
			batch.Queue(query, pokemon.pokemonName, pokemon.id, pokemon.type1, pokemon.type2, pokemon.sprite)
		} else {
			query := `
				INSERT INTO pokemon (name, pokedex_number, type_1_id, type_2_id, sprite_url)
				SELECT name, pokedex_number, type_1_id, type_2_id, sprite_url
				FROM (VALUES (
					$1, 
					$2::INTEGER,
					(SELECT id FROM pokemon_types WHERE type=$3),
					NULL::BIGINT,
					$4
				))
				AS data(name, pokedex_number, type_1_id, type_2_id, sprite_url)
				WHERE NOT EXISTS (
					SELECT 1
					FROM pokemon
					WHERE pokemon.name = data.name
					AND pokemon.pokedex_number = data.pokedex_number
				)
			`
			batch.Queue(query, pokemon.pokemonName, pokemon.id, pokemon.type1, pokemon.sprite)
		}
	}

	err = db.SendBatch(ctx, &batch).Close()
	if err != nil {
		utils.LogFailureError("Error sending batch inserts", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		utils.LogFailureError("Error committing db transaction", err)
	}
}

func main() {
	utils.LoadEnvVars("../../.env")
	conn := utils.ConnectToDb()

	names := readNameJson("../../data/pokemon_names.json")
	types := readTypeSpriteJson("../../data/pokemon_types_sprites.json")

	transformedData := combineData(names, types)

	insertIntoDb(context.Background(), conn, transformedData)
}
