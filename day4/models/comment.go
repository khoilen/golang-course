package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	UserID  uint   `json:"fk_user_id"`
	PostID  uint   `json:"fk_post_id"`
}
