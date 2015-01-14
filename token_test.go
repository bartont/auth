package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToken(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/token", tokenHandler).Methods("POST")

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	jsonStr := []byte(`{"email":"foo@you.com", "password":"zaq12wsx"}`)
	b := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", "/token", b)
	if err != nil {
		t.Fatal("'POST /token' request failed!")
	}

	r.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusCreated)
	}
}
