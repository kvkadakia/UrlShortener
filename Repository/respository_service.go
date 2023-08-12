package respository

import (
	"UrlShortener/handler"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

var ctx = context.TODO()
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

func FindValue(shortUrlCode string, field string, c *gin.Context) bson.M {
	var result bson.M
	queryErr := collection.FindOne(ctx, bson.D{{field, shortUrlCode}}).Decode(&result)

	if queryErr != nil {
		if queryErr != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		}
	}
	return result
}

func InsertDoc(newDoc *handler.UrlDoc) error {
	_, err := collection.InsertOne(ctx, newDoc)
	return err
}
