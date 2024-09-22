package main

import (
	"context"
	"os"
	"pokestocks/utils"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
)

func getPspIndexDocCount(elasticClient *elasticsearch.TypedClient) (*(count.Response), error) {
	count, err := elasticClient.Count().Index("pokemon_stock_pairs_index").Do(context.Background())

	if err != nil {
		return nil, err
	}

	return count, nil
}

func bulkAddFailureCallback(a context.Context, b esutil.BulkIndexerItem, c esutil.BulkIndexerResponseItem, err error) {
	utils.LogFailureError("Aborting from bulkAddFailureCallback()", err)
}

func main() {
	utils.LoadEnvVars("../../.env")
	conn := utils.ConnectToDb()
	typedElasticClient := utils.CreateTypedElasticClient("../../")
	regularElasticClient := utils.CreateRegularElasticClient("../../")

	count, err := getPspIndexDocCount(typedElasticClient)
	if err != nil {
		utils.LogFailureError("Error starting PSP indexing", err)
	}
	if count.Count == 1025 {
		utils.LogSuccess("Exiting as index already contains expected PSPs")
		os.Exit(0)
	}
	if count.Count > 0 {
		utils.LogFailure("Error starting PSP indexing: the index contains an unexpected number of PSPs")
	}

	query :=
		`
			WITH tab1 AS (
				SELECT
					pokemon_table.id,
					pokemon_table.name,
					pokemon_table.pokedex_number,
					pokemon_types1.type AS type_1,
					pokemon_types2.type AS type_2
				FROM pokemon AS pokemon_table
				INNER JOIN pokemon_types AS pokemon_types1 ON pokemon_types1.id = pokemon_table.type_1_id
				LEFT JOIN pokemon_types AS pokemon_types2 ON pokemon_types2.id = pokemon_table.type_2_id
			),
			tab2 AS (
				SELECT 
					stocks.id,
					stocks.symbol,
					stocks.name,
					stocks.active 
				FROM stocks
				INNER JOIN pokemon_stock_pairs AS psp ON psp.stock_id = stocks.id
				INNER JOIN tab1 ON tab1.id = psp.pokemon_id
			)
			SELECT JSON_BUILD_OBJECT(
				'id', psp.id,
				'pokemon', tab1.*, 
				'stock', tab2.*, 
				'active_season', seasons.active
			) FROM tab1
			INNER JOIN pokemon_stock_pairs AS psp ON psp.pokemon_id = tab1.id
			INNER JOIN tab2 ON psp.stock_id = tab2.id
			INNER JOIN seasons ON psp.season_id = seasons.id
			ORDER BY tab1.pokedex_number
		`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		utils.LogFailureError("Error querying PSPs", err)
	}
	defer rows.Close()

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "pokemon_stock_pairs_index",
		Client: regularElasticClient,
	})
	if err != nil {
		utils.LogFailureError("Error creating Elasticsearch bulk inserter", err)
	}

	for rows.Next() {
		jsonData := string(rows.RawValues()[0])

		err = bulkIndexer.Add(context.Background(),
			esutil.BulkIndexerItem{
				Action:    "index",
				Index:     "pokemon_stock_pairs_index",
				Body:      strings.NewReader(jsonData),
				OnFailure: bulkAddFailureCallback,
			},
		)
		if err != nil {
			utils.LogFailureError("Error queuing up bulk insert", err)
		}
	}

	if err = rows.Err(); err != nil {
		utils.LogFailureError("Error reading queried rows", err)
	}

	err = bulkIndexer.Close(context.Background())
	if err != nil {
		utils.LogFailureError("Error closing bulkIndexer", err)
	}

	bulkIndexFailed := bulkIndexer.Stats().NumFailed != 0
	if bulkIndexFailed {
		utils.LogFailure("Some PSPs failed to be properly indexed")
	}

	utils.LogSuccess("Successfully indexed PSPs into Elasticsearch")
}
