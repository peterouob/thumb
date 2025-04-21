package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"thumb/model"
	"time"

	"github.com/segmentio/kafka-go"
)

const batchSize = 100

var (
	batch       = make([]*kafka.Message, 0, batchSize)
	accumulator = make(map[string][3]int)
	ticker      = time.NewTicker(1 * time.Minute)
	ctx         = context.Background()
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"192.168.0.100:9092"},
		Topic:    "social",
		GroupID:  "social-action-consumer-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	defer ticker.Stop()

	flush := func() {
		if len(accumulator) == 0 {
			return
		}

		pipe := rdb.Pipeline()
		luaScript := `
    redis.call("HINCRBY", KEYS[1], "like", ARGV[1])
    redis.call("HINCRBY", KEYS[1], "forward", ARGV[2])
    redis.call("HINCRBY", KEYS[1], "dislike", ARGV[3])
    return 0
    `
		for postID, counts := range accumulator {
			key := fmt.Sprintf("post:%s", postID)
			pipe.Eval(ctx, luaScript, []string{key}, counts[0], counts[1], counts[2])
		}

		if _, err := pipe.Exec(ctx); err != nil {
			log.Printf("Redis pipeline error: %v", err)
		}

		batch = batch[:0]
		accumulator = make(map[string][3]int)
	}

	for {
		select {
		case <-ticker.C:
			flush()

		default:
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Failed to read message: %v", err)
				continue
			}

			batch = append(batch, &msg)

			var action model.SocialAction
			if err := json.Unmarshal(msg.Value, &action); err != nil {
				log.Printf("Failed to unmarshal: %v", err)
				continue
			}

			key := action.PostID
			acc := accumulator[key]
			acc[0] += action.Like
			acc[1] += action.Forward
			acc[2] += action.Dislike
			accumulator[key] = acc

			if len(batch) >= batchSize {
				flush()
			}
		}
	}

}
