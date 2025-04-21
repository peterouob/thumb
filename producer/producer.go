package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"thumb/model"

	"github.com/segmentio/kafka-go"
)

func main() {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("192.168.0.100:9092"),
		Topic:                  "social",
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}

	actions := model.SocialAction{
		PostID: "post_001", Like: 10, Forward: 5, Dislike: 2,
	}

	wg := sync.WaitGroup{}
	wg.Add(10)
	for range make([]struct{}, 10) {
		go func() {
			msg, err := json.Marshal(actions)
			if err != nil {
				log.Fatalf("Failed to marshal action: %v", err)
			}

			err = writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(actions.PostID),
					Value: msg,
				},
			)
			if err != nil {
				log.Fatalf("Failed to write message: %v", err)
			}
			fmt.Printf("Sent: %+v\n", actions)
			wg.Done()
		}()
	}
	wg.Wait()
	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}
}
