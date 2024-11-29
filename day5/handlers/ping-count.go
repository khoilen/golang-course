package handlers

import (
	"net/http"
	"redis-go/redis"

	"github.com/gin-gonic/gin"
)

func PingCount(c *gin.Context) {
	sessionID := c.GetHeader("Authorization")
	if sessionID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID is required"})
		return
	}

	username, err := redis.Client.Get(redis.Ctx, sessionID).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Session ID"})
		return
	}

	pingCountKey := "ping_count_" + username
	pingCount, err := redis.Client.Get(redis.Ctx, pingCountKey).Int()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := redis.Client.PFCount(redis.Ctx, "ping_users_hll").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ping_count": pingCount, "unique_users": count})
}
