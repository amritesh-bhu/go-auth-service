package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context) (*redis.Client, error) {

	redisURI := Load().RedisURI

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return rdb, redis.Nil
}
