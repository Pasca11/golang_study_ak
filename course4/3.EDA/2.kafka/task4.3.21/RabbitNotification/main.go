package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Notification service started")
	msgs, err := ch.Consume(
		"geoQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
	}
}
