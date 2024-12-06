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

type PostService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func NewPostService(db *gorm.DB, redisClient *redis.Client) *PostService {
	return &PostService{DB: db, RedisClient: redisClient}
}

func (service *PostService) CreatePost(post *models.Post) error {
	return service.DB.Create(post).Error
}

func (service *PostService) GetPostByID(postID uint) (*models.Post, error) {
	cacheKey := fmt.Sprintf("posts:post:%d", postID)
	val, err := service.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == redis.Nil {
		var post models.Post
		err := service.DB.Where("id = ?", postID).Find(&post).Error
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(post)
		if err != nil {
			return nil, err
		}
		service.RedisClient.Set(config.Ctx, cacheKey, data, 15*time.Minute).Err()
		return &post, err
	} else if err != nil {
		return nil, err
	}

	var post models.Post
	err = json.Unmarshal([]byte(val), &post)
	return &post, err
}

func (service *PostService) UpdatePost(post *models.Post) error {
	return service.DB.Save(post).Error
}

func (service *PostService) DeletePost(postID uint) error {
	return service.DB.Delete(&models.Post{}, postID).Error
}

func (service *PostService) GetPostsByAuthor(authorID uint) ([]models.Post, error) {
	var posts []models.Post
	err := service.DB.Where("author_id = ?", authorID).Find(&posts).Error
	return posts, err
}

func (service *PostService) GetAllPosts() ([]models.Post, error) {
	cacheKey := "posts:list"
	val, err := service.RedisClient.Get(config.Ctx, cacheKey).Result()

	if err == redis.Nil {
		var posts []models.Post
		err := service.DB.Find(&posts).Error

		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(posts)

		if err != nil {
			return nil, err
		}

		service.RedisClient.Set(config.Ctx, cacheKey, data, 15*time.Minute).Err()
		return posts, nil

	} else if err != nil {
		return nil, err
	}

	var posts []models.Post
	err = json.Unmarshal([]byte(val), &posts)
	return posts, err
}
