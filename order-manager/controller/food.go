package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"order-manager/entities"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DetailsRequestPayload struct {
	Quantity int `json:"quantity"`
	Distance int `json:"distance"` // In Kms
}

const Tax = 0.05

func GetOrderDetailHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := primitive.ObjectIDFromHex(vars["id"])
		if err != nil {
			log.Println("Invalid id")
		}
		result := database.Collection("foodItem").FindOne(ctx, bson.M{"_id": id})
		var foodItem entities.FoodItem
		err = result.Decode(&foodItem)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var payload DetailsRequestPayload
		json.NewDecoder(r.Body).Decode(&payload)

		var detail entities.FoodOrderDetail
		detail.FoodItem = foodItem
		detail.Quantity = payload.Quantity
		detail.BaseCost = foodItem.Cost * payload.Quantity // Rate * Quantity
		detail.Tax = float32(detail.BaseCost) * Tax
		detail.DeliveryCharge = payload.Distance
		detail.PreparationTime = foodItem.PreparationTime * payload.Quantity
		detail.DeliveryTime = (payload.Distance * 60) / 40 //mins
		json.NewEncoder(w).Encode(detail)
	}
}

func GetItemsHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		cur, err := database.Collection("foodItem").Find(ctx, bson.D{{}}, options.Find())
		defer cur.Close(ctx)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var results []*entities.FoodItem
		for cur.Next(ctx) {

			// create a value into which the single document can be decoded
			var elem entities.FoodItem
			err := cur.Decode(&elem)
			if err != nil {
				log.Println(err)
			}
			log.Println("ID: " + elem.Id.Hex())

			results = append(results, &elem)
		}
		if err := cur.Err(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(results)
	}
}
