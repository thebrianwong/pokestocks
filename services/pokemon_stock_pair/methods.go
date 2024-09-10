package pokemon_stock_pair

import (
	"context"
	"encoding/json"
	"fmt"
	"pokestocks/internal/structs"
	"sync"
	"time"

	common_pb "pokestocks/proto/common"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertDbRowToPokemonStockPair(rowDataMap map[string]any) *common_pb.PokemonStockPair {
	psp := common_pb.PokemonStockPair{
		Id: rowDataMap["pspId"].(int64),
		Pokemon: &common_pb.Pokemon{
			Id:            rowDataMap["pokemonId"].(int64),
			Name:          rowDataMap["pokemonName"].(string),
			PokedexNumber: rowDataMap["pokedexNumber"].(int32),
			CreatedAt:     timestamppb.New(rowDataMap["pokemonCreatedAt"].(time.Time)),
			UpdatedAt:     timestamppb.New(rowDataMap["pokemonUpdatedAt"].(time.Time)),
			Type1: &common_pb.PokemonType{
				Id:        rowDataMap["type1Id"].(int64),
				Type:      rowDataMap["type1Name"].(string),
				SpriteUrl: rowDataMap["type1SpriteUrl"].(string),
			},
			SpriteUrl: rowDataMap["pokemonSpriteUrl"].(string),
		},
		Stock: &common_pb.Stock{
			Id:        rowDataMap["stockId"].(int64),
			Symbol:    rowDataMap["stockSymbol"].(string),
			Name:      rowDataMap["stockName"].(string),
			CreatedAt: timestamppb.New(rowDataMap["stockCreatedAt"].(time.Time)),
			UpdatedAt: timestamppb.New(rowDataMap["stockUpdatedAt"].(time.Time)),
			Active:    rowDataMap["stockActive"].(bool),
		},
		Season: &common_pb.Season{
			Id:     rowDataMap["seasonId"].(int64),
			Name:   rowDataMap["seasonName"].(string),
			Active: rowDataMap["seasonActive"].(bool),
		},
	}

	// not all Pokemon have a second type
	type2Id, ok := rowDataMap["type2Id"].(int64)
	if ok {
		psp.Pokemon.Type2 = &common_pb.PokemonType{
			Id:        type2Id,
			Type:      rowDataMap["type2Name"].(string),
			SpriteUrl: rowDataMap["type2SpriteUrl"].(string),
		}
	}

	return &psp
}

func getStockPrice(symbol string) (float64, error) {
	client := marketdata.NewClient(marketdata.ClientOpts{})

	requestParams := marketdata.GetLatestTradeRequest{}

	data, err := client.GetLatestTrade(symbol, requestParams)
	if err != nil {
		return 0, err
	}

	stockPrice := data.Price
	return stockPrice, nil
}

func enrichWithStockPrices(psps []*common_pb.PokemonStockPair) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(psps))
	defer close(errChan)

	for _, psp := range psps {
		wg.Add(1)
		go func() {
			defer wg.Done()

			stockPrice, err := getStockPrice(psp.Stock.Symbol)
			if err != nil {
				errChan <- err
			}
			psp.Stock.Price = &stockPrice
		}()
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (s *Server) searchElasticIndex(searchValue string) (*search.Response, error) {
	two := float32(2.0)
	three := float32(3.0)
	fifteen := float32(15.0)

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
					},
				},
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
				},
			},
		},
	).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func convertPokemonStockPairElasticDocuments(elasticResponse *search.Response) ([]structs.PspElasticDocument, error) {
	var convertedDocuments []structs.PspElasticDocument

	docs := elasticResponse.Hits.Hits
	for _, doc := range docs {
		var convertedDocument structs.PspElasticDocument
		docData := doc.Source_
		err := json.Unmarshal(docData, &convertedDocument)
		if err != nil {
			return nil, err
		}
		convertedDocuments = append(convertedDocuments, convertedDocument)
	}

	return convertedDocuments, nil
}

func extractPokemonStockPairIds(docs []structs.PspElasticDocument) []string {
	var ids []string

	for _, doc := range docs {
		id := doc.Id
		ids = append(ids, fmt.Sprint(id))
	}

	return ids
}
