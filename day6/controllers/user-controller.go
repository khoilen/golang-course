package controllers

import (
	"context"
	pb "grpchomework/user-grpc/proto/user"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	GrpcClient pb.UserServiceClient
}

func (ctrl *UserController) Login(c *gin.Context) {
	var loginUser pb.LoginRequest

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grpcReq := &pb.LoginRequest{
		Username: loginUser.Username,
		Password: loginUser.Password,
	}

	log.Println("Sending Login request to gRPC server")
	grpcResp, err := ctrl.GrpcClient.Login(context.Background(), grpcReq)
	if err != nil {
		log.Printf("gRPC call failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	if grpcResp.Message != "Login successful" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": grpcResp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": grpcResp.Token})
}
