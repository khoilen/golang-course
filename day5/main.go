package main

import (
	"redis-go/handlers"
	"redis-go/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	redis.InitRedis()

	http := gin.Default()

	http.POST("/login", handlers.Login)
	http.GET("/ping", handlers.Ping)
	http.GET("/ping-count", handlers.PingCount)

	http.Run(":8080")
}
