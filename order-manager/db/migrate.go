package db

import (
	"bytes"
	"log"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(dbName string, client *mongo.Client) {

	d, err := mongodb.WithInstance(client, &mongodb.Config{
		DatabaseName:    dbName,
		TransactionMode: false,
	})

	cmds := []byte(`[
		{ "create":"order"},
		{ "create":"user"},
		{ "createIndexes": "user", "indexes": [{"key": {"name": 1},"name": "unique_user_name","unique": true}] },
		{ "create":"restaurant"},
		{ "createIndexes": "restaurant", "indexes": [{"key": {"name": 1},"name": "unique_restaurant_name","unique": true}] },
		{ "create":"foodItem"},
		{ "createIndexes": "foodItem", "indexes": [{"key": {"name": 1},"name": "unique_foodItem_name","unique": true}] },
		{ "create":"driver"},
		{ "createIndexes": "driver", "indexes": [{"key": {"name": 1},"name": "unique_driver_name","unique": true}] },
		{
		  "insert":"user",
		  "documents": [{"name":"Customer1", "password":"pass1"},{"name":"Customer2", "password":"pass2"},{"name":"Customer3", "password":"pass3"}]
		},
		{
		  "insert":"restaurant",
		  "documents": [{"name":"Restaurant1"},{"name":"Restaurant2"},{"name":"Restaurant3"}]
		},
		{
		  "insert":"foodItem",
		  "documents": [{"name":"FoodItem1","cost":50, "preparationTime":15},{"name":"FoodItem2","cost":70, "preparationTime":20},{"name":"FoodItem3","cost":100, "preparationTime":30}]
		},
		{
		  "insert":"driver",
		  "documents": [{"name":"Driver1"},{"name":"Driver2"},{"name":"Driver3"}]
		}
	]`)

	err = d.Run(bytes.NewReader(cmds))
	if err != nil {
		log.Println(err)
	}
}
