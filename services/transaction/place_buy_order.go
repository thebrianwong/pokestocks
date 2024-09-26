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

	portfolioId := in.PortfolioId
	pspId := []string{fmt.Sprint(in.PspId)}
	quantity := in.Quantity

	if quantity < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid quantity argument")
	}

	psps, err := cm.QueryPokemonStockPairs(ctx, pspId)
	if err != nil {
		errorCode := codes.Internal
		if err.Error() == "requested PSP does not exist" {
			errorCode = codes.NotFound
		}
		return nil, status.Errorf(errorCode, "error querying PSPs: %v", err)
	}

	err = cm.EnrichWithStockPrices(ctx, psps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying Alpaca for price data: %v", err)
	}

	psp := psps[0]
	stockPrice := *psp.Stock.Price

	totalPrice := stockPrice * float64(quantity)

	cash, err := cm.QueryPortfolioCash(ctx, portfolioId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, status.Error(codes.NotFound, "error querying portfolio cash: requested portfolio does not exist")
		}
		return nil, status.Errorf(codes.Internal, "error querying portfolio cash: %v", err)
	}

	hasSufficientCash := cash >= totalPrice

	if !hasSufficientCash {
		return nil, status.Error(codes.FailedPrecondition, "the portfolio contains insufficient cash to place this buy order")
	}

	err = cm.TransactBuyOrder(ctx, portfolioId, psp.Id, quantity, stockPrice)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error completing buy order transaction: %v", err)
	}

	return &transaction_pb.PlaceBuyOrderResponse{
		Message: fmt.Sprintf(
			`Success buying %d shares of %s corresponding to PSP id %d and Pokemon %s for %f a share. Total Price: %f. Current Cash: %f`, quantity, psp.Stock.Symbol, psp.Id, psp.Pokemon.Name, stockPrice, totalPrice, (cash - totalPrice),
		),
	}, nil
}
