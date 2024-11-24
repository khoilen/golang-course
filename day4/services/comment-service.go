package services

import (
	"user-auth/models"

	"gorm.io/gorm"
)

type CommentService struct {
	DB *gorm.DB
}

func (s *CommentService) AddComment(comment *models.Comment) error {
	return s.DB.Create(comment).Error
}

func (s *CommentService) GetCommentsByPost(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := s.DB.Where("post_id = ?", postID).Find(&comments).Error
	return comments, err
}
