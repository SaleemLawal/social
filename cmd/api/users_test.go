package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	testToken, _ := app.authenticator.GenerateToken(nil)

	t.Run("Should not allow unauthenticated requests", func(t *testing.T) {
		// check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := execRequest(req, mux)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("Should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer " + testToken)

		rr := execRequest(req, mux)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
