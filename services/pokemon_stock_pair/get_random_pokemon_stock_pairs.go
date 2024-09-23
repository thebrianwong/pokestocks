package pokemon_stock_pair

import (
	"context"
	"fmt"
	common_pb "pokestocks/proto/common"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	"time"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRandomPokemonStockPairs(ctx context.Context, in *psp_pb.GetRandomPokemonStockPairsRequest) (*psp_pb.GetRandomPokemonStockPairsResponse, error) {
	start := time.Now()
	db := s.DB

	idQuery := `
		SELECT id
		FROM pokemon_stock_pairs
	`

	idsRows, err := db.Query(context.Background(), idQuery)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while querying PSP ids: %v", err)
	}
	defer idsRows.Close()

	var ids []any

	for idsRows.Next() {
		var id int
		err = idsRows.Scan(&id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting queried PSP ids: %v", err)
		}

		ids = append(ids, id)
	}

	if err = idsRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading queried PSP ids: %v", err)
	}

	randomIndices := generateRandomIndices(5, len(ids))

	var pspIds []any

	for _, randomIndex := range randomIndices {
		pspIds = append(pspIds, randomIndex)
	}

	query := pspQueryString() + "WHERE psp.id IN ($1, $2, $3, $4, $5)"
	rows, err := db.Query(context.Background(), query, pspIds...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while querying PSPs: %v", err)
	}
	defer rows.Close()

	var psps []*common_pb.PokemonStockPair

	for rows.Next() {
		queriedData, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting queried PSPs to maps: %v", err)
		}

		psp := convertDbRowToPokemonStockPair(queriedData)
		psps = append(psps, psp)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading queried PSPs: %v", err)
	}

	err = s.enrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	fmt.Println("start", time.Since(start))

	return &psp_pb.GetRandomPokemonStockPairsResponse{Data: psps}, nil
}
