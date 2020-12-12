package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {

	// Obtain Token
	data, _ := json.Marshal(&User{Username: "pavan", Password: "abc123"})
	req, err := http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(signInHandler).ServeHTTP(rr, req)

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
