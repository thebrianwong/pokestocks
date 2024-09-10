package pokemon_stock_pair

import (
	pb "pokestocks/proto/pokemon_stock_pair"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pb.UnimplementedPokemonStockPairServiceServer
	DB            *pgxpool.Pool
	ElasticClient *elasticsearch.TypedClient
	AlpacaClient  *marketdata.Client
}
