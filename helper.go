package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func not_modified(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusNotModified)
	w.Header().Set("Content-Type", "application/json")
}

func ok_request(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, message)
}

func created_request(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, message)
}

func access_denied(w http.ResponseWriter, err error, message string) {
	log.Println(err, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, message)
}

func bad_request(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
}

func forbidden_request(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
}

func not_found(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
}

func gone(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusGone)
	w.Header().Set("Content-Type", "application/json")
}

func invalid_request(w http.ResponseWriter, err error, message string) {
	if err == nil {
		err = errors.New(message)
	}
	log.Println(err, message)
	http.Error(w, message, 422)
	w.Header().Set("Content-Type", "application/json")
}

func serveError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
}

func redirectHandler(path string) func(http.ResponseWriter, *http.Request) {
	// http://stackoverflow.com/a/9936937/588759
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)
	}
	// usage: http.HandleFunc("/1", redirectHandler("/one"))
}

func parseForm(r *http.Request) error {
	if r.Form != nil {
		return errors.New("Request body was parsed already.")
	}
	tmp := r.URL.RawQuery
	r.URL.RawQuery = ""
	if err := r.ParseForm(); err != nil {
		return err
	}
	r.URL.RawQuery = tmp
	return nil
}
