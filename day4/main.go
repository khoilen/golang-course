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
	db.AutoMigrate(&models.User{}, &models.Post{})

	userService := &services.UserService{DB: db}
	userController := &controllers.UserController{UserService: userService}

	postService := &services.PostService{DB: db}
	postController := &controllers.PostController{PostService: postService}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.POST("/logout", userController.Logout)
	router.GET("/posts", postController.ListPosts)

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/user/:username", userController.GetUser)
	auth.PUT("/user/:username", userController.UpdateUser)
	auth.POST("/user/:username/profile", userController.UploadProfile)

	auth.POST("/posts", postController.CreatePost)
	auth.GET("/posts/:postID", postController.GetPost)
	auth.PUT("/posts/:postID", postController.UpdatePost)
	auth.DELETE("/posts/:postID", postController.DeletePost)

	router.Run(":8080")
}
