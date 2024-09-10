package pokemon_stock_pair

import (
	"context"
	"encoding/json"
	"fmt"
	"pokestocks/internal/structs"
	"strings"
	"sync"
	"time"

	common_pb "pokestocks/proto/common"

	psp_pb "pokestocks/proto/pokemon_stock_pair"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Server) SearchPokemonStockPairs(ctx context.Context, in *psp_pb.SearchPokemonStockPairsRequest) (*psp_pb.SearchPokemonStockPairsResponse, error) {
	db := s.DB

	searchValue := in.SearchValue

	searchResults, err := s.searchElasticIndex(searchValue)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error searching data: %v", err)
	}

	elasticPsps, err := convertPokemonStockPairElasticDocuments(searchResults)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error formatting data: %v", err)
	}

	ids := extractPokemonStockPairIds(elasticPsps)

	if len(ids) == 0 {
		return &psp_pb.SearchPokemonStockPairsResponse{Data: nil}, nil
	}

	query := `
		SELECT 
			psp.id AS "pspId",
			pokemon.id AS "pokemonId",
			pokemon."name" AS "pokemonName",
			pokemon.pokedex_number AS "pokedexNumber",
			pokemon.created_at AS "pokemonCreatedAt",
			pokemon.updated_at AS "pokemonUpdatedAt",
			pokemon.sprite_url as "pokemonSpriteUrl",
			pokemon_types1.id AS "type1Id",
			pokemon_types1."type" AS "type1Name",
			pokemon_types1.sprite_url AS "type1SpriteUrl",
			pokemon_types2.id AS "type2Id",
			pokemon_types2."type" AS "type2Name",
			pokemon_types2.sprite_url AS "type2SpriteUrl",
			stocks.id AS "stockId",
			stocks.symbol AS "stockSymbol", 
			stocks."name" AS "stockName",
			stocks.created_at AS "stockCreatedAt",
			stocks.updated_at AS "stockUpdatedAt",
			stocks.active AS "stockActive",
			seasons.id AS "seasonId",
			seasons."name" AS "seasonName",
			seasons.active AS "seasonActive"
		FROM pokemon_stock_pairs AS psp
		INNER JOIN pokemon ON pokemon.id = psp.pokemon_id
		INNER JOIN pokemon_types AS pokemon_types1 ON pokemon_types1.id = pokemon.type_1_id
		LEFT JOIN pokemon_types AS pokemon_types2 ON pokemon_types2.id = pokemon.type_2_id
		INNER JOIN stocks ON stocks.id = psp.stock_id
		INNER JOIN seasons ON seasons.id = psp.season_id
	`

	queryArgs := []any{}
	positionalParams := []string{}

	for i, id := range ids {
		queryArgs = append(queryArgs, id)
		positionalParams = append(positionalParams, fmt.Sprintf("$%d", i+1))
	}

	orderByArgs := strings.Join(ids, ",")
	queryArgs = append(queryArgs, orderByArgs)

	positionalParamsString := strings.Join(positionalParams, ", ")
	query += fmt.Sprintf("WHERE psp.id IN (%s)", positionalParamsString)
	query += fmt.Sprintf("ORDER BY POSITION(psp.id::text IN $%d)", len(positionalParams)+1)

	rows, err := db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while querying for %v: %v", searchValue, err)
	}
	defer rows.Close()

	var psps []*common_pb.PokemonStockPair

	for rows.Next() {
		queriedData, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting queried data to map: %v", err)
		}

		psp := convertDbRowToPokemonStockPair(queriedData)
		psps = append(psps, psp)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading queried data: %v", err)
	}

	err = enrichWithStockPrices(psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)

	}
	return &psp_pb.SearchPokemonStockPairsResponse{Data: psps}, nil
}
