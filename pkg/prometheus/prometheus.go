package prom

import "github.com/prometheus/client_golang/prometheus"

var (
	LikeRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "like_system_http_requests_like_total",
			Help: "Total number of HTTP requests to the like service",
		},
		[]string{"path", "method", "status"},
	)
	KafkaRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "kafka_request_total",
		Help: "Total number of requests to the kafka service",
	}, []string{"topic", "status", "msg"})
)
