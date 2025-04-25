package utils

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

var (
	Writer *kafka.Writer
	Reader *kafka.Reader
)

func NewWriter(topic string) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("192.168.0.100:9092"),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
	log.Println("writer init done ...")
	Writer = writer
}

func NewReader(topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"192.168.0.100:9092"},
		Topic:    topic,
		GroupID:  fmt.Sprintf("%s_group", topic),
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	log.Println("reader init done ...")
	Reader = reader
}
