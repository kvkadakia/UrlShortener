package main

import (
	"UrlShortener/shortener"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("url_shortener").Collection("urls")
	log.Print("Database connected!")
}

func shorten(c *gin.Context) {
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

	var result bson.M
	queryErr := collection.FindOne(ctx, bson.D{{"urlCode", shortUrlCode}}).Decode(&result)
	if queryErr != nil {
		if queryErr != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		}
	}
	if len(result) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Short url already exists: %s", baseUrl+shortUrlCode)})
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
	_, err := collection.InsertOne(ctx, newDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"newUrl": newUrl})
}

func redirect(c *gin.Context) {
	code := c.Param("code")
	var result bson.M
	queryErr := collection.FindOne(ctx, bson.D{{"urlCode", code}}).Decode(&result)

	if queryErr != nil {
		if queryErr == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No URL with code: %s", code)})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
			return
		}
	}
	log.Print(result["longUrl"])
	updateUrlAccessDetails(code, result)
	var longUrl = fmt.Sprint(result["longUrl"])
	c.Redirect(http.StatusPermanentRedirect, longUrl)
}

func updateUrlAccessDetails(code string, result bson.M) {
	filter := bson.D{{"urlCode", code}}
	accessCount := result["accessCount"].(int64)
	accessCount += 1
	update := bson.D{{"$set", bson.D{{"accessCount", accessCount}}},
		{"$push", bson.D{{"accessedAt", time.Now()}}}}
	collection.UpdateOne(ctx, filter, update)
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"message": "Welcome to url shortener!"})
	})
	r.POST("/shorten", shorten)
	r.GET("/:code", redirect)
	r.Run(":8080")
}
