package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		invalid_request(w, nil, "missing request body")
		return
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		invalid_request(w, err, "user not created, database error")
		return
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	var dat map[string]interface{}
	json.NewDecoder(r.Body).Decode(&dat)

	if dat["password"] == nil || dat["email"] == nil {
		invalid_request(w, nil, "missing email or password")
		return
	}

	password := dat["password"].(string)
	passwordBytes := []byte(password)
	email := dat["email"]

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, 10)

	if err != nil {
		invalid_request(w, err, "user not created, unable to generate password hash")
		return
	}

	c := session.DB("UTx").C("registrations")

	var reg = Registration{
		Email:    email.(string),
		Password: string(hashedPassword),
	}

	upsertdata := bson.M{"$set": reg}

	info, err := c.UpsertId(reg.Email, upsertdata)

	if err != nil {
		invalid_request(w, err, fmt.Sprintf("user not created, unable to upsert, %v, upsertId: %v", err, info))
		return
	}

	result := Registration{}
	err = c.Find(bson.M{"email": email}).One(&result)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "user not created, unable to find user, %v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), passwordBytes)
	if err != nil {
		// password is not a match
		access_denied(w, err, fmt.Sprintf("user not created, password not a match, %v", err))

	} else {
		session, err := GetNewToken(w)
		if err != nil {
			access_denied(w, err, "unable to generate token")
		} else {
			created_request(w, string(session))
		}
	}

}
