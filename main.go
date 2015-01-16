package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	rootUrl    string
	urlMongo   string
	port       string
	privateKey []byte
	publicKey  []byte
)

type Registration struct {
	Email    string
	Password string
	UUID     string
}

func init() {
	if pk, err := ioutil.ReadFile(os.Getenv("UTX_PRIVATE_KEY")); err != nil {
		log.Fatal("Unable to read private key", err)
	} else {
		privateKey = pk
	}
	if pbk, err := ioutil.ReadFile(os.Getenv("UTX_PUBLIC_KEY")); err != nil {
		log.Fatal("Unable to read public key", err)
	} else {
		publicKey = pbk
	}
	rootUrl = os.Getenv("UTX_ROOT_URL_AUTH")
	urlMongo = os.Getenv("UTX_URL_MONGO")
	port = ":" + os.Getenv("UTX_PORT_AUTH")

	loadRoutes()
}

func loadRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/validate", validateHandler).Methods("PUT")
	r.HandleFunc("/token", tokenHandler).Methods("POST")
	r.HandleFunc("/registration", registrationHandler).Methods("POST")
	http.Handle("/", r)
}

func main() {
	log.Println("Listening on port " + port + ". Go to http://localhost" + port)
	log.Fatalf("ListenAndServe: %v", http.ListenAndServe(port, nil))
}
