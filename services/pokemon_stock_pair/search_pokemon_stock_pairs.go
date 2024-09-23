package pokemon_stock_pair

import (
	"context"
	psp_pb "pokestocks/proto/pokemon_stock_pair"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SearchPokemonStockPairs(ctx context.Context, in *psp_pb.SearchPokemonStockPairsRequest) (*psp_pb.SearchPokemonStockPairsResponse, error) {
	// startSearchPokemonStockPairs := time.Now()

	argumentProvided := in.ProtoReflect().Has(in.ProtoReflect().Descriptor().Fields().ByName("searchValue"))

	if !argumentProvided {
		return nil, status.Errorf(codes.InvalidArgument, "searchValue argument not provided")
	}

	searchValue := in.SearchValue

	pspIds, err := s.searchPokemonStockPairIds(ctx, searchValue)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error searching for PSP ids: %v", err)
	}
	if pspIds == nil {
		return &psp_pb.SearchPokemonStockPairsResponse{Data: nil}, nil
	}

	psps, err := s.queryPokemonStockPairs(ctx, pspIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying PSPs: %v", err)
	}

	err = s.enrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	// fmt.Println("startSearchPokemonStockPairs took:", time.Since(startSearchPokemonStockPairs))
	// fmt.Println("======")

	return &psp_pb.SearchPokemonStockPairsResponse{Data: psps}, nil
}
