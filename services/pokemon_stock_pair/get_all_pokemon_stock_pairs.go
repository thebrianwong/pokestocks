package pokemon_stock_pair

import (
	"context"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAllPokemonStockPairs(ctx context.Context, in *psp_pb.GetAllPokemonStockPairsRequest) (*psp_pb.GetAllPokemonStockPairsResponse, error) {
	db := s.DB

	query := pspQueryString() + "ORDER BY pokemon.id"

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while querying PSPs: %v", err)
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

	return &psp_pb.GetAllPokemonStockPairsResponse{Data: psps}, nil
}
