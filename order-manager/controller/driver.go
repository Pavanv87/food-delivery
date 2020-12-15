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

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"createTime", 1}}) // sort order by ascending order
		filter := bson.D{primitive.E{Key: "status", Value: entities.Prepared}}
		cur, err := database.Collection("order").Find(ctx, filter, findOptions) // returns all orders ready to be delivered
		defer cur.Close(ctx)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var results []*entities.Order
		for cur.Next(ctx) {

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
		json.NewEncoder(w).Encode(results) // TODO select first-in/oldest food-order for priority delivery, instead of returning all
	}
}

func GetDeliveryUpdateHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		type OrderUpdate struct { // food order to be updated
			Id     string `json:"id"`
			Status string `json:"status"`
		}
		var order OrderUpdate
		json.NewDecoder(r.Body).Decode(&order)

		id, err := primitive.ObjectIDFromHex(order.Id)
		if err != nil {
			log.Println("Invalid id")
		}

		// TODO: Driver can set only {PICKED, DELIVERED} Status
		// TODO restrict status being updated by other drivers/users
		driverName := r.Context().Value(0).(string)

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "status", Value: order.Status},
			primitive.E{Key: "driver", Value: driverName}}}}

		_, err = database.Collection(Order).UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		response := map[string]interface{}{"message": "Success"}
		json.NewEncoder(w).Encode(response)
	}
}
