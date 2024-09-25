package pokemon_stock_pair

import (
	cm "pokestocks/internal/client_manager"
	pb "pokestocks/proto/pokemon_stock_pair"
)

type Server struct {
	*pb.UnimplementedPokemonStockPairServiceServer
	*cm.ClientManager
}
