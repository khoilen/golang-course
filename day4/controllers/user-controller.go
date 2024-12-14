package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"
	"user-auth/config"
	"user-auth/models"
	"user-auth/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	key    = []byte("super-secret-key")
	store  = sessions.NewCookieStore(key)
	jwtKey = []byte("your_jwt_secret_key")
)

type UserController struct {
	UserService *services.UserService
	RedisClient *redis.Client
}

func (ctrl *UserController) SignUp(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if _, err := ctrl.UserService.GetUsersByUserName(newUser.UserName); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.Password = string(hashedPassword)

	if err := ctrl.UserService.AddUser(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_name":  newUser.UserName,
		"email":      newUser.Email,
		"first_name": newUser.FirstName,
		"last_name":  newUser.LastName,
		"birthday":   newUser.BirthDay,
	})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var loginUser models.User

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}

	user, err := ctrl.UserService.GetUsersByUserName(loginUser.UserName)
	if err != nil || user.UserName != loginUser.UserName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.CustomClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.UserName,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	redisKey := fmt.Sprintf("session:user:%s", user.UserName)
	fmt.Print(user.UserName)

	err = ctrl.RedisClient.Set(config.Ctx, redisKey, tokenString, 24*time.Hour).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"user_name":  user.UserName,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"birthday":   user.BirthDay,
		},
	})
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := ctrl.UserService.GetUsersByUserName(username)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	var updateUser models.User

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	username, _ := c.Get("username")

	user, err := ctrl.UserService.GetUsersByUserName(username.(string))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if updateUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.Password = string(hashedPassword)
	}

	if err := ctrl.UserService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User information updated",
		"user": gin.H{
			"user_name":  user.UserName,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"birthday":   user.BirthDay,
		},
	})
}

func (ctrl *UserController) UploadProfile(c *gin.Context) {
	username := c.Param("username")

	file, err := c.FormFile("profile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := ctrl.UserService.GetUsersByUserName(username)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	filename := filepath.Join("uploads", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Profile = filename
	if err := ctrl.UserService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile picture uploaded"})

}

func (ctrl *UserController) Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	session.Values["authenticated"] = false
	session.Save(c.Request, c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
