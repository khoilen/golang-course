package services

import (
	"user-auth/models"

	"gorm.io/gorm"
)

type PostService struct {
	DB *gorm.DB
}

func (s *PostService) CreatePost(post *models.Post) error {
	return s.DB.Create(post).Error
}

func (s *PostService) GetPostByID(postID uint) (*models.Post, error) {
	var post models.Post
	err := s.DB.First(&post, postID).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) UpdatePost(post *models.Post) error {
	return s.DB.Save(post).Error
}

func (s *PostService) DeletePost(postID uint) error {
	return s.DB.Delete(&models.Post{}, postID).Error
}

func (s *PostService) GetPostsByAuthor(authorID uint) ([]models.Post, error) {
	var posts []models.Post
	err := s.DB.Where("author_id = ?", authorID).Find(&posts).Error
	return posts, err
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := s.DB.Find(&posts).Error
	return posts, err
}
