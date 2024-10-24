package cache

import (
	"context"
	"fmt"
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
		util.Log().Panic("连接Redis不成功")
	}

	RedisClient = client
}

const cachePrefix = ""

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, cachePrefix+key, value, expiration).Err()
}

func Exists(ctx context.Context, key string) (bool, error) {
	n, err := RedisClient.Exists(ctx, cachePrefix+key).Result()
	return n > 0, err
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RedisClient.Expire(ctx, cachePrefix+key, expiration).Err()
}

func Del(ctx context.Context, key string) error {
	count, err := RedisClient.Del(ctx, cachePrefix+key).Result()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("key %s does not exist", key)
	}

	return nil
}

func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, cachePrefix+key).Result()
}

// New function to delete keys by prefix
func DelByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var deletedKeys int
	for {
		var keys []string
		var err error
		keys, cursor, err = RedisClient.Scan(ctx, cursor, cachePrefix+prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := RedisClient.Del(ctx, keys...).Err(); err != nil {
				return err
			}
			deletedKeys += len(keys)
		}
		if cursor == 0 {
			break
		}
	}
	if deletedKeys == 0 {
		return fmt.Errorf("no keys found matching prefix: %s", prefix)
	}
	return nil
}
