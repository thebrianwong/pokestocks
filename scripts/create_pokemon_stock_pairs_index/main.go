package main

import (
	"context"
	"pokestocks/utils"

	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func main() {
	utils.LoadEnvVars("../../.env")
	elasticClient := utils.CreateTypedElasticClient("../../http_ca.crt")

	three := 3
	twelve := 12
	fifteen := 15
	ngram_analyzer := "ngram_analyzer"
	full_name_analyzer := "full_name_analyzer"

	_, err := elasticClient.Indices.Create("pokemon_stock_pairs_index").Request(
		&create.Request{
			Settings: &types.IndexSettings{
				Analysis: &types.IndexSettingsAnalysis{
					Analyzer: map[string]types.Analyzer{
						ngram_analyzer: types.CustomAnalyzer{
							Tokenizer: "standard",
							Filter:    []string{"lowercase", "ngram_filter"},
						},
						full_name_analyzer: types.CustomAnalyzer{
							Tokenizer: "keyword",
							Filter:    []string{"lowercase"},
						},
					},
					Filter: map[string]types.TokenFilter{
						"ngram_filter": &types.NGramTokenFilter{
							Type:    "ngram",
							MinGram: &three,
							MaxGram: &fifteen,
						},
					},
				},
				MaxNgramDiff: &twelve,
			},
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"id": types.NewIntegerNumberProperty(),
					"pokemon": types.NestedProperty{
						Properties: map[string]types.Property{
							"id": types.NewIntegerNumberProperty(),
							"name": types.TextProperty{
								Type: "text",
								Fields: map[string]types.Property{
									"ngram": types.TextProperty{
										Type:     "text",
										Analyzer: &ngram_analyzer,
									},
								},
							},
							"pokedex_number": types.NewKeywordProperty(),
							// "created_at":     types.NewDateProperty(),
							// "updated_at":     types.NewDateProperty(),
							// "sprite_url":     types.NewKeywordProperty(),
							"type_1": types.NewKeywordProperty(),
							"type_2": types.NewKeywordProperty(),
							// "type_1": types.NestedProperty{
							// 	Properties: map[string]types.Property{
							// 		"id":         types.NewIntegerNumberProperty(),
							// 		"type":       types.NewKeywordProperty(),
							// 		"sprite_url": types.NewKeywordProperty(),
							// 	},
							// },
							// "type_2": types.NestedProperty{
							// 	Properties: map[string]types.Property{
							// 		"id":         types.NewIntegerNumberProperty(),
							// 		"type":       types.NewKeywordProperty(),
							// 		"sprite_url": types.NewKeywordProperty(),
							// 	},
							// },
						},
					},
					"stock": types.NestedProperty{
						Properties: map[string]types.Property{
							"id":     types.NewIntegerNumberProperty(),
							"symbol": types.NewKeywordProperty(),
							"name": types.TextProperty{
								Type: "text",
								Fields: map[string]types.Property{
									"ngram": types.TextProperty{
										Type:     "text",
										Analyzer: &ngram_analyzer,
									},
									"full_name": types.TextProperty{
										Type:     "keyword",
										Analyzer: &full_name_analyzer,
									},
								},
							},
							// "created_at": types.NewDateProperty(),
							// "updated_at": types.NewDateProperty(),
							"active": types.NewBooleanProperty(),
						},
					},
					"active_season": types.NewBooleanProperty(),
					// "season": types.NestedProperty{
					// 	Properties: map[string]types.Property{
					// 		"id":     types.NewIntegerNumberProperty(),
					// 		"name":   types.NewKeywordProperty(),
					// 		"active": types.NewBooleanProperty(),
					// 	},
					// },
				},
			},
		},
	).Do(context.Background())

	if err != nil {
		utils.LogFailureError("Error creating pokemon_stock_pairs_index", err)
	}
	utils.LogSuccess("Successfully created pokemon_stock_pairs_index")
}
