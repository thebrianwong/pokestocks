package client_manager

import (
	"context"
	"pokestocks/internal/helpers"
	redis_keys "pokestocks/redis"
	"pokestocks/utils"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/redis/go-redis/v9"
)

func (cc *ClientManager) searchElasticIndex(searchValue string) (*search.Response, error) {
	zeroPointOne := types.Float64(0.1)
	two := float32(2.0)
	three := float32(3.0)
	sevenPointFive := float32(7.5)
	fifteen := float32(15.0)
	twenty := float32(20.0)

	elasticClient := cc.ElasticClient

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
			MinScore: &zeroPointOne,
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Should: []types.Query{
						pokemonNestedQuery,
						stockNestedQuery,
					},
					Filter: []types.Query{
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

func (cc *ClientManager) SearchPokemonStockPairIds(ctx context.Context, searchValue string) ([]string, error) {
	redisClient := cc.RedisClient

	var pspIds []string

	cachedElasticIds, err := redisClient.ZRange(ctx, redis_keys.ElasticCacheKey(searchValue), 0, -1).Result()
	if err == nil && len(cachedElasticIds) != 0 {
		pspIds = cachedElasticIds
	} else {
		if err != nil {
			// if there is something wrong with Redis and it can't answer our request,
			// we can always just fallback to searching Elastic
			utils.LogWarningError("Error querying Redis key "+redis_keys.ElasticCacheKey(searchValue)+" for cached PSP ids. Falling back to Elastic", err)
		}
		searchResults, err := cc.searchElasticIndex(searchValue)
		if err != nil {
			return nil, err
		}

		elasticPsps, err := helpers.ConvertPokemonStockPairElasticDocuments(searchResults)
		if err != nil {
			return nil, err
		}
		if len(elasticPsps) == 0 {
			return nil, nil
		}

		pspIds = helpers.ExtractPokemonStockPairIds(elasticPsps)

		redisPipeline := redisClient.Pipeline()
		sortedSet := []redis.Z{}
		for i, id := range pspIds {
			sortedSetMember := redis.Z{
				Score:  float64(i),
				Member: id,
			}
			sortedSet = append(sortedSet, sortedSetMember)
		}
		midnightTomorrow := helpers.MidnightTomorrow()
		redisPipeline.ZAdd(ctx, redis_keys.ElasticCacheKey(searchValue), sortedSet...)
		redisPipeline.ExpireAt(ctx, redis_keys.ElasticCacheKey(searchValue), midnightTomorrow)

		_, err = redisPipeline.Exec(ctx)
		if err != nil {
			// don't return a gRPC response with an error
			// a response with data can still be generated even if we can't cache Elasticsearch results
			utils.LogWarningError("Error caching data to Redis for key "+redis_keys.ElasticCacheKey(searchValue)+". Skipping", err)
		}
	}

	return pspIds, nil
}
