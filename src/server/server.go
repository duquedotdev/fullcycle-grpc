package main

import (
	"github.com/duquedotdev/fullcycle-grpc/pb"
	"github.com/duquedotdev/fullcycle-grpc/src/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", "localhost:50052")

	if err != nil {
		log.Fatalf("Could not listen to port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err == nil {
		log.Fatalf("Could not serve: %v", err)
	}

}
