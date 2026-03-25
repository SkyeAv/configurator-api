package main

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func getHashHex(apiKey string) string {
	dataSha := []byte(apiKey)
	asIntSha := sha256.Sum256(dataSha)
	asHexSha := fmt.Sprintf("%x", asIntSha)

	dataX := []byte(asHexSha)
	asIntX := xxhash.Sum64(dataX)
	return fmt.Sprintf("%x", asIntX)
}

var ctx = context.Background()

func HypatiaAuth(c *gin.Context, username string, apiKey string) bool {
	if username == "" {
		c.JSON(400, gin.H{"error": "'username' is a required API parameter"})
		return false
	}

	if apiKey == "" {
		c.JSON(400, gin.H{"error": "'api-key' is a required API parameter"})
		return false
	}

	opts := &redis.Options{Addr: "localhost:6379", Password: "", DB: 0}
	rdb := redis.NewClient(opts)
	defer rdb.Close()

	hash, err := rdb.Get(ctx, username).Result()
	if err != nil {
		c.JSON(401, gin.H{"error": "No registered API Key associated with username."})
		return false
	}

	if hash == getHashHex(apiKey) {
		return true
	}

	c.JSON(401, gin.H{"error": "Username is associated with a different API Key."})
	return false
}
