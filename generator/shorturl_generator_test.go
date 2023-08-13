package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestShortLinkGenerator(t *testing.T) {
	initialLink_1 := "https://www.orkut.com"
	shortLink1 := GenerateShortLink(initialLink_1, UserId)

	initialLink_2 := "https://www.google.com"
	shortLink_2 := GenerateShortLink(initialLink_2, UserId)

	initialLink_3 := "https://www.twitter.com"
	shortLink_3 := GenerateShortLink(initialLink_3, UserId)

	assert.Equal(t, shortLink1, "bEUKZdTZ")
	assert.Equal(t, shortLink_2, "2dDEQAS1")
	assert.Equal(t, shortLink_3, "MQAzw2ab")
}
