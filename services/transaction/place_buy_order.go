package transaction

import (
	"context"
	"fmt"
	transaction_pb "pokestocks/proto/transaction"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) PlaceBuyOrder(ctx context.Context, in *transaction_pb.PlaceBuyOrderRequest) (*transaction_pb.PlaceBuyOrderResponse, error) {
	cm := s.ClientManager

	argumentsProvided :=
		in.ProtoReflect().Has(in.ProtoReflect().Descriptor().Fields().ByName("portfolioId")) &&
			in.ProtoReflect().Has(in.ProtoReflect().Descriptor().Fields().ByName("pspId")) &&
			in.ProtoReflect().Has(in.ProtoReflect().Descriptor().Fields().ByName("quantity"))

	if !argumentsProvided {
		return nil, status.Errorf(codes.InvalidArgument, "missing arguments")
	}

	stockQuantity := float64(in.Quantity)

	psps, err := cm.QueryPokemonStockPairs(ctx, []string{fmt.Sprint(in.PspId)})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying PSPs: %v", err)
	}

	err = cm.EnrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	psp := psps[0]
	stockPrice := *psp.Stock.Price

	totalPrice := stockPrice * stockQuantity

	cash, err := cm.QueryPortfolioCash(ctx, in.PortfolioId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying portfolio cash: %v", err)
	}

	hasSufficientCash := cash >= totalPrice

	return &transaction_pb.PlaceBuyOrderResponse{Message: fmt.Sprint(hasSufficientCash)}, nil
}
