package controller

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"order-manager/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNewOrdersHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		resturantName := r.Context().Value(0).(string)
		// Get new orders meant for this restaurant
		filter := bson.D{primitive.E{Key: "status", Value: entities.New}, primitive.E{Key: "foodOrderDetail.fooditem.restaurant", Value: resturantName}}
		cur, err := database.Collection("order").Find(ctx, filter, options.Find())
		defer cur.Close(ctx)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		results := make([]*entities.Order, 0) //empty slice
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
		json.NewEncoder(w).Encode(results)
	}
}

func GetOrderUpdateHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
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

		// TODO: In general status cannot rollback
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

func GetOrderInvoiceHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
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

		// Get order for this restaurant
		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		result := database.Collection("order").FindOne(ctx, filter)
		if err = result.Err(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		var order entities.Order
		err = result.Decode(&order)
		if err != nil {
			http.Error(w, "Cannot get order details", 500)
			return
		}
		// Return html template as invoice
		Invoice :=
			`
<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<title>INVOICE</title>
</head>
<body>
	<header>
		<h1>Invoice</h1>
	</header>
{{ .FoodOrderDetail.FoodItem.Restaurant }}

Invoice# 	{{ .Id }}
Date		{{ .CreateTime }}

Item		{{ .FoodOrderDetail.FoodItem.Name }}
</body>			
</html>							
`
		template.Must(template.New("").Parse(Invoice)).Execute(w, order)
	}
}
