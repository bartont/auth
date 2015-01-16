package main

import (
	"net/http"
)

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if uuid, email, isValid := checkCredentials(r); isValid {
		// tokenString Contains the actual token you should share with your client.
		session, err := GetNewToken(w, uuid, email)
		if err != nil {
			access_denied(w, err, "unable to generate token")
		} else {
			created_request(w, string(session))
		}

	} else {
		access_denied(w, nil, "invalid username or password")
	}
}
