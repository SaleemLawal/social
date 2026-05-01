package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/saleemlawal/social/internal/store"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				app.unauthorizedBasicAuthError(w, r, fmt.Errorf("No authorization header"))
				return
			}

			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicAuthError(w, r, fmt.Errorf("Invalid authorization header"))
				return
			}

			//decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicAuthError(w, r, err)
				return
			}

			username := app.config.auth.basic.username
			password := app.config.auth.basic.password
			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 || credentials[0] != username || credentials[1] != password {
				app.unauthorizedBasicAuthError(w, r, fmt.Errorf("Invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("No authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("Invalid authorization header"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		sub, ok := claims["sub"].(float64)
		if !ok {
			app.unauthorizedError(w, r, fmt.Errorf("invalid subject claim"))
			return
		}
		userId := int64(sub)

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, UserKeyContext, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkPostOwnershipMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post, ok := getPostFromCtx(r)
		if !ok {
			app.unauthorizedError(w, r, fmt.Errorf("Post not found"))
			return
		}

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
		}

		if !allowed {
			app.forbiddenResponse(w, r, err)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}
