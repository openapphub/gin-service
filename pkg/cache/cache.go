package cache

import (
	"context"
	"openapphub/internal/util"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:       os.Getenv("REDIS_ADDR"),
		Password:   os.Getenv("REDIS_PW"),
		DB:         int(db),
		MaxRetries: 1,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		util.Log().Panic("连接Redis不成功", err)
	}

	RedisClient = client
}

const cachePrefix = "cache_"

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, cachePrefix+key, value, expiration).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, cachePrefix+key).Result()
}

func Del(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, cachePrefix+key).Err()
}
