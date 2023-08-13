package shortener

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"net/http/httptest"
	"testing"
)

func init() {
	InitializeDb()
}

func TestInsertionRetrievalAndDeletion(t *testing.T) {
	initialLink := "https://www.google.com"
	shortURL := "http://localhost:8080/GMH9monD"
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	SaveUrlMapping("GMH9monD", initialLink, shortURL, ctx)
	var result bson.M
	result = RetrieveInitialUrl(ctx, "GMH9monD")
	DeleteShortUrl("GMH9monD")
	assert.Equal(t, shortURL, result["shortUrl"])
}
