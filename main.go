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

// set up rsa keys for authentication tokens
func init() {
	pk, err := ioutil.ReadFile(os.Getenv("UTX_PRIVATE_KEY"))
	if err != nil {
		log.Fatal("Unable to read private key", err)
	} else {
		privateKey = pk
	}
	pbk, err := ioutil.ReadFile(os.Getenv("UTX_PUBLIC_KEY"))
	if err != nil {
		log.Fatal("Unable to read public key", err)
	} else {
		publicKey = pbk
	}
}

// start the server
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
