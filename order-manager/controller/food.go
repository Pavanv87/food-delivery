package controller

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"

	"order-manager/entities"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DetailsRequestPayload struct {
	Quantity int              `json:"quantity"`
	Address  entities.Address `json:"address"`
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
			http.Error(w, "Cannot find Food Item", 500)
			return
		}

		var payload DetailsRequestPayload
		json.NewDecoder(r.Body).Decode(&payload)

		var detail entities.FoodOrderDetail
		detail.FoodItem = foodItem
		detail.Quantity = payload.Quantity
		detail.BaseCost = foodItem.Cost * payload.Quantity // Rate * Quantity
		detail.Tax = float32(detail.BaseCost) * Tax

		result2 := database.Collection("restaurant").FindOne(ctx, bson.M{"name": foodItem.Restaurant})
		var restaurant entities.Restaurant
		err = result2.Decode(&restaurant)
		if err != nil {
			http.Error(w, "Cannot find Restaurant", 500)
			return
		}
		dist := calculateDistance(payload.Address.Coordinates, restaurant.Address.Coordinates)
		detail.Address = payload.Address
		detail.DeliveryCharge = dist
		detail.PreparationTime = foodItem.PreparationTime * payload.Quantity
		detail.DeliveryTime = (dist * 60) / 40 //mins
		json.NewEncoder(w).Encode(detail)
	}
}

func calculateDistance(custCoord, restCoord entities.Coordinates) int {
	first := math.Pow(float64(restCoord.X-custCoord.X), 2)
	second := math.Pow(float64(restCoord.Y-custCoord.Y), 2)
	return int(math.Sqrt(first + second))
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
