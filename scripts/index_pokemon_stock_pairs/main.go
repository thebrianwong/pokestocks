package main

import (
	"context"
	"pokestocks/internal/structs"
	"pokestocks/utils"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	"github.com/jackc/pgx/v5"
)

func getPspIndexDocCount(elasticClient *elasticsearch.TypedClient) (*(count.Response), error) {
	count, err := elasticClient.Count().Index("pokemon_stock_pairs_index").Do(context.Background())

	if err != nil {
		return nil, err
	}

	return count, nil
}

func convertDbRowToIndexPayload(rowDataMap map[string]any) structs.PspElasticDocument {
	payload := structs.PspElasticDocument{
		Id: rowDataMap["pspId"].(int64),
		Pokemon: structs.PspNestedPokemon{
			Id:            rowDataMap["pokemonId"].(int64),
			Name:          rowDataMap["pokemonName"].(string),
			PokedexNumber: rowDataMap["pokedexNumber"].(int32),
			Type1:         rowDataMap["type1Name"].(string),
		},
		Stock: structs.PspNestedStock{
			Id:     rowDataMap["stockId"].(int64),
			Symbol: rowDataMap["stockSymbol"].(string),
			Name:   rowDataMap["stockName"].(string),
			Active: rowDataMap["stockActive"].(bool),
		},
		ActiveSeason: rowDataMap["seasonActive"].(bool),
	}

	hasSecondType := rowDataMap["type2Name"]
	if hasSecondType != nil {
		payload.Pokemon.Type2 = rowDataMap["type2Name"].(string)
	}

	return payload
}

func indexPspPayload(elasticClient *elasticsearch.TypedClient, payload structs.PspElasticDocument) error {
	_, err := elasticClient.Index("pokemon_stock_pairs_index").Request(payload).Do(context.Background())

	return err
}

func main() {
	utils.LoadEnvVars("../../.env")
	conn := utils.ConnectToDb()
	elasticClient := utils.ConnectToElastic("../../http_ca.crt")

	count, err := getPspIndexDocCount(elasticClient)
	if err != nil {
		utils.LogFailureError("Error starting PSP indexing", err)
	}
	if count.Count > 0 {
		utils.LogFailure("Error starting PSP indexing: the index already contains PSPs")
	}

	query := `
		SELECT 
			psp.id AS "pspId",
			pokemon.id AS "pokemonId",
			pokemon."name" AS "pokemonName",
			pokemon.pokedex_number AS "pokedexNumber",
			pokemon_types1."type" AS "type1Name",
			pokemon_types2."type" AS "type2Name",
			stocks.id AS "stockId",
			stocks.symbol AS "stockSymbol", 
			stocks."name" AS "stockName",
			stocks.active AS "stockActive",
			seasons.active AS "seasonActive"
		FROM pokemon_stock_pairs AS psp
		INNER JOIN pokemon ON pokemon.id = psp.pokemon_id
		INNER JOIN pokemon_types AS pokemon_types1 ON pokemon_types1.id = pokemon.type_1_id
		LEFT JOIN pokemon_types AS pokemon_types2 ON pokemon_types2.id = pokemon.type_2_id
		INNER JOIN stocks ON stocks.id = psp.stock_id
		INNER JOIN seasons ON seasons.id = psp.season_id
		ORDER BY pokemon.id
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		utils.LogFailureError("Error querying PSPs", err)
	}
	defer rows.Close()

	var payloads []structs.PspElasticDocument

	for rows.Next() {
		queriedData, err := pgx.RowToMap(rows)
		if err != nil {
			utils.LogFailureError("Error reading individual row", err)
		}

		payload := convertDbRowToIndexPayload(queriedData)
		payloads = append(payloads, payload)
	}

	if err = rows.Err(); err != nil {
		utils.LogFailureError("Error reading queried rows", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(payloads))
	defer close(errChan)

	for _, payload := range payloads {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := indexPspPayload(elasticClient, payload)
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()

	select {
	case err := <-errChan:
		utils.LogFailureError("Error indexing PSPs into Elasticsearch", err)
	default:
		utils.LogSuccess("Successfully indexed PSPs into Elasticsearch")
	}
}
