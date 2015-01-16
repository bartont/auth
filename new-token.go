package main

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func GetNewToken(w http.ResponseWriter, uuid string, email string) (string, error) {
	// Create a Token that will be signed with RSA 256.
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	exp := time.Now().Unix() + (60 * 60) // 1 hour
	etm := time.Unix(exp, 0)

	claims := token.Claims
	claims["exp"] = exp
	claims["created"] = time.Now()
	claims["expires"] = etm
	claims["uuid"] = uuid
	claims["email"] = email

	tokenString, err := token.SignedString(privateKey)
	// The claims object allows you to store information in the actual token.
	if err != nil {
		return "", err
	}

	claims["token"] = tokenString
	session, err := json.Marshal(claims)
	return string(session), err
}
