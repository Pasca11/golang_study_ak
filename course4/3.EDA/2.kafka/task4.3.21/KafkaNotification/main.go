package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	cfg := kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "my_topic",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
		MaxWait:   1 * time.Second,
	}

	reader := kafka.NewReader(cfg)
	defer reader.Close()
	defer log.Println("Kafka consumer stopped")

	log.Println("Kafka listener started")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("reader error:", err)
		} else {
			log.Println("reader message:", string(m.Value))
		}
	}
}
