package controllers

import (
	"net/http"
	"strconv"
	"user-auth/models"
	"user-auth/services"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	PostService *services.PostService
}

func (ctrl *PostController) CreatePost(c *gin.Context) {
	var newPost models.Post
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	newPost.AuthorID = authorID.(uint)

	if err := ctrl.PostService.CreatePost(&newPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

func (ctrl *PostController) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := ctrl.PostService.GetPostByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (ctrl *PostController) UpdatePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var updatePost models.Post
	if err := c.ShouldBindJSON(&updatePost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := ctrl.PostService.GetPostByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	authorID, exists := c.Get("userID")
	if !exists || post.AuthorID != authorID.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to update this post"})
		return
	}

	post.Title = updatePost.Title
	post.Content = updatePost.Content

	if err := ctrl.PostService.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func (ctrl *PostController) DeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := ctrl.PostService.GetPostByID(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	authorID, exists := c.Get("userID")
	if !exists || post.AuthorID != authorID.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to delete this post"})
		return
	}

	if err := ctrl.PostService.DeletePost(uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (ctrl *PostController) ListPosts(c *gin.Context) {
	posts, err := ctrl.PostService.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}
