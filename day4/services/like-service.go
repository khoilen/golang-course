package services

import (
	"user-auth/models"

	"gorm.io/gorm"
)

type LikeService struct {
	DB *gorm.DB
}

func (s *LikeService) AddLike(like *models.Like) error {
	return s.DB.Create(like).Error
}

func (s *LikeService) RemoveLike(like *models.Like) error {
	return s.DB.Where("user_id = ? AND post_id = ?", like.UserID, like.PostID).Delete(&models.Like{}).Error
}
