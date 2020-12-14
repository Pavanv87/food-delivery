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

const (
	Order string = "order"
)

func GetOrderPlaceHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload entities.FoodOrderDetail
		json.NewDecoder(r.Body).Decode(&payload)

		var order entities.Order
		order.FoodOrderDetail = payload
		order.Status = entities.New
		order.Customer = r.Context().Value(0).(string)

		result, err := database.Collection(Order).InsertOne(ctx, order, options.InsertOne())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		response := map[string]interface{}{"id": result.InsertedID.(primitive.ObjectID).Hex()}
		json.NewEncoder(w).Encode(response)
	}
}

func GetOrderCancelHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		type OrderId struct {
			Id string `json:"id"`
		}
		var orderId OrderId
		json.NewDecoder(r.Body).Decode(&orderId)

		id, err := primitive.ObjectIDFromHex(orderId.Id)
		if err != nil {
			log.Println("Invalid id")
		}

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "status", Value: entities.Cancelled}}}}

		_, err = database.Collection(Order).UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		response := map[string]interface{}{"message": "Success"}
		json.NewEncoder(w).Encode(response)
	}
}
