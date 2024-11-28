package handlers

import (
	"net/http"
	"redis-go/redis"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	sessionID := uuid.New().String()
	redis.Client.Set(redis.Ctx, sessionID, username, time.Hour)
	c.JSON(http.StatusOK, gin.H{"session_id": sessionID, "username": username})
}
