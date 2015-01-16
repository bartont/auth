package main

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		access_denied(w, nil, err.Error())
		return
	} else if token.Valid {
		claims := token.Claims

		tokenInfo, err := json.Marshal(token.Claims)
		if err != nil {
			access_denied(w, err, "error parsing marshalling JSON")
			return
		}

		session, err := mgo.Dial(urlMongo)
		if err != nil {
			access_denied(w, err, "unable to find user")
			return
		}
		defer session.Close()

		session.SetMode(mgo.Monotonic, true)

		connection := session.DB("UTx").C("registrations")
		result := Registration{}

		if err = connection.Find(bson.M{"email": claims["email"]}).One(&result); err != nil {
			access_denied(w, err, "unable to find user")
			return
		}

		// check that the token's uuid is the user's uuid
		if claims["uuid"] == result.UUID {
			ok_request(w, string(tokenInfo))
		} else {
			access_denied(w, err, "unable to validate token")
		}
	} else {
		access_denied(w, err, "unable to validate token")
	}
}
