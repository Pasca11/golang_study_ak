package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func prepareRouter() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/1", first)
	mux.Get("/2", second)
	mux.Get("/3", third)

	return mux
}

func TestRouter(t *testing.T) {
	router := prepareRouter()

	t.Run("Test /1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Invalid response code")
		assert.Equal(t, rec.Body.String(), `Hello world`, "Invalid response body")
	})
	t.Run("Test /2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/2", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Invalid response code")
		assert.Equal(t, rec.Body.String(), `Hello world2`, "Invalid response body")
	})
	t.Run("Test /3", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/3", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Invalid response code")
		assert.Equal(t, rec.Body.String(), `Hello world3`, "Invalid response body")
	})
}
