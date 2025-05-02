package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"thumb/service"
)

func InitRouter(r *gin.Engine) {

	r.GET("/thumb", service.Thumb)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

}
