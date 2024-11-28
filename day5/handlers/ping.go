package handlers

import (
	"net/http"
	"redis-go/redis"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var lock = &sync.Mutex{}
var isLocked = false

func Ping(c *gin.Context) {
	sessionId := c.GetHeader("Authorization")
	if sessionId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Session ID"})
	}

	username, err := redis.Client.Get(redis.Ctx, sessionId).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Session ID"})
		return
	}

	if !tryLock() {
		c.JSON(http.StatusOK, gin.H{"message": "deny"})
		return
	}

	defer unlock()

	time.Sleep(5 * time.Second)

	redis.Client.Incr(redis.Ctx, "ping_count_"+username)

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func tryLock() bool {
	locked := lock.TryLock()
	if locked {
		isLocked = true
	}
	return locked
}

func unlock() {
	if isLocked {
		lock.Unlock()
		isLocked = false
	}
}
