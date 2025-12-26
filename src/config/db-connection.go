package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context) (*mongo.Client, error) {
	log.Println("Connecting to DB....")

	cfg := Load()
	if cfg.MongoURI == "" {
		log.Fatal("MongoURI is not set!")
	}

	var err error
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Mongo connection error; ", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Mongo ping failed: ", err)
		return nil, err
	}

	return client, nil
}
