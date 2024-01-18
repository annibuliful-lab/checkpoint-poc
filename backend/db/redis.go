package db

import (
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	onceRedis   sync.Once
)

func GetRedisClient() *redis.Client {
	onceRedis.Do(func() {

		redisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})

	return redisClient
}
