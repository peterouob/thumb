package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
	prom "thumb/pkg/prometheus"
	"thumb/router"
	"thumb/utils"
)

var once sync.Once

func main() {
	go registerMetrics()

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
