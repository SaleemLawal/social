package main

import (
	"net/http"
)

func (app *application) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: get user ID from context
	feeds, err := app.store.Posts.GetFeeds(r.Context(), int64(6))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feeds); err != nil {
		app.internalServerError(w, r, err)
	}
}
