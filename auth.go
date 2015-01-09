package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	privateKey []byte
	publicKey  []byte
)

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

	r.HandleFunc("/validate", ValidateHandler).Methods("PUT")
	r.HandleFunc("/token", TokenHandler).Methods("POST")
	http.Handle("/", r)

	log.Println("Listening on port 8000. Go to http://localhost:8000/")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := jwt.New(jwt.GetSigningMethod("RS256")) // Create a Token that will be signed with RSA 256.

	exp := time.Now().Unix() + (60 * 60) // 1 hour
	etm := time.Unix(exp, 0)

	info := token.Claims
	info["id"] = "This is my super fake ID"
	info["exp"] = exp
	info["created"] = time.Now()
	info["expires"] = etm

	// The claims object allows you to store information in the actual token.
	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, err.Error())
		log.Println("Unable to get token", err)
		return
	}

	info["token"] = tokenString
	session, _ := json.Marshal(info)
	// tokenString Contains the actual token you should share with your client.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(session))
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	fmt.Println(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) // Default is unauthorized
		log.Println(err)
		fmt.Fprintf(w, err.Error())
	} else if token.Valid {
		w.WriteHeader(http.StatusOK)
		tokenInfo, _ := json.Marshal(token.Claims)
		fmt.Fprintf(w, string(tokenInfo))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		fmt.Fprintf(w, "unable to validate token")
	}
}
