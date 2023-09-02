package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

func getRedisClient(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	for {
		// Check Redis connection using PING command
		_, err := rdb.Ping(ctx).Result()
		if err == nil {
			// Successfully connected to Redis
			fmt.Println("Connected to Redis!")
			return rdb, nil
		}

		// If connection failed, wait for 5 seconds before retrying
		fmt.Println("Failed to connect to Redis. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
