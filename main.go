package main

import (
	"log"
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

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %v: %v", port, err)
	}

	s := grpc.NewServer()
	psp_pb.RegisterPokemonStockPairServiceServer(s, &psp_service.Server{DB: conn})

	reflection.Register(s)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
