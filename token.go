package main

import (
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

func CheckCredentials(r *http.Request) bool {
	if r.Body == nil {
		return false
	}

	var dat map[string]interface{}
	json.NewDecoder(r.Body).Decode(&dat)

	if dat["password"] == nil || dat["email"] == nil {
		return false
	}

	password := dat["password"].(string)
	email := dat["email"]
	passwordBytes := []byte(password)

	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("unable to connect to mongo")
		return false
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	connection := session.DB("UTx").C("registrations")
	result := Registration{}
	err = connection.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		fmt.Println("unable to find user")
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), passwordBytes)
	if err != nil {
		// password is not a match
		fmt.Println("passwords not a match")
		return false
	} else {
		return true
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if CheckCredentials(r) {
		w.Header().Set("Content-Type", "application/json")
		token := jwt.New(jwt.GetSigningMethod("RS256")) // Create a Token that will be signed with RSA 256.

		exp := time.Now().Unix() + (60 * 60) // 1 hour
		etm := time.Unix(exp, 0)

		claims := token.Claims
		claims["exp"] = exp
		claims["created"] = time.Now()
		claims["expires"] = etm

		// The claims object allows you to store information in the actual token.
		tokenString, err := token.SignedString(privateKey)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, err.Error())
			log.Println("Unable to get token", err)
			return
		}

		claims["token"] = tokenString
		session, _ := json.Marshal(claims)
		// tokenString Contains the actual token you should share with your client.
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, string(session))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "invalid username or password")
	}
}
