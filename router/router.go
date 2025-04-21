package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"thumb/service"
)

func InitRouter(r *gin.Engine) {

	r.GET("/like", service.Like)
	r.GET("/save", service.Save)
	r.GET("/data", service.Data)
	r.GET("/alert", service.Alert)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

}
