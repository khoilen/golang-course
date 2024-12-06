package main

import (
	"context"
	"grpchomework/user-grpc/config"
	pb "grpchomework/user-grpc/proto/user"
	"grpchomework/user-grpc/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	userService *services.UserService
}

func NewGRPCServer() *server {
	db := config.InitDB()
	redisClient := config.NewRedisClient()
	return &server{
		userService: &services.UserService{DB: db, RedisClient: redisClient},
	}
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.userService.Login(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := NewGRPCServer()
	pb.RegisterUserServiceServer(s, srv)
	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
