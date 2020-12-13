package main

import (
	"context"
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
	router.HandleFunc("/food/order", controller.GetOrderHandler(ctx, database)).Methods("POST")

	fmt.Println("Auth Service Starting...")
	http.ListenAndServe(":8080", router)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		next.ServeHTTP(w, r) // call original
		log.Println("After")
	})
}
