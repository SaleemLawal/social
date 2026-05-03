package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saleemlawal/social/internal/auth"
	"github.com/saleemlawal/social/internal/store"
	"github.com/saleemlawal/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()
	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStorage()
	mockCache := cache.NewMockCache()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCache,
		authenticator: testAuth,
	}
}

func execRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}
