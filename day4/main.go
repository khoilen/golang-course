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
	userController := &controllers.UserController{UserService: userService, RedisClient: redisClient}

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
	v1 := router.Group("/v1")
	{
		v1.POST("/users", userController.SignUp)
		v1.POST("/sessions", userController.Login)
		v1.GET("/posts", postController.ListPosts)
		v1.POST("/logout", userController.Logout)

		auth := router.Group("/")
		auth.Use(middleware.AuthMiddleware())
		v1.GET("/users/:username", userController.GetUser)
		v1.PUT("/users/", userController.UpdateUser)
		v1.POST("/users/:username/profile", userController.UploadProfile)

		v1.POST("/posts", postController.CreatePost)
		v1.GET("/posts/:postID", postController.GetPost)
		v1.PUT("/posts/:postID", postController.UpdatePost)
		v1.DELETE("/posts/:postID", postController.DeletePost)

		v1.POST("/posts/:postID/like", likeController.LikePost)
		v1.DELETE("/posts/:postID/like", likeController.UnlikePost)

		v1.POST("/posts/:postID/comments", commentController.AddComment)
		v1.GET("/posts/:postID/comments", commentController.GetComments)
	}

	router.Run(":8080")
}
