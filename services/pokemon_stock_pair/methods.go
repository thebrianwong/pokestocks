package pokemon_stock_pair

import (
	"context"
	"log"
	"strconv"
	"time"

	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"

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

func (s *Server) GetAllPokemonStockPairs(ctx context.Context, in *psp_pb.GetAllPokemonStockPairsRequest) (*psp_pb.GetAllPokemonStockPairsResponse, error) {
	db := s.DB

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
		ORDER BY pokemon.id
	`
	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Fatalf("Error querying PSPs: %v", err)
	}
	defer rows.Close()

	var psps []*common_pb.PokemonStockPair

	for rows.Next() {
		queriedData, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, err
		}

		psp := convertDbRowToPokemonStockPair(queriedData)
		psps = append(psps, psp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &psp_pb.GetAllPokemonStockPairsResponse{Data: psps}, nil
}

func (s *Server) GetPokemonStockPair(ctx context.Context, in *psp_pb.GetPokemonStockPairRequest) (*psp_pb.GetPokemonStockPairResponse, error) {
	db := s.DB

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

	var searchValue string

	switch input := in.Input.(type) {
	case *psp_pb.GetPokemonStockPairRequest_PokemonName:
		searchValue = input.PokemonName
		query = query + `WHERE LOWER(pokemon."name") = LOWER($1)`
	case *psp_pb.GetPokemonStockPairRequest_PokedexNumber:
		searchValue = strconv.FormatInt(int64(input.PokedexNumber), 10)
		query = query + "WHERE pokemon.pokedex_number = $1"
	case *psp_pb.GetPokemonStockPairRequest_StockName:
		searchValue = "%" + input.StockName + "%"
		query = query + `WHERE LOWER(stocks."name") LIKE LOWER($1)`
	case *psp_pb.GetPokemonStockPairRequest_StockSymbol:
		searchValue = input.StockSymbol
		query = query + "WHERE LOWER(stocks.symbol) = LOWER($1)"
	default:
		return nil, status.Errorf(codes.InvalidArgument, "input not provided")
	}

	rows, err := db.Query(ctx, query, searchValue)
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

	return &psp_pb.GetPokemonStockPairResponse{Data: psps}, nil
}
