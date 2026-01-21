package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func NewRedisClient(ctx context.Context) (*redis.Client, error) {

	if rdb != nil {
		return rdb, nil
	}

	redisURI := Load().RedisURI

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return rdb, nil
}
