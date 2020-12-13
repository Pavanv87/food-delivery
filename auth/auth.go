package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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
	Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/auth/signin", signInHandler).Methods("POST")
	router.HandleFunc("/auth/validate", authHandler).Methods("GET")

	fmt.Println("Auth Service Starting...")
	http.ListenAndServe(":8080", router)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
	})
	tokenString, error := token.SignedString([]byte(SecretKey))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
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
		json.NewEncoder(w).Encode(Response{Message: "Success"})
	}
}
