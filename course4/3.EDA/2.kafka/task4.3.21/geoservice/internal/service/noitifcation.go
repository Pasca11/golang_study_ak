package service

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Notificator interface {
	Notify(message string) error
}

type RabbitService struct {
	channel *amqp.Channel
	queue   *amqp.Queue
}

type KafkaService struct {
	writer *kafka.Writer
}

func NewNotificationService() (Notificator, error) {
	godotenv.Load()

	broker := os.Getenv("BROKER")
	if broker == "rabbit" {
		return createRabbit()
	}
	if broker == "kafka" {
		return createKafka()
	}
	return nil, fmt.Errorf("unknown broker type: %s", broker)
}

func (n *RabbitService) Notify(message string) error {
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

func (s *KafkaService) Notify(message string) error {
	err := s.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Value: []byte(message),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func createRabbit() (Notificator, error) {
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
	return &RabbitService{channel: ch, queue: &q}, nil
}

func createKafka() (Notificator, error) {
	writer := &kafka.Writer{
		Topic:    "my_topic",
		Balancer: &kafka.LeastBytes{},
		Addr:     kafka.TCP("kafka:9092"),
	}
	stats := writer.Stats()
	log.Printf("Kafka stats: %+v", stats)
	return &KafkaService{writer: writer}, nil
}
