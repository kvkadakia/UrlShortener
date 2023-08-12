package shortener

import (
	"UrlShortener/generator"
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
	UrlCode   string    `bson:"urlCode"`
	LongUrl   string    `bson:"longUrl"`
	ShortUrl  string    `bson:"shortUrl"`
	CreatedAt time.Time `bson:"createdAt"`
	ExpiresAt time.Time `bson:"expiresAt"`
}

var ctx = context.TODO()
var baseUrl = "http://localhost:8080/"
var urlCollection *mongo.Collection
var accessLogsCollection *mongo.Collection

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
	urlCollection = client.Database("url_shortener").Collection("urls")
	accessLogsCollection = client.Database("url_shortener").Collection("access_timestamps")
	log.Print("Database connected!")
}

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

	shortUrlCode := generator.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)

	var result bson.M
	queryErr := urlCollection.FindOne(ctx, bson.D{{"urlCode", shortUrlCode}}).Decode(&result)
	if queryErr != nil {
		if queryErr != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		}
	}
	if len(result) > 0 {
		calculateAccessDetails(shortUrlCode, c)
		return
	}

	var shortUrl = baseUrl + shortUrlCode
	newDoc := &UrlDoc{
		UrlCode:   shortUrlCode,
		LongUrl:   creationRequest.LongUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now(),
	}
	_, err := urlCollection.InsertOne(ctx, newDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"shortUrl": shortUrl})
}

func Redirect(c *gin.Context) {
	code := c.Param("code")
	var result bson.M
	queryErr := urlCollection.FindOne(ctx, bson.D{{"urlCode", code}}).Decode(&result)
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
	updateUrlAccessDetails(code)
	var longUrl = fmt.Sprint(result["longUrl"])
	c.Redirect(http.StatusPermanentRedirect, longUrl)
}

func calculateAccessDetails(shortUrlCode string, c *gin.Context) {
	currTime := time.Now()
	pastTwentyFourHoursTime := currTime.Add(time.Duration(-24) * time.Hour)
	pastWeekTime := currTime.Add(time.Duration(-24*7) * time.Hour)
	pastWeekAccessLogsFilter := bson.M{
		"accessedAt": bson.M{
			"$gt": pastWeekTime,
			"$lt": currTime,
		},
		"urlCode": shortUrlCode,
	}
	pastTwentyFourHoursAccessLogsFilter := bson.M{
		"accessedAt": bson.M{
			"$gt": pastTwentyFourHoursTime,
			"$lt": currTime,
		},
		"urlCode": shortUrlCode,
	}
	allTimeAccessLogsFilter := bson.M{
		"urlCode": shortUrlCode,
	}
	countPastWeekAccessLogs, _ := accessLogsCollection.CountDocuments(ctx, pastWeekAccessLogsFilter)
	countPastTwentyFourHoursAccessLogs, _ := accessLogsCollection.CountDocuments(ctx, pastTwentyFourHoursAccessLogsFilter)
	countAllTimeAccessLogs, _ := accessLogsCollection.CountDocuments(ctx, allTimeAccessLogsFilter)
	c.JSON(http.StatusOK, gin.H{"info": fmt.Sprintf("Short url already exists: %s | totalAccessCount : %v, pastTwentyFourHoursAccessCount : %v, pastWeekAccessCount : %v", baseUrl+shortUrlCode, countAllTimeAccessLogs, countPastTwentyFourHoursAccessLogs, countPastWeekAccessLogs)})
}

func updateUrlAccessDetails(code string) {
	data := bson.D{{"urlCode", code}, {"accessedAt", time.Now()}}
	_, err := accessLogsCollection.InsertOne(ctx, data)
	if err != nil {
		return
	}
}
