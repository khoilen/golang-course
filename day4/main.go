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
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Like{}, &models.Comment{}, &models.Follow{})

	followService := &services.FollowService{DB: db}
	followController := &controllers.FollowController{FollowService: followService}

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

		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/users/:username", userController.GetUser)
			auth.PUT("/users", userController.UpdateUser)
			auth.POST("/users/:username/profile", userController.UploadProfile)
			auth.GET("/friends/:user_id", followController.GetFollows)
			auth.POST("/friends/:user_id", followController.FollowUser)
			auth.DELETE("/friends/:user_id", followController.UnFollowUser)
			auth.GET("/friends/:user_id/posts", followController.GetFollowsPosts)

			auth.POST("/posts", postController.CreatePost)
			auth.GET("/posts/:postID", postController.GetPost)
			auth.PUT("/posts/:postID", postController.UpdatePost)
			auth.DELETE("/posts/:postID", postController.DeletePost)

			auth.POST("/posts/:postID/like", likeController.LikePost)
			auth.DELETE("/posts/:postID/like", likeController.UnlikePost)

			auth.POST("/posts/:postID/comments", commentController.AddComment)
			auth.GET("/posts/:postID/comments", commentController.GetComments)
		}
	}

	router.Run(":8080")
}
