package utils

import "github.com/prometheus/client_golang/prometheus"

var (
	LikeRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "like_system_http_requests_like_total",
			Help: "Total number of HTTP requests to the like service",
		},
		[]string{"path", "method", "status"},
	)
	SaveRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "save_system_http_requests_save_total",
			Help: "Total number of HTTP requests to the save service",
		},
		[]string{"path", "method", "status"},
	)
	DataRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "data_system_http_requests_data_total",
			Help: "Total number of HTTP requests to the data service",
		},
		[]string{"path", "method", "status"},
	)
)
