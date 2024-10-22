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
		util.Log().Panic("连接Redis不成功")
	}

	RedisClient = client
}

const cachePrefix = "cache_"

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, cachePrefix+key, value, expiration).Err()
}

var usingRedis bool

// func InitCache(useRedis bool) {
// 	usingRedis = useRedis
// 	if useRedis {
// 		initRedis()
// 	} else {
// 		initMemory()
// 	}
// }

func IsUsingRedis() bool {
	return usingRedis
}
func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, cachePrefix+key).Result()
}

func Del(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, cachePrefix+key).Err()
}

// New function to delete keys by prefix
func DelByPrefix(ctx context.Context, prefix string) error {
	iter := RedisClient.Scan(ctx, 0, cachePrefix+prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		err := RedisClient.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

// New function to check if a key exists
func Exists(ctx context.Context, key string) (bool, error) {
	exists, err := RedisClient.Exists(ctx, cachePrefix+key).Result()
	return exists == 1, err
}

// New function to set the expiration of a key
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RedisClient.Expire(ctx, cachePrefix+key, expiration).Err()
}
