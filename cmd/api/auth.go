package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/saleemlawal/social/internal/store"
)

type registerUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"` // TODO: add password complexity rules
}

type userWithToken struct {
	User  *store.User `json:"user"`
	Token string      `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Register a new user
//	@Description	Registers a new user with the provided username and email
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body		registerUserRequest	true	"User to register"
//	@Success		201		{object}	userWithToken
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input registerUserRequest

	if err := readJSON(w, r, &input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: input.Username,
		Email:    input.Email,
	}

	// hash passsword
	if err := user.Password.Set(input.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// store the user
	plainToken := uuid.New().String()

	// hash the token for storeage but keep the plain token for sending on email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(r.Context(), user, string(hashToken), app.config.mail.exp)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail, store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// send the invitation email

	// rollback if failure

	userWithToken := &userWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
