package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"order-manager/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetItemsHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		cur, err := database.Collection("foodItem").Find(ctx, bson.D{{}}, options.Find())
		defer cur.Close(ctx)

		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), 500)
		}

		var results []*entities.FoodItem
		for cur.Next(ctx) {

			// create a value into which the single document can be decoded
			var elem entities.FoodItem
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}

			results = append(results, &elem)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), 500)
		}
		json.NewEncoder(w).Encode(results)
	}
}
