package main

import (
	"net/http"

	"github.com/saleemlawal/social/internal/store"
)

// getFeedsHandler godoc
//
//	@Summary		Get user feed
//	@Description	Returns a paginated feed of posts from the user and users they follow
//	@Tags			feed
//	@Produce		json
//	@Param			limit	query		int		false	"Number of posts to return"	minimum(1)			maximum(20)	default(20)
//	@Param			offset	query		int		false	"Number of posts to skip"	minimum(0)			default(0)
//	@Param			sort	query		string	false	"Sort order by created_at"	Enums(asc, desc)	default(desc)
//	@Param			tags	query		string	false	"Comma-separated list of tags to filter by"
//	@Param			search	query		string	false	"Search term to filter by title or content"
//	@Param			since	query		string	false	"Return posts created after this time (RFC3339)"
//	@Param			until	query		string	false	"Return posts created before this time (RFC3339)"
//	@Success		200		{array}		store.Feed
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feeds [get]
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
