package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type FoodItem struct {
	Id              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Cost            int                `json:"cost" bson:"cost"`
	PreparationTime int                `json:"preparationTime" bson:"preparationTime"` //In Minutes
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
	Id   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

type FoodOrderDetail struct {
	FoodItem        FoodItem `json:"foodItem"`
	BaseCost        int      `json:"baseCost"`
	Tax             float32  `json:"tax"`
	Quantity        int      `json:"quantity"`
	DeliveryCharge  int      `json:"deliveryCharge"`
	PreparationTime int      `json:"preparationTime"`
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
	Status          OrderStatus        `json:"status" bson:"status"`
	Customer        string             `json:"customer" bson:"customer"`
	FoodOrderDetail FoodOrderDetail    `json:"foodOrderDetail" bson:"foodOrderDetail"`
	RestaurantId    primitive.ObjectID `json:"restaurantId" bson:"restaurantId"`
}

type Delivery struct {
	Status DeliveryStatus `json:"status" bson:"status"`
	Order  Order          `json:"order" bson:"order"`
	Driver Driver         `json:"driver" bson:"driver"`
}
