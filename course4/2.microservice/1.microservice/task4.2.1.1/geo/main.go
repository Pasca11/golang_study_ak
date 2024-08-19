package main

import (
	"geo/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	server := grpc.NewServer()
	service.Register(server)

	l, err := net.Listen("tcp", ":3333")
	if err != nil {
		panic(err)
	}
	go func() {
		err := server.Serve(l)
		if err != nil {
			panic(err)
		}
	}()
	log.Println("Geo gRPC server listening on :3333")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	server.GracefulStop()
	log.Println("Shutting down server...")
}
