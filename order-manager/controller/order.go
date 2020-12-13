package controller

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetOrderHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
