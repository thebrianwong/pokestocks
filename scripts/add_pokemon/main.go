package main

import (
	"encoding/json"
	"log"
	"os"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func readNameJson(file string) []string {
	data := []string{}

	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening json file: %v", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&data)
	if err != nil {
		log.Fatalf("Error reading name data into memory: %v", err)
	}

	return data
}

type typeSpriteData struct {
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
		log.Fatalf("Error opening json file: %v", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	caser := cases.Title(language.English, cases.NoLower)

	rawData := []rawData{}
	err = dec.Decode(&rawData)
	if err != nil {
		log.Fatalf("Error reading type and sprite data into memory: %v", err)
	}
	for _, obj := range rawData {
		var type1 string
		var type2 string
		var sprite string

		types := obj.PokemonAggregate.Pokemon[0].PokemonTypes
		type1 = types[0].PokemonTypeName.Name
		type1 = caser.String(type1)
		if len(types) == 2 {
			type2 = types[1].PokemonTypeName.Name
			type2 = caser.String(type2)
		}
		sprite = obj.PokemonAggregate.Pokemon[0].SpritesAggregate.Sprites[0].SpriteUrl

		objData := typeSpriteData{type1, type2, sprite}
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

func main() {
	names := readNameJson("../../data/pokemon_names.json")
	types := readTypeSpriteJson("../../data/pokemon_types_sprites.json")

	transformedData := combineData(names, types)

}
