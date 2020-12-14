package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"order-manager/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNewDeliveriesHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		filter := bson.D{primitive.E{Key: "status", Value: entities.Prepared}}
		cur, err := database.Collection("order").Find(ctx, filter, options.Find())
		defer cur.Close(ctx)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var results []*entities.Order
		for cur.Next(ctx) {

			// create a value into which the single document can be decoded
			var elem entities.Order
			err := cur.Decode(&elem)
			if err != nil {
				log.Println(err)
			}

			results = append(results, &elem)
		}
		if err := cur.Err(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(results)
	}
}

func GetDeliveryUpdateHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		type OrderUpdate struct {
			Id     string `json:"id"`
			Status string `json:"status"`
		}
		var order OrderUpdate
		json.NewDecoder(r.Body).Decode(&order)

		id, err := primitive.ObjectIDFromHex(order.Id)
		if err != nil {
			log.Println("Invalid id")
		}

		// TODO: after preparing, Status cannot be updated by other restaurants
		resturantName := r.Context().Value(0).(string)

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "status", Value: order.Status},
			primitive.E{Key: "restaurant", Value: resturantName}}}}

		_, err = database.Collection(Order).UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		response := map[string]interface{}{"message": "Success"}
		json.NewEncoder(w).Encode(response)
	}
}
