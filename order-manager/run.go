package main

import (
	"context"
	"order-manager/db"
)

func main() {
	conf := db.Config{Host: "localhost", Port: "27017", Username: "admin", Password: "password", Database: "food-delivery"}

	ctx := context.TODO()
	mClient := conf.NewClient(ctx) // Mongo client
	db.Migrate(conf.Database, mClient)
}
