package pokemon_stock_pair

import (
	"context"
	"fmt"
	"pokestocks/internal/helpers"
	psp_pb "pokestocks/proto/pokemon_stock_pair"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRandomPokemonStockPairs(ctx context.Context, in *psp_pb.GetRandomPokemonStockPairsRequest) (*psp_pb.GetRandomPokemonStockPairsResponse, error) {
	start := time.Now()
	cm := s.ClientManager
	db := cm.DB

	idQuery := `
		SELECT id
		FROM pokemon_stock_pairs
	`

	idsRows, err := db.Query(context.Background(), idQuery)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while querying PSP ids: %v", err)
	}
	defer idsRows.Close()

	var ids []string

	for idsRows.Next() {
		var id string
		err = idsRows.Scan(&id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error converting queried PSP ids: %v", err)
		}

		ids = append(ids, id)
	}

	if err = idsRows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "error reading queried PSP ids: %v", err)
	}

	randomIndices := helpers.GenerateRandomIndices(5, len(ids))

	var pspIds []string

	for _, randomIndex := range randomIndices {
		pspIds = append(pspIds, fmt.Sprint(randomIndex))
	}

	psps, err := cm.QueryPokemonStockPairs(ctx, pspIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying PSPs: %v", err)
	}

	err = cm.EnrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	fmt.Println("start", time.Since(start))

	return &psp_pb.GetRandomPokemonStockPairsResponse{Data: psps}, nil
}
