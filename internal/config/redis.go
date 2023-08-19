package config

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func getRedisClient(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	// Check Redis connection using PING command
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		// Could not connect to Redis, return the error
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis!")
	return rdb, nil
}
