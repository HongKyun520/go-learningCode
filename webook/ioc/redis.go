package ioc

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v", err)
	}

	return rdb
}
