package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/saleemlawal/social/internal/store"
)


type contextKey string

const postContextKey contextKey = "post"

type CreatePostPayload struct {
	Title string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=10,max=1000"`
	Tags []string `json:"tags" validate:"required,min=1,max=5"`
}

type UpdatePostPayload struct {
	Title *string `json:"title" validate:"omitempty,min=3,max=255"`
	Content *string `json:"content" validate:"omitempty,min=10,max=1000"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	userId := 1

	post := &store.Post{
		Title: payload.Title,
		Content: payload.Content,
		Tags: payload.Tags,
		// TODO: Change after authentication
		UserID: int64(userId),
	}
	

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post, _ := getPostFromCtx(r)
	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post, _ := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postId := chi.URLParam(r, "postId")

	id, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		postId := chi.URLParam(r, "postId")

		id, err := strconv.ParseInt(postId, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		post, err := app.store.Posts.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrRecordNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postContextKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) (*store.Post, bool) {
	post, ok := r.Context().Value(postContextKey).(*store.Post)
	return post, ok
}