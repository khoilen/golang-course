package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName  string `json:"username" gorm:"uniqueIndex;size:191"`
	Password  string `json:"password" gorm:"size:191"`
	Profile   string `json:"profile" gorm:"size:191"`
	Email     string `json:"email" gorm:"size:191"`
	FirstName string `json:"first_name" gorm:"size:191"`
	LastName  string `json:"last_name" gorm:"size:191"`
	BirthDay  string `json:"birthday" gorm:"size:191"`
}

type CustomClaims struct {
	UserID uint `json:"userID"`
	jwt.StandardClaims
}
