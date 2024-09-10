package main

import (
	"net"
	"os"
	"pokestocks/utils"

	psp_pb "pokestocks/proto/pokemon_stock_pair"
	psp_service "pokestocks/services/pokemon_stock_pair"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	utils.LoadEnvVars("")
	port := os.Getenv("GRPC_PORT")
	conn := utils.ConnectToDb()
	elasticClient := utils.CreateTypedElasticClient("")
	alpacaClient := utils.CreateAlpacaClient()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		utils.LogFailureError("Failed to listen on port "+port, err)
	}

	s := grpc.NewServer()
	psp_pb.RegisterPokemonStockPairServiceServer(
		s,
		&psp_service.Server{
			DB:            conn,
			ElasticClient: elasticClient,
			AlpacaClient:  alpacaClient,
		},
	)

	reflection.Register(s)

	utils.LogSuccess("gRPC server listening at " + lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		utils.LogFailureError("Failed to serve", err)
	}
}
