package pokemon_stock_pair

import (
	"context"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	"strconv"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetPokemonStockPair(ctx context.Context, in *psp_pb.GetPokemonStockPairRequest) (*psp_pb.GetPokemonStockPairResponse, error) {
	db := s.DB

	query := pspQueryString()

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

	err = s.enrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)

	}
	return &psp_pb.GetPokemonStockPairResponse{Data: psps}, nil
}
