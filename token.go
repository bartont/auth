package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := jwt.New(jwt.GetSigningMethod("RS256")) // Create a Token that will be signed with RSA 256.

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
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, err.Error())
		log.Println("Unable to get token", err)
		return
	}

	info["token"] = tokenString
	session, _ := json.Marshal(info)
	// tokenString Contains the actual token you should share with your client.
	w.WriteHeader(http.StatusCreated)

	log.Println(string(session))
	fmt.Fprintf(w, string(session))
}
