package main

import (
	"context"
	"errors"
	"fmt"
	"geoservice-rabbit/internal/controller"
	"geoservice-rabbit/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
)

// @title GeoService
// @version 1.0
// @description Simple GeoService.

// @host localhost:8080
// @BasePath /

const (
	adress = ":8080"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	r := makeRouter()
	metrics.RegisterAll()
	server := &http.Server{
		Addr:         adress,
		Handler:      r,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	go func() {
		log.Println("Starting server on", adress)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	sig := <-sigChan
	fmt.Printf("Recieved signal: %v. Starting shutting down\n", sig)

	shuttingDownTime := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shuttingDownTime)
	defer cancel()

	err := server.Shutdown(ctx)

	if err == nil {
		log.Println("Server stopped gracefully")
	}
}

func makeRouter() *chi.Mux {
	r := chi.NewRouter()
	tokenAuth = jwtauth.New("HS256", []byte("mysecretkey"), nil)
	r.Use(middleware.Logger)
	responder, err := controller.NewController(tokenAuth)
	if err != nil {
		log.Fatal(err)
	}

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/api/profile", pprof.Index)
		r.Post("/api/address/search", responder.SearchAnswer)
		r.Post("/api/address/geocode", responder.GeocodeAnswer)
	})

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/api/users/{id}", responder.GetUserByID)
	r.Post("/api/register", responder.RegisterUser)
	r.Post("/api/login", responder.LoginUser)
	r.NotFound(responder.NotFoundAnswer)
	return r
}
