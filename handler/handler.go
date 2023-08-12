package handler

import (
	respository "UrlShortener/Repository"
	"UrlShortener/shortener"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

type UrlDoc struct {
	UrlCode     string      `bson:"urlCode"`
	LongUrl     string      `bson:"longUrl"`
	ShortUrl    string      `bson:"shortUrl"`
	CreatedAt   time.Time   `bson:"createdAt"`
	ExpiresAt   time.Time   `bson:"expiresAt"`
	AccessedAt  []time.Time `bson:"accessedAt"`
	AccessCount int64       `bson:"accessCount"`
}

var ctx = context.TODO()
var baseUrl = "http://localhost:8080/"

func Shorten(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.BindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, urlErr := url.ParseRequestURI(creationRequest.LongUrl)
	if urlErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": urlErr.Error()})
	}

	shortUrlCode := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)

	result := respository.FindValue(shortUrlCode, "urlcode", c)

	if len(result) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Code in use: %s", shortUrlCode)})
		return
	}

	var newUrl = baseUrl + shortUrlCode
	newDoc := &UrlDoc{
		UrlCode:    shortUrlCode,
		LongUrl:    creationRequest.LongUrl,
		ShortUrl:   newUrl,
		CreatedAt:  time.Now(),
		AccessedAt: []time.Time{time.Now()},
	}
	err := respository.InsertDoc(newDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"newUrl": newUrl,
	})

}
