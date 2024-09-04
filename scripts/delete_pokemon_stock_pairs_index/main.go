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
		log.Fatalf("Error deleting pokemon_stock_pairs_index: %v", err)
	}
	log.Println("Successfully deleted pokemon_stock_pairs_index")
}
