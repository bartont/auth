package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func checkCredentials(r *http.Request) (string, string, bool) {
	if r.Body == nil {
		return "", "", false
	}

	var dat map[string]interface{}
	json.NewDecoder(r.Body).Decode(&dat)

	if dat["password"] == nil || dat["email"] == nil {
		return "", "", false
	}

	password := dat["password"].(string)
	email := dat["email"]
	passwordBytes := []byte(password)

	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("unable to connect to mongo")
		return "", "", false
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	connection := session.DB("UTx").C("registrations")
	result := Registration{}

	if err = connection.Find(bson.M{"email": email}).One(&result); err != nil {
		fmt.Println("unable to find user")
		return "", "", false
	}

	if err = bcrypt.CompareHashAndPassword([]byte(result.Password), passwordBytes); err != nil {
		// password is not a match
		fmt.Println("passwords not a match")
		return "", "", false
	} else {
		return result.UUID, result.Email, true
	}
}
