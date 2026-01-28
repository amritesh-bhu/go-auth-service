package helper

import (
	"context"

	"github.com/go-auth-service/src/config"
	"github.com/go-auth-service/src/services"
)

func GetTokens(ctx context.Context, userId string) (string, string, error) {
	refreshKey := config.Load().RefreshSecret
	refreshToken, err := GenerateRefreshToken(ctx, userId, refreshKey)
	if err != nil {
		return "", "", err
	}

	sessionId, err := services.SaveRefreshToken(ctx, refreshToken, userId)
	if err != nil {
		return "", "", err
	}

	accessKey := config.Load().AccessSecret
	accessToken, err := GenerateAccessToken(ctx, userId, accessKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, sessionId, nil
}
