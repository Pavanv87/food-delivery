package main

import (
	"auth/db"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {

	// Obtain Token
	data, _ := json.Marshal(&User{Username: "pavan", Password: "abc123"})
	req, err := http.NewRequest("POST", "/auth/user/signin", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	// Not Recommended, probably use dockertest
	conf := db.Config{Host: "localhost", Port: "27017", Username: "admin", Password: "password", Database: "food-delivery"}

	ctx := context.TODO()
	mClient := conf.NewClient(ctx) // Mongo client

	database := mClient.Database(conf.Database)

	rr := httptest.NewRecorder()
	http.HandlerFunc(GetSignInHandler(ctx, database)).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK { // Check status code.
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	jwt := &JwtToken{}
	err = json.Unmarshal(rr.Body.Bytes(), jwt)
	if err != nil || jwt.Token == "" { // Check the response body is what we expecfor token.
		t.Error("handler returned empty response", rr.Body.String())
	}

	// Token validation
	req, err = http.NewRequest("GET", "/auth/validate", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt.Token+"123")

	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler)
	handler.ServeHTTP(rr, req)

	// Check status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
