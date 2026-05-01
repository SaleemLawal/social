package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/saleemlawal/social/internal/mailer"
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

type createUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
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
	isProduction := app.config.env == "production"
	activationURL := fmt.Sprintf("%s/authentication/activate?token=%s", app.config.frontendUrl, plainToken)

	email := &mailer.Email{
		Username:      user.Username,
		ToEmail:       user.Email,
		ActivationURL: activationURL,
	}

	err = app.mailer.Send(mailer.UserInvitationTemplate, email, isProduction)

	if err != nil {
		app.logger.Errorw("Failed to send invitation email", "error", err.Error())

		// rollback user creation if email sending fails
		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("Failed to rollback user creation", "error", err.Error())
			app.internalServerError(w, r, err)
			return
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent successfully")

	userWithToken := &userWithToken{
		User:  user,
		Token: plainToken,
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// createTokenHandler godoc
//
//	@Summary		Create a new token
//	@Description	Creates a new token for the user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserTokenPayload	true	"User to register"
//	@Success		201		{object}	string					"token"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input createUserTokenPayload

	if err := readJSON(w, r, &input); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	// fetch the user (check if the user exist) from the payload
	user, err := app.store.Users.GetByEmail(r.Context(), input.Email)
	if err != nil {
		switch err {
		case store.ErrRecordNotFound:
			app.unauthorizedError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(input.Password); err != nil {
		app.unauthorizedError(w, r, fmt.Errorf("Invalid credentials"))
		return
	}

	// generate a new token -> add claims to the token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"aud": app.config.auth.token.audience,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(jwt.Claims(claims))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	// return the token
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
