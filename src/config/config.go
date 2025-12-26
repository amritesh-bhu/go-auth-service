package config

import (
	"log"
	"os"
)

type Config struct {
	Port          string
	MongoURI      string
	AccessSecret  string
	RefreshSecret string
}

func Load() *Config {
	c := &Config{
		Port:          getEnv("PORT"),
		MongoURI:      getEnv("MONGO_URI"),
		AccessSecret:  getEnv("ACCESS_SECRET"),
		RefreshSecret: getEnv("REFRESH_SECRET"),
	}

	// if c.AccessSecret == "" || c.RefreshSecret == "" {
	// log.Println("Warn: Access/Refresh secrets are empty — set them in .env for production")
	// }

	if c.MongoURI == "" || c.Port == "" {
		log.Println("Warn: MONGO_URI/Port are empty — set them in .env for production")
	}

	return c
}

func getEnv(key string) string {
	v := os.Getenv(key)
	return v
}
