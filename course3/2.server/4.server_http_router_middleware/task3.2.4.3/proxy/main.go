package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	mux := chi.NewRouter()

	URL, err := url.Parse("http://hugo:1313")
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(URL)

	mux.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			proxy.ServeHTTP(w, r)
		} else {
			w.Write([]byte("Hello from api"))
		}
	})

	http.ListenAndServe(":8080", mux)
}
