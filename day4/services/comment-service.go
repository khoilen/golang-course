package services

import (
	"encoding/json"
	"fmt"
	"time"
	"user-auth/config"
	"user-auth/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CommentService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func NewCommentService(db *gorm.DB, redisClient *redis.Client) *CommentService {
	return &CommentService{DB: db, RedisClient: redisClient}
}

func (s *CommentService) AddComment(comment *models.Comment) error {
	err := s.DB.Create(comment).Error
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("comments:post:%d", comment.PostID)
	s.RedisClient.Del(config.Ctx, cacheKey)

	return nil
}

func (s *CommentService) GetCommentsByPost(postID uint) ([]models.Comment, error) {
	cacheKey := fmt.Sprintf("comments:post:%d", postID)

	val, err := s.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == redis.Nil {
		var comments []models.Comment
		err := s.DB.Where("post_id = ?", postID).Find(&comments).Error
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(comments)
		if err != nil {
			return nil, err
		}
		s.RedisClient.Set(config.Ctx, cacheKey, data, 5*time.Minute).Err()
		return comments, nil
	} else if err != nil {
		return nil, err
	}
	var comments []models.Comment
	err = json.Unmarshal([]byte(val), &comments)
	return comments, err
}
