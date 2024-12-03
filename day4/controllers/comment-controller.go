package controllers

import (
	"net/http"
	"strconv"
	"user-auth/config"
	"user-auth/models"
	"user-auth/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct {
	CommentService *services.CommentService
	PostService    *services.PostService
}

func NewCommentController(db *gorm.DB) *CommentController {
	redisClient := config.NewRedisClient()
	commentService := services.NewCommentService(db, redisClient)
	postService := services.NewPostService(db)
	return &CommentController{
		CommentService: commentService,
		PostService:    postService,
	}
}

func (ctrl *CommentController) AddComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	post, err := ctrl.PostService.GetPostByID(uint(postID))
	if err != nil || post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	comment.PostID = uint(postID)
	comment.UserID = userID.(uint)

	if err := ctrl.CommentService.AddComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully"})
}

func (ctrl *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := ctrl.PostService.GetPostByID(uint(postID))
	if err != nil || post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comments, err := ctrl.CommentService.GetCommentsByPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
