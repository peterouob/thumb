package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strings"
)

var (
	rdb     *redis.Client
	ErrPipe = errors.New("error in exec use pipe")
)

func NewRdb() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Println("redis init failed ...")
	}

	log.Println("redis init ...")
}

func RunScript(ctx context.Context, counts map[string]int) error {
	pipe := rdb.Pipeline()

	luaScript := `
    return redis.call("HINCRBY", KEYS[1], ARGV[1], ARGV[2])
    `
	for k, v := range counts {
		postId, thumbType, _ := strings.Cut(k, ":")
		pipe.Eval(ctx, luaScript, []string{fmt.Sprintf("count_%s", postId)}, thumbType, v)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return ErrPipe
	}
	return nil
}
