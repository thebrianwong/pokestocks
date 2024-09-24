package transaction

import (
	"pokestocks/internal/structs"
	pb "pokestocks/proto/transaction"
)

type Server struct {
	*pb.UnimplementedTransactionServiceServer
	*structs.ClientConfig
}
