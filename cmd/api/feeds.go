package main

import (
	"net/http"

	"github.com/saleemlawal/social/internal/store"
)

func (app *application) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	fq := &store.PaginationFeedsQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
	}

	fq, err := fq.Parse(r)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// TODO: get user ID from context
	feeds, err := app.store.Posts.GetFeeds(r.Context(), int64(6), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feeds); err != nil {
		app.internalServerError(w, r, err)
	}
}
