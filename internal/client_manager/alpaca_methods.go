package client_manager

import (
	"context"
	common_pb "pokestocks/proto/common"
	redis_keys "pokestocks/redis"
	"pokestocks/utils"
	"strconv"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/redis/go-redis/v9"
)

func (cc *ClientManager) GetAlpacaClock() (*alpaca.Clock, error) {
	clock, err := cc.AlpacaTradingClient.GetClock()
	if err != nil {
		return nil, err
	}
	return clock, nil
}

func (cc *ClientManager) IsMarketOpen(ctx context.Context) (bool, error) {
	redisClient := cc.RedisClient
	redisPipeline := redisClient.Pipeline()

	cachedMarketStatus, err := redisClient.Get(ctx, redis_keys.MarketStatusKey()).Result()
	if err == nil {
		return cachedMarketStatus == "open", nil
	} else {
		clock, err := cc.GetAlpacaClock()
		if err != nil {
			utils.LogWarningError("Error calling Alpaca clock API", err)
			return false, err
		}

		marketIsOpen := clock.IsOpen
		if marketIsOpen {
			marketCloseTime := clock.NextClose
			redisPipeline.Set(ctx, redis_keys.MarketStatusKey(), "open", 0)
			redisPipeline.ExpireAt(ctx, redis_keys.MarketStatusKey(), marketCloseTime)
		} else {
			marketOpenTime := clock.NextOpen
			redisPipeline.Set(ctx, redis_keys.MarketStatusKey(), "close", 0)
			redisPipeline.ExpireAt(ctx, redis_keys.MarketStatusKey(), marketOpenTime)
		}

		_, err = redisPipeline.Exec(ctx)
		if err != nil {
			utils.LogWarningError("Error caching market status to Redis", err)
		}

		return marketIsOpen, nil
	}
}

func (cc *ClientManager) EnrichWithStockPrices(ctx context.Context, psps []*common_pb.PokemonStockPair) error {
	alpacaMarketDataClient := cc.AlpacaMarketDataClient
	redisClient := cc.RedisClient
	redisPipeline := redisClient.Pipeline()

	marketIsOpen, err := cc.IsMarketOpen(ctx)
	if err != nil {
		utils.LogWarningError("Error checking if market is open", err)
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
			utils.LogWarningError("Error checking if cached stock prices exists. Falling back to retrieving data from Alpaca", err)
		}
	}

	nonCachedData := map[string]marketdata.Trade{}
	var nextMarketOpen time.Time

	// this will always run when market is open
	// skip querying Alpaca if there are no cache misses when market is close
	if len(symbols) > 0 {
		requestParams := marketdata.GetLatestTradeRequest{}
		nonCachedData, err = alpacaMarketDataClient.GetLatestTrades(symbols, requestParams)
		if err != nil {
			return err
		}

		cachedMarketOpen, err := redisClient.Get(ctx, redis_keys.NextMarketOpenKey()).Result()
		if err != nil && err != redis.Nil {
			utils.LogWarningError("Error checking Redis for next market open", err)
		}

		if cachedMarketOpen == "" {
			clock, err := cc.GetAlpacaClock()
			if err != nil {
				utils.LogWarningError("Error parsing date string to time.Time. Defaulting to 3 hours", err)
				nextMarketOpen = time.Now().Add(time.Hour * 3)
			} else {
				nextMarketOpen = clock.NextOpen
				redisPipeline.Set(ctx, redis_keys.NextMarketOpenKey(), nextMarketOpen, 0)
				redisPipeline.ExpireAt(ctx, redis_keys.NextMarketOpenKey(), nextMarketOpen)

				_, err = redisPipeline.Exec(ctx)
				if err != nil {
					utils.LogWarningError("Error caching next market open to Redis", err)
				}
			}
		} else {
			nextMarketOpen, err = time.Parse("2006-01-02T15:04:05-07:00", cachedMarketOpen)
			if err != nil {
				utils.LogWarningError("Error parsing date string to time.Time. Defaulting to 3 hours", err)
				nextMarketOpen = time.Now().Add(time.Hour * 3)
			}
		}
	}

	for _, psp := range psps {
		priceData, ok := nonCachedData[psp.Stock.Symbol]
		if ok {
			psp.Stock.Price = &priceData.Price
			if !marketIsOpen {
				redisPipeline.Set(ctx, redis_keys.StockSymbolKey(psp.Stock.Symbol), priceData.Price, 0)
				redisPipeline.ExpireAt(ctx, redis_keys.StockSymbolKey(psp.Stock.Symbol), nextMarketOpen)
			}
		} else {
			psp.Stock.Price = cachedSymbolsMap[psp.Stock.Symbol]
		}
	}

	if !marketIsOpen {
		_, err = redisPipeline.Exec(ctx)
		if err != nil {
			utils.LogWarning("Error caching stock prices to Redis")
		}
	}

	return nil
}
