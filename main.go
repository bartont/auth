package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	port = ":8000"
)

var (
	privateKey []byte
	publicKey  []byte
)

type Registration struct {
	Email    string
	Password string
}

func init() {
	pk, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/demo.rsa")
	if err != nil {
		log.Fatal("Unable to read private key", err)
	} else {
		privateKey = pk
	}
	pbk, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/demo.rsa.pub")
	if err != nil {
		log.Fatal("Unable to read public key", err)
	} else {
		publicKey = pbk
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/validate", validateHandler).Methods("PUT")
	r.HandleFunc("/token", tokenHandler).Methods("POST")
	r.HandleFunc("/registration", registrationHandler).Methods("POST")
	http.Handle("/", r)

	log.Println("Listening on port 8000. Go to http://localhost:8000/")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
