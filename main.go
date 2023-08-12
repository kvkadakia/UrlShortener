package main

import (
	"UrlShortener/shortener"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "Welcome to url generator!"})
	})
	r.POST("/shorten", func(c *gin.Context) {
		shortener.Shorten(c)
	})
	r.GET("/:code", func(c *gin.Context) {
		shortener.Redirect(c)
	})
	r.Run(":8080")
}
