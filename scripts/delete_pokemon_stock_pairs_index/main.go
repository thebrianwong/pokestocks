package main

import (
	"context"
	"os"
	"pokestocks/utils"
)

func main() {
	utils.LoadEnvVars("../../.env")
	elasticClient := utils.CreateTypedElasticClient("../../")

	indexExists, err := elasticClient.Indices.Exists("pokemon_stock_pairs_index").Do(context.Background())
	if err != nil {
		utils.LogFailureError("Error checking if pokemon_stock_pairs_index exists", err)
	}
	if !indexExists {
		utils.LogSuccess("Exiting as pokemon_stock_pairs_index does not exist")
		os.Exit(0)
	}

	_, err = elasticClient.Indices.Delete("pokemon_stock_pairs_index").Do(context.Background())

	if err != nil {
		utils.LogFailureError("Error deleting pokemon_stock_pairs_index", err)
	}
	utils.LogSuccess("Successfully deleted pokemon_stock_pairs_index")
}
