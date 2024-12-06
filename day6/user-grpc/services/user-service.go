package services

import (
	"context"
	"fmt"
	"grpchomework/user-grpc/config"
	"grpchomework/user-grpc/models"
	"time"

	userPb "grpchomework/user-grpc/proto/user"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	userPb.UnimplementedUserServiceServer
}

type CustomClaims struct {
	UserID uint `json:"userID"`
	jwt.StandardClaims
}

func (s *UserService) AddUser(user *userPb.LoginRequest) error {
	return s.DB.Create(user).Error
}

func (s *UserService) GetUsersByUserName(username string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("user_name = ?", username).First(&user).Error
	return &user, err
}

func (s *UserService) UpdateUser(user *userPb.LoginRequest) error {
	return s.DB.Save(user).Error
}

func (s *UserService) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	user, err := s.GetUsersByUserName(req.Username)
	if err != nil || user.UserName != req.Username {
		return &userPb.LoginResponse{Message: "Invalid username or password"}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &userPb.LoginResponse{Message: "Invalid username or password"}, nil
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &CustomClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.UserName,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_jwt_secret_key"))
	if err != nil {
		return &userPb.LoginResponse{Message: "Failed to generate token"}, err
	}

	redisKey := fmt.Sprintf("session:%s", user.UserName)
	err = s.RedisClient.Set(config.Ctx, redisKey, tokenString, 24*time.Hour).Err()
	if err != nil {
		return &userPb.LoginResponse{Message: "Failed to store session"}, err
	}

	return &userPb.LoginResponse{Token: tokenString, Message: "Login successful"}, nil
}
