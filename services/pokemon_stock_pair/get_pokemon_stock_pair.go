package pokemon_stock_pair

import (
	"context"
	"fmt"
	psp_pb "pokestocks/proto/pokemon_stock_pair"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetPokemonStockPair(ctx context.Context, in *psp_pb.GetPokemonStockPairRequest) (*psp_pb.GetPokemonStockPairResponse, error) {
	cm := s.ClientManager

	argumentProvided := in.ProtoReflect().Has(in.ProtoReflect().Descriptor().Fields().ByName("id"))

	if !argumentProvided {
		return nil, status.Errorf(codes.InvalidArgument, "id argument not provided")
	}

	pspId := in.Id

	psp, err := s.queryPokemonStockPairs(ctx, []string{fmt.Sprint(pspId)})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying PSPs: %v", err)
	}

	err = s.enrichWithStockPrices(ctx, psp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	return &psp_pb.GetPokemonStockPairResponse{Data: psp}, nil
}
