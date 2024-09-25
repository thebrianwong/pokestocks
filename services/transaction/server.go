package transaction

import (
	cm "pokestocks/internal/client_manager"
	pb "pokestocks/proto/transaction"
)

type Server struct {
	*pb.UnimplementedTransactionServiceServer
	*cm.ClientManager
}
