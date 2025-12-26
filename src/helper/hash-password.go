package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hash, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(plainPassword),
	)

	return err == nil
}
