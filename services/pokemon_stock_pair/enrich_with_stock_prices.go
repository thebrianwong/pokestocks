package pokemon_stock_pair

import (
	common_pb "pokestocks/proto/common"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

func (s *Server) enrichWithStockPrices(psps []*common_pb.PokemonStockPair) error {
	alpacaClient := s.AlpacaMarketDataClient

	symbols := []string{}

	for _, psp := range psps {
		symbols = append(symbols, psp.Stock.Symbol)
	}

	requestParams := marketdata.GetLatestTradeRequest{}
	data, err := alpacaClient.GetLatestTrades(symbols, requestParams)
	if err != nil {
		return err
	}

	for _, psp := range psps {
		priceData := data[psp.Stock.Symbol]
		psp.Stock.Price = &priceData.Price
	}

	return nil
}
