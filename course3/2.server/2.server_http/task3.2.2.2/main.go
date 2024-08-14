package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		fmt.Println("port env var not set")
	}
	log.Printf("Сервер запущен на порту %s", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
