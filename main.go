package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "train-seat-allocation/proto/pb"
	"train-seat-allocation/repositories"
	"train-seat-allocation/services"
	"train-seat-allocation/transport"
)

func main() {
	repo := repositories.NewInMemoryTicketRepository()
	ticketService := services.NewTicketService(repo)
	grpcServer := grpc.NewServer()

	pb.RegisterTrainTicketServiceServer(grpcServer, transport.NewGRPCServer(ticketService))

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server running on port :5001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
