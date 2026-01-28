package services

import (
	"context"
	"time"

	"github.com/go-auth-service/src/config"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func SaveRefreshToken(ctx context.Context, userId, refreshToken string) (string, error) {

	client, err := config.NewRedisClient(ctx)
	if err != nil {
		return "", err
	}

	key := "user:session:" + userId
	sessionId, err := gonanoid.New()
	if err != nil {
		return "", err
	}
	err = client.HSet(ctx, key, sessionId, refreshToken, 7*24*time.Hour).Err()

	return sessionId, nil
}

func GetRefreshToken(ctx context.Context, sessionId, userId string) (string, error) {
	client, err := config.NewRedisClient(ctx)
	if err != nil {
		return "", err
	}
	key := "user:session:" + userId

	return client.HGet(ctx, key, sessionId).Result()
}

func DeleteRefreshToken(ctx context.Context, sessionId, userId string) error {
	client, err := config.NewRedisClient(ctx)
	if err != nil {
		return err
	}
	key := "user:session:" + userId
	return client.Del(ctx, key, sessionId).Err()
}

func LogOutFromMultipleDevice(ctx context.Context, userId string) error {
	client, err := config.NewRedisClient(ctx)
	if err != nil {
		return err
	}
	key := "user:session:" + userId
	return client.Del(ctx, key).Err()
}
