package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getAuthHeader(t *testing.T) string {
	r := mux.NewRouter()
	r.HandleFunc("/token", tokenHandler).Methods("POST")

	respRec := httptest.NewRecorder()

	jsonStr := []byte(`{"email":"foo@you.com", "password":"zaq12wsx"}`)
	b := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest("POST", "/token", b)
	if err != nil {
		t.Fatal("'POST /token' request failed!")
	}

	r.ServeHTTP(respRec, req)

	var dat map[string]interface{}
	json.NewDecoder(respRec.Body).Decode(&dat)

	return "Bearer " + dat["token"].(string)
}

func TestValidate(t *testing.T) {
	authorizationHeader := getAuthHeader(t)

	r := mux.NewRouter()
	r.HandleFunc("/validate", validateHandler).Methods("PUT")

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	req, err := http.NewRequest("PUT", "/validate", nil)

	req.Header.Set("Authorization", authorizationHeader)
	if err != nil {
		t.Fatal("'PUT /validateHandler' request failed!")
	}

	r.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", respRec.Code, " instead of ", http.StatusOK)
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

func TestCleanUp(t *testing.T) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("unable to connect to mongo")
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	connection := session.DB("UTx").C("registrations")
	err = connection.Remove(bson.M{"email": "foo@you.com"})
}
