package main

import (
	"os"
	"user-auth/config"
	"user-auth/controllers"
	"user-auth/middleware"
	"user-auth/models"
	"user-auth/services"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	redisClient := config.NewRedisClient()
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Like{}, &models.Comment{})

	userService := &services.UserService{DB: db}
	userController := &controllers.UserController{UserService: userService}

	postService := &services.PostService{DB: db, RedisClient: redisClient}
	postController := &controllers.PostController{PostService: postService}

	likeService := &services.LikeService{DB: db}
	likeController := &controllers.LikeController{LikeService: likeService, PostService: postService}

	commentService := &services.CommentService{DB: db, RedisClient: redisClient}
	commentController := &controllers.CommentController{CommentService: commentService, PostService: postService}

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.GET("/posts", postController.ListPosts)
	router.POST("/logout", userController.Logout)

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/user/:username", userController.GetUser)
	auth.PUT("/user/:username", userController.UpdateUser)
	auth.POST("/user/:username/profile", userController.UploadProfile)

	auth.POST("/posts", postController.CreatePost)
	auth.GET("/posts/:postID", postController.GetPost)
	auth.PUT("/posts/:postID", postController.UpdatePost)
	auth.DELETE("/posts/:postID", postController.DeletePost)

	auth.POST("/posts/:postID/like", likeController.LikePost)
	auth.DELETE("/posts/:postID/like", likeController.UnlikePost)

	auth.POST("/posts/:postID/comments", commentController.AddComment)
	auth.GET("/posts/:postID/comments", commentController.GetComments)

	router.Run(":8080")
}
