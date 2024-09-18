package pokemon_stock_pair

import (
	"context"
	"fmt"
	common_pb "pokestocks/proto/common"
	redis_keys "pokestocks/redis"
	"pokestocks/utils"
	"strconv"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/redis/go-redis/v9"
)

func (s *Server) enrichWithStockPrices(ctx context.Context, psps []*common_pb.PokemonStockPair) error {
	alpacaClient := s.AlpacaMarketDataClient
	redisClient := s.RedisClient
	redisPipeline := redisClient.Pipeline()

	marketIsOpen, err := s.isMarketOpen(ctx)
	if err != nil {
		utils.LogWarning(fmt.Sprintf("Error checking if market is open: %v", err))
	}

	symbols := []string{}
	cachedSymbolsMap := map[string]*float64{}
	nonCachedSymbolsArr := []string{}

	for _, psp := range psps {
		symbols = append(symbols, psp.Stock.Symbol)
	}

	if !marketIsOpen {
		for _, symbol := range symbols {
			redisPipeline.Get(ctx, redis_keys.StockSymbolKey(symbol))
		}
		results, err := redisPipeline.Exec(ctx)

		// if err is due to cache misses, that's ok and we can query Alpaca for the stocks that missed
		// if err is due to any other problem (network, redis down), we should query Alpaca for every stock
		if err == nil || err == redis.Nil {
			for _, result := range results {
				key := result.(*redis.StringCmd).Args()[1]
				stockSymbol := redis_keys.GetSymbolFromStockSymbolKey(key.(string))

				rawValue := result.(*redis.StringCmd).Val()
				if rawValue == "" {
					nonCachedSymbolsArr = append(nonCachedSymbolsArr, stockSymbol)
					continue
				}

				stockPrice, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					nonCachedSymbolsArr = append(nonCachedSymbolsArr, stockSymbol)
					continue
				}

				cachedSymbolsMap[stockSymbol] = &stockPrice
			}

			symbols = nonCachedSymbolsArr
		} else {
			utils.LogWarning(fmt.Sprintf("Error checking if cached data exists. Falling back to hitting the Alpaca API for all stock prices: %v", err))
		}
	}

	nonCachedData := map[string]marketdata.Trade{}

	// skip querying Alpaca if there are no cache misses
	if len(symbols) > 0 {
		fmt.Println("symbols to get from alpaca", symbols)
		fmt.Println("hitting alpaca")
		requestParams := marketdata.GetLatestTradeRequest{}
		nonCachedData, err = alpacaClient.GetLatestTrades(symbols, requestParams)
		if err != nil {
			return err
		}
	}

	for _, psp := range psps {
		priceData, ok := nonCachedData[psp.Stock.Symbol]
		if ok {
			psp.Stock.Price = &priceData.Price
			redisPipeline.Set(ctx, redis_keys.StockSymbolKey(psp.Stock.Symbol), priceData.Price, time.Second*20)
		} else {
			psp.Stock.Price = cachedSymbolsMap[psp.Stock.Symbol]
		}
	}

	_, err = redisPipeline.Exec(ctx)
	if err != nil {
		utils.LogWarning("Error caching stock prices")
	}

	return nil
}
