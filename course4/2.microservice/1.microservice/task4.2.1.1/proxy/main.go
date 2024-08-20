package main

import (
	chi "github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"proxy/internal/controller"
	"proxy/internal/service"
)

func main() {
	mux := chi.NewRouter()
	conn, _ := grpc.NewClient(net.JoinHostPort("geo", "3333"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn2, err := grpc.NewClient(net.JoinHostPort("auth", "4444"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("new client error:", err)
	}
	serv := service.NewProxyService(conn, conn2)
	ctrl := controller.NewController(serv)

	mux.Group(func(router chi.Router) {
		router.Use(ctrl.AuthMiddleware)

		router.Post("/api/address/search", ctrl.SearchAnswer)
		router.Post("/aoi/address/answer", ctrl.GeocodeAnswer)
	})

	mux.Post("/login", ctrl.Login)
	mux.Post("/register", ctrl.Register)

	log.Println("Proxy server started on 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
