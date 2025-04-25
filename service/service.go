package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"net/http"
	"strconv"
	"thumb/model"
	"thumb/utils"
)

func init() {
	if utils.Writer == nil || utils.Reader == nil {
		utils.NewReader("social")
		utils.NewWriter("social")
	}
}

func Like(c *gin.Context) {

	var m model.SocialAction
	if err := c.BindJSON(&m); err != nil {
		utils.LikeRequestsTotal.WithLabelValues(c.Request.URL.String(),
			strconv.Itoa(http.StatusInternalServerError), c.Request.Method).Inc()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("bind json error %v", err),
		})
		return
	}

	msg, err := json.Marshal(m)
	if err != nil {
		countLike(c, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("json marsharl error %v", err),
		})
		return
	}

	if err := utils.Writer.WriteMessages(c, kafka.Message{
		Key:   []byte(m.PostID),
		Value: msg,
	}); err != nil {
		countLike(c, http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("json marsharl error %v", err),
		})
		return
	}
	countLike(c, http.StatusOK)
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

func countLike(c *gin.Context, status int) {
	utils.LikeRequestsTotal.WithLabelValues(c.Request.URL.String(),
		strconv.Itoa(status), c.Request.Method).Inc()
}
