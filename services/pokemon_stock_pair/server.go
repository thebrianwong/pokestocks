package pokemon_stock_pair

import (
	pb "pokestocks/proto/pokemon_stock_pair"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pb.UnimplementedPokemonStockPairServiceServer
	DB            *pgxpool.Pool
	ElasticClient *elasticsearch.TypedClient
}
