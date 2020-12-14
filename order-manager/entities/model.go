package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type FoodItem struct {
	Id              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Cost            int                `json:"cost" bson:"cost"`
	PreparationTime int                `json:"preparationTime" bson:"preparationTime"` //In Minutes
	Restaurant      string             `json:"restaurant" bson:"restaurant"`
}

type Customer struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Password string             `json:"password" bson:"password"`
}
type Driver struct {
	Id   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}
type Restaurant struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Address Address            `json:"address" bson:"address"`
}

type FoodOrderDetail struct {
	FoodItem        FoodItem `json:"foodItem"`
	BaseCost        int      `json:"baseCost"`
	Tax             float32  `json:"tax"`
	Quantity        int      `json:"quantity"`
	DeliveryCharge  int      `json:"deliveryCharge"`
	PreparationTime int      `json:"preparationTime"`
	Address         Address  `json:"address"`
	DeliveryTime    int      `json:"deliveryTime"`
}

type OrderStatus string
type DeliveryStatus string

const (
	New       OrderStatus    = "NEW"
	Preparing OrderStatus    = "PREPARING"
	Prepared  OrderStatus    = "PREPARED"
	Cancelled OrderStatus    = "CANCELLED"
	Picked    DeliveryStatus = "PICKED"
	Delivered DeliveryStatus = "DELIVERED"
)

type Order struct {
	Id              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Status          OrderStatus        `json:"status" bson:"status"`
	Customer        string             `json:"customer" bson:"customer"`
	FoodOrderDetail FoodOrderDetail    `json:"foodOrderDetail" bson:"foodOrderDetail"`
	Address         Address            `json:"address" bson:"address"`
	Driver          string             `json:"driver,omitempty" bson:"restaurant"`
}

type Coordinates struct {
	X float32 `json:"x" bson:"x"`
	Y float32 `json:"y" bson:"y"`
}
type Address struct {
	DoorNo      string      `json:"doorNo" bson:"doorNo"`
	Street      string      `json:"street" bson:"street"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
}
