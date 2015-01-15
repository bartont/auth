package main

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		access_denied(w, err, err.Error())
	} else if token.Valid {
		tokenInfo, err := json.Marshal(token.Claims)
		if err != nil {
			access_denied(w, err, "error parsing marshalling JSON")
		}
		ok_request(w, string(tokenInfo))
	} else {
		access_denied(w, err, "unable to validate token")
	}
}
