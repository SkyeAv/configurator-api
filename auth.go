package main

import (
	"context"
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func getHashHex(apiKey string) string {
	data := []byte(apiKey)
	asInt := xxhash.Sum64(data)

	return fmt.Sprintf("%x", asInt)
}

var ctx = context.Background()

func HypatiaAuth(c *gin.Context, username string, apiKey string) bool {
	opts := &redis.Options{Addr: "localhost:6379", Password: "", DB: 0}
	rdb := redis.NewClient(opts)
	defer rdb.Close()

	hash, err := rdb.Get(ctx, username).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "No registered API Key associated with username."})
		return false
	}

	if hash == getHashHex(apiKey) {
		return true
	}

	c.JSON(501, gin.H{"error": "Username is associated with a different API Key."})
	return false
}
