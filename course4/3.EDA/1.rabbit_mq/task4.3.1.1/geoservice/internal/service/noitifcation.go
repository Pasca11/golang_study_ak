package service

import (
	"github.com/streadway/amqp"
	"log"
)

type Notificator interface {
	Notify(message string) error
}

type NotificationService struct {
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewNotificationService() (Notificator, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		"geoQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &NotificationService{channel: ch, queue: &q}, nil
}

func (n *NotificationService) Notify(message string) error {
	err := n.channel.Publish(
		"",
		n.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Message %s was sent successfully", message)
	return nil
}
