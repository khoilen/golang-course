package controllers

import (
	"net/http"
	"strconv"
	"user-auth/services"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	FollowService *services.FollowService
}

func (uc *FollowController) GetFollows(c *gin.Context) {
	userID := c.Param("user_id")

	follows, err := uc.FollowService.GetFollowsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch follows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"follows": follows})
}

func (fc *FollowController) FollowUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	followerIDUint, ok := followerID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse follower ID"})
		return
	}

	err = fc.FollowService.FollowUser(followerIDUint, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully followed the user"})
}
