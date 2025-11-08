package configs

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// ConnectREDISDB connects to Redis
func ConnectREDISDB() error {
	redisURL := RedisURL()
	
	// 0 ALLAH MKENI TRANU PASHA ZOTIN, QEESHTU POM SHTINI ME BA SENE
	opt, err := redis.ParseURL("redis://syn_redis:6379")
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}


    log.Printf("Connecting to Redis at: %s", redisURL)
	
	redisClient = redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetRedisClient returns the Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}