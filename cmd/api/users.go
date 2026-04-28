package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/saleemlawal/social/internal/store"
)

type UserKey string
const UserKeyContext UserKey = "user"

type FollowUser struct {
	UserId int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the user id from the url
	followerUser := getUserFromCtx(r)

	// revert back to auth user
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.store.Users.Follow(r.Context(), followerUser.ID, payload.UserId); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictError(w, r, err)
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}		
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
		// get the user id from the url
	unFollowedUser := getUserFromCtx(r)

	// revert back to auth user
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.store.Users.Unfollow(r.Context(), unFollowedUser.ID, payload.UserId); err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictError(w, r, err)
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			switch err {
			case store.ErrRecordNotFound:
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, UserKeyContext, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) (*store.User) {
	user, _ := r.Context().Value(UserKeyContext).(*store.User)
	return user
}