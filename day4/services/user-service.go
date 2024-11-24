package services

import (
	"user-auth/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (s *UserService) AddUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *UserService) GetUsersByUserName(username string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("user_name = ?", username).First(&user).Error
	return &user, err
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.DB.Save(user).Error
}
