package models

import (
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	UserID     uint `json:"fk_user_id" gorm:"index:idx_user_id"`
	FollowerID uint `json:"fk_follower_id" gorm:"index:idx_user_id"`
}
