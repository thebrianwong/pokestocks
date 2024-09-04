package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"pokestocks/utils"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func main() {
	// refactor this into a util function
	// ===
	utils.LoadEnvVars("../../.env")

	elasticUsername := os.Getenv("ELASTIC_USERNAME")
	elasticPassword := os.Getenv("ELASTIC_PASSWORD")
	// elasticApiKey := os.Getenv("ELASTIC_API_KEY")
	elasticEndpoint := os.Getenv(("ELASTIC_ENDPOINT"))
	cert, err := os.ReadFile("../../http_ca.crt")
	if err != nil {
		fmt.Println(err)
	}
	elasticClient, _ := elasticsearch.NewTypedClient(elasticsearch.Config{
		// APIKey:    elasticApiKey,
		Addresses: []string{elasticEndpoint},
		Username:  elasticUsername,
		Password:  elasticPassword,
		CACert:    cert,
	})
	// ===

	_, err = elasticClient.Indices.Create("pokemon_stock_pairs_index").Request(
		&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"id": types.NewIntegerNumberProperty(),
					"pokemon": types.NestedProperty{
						Properties: map[string]types.Property{
							"id":             types.NewIntegerNumberProperty(),
							"name":           types.NewTextProperty(),
							"pokedex_number": types.NewIntegerNumberProperty(),
							"created_at":     types.NewDateProperty(),
							"updated_at":     types.NewDateProperty(),
							"sprite_url":     types.NewKeywordProperty(),
							"type_1": types.NestedProperty{
								Properties: map[string]types.Property{
									"id":         types.NewIntegerNumberProperty(),
									"type":       types.NewKeywordProperty(),
									"sprite_url": types.NewKeywordProperty(),
								},
							},
							"type_2": types.NestedProperty{
								Properties: map[string]types.Property{
									"id":         types.NewIntegerNumberProperty(),
									"type":       types.NewKeywordProperty(),
									"sprite_url": types.NewKeywordProperty(),
								},
							},
						},
					},
					"stock": types.NestedProperty{
						Properties: map[string]types.Property{
							"id":         types.NewIntegerNumberProperty(),
							"symbol":     types.NewKeywordProperty(),
							"name":       types.NewTextProperty(),
							"created_at": types.NewDateProperty(),
							"updated_at": types.NewDateProperty(),
							"active":     types.NewBooleanProperty(),
						},
					},
					"season": types.NestedProperty{
						Properties: map[string]types.Property{
							"id":     types.NewIntegerNumberProperty(),
							"name":   types.NewKeywordProperty(),
							"active": types.NewBooleanProperty(),
						},
					},
				},
			},
		},
	).Do(context.TODO())

	if err != nil {
		log.Fatalf("Error creating pokemon_stock_pair_index: %v", err)
	}
}
