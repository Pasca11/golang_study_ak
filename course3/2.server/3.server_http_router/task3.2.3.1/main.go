package main

import (
	chi "github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	mux := chi.NewRouter()

	mux.Get("/1", first)
	mux.Get("/2", second)
	mux.Get("/3", third)

	http.ListenAndServe(":8080", mux)
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
