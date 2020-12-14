package main

import (
	"auth/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var SecretKey = []byte("YouAreWelcome")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Response struct {
	Message string     `json:"message"`
	Claim   jwt.Claims `json:"claim"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {}

func main() {

	conf := db.Config{Host: "localhost", Port: "27017", Username: "admin", Password: "password", Database: "food-delivery"}

	ctx := context.TODO()
	mClient := conf.NewClient(ctx) // Mongo client
	database := mClient.Database(conf.Database)

	router := mux.NewRouter()
	router.HandleFunc("/auth/{collection}/signin", GetSignInHandler(ctx, database)).Methods("POST")
	router.HandleFunc("/auth/verify", authHandler).Methods("GET")

	fmt.Println("Auth Service Starting...")
	http.ListenAndServe(":8081", router)
}

func GetSignInHandler(ctx context.Context, database *mongo.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)

		vars := mux.Vars(r)
		collection := vars["collection"]
		result := database.Collection(collection).FindOne(ctx, bson.M{"name": user.Username})
		var dbUser User
		err := result.Decode(&dbUser)
		if err != nil {
			http.Error(w, "User Not Found", 500)
			return
		}

		if dbUser.Password != user.Password {
			http.Error(w, "UserName or Password mismatch", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
		})
		tokenString, err := token.SignedString([]byte(SecretKey))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var token string

	// Bearer token from Authorization header
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}

	if token == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// parse token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, msg
		}
		return SecretKey, nil
	})
	if err != nil {
		http.Error(w, "Error parsing token", http.StatusUnauthorized)
		return
	}
	// Check token is valid
	if parsedToken != nil && parsedToken.Valid {
		json.NewEncoder(w).Encode(Response{Message: "Success", Claim: parsedToken.Claims})
	}
}
