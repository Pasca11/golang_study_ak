package main

import (
	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	mux := chi.NewRouter()

	mux.Use(LoggerMiddleware)

	mux.Get("/1", first)
	mux.Get("/2", second)
	mux.Get("/3", third)

	http.ListenAndServe(":8080", mux)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("Request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
		next.ServeHTTP(w, r)
	})
}

func first(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func second(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world2"))
}

func third(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world3"))
}
