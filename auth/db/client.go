package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//github.com/google/wire
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func (dbConf *Config) NewClient(ctx context.Context) *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s/%s?connect=direct", dbConf.Host, dbConf.Port, dbConf.Database)

	// Set client options
	clientOptions := options.Client().
		SetAuth(options.Credential{Username: dbConf.Username, Password: dbConf.Password}).
		ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}
