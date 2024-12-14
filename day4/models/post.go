package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Image    string `json:"image_path"`
	Content  string `json:"content"`
	AuthorID uint   `json:"fk_user_id"`
	Visible  bool   `json:"visible"`
}
