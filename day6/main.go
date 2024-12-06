package main

import (
	"fmt"
	"grpchomework/controllers"
	pb "grpchomework/user-grpc/proto/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect to gRPC server: %v\n", err)
		return
	}
	defer conn.Close()

	grpcClient := pb.NewUserServiceClient(conn)

	userController := &controllers.UserController{GrpcClient: grpcClient}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/login", userController.Login)

	router.Run(":8080")
}
