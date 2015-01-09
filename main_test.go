package main

import (
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

	req, err := http.NewRequest("POST", "/token", nil)
	if err != nil {
		t.Fatal("'POST /token' request failed!")
	}

	r.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusCreated {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusCreated)
	}
}

func TestValidateFails(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/validate", validateHandler).Methods("PUT")

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	req, err := http.NewRequest("PUT", "/validate", nil)
	if err != nil {
		t.Fatal("'PUT /validateHandler' request failed!")
	}

	r.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusUnauthorized {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusCreated)
	}
}
