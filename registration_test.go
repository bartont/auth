package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistration(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/registration", registrationHandler).Methods("POST")

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	jsonStr := []byte(`{"email":"foo@you.com", "password":"zaq12wsx"}`)
	b := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", "/registration", b)
	if err != nil {
		t.Fatal("'POST /registration' request failed!")
	}

	r.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusCreated)
	}
}
