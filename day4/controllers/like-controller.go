package controllers

import (
	"net/http"
	"strconv"
	"user-auth/models"
	"user-auth/services"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	LikeService *services.LikeService
	PostService *services.PostService
}

func (ctrl *LikeController) LikePost(c *gin.Context) {
	var like models.Like
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

	like.PostID = uint(postID)
	like.UserID = userID.(uint)

	if err := ctrl.LikeService.AddLike(&like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post liked successfully"})
}

func (ctrl *LikeController) UnlikePost(c *gin.Context) {
	var like models.Like
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
	like.PostID = uint(postID)
	like.UserID = userID.(uint)

	post, err := ctrl.PostService.GetPostByID(uint(postID))
	if err != nil || post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	like.PostID = uint(postID)

	if err := ctrl.LikeService.RemoveLike(&like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}
