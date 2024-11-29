package handlers

import (
	"net/http"
	"redis-go/redis"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var lock = &sync.Mutex{}

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

	rateLimitKey := "rate_limit_" + username
	val, _ := redis.Client.Get(redis.Ctx, rateLimitKey).Int()

	if val >= 2 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. You can only call /ping 2 times in 60 seconds."})
		return
	}

	redis.Client.Incr(redis.Ctx, rateLimitKey)

	if val == 0 {
		redis.Client.Expire(redis.Ctx, rateLimitKey, 60*time.Second)
	}

	lock.Lock()
	defer lock.Unlock()

	time.Sleep(5 * time.Second)

	redis.Client.Incr(redis.Ctx, "ping_count_"+username)
	redis.Client.PFAdd(redis.Ctx, "ping_users_hll", username)

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
