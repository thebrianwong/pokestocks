package main

import (
	"context"
	"log"
	"pokestocks/utils"

	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func main() {
	utils.LoadEnvVars("../../.env")
	elasticClient := utils.ConnectToElastic("../../http_ca.crt")

	_, err := elasticClient.Indices.Create("pokemon_stock_pairs_index").Request(
		&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"id": types.NewIntegerNumberProperty(),
					"pokemon": types.NestedProperty{
						Properties: map[string]types.Property{
							"id":             types.NewIntegerNumberProperty(),
							"name":           types.NewTextProperty(),
							"pokedex_number": types.NewKeywordProperty(),
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
		log.Fatalf("Error creating pokemon_stock_pairs_index: %v", err)
	}
	log.Println("Successfully created pokemon_stock_pairs_index")
}
