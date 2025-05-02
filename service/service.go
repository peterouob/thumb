package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"thumb/model"
	"thumb/pkg/kafka"
	"thumb/pkg/prometheus"
)

func Thumb(c *gin.Context) {

	var m model.SocialAction
	if err := c.BindJSON(&m); err != nil {
		prom.LikeRequestsTotal.WithLabelValues(c.Request.URL.String(),
			strconv.Itoa(http.StatusInternalServerError), c.Request.Method).Inc()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("bind json error %v", err),
		})
		return
	}

	pd := kafka.NewProducer()
	pd.Send(m.Topic, m)

	countLike(c, http.StatusOK)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"url":    c.Request.RequestURI,
	})
}

func countLike(c *gin.Context, status int) {
	prom.LikeRequestsTotal.WithLabelValues(c.Request.URL.String(),
		strconv.Itoa(status), c.Request.Method).Inc()
}
