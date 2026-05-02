package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "The resource you are looking for might have been removed, had its name changed, or is temporarily unavailable.")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
}

func (app *application) unauthorizedBasicAuthError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Unauthorized basic auth error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "Unauthorized. Please provide valid credentials.")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("Forbidden error", "method", r.Method, "path", r.URL.Path, "error")
	writeJSONError(w, http.StatusForbidden, "Forbidden")
}
