// File: fitness/test/redis/redis_test.go
package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisConnection(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := rdb.Set(ctx, "test_key", "test_value", time.Second*10).Err()
	if err != nil {
		t.Fatalf("Redis SET failed: %v", err)
	}

	val, err := rdb.Get(ctx, "test_key").Result()
	if err != nil || val != "test_value" {
		t.Fatalf("Redis GET failed: got '%s', err: %v", val, err)
	}
}