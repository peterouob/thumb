package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
	"thumb/pkg/kafka"
	prom "thumb/pkg/prometheus"
	"thumb/router"
	"thumb/utils"
	"time"
)

var once sync.Once

func main() {
	kafka.NewProducer()
	consumer := kafka.NewConsumer("thumb_groub", []string{"thumb"}, 10, 10*time.Minute)
	go func() {
		registerMetrics()
		consumer.StartConsume()
	}()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(utils.Logger())
	router.InitRouter(r)
	r.Run(":8082")
}

func registerMetrics() {
	once.Do(func() {
		prometheus.MustRegister(prom.LikeRequestsTotal, prom.KafkaRequestTotal)
	})
	log.Println("register metrics")
}
