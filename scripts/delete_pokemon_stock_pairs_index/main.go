package main

import (
	"context"
	"pokestocks/utils"
)

func main() {
	utils.LoadEnvVars("../../.env")
	elasticClient := utils.ConnectToElastic("../../http_ca.crt")

	_, err := elasticClient.Indices.Delete("pokemon_stock_pairs_index").Do(context.TODO())

	if err != nil {
		utils.LogFailureError("Error deleting pokemon_stock_pairs_index", err)
	}
	utils.LogSuccess("Successfully deleted pokemon_stock_pairs_index")
}
