package services

import (
	"user-auth/models"

	"gorm.io/gorm"
)

type FollowService struct {
	DB *gorm.DB
}

func (s *FollowService) GetFollowsByUserID(userID string) ([]models.User, error) {
	var follows []models.User

	err := s.DB.Table("users").Select("users.*").
		Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.user_id = ?", userID).
		Find(&follows).Error

	if err != nil {
		return nil, err
	}
	return follows, nil
}

func (s *FollowService) FollowUser(followerID, userID uint) error {
	follow := models.Follow{
		UserID:     userID,
		FollowerID: followerID,
	}
	return s.DB.Create(&follow).Error
}
