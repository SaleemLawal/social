package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s path: %s error: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request error: %s path: %s error: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found error: %s path: %s error: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "The resource you are looking for might have been removed, had its name changed, or is temporarily unavailable.")
}