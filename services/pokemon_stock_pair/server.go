package pokemon_stock_pair

import (
	pb "pokestocks/proto/pokemon_stock_pair"

	"github.com/jackc/pgx/v5"
)

type Server struct {
	pb.UnimplementedPokemonStockPairServiceServer
	DB *pgx.Conn
}
