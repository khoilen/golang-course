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
		Where("follows.deleted_at IS NULL").
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

func (s *FollowService) UnFollowUser(followerID, userID uint) error {
	return s.DB.Where("user_id = ? AND follower_id = ?", userID, followerID).Delete(&models.Follow{}).Error
}

func (s *FollowService) IsFollowing(followerID, userID uint) (bool, error) {
	var count int64
	err := s.DB.Model(&models.Follow{}).
		Where("user_id = ? AND follower_id = ?", userID, followerID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *FollowService) GetFollowsPosts(followerID, userID uint) ([]models.Post, error) {
	var posts []models.Post

	err := s.DB.Table("posts").Select("posts.*").
		Joins("JOIN follows ON follows.user_id = posts.author_id").
		Where("follows.follower_id = ?", followerID).
		Where("follows.user_id = ?", userID).
		Where("follows.deleted_at IS NULL").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}
	return posts, nil
}
