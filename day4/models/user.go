package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
}
