package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Get("/api/1", first)
	mux.Get("/api/2", second)
	mux.Get("/api/3", third)

	http.ListenAndServe(":1313", mux)
}

func first(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in one")
	w.Write([]byte("Hello world"))
}

func second(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world2"))
}

func third(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world3"))
}
