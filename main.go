package main

import (
	"UrlShortener/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "Hello World!"})
	})
	r.POST("/shorten", handler.Shorten)
	r.Run(":8000")
}
