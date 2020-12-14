package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-manager/controller"
	"order-manager/db"

	"github.com/gorilla/mux"
)

func main() {
	conf := db.Config{Host: "localhost", Port: "27017", Username: "admin", Password: "password", Database: "food-delivery"}

	ctx := context.TODO()
	mClient := conf.NewClient(ctx) // Mongo client
	db.Migrate(conf.Database, mClient)

	database := mClient.Database(conf.Database)

	router := mux.NewRouter()
	router.HandleFunc("/food/items", controller.GetItemsHandler(ctx, database)).Methods("GET")
	router.HandleFunc("/food/{id}/detail", controller.GetOrderDetailHandler(ctx, database)).Methods("POST")

	router.HandleFunc("/customer/order", authMiddleware(controller.GetOrderPlaceHandler(ctx, database))).Methods("POST")
	router.HandleFunc("/customer/order/cancel", authMiddleware(controller.GetOrderCancelHandler(ctx, database))).Methods("POST")

	router.HandleFunc("/restaurant/orders", authMiddleware(controller.GetNewOrdersHandler(ctx, database))).Methods("GET")
	router.HandleFunc("/restaurant/order/update", authMiddleware(controller.GetOrderUpdateHandler(ctx, database))).Methods("POST")

	fmt.Println("Auth Service Starting...")
	log.Fatalf("ListenAndServe Error: %s", http.ListenAndServe(":8082", router).Error())
}

// type contextKey int

// const (
// 	CustomerKey contextKey = iota
// )

func authMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://localhost:8081/auth/verify", nil)
		req.Header = r.Header
		resp, _ := client.Do(req)

		type Auth struct {
			Claim struct {
				Name string `json:"username"`
			} `json:"claim"`
		}
		var auth Auth
		json.NewDecoder(resp.Body).Decode(&auth)
		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Authentication Failed", http.StatusUnauthorized)
			return
		}
		ctxWithUser := context.WithValue(r.Context(), 0, auth.Claim.Name)
		r.WithContext(ctxWithUser)
		http.HandlerFunc(next).ServeHTTP(w, r.WithContext(ctxWithUser))
	}
}
