package main

import (
	"auth/internal/repository"
	"auth/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	grpcServer := grpc.NewServer()
	//service.RegServer(grpcServer)

	rep, err := repository.NewStorage()
	if err != nil {
		panic(err)
	}
	service.NewAuthService(rep, []byte("secret"), grpcServer)

	l, err := net.Listen("tcp", ":4444")
	if err != nil {
		panic(err)
	}
	go func() {
		err := grpcServer.Serve(l)
		if err != nil {
			panic(err)
		}
	}()
	log.Println("Auth gRPC server listening on :4444")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	grpcServer.GracefulStop()
	log.Println("Shutting down server...")
}
