package handlers

import (
	"net/http"
	"redis-go/redis"
	"sort"

	"github.com/gin-gonic/gin"
)

func TopHandler(c *gin.Context) {
	keys, err := redis.Client.Keys(redis.Ctx, "ping_count_*").Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	counts := make(map[string]int)
	for _, key := range keys {
		count, err := redis.Client.Get(redis.Ctx, key).Int()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		username := key[len("ping_count_"):]
		counts[username] = count
	}

	type KeyValue struct {
		Key   string
		Value int
	}

	var kvSlice []KeyValue
	for k, v := range counts {
		kvSlice = append(kvSlice, KeyValue{Key: k, Value: v})
	}

	sort.Slice(kvSlice, func(i, j int) bool { return kvSlice[i].Value > kvSlice[j].Value })

	if len(kvSlice) > 10 {
		kvSlice = kvSlice[:10]
	}

	c.JSON(http.StatusOK, kvSlice)

}
