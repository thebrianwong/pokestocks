package client_manager

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ClientManager struct {
	DB                     *pgxpool.Pool
	ElasticClient          *elasticsearch.TypedClient
	AlpacaMarketDataClient *marketdata.Client
	AlpacaTradingClient    *alpaca.Client
	RedisClient            *redis.Client
}
