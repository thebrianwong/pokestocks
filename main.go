package main

import (
	"net"
	"os"
	cm "pokestocks/internal/client_manager"
	"pokestocks/utils"

	psp_pb "pokestocks/proto/pokemon_stock_pair"
	psp_service "pokestocks/services/pokemon_stock_pair"

	transaction_pb "pokestocks/proto/transaction"
	transaction_service "pokestocks/services/transaction"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	utils.LoadEnvVars("")
	port := os.Getenv("GRPC_PORT")
	conn := utils.ConnectToDb()
	elasticClient := utils.CreateTypedElasticClient("")
	alpacaMarketDataClient := utils.CreateAlpacaMarketDataClient()
	alpacaTradingClient := utils.CreateAlpacaTradingClient()
	redisClient := utils.CreateRedisClient()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		utils.LogFailureError("Failed to listen on port "+port, err)
	}

	clientManager := cm.ClientManager{
		DB:                     conn,
		ElasticClient:          elasticClient,
		AlpacaMarketDataClient: alpacaMarketDataClient,
		AlpacaTradingClient:    alpacaTradingClient,
		RedisClient:            redisClient,
	}

	s := grpc.NewServer()
	psp_pb.RegisterPokemonStockPairServiceServer(
		s,
		&psp_service.Server{
			UnimplementedPokemonStockPairServiceServer: &psp_pb.UnimplementedPokemonStockPairServiceServer{},
			ClientManager: &clientManager,
		},
	)
	transaction_pb.RegisterTransactionServiceServer(
		s,
		&transaction_service.Server{
			UnimplementedTransactionServiceServer: &transaction_pb.UnimplementedTransactionServiceServer{},
			ClientManager:                         &clientManager,
		},
	)

	reflection.Register(s)

	utils.LogSuccess("gRPC server listening at " + lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		utils.LogFailureError("Failed to serve", err)
	}
}
