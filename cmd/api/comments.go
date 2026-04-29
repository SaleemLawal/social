package main

import (
	"net/http"

	"github.com/saleemlawal/social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,min=10,max=1000"`
	Likes   int    `json:"likes" validate:"omitempty,min=0"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post, ok := getPostFromCtx(r)

	if !ok {
		app.notFoundError(w, r, store.ErrRecordNotFound)
		return
	}

	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: payload.Content,
		Likes:   payload.Likes,
		PostID:  post.ID,
		UserID:  1, // TODO: Change after authentication
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
