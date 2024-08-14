package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	mux := chi.NewRouter()
	mux.Route("/group1", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			id := chi.URLParam(r, "id")
			w.Write([]byte("Group 1 Hello World" + id))
		})
	})
	mux.Route("/group2", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			id := chi.URLParam(r, "id")
			w.Write([]byte("Group 2 Hello World" + id))
		})
	})
	mux.Route("/group3", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			id := chi.URLParam(r, "id")
			w.Write([]byte("Group 3 Hello World" + id))
		})
	})
	http.ListenAndServe(":8080", mux)
}
