package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	UserID uint `json:"fk_user_id" gorm:"index:idx_user_post"`
	PostID uint `json:"fk_post_id" gorm:"index:idx_user_post"`
}
