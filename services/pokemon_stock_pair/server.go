package pokemon_stock_pair

import (
	"pokestocks/internal/structs"
	pb "pokestocks/proto/pokemon_stock_pair"
)

type Server struct {
	*pb.UnimplementedPokemonStockPairServiceServer
	*structs.ClientConfig
}
