package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) // Default is unauthorized
		log.Println(err)
		fmt.Fprintf(w, err.Error())
	} else if token.Valid {
		tokenInfo, err := json.Marshal(token.Claims)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println(err)
			fmt.Fprintf(w, "error parsing marshalling JSON")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(tokenInfo))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println(err)
		fmt.Fprintf(w, "unable to validate token")
	}
}
