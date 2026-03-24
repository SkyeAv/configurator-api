package main

import (
	"context"
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/redis/go-redis/v9"
)

func getHashHex(apiKey string) string {
	data := []byte(apiKey)
	asInt := xxhash.Sum64(data)

	return fmt.Sprintf("%x", asInt)
}

var ctx = context.Background()

func auth(username string, apiKey string) bool {
	opts := &redis.Options{Addr: "localhost:6379", Password: "", DB: 0}
	rdb := redis.NewClient(opts)
	defer rdb.Close()

	hash, err := rdb.Get(ctx, username).Result()
	CheckErr(1, err)

	if hash == getHashHex(apiKey) {
		return true
	}

	return false
}
