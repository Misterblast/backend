package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")
	db := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  3 * time.Second,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis: %v", err)
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info("Redis connected successfully")
	return rdb, nil
}

func Get(ctx context.Context, redisKey string, rdb *redis.Client) (string, error) {
	val, err := rdb.Get(ctx, redisKey).Result()
	if err != nil {
		// log.Error("Failed to get value from Redis: %v", err)
		return "", fmt.Errorf("failed to get value from redis: %w", err)
	}
	return val, nil
}

func Set(ctx context.Context, redisKey string, value string, rdb *redis.Client) error {
	err := rdb.Set(ctx, redisKey, value, 10*time.Minute).Err()
	if err != nil {
		log.Error("Failed to set value in Redis: %v", err)
		return fmt.Errorf("failed to set value in redis: %w", err)
	}
	return nil
}
