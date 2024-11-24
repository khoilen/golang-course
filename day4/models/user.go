package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `json:"username" gorm:"uniqueIndex"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
}

type CustomClaims struct {
	UserID uint `json:"userID"`
	jwt.StandardClaims
}
