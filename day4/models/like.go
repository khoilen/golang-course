package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"index:idx_user_post"`
	PostID uint `json:"post_id" gorm:"index:idx_user_post"`
}
