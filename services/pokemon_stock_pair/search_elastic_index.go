package pokemon_stock_pair

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func (s *Server) searchElasticIndex(searchValue string) (*search.Response, error) {
	two := float32(2.0)
	three := float32(3.0)
	sevenPointFive := float32(7.5)
	fifteen := float32(15.0)
	twenty := float32(20.0)

	elasticClient := s.ElasticClient

	pokemonNestedQuery := types.Query{
		Nested: &types.NestedQuery{
			Path: "pokemon",
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Should: []types.Query{
						{
							Match: map[string]types.MatchQuery{
								"pokemon.name": {
									Query:     searchValue,
									Boost:     &fifteen,
									Fuzziness: 1,
								},
							},
						},
						{
							Match: map[string]types.MatchQuery{
								"pokemon.name.ngram": {
									Query:     searchValue,
									Boost:     &two,
									Fuzziness: 1,
								},
							},
						},
						{
							Prefix: map[string]types.PrefixQuery{
								"pokemon.name": {
									Value: searchValue,
									Boost: &three,
								},
							},
						},
						{
							Match: map[string]types.MatchQuery{
								"pokemon.pokedex_number": {
									Query: searchValue,
									Boost: &twenty,
								},
							},
						},
						{
							Match: map[string]types.MatchQuery{
								"pokemon.type_1.text": {
									Query: searchValue,
									Boost: &sevenPointFive,
								},
							},
						},
						{
							Match: map[string]types.MatchQuery{
								"pokemon.type_2.text": {
									Query: searchValue,
									Boost: &sevenPointFive,
								},
							},
						},
					},
				},
			},
		},
	}

	stockNestedQuery := types.Query{
		Nested: &types.NestedQuery{
			Path: "stock",
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Should: []types.Query{
						{
							Match: map[string]types.MatchQuery{
								"stock.name": {
									Query:     searchValue,
									Boost:     &fifteen,
									Fuzziness: 1,
								},
							},
						},
						{
							Match: map[string]types.MatchQuery{
								"stock.name.ngram": {
									Query:     searchValue,
									Boost:     &two,
									Fuzziness: 1,
								},
							},
						},
						{
							Prefix: map[string]types.PrefixQuery{
								"stock.name": {
									Value: searchValue,
									Boost: &three,
								},
							},
						},
						{
							Prefix: map[string]types.PrefixQuery{
								"stock.name.full_name": {
									Value: searchValue,
									Boost: &fifteen,
								},
							},
						},
						{
							Match: map[string]types.MatchQuery{
								"stock.symbol": {
									Query: strings.ToUpper(searchValue),
									Boost: &twenty,
								},
							},
						},
					},
				},
			},
		},
	}

	activeStockFilter := types.Query{
		Nested: &types.NestedQuery{
			Path: "stock",
			Query: &types.Query{
				Term: map[string]types.TermQuery{
					"stock.active": {
						Value: true,
					},
				},
			},
		},
	}

	activeSeasonFilter := types.Query{
		Term: map[string]types.TermQuery{
			"active_season": {
				Value: true,
			},
		},
	}

	res, err := elasticClient.Search().Index("pokemon_stock_pairs_index").Request(
		&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Should: []types.Query{
						pokemonNestedQuery,
						stockNestedQuery,
					},
					Must: []types.Query{
						activeStockFilter,
						activeSeasonFilter,
					},
				},
			},
		},
	).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}
