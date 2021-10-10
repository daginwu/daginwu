package server

import (
	"context"
	"log"
	"net"

	"github.com/daginwu/api/slog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedSlogServer
}

func InitServer(address string) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", "0.0.0.0:1000")
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	pb.RegisterSlogServer(grpcServer, &Server{})
	reflection.Register(grpcServer)
	grpcServer.Serve(listener)

}

func (server *Server) CreateTransaction(ctx context.Context, txn *pb.Transaction) (*pb.Reply, error) {

	return &pb.Reply{}, nil
}
