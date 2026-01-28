package helper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateRefreshToken(ctx context.Context, userId, refreshKey string) (string, error) {

	expirationTime := time.Now().Add(24 * 7 * time.Hour)

	claims := &Claim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(refreshKey))
	if err != nil {
		fmt.Println("refreshToken", err)
		return "", err
	}

	return tokenString, nil
}

func GenerateAccessToken(ctx context.Context, userId, accessKey string) (string, error) {

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(accessKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateKey(tokenString, key string) (*Claim, error) {
	claims := &Claim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
			}
			return []byte(key), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}
