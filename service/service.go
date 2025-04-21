package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"thumb/utils"
)

func Like(c *gin.Context) {
	utils.LikeRequestsTotal.WithLabelValues(c.Request.URL.String(),
		strconv.Itoa(http.StatusOK), c.Request.Method).Inc()
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"url":    c.Request.RequestURI,
	})
}

func Save(c *gin.Context) {
	utils.SaveRequestsTotal.WithLabelValues(c.Request.URL.String(),
		strconv.Itoa(http.StatusOK), c.Request.Method).Inc()
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func Alert(c *gin.Context) {}

func Data(c *gin.Context) {}
