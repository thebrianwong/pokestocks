package main

import (
	"context"
	"log"
	"pokestocks/utils"
)

func main() {
	utils.LoadEnvVars("../../.env")
	elasticClient := utils.ConnectToElastic("../../http_ca.crt")

	_, err := elasticClient.Indices.Delete("pokemon_stock_pairs_index").Do(context.TODO())

	if err != nil {
		utils.LogFailureError("Error deleting pokemon_stock_pairs_index", err)
	}
	log.Println("Successfully deleted pokemon_stock_pairs_index")
}
