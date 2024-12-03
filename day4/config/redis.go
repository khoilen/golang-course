package config

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var Ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
