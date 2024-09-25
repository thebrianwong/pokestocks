package helpers

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"pokestocks/internal/psp_elastic_document"
	common_pb "pokestocks/proto/common"
	"slices"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PspQueryString() string {
	return `
		SELECT 
			psp.id AS "pspId",
			pokemon.id AS "pokemonId",
			pokemon."name" AS "pokemonName",
			pokemon.pokedex_number AS "pokedexNumber",
			pokemon.created_at AS "pokemonCreatedAt",
			pokemon.updated_at AS "pokemonUpdatedAt",
			pokemon.sprite_url as "pokemonSpriteUrl",
			pokemon_types1.id AS "type1Id",
			pokemon_types1."type" AS "type1Name",
			pokemon_types1.sprite_url AS "type1SpriteUrl",
			pokemon_types2.id AS "type2Id",
			pokemon_types2."type" AS "type2Name",
			pokemon_types2.sprite_url AS "type2SpriteUrl",
			stocks.id AS "stockId",
			stocks.symbol AS "stockSymbol", 
			stocks."name" AS "stockName",
			stocks.created_at AS "stockCreatedAt",
			stocks.updated_at AS "stockUpdatedAt",
			stocks.active AS "stockActive",
			seasons.id AS "seasonId",
			seasons."name" AS "seasonName",
			seasons.active AS "seasonActive"
		FROM pokemon_stock_pairs AS psp
		INNER JOIN pokemon ON pokemon.id = psp.pokemon_id
		INNER JOIN pokemon_types AS pokemon_types1 ON pokemon_types1.id = pokemon.type_1_id
		LEFT JOIN pokemon_types AS pokemon_types2 ON pokemon_types2.id = pokemon.type_2_id
		INNER JOIN stocks ON stocks.id = psp.stock_id
		INNER JOIN seasons ON seasons.id = psp.season_id
	`
}

func ExtractPokemonStockPairIds(docs []psp_elastic_document.PspElasticDocument) []string {
	var ids []string

	for _, doc := range docs {
		id := doc.Id
		ids = append(ids, fmt.Sprint(id))
	}

	return ids
}

func ConvertPokemonStockPairElasticDocuments(elasticResponse *search.Response) ([]psp_elastic_document.PspElasticDocument, error) {
	var convertedDocuments []psp_elastic_document.PspElasticDocument

	docs := elasticResponse.Hits.Hits
	for _, doc := range docs {
		var convertedDocument psp_elastic_document.PspElasticDocument
		docData := doc.Source_
		err := json.Unmarshal(docData, &convertedDocument)
		if err != nil {
			return nil, err
		}
		convertedDocuments = append(convertedDocuments, convertedDocument)
	}

	return convertedDocuments, nil
}

func MidnightTomorrow() time.Time {
	today := time.Now()
	midnightTomorrow := time.Date(today.Year(), today.Month(), today.Day()+1, 0, 0, 0, 0, today.Location())
	return midnightTomorrow
}

func GenerateRandomIndices(indicesCount int, sliceLength int) []int {
	var indices []int

	for len(indices) < indicesCount {
		index := rand.IntN(sliceLength)
		if !slices.Contains(indices, index) {
			indices = append(indices, index)
		}
	}

	return indices
}

func ConvertDbRowToPokemonStockPair(rowDataMap map[string]any) *common_pb.PokemonStockPair {
	psp := common_pb.PokemonStockPair{
		Id: rowDataMap["pspId"].(int64),
		Pokemon: &common_pb.Pokemon{
			Id:            rowDataMap["pokemonId"].(int64),
			Name:          rowDataMap["pokemonName"].(string),
			PokedexNumber: rowDataMap["pokedexNumber"].(int32),
			CreatedAt:     timestamppb.New(rowDataMap["pokemonCreatedAt"].(time.Time)),
			UpdatedAt:     timestamppb.New(rowDataMap["pokemonUpdatedAt"].(time.Time)),
			Type1: &common_pb.PokemonType{
				Id:        rowDataMap["type1Id"].(int64),
				Type:      rowDataMap["type1Name"].(string),
				SpriteUrl: rowDataMap["type1SpriteUrl"].(string),
			},
			SpriteUrl: rowDataMap["pokemonSpriteUrl"].(string),
		},
		Stock: &common_pb.Stock{
			Id:        rowDataMap["stockId"].(int64),
			Symbol:    rowDataMap["stockSymbol"].(string),
			Name:      rowDataMap["stockName"].(string),
			CreatedAt: timestamppb.New(rowDataMap["stockCreatedAt"].(time.Time)),
			UpdatedAt: timestamppb.New(rowDataMap["stockUpdatedAt"].(time.Time)),
			Active:    rowDataMap["stockActive"].(bool),
		},
		Season: &common_pb.Season{
			Id:     rowDataMap["seasonId"].(int64),
			Name:   rowDataMap["seasonName"].(string),
			Active: rowDataMap["seasonActive"].(bool),
		},
	}

	// not all Pokemon have a second type
	type2Id, ok := rowDataMap["type2Id"].(int64)
	if ok {
		psp.Pokemon.Type2 = &common_pb.PokemonType{
			Id:        type2Id,
			Type:      rowDataMap["type2Name"].(string),
			SpriteUrl: rowDataMap["type2SpriteUrl"].(string),
		}
	}

	return &psp
}
