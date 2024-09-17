package pokemon_stock_pair

import (
	pb "pokestocks/proto/pokemon_stock_pair"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	pb.UnimplementedPokemonStockPairServiceServer
	DB                     *pgxpool.Pool
	ElasticClient          *elasticsearch.TypedClient
	AlpacaMarketDataClient *marketdata.Client
	AlpacaBrokerClient     *alpaca.Client
	RedisClient            *redis.Client
}
