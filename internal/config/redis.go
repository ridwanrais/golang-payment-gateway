package config

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func getRedisClient() *redis.Client {
	// Replace 'YOUR_REDIS_HOST:PORT' with the address of your Redis server
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	fmt.Println("Connected to Redis!")
	return rdb
}
