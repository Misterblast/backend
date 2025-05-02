package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/redis/go-redis/v9"
)

const (
	LongEXP     = 24 * time.Hour
	StandardEXP = 10 * time.Minute
	FastEXP     = 5 * time.Minute
	BlazingEXP  = 3 * time.Minute
	InstantEXP  = 1 * time.Minute
)

type ExpirationType string

const (
	ExpFast     ExpirationType = "fast"
	ExpStandard ExpirationType = "standard"
	ExpLong     ExpirationType = "long"
	ExpBlazing  ExpirationType = "blazing"
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
		return "", fmt.Errorf("failed to get value from redis: %w", err)
	}
	return val, nil
}

func Set(ctx context.Context, redisKey string, value string, rdb *redis.Client, expType ExpirationType) error {
	var ttl time.Duration

	switch expType {
	case ExpFast:
		ttl = FastEXP
	case ExpStandard:
		ttl = StandardEXP
	case ExpLong:
		ttl = LongEXP
	case ExpBlazing:
		ttl = BlazingEXP
	default:
		ttl = StandardEXP // fallback default
	}

	err := rdb.Set(ctx, redisKey, value, ttl).Err()
	if err != nil {
		log.Error("Failed to set value in Redis: %v", err)
		return fmt.Errorf("failed to set value in redis: %w", err)
	}
	return nil
}
