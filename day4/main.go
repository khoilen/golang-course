package main

import (
	"os"
	"user-auth/controllers"
	"user-auth/database"
	"user-auth/middleware"
	"user-auth/models"
	"user-auth/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.InitDB()
	db.AutoMigrate(&models.User{})

	userService := &services.UserService{DB: db}
	userController := &controllers.UserController{UserService: userService}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/logout", userController.Logout)

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/user/:username", userController.GetUser)
	auth.PUT("/user/:username", userController.UpdateUser)
	auth.POST("/user/:username/profile", userController.UploadProfile)

	router.Run(":8080")
}
