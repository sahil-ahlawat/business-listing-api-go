package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func Init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})
}

func Get(key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}